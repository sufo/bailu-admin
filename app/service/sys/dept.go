/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 部门
 */

package sys

import (
	"bailu/app/domain/entity"
	"bailu/app/domain/repo"
	"bailu/app/domain/resp"
	respErr "bailu/pkg/exception"
	"context"
	"github.com/google/wire"
	"strconv"
	"strings"
)

var DeptSet = wire.NewSet(wire.Struct(new(DeptService), "*"))

type DeptService struct {
	DeptRepo  *repo.DeptRepo
	UserRepo  *repo.UserRepo
	TransRepo *repo.Trans
}

func (d *DeptService) List(ctx context.Context, name string, status string) (*resp.PageResult[entity.Dept], error) {
	result, err := d.DeptRepo.FindDepts(ctx, name, status)
	return result.(*resp.PageResult[entity.Dept]), err
}

func (d *DeptService) Tree(ctx context.Context, name string, status string) ([]*entity.Dept, error) {
	result, err := d.DeptRepo.FindDepts(ctx, name, status)
	if err != nil {
		return nil, err
	}
	return d.BuildTree(result.([]*entity.Dept)), nil
}

func (d *DeptService) BuildTree(depts []*entity.Dept) []*entity.Dept {
	var tree = make([]*entity.Dept, 0)
	for _, ele := range depts {
		if ele.ID == 0 || !ele.HasParentNode(depts) {
			tree = append(tree, ele)
		}
	}
	for _, t := range tree {
		t.Children = d.recursion(depts, t.ID)
	}
	return tree
}

func (d *DeptService) recursion(depts []*entity.Dept, pid uint64) []*entity.Dept {
	var tree []*entity.Dept
	for _, ele := range depts {
		if ele.Pid == pid {
			tree = append(tree, ele)
		}
	}
	for _, t := range tree {
		t.Children = d.recursion(depts, t.ID)
	}
	return tree
}

// 新增部门
func (d *DeptService) Create(ctx context.Context, dept *entity.Dept) error {
	d.checkName(ctx, dept.Name, dept.Pid, 0)
	//判断父级部门状态
	if dept.Pid != 0 {
		p_dept, err := d.DeptRepo.FindById(ctx, dept.Pid)
		if err != nil {
			panic(respErr.InternalServerErrorWithError(err))
		}
		if p_dept.Status != 1 {
			panic(respErr.WrapLogicResp("部门停用，不允许新增"))
		}

		var p_ancestors = p_dept.Ancestors
		if p_ancestors != "" {
			p_ancestors = p_ancestors + ","
		}
		dept.Ancestors = p_ancestors + strconv.FormatUint(p_dept.ID, 10)
	}
	return d.DeptRepo.Create(ctx, dept)
}

func (d *DeptService) Update(ctx context.Context, dept *entity.Dept) error {
	d.checkName(ctx, dept.Name, dept.Pid, dept.ID)

	oldDept, err := d.DeptRepo.FindById(ctx, dept.ID)
	if err != nil {
		panic(respErr.InternalServerErrorWithError(err))
	}

	//说明该部门重新选择了新的父级部门
	if oldDept.Pid != dept.Pid {

		//当前dept不是顶级部门
		if dept.Pid != 0 {
			p_dept, err := d.DeptRepo.FindById(ctx, dept.Pid)
			if err != nil {
				panic(respErr.InternalServerErrorWithError(err))
			}
			//ancestors变量处理
			dept.Ancestors = p_dept.Ancestors + "," + strconv.FormatUint(p_dept.ID, 10)
		} else {
			dept.Ancestors = ""
		}

		//开启事务
		return d.TransRepo.Exec(ctx, func(ctx context.Context) error {
			//更新dept
			err := d.DeptRepo.Update(ctx, dept)
			if err != nil {
				return err
			}

			//更新所有子节点Ancestors
			err = d.updateChildrenAncestors(ctx, dept.Ancestors, oldDept.Ancestors)
			if err != nil {
				return err
			}

			// 如果该部门是启用状态，则启用该部门的所有上级部门
			if dept.Status == 1 && oldDept.Status != dept.Status {
				return d.enableParents(ctx, dept.Ancestors)
			}
			return nil
		})

	} else {
		return d.TransRepo.Exec(ctx, func(ctx context.Context) error {
			//更新dept
			if err := d.DeptRepo.Update(ctx, dept); err == nil {
				// 如果该部门是启用状态，则启用该部门的所有上级部门
				if dept.Status == 1 && oldDept.Status != dept.Status {
					return d.enableParents(ctx, dept.Ancestors)
				}
				return nil
			}
			return err
		})
	}
}

// 这里巧妙的采用ancestors的匹配，来找到所有的子孙子元素，然后替换更改的部分
// 解决了需要递归的问题
// 操作最好放在事务中
func (d *DeptService) updateChildrenAncestors(ctx context.Context, newAncestors string, oldAncestors string) error {
	depts, err := d.DeptRepo.FindBy(ctx, "ancestors like ?", oldAncestors+"%")
	if err != nil {
		panic(respErr.InternalServerErrorWithError(err))
	}
	for _, dept := range depts {
		ancestors := strings.Replace(dept.Ancestors, oldAncestors, newAncestors, 1)
		err := d.DeptRepo.Where(ctx, "id", dept.ID).Update("ancestors", ancestors).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// 启用所有父级节点
func (d *DeptService) enableParents(ctx context.Context, ancestors string) error {
	idsStr := strings.Split(ancestors, ",")
	size := len(idsStr)
	//id []string转[]uint64
	if size > 0 {
		var ids = make([]uint64, size)
		for _, idStr := range idsStr {
			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				return err
			}
			ids = append(ids, id)
		}
		return d.DeptRepo.Where(ctx, "id IN ?", ids).Update("status", true).Error
	}
	return nil
}

func (d *DeptService) Delete(ctx context.Context, id uint64) error {
	return d.TransRepo.Exec(ctx, func(ctx context.Context) error {
		err := d.DeptRepo.Delete(ctx, id)
		if err != nil {
			return err
		}
		//删除部门后，将用户中包含有该部门的用户更新
		return d.UserRepo.Where(ctx, "dept_id=?", id).Update("dept_id", nil).Error
	})
}

func (d *DeptService) DeptIdsByRoleId(ctx context.Context, roleId uint64) ([]uint64, error) {
	return d.DeptRepo.DeptIdsByRoleId(ctx, roleId, true)
}

// 检查名称是为重复
//func (d *DeptService) CheckDeptUnique(ctx context.Context, dept *entity.Dept) bool {
//	var info entity.Dept
//	deptId := dept.GetID()
//	if e := d.DeptRepo.Where(ctx, "dept_name=?", dept.Name).Where("pid=?", dept.Pid).First(&info).Error; e == nil {
//		return info.GetID() != deptId
//	} else if e == gorm.ErrRecordNotFound {
//		return true
//	} else {
//		panic(respErr.InternalServerErrorWithMsg(e.Error()))
//	}
//}

func (m *DeptService) HasChildByID(ctx context.Context, id uint64) (bool, error) {
	return m.DeptRepo.IsExist(ctx, "pid=?", id)
}

// name check
// 更新必须传id
func (m *DeptService) checkName(ctx context.Context, name string, pid uint64, id uint64) {
	var query = "name=? AND pid=?"
	var args = []interface{}{name, pid}
	if id != 0 {
		query += " AND id <> ?"
		args = append(args, id)
	}
	isExist, err := m.DeptRepo.IsExist(ctx, query, args...)
	if err != nil {
		panic(respErr.InternalServerErrorWithError(err))
	}
	if isExist {
		panic(respErr.BadRequestErrorWithMsg("部门名重复！"))
	}
}
