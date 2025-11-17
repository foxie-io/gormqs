package models

type Item struct {
	Base
	Product  string  `json:"product"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func (Item) TableName() string {
	return "items"
}
