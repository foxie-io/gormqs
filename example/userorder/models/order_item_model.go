package models

type OrderItem struct {
	OrderID  uint    `json:"orderId" gorm:"primaryKey"`
	ItemID   uint    `json:"itemId" gorm:"primaryKey"`
	Quantity uint    `json:"multiple"`
	Price    float64 `json:"price"`
	Discount float64 `json:"discount"`

	Item  *Item  `json:"item,omitempty"`
	Order *Order `json:"order,omitempty"`
}

func (item OrderItem) PayAmount() float64 {
	total := item.TotalAmount()
	discount := item.DiscountAmount()

	if total < discount {
		return 0
	}

	return total - discount
}

func (item OrderItem) TotalAmount() float64 {
	price := item.Price
	quantity := float64(item.Quantity)

	if price <= 0 || quantity <= 0 {
		return 0
	}

	return price * quantity
}

func (item OrderItem) DiscountAmount() float64 {
	total := item.TotalAmount()
	discount := item.Discount

	if total <= 0 || discount <= 0 || discount >= 1 {
		return 0
	}

	return total * discount
}

func (OrderItem) TableName() string {
	return "order_items"
}
