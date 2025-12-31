/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 登录访问日志
 */

package sys

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"bailu/app/domain/dto"
	"bailu/app/domain/entity"
	"bailu/app/domain/repo"
	"bailu/app/domain/repo/base"
	"bailu/app/domain/resp"
	"bailu/utils"
	"bailu/utils/types"
	"time"
)

type LoginLogService struct {
	Repo *repo.LoginLogRepo
}

func NewLoginLogService(repo *repo.LoginLogRepo) *LoginLogService {
	return &LoginLogService{repo}
}

func (l *LoginLogService) List(ctx context.Context, params dto.LoginLogParams) (*resp.PageResult[entity.LoginInfo], error) {
	builder := base.NewQueryBuilder()
	builder.WithWhereStruct(params).
		WithPagination(ctx)

	if result, err := l.Repo.FindByBuilder(ctx, builder); err != nil {
		return nil, err
	} else {
		pageRecord := result.(*resp.PageResult[entity.LoginInfo])
		return pageRecord, nil
	}
}

func (l *LoginLogService) FindByName(ctx context.Context, userName string) (*entity.LoginInfo, error) {
	builder := base.NewQueryBuilder()
	builder.WithWhere("username=?", userName)
	var loginInfo = &entity.LoginInfo{}
	if err := l.Repo.FindModelByBuilder(ctx, builder, loginInfo); err != nil {
		return nil, err
	} else {
		return loginInfo, nil
	}
}

//func (l *LoginLogService) Create(ctx context.Context, log entity.LoginInfo) (entity.LoginInfo, error) {
//	err := l.Repo.Create(ctx, &log)
//	return log, err
//}

func (l *LoginLogService) Create(c *gin.Context, username string, status int, msg string) error {
	var log = entity.LoginInfo{Username: username, Status: status}
	ua := user_agent.New(c.Request.UserAgent())
	//name, version := ua.Engine()
	//log.Browser = name + " " + version
	browser, v := ua.Browser()
	log.Browser = browser + " " + v
	log.Os = ua.OS()
	log.Ip = c.ClientIP()
	log.Addr = utils.GetAddr(log.Ip)
	log.LoginTime = types.JSONTime{time.Now()}
	log.Msg = msg
	err := l.Repo.Create(c.Request.Context(), &log)
	return err
}

func (l *LoginLogService) Destroy(ctx context.Context, ids []uint64) error {
	return l.Repo.Delete(ctx, ids)
}

// 清空
func (l *LoginLogService) Clean(ctx context.Context) error {
	return l.Repo.Truncate(ctx)
}
