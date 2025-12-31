/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package dto

type NoticeParams struct {
	Title string `json:"title" form:"title" query:"title,like"`
	//Sender string `json:"sender" form:"sender" query:"sender,like"`
	ReadStatus  *int   `json:"status" form:"status" query:"-"`                   //读取状态
	SendStatus  string `json:"sendStatus" form:"sendStatus" query:"send_status"` //发布状态
	IfScheduled string `json:"ifScheduled" form:"ifScheduled" query:"-"`         //是否定时任务
}
