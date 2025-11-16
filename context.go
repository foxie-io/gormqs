package gormqs

import (
	"context"

	"gorm.io/gorm"
)

// ReplaceContext replace gorm.DB instance in context
func ReplaceContext(tx *gorm.DB) context.Context {
	tx.Statement.Context = ContextWithValue(tx.Statement.Context, tx)
	return tx.Statement.Context
}

// WrapContext wrap gorm.DB instance in context
func WrapContext(db *gorm.DB) context.Context {
	exist := ContextValue[*gorm.DB](db.Statement.Context, nil)
	if exist != nil {
		return exist.Statement.Context
	}

	return ReplaceContext(db)
}
