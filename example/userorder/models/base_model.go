package models

import "time"

type Base struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

func NewBase() Base {
	return Base{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
