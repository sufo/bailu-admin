/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 菜单（按钮）=> api
 */

package repo

import (
	"bailu/app/domain/entity"
	"bailu/app/domain/repo/base"
	"context"
	"gorm.io/gorm"
)

func NewMenuApiRepo(db *gorm.DB) *MenuApiRepo {
	return &MenuApiRepo{base.Repository[entity.MenuApi]{db}}
}

type MenuApiRepo struct {
	base.Repository[entity.MenuApi]
}

// 硬删除
func (m *MenuApiRepo) Delete(ctx context.Context, menuId uint64) error {
	return m.Where(ctx, "menu_id=?", menuId).Delete(&entity.MenuApi{}).Error
}
