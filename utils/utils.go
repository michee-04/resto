package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type JsonResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ParseBody(r *http.Request, x interface{}) {
	if b, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal([]byte(b), x); err != nil {
			return
		}
	}
}

func HashedPassword(p string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(h), nil
}

func CheckPasswordHash(p, h string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(p))

	return err == nil
}

func GenerateVerificationToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func ResponseWithJson(w http.ResponseWriter, status int, message string, data interface{}) {

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	response := JsonResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
	r, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error processing response", http.StatusInternalServerError)
		return
	}

	w.Write(r)
}
