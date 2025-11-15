package queries

import (
	"context"
	"example/userorder/models"

	"github.com/foxie-io/gormqs"
	"gorm.io/gorm"
)

var (
	_ gormqs.Model   = (*models.Item)(nil)
	_ gormqs.Querier = (*ItemQueries)(nil)
)

type (
	ItemQueries struct {
		gormqs.Queries[models.Item, *ItemQueries]
		db    *gorm.DB
		model models.Item
	}
)

func (qr ItemQueries) DBInstance(ctx context.Context) *gorm.DB {
	db := gormqs.ContextValue(ctx, qr.db)
	return db.WithContext(ctx).Model(qr.model)
}

func NewItemQueries(db *gorm.DB) *ItemQueries {
	qs := &ItemQueries{db: db}
	qs.Queries = gormqs.NewQueries[models.Item](qs)
	return qs
}
