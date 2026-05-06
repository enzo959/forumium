package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/enzo959/forumium/config"
	"github.com/enzo959/forumium/models"
	"github.com/enzo959/forumium/services"
)

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email invalide"})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Si cet email existe, un lien de reset a été envoyé"})
		return
	}

	if err := services.GeneratePasswordReset(user.ID, user.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'envoi de l'email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Si cet email existe, un lien de reset a été envoyé"})
}