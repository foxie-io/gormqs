package gormqs

import (
	"context"

	"gorm.io/gorm"
)

type Option func(*gorm.DB) *gorm.DB

func Options(opts ...Option) Option {
	return func(db *gorm.DB) *gorm.DB {
		for _, opt := range opts {
			db = opt(db)
		}
		return db
	}
}

func Apply(q *gorm.DB, options []Option) *gorm.DB {
	for _, opt := range options {
		q = opt(q)
	}
	return q
}

// gorm model should contain TableName
type Model interface {
	TableName() string
}

type Querier interface {
	// instance use to build query
	DBInstance(ctx context.Context) *gorm.DB
}

type Queries[M Model, Q any] interface {
	// get the querier instance
	Querier() Q

	CreateOne(ctx context.Context, record *M) error
	CreateMany(ctx context.Context, record *[]*M) error

	GetOne(ctx context.Context, opts ...Option) (result *M, err error)
	GetMany(ctx context.Context, opts ...Option) (result []*M, err error)

	// update one or many is base on opt select
	Updates(ctx context.Context, record *M, opt Option, opts ...Option) (affectedRow int64, err error)
	Count(ctx context.Context, opt Option, opts ...Option) (count int64, err error)
	Delete(ctx context.Context, opt Option, opts ...Option) (affectedRow int64, err error)

	// scan pattern for custom type
	GetOneTo(ctx context.Context, result Model, opts ...Option) error
	GetManyTo(ctx context.Context, resultList any, opts ...Option) error
}

type defaultQueries[M Model, Querier any] struct {
	querier Querier
	model   M
}

func NewQueries[M Model, Q any](querier Q) Queries[M, Q] {
	qs := &defaultQueries[M, Q]{querier: querier}

	// interface check
	_ = qs.asQuerier()
	return qs
}

func (qs *defaultQueries[M, Q]) Querier() Q {
	return qs.querier
}

func (qs *defaultQueries[M, Q]) asQuerier() Querier {
	return any(qs.querier).(Querier)
}

func (qs *defaultQueries[M, Q]) dbInstance(ctx context.Context, opts ...Option) *gorm.DB {
	query := qs.asQuerier().DBInstance(ctx)
	return Apply(query, opts)
}

func (qs *defaultQueries[M, Q]) CreateOne(ctx context.Context, record *M) error {
	return qs.dbInstance(ctx).Create(record).Error
}

func (qs *defaultQueries[M, Q]) CreateMany(ctx context.Context, records *[]*M) error {
	return qs.dbInstance(ctx).Create(records).Error
}

func (qs *defaultQueries[M, Q]) GetOne(ctx context.Context, opts ...Option) (*M, error) {
	var result M
	err := qs.dbInstance(ctx, opts...).First(&result).Error
	return &result, err
}

func (qs *defaultQueries[M, Q]) GetMany(ctx context.Context, opts ...Option) ([]*M, error) {
	var result []*M
	err := qs.dbInstance(ctx, opts...).Find(&result).Error
	return result, err
}

func (qs *defaultQueries[M, Q]) Updates(ctx context.Context, record *M, opt Option, opts ...Option) (affectedRow int64, err error) {
	cmd := qs.dbInstance(ctx, opt, Options(opts...)).Model(record).Updates(record)
	affectedRow, err = cmd.RowsAffected, cmd.Error
	return
}

func (qs *defaultQueries[M, Q]) Count(ctx context.Context, opt Option, opts ...Option) (count int64, err error) {
	err = qs.dbInstance(ctx, opt, Options(opts...)).Count(&count).Error
	return
}

func (qs *defaultQueries[M, Q]) Delete(ctx context.Context, opt Option, opts ...Option) (affectedRow int64, err error) {
	cmd := qs.dbInstance(ctx, opt, Options(opts...)).Delete(qs.model)
	affectedRow, err = cmd.RowsAffected, cmd.Error
	return
}

func (qs *defaultQueries[M, Q]) GetOneTo(ctx context.Context, result Model, opts ...Option) error {
	return qs.dbInstance(ctx, opts...).Find(result).Error
}

func (qs *defaultQueries[M, Q]) GetManyTo(ctx context.Context, resultList any, opts ...Option) error {
	return qs.dbInstance(ctx, opts...).Find(resultList).Error
}
