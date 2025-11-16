package models

type User struct {
	ID             uint
	Username       string  `json:"username"`
	Balance        float64 `json:"balance"`
	BlockedBalance float64 `json:"blockedBalance"`
}

func (User) TableName() string {
	return "users"
}
