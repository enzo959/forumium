package repositories

import (
	"github.com/enzo959/forumium/config"
	"github.com/enzo959/forumium/models"
)

func CreatePasswordReset(reset *models.PasswordReset) error {
	return config.DB.Create(reset).Error
}

func FindPasswordResetByToken(token string) (*models.PasswordReset, error) {
	var reset models.PasswordReset
	err := config.DB.Where("token = ? AND used = ? AND expires_at > NOW()", token, false).First(&reset).Error
	if err != nil {
		return nil, err
	}
	return &reset, nil
}

func MarkPasswordResetAsUsed(id uint) error {
	return config.DB.Model(&models.PasswordReset{}).Where("id = ?", id).Update("used", true).Error
}