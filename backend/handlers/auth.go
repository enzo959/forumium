package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/enzo959/forumium/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/enzo959/forumium/config"
)

type RegisterInput struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "JSON invalide",
		})
		return
	}

	if input.Email == "" || input.Username == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email, username et password requis",
		})
		return
	}

	if len(input.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "mot de passe trop court (min 8 caractères)",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erreur lors du hash du mot de passe",
		})
		return
	}

	var existingUser models.User
	result := config.DB.Where("email = ?", input.Email).First(&existingUser)

	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "email déjà utilisé",
		})
		return
	}

	user := models.User{
		UserName: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		Avatar:   "",
	}

	result = config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erreur création utilisateur",
		})
		return
	}

	userID := strconv.Itoa(int(user.ID))

	accessToken, err := config.GenerateAccessToken(userID, input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erreur génération access token",
		})
		return
	}

	refreshToken, err := config.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erreur génération refresh token",
		})
		return
	}

	hashedRefreshToken, err := bcrypt.GenerateFromPassword(
		[]byte(refreshToken),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erreur hash refresh token",
		})
		return
	}

	tokenRecord := models.RefreshToken{
		UserID:    user.ID,
		Token:     string(hashedRefreshToken),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
	}

	result = config.DB.Create(&tokenRecord)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erreur sauvegarde refresh token",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "utilisateur créé",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "JSON invalide",
		})
		return
	}

	if input.Email == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email et password requis",
		})
		return
	}

	var user models.User
	result := config.DB.Where("email = ?", input.Email).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "identifiants invalides",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(input.Password),
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "identifiants invalides",
		})
		return
	}

	userID := strconv.Itoa(int(user.ID))

	accessToken, err := config.GenerateAccessToken(userID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erreur génération access token",
		})
		return
	}

	refreshToken, err := config.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erreur génération refresh token",
		})
		return
	}

	hashedRefreshToken, err := bcrypt.GenerateFromPassword(
		[]byte(refreshToken),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erreur hash refresh token",
		})
		return
	}

	result = config.DB.Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked = ?", user.ID, false).
		Update("revoked", true)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erreur révocation anciens tokens",
		})
		return
	}

	tokenRecord := models.RefreshToken{
		UserID:    user.ID,
		Token:     string(hashedRefreshToken),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
	}

	result = config.DB.Create(&tokenRecord)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erreur sauvegarde refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":       user.ID,
			"username": user.UserName,
			"email":    user.Email,
		},
	})
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token"`
	UserID       uint   `json:"user_id"`
}

func Refresh(c *gin.Context) {
	var input RefreshInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "JSON invalide",
		})
		return
	}

	if input.RefreshToken == "" || input.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "refresh_token et user_id requis",
		})
		return
	}

	var tokens []models.RefreshToken

	result := config.DB.Where(
		"user_id = ? AND revoked = ? AND expires_at > ?",
		input.UserID, false, time.Now(),
	).Find(&tokens)

	if result.Error != nil || len(tokens) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "token invalide",
		})
		return
	}

	var validToken *models.RefreshToken

	for _, t := range tokens {
		err := bcrypt.CompareHashAndPassword(
			[]byte(t.Token),
			[]byte(input.RefreshToken),
		)

		if err == nil {
			validToken = &t
			break
		}
	}

	if validToken == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "token invalide",
		})
		return
	}

	if validToken.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "token expiré",
		})
		return
	}

	result = config.DB.Model(validToken).Update("revoked", true)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erreur révocation token",
		})
		return
	}

	var user models.User
	result = config.DB.First(&user, validToken.UserID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "utilisateur introuvable",
		})
		return
	}

	userID := strconv.Itoa(int(user.ID))

	accessToken, err := config.GenerateAccessToken(userID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erreur génération access token",
		})
		return
	}

	refreshToken, err := config.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erreur génération refresh token",
		})
		return
	}

	hashedRefreshToken, err := bcrypt.GenerateFromPassword(
		[]byte(refreshToken),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erreur hash refresh token",
		})
		return
	}

	newToken := models.RefreshToken{
		UserID:    user.ID,
		Token:     string(hashedRefreshToken),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
	}

	result = config.DB.Create(&newToken)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erreur sauvegarde refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// ================================
// 🚪 LOGOUT
// ================================

type LogoutInput struct {
	RefreshToken string `json:"refresh_token"`
	UserID       uint   `json:"user_id"`
}

func Logout(c *gin.Context) {
	var input LogoutInput

	// ÉTAPE 1 - Parser JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "JSON invalide",
		})
		return
	}

	// ÉTAPE 2 - Validation
	if input.RefreshToken == "" || input.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "refresh_token et user_id requis",
		})
		return
	}

	// ÉTAPE 3 - Récupérer tokens valides
	var tokens []models.RefreshToken

	result := config.DB.Where(
		"user_id = ? AND revoked = ? AND expires_at > ?",
		input.UserID, false, time.Now(),
	).Find(&tokens)

	if result.Error != nil || len(tokens) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "token invalide",
		})
		return
	}

	// ÉTAPE 4 - Trouver le bon token (bcrypt)
	var validToken *models.RefreshToken

	for _, t := range tokens {
		err := bcrypt.CompareHashAndPassword(
			[]byte(t.Token),
			[]byte(input.RefreshToken),
		)

		if err == nil {
			validToken = &t
			break
		}
	}

	if validToken == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "token invalide",
		})
		return
	}

	// ÉTAPE 5 - Révoquer le token
	result = config.DB.Model(validToken).Update("revoked", true)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erreur révocation token",
		})
		return
	}

	// ÉTAPE 6 - Réponse
	c.JSON(http.StatusOK, gin.H{
		"message": "déconnecté avec succès",
	})
}
