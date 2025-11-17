# Gorm Queries

Simple Gorm Queries Wrapper

## Features

- basic queries without bloated your
- extendable queries and option (customize your need in queries or option)

## Question

- why not just use [Gorm](https://github.com/go-gorm/gorm)?
  - Gorm is a good ORM, but it lacks type safety
- why not [Ent](https://github.com/ent/ent)?
  - Ent is type safe, but it bloats your code
- why [GormQs](https://github.com/foxie-io/gormqs)?
  - It is in the middle between Gorm and Ent
  - Reduce usage of go type interface{} | any of `Gorm`

## Queries

- basic queries interface so, it won't have to much black magic
- support with options (use to build query like `where id =`, `select col,col2`, `preload`, `join`)

```go
type Queries[M Model, Q any] interface {
	CreateOne(ctx context.Context, record *M) error
	CreateMany(ctx context.Context, record *[]*M) error

	GetOne(ctx context.Context, opts ...Option) (result *M, err error)
	GetMany(ctx context.Context, opts ...Option) (result []*M, err error)

	// update one or many is base on at least one opt

	Updates(ctx context.Context, record *M, opt Option, opts ...Option) (affectedRow int64, err error)
	Count(ctx context.Context, opt Option, opts ...Option) (count int64, err error)
	Delete(ctx context.Context, opt Option, opts ...Option) (affectedRow int64, err error)

	// scan pattern for custom type, and dynamic select without mapping to struct again
	GetOneTo(ctx context.Context, result Model, opts ...Option) error
	GetManyTo(ctx context.Context, resultList any, opts ...Option) error

	// list with count | list ony | count ony with same option meaning same query filter
	GetManyWithCount(ctx context.Context, r ManyWithCountResulter, options ...Option) error
}
```

### Model

```go
// models/user.go
package models

type User struct {
	...
	Username string  `json:"username"`
	Balance  float64 `json:"balance"`
}

func (User) TableName() string {
	return "users"
}
```

### Extended `queries`

```go
// queries/user.go
package queries

type UserQueries struct {
	gormqs.Queries[models.User, *UserQueries] // basic queries, GetOne,GetMany, etc...
	db    *gorm.DB
	model models.User
}

// provider db instance for gormqs.Queries to use
func (qr *UserQueries) DBInstance(ctx context.Context) *gorm.DB {
	db := gormqs.ContextValue(ctx, qr.db) // use instance from context for transaction
	return db.WithContext(ctx).Table(qr.model.TableName()).Model(qr.model)
}

func NewUserQueries(db *gorm.DB) *UserQueries {
	qs := &UserQueries{db: db}
	qs.Queries = gormqs.NewQueries[models.User](qs)
	return qs
}

func (qr *UserQueries) UpdateUserByUsername(ctx context.Context, username string, newVal any) error {
	// implement your own query
}
```

### Implement `queries option` (type safe and reusable)

- recommend to write your own `codegen` for minimum options to avoid code bloated

```go
// queries/options/user.go
package qopt

type UserColumn string

type UserSchema struct {
	ID             UserColumn
	CreatedAt      UserColumn
	UpdatedAt      UserColumn
	Username       UserColumn
	Balance        UserColumn
	BlockedBalance UserColumn
}

var USER = UserSchema{
	ID:             "id",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	Username:       "username",
	Balance:        "balance",
	BlockedBalance: "blocked_balance",
}

func (s UserSchema) Where(col UserColumn, operation, value any) gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		query := fmt.Sprintf("%s %s ?", gormqs.WithTable(string(col), db), operation)
		return db.Where(query, value)
	}
}

func (s UserSchema) WhereID(id uint) gormqs.Option {
	return s.Where(s.ID, "=", id)
}

func (s UserSchema) SelectAll() gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select("*")
	}
}

func (s UserSchema) Select(cols ...UserColumn) gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		columns := make([]string, len(cols))
		for i, col := range cols {
			columns[i] = gormqs.WithTable(string(col), db)
		}
		return db.Select(columns)
	}
}


```

### Usage with transaction

```go
user_qs := queries.NewUserQueries(db)
user1, _ := user_qs.GetOne(ctx,  qopt.Where(qopt.USER.Username,"=","user1"))

db.Transaction(func(tx *gorm.DB) error {
	ctx := gormqs.WrapContext(tx)

	user, err := user_qs.GetOne(ctx, gormqs.LockForUpdate(), qopt.USER.WhereID(user1.ID))
	if err != nil {
		return err
	}

	user.Balance += 100
	// update only balance
	if _, err := user_qs.Updates(ctx, user, qopt.USER.Select(qopt.USER.Balance)); err != nil {
		return err
	}

    // update all columns on user
	if _, err := user_qs.Updates(ctx, user, qopt.USER.SelectAll); err != nil {
		return err
	}

	// seperate result and value
	var result models.User
	if _, err := user_qs.Updates(ctx, user, qopt.USER.SelectAll, gormq.WithModel(&result)); err != nil {
		return err
	}

	// use gorm directly with option
	if err := tx.Scopes(gormq.WithModel(&user), qopt.USER.SelectAll).Updates(user).Err;err != nil {
		return err
	}

	return nil
})

```
