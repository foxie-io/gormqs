# gormqs

Simple Gorm Queries prod

## Features

- Simple interface for gorm queries
- Use default queries if not provide
- Support custom queries
- Support option for queries

## Queries

```
type Queries[M Model, Q any] interface {
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
```

## Usage

### Declare Model

```
// models/user_model.go

type User struct {
	...
	Username string  `json:"username"`
	Balance  float64 `json:"balance"`
}

func (User) TableName() string {
	return "users"
}

```

### Declare UserQueries

```
// queries/user_query.go

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
```

### Declare Your common option

```
// queries/options/option.go

func WhereID(id uint) gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func LockForUpdate(mode ...clause.Locking) gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		if len(mode) > 0 {
			return db.Clauses(mode[0])
		}
		return db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
}
```

### Declare your own option for type safe

```
// queries/options/user_option.go

type UserColumn string

const (
	UserID        UserColumn = "id"
	UserCreatedAt UserColumn = "created_at"
	UserUpdatedAt UserColumn = "updated_at"
	UserUsername  UserColumn = "username"
	UserBalance   UserColumn = "balance"
)

func UserWhere(col UserColumn, operation, value any) gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		query := fmt.Sprintf("%s %s ?", gormqs.WithTable(string(col), db), operation)
		return db.Where(query, value)
	}
}

func UserSelect(cols ...UserColumn) gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		for _, col := range cols {
			db.Statement.Selects = append(db.Statement.Selects, gormqs.WithTable(string(col), db))
		}
		return db
	}
}

```

### Usage with or without transaction

```
user_qs := queries.NewUserQueries(db)
user1, err := user_qs.GetOne(ctx,  qsopt.WhereID(user1.ID))

err := db.Transaction(func(tx *gorm.DB) error {
		ctx := gormqs.ContextWithValue(tx.Statement.Context, tx)

		user, err := user_qs.GetOne(ctx, qsopt.LockForUpdate(), qsopt.WhereID(user1.ID))
		if err != nil {
			return err
		}

		user.Balance += 100
		if _, err := user_qs.Updates(ctx, user, qsopt.UserSelect(qsopt.UserBalance)); err != nil {
			return err
		}
		return nil
	})

```
