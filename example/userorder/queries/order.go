package queries

import (
	"context"
	"example/userorder/models"
	qopt "example/userorder/queries/options"

	"github.com/foxie-io/gormqs"
	"gorm.io/gorm"
)

var (
	_ gormqs.Model   = (*models.Order)(nil)
	_ gormqs.Querier = (*OrderQueries)(nil)
)

type (
	OrderQueries struct {
		gormqs.Queries[models.Order, *OrderQueries]
		db    *gorm.DB
		model models.Order
	}
)

func (qs *OrderQueries) DBInstance(ctx context.Context) *gorm.DB {
	db := gormqs.ContextValue(ctx, qs.db)
	return db.WithContext(ctx).Table(qs.model.TableName()).Model(qs.model)
}

func NewOrderQueries(db *gorm.DB) *OrderQueries {
	qs := &OrderQueries{db: db}
	qs.Queries = gormqs.NewQueries[models.Order](qs)
	return qs
}

/*
	Add Custom Query to OrderQueries

alternative:

	qr.DBInstance(ctx).
		Where("id = ?", orderID).
		Joins("OrderItems").
		First(&order)
*/
func (qs *OrderQueries) GetOneWithDetails(ctx context.Context, orderID uint) (*models.Order, error) {
	return qs.GetOne(ctx,
		qopt.ORDER.Where(qopt.ORDER.ID, "=", orderID),
		qopt.ORDER.PreloadOrderItems(),
		qopt.ORDER.PreloadUser(),
	)
}
