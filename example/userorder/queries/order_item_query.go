package queries

import (
	"context"
	"example/userorder/models"

	"github.com/foxie-io/gormqs"
	"gorm.io/gorm"
)

var (
	_ gormqs.Model   = (*models.OrderItem)(nil)
	_ gormqs.Querier = (*OrderItemQueries)(nil)
)

type (
	OrderItemQueries struct {
		gormqs.Queries[models.OrderItem, *OrderItemQueries]
		db    *gorm.DB
		model models.OrderItem
	}
)

func (qs OrderItemQueries) DBInstance(ctx context.Context) *gorm.DB {
	db := gormqs.ContextValue(ctx, qs.db)
	return db.WithContext(ctx).Model(qs.model)
}

func NewOrderItemQueries(db *gorm.DB) *OrderItemQueries {
	qs := &OrderItemQueries{db: db}
	qs.Queries = gormqs.NewQueries[models.OrderItem](qs)
	return qs
}
