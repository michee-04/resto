package models

import "time"

type Invoice struct {
	InvoiceId     string `gorm:"primary_key;not null;unique;column:invoice_id" json:"invoice_id"`
	PaymentMethod string `gorm:"column:payment_method" json:"payment_method"`
	PaymentStatus string `gorm:"column:payment_status" json:"payment_status"`
	Date          time.Time
	Order         []*Order
}
