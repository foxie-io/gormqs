package queries

import (
	"context"
	"example/userorder/models"
	qsopt "example/userorder/queries/options"

	"github.com/foxie-io/gormqs"
	"gorm.io/gorm"
)

var (
	_ gormqs.Model   = (*models.Order)(nil)
	_ gormqs.Querier = (*OrderQuerier)(nil)
)

type (
	OrderQueries interface {
		gormqs.Queries[models.Order, *OrderQuerier]
	}

	OrderQuerier struct {
		queries gormqs.Queries[models.Order, *OrderQuerier]
		db      *gorm.DB
		model   models.Order
	}
)

func (qr *OrderQuerier) DBInstance(ctx context.Context) *gorm.DB {
	db := gormqs.ContextValue(ctx, qr.db)
	return db.WithContext(ctx).Table(qr.model.TableName()).Model(qr.model)
}

func NewOrderQueries(db *gorm.DB) OrderQueries {
	querier := &OrderQuerier{db: db}
	querier.queries = gormqs.NewQueries[models.Order](querier)
	return querier.queries
}

/*
	Add Custom Query

use gorm:

	qr.DBInstance(ctx).
		Where("id = ?", orderID).
		Joins("Items").
		First(&order)

	=

use queries: for reusable options and type safe

	qr.queries.GetOne(ctx,
		qsopt.OrderWhere(qsopt.OrderID, " = ", orderID),
		qsopt.OrderJoinItems(),
	)
*/
func (qr *OrderQuerier) GetOneWithItems(ctx context.Context, orderID uint) (*models.Order, error) {
	return qr.queries.GetOne(ctx,
		qsopt.OrderWhere(qsopt.OrderID, "=", orderID),
		qsopt.OrderJoinItems(),
	)
}
