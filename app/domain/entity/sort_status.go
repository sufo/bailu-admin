/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 包含排序和是否启用
 */

package entity

type SortAndStatus struct {
	Sort   *uint `json:"sort" gorm:"comment:排序;default:0"`
	Status uint8 `json:"status" gorm:"type:tinyint(4);comment:是否启用 (1:启用 2:禁用);default:0"`
}
