/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package repo

import (
	"bailu/app/domain/entity"
	"bailu/app/domain/repo/base"
	"gorm.io/gorm"
)

type LoginLogRepo struct {
	base.Repository[entity.LoginInfo]
}

func NewLoginLogRepo(db *gorm.DB) *LoginLogRepo {
	return &LoginLogRepo{base.Repository[entity.LoginInfo]{db}}
}
