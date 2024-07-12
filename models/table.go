package models

type Table struct {
	TableId     string `gorm:"primary_key;unique;column:table_id" json:"table_id"`
	NumberGuest int    `gorm:"column:number_guest" json:"number_guest"`
	TableNumber int    `gorm:"column:table_number" json:"table_number"`
	Orders      []Order
}
