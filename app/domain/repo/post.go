/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 岗位
 */

package repo

import (
	"bailu/app/domain/entity"
	"bailu/app/domain/repo/base"
	"gorm.io/gorm"
)

// var PostSet = wire.NewSet(wire.Struct(new(PostRepo), "*"))
func NewPostRepo(db *gorm.DB) *PostRepo {
	r := base.Repository[entity.Post]{db}
	return &PostRepo{r}
}

type PostRepo struct {
	base.Repository[entity.Post]
}
