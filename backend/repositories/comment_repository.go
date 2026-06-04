package repositories

import (
	"github.com/enzo959/forumium/config"
	"github.com/enzo959/forumium/models"
)

func GetCommentsByPostID(postID uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := config.DB.Preload("User").Where("post_id = ?", postID).Order("created_at ASC").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func CreateComment(comment *models.Comment) error {
	return config.DB.Create(comment).Error
}

func GetCommentByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := config.DB.First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func DeleteComment(comment *models.Comment) error {
	return config.DB.Delete(comment).Error
}