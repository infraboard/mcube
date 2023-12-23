package datasource

import (
	"context"

	"gorm.io/gorm"
)

type TransactionCtxKey struct{}

func WithTransactionCtx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, TransactionCtxKey{}, tx)
}

// 开启事务
func BeginTransaction(ctx context.Context) *gorm.DB {
	return DB(ctx).Begin().WithContext(ctx)
}

// 结束事务
func EndTransaction(tx *gorm.DB, err error) error {
	if err != nil {
		return tx.Rollback().Error
	}

	return tx.Commit().Error
}

func GetTransactionFromCtx(ctx context.Context) *gorm.DB {
	if ctx != nil {
		tx, ok := ctx.Value(TransactionCtxKey{}).(*gorm.DB)
		if ok {
			return tx
		}
	}

	return nil
}
