/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package dto

type PaginationParam struct {
	//PageIndex int `form:"pageIndex,default=1"`
	//PageSize  int `form:"pageSize,default=10" binding:"max=100"`
	PageIndex int `form:"pageIndex" json:"pageIndex" binding:"omitempty,min=0"`
	PageSize  int `form:"pageSize" json:"pageSize" binding:"omitempty,min=0,max=100"`
}

type StatusParam struct {
	ID     uint64 `json:"id" binding:"required"`
	Status uint8  `json:"status" binding:"required"`
}
