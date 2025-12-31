/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc menu参数
 */

package dto

type MenuParams struct {
	Name        string `json:"name" form:"name" query:"name,like"`             //菜单名字  --title
	Status      *int   `json:"status" form:"status" query:"status"`            //是否开启
	ExcludeType string `json:"excludeType" form:"excludeType" query:"type,!="` //排除的菜单类型
	Pid         *int   `json:"pid" form:"pid" query:"pid"`                     //父级id
}
