package transactions

import (
	"context"

	"gorm.io/gorm"
)

type DatabaseTransaction interface {
	RunTransactions(ctx context.Context, repos ...RepositoryHandleFunc) error
}

type RepositoryHandleFunc func(context.Context, *gorm.DB) error

type GormTxDB struct {
	db *gorm.DB
}

func New(db *gorm.DB) DatabaseTransaction {
	return &GormTxDB{
		db: db,
	}
}

func (t *GormTxDB) RunTransactions(ctx context.Context, repos ...RepositoryHandleFunc) error {
	tx := t.db.Begin()
	if err := tx.Error; err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, repo := range repos {
		if err := repo(ctx, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
