package models

type User struct {
	Base
	Username string  `json:"username"`
	Balance  float64 `json:"balance"`
}

func (User) TableName() string {
	return "users"
}
