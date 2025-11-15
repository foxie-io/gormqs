package models

type Item struct {
	Base
	OrderID  uint    `json:"orderId"`
	Order    *Order  `json:"order" gorm:"foreignKey:OrderID;references:ID" `
	Product  string  `json:"product"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func (Item) TableName() string {
	return "items"
}
