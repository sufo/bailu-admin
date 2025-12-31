/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 上传文件管理
 */

package entity

type FileInfo struct {
	ID         uint64 `json:"id,string" gorm:"primarykey" mapstructure:"id"`
	Url        string `json:"url" gorm:"size:255;comment:地址,本地存储则保存的是路径"`
	CategoryId uint64 `json:"categoryId" gorm:"not null;comment:分类ID;default:0"`
	Name       string `json:"name" gorm:"size:80;comment:文件名"`
	Size       int64  `json:"size" gorm:"size:255;comment:文件大小（KB）"`
	OriginName string `json:"originName" gorm:"size:80;comment:原始文件名"`
	MIME       string `json:"type" gorm:"column:mime;size:30;comment:文件类型"`
	Path       string `json:"path" gorm:"size:200;comment:文件路径"`
	Tags       string `json:"tags" gorm:"comment:多tag以逗号分割"`
	BaseEntity
}

var FileTN = "files"

func (i FileInfo) GetID() uint64 {
	return i.ID
}

func (FileInfo) TableName() string {
	return FileTN
}

type FileCategory struct {
	ID   uint64 `json:"value,string" gorm:"primarykey"`
	Name string `json:"label" gorm:"size:50;comment:文件名" binding:"required"`
	BaseEntity
}

var CategoryTN = "file_category"

func (i FileCategory) GetID() uint64 {
	return i.ID
}

func (FileCategory) TableName() string {
	return CategoryTN
}
