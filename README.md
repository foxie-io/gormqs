# Gorm Queries (gormqs)

`gormqs` is a lightweight and extendable wrapper for Gorm, designed to simplify query building while maintaining type safety. It strikes a balance between the flexibility of Gorm and the strict type safety of libraries like Ent.

---

## Features

- **Simplified Queries**: Perform basic queries without unnecessary complexity.
- **Extendable Options**: Customize queries with reusable and type-safe options.
- **Type Safety**: Minimize the use of `interface{}` or `any` for better code reliability.
- **Transaction Support**: Easily integrate with Gorm transactions.

---

## Why Use `gormqs`?

### Why not just use [Gorm](https://github.com/go-gorm/gorm)?

Gorm is a powerful ORM, but it lacks type safety, which can lead to runtime errors.

### Why not [Ent](https://github.com/ent/ent)?

Ent provides type safety but can introduce significant boilerplate code.

### Why `gormqs`?

`gormqs` offers a middle ground, combining the best of both worlds:

- Simplifies query building.
- Reduces boilerplate code.
- Enhances type safety.

---

## Getting Started

### Installation

Add `gormqs` to your project:

```bash
go get github.com/foxie-io/gormqs
```

Ensure you have Gorm installed as well:

```bash
go get gorm.io/gorm
```

---

### Basic Usage

#### Define Your Model

```go
// models/user.go
package models

type User struct {
	ID       uint    `json:"id"`
	Username string  `json:"username"`
	Balance  float64 `json:"balance"`
}

func (User) TableName() string {
	return "users"
}
```

#### Create Queries

```go
// queries/user.go
package queries

import (
	"context"
	"gorm.io/gorm"
	"github.com/foxie-io/gormqs"
	"your_project/models"
)

type UserQueries struct {
	gormqs.Queries[models.User, *UserQueries]
	db    *gorm.DB
	model models.User
}

func (qr *UserQueries) DBInstance(ctx context.Context) *gorm.DB {
	db := gormqs.ContextValue(ctx, qr.db)
	return db.WithContext(ctx).Table(qr.model.TableName()).Model(qr.model)
}

func NewUserQueries(db *gorm.DB) *UserQueries {
	qs := &UserQueries{db: db}
	qs.Queries = gormqs.NewQueries[models.User](qs)
	return qs
}
```

#### Use Queries

```go
userQueries := queries.NewUserQueries(db)

// Fetch a user by username
user, err := userQueries.GetOne(ctx, qopt.USER.Where(qopt.USER.Username, "=", "example"))
if err != nil {
	log.Fatal(err)
}

fmt.Println("Fetched User:", user)
```

---

## Advanced Features

### Custom Query Options

Define reusable and type-safe query options:

```go
// queries/options/user.go
package qopt

type UserColumn string

type UserSchema struct {
	ID       UserColumn
	Username UserColumn
	Balance  UserColumn
}

var USER = UserSchema{
	ID:       "id",
	Username: "username",
	Balance:  "balance",
}

func (s UserSchema) Where(col UserColumn, operation, value any) gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		query := fmt.Sprintf("%s %s ?", gormqs.WithTable(string(col), db), operation)
		return db.Where(query, value)
	}
}
```

### Transactions

Integrate with Gorm transactions:

```go
db.Transaction(func(tx *gorm.DB) error {
	ctx := gormqs.WrapContext(tx)

	user, err := userQueries.GetOne(ctx, gormqs.LockForUpdate(), qopt.USER.WhereID(1))
	if err != nil {
		return err
	}

	user.Balance += 100
	if _, err := userQueries.Update(ctx, user, qopt.USER.Select(qopt.USER.Balance)); err != nil {
		return err
	}

	return nil
})
```

---

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

---

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
