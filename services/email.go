package services

import (
	"fmt"
	"os"

	"github.com/michee-04/resto/models"
	"gopkg.in/gomail.v2"
)

// Fonction pour l'envoi de l'email sur l'adresse de l'utilisateur pour activer son compte
func SendVerificationAccount(u *models.User) {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "voteprojet@gmail.com")
	m.SetHeader("To", u.Email)
	m.SetHeader("Subject", "Veuillez activer votre compte d'utilisateur")
	m.SetBody("text/html", fmt.Sprintf("Cliquer sur <a href=\"%s/auth/verify?token=%s\">Ici</a> pour verifier votre adresse email", baseURL, u.TokenAccount))

	d := gomail.NewDialer("smtp.gmail.com", 587, "voteprojet@gmail.com", "jmbd aicq hdov mvyq")
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

// Fonction pour la reinitialisation du mot de passe de l'utilisateur par email
func SendResetPasswordAccoubt(u *models.User) {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "voteprojet@gmail.com")
	m.SetHeader("To", u.Email)
	m.SetHeader("Subject", "Reinitialisation du mot de passe")
	m.SetBody("text/html", fmt.Sprintf("Cliquez <a href=\"%s/auth/reset-password-email?token=%s\">ici</a> pour reinitialiser le mot de passe. Ce lien est valide pour une duree d'une heure.", baseURL, u.TokenPassword))

	d := gomail.NewDialer("smtp.gmail.com", 587, "voteprojet@gmail.com", "jmbd aicq hdov mvyq")
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
