package services

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendPasswordResetEmail(toEmail string, token string) error {
	host := os.Getenv("MAILTRAP_HOST")
	port := os.Getenv("MAILTRAP_PORT")
	user := os.Getenv("MAILTRAP_USER")
	pass := os.Getenv("MAILTRAP_PASS")

	resetLink := fmt.Sprintf("http://localhost:5173/reset-password?token=%s", token)
	subject := "Réinitialisation de votre mot de passe"
	body := fmt.Sprintf("Bonjour,\n\nCliquez sur ce lien pour réinitialiser votre mot de passe :\n%s\n\nCe lien expire dans 1 heure.", resetLink)
	message := fmt.Sprintf("Subject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s", subject, body)

	auth := smtp.PlainAuth("", user, pass, host)
	addr := fmt.Sprintf("%s:%s", host, port)

	err := smtp.SendMail(addr, auth, "noreply@forumium.com", []string{toEmail}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}