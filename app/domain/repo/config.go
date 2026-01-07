/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 系统参数配置
 */

package repo

import (
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/repo/base"
	"gorm.io/gorm"
)

func NewSysConfigRepo(db *gorm.DB) *SysConfigRepo {
	r := base.Repository[entity.SysConfig]{db}
	return &SysConfigRepo{r}
}

type SysConfigRepo struct {
	base.Repository[entity.SysConfig]
}
