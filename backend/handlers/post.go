package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/enzo959/forumium/config"
	"github.com/enzo959/forumium/dtos"
	"github.com/enzo959/forumium/models"
	"github.com/enzo959/forumium/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func GetPostByID(c *gin.Context) {

	rawID := c.Param("id")

	parsedID, err := strconv.ParseUint(rawID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	id := uint(parsedID)

	var post models.Post

	result := config.DB.
		Preload("User").
		Preload("Categories").
		Preload("Reactions").
		First(&post, id)

	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Post not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	likes := 0
	dislikes := 0

	for _, reaction := range post.Reactions {

		if reaction.Type == "like" {
			likes++
		} else if reaction.Type == "dislike" {
			dislikes++
		}
	}

	categories := []dtos.CategoryResponse{}

	for _, category := range post.Categories {
		categories = append(categories, dtos.CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
		})
	}

	response := dtos.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		Image:     post.Image,
		CreatedAt: post.CreatedAt,

		Author: dtos.AuthorResponse{
			ID:       post.User.ID,
			Username: post.User.UserName,
			Avatar:   post.User.Avatar,
		},

		Categories: categories,

		Likes:    likes,
		Dislikes: dislikes,
	}

	c.JSON(http.StatusOK, response)
}
