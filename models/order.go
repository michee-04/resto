package models

import "time"

type Order struct {
	OrderId   string `gorm:"primary_key;unique;column:order_id" json:"order_id"`
	Date      time.Time
	InvoiceId string `gorm:"column:invoice_id" json:"invoice_id"`
	TableId   string `gorm:"column:table_id" json:"table_id"`
	Invoice   *Invoice
	Table     *Table
	OrderItem []*OrderItem
}
