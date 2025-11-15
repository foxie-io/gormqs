package queries

import (
	"context"
	"example/userorder/models"

	"github.com/foxie-io/gormqs"
	"gorm.io/gorm"
)

var (
	_ gormqs.Model   = (*models.User)(nil)
	_ gormqs.Querier = (*UserQuerier)(nil)
)

type (
	UserQuerier struct {
		gormqs.Queries[models.User, *UserQuerier]
		db    *gorm.DB
		model models.User
	}
)

func (qr UserQuerier) DBInstance(ctx context.Context) *gorm.DB {
	// if ctx has db instance will use that if not use default
	dbOrTx := gormqs.ContextValue(ctx, qr.db)
	return dbOrTx.WithContext(ctx).Table(qr.model.TableName()).Model(qr.model)
}

func NewUserQueries(db *gorm.DB) *UserQuerier {
	querier := &UserQuerier{db: db}
	querier.Queries = gormqs.NewQueries[models.User](querier)
	return querier
}
