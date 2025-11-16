package queries

import (
	"context"
	"example/pagination/models"

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
func (qs *UserQueries) DBInstance(ctx context.Context) *gorm.DB {
	db := gormqs.ContextValue(ctx, qs.db)
	return db.WithContext(ctx).Table(qs.model.TableName()).Model(qs.model)
}

func NewUserQueries(db *gorm.DB) *UserQueries {
	qs := &UserQueries{db: db}
	qs.Queries = gormqs.NewQueries[models.User](qs)
	return qs
}
