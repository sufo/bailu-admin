/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 在线用户
 */

package sys

import (
	"context"
	"encoding/json"
	"github.com/google/wire"
	"github.com/jinzhu/copier"
	"github.com/mssola/user_agent"
	gUtil "gorm.io/gorm/utils"
	"bailu/app/config"
	"bailu/app/core/appctx"
	"bailu/app/domain/entity"
	"bailu/app/domain/resp"
	respErr "bailu/pkg/exception"
	"bailu/pkg/jwt"
	"bailu/pkg/store"
	"bailu/utils"
	"bailu/utils/page"
	"bailu/utils/types"
	"net/http"
	"sort"
	"strings"
	"time"
)

var OnlineSet = wire.NewSet(wire.Struct(new(OnlineService), "*"))

type OnlineService struct {
	DB store.IStore
}

// 获取用户信息
func (online *OnlineService) GetOne(key string) (string, error) {
	return online.DB.Get(key)
}

func (online *OnlineService) GetAll(ctx context.Context) ([]*entity.OnlineUserDto, error) {
	res, err := online.DB.Scan(config.Conf.JWT.OnlineKey)
	if err != nil {
		return nil, err
	}
	users := make([]*entity.OnlineUserDto, 0)
	for _, item := range res {
		var u *entity.OnlineUserDto
		if err := json.Unmarshal([]byte(item.V), u); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (online *OnlineService) List(ctx context.Context, userName string, addr string) (*resp.PageResult[entity.OnlineUserDto], error) {
	res, err := online.DB.Scan(config.Conf.JWT.OnlineKey)
	if err != nil {
		return nil, err
	}
	users := make([]*entity.OnlineUserDto, 0)
	//条件筛选
	for _, item := range res {
		var u = &entity.OnlineUserDto{}
		if err := json.Unmarshal([]byte(item.V), u); err != nil {
			return nil, err
		}
		var flag = true
		if userName != "" {
			flag = strings.Contains(u.Username, userName)
			if !flag {
				continue
			}
		}
		if addr != "" {
			flag = strings.Contains(u.Addr, addr)
			if !flag {
				continue
			}
		}

		if flag {
			//列表不需要角色信息,这里清除
			u.Roles = make([]entity.Role, 0)

			users = append(users, u)
		}
	}
	//排序
	sort.SliceStable(users, func(i, j int) bool { return users[i].LoginTime.After(users[i].LoginTime.Time) })

	var pageIndex, pageSize int
	if p, exist := appctx.GetPageCtx[page.Pagination](ctx); exist {
		pageIndex = p.PageIndex
		pageSize = p.PageSize
	} else {
		panic(respErr.BadRequestErrorWithMsg("必须传分页参数"))
	}
	//分页
	var pageResult = resp.PageResult[entity.OnlineUserDto]{
		PageSize:  pageSize,
		PageIndex: pageIndex,
		ItemCount: int64(len(users)),
	}
	var size = len(users)
	var start = (pageIndex - 1) * pageSize
	if size <= start {
		l := make([]*entity.OnlineUserDto, 0)
		pageResult.List = l
	} else {
		var end = pageIndex * pageSize
		if size < pageIndex*pageSize {
			end = size
		}
		pageResult.List = users[(pageIndex-1)*pageSize : end]
	}

	return &pageResult, nil
}

// save online user
// token: userKey
func (online *OnlineService) Save(user *entity.User, request *http.Request, token string) (*entity.OnlineUserDto, error) {
	onlineUserDto := entity.OnlineUserDto{}
	err := copier.Copy(&onlineUserDto, *user)
	if err != nil {
		return nil, err
	}
	ua := user_agent.New(request.UserAgent())
	//name, version := ua.Engine()
	//onlineUserDto.Browser = name + " " + version
	browser, v := ua.Browser()
	onlineUserDto.Browser = browser + " " + v

	onlineUserDto.Os = ua.OS()
	onlineUserDto.Addr = utils.GetAddr(user.Ip)
	onlineUserDto.Token = token
	onlineUserDto.LoginTime = types.JSONTime{time.Now()}
	//var roles []uint64
	//for i, e := range user.Roles {
	//	roles[i] = e.ID
	//}
	//onlineUserDto.RoleIds = roles
	conf := config.Conf.JWT
	userKey := jwt.UserKey(request.UserAgent(), gUtil.ToString(user.ID))
	return &onlineUserDto, online.DB.Set(conf.OnlineKey+userKey, onlineUserDto, time.Duration(conf.Expired)*time.Second)
}

/**
 * 踢出用户
 * @param key /
 */
//func (online *OnlineService) KickOut(key string) error {
//	res, err := online.DB.Find(config.Conf.JWT.OnlineKey, key)
//	if err != nil {
//		return err
//	}
//	return online.DB.Del(res.K)
//}

func (online *OnlineService) KickOut(key string) error {
	var onlineKey = config.Conf.JWT.OnlineKey + key
	exist, err := online.DB.Check(onlineKey)
	if err != nil {
		return err
	}
	if exist {
		return online.DB.Del(onlineKey)
	}
	return nil
}

func (online *OnlineService) BatchKickOut(request *http.Request, keys []string) error {
	isM := utils.IsMobile(request.UserAgent())
	for _, key := range keys {
		userKey := jwt.UserKeyBy(isM, key)
		if err := online.DB.Del(config.Conf.JWT.OnlineKey + userKey); err != nil {
			return err
		}
	}
	return nil
}
