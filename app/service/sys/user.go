/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package sys

import (
	"bailu/app/config"
	"bailu/app/core/appctx"
	"bailu/app/domain/dto"
	"bailu/app/domain/entity"
	"bailu/app/domain/repo"
	"bailu/app/domain/repo/base"
	"bailu/app/domain/resp"
	base2 "bailu/app/service/base"
	"bailu/global/consts"
	respErr "bailu/pkg/exception"
	"bailu/pkg/i18n"
	"bailu/pkg/log"
	"bailu/pkg/rsa"
	"bailu/pkg/sms"
	"bailu/pkg/store"
	"bailu/utils"
	"context"
	"github.com/google/wire"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"
)

var UserSet = wire.NewSet(wire.Struct(new(UserOption), "*"), NewUserService)

type UserService struct {
	base2.BaseService[entity.User]
	UserOption
}

type UserOption struct {
	UserRepo *repo.UserRepo
	TranRepo *repo.Trans
	MenuRepo *repo.MenuRepo
	Sms      *sms.AliyunClient
	Store    store.IStore
}

func NewUserService(opt UserOption) *UserService {
	return &UserService{base2.BaseService[entity.User]{opt.UserRepo.Repository}, opt}
}

func (u *UserService) List(ctx context.Context, params dto.UserQueryParams) (*resp.PageResult[entity.User], error) {
	builder := base.NewQueryBuilder()
	builder.WithTable("sys_user u").
		WithJoin("left join sys_dept as d on d.id=u.dept_id").
		WithWhereStructAndAlias(params, "sys_dept")

	if params.DeptId != "" {
		did, err := strconv.ParseUint(params.DeptId, 10, 64)
		if err != nil {
			panic(respErr.BadRequestErrorWithError(err))
		}
		builder.WithWhere("(u.dept_id=? or u.dept_id in (select sys_dept.id from sys_dept where find_in_set(?, ancestors)))", did, did)
	}
	builder.WithPreload("Roles").
		WithPreload("Posts").
		WithDataScope(ctx, "d", "u").
		WithPagination(ctx)
	return u.UserRepo.ListByBuilder(ctx, builder)
}

// 新增用户
func (u *UserService) Create(ctx context.Context, user dto.UserDto) error {
	user.Password = u.decryptAndHash(user.Password)
	return u.TranRepo.Exec(ctx, func(ctx context.Context) error {
		tx, _ := appctx.FromTrans(ctx)
		err := tx.Table(entity.UserTN).Create(&user).Error
		if err != nil {
			return err
		}
		//插入用户角色
		if user.RoleIds != nil && len(user.RoleIds) > 0 {
			var userRoles = make([]map[string]any, 0)
			for _, rid := range user.RoleIds {
				var userRole = make(map[string]any)
				userRole["user_id"] = user.ID
				userRole["role_id"] = rid
				userRoles = append(userRoles, userRole)
			}
			err := tx.Table("sys_user_role").Create(&userRoles).Error
			if err != nil {
				return err
			}
		}
		//插入岗位
		if user.PostIds != nil && len(user.PostIds) > 0 {
			var userPosts = make([]map[string]any, 0, len(user.PostIds))
			for _, rid := range user.PostIds {
				userPost := map[string]any{"user_id": user.ID, "post_id": rid}
				userPosts = append(userPosts, userPost)
			}
			err := tx.Table("sys_user_post").Create(&userPosts).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// 用户注册
func (u *UserService) Register(ctx context.Context, user dto.UserDto) error {
	//user.Password = u.CheckPwdAndHash(user.Password)
	return u.UserRepo.Table(ctx, entity.UserTN).Create(&user).Error
}

// 删除用户与角色关联
func (u *UserService) DelUserRoleByUserId(ctx context.Context, userIds []uint64) error {
	var userRole entity.UserRole
	return u.UserRepo.GetDB(ctx).Where("user_id in ?", userIds).Delete(&userRole).Error
}

// 删除用户与岗位关联
func (u *UserService) DelUserPostByUserId(ctx context.Context, userIds []uint64) error {
	var userPost entity.UserPost
	return u.UserRepo.GetDB(ctx).Where("user_id in ?", userIds).Delete(&userPost).Error
}

func (u *UserService) FindById(ctx context.Context, userId uint64) (*entity.User, error) {
	var user entity.User
	err := u.UserRepo.Where(ctx, "id=?", userId).
		Preload("Roles").Preload("Posts").Find(&user).Error
	return &user, err
}

// 获取用户信息（含permission）
func (u *UserService) FindUserInfoById(req *http.Request, userId uint64) (*entity.User, error) {
	var ctx = req.Context()
	var user entity.User
	err := u.UserRepo.Where(ctx, "sys_user.id=?", userId).
		Joins("left join sys_dept d on d.id=sys_user.id").
		Select("sys_user.*,d.name as deptName").
		Preload("Roles").Preload("Posts").Find(&user).Error
	if err != nil {
		return nil, err
	}

	//处理相关权限
	var permissions []string
	//超管
	if user.IsSuper() {
		permissions = append(permissions, consts.SUPER_PERMISSION)
	} else {
		roleIds := make([]uint64, len(user.Roles))
		for _, role := range user.Roles {
			roleIds = append(roleIds, role.ID)
		}
		var menus []*entity.Menu
		if menus, err = u.MenuRepo.FindByRoleIds(ctx, roleIds); err == nil {
			for _, menu := range menus {
				permissions = append(permissions, *menu.Permission)
			}
		}
	}
	user.Permissions = permissions
	if user.Avatar != "" && !strings.HasPrefix(user.Avatar, "http") {
		user.Avatar = base2.FileUrl(req, user.Avatar)
	}
	return &user, err
}

func (u *UserService) Update(ctx context.Context, user entity.User) error {
	//注意修改用户没有提供修改密码框
	return u.TranRepo.Exec(ctx, func(ctx context.Context) error {
		tx, _ := appctx.FromTrans(ctx)
		err := tx.Where("id=?", user.ID).Updates(&user).Error
		if err != nil {
			return err
		}

		//插入用户角色
		if user.RoleIds != nil { //为nil，表示角色没有发生修改，如果是清空角色则要传空数组
			//删除之前用户和角色关联
			if err := tx.Where("user_id=?", user.ID).Unscoped().Delete(&entity.UserRole{}).Error; err != nil {
				return err
			}
			if len(user.RoleIds) > 0 {
				var userRoles = make([]map[string]any, 0)
				for _, rid := range user.RoleIds {
					var userRole = make(map[string]any)
					userRole["user_id"] = user.ID
					userRole["role_id"] = rid
					userRoles = append(userRoles, userRole)
				}
				err := tx.Table("sys_user_role").Create(&userRoles).Error
				if err != nil {
					return err
				}
			}
		}

		//插入岗位
		if user.PostIds != nil {
			//删除之前用户所在岗位
			if err := tx.Where("user_id=?", user.ID).Unscoped().Delete(&entity.UserPost{}).Error; err != nil {
				return err
			}
			if len(user.PostIds) > 0 {
				var userPosts = make([]map[string]any, 0, len(user.PostIds))
				for _, rid := range user.PostIds {
					userPost := map[string]any{"user_id": user.ID, "post_id": rid}
					userPosts = append(userPosts, userPost)
				}
				err := tx.Table("sys_user_post").Create(&userPosts).Error
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (u *UserService) Delete(ctx context.Context, ids []uint64) error {
	//检查用户能否被操作
	for _, id := range ids {
		CheckUserAllowed(entity.User{ID: id})
	}
	//检查数据权限
	u.checkUserDataScope(ctx, ids)

	return u.TranRepo.Exec(ctx, func(ctx context.Context) error {
		err := u.UserRepo.Delete(ctx, ids)
		if err != nil {
			return err
		}
		//解除用户和角色关系
		err = u.DelUserRoleByUserId(ctx, ids)
		if err != nil {
			return err
		}
		//接触用户和岗位联系
		return u.DelUserPostByUserId(ctx, ids)
	})
}

func (u *UserService) ChangeStatus(ctx context.Context, id uint64, status uint8) error {
	return u.UserRepo.Where(ctx, "id=?", id).UpdateColumn("status", status).Error
}

func (u *UserService) ChangePassword(ctx context.Context, params dto.ChangePwdParams) error {
	decryptOldPwdStr, err := rsa.PrivateDecrypt(params.Password, config.Conf.RSA.PrivateKey)
	if err != nil {
		return err
	}
	decryptNewPwdStr, err := rsa.PrivateDecrypt(params.NewPassword, config.Conf.RSA.PrivateKey)
	if err != nil {
		return err
	}
	//新密码校验，并返回hash
	NewPwdHash := u.CheckPwdAndHash(decryptNewPwdStr)

	user, err := u.UserRepo.FindById(ctx, params.Id)
	if err != nil {
		return err
	}

	if user.Status != 1 {
		panic(respErr.UserDisableError)
	} else { //校验原始密码是否正确
		hash, err := bcrypt.GenerateFromPassword([]byte(decryptOldPwdStr), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		checkErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(hash))
		if checkErr != nil {
			panic(respErr.WrapLogicResp(i18n.DefTr("admin.pwdIncorrect")))
		}
	}
	return u.UserRepo.Where(ctx, "id=?", params.Id).UpdateColumn("password", NewPwdHash).Error
}

func (u *UserService) ResetPasswordByPhone(ctx context.Context, phone string, password string) error {
	//密码处理
	hash := u.CheckPwdAndHash(password)
	return u.UserRepo.Where(ctx, "phone=?", phone).UpdateColumn("password", hash).Error
}

// 检查用户是否允许操作
func CheckUserAllowed(user entity.User) {
	if user.IsSuper() {
		panic(respErr.BadRequestErrorWithMsg("不允许操作超级管理员用户"))
	}
}

// 检查用户是否有数据权限
func (u *UserService) checkUserDataScope(ctx context.Context, userIds []uint64) {
	var users = make([]entity.User, 0)
	err := u.UserRepo.Where(ctx, "id in ?", userIds).Find(users).Error
	if err != nil {
		log.L.Error(err)
		panic(respErr.InternalServerError)
	}
	if len(users) == 0 {
		panic(respErr.WrapLogicResp("没有权限访问用户数据！"))
	}
}

// 密码解密、检验并hash返回
func (u *UserService) CheckPwdAndHash(pwd string) string {
	//处理密码
	decryptStr, err := rsa.PrivateDecrypt(pwd, config.Conf.RSA.PrivateKey)
	if err != nil {
		panic(respErr.InternalServerErrorWithMsg(i18n.DefTr("admin.AccountOrPwdErr")))
	}
	//密码强度校验
	if !utils.PasswordStrength(decryptStr) {
		panic(respErr.BadRequestErrorWithMsg(i18n.DefTr("tip.pwdFormatTip", 6, 18, 2)))
	}
	//加密密码 采用bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(decryptStr), bcrypt.DefaultCost)
	if err != nil {
		log.L.Error(err)
		panic(respErr.InternalServerError)
	}
	return string(hash)
}

// 解密并哈希
func (u *UserService) decryptAndHash(pwd string) string {
	//处理密码
	decryptStr, err := rsa.PrivateDecrypt(pwd, config.Conf.RSA.PrivateKey)
	if err != nil {
		if config.Conf.IsDebug() {
			panic(err)
		} else {
			panic(respErr.InternalServerErrorWithMsg(i18n.DefTr("admin.AccountOrPwdErr")))
		}
	}
	if decryptStr == "" { //客户端密码可能没加密
		if config.Conf.IsDebug() {
			panic(respErr.BadRequestErrorWithMsg("密码可能没有加密"))
		} else {
			panic(respErr.InternalServerErrorWithMsg(i18n.DefTr("admin.AccountOrPwdErr")))
		}
	}
	if len(decryptStr) < 6 {
		panic(respErr.BadRequestErrorWithMsg(i18n.DefTr("tip.pwdLengthTip")))
	}
	//加密密码 采用bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(decryptStr), bcrypt.DefaultCost)
	if err != nil {
		log.L.Error(err)
		panic(respErr.InternalServerError)
	}
	return string(hash)
}
