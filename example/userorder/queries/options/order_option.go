package qsopt

import (
	"fmt"

	"github.com/foxie-io/gormqs"
	"gorm.io/gorm"
)

type OrderColumn string

const (
	OrderID       OrderColumn = "id"
	OrderAmount   OrderColumn = "amount"
	OrderDiscount OrderColumn = "discount"
	OrderUserID   OrderColumn = "user_id"
)

func OrderWhere(col OrderColumn, operation, value any) gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		query := fmt.Sprintf("%s %s ?", gormqs.WithTable(string(col), db), operation)
		return db.Where(query, value)
	}
}

func OrderJoinItems() gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins("Items")
	}
}

func OrderSelect(cols ...OrderColumn) gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		for _, col := range cols {
			db.Statement.Selects = append(db.Statement.Selects, gormqs.WithTable(string(col), db))
		}
		return db
	}
}
