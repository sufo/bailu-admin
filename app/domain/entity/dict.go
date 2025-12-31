/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 字典
 */

package entity

var _ IModel = (*Dict)(nil)

type Dict struct {
	ID          uint64 `json:"id" gorm:"primarykey"`
	Name        string `json:"name" gorm:"size:100;default:'';comment：字典名称" binding:"required"`
	Code        string `json:"code" gorm:"size:100;comment:字典类型/编码;unique" binding:"required"`
	Description string `json:"description" gorm:"size:255;comment:字典描述"`
	//SysFlag     uint8      `json:"sysFlag" gorm:"size:tinyint(4);comment:系统标志(1:系统类， 2:业务类)"`
	//Items []DictItem `json:"-" gorm:"foreignkey:code;"`
	BaseEntity
}

var DictTN = "sys_dict"

func (d Dict) GetID() uint64 {
	return d.ID
}

func (Dict) TableName() string {
	return DictTN
}
