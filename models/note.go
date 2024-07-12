package models

import (
	"github.com/google/uuid"
	"github.com/michee-04/resto/database"
	"gorm.io/gorm"
)

type Note struct {
	NoteId string `gorm:"primary_key;not null;unique;column:note_id" json:"note_id"`
	Text   string `gorm:"column:text" json:"text"`
}

func init() {
	database.ConnectDB()
	Db = database.GetBD()
	// Db.AutoMigrate(&Note{})
}

func (n *Note) CreateNote() *Note {
	n.NoteId = uuid.New().String()
	Db.Create(&n)
	return n
}

func GetNote() []Note {
	var n []Note
	Db.Find(&n)
	return n
}

func GetNoteById(Id string) (*Note, *gorm.DB) {
	var n Note
	db := Db.Where("note_id=?", Id).First(&n)
	return &n, db
}

func DeleteNote(Id string) Note {
	var n Note
	Db.Where("note_id=?", Id).Delete(&n)
	return n
}