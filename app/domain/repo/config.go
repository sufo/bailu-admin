/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 系统参数配置
 */

package repo

import (
	"bailu/app/domain/entity"
	"bailu/app/domain/repo/base"
	"gorm.io/gorm"
)

func NewSysConfigRepo(db *gorm.DB) *SysConfigRepo {
	r := base.Repository[entity.SysConfig]{db}
	return &SysConfigRepo{r}
}

type SysConfigRepo struct {
	base.Repository[entity.SysConfig]
}
