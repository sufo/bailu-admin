/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 基础类
 */

package util

import (
	"github.com/sufo/bailu-admin/app/core/appctx"
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 这里是否考虑使用session
// Get gorm.DB from context
func GetDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	tx, ok := appctx.FromTrans(ctx)
	if ok && !appctx.FromNoTrans(ctx) {
		if appctx.FromTransLock(ctx) {
			tx = tx.Clauses(clause.Locking{Strength: "UPDATE"})
		}
		return tx
	}
	//return defDB
	//钩子函数获取context,tx.Statement.Context
	return defDB.WithContext(ctx)
}

// Get gorm.DB.Model from context
func GetDBWithModel[T any](ctx context.Context, defDB *gorm.DB) *gorm.DB {
	var t T
	return GetDB(ctx, defDB).Model(&t)
}

// table=="" => GetDBWithModel
func GetDBWithTable[T any](ctx context.Context, defDB *gorm.DB, table string) *gorm.DB {
	if table == "" {
		return GetDB(ctx, defDB).Table(table)
	} else {
		return GetDBWithModel[T](ctx, defDB)
	}
}

// Define transaction execute function
// import cycle not allowed
//type TransFunc func(context.Context) error
//
//func ExecTrans(ctx context.Context, db *gorm.DB, fn TransFunc) error {
//	transModel := &repo.Trans{DB: db}
//	return transModel.Exec(ctx, fn)
//}
//
//func ExecTransWithLock(ctx context.Context, db *gorm.DB, fn TransFunc) error {
//	if !appctx.FromTransLock(ctx) {
//		ctx = appctx.NewTransLock(ctx)
//	}
//	return ExecTrans(ctx, db, fn)
//}
