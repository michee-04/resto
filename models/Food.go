package models

import (
	"github.com/google/uuid"
	"github.com/michee-04/resto/database"
	"gorm.io/gorm"
)

type Food struct {
	FoodId    string `gorm:"unique;not null;column:food_id" json:"food_id"`
	Name      string `gorm:"column:name" json:"name"`
	Price     string `gorm:"not null; column:price" json:"price"`
	FoodImage string `gorm:"not null; column:food_image" json:"food_image"`
	MenuId    string `gorm:"not null; index;foreignKey:MenuId"`
	Menu      *Menu
}

func init() {
	database.ConnectDB()
	Db = database.GetBD()
	// Db.AutoMigrate(&Food{})
}

func (f *Food) CreateFood() *Food {
	f.FoodId = uuid.New().String()
	Db.Create(&f)
	return f
}

func GetFood() []Food {
	var f []Food
	Db.Preload("Menu").Find(&f)
	return f
}

func GetFoodById(Id string) (*Food, *gorm.DB) {
	var f Food
	db := Db.Preload("Menu").Where("food_id=?", Id).First(&f)
	return &f, db
}

func DeleteFood(Id string) Food {
	var f Food
	Db.Preload("Menu").Where("food_id").Delete(&f)
	return f
} 