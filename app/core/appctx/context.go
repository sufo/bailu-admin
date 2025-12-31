/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package appctx

import (
	"context"
	"gorm.io/gorm"
)

type (
	transCtx     struct{}
	noTransCtx   struct{}
	transLockCtx struct{}
	userIDCtx    struct{}
	traceIDCtx   struct{}
	tagCtx       struct{}
	userCtx      struct{}
	pageCtx      struct{}
)

var ctx = context.WithValue(context.Background(), struct{}{}, "APP")

// Wrap transaction context
func NewTrans(ctx context.Context, trans interface{}) context.Context {
	return context.WithValue(ctx, transCtx{}, trans)
}

func FromTrans(ctx context.Context) (*gorm.DB, bool) {
	//v := ctx.Value(transCtx{}).(T)
	//return v, &v != nil

	//v := ctx.Value(transCtx{}).(*gorm.DB)
	v := ctx.Value(transCtx{})
	if v == nil {
		return nil, false
	} else {
		return v.(*gorm.DB), true
	}
}

func NewNoTrans(ctx context.Context) context.Context {
	return context.WithValue(ctx, noTransCtx{}, true)
}

func FromNoTrans(ctx context.Context) bool {
	v := ctx.Value(noTransCtx{})
	return v != nil && v.(bool)
}

func NewTransLock(ctx context.Context) context.Context {
	return context.WithValue(ctx, transLockCtx{}, true)
}

func FromTransLock(ctx context.Context) bool {
	v := ctx.Value(transLockCtx{})
	return v != nil && v.(bool)
}

func SetAuth(ctx context.Context, user interface{}) context.Context {
	return context.WithValue(ctx, userCtx{}, user)
}

func GetAuthUser[T any](ctx context.Context) *T {
	v := ctx.Value(userCtx{})
	user, ok := v.(*T)
	if ok {
		return user
	} else {
		return nil
	}
}

func NewTagCtx(ctx context.Context, tag string) context.Context {
	return context.WithValue(ctx, tagCtx{}, tag)
}
func FromTagContext(ctx context.Context) string {
	v := ctx.Value(tagCtx{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func NewPageCtx(ctx context.Context, val interface{}) context.Context {
	return context.WithValue(ctx, pageCtx{}, val)
}
func GetPageCtx[T any](ctx context.Context) (*T, bool) {
	v := ctx.Value(pageCtx{})
	if v != nil {
		page, ok := v.(*T)
		return page, ok
	} else {
		return nil, false
	}
}
