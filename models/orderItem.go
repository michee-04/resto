package models

type OrderItem struct {
	OrderItemId string `gorm:"primary_key;not null;unique;column:order_item_id" json:"order_item_id"`
	Quantity    string `gorm:"column:quantity" json:"quantity"`
	UnitPrice   string `gorm:"column:unit_price" json:"unit_price"`
	OrderId     string `gorm:"column:order_id" json:"order_id"`
	Order       *Order
}
