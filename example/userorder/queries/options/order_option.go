package qsopt

import (
	"fmt"

	"github.com/foxie-io/gormqs"
	"gorm.io/gorm"
)

type OrderColumn string

type OrderSchema struct {
	ID        OrderColumn
	CreatedAt OrderColumn
	UpdatedAt OrderColumn
	Discount  OrderColumn
	UserID    OrderColumn
}

var ORDER = OrderSchema{
	ID:        "id",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	Discount:  "discount",
	UserID:    "user_id",
}

func (s OrderSchema) Where(col OrderColumn, operation, value any) gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		query := fmt.Sprintf("%s %s ?", gormqs.WithTable(string(col), db), operation)
		return db.Where(query, value)
	}
}

func (s OrderSchema) Select(cols ...OrderColumn) gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		columns := make([]string, len(cols))
		for i, col := range cols {
			columns[i] = gormqs.WithTable(string(col), db)
		}

		return db.Select(columns)
	}
}

func (s OrderSchema) WhereID(id uint) gormqs.Option {
	return s.Where(s.ID, "=", id)
}

func (s OrderSchema) PreloadOrderItems() gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload("OrderItems")
	}
}

func (s OrderSchema) PreloadUser() gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload("User")
	}
}
