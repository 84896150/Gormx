package gormx

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type DataBaseSoftDelete int8

const (
	Normal DataBaseSoftDelete = iota + 1
	Deleted
)

type ModelBaseInfo struct {
	CreateAt time.Time          `gorm:"column:create_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_at"`                             // 创建时间
	UpdateAt time.Time          `gorm:"column:update_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;NOT NULL" json:"update_at"` // 最后修改时间
	Deleted  DataBaseSoftDelete `gorm:"column:deleted;default:1;NOT NULL" json:"deleted"`                                                 // 是否删除 1-正常 2-删除
}

type Transaction interface {
	InTx(context.Context, func(ctx context.Context) error) error
}

// Data .
type Data struct {
	GORMDB *gorm.DB
}

type contextTxKey struct{}

// DB 和事务相关的db操作，在取db连接时均采用此方法
func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.GORMDB.WithContext(ctx)
}

// InTx fn是包含了事务操作的方法，只要fn里面有异常，里面的db操作都会回滚
func (d *Data) InTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.GORMDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

// NewTransaction .
func NewTransaction(d *Data) Transaction {
	return d
}
