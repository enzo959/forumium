package handlers

import (
	"net/http"

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
	//hashedPassword
	_, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erreur lors du hash du mot de passe",
		})
		return
	}

	// ÉTAPE 4 - Créer l'user en base
	// TODO: sauvegarder l'utilisateur en base
	// ex: user := models.User{Email: input.Email, Password: string(hashedPassword)}
	// db.Create(&user)

	// FAKE userID pour l'instant
	userID := "1"

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

	c.JSON(http.StatusCreated, gin.H{
		"message":       "utilisateur créé",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
