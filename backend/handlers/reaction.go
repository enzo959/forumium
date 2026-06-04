package handlers

import (
	"net/http"
	"strconv"

	"github.com/enzo959/forumium/models"
	"github.com/enzo959/forumium/repositories"
	"github.com/gin-gonic/gin"
)

func ReactToPost(c *gin.Context) {
	rawPostID := c.Param("id")
	parsedPostID, err := strconv.ParseUint(rawPostID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	postID := uint(parsedPostID)

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

	var req struct {
		Type string `json:"type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if req.Type != "like" && req.Type != "dislike" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Type must be 'like' or 'dislike'"})
		return
	}
	newType := models.ReactionType(req.Type)

	existing, err := repositories.GetReactionByUserAndPost(userID, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if existing != nil {
		if existing.Type == newType {
			// Même type → annuler la réaction
			if err := repositories.DeleteReaction(existing); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete reaction"})
				return
			}
		} else {
			// Type différent → mettre à jour
			if err := repositories.UpdateReaction(existing, newType); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reaction"})
				return
			}
		}
	} else {
		reaction := models.Reaction{
			UserID: userID,
			PostID: postID,
			Type:   newType,
		}
		if err := repositories.CreateReaction(&reaction); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reaction"})
			return
		}
	}

	likes, dislikes, err := repositories.GetReactionCountsByPostID(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reaction counts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"likes":    likes,
		"dislikes": dislikes,
	})
}