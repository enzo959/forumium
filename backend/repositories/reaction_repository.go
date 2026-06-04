package repositories

import (
	"errors"

	"github.com/enzo959/forumium/config"
	"github.com/enzo959/forumium/models"
	"gorm.io/gorm"
)

func GetReactionByUserAndPost(userID uint, postID uint) (*models.Reaction, error) {
	var reaction models.Reaction
	err := config.DB.Unscoped().Where("user_id = ? AND post_id = ?", userID, postID).First(&reaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &reaction, nil
}

func CreateReaction(reaction *models.Reaction) error {
	return config.DB.Create(reaction).Error
}

func UpdateReaction(reaction *models.Reaction, newType models.ReactionType) error {
	return config.DB.Model(reaction).Update("type", newType).Error
}

func DeleteReaction(reaction *models.Reaction) error {
	return config.DB.Unscoped().Where("id = ?", reaction.ID).Delete(&models.Reaction{}).Error
}

func GetReactionCountsByPostID(postID uint) (likes int64, dislikes int64, err error) {
	config.DB.Model(&models.Reaction{}).Where("post_id = ? AND type = ?", postID, models.Like).Count(&likes)
	config.DB.Model(&models.Reaction{}).Where("post_id = ? AND type = ?", postID, models.Dislike).Count(&dislikes)
	return likes, dislikes, nil
}