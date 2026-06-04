package handlers

import (
	"net/http"
	"strconv"

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