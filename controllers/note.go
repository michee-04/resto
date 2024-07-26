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

func CreateNote(w http.ResponseWriter, r *http.Request) {
	note := &models.Note{}
	utils.ParseBody(r, &note)
	n := note.CreateNote()
	res, _ := json.Marshal(n)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetNote(w http.ResponseWriter, r *http.Request) {
	n := models.GetNote()
	res, _ := json.Marshal(n)
	w.Header().Set("content-type", "appplication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetNoteId(w http.ResponseWriter, r *http.Request) {
	noteId := chi.URLParam(r, "noteId")
	n, _ := models.GetNoteById(noteId)
	res, _ := json.Marshal(n)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateNote(w http.ResponseWriter, r *http.Request) {
	note := &models.Note{}
	utils.ParseBody(r, note)
	noteId := chi.URLParam(r, "noteId")

	var n models.Note
	err := models.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("note_id=?", noteId).First(&n).Error; err != nil {
			return err
		}
		if note.Text != "" {
			n.Text = note.Text
		}
		return tx.Save(&n).Error
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "note not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update note: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	utils.ResponseWithJson(w, http.StatusOK, "note update successful", n)
}

func DeleteNote(w http.ResponseWriter, r *http.Request) {
	noteId := chi.URLParam(r, "noteId")
	n := models.DeleteFood(noteId)
	utils.ResponseWithJson(w, http.StatusOK, "Note delete successful", n)
}
