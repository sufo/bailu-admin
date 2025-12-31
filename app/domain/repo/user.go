/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package repo

import (
	"bailu/app/domain/entity"
	"bailu/app/domain/repo/base"
	"bailu/app/domain/repo/util"
	"context"
	"errors"
	"gorm.io/gorm"
)

//type UserRepo struct {
//	gorme.Repository[entity.User]
//}
//
//func NewUserRepo(db *gorm.DB) *UserRepo {
//	repo := UserRepo{}
//	repo.SetDB(db)
//	return &repo
//}

//var UserSet = wire.NewSet(wire.Struct(new(UserRepo), "*"))

func NewUserRepo(db *gorm.DB) *UserRepo {
	r := base.Repository[entity.User]{db}
	return &UserRepo{r}
}

type UserRepo struct {
	base.Repository[entity.User]
}

func (u *UserRepo) FindByName(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := util.GetDBWithModel[entity.User](ctx, u.DB).
		Preload("Posts").Preload("Roles").
		Joins("left join sys_dept d on d.id=sys_user.id").
		Select("sys_user.*,d.name as deptName").
		Where("username=?", username).
		First(&user).Error
	return &user, err
}

// 这里更新iP
func (u *UserRepo) UpdateIp(ctx context.Context, id uint64, ip string) error {
	return u.Where(ctx, "id=?", id).UpdateColumn("ip", ip).Error
}

func (u *UserRepo) Updates(ctx context.Context, id uint64, columns []string, args ...any) error {
	if len(columns) != len(args) {
		return errors.New("columns数组长度args个数必须一致")
	}
	values := make(map[string]any)
	db := u.Where(ctx, "id=?", id)
	for index, col := range columns {
		values[col] = args[index]
	}
	db.UpdateColumns(values)
	return db.Error
}

// 解除用户和岗位关系
func (u *UserRepo) UntiedPost(ctx context.Context, postIds []uint64) error {
	var userPost entity.UserPost
	return u.GetDB(ctx).Where("post_id in ?", postIds).Unscoped().Delete(&userPost).Error
}

//func (u *UserRepo) Query(ctx context.Context, params dto.UserQueryParams) (*resp.PageResult[entity.User], error) {
//	db := base.GetDB(ctx, u.DB).Debug().Table(fmt.Sprintf("%s as u", entity.UserTN)).Preload("Roles").Preload("Posts").
//		Joins(fmt.Sprintf("left join %s as d on u.dept_id=d.id", entity.DeptTN))
//	if v := params.Username; v != "" {
//		v = "%" + v + "%"
//		db = db.Where("u.username LIKE ? OR u.nick_name LIKE ?", v, v)
//	}
//	if v := params.DeptId; v != nil {
//		db = db.Where("u.dept_id=?", v)
//	}
//
//	if v := params.Enable; v != nil {
//		db = db.Where("u.status=?", v)
//	}
//	if v := params.Phone; v != "" {
//		v = "%" + v + "%"
//		db = db.Where("u.phone LIKE ? ", v)
//	}
//
//	//处理日期范围
//	if params.BeginDate != "" && params.EndDate != "" {
//		bDate, err := time.Parse("20060102", params.BeginDate)
//		if err != nil {
//			panic(respErr.BadRequestErrorWithError(err))
//		}
//		eDate, err := time.Parse("20060102", params.EndDate)
//		if err != nil {
//			panic(respErr.BadRequestErrorWithError(err))
//		}
//		db = db.Where("u.created_at BETWEEN ? AND ?", bDate, eDate)
//
//	} else if params.BeginDate != "" {
//		bDate, err := time.Parse("20060102", params.BeginDate)
//		if err != nil {
//			panic(respErr.BadRequestErrorWithError(err))
//		}
//		db.Where("u.updated_at > ?", bDate)
//	} else if params.EndDate != "" {
//		eDate, err := time.Parse("20060102", params.EndDate)
//		if err != nil {
//			panic(respErr.BadRequestErrorWithError(err))
//		}
//		db.Where("u.updated_at < ?", eDate)
//	}
//	db = util.DataScope(ctx, db, "u", "d")
//	return util.Paginate[entity.User](db, params.PageIndex, params.PageSize, "ID ASC")
//}
