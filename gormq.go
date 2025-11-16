package gormqs

import (
	"context"
	"errors"
	"strings"

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

/*
Gorm Model implement

	type Order struct{
		ID uint
		Price float64
	}

	func (Order) TableName() string {
	 	return "orders"
	}
*/
type Model interface {
	TableName() string
}

/*
Resulter implement

	type Order struct{
		ID uint
		Price float64
	}

	func (Order) TableName() string {
	 	return "orders"
	}

	type Response struct {
		Total int
		Data []*Order
	}

	func (f *Filter) QsList() any {
		return &f.Data
	}

	func (f *Filter) QsCount() *int64 {
		return &f.Total
	}
*/
type ManyWithCountResulter interface {
	QsList() any
	QsCount() *int64
}

type ManyWithCountResult[T Model] struct {
	Select string `json:"select"`
	List   *[]*T  `json:"list,omitempty"`
	Count  *int64 `json:"count,omitempty"`
}

func NewManyWithCountResulter[T Model](selects string) *ManyWithCountResult[T] {
	return &ManyWithCountResult[T]{
		Select: selects,
	}
}

func (r *ManyWithCountResult[T]) QsList() any {
	if strings.Contains(r.Select, "list") {
		return r.List
	}
	return nil
}

func (r *ManyWithCountResult[T]) QsCount() *int64 {
	if strings.Contains(r.Select, "count") {
		return r.Count
	}
	return nil
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

	// update one or many is base on at least one option
	Updates(ctx context.Context, record *M, opt Option, opts ...Option) (affectedRow int64, err error)

	// must have at least one option
	Count(ctx context.Context, opt Option, opts ...Option) (count int64, err error)

	// must have at least one option
	Delete(ctx context.Context, opt Option, opts ...Option) (affectedRow int64, err error)

	// scan pattern for custom type without mapping to struct again

	GetOneTo(ctx context.Context, result Model, opts ...Option) error
	GetManyTo(ctx context.Context, resultList any, opts ...Option) error

	// list with count | list ony | count ony with same option meaning same query filter
	GetManyWithCount(ctx context.Context, r ManyWithCountResulter, options ...Option) error
}

type queries[M Model, Querier any] struct {
	querier Querier
	model   M
}

func NewQueries[M Model, Q any](querier Q) Queries[M, Q] {
	qs := &queries[M, Q]{querier: querier}

	// interface check
	_ = qs.asQuerier()
	return qs
}

func (qs *queries[M, Q]) Querier() Q {
	return qs.querier
}

func (qs *queries[M, Q]) asQuerier() Querier {
	return any(qs.querier).(Querier)
}

func (qs *queries[M, Q]) dbInstance(ctx context.Context, opts ...Option) *gorm.DB {
	query := qs.asQuerier().DBInstance(ctx)
	return Apply(query, opts)
}

func (qs *queries[M, Q]) CreateOne(ctx context.Context, record *M) error {
	return qs.dbInstance(ctx).Create(record).Error
}

func (qs *queries[M, Q]) CreateMany(ctx context.Context, records *[]*M) error {
	return qs.dbInstance(ctx).Create(records).Error
}

func (qs *queries[M, Q]) GetOne(ctx context.Context, opts ...Option) (*M, error) {
	var result M
	if err := qs.dbInstance(ctx, opts...).Model(&result).First(&result).Error; err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &result, nil
}

func (qs *queries[M, Q]) GetMany(ctx context.Context, opts ...Option) ([]*M, error) {
	var result []*M
	err := qs.dbInstance(ctx, opts...).Find(&result).Error
	return result, err
}

func (qs *queries[M, Q]) Updates(ctx context.Context, record *M, opt Option, opts ...Option) (affectedRow int64, err error) {
	cmd := qs.dbInstance(ctx, opt, Options(opts...)).Model(record).Updates(record)
	affectedRow, err = cmd.RowsAffected, cmd.Error
	return
}

func (qs *queries[M, Q]) Count(ctx context.Context, opt Option, opts ...Option) (count int64, err error) {
	err = qs.dbInstance(ctx, opt, Options(opts...)).Count(&count).Error
	return
}

func (qs *queries[M, Q]) Delete(ctx context.Context, opt Option, opts ...Option) (affectedRow int64, err error) {
	cmd := qs.dbInstance(ctx, opt, Options(opts...)).Delete(qs.model)
	affectedRow, err = cmd.RowsAffected, cmd.Error
	return
}

func (qs *queries[M, Q]) GetOneTo(ctx context.Context, result Model, opts ...Option) error {
	return qs.dbInstance(ctx, opts...).Find(result).Error
}

func (qs *queries[M, Q]) GetManyTo(ctx context.Context, resultList any, opts ...Option) error {
	return qs.dbInstance(ctx, opts...).Find(resultList).Error
}

func (qs *queries[M, Q]) GetManyWithCount(ctx context.Context, r ManyWithCountResulter, opts ...Option) error {
	var (
		ListOnly     = r.QsList() != nil && r.QsCount() == nil
		CountOnly    = r.QsList() == nil && r.QsCount() != nil
		ListAndCount = r.QsList() != nil && r.QsCount() != nil
	)

	switch {

	case CountOnly:
		total, err := qs.Count(ctx, Options(opts...), WithoutLimitAndOffset())
		if err != nil {
			return err
		}
		*r.QsCount() = total
		return nil

	case ListOnly:
		return qs.GetManyTo(ctx, r.QsList(), WithModel(qs.model), Options(opts...))

	case ListAndCount:
		return qs.GetManyTo(ctx, r.QsList(), Options(opts...), Count(r.QsCount(), WithModel(qs.model)))
	}

	return errors.New("not support operation, one of resp or count must not nil")
}
