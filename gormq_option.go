package gormqs

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

func WhereID(id any) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func Select(query interface{}, args ...interface{}) Option {
	return func(q *gorm.DB) *gorm.DB {
		return q.Select(query, args...)
	}
}

func WithDebug() Option {
	return func(q *gorm.DB) *gorm.DB {
		return q.Session(&gorm.Session{
			Logger: logger.Default,
		}).Debug()
	}
}

func Omit(columns ...string) Option {
	return func(q *gorm.DB) *gorm.DB {
		if len(columns) == 0 {
			return q
		}

		return q.Omit(columns...)
	}
}

func HardDelete() Option {
	return func(q *gorm.DB) *gorm.DB {
		return q.Unscoped()
	}
}

func WithoutLimitAndOffset() Option {
	return func(q *gorm.DB) *gorm.DB {
		return q.Offset(-1).Limit(-1)
	}
}

func LimitAndOffset(limit int, offset int) Option {
	return func(q *gorm.DB) *gorm.DB {
		return q.Limit(limit).Offset(offset)
	}
}

func Count(count *int64, countOpts ...Option) Option {
	return func(q *gorm.DB) *gorm.DB {
		newSession := q.Session(&gorm.Session{Initialized: true})
		for _, opt := range countOpts {
			opt(newSession)
		}
		WithoutLimitAndOffset()(newSession).Count(count)
		return q
	}
}

func LockForUpdate(lock ...clause.Locking) Option {
	return func(q *gorm.DB) *gorm.DB {
		if len(lock) > 0 {
			return q.Clauses(lock[0])
		}
		return q.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate})
	}
}

func If(condition bool, opt Option) Option {
	return func(q *gorm.DB) *gorm.DB {
		if condition {
			q = opt(q)
		}
		return q
	}
}

func IfF(conditionFn func() bool, opt Option) Option {
	return If(conditionFn(), opt)
}

func WithModel(model interface{}) Option {
	return func(q *gorm.DB) *gorm.DB {
		return q.Model(model)
	}
}
