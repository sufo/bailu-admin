/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package dto

type OperParams struct {
	Path      string `json:"path" query:"path,like"`          //
	Method    string `json:"method" query:"method"`           //
	OperName  string `json:"OperName" query:"oper_name,like"` // 用户名
	Status    string `json:"status" query:"-"`                //状态 0成功 1失败 ,这个条件要单独处理
	BeginDate string `json:"beginDate" query:"created_at,between EndDate"`
	EndDate   string `json:"endDate"`
}

type LoginLogParams struct {
	Addr      string `json:"addr" query:"addr,like"`
	Username  string `json:"username" query:"username,like"`
	Status    *int   `json:"status" query:"status"`
	BeginDate string `json:"beginDate" query:"created_at,between EndDate"`
	EndDate   string `json:"endDate"`
}
