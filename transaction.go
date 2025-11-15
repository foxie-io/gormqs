package gormqs

import (
	"context"

	"gorm.io/gorm"
)

type TxHandler func(*gorm.DB) error

// merge multiple handlers
func MergeTx(handlers ...TxHandler) TxHandler {
	return func(d *gorm.DB) error {
		for _, handler := range handlers {
			if err := handler(d); err != nil {
				return err
			}
		}
		return nil
	}
}

// replace context use in transaction
func TxCtx(newCtx context.Context) TxHandler {
	return func(d *gorm.DB) error {
		d.Statement.Context = newCtx
		return nil
	}
}

// easier to seperate logic into smaller functions
func Tx(handlers ...TxHandler) func(*gorm.DB) error {
	return func(tx *gorm.DB) error {
		tx.Statement.Context = ContextWithValue(tx.Statement.Context, tx)

		for _, handler := range handlers {
			if err := handler(tx); err != nil {
				return err
			}
		}

		return nil
	}
}
