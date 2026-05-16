package handlers

import (
	"net/http"

	"github.com/enzo959/forumium/repositories"
	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	categories, err := repositories.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur lors de la récupération des catégories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}