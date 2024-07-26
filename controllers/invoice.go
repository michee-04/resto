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

func CreateInvoice(w http.ResponseWriter, r *http.Request) {
	invoice := &models.Invoice{}
	utils.ParseBody(r, &invoice)
	userId := chi.URLParam(r, "userId")
	invoice.UserId = userId
	i := invoice.CreateInvoice()
	res, _ := json.Marshal(i)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetInvoice(w http.ResponseWriter, r *http.Request) {
	i := models.GetInvoice()
	res, _ := json.Marshal(i)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetInvoiceId(w http.ResponseWriter, r *http.Request) {
	invoiceId := chi.URLParam(r, "invoiceId")
	i, _ := models.GetInvoiceId(invoiceId)
	res, _ := json.Marshal(i)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	invoiceUpdate := &models.Invoice{}
	utils.ParseBody(r, invoiceUpdate)
	invoiceId := chi.URLParam(r, "invoiceId")

	var i models.Invoice
	err := models.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("invoice_id=?", invoiceId).First(&i).Error; err != nil {
			return err
		}
		if invoiceUpdate.PaymentMethod != "" {
			i.PaymentMethod = invoiceUpdate.PaymentMethod
		}
		if invoiceUpdate.Status != "" {
			i.Status = invoiceUpdate.Status
		}
		return tx.Save(&i).Error
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "invoice not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update invoice"+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, "invoice update", i)
}

func DeleteInvoice(w http.ResponseWriter, r *http.Request) {
	invoiceId := chi.URLParam(r, "invoiceId")
	i := models.DeleteInvoice(invoiceId)
	utils.ResponseWithJson(w, http.StatusOK, "invoice delete successful", i)
}
