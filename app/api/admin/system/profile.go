/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 个人信息处理
 */

package system

import (
	"bailu/app/api/admin"
	"bailu/app/domain/dto"
	"bailu/app/domain/entity"
	"bailu/app/domain/resp"
	"bailu/app/domain/resp/status"
	"bailu/app/domain/vo"
	"bailu/app/service/sys"
	"bailu/global/consts"
	respErr "bailu/pkg/exception"
	"bailu/pkg/i18n"
	"bailu/pkg/log"
	"bailu/utils/upload"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/jinzhu/copier"
)

var ProfileSet = wire.NewSet(wire.Struct(new(ProfileApi), "*"))

type ProfileApi struct {
	UserSrv *sys.UserService
	//OnlineSrv *service.OnlineService
	Oss upload.OSS //已经是指针
}

// @title 修改个人信息
// @Summary 修改个人信息接口
// @Description 修改个人信息接口
// @Tags Profile
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body entity.User true "个人信息"
// @Success 200 {object} resp.Response[any]
// @Router /api/user/profile [put]
func (p *ProfileApi) Edit(c *gin.Context) {
	ctx := c.Request.Context()
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		panic(respErr.BadRequestError)
	}
	//检查手机号
	if len(user.Phone) > 0 {
		isUnique := p.UserSrv.CheckUnique(ctx, "phone=? and id!=?", user.Phone, user.ID)
		if !isUnique {
			resp.FailWithMsg(c, i18n.DefTr("admin.existed", user.Username))
			return
		}
	}
	//检察邮箱
	if len(user.Email) > 0 {
		isUnique := p.UserSrv.CheckUnique(ctx, "email=? and id!=?", user.Email, user.ID)
		if !isUnique {
			resp.FailWithMsg(c, i18n.DefTr("admin.existed", user.Email))
			return
		}
	}
	err := p.UserSrv.Update(ctx, user)
	if err != nil {
		resp.FailWithError(c, err)
		return
	}

	//返回User信息
	info, err := p.UserSrv.FindUserInfoById(c.Request, user.ID)
	if err != nil {
		panic(respErr.InternalServerErrorWithMsg(err.Error()))
	}
	//处理userVo roles
	var userVo = &vo.User{}
	var roles = make([]vo.KV, 0)
	copier.Copy(userVo, info)
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
}

// 通过旧密码来修改密码
// @title 修改密码
// @Summary 修改密码接口
// @Description 通过旧密码来设置新密码
// @Tags Profile
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body dto.ChangePwdParams true "用户id、原密码、新密码"
// @Success 200 {object} resp.Response[any]
// @Router /api/user/changePwd [put]
func (p *ProfileApi) ChangePwd(c *gin.Context) {
	var param dto.ChangePwdParams
	if err := c.ShouldBindJSON(&param); err != nil {
		panic(respErr.BadRequestErrorWithMsg(status.StatusText(status.StatusBadRequest)))
	}

	if err := p.UserSrv.ChangePassword(c.Request.Context(), param); err != nil {
		log.L.Error(err)
		resp.InternalServerError(c)
		return
	}
}

// avtar upload
// @title 上传头像
// @Summary 上传头像接口
// @Description 上传头像接口
// @Tags Profile
// @Accept json
// @Produce json
// @Security Bearer
// @Param file formData file true "file"
// @Success 200 {object} resp.Response[any]
// @Router /api/user/avatar [put]
func (p *ProfileApi) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		respErr.BadRequestErrorWithError(err)
	}
	fPath, _, err := p.Oss.UploadFileToDir(file, consts.IMG_DIR)
	if err != nil {
		respErr.InternalServerErrorWithError(err)
	}
	onlineUser := admin.GetLoginUser(c)

	//删除之前文件
	user, err := p.UserSrv.FindById(c.Request.Context(), onlineUser.ID)
	if err != nil {
		respErr.InternalServerErrorWithError(err)
	}
	err = p.Oss.DeleteFileInDir(user.Avatar, consts.IMG_DIR)
	if err != nil {
		respErr.InternalServerErrorWithError(err)
	}

	err = p.UserSrv.UserRepo.UpdateColumn(c.Request.Context(), onlineUser.ID, "avatar", fPath).Error
	if err != nil {
		respErr.InternalServerErrorWithError(err)
	}
	r := c.Request
	url := admin.FileUrl(r, fPath)
	resp.OKWithData(c, map[string]string{
		"url": url,
	})
}
