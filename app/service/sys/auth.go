/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package sys

import (
	"bailu/app/domain/dto"
	"bailu/app/domain/entity"
	"bailu/app/domain/repo"
	respErr "bailu/pkg/exception"
	"bailu/pkg/i18n"
	"bailu/pkg/log"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var AuthSet = wire.NewSet(wire.Struct(new(AuthService), "*"))

type AuthService struct {
	UserRepo *repo.UserRepo
	MenuSrv  *MenuService
}

func (a *AuthService) Login(c *gin.Context, params dto.LoginUser) (*entity.User, error) {

	user, err := a.UserRepo.FindByName(c.Request.Context(), params.Username)
	if err == nil {
		if user.Status != 1 {
			return nil, respErr.UserDisableError
		} else if checkErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); checkErr != nil {
			return nil, respErr.WrapLogicResp(i18n.DefTr("admin.pwdIncorrect"))
		}

		ip := c.ClientIP()
		//更新ip和登录时间
		if err := a.UserRepo.Updates(c.Request.Context(), user.ID, []string{"ip", "last_login_time"}, ip, time.Now()); err != nil {
			log.L.Warn(err)
			return nil, err
		}
		user.Ip = ip
		err = a.withPermissions(c.Request.Context(), user)
		return user, err
	}
	return nil, err
}

// 处理相关权限
func (a *AuthService) withPermissions(ctx context.Context, user *entity.User) error {
	if perms, err := a.MenuSrv.FindPermissions(ctx, user); err != nil {
		return err
	} else {
		user.Permissions = perms
		return nil
	}
}
