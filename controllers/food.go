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

func CreateFodd(w http.ResponseWriter, r *http.Request) {
	food := &models.Food{}
	utils.ParseBody(r, &food)

	menuId := chi.URLParam(r, "menuId")
	food.MenuId = menuId
	f := food.CreateFood()
	res, _ := json.Marshal(f)
	w.Header().Set("conten-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetFood(w http.ResponseWriter, r *http.Request) {
	f := models.GetFood()
	res, _ := json.Marshal(f)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetfoodId(w http.ResponseWriter, r *http.Request) {
	foodId := chi.URLParam(r, "foodId")
	f, _ := models.GetFoodById(foodId)
	res, _ := json.Marshal(f)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateFood(w http.ResponseWriter, r *http.Request) {
	foodUpdate := &models.Food{}
	utils.ParseBody(r, foodUpdate)
	foodId := chi.URLParam(r, "foodId")

	var f models.Food
	err := models.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("food_id=?", foodId).First(&f).Error; err != nil {
			return err
		}

		if foodUpdate.Name != "" {
			f.Name = foodUpdate.Name
		}
		if foodUpdate.Price != "" {
			f.Price = foodUpdate.Price
		}
		if foodUpdate.FoodImage != "" {
			f.FoodImage = foodUpdate.FoodImage
		}
		return tx.Save(&f).Error
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "food not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update food: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, "food update successful", f)
}

func Deletefood(w http.ResponseWriter, r *http.Request) {
	foodId := chi.URLParam(r, "foodId")
	f := models.DeleteFood(foodId)
	utils.ResponseWithJson(w, http.StatusOK, "food delete successful", f)
}
