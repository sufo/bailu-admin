/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 岗位
 */

package entity

var _ IModel = (*Post)(nil)

type Post struct {
	//ID       uint64 `json:"id,string" gorm:"primarykey"`
	ID       uint64 `json:"id" gorm:"primarykey"`
	PostCode string `json:"postCode" gorm:"not null;size:64;gorm:comment:岗位编号"`
	Name     string `json:"name" gorm:"not null;size:50;comment:岗位名称"`
	Remark   string `json:"remark" gorm:"size:500;comment:备注"`
	SortAndStatus
	BaseEntity
}

var PostTN = "sys_post"

func (p Post) GetID() uint64 {
	return p.ID
}

func (Post) TableName() string {
	return PostTN
}
