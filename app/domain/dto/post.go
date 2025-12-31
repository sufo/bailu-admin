/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package dto

type PostParams struct {
	Name     string `json:"name" form:"name" query:"name,like"`
	PostCode string `json:"postCode" form:"postCode" query:"post_code,like"`
	Status   *int   `json:"status" form:"status" query:"status"`
}
