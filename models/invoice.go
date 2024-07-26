package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/michee-04/resto/database"
	"gorm.io/gorm"
)

type Invoice struct {
	InvoiceId     string `gorm:"primary_key;not null;unique;column:invoice_id" json:"invoice_id"`
	PaymentMethod string `gorm:"column:payment_method" json:"payment_method"`
	Status        string `gorm:"column:status" json:"status"`
	UserId        string `gorm:"not null; index;column:user_id" json:"user_id"`
	Date          time.Time
	User          *User
	Order         []*Order
}

func init() {
	database.ConnectDB()
	Db = database.GetBD()
	// Db.Migrator().DropTable(&Invoice{})
	// Db.AutoMigrate(&Invoice{})
}

func (i *Invoice) CreateInvoice() *Invoice {
	i.InvoiceId = uuid.New().String()
	Db.Create(&i)
	return i
}

func GetInvoice() []Invoice {
	var i []Invoice
	Db.Preload("User").Preload("Order").Find(&i)
	return i
}

func GetInvoiceId(Id string) (*Invoice, *gorm.DB) {
	var i Invoice
	db := Db.Preload("User").Preload("Order").Where("invoice_id=?", Id).First(&i)
	return &i, db
}

func DeleteInvoice(id string) Invoice {
	var i Invoice
	Db.Preload("User").Preload("Order").Where("invoice_id=?", id).Find(&i)
	for _, order := range i.Order {
		Db.Delete(&order)
	}
	Db.Delete(&i)
	return i
}
