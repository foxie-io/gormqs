package queries

import (
	"context"
	"example/userorder/models"

	"github.com/foxie-io/gormqs"
	"gorm.io/gorm"
)

var (
	_ gormqs.Model   = (*models.Item)(nil)
	_ gormqs.Querier = (*ItemQuerier)(nil)
)

type (
	ItemQueries interface {
		gormqs.Queries[models.Item, *ItemQuerier]
	}

	ItemQuerier struct {
		queries gormqs.Queries[models.Item, *ItemQuerier]
		db      *gorm.DB
		model   models.Item
	}
)

func (qr ItemQuerier) DBInstance(ctx context.Context) *gorm.DB {
	db := gormqs.ContextValue(ctx, qr.db)
	return db.WithContext(ctx).Table(qr.model.TableName()).Model(qr.model)
}

func NewItemQueries(db *gorm.DB) ItemQueries {
	querier := &ItemQuerier{db: db}
	querier.queries = gormqs.NewQueries[models.Item](querier)
	return querier.queries
}
