package qsopt

import (
	"fmt"

	"github.com/foxie-io/gormqs"
	"gorm.io/gorm"
)

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

func (s UserSchema) Select(cols ...UserColumn) gormqs.Option {
	if len(cols) == 0 {
		return s.Select("*")
	}

	return func(db *gorm.DB) *gorm.DB {
		columns := make([]string, len(cols))
		for i, col := range cols {
			columns[i] = gormqs.WithTable(string(col), db)
		}

		return db.Select(columns)
	}
}

func (s UserSchema) WhereID(id uint) gormqs.Option {
	return s.Where(s.ID, "=", id)
}
