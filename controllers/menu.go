package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/michee-04/resto/models"
	"github.com/michee-04/resto/utils"
	"gorm.io/gorm"
)


func CreateMenu(w http.ResponseWriter, r *http.Request) {
	menu := &models.Menu{}
	utils.ParseBody(r, menu)
	m := menu.CreateMenu()
	res, _ := json.Marshal(m)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetMenu(w http.ResponseWriter, r *http.Request) {
	m := models.GetMenu()
	res, _ := json.Marshal(m)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetMenuById(w http.ResponseWriter, r *http.Request) {
	menuId := chi.URLParam(r, "menuId")
	m, _ := models.GetMenuById(menuId)
	res, _ := json.Marshal(m)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateMenu(w http.ResponseWriter, r *http.Request) {
	menuUpdate := &models.Menu{}
	utils.ParseBody(r, menuUpdate)
	menuId := chi.URLParam(r, "menuId")

	var m models.Menu
	err := models.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("menu_id=?", menuId).First(&m).Error; err != nil {
			return err
		}

		if menuUpdate.Category != "" {
			m.Category = menuUpdate.Category
		}
		if menuUpdate.Name != "" {
			m.Name = menuUpdate.Name
		}
		return tx.Save(&m).Error
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "menu not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update menu: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, "Menu update successful", nil)
	res, _ := json.Marshal(m)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func Deletemenu(w http.ResponseWriter, r *http.Request) {
	menuId := chi.URLParam(r, "menuId")
	m := models.DeleteMenu(menuId)
	utils.ResponseWithJson(w, http.StatusOK, "Menu delete successful", m)
}