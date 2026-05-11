package handlers

import (
	"net/http"
	"strconv"

	"github.com/enzo959/forumium/repositories"
	"github.com/gin-gonic/gin"
)

func GetPosts(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 50 {
		limit = 10
	}

	posts, total, err := repositories.GetPosts(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur lors de la récupération des posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}