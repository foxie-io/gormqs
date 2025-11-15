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
	UserQueries interface {
		gormqs.Queries[models.User, *UserQuerier]
	}

	UserQuerier struct {
		gormqs.Queries[models.User, *UserQuerier]
		db    *gorm.DB
		model models.User
	}
)

func (qr UserQuerier) DBInstance(ctx context.Context) *gorm.DB {
	db := gormqs.ContextValue(ctx, qr.db)
	return db.WithContext(ctx).Table(qr.model.TableName()).Model(qr.model)
}

func NewUserQueries(db *gorm.DB) UserQueries {
	querier := &UserQuerier{db: db}
	querier.Queries = gormqs.NewQueries[models.User](querier)
	return querier
}
