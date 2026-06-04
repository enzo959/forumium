package repositories

import (
	"github.com/enzo959/forumium/config"
	"github.com/enzo959/forumium/models"
)

func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := config.DB.
		Preload("Posts").
		Preload("Comments").
		Preload("Reactions").
		First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}