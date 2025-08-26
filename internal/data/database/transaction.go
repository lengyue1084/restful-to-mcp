package database

import (
	"context"
	"flow-bridge-mcp/internal/biz"
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
func (t *transaction) ExecTx(ctx context.Context, fn func(ctx context.Context) error) (err error) {

	err = t.data.Db.Transaction(func(tx *gorm.DB) error {
		// 将事务对象放入 context，在data层操作数据库时，使用 context.Value("tx") 获取事务对象，不然会出现tx对象不一致
		ctx = context.WithValue(ctx, "tx", tx)
		// 执行业务逻辑
		return fn(ctx)
		// 如果 fn 返回 error，GORM 自动回滚
		// 如果 fn 返回 nil，GORM 自动提交
	})
	return
}
