package tools

import (
	"context"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/stores"
)

// AfterFunc 是一个用于处理事务后执行函数的工具结构体，支持在事务执行完毕后运行特定的函数，包括事务内函数和异步函数。
type AfterFunc struct {
	txFuncs []func(db *stores.DB) error //业务处理完之后需要在事务内处理
	funcs   []func(ctx context.Context) //执行完之后异步处理,事务实行完之后处理就行
}

func (a *AfterFunc) AddFunc(f func(ctx context.Context)) {
	a.funcs = append(a.funcs, f)
}

func (a *AfterFunc) AddTxFunc(f func(tx *stores.DB) error) {
	a.txFuncs = append(a.txFuncs, f)
}

func (a *AfterFunc) Handle(ctx context.Context, handle func(tx *stores.DB) error) error {
	err := stores.GetTenantConn(ctx).Transaction(func(tx *stores.DB) error {
		err := handle(tx)
		if err != nil {
			return err
		}
		for _, f := range a.txFuncs {
			err := f(tx)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, f := range a.funcs {
		ctxs.GoNewCtx(ctx, func(ctx context.Context) {
			f(ctx)
		})
	}
	return nil
}
