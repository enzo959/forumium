package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/enzo959/forumium/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUserProfile(c *gin.Context) {
	rawID := c.Param("id")
	parsedID, err := strconv.ParseUint(rawID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	userID := uint(parsedID)

	user, err := repositories.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	type PostSummary struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		CreatedAt string `json:"created_at"`
	}
	posts := []PostSummary{}
	for _, p := range user.Posts {
		posts = append(posts, PostSummary{
			ID:        p.ID,
			Title:     p.Title,
			CreatedAt: p.CreatedAt.String(),
		})
	}

	type CommentSummary struct {
		ID        uint   `json:"id"`
		Content   string `json:"content"`
		PostID    uint   `json:"post_id"`
		CreatedAt string `json:"created_at"`
	}
	comments := []CommentSummary{}
	for _, cm := range user.Comments {
		comments = append(comments, CommentSummary{
			ID:        cm.ID,
			Content:   cm.Content,
			PostID:    cm.PostID,
			CreatedAt: cm.CreatedAt.String(),
		})
	}

	type ReactionSummary struct {
		PostID uint   `json:"post_id"`
		Type   string `json:"type"`
	}
	reactions := []ReactionSummary{}
	for _, r := range user.Reactions {
		reactions = append(reactions, ReactionSummary{
			PostID: r.PostID,
			Type:   string(r.Type),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"username":   user.UserName,
		"avatar":     user.Avatar,
		"created_at": user.CreatedAt,
		"posts":      posts,
		"comments":   comments,
		"reactions":  reactions,
	})
}