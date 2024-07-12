package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/michee-04/resto/models"
	"github.com/michee-04/resto/utils"
)


func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	utils.ParseBody(r, user)
	hashedPassword, _ := utils.HashedPassword(user.Password)
	user.Password = hashedPassword
	// services.SendSms(user)

	u := user.CreateUser()
	res, _ := json.Marshal(u)
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