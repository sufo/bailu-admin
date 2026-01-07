/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package repo

import (
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/repo/base"
	"gorm.io/gorm"
)

type LoginLogRepo struct {
	base.Repository[entity.LoginInfo]
}

func NewLoginLogRepo(db *gorm.DB) *LoginLogRepo {
	return &LoginLogRepo{base.Repository[entity.LoginInfo]{db}}
}
