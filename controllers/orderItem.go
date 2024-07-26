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

func CreateOrderItem(w http.ResponseWriter, r *http.Request) {
	orderItem := &models.OrderItem{}
	utils.ParseBody(r, &orderItem)
	o := orderItem.CreateOrderItem()
	res, _ := json.Marshal(o)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetOrderItem(w http.ResponseWriter, r *http.Request) {
	o := models.GetOrderItem()
	res, _ := json.Marshal(o)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetOrderItemId(w http.ResponseWriter, r *http.Request) {
	orderItemId := chi.URLParam(r, "orderItemId")
	o, _ := models.GetOrderItemId(orderItemId)
	res, _ := json.Marshal(o)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateOrderItem(w http.ResponseWriter, r *http.Request) {
	prderItemUpdate := &models.OrderItem{}
	utils.ParseBody(r, prderItemUpdate)
	orderItemId := chi.URLParam(r, "orderItemId")

	var o models.OrderItem
	err := models.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("food_id=?", orderItemId).First(&o).Error; err != nil {
			return err
		}

		if prderItemUpdate.Quantity != "" {
			o.Quantity = prderItemUpdate.Quantity
		}
		if prderItemUpdate.UnitPrice != "" {
			o.UnitPrice = prderItemUpdate.UnitPrice
		}
		return tx.Save(&o).Error
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "food not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update food: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, "food update successful", o)
}

func DeleteOrderIten(w http.ResponseWriter, r *http.Request) {
	orderItemId := chi.URLParam(r, "orderItemId")
	o := models.DeleteOrderItem(orderItemId)
	utils.ResponseWithJson(w, http.StatusOK, "Order item delete successful", o)
}
