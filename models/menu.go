package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/michee-04/resto/database"
	"gorm.io/gorm"
)

type Menu struct {
	MenuId    string    `gorm:"not null;primary_key;unique;column:menu_id" json:"menu_id"`
	Name      string    `gorm:"not null;unique;column:name" json:"name"`
	Category  string    `gorm:"column:category" json:"category"`
	StartDate time.Time `gorm:"column:start_date" json:"start_date"`
	EndDate   time.Time `gorm:"column:end_date" json:"end_date"`
	Food      []*Food
}

func init() {
	database.ConnectDB()
	Db = database.GetBD()
	// Db.AutoMigrate(&Menu{})
}

func (m *Menu) CreateMenu() *Menu {
	m.MenuId = uuid.New().String()
	Db.Create(&m)
	return m
}

func GetMenu() []Menu{
	var m []Menu
	Db.Preload("Food").Find(&m)
	return m
}

func GetMenuById(Id string) (*Menu, *gorm.DB) {
	var m Menu
	db := Db.Preload("Food").Where("menu_id=?", Id).First(&m)
	return &m, db
}

func DeleteMenu(Id string) Menu {
	var m Menu
	Db.Preload("Food").Where("menu_id=?", Id).Find(&m)
	for _, food := range m.Food {
		Db.Delete(&food)
	}
	Db.Delete(&m)
	return m
}