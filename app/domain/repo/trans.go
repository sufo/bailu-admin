/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc  事务
 */

package repo

import (
	"context"
	"github.com/google/wire"
	"gorm.io/gorm"
	"bailu/app/core/appctx"
)

var TransSet = wire.NewSet(wire.Struct(new(Trans), "*"))

type TransFunc func(context.Context) error

type Trans struct {
	DB *gorm.DB
}

func (a *Trans) Exec(ctx context.Context, fn func(context.Context) error) error {
	if _, ok := appctx.FromTrans(ctx); ok {
		return fn(ctx)
	}

	return a.DB.Transaction(func(db *gorm.DB) error {
		//return fn(appctx.NewTrans(ctx, db))
		return fn(appctx.NewTrans(ctx, db.WithContext(ctx))) //db中放入ctx
	})
}

func (a *Trans) ExecTransWithLock(ctx context.Context, db *gorm.DB, fn TransFunc) error {
	if !appctx.FromTransLock(ctx) {
		ctx = appctx.NewTransLock(ctx)
	}
	return a.Exec(ctx, fn)
}
