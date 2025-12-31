/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 部门
 */

package entity

var _ IModel = (*Dept)(nil)

type Dept struct {
	ID        uint64  `json:"id,string" gorm:"primarykey" mapstructure:"id"`
	Pid       uint64  `json:"pid,string" gorm:"default:0; comment:父部门ID" mapstructure:"pid"`
	Ancestors string  `json:"ancestors" gorm:"default:''; comment:上级部门列表，逗号隔开。dataScope会用到" mapstructure:"ancestors"`
	Name      string  `json:"name" gorm:"not null;size:30;comment:部门名称" binding:"required" mapstructure:"name"`
	Sort      int     `json:"sort" gorm:"default:0;comment:排序" mapstructure:"sort"`
	Leader    string  `json:"leader" gorm:"size:20;comment:部门领导" mapstructure:"leader"`
	Phone     string  `json:"phone" gorm:"size:12;comment:联系电话" mapstructure:"phone"`
	Email     string  `json:"email" gorm:"size:50;comment:邮箱" mapstructure:"email"`
	Status    uint8   `json:"status" gorm:"type:tinyint(4);comment:是否启用(1:启用 2:禁用);default 0" binding:"oneof=0 1" mapstructure:"status"`
	Children  []*Dept `json:"children" gorm:"-"`
	BaseEntity
}

var DeptTN = "sys_dept"

func (Dept) TableName() string {
	return DeptTN
}

func (d Dept) GetID() uint64 {
	return d.ID
}

// whether it has parent node or not
func (d *Dept) HasParentNode(depts []*Dept) (has bool) {
	for _, ele := range depts {
		has = d.Pid == ele.ID
		if has {
			break
		}
	}
	return
}
