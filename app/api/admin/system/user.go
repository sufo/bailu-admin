/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 用户管理
 */

package system

import (
	"github.com/sufo/bailu-admin/app/api/admin"
	"github.com/sufo/bailu-admin/app/config"
	"github.com/sufo/bailu-admin/app/domain/dto"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/resp"
	"github.com/sufo/bailu-admin/app/domain/resp/status"
	"github.com/sufo/bailu-admin/app/domain/vo"
	"github.com/sufo/bailu-admin/app/service/sys"
	"github.com/sufo/bailu-admin/global/consts"
	respErr "github.com/sufo/bailu-admin/pkg/exception"
	"github.com/sufo/bailu-admin/pkg/i18n"
	"github.com/sufo/bailu-admin/pkg/jwt"
	"github.com/sufo/bailu-admin/pkg/log"
	"github.com/sufo/bailu-admin/pkg/rsa"
	"github.com/sufo/bailu-admin/pkg/store"
	"github.com/sufo/bailu-admin/utils"
	"github.com/sufo/bailu-admin/utils/page"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

var UserSet = wire.NewSet(wire.Struct(new(UserApi), "*"))

type UserApi struct {
	UserSrv *sys.UserService
	Store   store.IStore
}

// @title 用户列表
// @Summary 用户列表接口
// @Description 按条件查询用户列表
// @Tags User
// @Accept json
// @Produce json
// @Param query query dto.UserQueryParams false "查询条件"
// //@Success 200 {object} resp.Response{data=resp.PageResult[entity.User]}
// @Success 200 {object} resp.Response[resp.PageResult[entity.User]]
// @Security Bearer
// @Router /api/user [get]
func (u *UserApi) Index(c *gin.Context) {
	var queryParams dto.UserQueryParams
	if err := c.ShouldBind(&queryParams); err != nil {
		panic(respErr.BadRequestError)
	}

	page.StartPage(c)
	result, err := u.UserSrv.List(c.Request.Context(), queryParams)
	if err != nil {
		panic(respErr.InternalServerErrorWithMsg(err.Error()))
	}
	resp.OKWithData(c, result)
}

// @title 当前用户信息
// @Summary 获取当前用户信息接口
// @Description 获取当前用户信息
// @Tags User
// @Accept json
// @Produce json
// //@Success 200 {object} resp.Response{data=vo.User}
// @Success 200 {object} resp.Response[vo.User]
// @Router /api/user/info [get]
// @Security Bearer
func (u *UserApi) GetInfo(c *gin.Context) {
	if onlineUser, exist := c.Get(consts.REQUEST_USER); exist {
		user, err := u.UserSrv.FindUserInfoById(c.Request, onlineUser.(*entity.OnlineUserDto).ID)
		if err != nil {
			panic(respErr.InternalServerErrorWithMsg(err.Error()))
		}
		//处理userVo roles
		var userVo = &vo.User{}
		var roles = make([]vo.KV, 0)
		copier.Copy(userVo, user)
		for _, r := range user.Roles {
			op := vo.KV{r.RoleKey, r.Name}
			roles = append(roles, op)
		}
		userVo.Roles = roles
		//post
		var posts = make([]vo.KV, 0)
		for _, p := range user.Posts {
			op := vo.KV{p.PostCode, p.Name}
			posts = append(posts, op)
		}
		userVo.Posts = posts

		//response
		resp.OKWithData(c, userVo)
	} else {
		resp.FailWithMsg(c, i18n.DefTr("admin.userNotExit"))
	}
}

// @title 创建用户
// @Summary 创建用户接口
// @Description 创建用户接口
// @Tags User
// @Accept json
// @Produce json
// @Param body body entity.User true "用户信息"
// @Success 200 {object} resp.Response[any]
// @Security Bearer
// @Router /api/user [post]
func (u *UserApi) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var user dto.UserDto
	if err := c.ShouldBind(&user); err != nil {
		panic(respErr.BadRequestError)
	}

	//检查用户名
	isUnique := u.UserSrv.CheckUnique(ctx, "username=?", user.Username)
	if !isUnique {
		resp.FailWithMsg(c, i18n.DefTr("admin.existed", user.Username))
		return
	}
	//检查手机号
	if len(user.Phone) > 0 {
		isUnique := u.UserSrv.CheckUnique(ctx, "phone=?", user.Phone)
		if !isUnique {
			resp.FailWithMsg(c, i18n.DefTr("admin.existed", user.Username))
			return
		}
	}
	//检察邮箱
	if len(user.Email) > 0 {
		isUnique := u.UserSrv.CheckUnique(ctx, "email=?", user.Email)
		if !isUnique {
			resp.FailWithMsg(c, i18n.DefTr("admin.existed", user.Email))
			return
		}
	}

	err := u.UserSrv.Create(ctx, user)
	if err != nil {
		resp.FailWithError(c, err)
		return
	}
	resp.Ok(c)
}

// @title 修改用户
// @Summary 修改用户接口
// @Description 修改用户接口
// @Tags User
// @Accept json
// @Produce json
// @Param body body entity.User true "用户信息"
// @Success 200 {object} resp.Response[any]
// @Security Bearer
// @Router /api/user [put]
func (u *UserApi) Edit(c *gin.Context) {
	ctx := c.Request.Context()
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		panic(respErr.BadRequestError)
	}
	//检查手机号
	if len(user.Phone) > 0 {
		isUnique := u.UserSrv.CheckUnique(ctx, "phone=? and id!=?", user.Phone, user.ID)
		if !isUnique {
			resp.FailWithMsg(c, i18n.DefTr("admin.existed", user.Username))
			return
		}
	}
	//检察邮箱
	if len(user.Email) > 0 {
		isUnique := u.UserSrv.CheckUnique(ctx, "email=? and id!=?", user.Email, user.ID)
		if !isUnique {
			resp.FailWithMsg(c, i18n.DefTr("admin.existed", user.Email))
			return
		}
	}
	err := u.UserSrv.Update(ctx, user)
	if err != nil {
		resp.FailWithError(c, err)
		return
	}
	resp.Ok(c)
}

// @title 批量删除用户
// @Summary 批量删除用户接口
// @Description 批量删除用户接口
// @Tags User
// @Accept json
// @Produce json
// @Param userIds path string true "用户id集合"
// @Success 200 {object} resp.Response[any]
// @Security Bearer
// @Router /api/user/{userIds} [delete]
func (u *UserApi) Destroy(c *gin.Context) {
	v := admin.ParseParamIDs(c, "userIds")
	logigUser := admin.GetLoginUser(c)
	if utils.Includes(v, logigUser.ID) {
		resp.FailWithMsg(c, "当前用户不能删除")
		return
	}
	err := u.UserSrv.Delete(c.Request.Context(), v)
	if err != nil {
		resp.FailWithError(c, err)
		return
	}
	resp.Ok(c)
}

//func (u *UserApi) _isExist(c *gin.Context, paramName string) {
//	name := c.Param(paramName)
//	if name == "" {
//		resp.FailWithError(c, respErr.BadRequestError)
//		return
//	}
//	//得到数据库字段名
//	fieldName := utils.Camel2Case(paramName)
//	if isExist, err := u.UserSrv.UserRepo.IsExist(c, fieldName+"=?", name); err != nil {
//		resp.FailWithError(c, respErr.InternalServerError)
//		return
//	} else {
//		var exist = 0
//		if isExist {
//			exist = 1
//		}
//		//data := map[string]int{
//		//	"exist": exist,
//		//}
//		resp.Result(c, resp.Data(exist), resp.Msg(i18n.DefTr("admin.existed", name)))
//	}
//}

func (u *UserApi) _isExist(c *gin.Context, paramName ...string) {
	vals := make([]string, len(paramName))
	fieldNames := make([]string, len(paramName))
	var query = ""
	for index, value := range paramName {
		//获取路径参数
		vals[index] = c.Param(value)
		if vals[index] == "" {
			panic(respErr.BadRequestError)
		}
		//得到数据库字段名
		fieldNames[index] = utils.Camel2Case(value)

		query += fieldNames[index] + "=? and"
	}
	query = query[0 : len(query)-4]

	if isExist, err := u.UserSrv.UserRepo.IsExist(c.Request.Context(), query, vals); err != nil {
		resp.FailWithError(c, respErr.InternalServerError)
		return
	} else {
		var exist = 0
		if isExist {
			exist = 1
		}
		resp.Result(c, resp.Data(exist), resp.Msg(i18n.DefTr("admin.existed", vals[len(vals)-1])))
	}
}

// 用户名是否存在
// @title 用户名是否存在
// @Summary 用户名是否存在
// @Description 用户名是否存在
// @Tags User
// @Accept json
// @Produce json
// @Param username path string true "username"
// @Success 200 {object} resp.Response[any]
// @Security Bearer
// @Router /api/username/{username} [get]
func (u *UserApi) NameExist(c *gin.Context) {
	u._isExist(c, "username")
}

// 手机号是否存在
// @title 手机号是否存在
// @Summary 手机号是否存在
// @Description 手机号是否存在
// @Tags User
// @Accept json
// @Produce json
// @Param dialCode path string true "国际电话区号"
// @Param phone path string true "手机号"
// @Success 200 {object} resp.Response[any]
// @Security Bearer
// @Router /api/phone/{dialCode}/{phone} [get]
func (u *UserApi) PhoneExist(c *gin.Context) {
	u._isExist(c, "dialCode", "phone")
}

// 用户注册
// 注册需要短信验证码
// @title 用户注册
// @Summary 用户注册接口
// @Description 用户注册，需要短信验证码
// @Tags User
// @Accept json
// @Produce json
// @Param body body dto.RegisterParams true "注册信息"
// @Success 200 {object} resp.Response[any]
// @Security Bearer
// @Router /api/register [post]
func (u *UserApi) Register(c *gin.Context) {
	var r dto.RegisterParams
	if err := c.ShouldBindJSON(&r); err != nil {
		panic(respErr.BadRequestErrorWithMsg(status.StatusText(status.StatusBadRequest)))
	}

	isDebug := config.Conf.Server.Mode == "debug"

	//验证码处理
	if isDebug {
		if r.SMSCode != "123456" {
			resp.FailWithMsg(c, i18n.DefTr("tip.smsCodeInvalid"))
		}
	} else {
		code, err := u.Store.Get(r.Phone)
		if err != nil || r.SMSCode != code {
			resp.FailWithMsg(c, i18n.DefTr("tip.smsCodeInvalid"))
			return
		}
	}

	//处理密码
	decryptStr, err := rsa.PrivateDecrypt(r.Password, config.Conf.RSA.PrivateKey)
	if err != nil {
		resp.FailWithMsg(c, i18n.DefTr("admin.AccountOrPwdErr"))
		return
	}
	//密码强度校验
	if !utils.PasswordStrength(decryptStr) {
		resp.FailWithMsg(c, i18n.DefTr("tip.pwdFormatTip", 6, 18, 2))
		return
	}

	//检查数据库是否存在重复

	//加密密码(哈希) 采用bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(decryptStr), bcrypt.DefaultCost)
	var user = dto.UserDto{
		Username: r.Username,
		Password: string(hash),
		Phone:    r.Phone,
	}
	if err := u.UserSrv.Register(c.Request.Context(), user); err != nil {
		log.L.Error(err)
		//resp.InternalServerError(c)
		panic(err)
		return
	}
	fmt.Sprintf("%v", user)
	resp.Ok(c)
}

// 通过管理员来直接重置用户密码 需要校验权限
// @title 管理员重置用户密码
// @Summary 通过管理员来直接重置用户密码
// @Description 通过管理员来直接重置用户密码 需要校验权限
// @Tags User
// @Accept json
// @Produce json
// @Param body body dto.SetPwdParams true "请求参数"
// @Success 200 {object} resp.Response[any]
// @Security Bearer
// @Router /api/user/resetPassword [put]
func (u *UserApi) ResetPwd(c *gin.Context) {
	var param dto.SetPwdParams

	if err := c.ShouldBindJSON(&param); err != nil {
		panic(respErr.BadRequestErrorWithMsg(status.StatusText(status.StatusBadRequest)))
	}
	if err := u.UserSrv.Update(c.Request.Context(), entity.User{ID: param.Id, Password: param.Password}); err != nil {
		log.L.Error(err)
		resp.InternalServerError(c)
		return
	}
	resp.Ok(c)
}

// @title 启用/禁用用户
// @Summary 启用/禁用用户接口
// @Description 启用/禁用用户
// @Tags User
// @Accept json
// @Produce json
// @Param body body object{id=integer,status=integer} true "id和status"
// @Success 200 {object} resp.Response[any]
// @Security Bearer
// @Router /api/user/status [patch]
func (u *UserApi) Status(c *gin.Context) {

	var params = make(map[string]string)
	if err := c.ShouldBindJSON(&params); err != nil {
		panic(respErr.BadRequestError)
	}
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		panic(respErr.BadRequestError)
	}
	status, err := strconv.ParseUint(params["status"], 10, 8)
	if err != nil {
		panic(respErr.BadRequestError)
	}
	err = u.UserSrv.ChangeStatus(c.Request.Context(), id, uint8(status))
	if err != nil {
		log.L.Error(err)
		resp.InternalServerError(c)
		return
	}
	resp.Ok(c)
}

// @title 短信验证修改用户密码
// @Summary 短信验证修改用户密码接口
// @Description 通过手机短信验证码修改用户密码
// @Tags User
// @Accept json
// @Produce json
// @Param body body dto.ResetParams true "请求参数"
// @Success 200 {object} resp.Response[any]
// @Security Bearer
// @Router /api/user/resetPasswordBySmsCode [put]
func (u *UserApi) ResetPwdBySMSCode(c *gin.Context) {
	var r dto.ResetParams
	if err := c.ShouldBindJSON(&r); err != nil {
		panic(respErr.BadRequestErrorWithMsg(status.StatusText(status.StatusBadRequest)))
	}

	isDebug := config.Conf.Server.Mode == "debug"

	//验证码处理
	if isDebug {
		if r.SMSCode != "123456" {
			resp.FailWithMsg(c, i18n.DefTr("tip.smsCodeInvalid"))
		}
	} else {
		code, err := u.Store.Get(r.Phone)
		if err != nil || r.SMSCode != code {
			resp.FailWithMsg(c, i18n.DefTr("tip.smsCodeInvalid"))
			return
		}
	}
	if err := u.UserSrv.ResetPasswordByPhone(c.Request.Context(), r.Phone, r.Password); err != nil {
		log.L.Error(err)
		resp.InternalServerError(c)
		return
	}
	//使token失效
	token, _ := c.Get(consts.REQ_TOKEN)
	u.Store.Del(config.Conf.JWT.OnlineKey + token.(string))

	resp.Ok(c)
}

// @title 退出登录
// @Summary 退出登录接口
// @Description 退出当前登录帐号
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} resp.Response[any]
// @Security Bearer
// @Router /api/user/logout [post]
func (u *UserApi) Logout(c *gin.Context) {
	//token, _ := c.Get(consts.REQ_TOKEN)
	//err := u.Store.Del(config.Conf.JWT.OnlineKey + token.(string))
	user := admin.GetLoginUser(c)
	userKey := jwt.UserKey(c.Request.UserAgent(), utils.Strval(user.ID))
	err := u.Store.Del(config.Conf.JWT.OnlineKey + userKey)
	if err != nil {
		resp.InternalServerError(c)
	} else {
		resp.Ok(c)
	}
}
