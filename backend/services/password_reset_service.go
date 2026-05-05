package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/enzo959/forumium/models"
	"github.com/enzo959/forumium/repositories"
)

func GeneratePasswordReset(userID uint, email string) error {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}
	token := hex.EncodeToString(bytes)

	reset := &models.PasswordReset{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}

	if err := repositories.CreatePasswordReset(reset); err != nil {
		return fmt.Errorf("failed to save reset token: %w", err)
	}

	if err := SendPasswordResetEmail(email, token); err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	return nil
}