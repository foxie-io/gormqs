package queries

import (
	"context"
	"example/userorder/models"

	"github.com/foxie-io/gormqs"
	"gorm.io/gorm"
)

var (
	_ gormqs.Model   = (*models.User)(nil)
	_ gormqs.Querier = (*UserQueries)(nil)
)

type (
	UserQueries struct {
		gormqs.Queries[models.User, *UserQueries]
		db    *gorm.DB
		model models.User
	}
)

// provider db instance for gormqs.Queries to use
func (qr UserQueries) DBInstance(ctx context.Context) *gorm.DB {
	db := gormqs.ContextValue(ctx, qr.db)
	return db.WithContext(ctx).Table(qr.model.TableName()).Model(qr.model)
}

func NewUserQueries(db *gorm.DB) *UserQueries {
	qs := &UserQueries{db: db}
	qs.Queries = gormqs.NewQueries[models.User](qs)
	return qs
}
