/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc base service
 */
package base

import (
	"context"
	"github.com/sufo/bailu-admin/app/config"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/repo/base"
	"github.com/sufo/bailu-admin/global/consts"
	respErr "github.com/sufo/bailu-admin/pkg/exception"
	"net/http"
)

//func NewBaseService[T entity.IModel](db *gorm.DB) *BaseService[T] {
//	r := repo.Repository[T]{db}
//	return &BaseService[T]{r}
//}

type BaseService[T entity.IModel] struct {
	Repo base.Repository[T]
}

// 检查字段唯一性
func (b *BaseService[T]) CheckUnique(ctx context.Context, query string, args ...any) bool {
	var count int64
	b.Repo.Where(ctx, query, args...).Count(&count)
	return count == 0
}

func (b *BaseService[T]) ContainSuper(ctx context.Context, userIds []uint64) bool {
	var users []*entity.User
	err := b.Repo.Where(ctx, "id in ?", userIds).Preload("Roles").Find(&users).Error
	if err != nil {
		panic(respErr.InternalServerError)
	}
	for _, u := range users {
		for _, r := range u.Roles {
			if consts.SUPER_ROLE_ID == r.ID {
				return true
			}
		}
	}
	return false
}

func ReqSchema(r *http.Request) string {
	// Check for the X-Forwarded-Proto header first, which is the standard for reverse proxies.
	if proto := r.Header.Get("X-Forwarded-Proto"); proto == "https" {
		return "https"
	}
	// Fallback for direct TLS connections.
	if r.TLS != nil {
		return "https"
	}
	// Default to http.
	return "http"
}

// 获取文件url
func FileUrl(r *http.Request, path string) string {
	var uploadConf = config.Conf.Upload
	if uploadConf.Model == "local" {
		return ReqSchema(r) + "://" + r.Host + config.Conf.Local.Path + "/" + path
	} else {
		return path
	}
}
