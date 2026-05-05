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
