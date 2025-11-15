package models

type Order struct {
	Base
	PayAmount      float64      `json:"payAmount"`
	DiscountAmount float64      `json:"discount"`
	UserID         uint         `json:"userId"`
	OrderItems     []*OrderItem `json:"orderItems,omitempty" gorm:"foreignKey:OrderID"`
	User           *User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (Order) TableName() string {
	return "orders"
}
