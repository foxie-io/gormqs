package qsopt

import (
	"fmt"

	"github.com/foxie-io/gormqs"
	"gorm.io/gorm"
)

type UserColumn string

const (
	UserID       UserColumn = "id"
	UserUsername UserColumn = "username"
	UserBalance  UserColumn = "balance"
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
