/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package dto

type RoleUsers struct {
	RoleId  uint64   `json:"roleId" binding:"required" query:"role_id"`
	UserIds []uint64 `json:"userIds" binding:"required,min=1" query:"user_id,in"`
}
