package handlers

import (
	"net/http"
	"time"

	"github.com/enzo959/forumium/config"
	"github.com/enzo959/forumium/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type ResetPasswordInput struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func ResetPassword(c *gin.Context) {
	var input ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token et password requis"})
		return
	}

	if len(input.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "mot de passe trop court (min 8 caractères)"})
		return
	}

	var reset models.PasswordReset
	result := config.DB.Where(
		"token = ? AND used = ? AND expires_at > ?",
		input.Token, false, time.Now(),
	).First(&reset)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token invalide ou expiré"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur lors du hash du mot de passe"})
		return
	}

	result = config.DB.Model(&models.User{}).
		Where("id = ?", reset.UserID).
		Update("password", string(hashedPassword))

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur lors de la mise à jour du mot de passe"})
		return
	}

	config.DB.Model(&reset).Update("used", true)

	c.JSON(http.StatusOK, gin.H{"message": "mot de passe mis à jour avec succès"})
}