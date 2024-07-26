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

func CreateTable(w http.ResponseWriter, r *http.Request) {
	table := &models.Table{}
	utils.ParseBody(r, &table)
	t := table.CreateTable()
	res, _ := json.Marshal(t)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetTable(w http.ResponseWriter, r *http.Request) {
	t := models.GetTable()
	res, _ := json.Marshal(t)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetTableId(w http.ResponseWriter, r *http.Request) {
	tableId := chi.URLParam(r, "tableId")
	t, _ := models.GetTableById(tableId)
	res, _ := json.Marshal(t)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateTable(w http.ResponseWriter, r *http.Request) {
	tableUpdate := &models.Table{}
	utils.ParseBody(r, tableUpdate)

	tableId := chi.URLParam(r, "tableId")
	table, _ := models.GetTableById(tableId)

	var t models.Table
	err := models.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("table_id", tableId).First(&t).Error; err != nil {
			return err
		}

		if tableUpdate.TableNumber != table.TableNumber {
			t.TableNumber = tableUpdate.TableNumber
		}
		return tx.Save(&t).Error
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ResponseWithJson(w, http.StatusNotFound, "table not found", nil)
		} else {
			http.Error(w, "Failed to update table: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	utils.ResponseWithJson(w, http.StatusOK, "table update successful", t)
}

func DeletTable(w http.ResponseWriter, r *http.Request) {
	tableId := chi.URLParam(r, "tableId")
	t := models.DeleteTable(tableId)

	utils.ResponseWithJson(w, http.StatusOK, "table delete", t)
}
