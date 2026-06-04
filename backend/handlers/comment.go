package handlers

import (
	"net/http"
	"strconv"

	"github.com/enzo959/forumium/models"
	"github.com/enzo959/forumium/repositories"
	"github.com/gin-gonic/gin"
)

func GetCommentsByPostID(c *gin.Context) {
	rawID := c.Param("id")

	parsedID, err := strconv.ParseUint(rawID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	postID := uint(parsedID)

	comments, err := repositories.GetCommentsByPostID(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	type CommentResponse struct {
		ID        uint   `json:"id"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
		Author    struct {
			ID       uint   `json:"id"`
			Username string `json:"username"`
			Avatar   string `json:"avatar"`
		} `json:"author"`
	}

	result := []CommentResponse{}

	for _, comment := range comments {
		cr := CommentResponse{
			ID:        comment.ID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt.String(),
		}
		cr.Author.ID = comment.User.ID
		cr.Author.Username = comment.User.UserName
		cr.Author.Avatar = comment.User.Avatar
		result = append(result, cr)
	}

	c.JSON(http.StatusOK, gin.H{"comments": result})
}

func CreateComment(c *gin.Context) {
	// Récupère l'ID du post depuis l'URL
	rawPostID := c.Param("id")
	parsedPostID, err := strconv.ParseUint(rawPostID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	postID := uint(parsedPostID)

	// Récupère l'ID de l'utilisateur connecté depuis le middleware
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDStr, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}
	parsedUserID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}
	userID := uint(parsedUserID)

	// Valide le body
	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content is required"})
		return
	}

	comment := models.Comment{
		UserID:  userID,
		PostID:  postID,
		Content: req.Content,
	}

	if err := repositories.CreateComment(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         comment.ID,
		"content":    comment.Content,
		"created_at": comment.CreatedAt,
	})
}