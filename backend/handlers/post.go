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
		Likes:      likes,
		Dislikes:   dislikes,
	}

	c.JSON(http.StatusOK, response)
}

func CreatePost(c *gin.Context) {

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	userIDStr, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	parsedID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	userID := uint(parsedID)

	var req dtos.CreatePostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Title is required",
		})
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Content is required",
		})
		return
	}

	if len(req.CategoryIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "At least one category is required",
		})
		return
	}

	var categories []models.Category

	result := config.DB.
		Where("id IN ?", req.CategoryIDs).
		Find(&categories)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch categories",
		})
		return
	}

	if len(categories) != len(req.CategoryIDs) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "One or more categories not found",
		})
		return
	}

	post := models.Post{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
		Image:   req.Image,
	}

	result = config.DB.Create(&post)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create post",
		})
		return
	}

	err = config.DB.
		Model(&post).
		Association("Categories").
		Replace(categories)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to associate categories",
		})
		return
	}

	err = config.DB.
		Preload("User").
		Preload("Categories").
		First(&post, post.ID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Post not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load post",
		})
		return
	}

	categories2 := []dtos.CategoryResponse{}

	for _, category := range post.Categories {
		categories2 = append(categories2, dtos.CategoryResponse{
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
		Categories: categories2,
		Likes:      0,
		Dislikes:   0,
	}

	c.JSON(http.StatusCreated, response)
}

// À ajouter à la fin de post.go (en dehors de CreatePost !)

func UpdatePost(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	parsedUserID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}
	userID := uint(parsedUserID)


	postIDRaw := c.Param("id")
	parsedPostID, err := strconv.ParseUint(postIDRaw, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	postID := uint(parsedPostID)

	var post models.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this post"})
		return
	}

	var req dtos.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	post.Image = req.Image 

	if err := config.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	if len(req.CategoryIDs) > 0 {
		var categories []models.Category
		result := config.DB.Where("id IN ?", req.CategoryIDs).Find(&categories)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
			return
		}

		if len(categories) != len(req.CategoryIDs) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "One or more categories not found"})
			return
		}

		if err := config.DB.Model(&post).Association("Categories").Replace(categories); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update categories association"})
			return
		}
	}

	if err := config.DB.Preload("User").Preload("Categories").Preload("Reactions").First(&post, post.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload updated post"})
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

	categoriesResponse := []dtos.CategoryResponse{}
	for _, category := range post.Categories {
		categoriesResponse = append(categoriesResponse, dtos.CategoryResponse{
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
		Categories: categoriesResponse,
		Likes:      likes,
		Dislikes:   dislikes,
	}

	c.JSON(http.StatusOK, response)
}

func DeletePost(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	parsedUserID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}
	userID := uint(parsedUserID)

	postIDRaw := c.Param("id")
	parsedPostID, err := strconv.ParseUint(postIDRaw, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	postID := uint(parsedPostID)

	var post models.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete this post"})
		return
	}

	if err := config.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post successfully deleted",
	})
}