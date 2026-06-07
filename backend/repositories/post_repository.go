package repositories

import (
	"github.com/enzo959/forumium/config"
	"github.com/enzo959/forumium/models"
)

func GetPosts(page int, limit int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	offset := (page - 1) * limit
	config.DB.Model(&models.Post{}).Count(&total)

	err := config.DB.
		Preload("User").
		Preload("Categories").
		Preload("Reactions").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error

	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}