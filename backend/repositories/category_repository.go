package repositories

import (
	"github.com/enzo959/forumium/config"
	"github.com/enzo959/forumium/models"
)

func GetAllCategories() ([]models.Category, error) {
	var categories []models.Category

	err := config.DB.Order("name ASC").Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}