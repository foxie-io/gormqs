package queries

import (
	"context"
	"example/userorder/models"
	qsopt "example/userorder/queries/options"

	"github.com/foxie-io/gormqs"
	"gorm.io/gorm"
)

var (
	_ gormqs.Model   = (*models.User)(nil)
	_ gormqs.Querier = (*UserQueries)(nil)
)

type (
	UserQueries struct {
		gormqs.Queries[models.User, *UserQueries]
		db    *gorm.DB
		model models.User
	}
)

// provider db instance for gormqs.Queries to use
func (qs *UserQueries) DBInstance(ctx context.Context) *gorm.DB {
	db := gormqs.ContextValue(ctx, qs.db)
	return db.WithContext(ctx).Table(qs.model.TableName()).Model(qs.model)
}

func NewUserQueries(db *gorm.DB) *UserQueries {
	qs := &UserQueries{db: db}
	qs.Queries = gormqs.NewQueries[models.User](qs)
	return qs
}

func (qs *UserQueries) LockForUpdate(ctx context.Context, userId uint, updateUser func(u models.User) models.User, updateColumns ...qsopt.UserColumn) (returnUser *models.User, returnErr error) {
	returnErr = qs.DBInstance(ctx).Transaction(func(tx *gorm.DB) error {
		ctx := gormqs.ContextWithValue(tx.Statement.Context, tx)
		user, err := qs.GetOne(ctx, gormqs.LockForUpdate(), qsopt.USER.WhereID(userId))
		if err != nil {
			return err
		}

		newValue := updateUser(*user)
		newValue.ID = user.ID

		_, err = qs.Updates(ctx, &newValue, qsopt.USER.Select(updateColumns...))
		returnUser = &newValue
		return err
	})

	return
}

func (qs *UserQueries) BlockBalance(ctx context.Context, userId uint, blockingAmount float64) (returnUser *models.User, returnErr error) {
	return qs.LockForUpdate(ctx, userId, func(u models.User) models.User {
		return models.User{
			Balance:        u.Balance - blockingAmount,
			BlockedBalance: u.BlockedBalance + blockingAmount,
		}
	},
		qsopt.USER.Balance,
		qsopt.USER.BlockedBalance,
	)
}

func (qs *UserQueries) UnblockBalance(ctx context.Context, userId uint, blockingAmount float64) (returnUser *models.User, returnErr error) {
	return qs.LockForUpdate(ctx, userId, func(u models.User) models.User {
		return models.User{
			Balance:        u.Balance + blockingAmount,
			BlockedBalance: u.BlockedBalance - blockingAmount,
		}
	},
		qsopt.USER.Balance,
		qsopt.USER.BlockedBalance,
	)
}

func (qs *UserQueries) CommitBlockedBalance(ctx context.Context, userId uint, blockingAmount float64) (returnUser *models.User, returnErr error) {
	return qs.LockForUpdate(ctx, userId, func(u models.User) models.User {
		return models.User{
			Balance:        u.Balance + blockingAmount,
			BlockedBalance: u.BlockedBalance - blockingAmount,
		}
	},
		qsopt.USER.BlockedBalance,
	)
}
