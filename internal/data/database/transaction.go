package database

import (
	"context"
	"flow-bridge-mcp/internal/biz"
	"flow-bridge-mcp/internal/data/model"
	"gorm.io/gorm"
)

type transaction struct {
	data *Data
}

func NewTransaction(d *Data) biz.Transaction {
	return &transaction{
		data: d,
	}
}

// ExecTx 该方法是为了biz层使用事务定义的，普通事务无需使用该方法
// ctx: 上下文对象，用于传递请求作用域的数据和控制请求生命周期
// fn: 在事务中执行的函数，接收上下文对象作为参数，返回错误信息
// 返回值: 执行过程中发生的错误，如果执行成功则返回nil
func (t *transaction) ExecTx(ctx context.Context, fn func(ctx context.Context) error) (err error) {

	// 在GORM事务中执行函数，确保数据库操作的原子性
	err = t.data.Db.Transaction(func(tx *gorm.DB) error {
		// 将事务对象存储到上下文中，以便在fn函数中可以通过上下文获取事务对象进行数据库操作
		ctx = context.WithValue(ctx, model.TxKey, tx)
		return fn(ctx)

	})
	return
}
