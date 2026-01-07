/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 系统配置表
 */

package entity

import "time"

type SysConfig struct {
	ID        uint64         `json:"id" gorm:"primarykey"`
	Name      string         `json:"name" gorm:"size:100;comment:参数名称" binding:"required"`
	Key       string         `json:"key" gorm:"size:100;comment:参数键名" binding:"required"`
	Value     string         `json:"value" gorm:"size:500;comment:参数值" binding:"required"`
	Type      string         `json:"type" gorm:"type:tinyint(1);comment:1 系统类 2 业务类" binding:"required"`
	Status    int            `json:"status" gorm:"type:tinyint(4);comment:是否启用 (1:启用 2:禁用);default:0"`
	Remark    string         `json:"remark" gorm:"size:500;comment:备注"`
	CreateBy  uint64         `json:"-" gorm:"column:create_by;default:0;comment:创建者"`
	UpdateBy  uint64         `json:"-" gorm:"column:update_by;default:0;comment:更新者"`
	CreatedAt time.Time `json:"createdAt"` //`gorm:"default:null"`
	UpdatedAt time.Time `json:"updatedAt"` //`gorm:"default:null"`
}

var ConfigTN = "sys_config"

func (t SysConfig) TableName() string {
	return ConfigTN
}

func (t SysConfig) GetID() uint64 {
	return t.ID
}
