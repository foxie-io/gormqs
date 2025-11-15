package qsopt

import (
	"fmt"

	"github.com/foxie-io/gormqs"
	"gorm.io/gorm"
)

type OrderColumn string

const (
	OrderID        OrderColumn = "id"
	OrderCreatedAt OrderColumn = "created_at"
	OrderUpdatedAt OrderColumn = "updated_at"
	OrderAmount    OrderColumn = "pay_amount"
	OrderDiscount  OrderColumn = "discount"
	OrderUserID    OrderColumn = "user_id"
)

func OrderWhere(col OrderColumn, operation, value any) gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		query := fmt.Sprintf("%s %s ?", gormqs.WithTable(string(col), db), operation)
		return db.Where(query, value)
	}
}

func OrderPreloadOrderItems() gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload("OrderItems")
	}
}

func OrderSelect(cols ...OrderColumn) gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		for _, col := range cols {
			db.Select(gormqs.WithTable(string(col), db))
		}

		return db
	}
}
