package models

type Order struct {
	Base
	Amount   float64 `json:"amount"`
	Discount float64 `json:"discount"`
	UserID   uint    `json:"userId"`

	Items []*Item `json:"items" gorm:"foreignKey:OrderID;references:ID"`
	User  *User   `json:"user" gorm:"foreignKey:UserID"`
}

func (Order) TableName() string {
	return "orders"
}
