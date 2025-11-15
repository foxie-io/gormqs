package gormqs

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type (
	contextKey string
)

func ContextWithValue(ctx context.Context, values ...any) context.Context {
	for _, val := range values {
		key := contextKey(fmt.Sprintf("%T", val))
		ctx = context.WithValue(ctx, key, val)
	}
	return ctx
}

func ContextValue[T any](ctx context.Context, fallback T) T {
	var (
		zero T
		key  = contextKey(fmt.Sprintf("%T", zero))
	)
	casted, ok := ctx.Value(key).(T)
	if !ok {
		return fallback
	}

	return casted
}

// safeTextForSql sanitizes the input text by replacing unsafe characters for SQL queries.
func SafeTextForSql(text string) string {
	// Replace wildcard asterisk (*) with SQL's LIKE wildcard (%)
	text = strings.Split(text, ";")[0]

	// Remove potentially harmful characters (e.g., semicolon, hyphen)
	text = strings.Map(func(r rune) rune {
		switch r {
		case '*':
			return '%'
		case '/', '\\', '[', ']', '(', ')', ';', '`', '|', '&', '^', '%', '$', '#', '@', '!', '<', '>', '?', ',', '+', '-':
			return rune(0)
		}
		return r
	}, text)

	text = strings.ReplaceAll(text, string(rune(0)), "")
	return text
}

func WithTable(col string, db *gorm.DB) string {
	if db == nil {
		return col
	}

	if db.Statement.Table == "" {
		return col
	}

	if strings.Contains(col, ".") {
		return col
	}

	return fmt.Sprintf("`%s`.`%s`", db.Statement.Table, col)
}
