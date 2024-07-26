package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/michee-04/resto/models"
	"github.com/michee-04/resto/provider"
	"github.com/michee-04/resto/services"
	"github.com/michee-04/resto/utils"
	"gorm.io/gorm"
)

type UserResponse struct {
	UserId      string `json:"user_id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	IsAdmin     bool   `json:"is_admin"`
	Token       string `json:"token"`
	EmailVerify bool   `json:"emil_verify"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	utils.ParseBody(r, &user)

	emailExists, err := models.GetEmailExists(user.Email)
	if err != nil {
		utils.ResponseWithJson(w, http.StatusInternalServerError, "Error", err)
		return
	}
	if emailExists {
		utils.ResponseWithJson(w, http.StatusConflict, "email already exists", err)
		return
	}

	hashedPassword, _ := utils.HashedPassword(user.Password)
	emailToken := utils.GenerateVerificationToken()
	user.Password = hashedPassword
	user.TokenAccount = emailToken

	u := user.CreateUser()
	res, _ := json.Marshal(u)
	services.SendVerificationAccount(u)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	u := models.GetUser()
	res, _ := json.Marshal(u)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	u, _ := models.GetUserById(userId)
	res, _ := json.Marshal(u)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	u := models.DeleteUser(userId)

	utils.ResponseWithJson(w, http.StatusOK, "Delete user successful", u)
}

func UpddateUser(w http.ResponseWriter, r *http.Request) {
	userUpdate := &models.User{}
	utils.ParseBody(r, userUpdate)

	userId := chi.URLParam(r, "userId")

	if !provider.VerificationToken(r, userId) {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var u models.User
	err := models.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id=?", userId).First(&u).Error; err != nil {
			return err
		}

		if userUpdate.Username != "" {
			u.Username = userUpdate.Username
		}
		if userUpdate.Password != "" {
			hashedPassword, _ := utils.HashedPassword(userUpdate.Password)
			u.Password = hashedPassword
		}
		if userUpdate.Email != "" {
			u.Email = userUpdate.Email
			services.SendVerificationAccount(&u)
		}

		return tx.Save(&u).Error
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res, err := json.Marshal(u)
	if err != nil {
		http.Error(w, "Failed update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByEmail(loginReq.Email)
	if err != nil {
		http.Error(w, "Database querry failed", http.StatusInternalServerError)
		return
	}

	if user == nil {
		utils.ResponseWithJson(w, http.StatusUnauthorized, "Invalid email", nil)
	}

	if !user.EmailVerify {
		utils.ResponseWithJson(w, http.StatusUnauthorized, "Email not verified", nil)
		return
	}
	if user.Password == "" {
		http.Error(w, "Add Password", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(loginReq.Password, user.Password) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := provider.GenerateJWT(user.UserId, user.IsAdmin)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	user.Token = token
	models.Db.Save(&user)
	userResponse := UserResponse{
		UserId:      user.UserId,
		Username:    user.Username,
		Email:       user.Email,
		IsAdmin:     user.IsAdmin,
		Token:       user.Token,
		EmailVerify: user.EmailVerify,
	}

	utils.ResponseWithJson(w, http.StatusOK, "Login successful", userResponse)
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	u, _ := models.GetUserById(userId)

	if u == nil {
		http.Error(w, "user not found", http.StatusInternalServerError)
		return
	}

	if !provider.VerificationToken(r, userId) {
		http.Error(w, "Invalid token", http.StatusInternalServerError)
		return
	}

	u.Logout()

	utils.ResponseWithJson(w, http.StatusOK, "Logout successful", nil)
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	user, err := models.FindUserByToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	err = user.Verify()
	if err != nil {
		log.Printf("Erreur lors de la vérification de l'utilisateur: %v\n", err)
		http.Error(w, "Unable to verify user", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("template/emailverification.tmpl"))

	w.Header().Set("content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := models.Db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := user.GeneratePasswordToken(); err != nil {
		http.Error(w, "Failed to generate reset token", http.StatusInternalServerError)
		return
	}

	services.SendResetPasswordAccoubt(&user)

	utils.ResponseWithJson(w, http.StatusOK, "Veuillez verifier votre email pour la reinitialisation du mot de passe", user.TokenPassword)

}

func ResetPasswordEmail(w http.ResponseWriter, r *http.Request) {
	// Parse le fichier HTML
	tmpl := template.Must(template.ParseFiles("template/emailPassword.tmpl"))

	// Définir le content-type à text/html
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	// Exécuter le template avec une structure de données vide
	tmpl.Execute(w, nil)
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TokenPassword string `json:"tokenPassword"`
		Password      string `json:"password"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	fmt.Println("Received body:", string(body))

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	fmt.Println("Received tokenPassword:", req.TokenPassword)

	// Récupérer l'utilisateur par le token de réinitialisation de mot de passe
	user, err := models.FindUserPasswordToken(req.TokenPassword)
	if err != nil {
		if err.Error() == "reset token has expired" {
			http.Error(w, "Reset token has expired", http.StatusUnauthorized)
		} else {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		}
		return
	}

	// Vérifier si le token récupéré correspond au token fourni dans la requête
	if user.TokenPassword != req.TokenPassword {
		http.Error(w, "Invalid reset token", http.StatusUnauthorized)
		return
	}

	// Mettre à jour le mot de passe de l'utilisateur
	err = user.UpdatePassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, "Password updated successfully", nil)
}

func SetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère l'email depuis le header ou les paramètres de la requête
	email := r.Header.Get("X-User-Email")
	if email == "" {
		email = chi.URLParam(r, "email")
		if email == "" {
			http.Error(w, "Email not provided", http.StatusBadRequest)
			return
		}
	}

	var req struct {
		Password string `json:"password"`
	}

	// Parse le corps de la requête pour obtenir le mot de passe
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Recherche l'utilisateur par email
	user, err := models.GetUserByEmail(email)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}

	// Vérifie si l'utilisateur n'a pas de mot de passe
	if user.Password != "" {
		http.Error(w, "Password already set", http.StatusBadRequest)
		return
	}

	// Hache le nouveau mot de passe
	hashedPassword, err := utils.HashedPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Met à jour le mot de passe de l'utilisateur
	user.Password = hashedPassword
	if err := models.Db.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, "Password set successfully", nil)
}
