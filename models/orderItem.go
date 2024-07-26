package models

import (
	"github.com/google/uuid"
	"github.com/michee-04/resto/database"
	"gorm.io/gorm"
)

type OrderItem struct {
	OrderItemId string `gorm:"primary_key;not null;unique;column:order_item_id" json:"order_item_id"`
	Quantity    string `gorm:"column:quantity" json:"quantity"`
	UnitPrice   string `gorm:"column:unit_price" json:"unit_price"`
	OrderId     string `gorm:"not null;index;column:order_id" json:"order_id"`
	Order       *Order
}

func init() {
	database.ConnectDB()
	Db = database.GetBD()
	// Db.AutoMigrate(&OrderItem{})
}

func (o *OrderItem) CreateOrderItem() *OrderItem {
	o.OrderItemId = uuid.New().String()
	Db.Create(&o)
	return o
}

func GetOrderItem() []OrderItem {
	var o []OrderItem
	Db.Preload("Order").Find(&o)
	return o
}

func GetOrderItemId(id string) (*OrderItem, *gorm.DB) {
	var o OrderItem
	db := Db.Preload("Order").Where("order_item_id=?", id).First(&o)
	return &o, db
}

func DeleteOrderItem(Id string) OrderItem {
	var o OrderItem
	Db.Preload("Order").Where("order_item_id=?", Id).Delete(&o)
	return o
}
