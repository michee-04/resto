package models

import (
	"github.com/google/uuid"
	"github.com/michee-04/resto/database"
	"gorm.io/gorm"
)

type Table struct {
	TableId     string `gorm:"primary_key;unique;column:table_id" json:"table_id"`
	NumberGuest int    `gorm:"column:number_guest" json:"number_guest"`
	TableNumber int    `gorm:"column:table_number" json:"table_number"`
	Order       []Order
}

func init() {
	database.ConnectDB()
	Db = database.GetBD()
	// Db.AutoMigrate(&Table{})
}

func (t *Table) CreateTable() *Table {
	t.TableId = uuid.New().String()
	Db.Create(&t)
	return t
}

func GetTable() []Table {
	var t []Table
	Db.Find(&t)
	return t
}

func GetTableById(Id string) (*Table, *gorm.DB) {
	var t Table
	db := Db.Preload("Order").Where("table_id=?", Id).First(&t)
	return &t, db
}

func DeleteTable(Id string) Table {
	var t Table
	Db.Preload("Order").Where("table_id=?", Id).Find(&t)
	for _, order := range t.Order {
		Db.Delete(&order)
	}
	Db.Delete(&t)
	return t
}
