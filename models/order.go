package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/michee-04/resto/database"
	"gorm.io/gorm"
)

type Order struct {
	OrderId   string `gorm:"primary_key;unique;column:order_id" json:"order_id"`
	Date      time.Time
	InvoiceId string `gorm:"not null;index;column:invoice_id" json:"invoice_id"`
	Invoice   *Invoice
	OrderItem []*OrderItem
}

func init() {
	database.ConnectDB()
	Db = database.GetBD()
	// Db.AutoMigrate(&Order{})
}

func (o *Order) CreateOrder() *Order {
	o.OrderId = uuid.New().String()
	Db.Create(&o)
	return o
}

func GetOrder() []Order {
	var o []Order
	Db.Find(&o)
	return o
}

func GetOrderId(Id string) (*Order, *gorm.DB) {
	var o Order
	db := Db.Preload("OrderItem").Where("order_id", Id).First(&o)
	return &o, db
}

func DeleteOrder(Id string) Order {
	var o Order
	Db.Preload("OrderItem").Where("order_id=?", Id).Find(&o)
	for _, orderitem := range o.OrderItem {
		Db.Delete(&orderitem)
	}
	Db.Delete(&o)
	return o
}
