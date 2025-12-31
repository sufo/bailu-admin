/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 我的消息
 */

package dto

type MessageParams struct {
	Type  string `json:"type" form:"type" query:"-" binding:"required"` //查询类型 notice、event、chat、msg或空查所有
	Title string `json:"title" form:"title" query:"title,like"`
	//ReadFlag string `json:"readFlag" query:"-"` //0未读 1已读
}
