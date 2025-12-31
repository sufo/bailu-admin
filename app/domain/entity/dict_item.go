/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 字典项
 */

package entity

var _ IModel = (*DictItem)(nil)

type DictItem struct {
	ID uint64 `json:"id" gorm:"primarykey"`
	//DictId uint64 `json:"dictId" gorm:"comment:字典ID 外键"`
	Label string `json:"label" gorm:"default:'';size:100;comment:字典标签;" binding:"required"`
	Value string `json:"value" gorm:"default:'';size:100;comment:字典值" binding:"required"`
	Code  string `json:"code" gorm:"default:'';size:100;comment:编码" binding:"required"`
	//Remark    types.NullString `json:"remark" gorm:"default:null;size:255;comment:描述"`
	Remark *string `json:"remark" gorm:"default:null;size:255;comment:描述"`
	//IsDefault int   `json:"isDefault" gorm:"type:tinyint(4);default:0;comment:是否默认选中 (1:选中 2:不选)"`
	IsDefault *bool `json:"isDefault" gorm:"default:false;comment:是否默认选中"`
	//Fixed     int   `json:"fixed" gorm:"type:tinyint(4);default:1;comment:是否固定（固定的字典不提供编辑功能） (1:不固定 2：固定)"`
	Fixed *bool `json:"fixed" gorm:"default:false;comment:是否固定（固定的字典不提供编辑功能）"`
	SortAndStatus
	BaseEntity
}

var DictItemTN = "sys_dict_item"

func (i DictItem) GetID() uint64 {
	return i.ID
}

func (DictItem) TableName() string {
	return DictItemTN
}
