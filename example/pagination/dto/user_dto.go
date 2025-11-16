package dto

import (
	"github.com/foxie-io/gormqs"
	"gorm.io/gorm"
)

// type check
var (
	_ gormqs.Model = (*BaseUser)(nil)
)

type BaseUser struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

func (BaseUser) TableName() string {
	return "users"
}

type UserFilter struct {
	Username *string
}

func (u UserFilter) DBApply() gormqs.Option {
	return func(d *gorm.DB) *gorm.DB {
		if u.Username != nil {
			d = d.Where("username = ?", *u.Username)
		}

		return d
	}
}
