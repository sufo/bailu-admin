/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package entity

import (
	"github.com/sufo/bailu-admin/app/core/appctx"
	"github.com/sufo/bailu-admin/pkg/log"
	"gorm.io/gorm"
	"reflect"
	"time"
)

const (
	Enabled  uint8 = 1 // 启用
	Running  uint8 = 1 // 运行中
	Finish   uint8 = 2 // 完成
	Disabled uint8 = 2 // 禁用
	Failure  uint8 = 3 // 失败
	Cancel   uint8 = 4 // 取消
)

type BaseEntity struct {
	Model
	CreateBy uint64 `json:"-" gorm:"column:create_by;default:0;comment:创建者"`
	UpdateBy uint64 `json:"-" gorm:"column:update_by;default:0;comment:更新者"`
}

// hook函数处理CreateBy
func (r *BaseEntity) BeforeCreate(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	if ctx != nil {
		user := appctx.GetAuthUser[OnlineUserDto](ctx)
		if user != nil {
			r.CreateBy = user.ID
			
		} else {
			typeof := reflect.TypeOf(tx.Statement.Model)
			log.L.Errorf("%s BeforeCreate: user is nil", typeof.Elem().Name())
		}
	}
	return nil
}

// hook函数处理UpdateBy
func (r *BaseEntity) BeforeUpdate(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	if ctx != nil {
		user := appctx.GetAuthUser[OnlineUserDto](ctx)
		if user != nil {
			r.UpdateBy = user.ID

		} else {
			typeof := reflect.TypeOf(tx.Statement.Model)
			log.L.Error("%s BeforeUpdate: user is nil ", typeof.Elem().Name())
		}
	}
	return nil
}

/*
*
gorm支持对model中的时间字段进行默认操作：
CreatedAt：创建时间
UpdatedAt：更新时间
DeletedAt: 删除时间
也就是说，若model中含有上述字段，在CRUD操作时，gorm会自动更新上述的时间字段的值为now()。
*/
type Model struct {
	//ID        uint64         `json:"id,string" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"` //`gorm:"default:null"`
	UpdatedAt time.Time `json:"updatedAt"` //`gorm:"default:null"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

// 其他实体需要实现
type IModel interface {
	TableName() string
	GetID() uint64
}
