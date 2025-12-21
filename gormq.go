package gormqs

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
type ListOrCountResulter interface {
	QsList() any
	QsCount() *int64
}

var (
	// ListResulter implement interface ListOrCountResulter
	_ ListOrCountResulter = (*ListResulter[Model])(nil)
)

type ListResulter[T Model] struct {
	selects string
	List    *[]*T  `json:"list,omitempty"`
	Count   *int64 `json:"count,omitempty"`
}

/*
ListResulter

selects := "list,count" | "list" | "count"
*/
func NewListResulter[T Model](selects string) *ListResulter[T] {
	return &ListResulter[T]{
		selects: selects,
	}
}

func (r *ListResulter[T]) QsList() any {
	if strings.Contains(r.selects, "list") {
		r.List = new([]*T)
		return r.List
	}
	return nil
}

func (r *ListResulter[T]) QsCount() *int64 {
	if strings.Contains(r.selects, "count") {
		r.Count = new(int64)
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

	/*
		Update

			// after update return value will mutate user
			var user models.User
			user.money = 100
			qs.Update(ctx, &user)

			// after update return value will mutate updatedUser
			var (
				userNewValue models.User
				updatedUser models.User
			)
			user.money = 100
			qs.Update(ctx, &userNewValue, WithModel(&updatedUser) // after update return value update to user
	*/
	Update(ctx context.Context, record *M, opt Option, opts ...Option) (affectedRow int64, err error)

	/*update with expr

	var user models.User
	// user id=3, balacne =100, money = 100

	qs.UpdateWithExpr(ctx,
		map[string]clause.Expr{
			"balance": gorm.Expr("balance + ?", 100),
			"money":   gorm.Expr("money - ?", 100),
		},
		WithModel(&user),
	) // SQL: Update "users" SET "balance" = "balance" + 100, "money" = "money" - 100 WHERE "users"."id" = 3

	log.Println(user) // user.balance = 200, user.money = 0


	qs.UpdateWithExpr(ctx,
		map[string]clause.Expr{
			"balance": gorm.Expr("balance + ?", 20),
			"money":   gorm.Expr("money - ?", 20),
		},
		Where("id = ?", 2),
	) // SQL: UPDATE "users" SET "balance" = "balance" + 20, "money" = "money" - 20 WHERE "users"."id" = 2
	*/
	UpdateWithExpr(ctx context.Context, values map[string]clause.Expr, opt Option, opts ...Option) (affectedRow int64, err error)

	/*Count need atleast one option

	total,err := qs.Count(ctx, Where(id > 0)) // SQL: SELECT count(*) FROM "users" WHERE "users"."id" > 0
	*/
	Count(ctx context.Context, opt Option, opts ...Option) (count int64, err error)

	/*Delete need atleast one option

	count, err := qs.Delete(ctx, Where(id > 0)) // SQL: DELETE FROM "users" WHERE "users"."id" > 0
	*/
	Delete(ctx context.Context, opt Option, opts ...Option) (affectedRow int64, err error)

	// scan pattern for custom type without mapping to struct again

	/*get one to struct directly

	type BasicUser struct {
		ID   int
		Name string
	}
	var user BasicUser
	err := qs.GetOneTo(ctx, &user) // SQL: SELECT "users"."id", "users"."name" FROM "users" LIMIT 1

	type User {
		ID int
		Name string
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	var basicUser User
	err := qs.GetOneTo(ctx, &basicUser) // SQL: SELECT "users"."id", "users"."name", "users"."created_at", "users"."updated_at" FROM "users" LIMIT 1
	*/
	GetOneTo(ctx context.Context, result Model, opts ...Option) error

	/*get many to struct directly


	type BasicUser struct {
		ID   int
		Name string
	}
	var basicUsers []BasicUser
	err := qs.GetOneTo(ctx, &basicUsers) // SQL: SELECT "users"."id", "users"."name" FROM "users"


	type User {
		ID int
		Name string
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	var users []User
	err := qs.GetOneTo(ctx, &basicUsers) // SQL: SELECT "users"."id", "users"."name", "users"."created_at", "users"."updated_at" FROM "users"
	*/
	GetManyTo(ctx context.Context, resultList any, opts ...Option) error

	/* get many ony, count only or both


	var resulter NewManyWithCountResulter[User]("many,count") // list and count
	err := qs.GetList(ctx, resulter)


	var resulter NewListResulter[User]("count") // list only
	err := qs.GetList(ctx, resulter)

	var resulter NewListResulter[User]("list") // list only
	err := qs.GetList(ctx, resulter)
	*/
	GetListTo(ctx context.Context, listResulter ListOrCountResulter, options ...Option) error
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
	err := qs.dbInstance(ctx, opts...).Model(&result).First(&result).Error
	return &result, err
}

func (qs *queries[M, Q]) GetMany(ctx context.Context, opts ...Option) ([]*M, error) {
	var result []*M
	err := qs.dbInstance(ctx, opts...).Find(&result).Error
	return result, err
}

func (qs *queries[M, Q]) Update(ctx context.Context, record *M, opt Option, opts ...Option) (affectedRow int64, err error) {
	defaultOpt := Options(WithModel(record), opt)
	cmd := qs.dbInstance(ctx, defaultOpt, Options(opts...)).Updates(record)
	affectedRow, err = cmd.RowsAffected, cmd.Error
	return
}

func (qs *queries[M, Q]) UpdateWithExpr(ctx context.Context, values map[string]clause.Expr, opt Option, opts ...Option) (affectedRow int64, err error) {
	cmd := qs.dbInstance(ctx, opt, Options(opts...)).Updates(values)
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

func (qs *queries[M, Q]) GetOneTo(ctx context.Context, r Model, opts ...Option) error {
	return qs.dbInstance(ctx, opts...).First(r).Error
}

func (qs *queries[M, Q]) GetManyTo(ctx context.Context, rList any, opts ...Option) error {
	return qs.dbInstance(ctx, opts...).Find(rList).Error
}

func (qs *queries[M, Q]) GetListTo(ctx context.Context, r ListOrCountResulter, opts ...Option) error {
	switch {
	case r.QsList() != nil && r.QsCount() != nil:

		return qs.GetManyTo(ctx, r.QsList(), Options(opts...), Count(r.QsCount(), WithModel(qs.model)))

	case r.QsList() != nil:
		return qs.GetManyTo(ctx, r.QsList(), WithModel(qs.model), Options(opts...))

	case r.QsCount() != nil:
		total, err := qs.Count(ctx, Options(opts...), WithoutLimitAndOffset())
		if err != nil {
			return err
		}
		*r.QsCount() = total
		return nil

	}

	return errors.New("not support operation, one of resp or count must not nil")
}
