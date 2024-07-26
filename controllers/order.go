package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/michee-04/resto/models"
	"github.com/michee-04/resto/utils"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	order := &models.Order{}
	utils.ParseBody(r, &order)
	invoiceId := chi.URLParam(r, "invoiceId")
	order.InvoiceId = invoiceId
	o := order.CreateOrder()
	res, _ := json.Marshal(o)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	o := models.GetOrder()
	res, _ := json.Marshal(o)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetOrderId(w http.ResponseWriter, r *http.Request) {
	orderId := chi.URLParam(r, "orderId")
	o, _ := models.GetOrderId(orderId)
	res, _ := json.Marshal(o)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	orderId := chi.URLParam(r, "orderId")
	o := models.DeleteOrder(orderId)
	utils.ResponseWithJson(w, http.StatusOK, "order delete successful", o)
}
