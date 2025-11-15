package qsopt

import (
	"github.com/foxie-io/gormqs"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// commond usage

func Limit(limit int) gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

func Page(page int, limit int) gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit).Offset(page + 1)
	}
}

func Count(result *int64) gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Count(result)
	}
}

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

		return db.Clauses(clause.Locking{})
	}
}
