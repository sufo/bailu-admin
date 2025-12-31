/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 分页响应对象
 */

package resp

type PageResult[T any] struct {
	PageSize  int   `json:"pageSize"`  //每页条数
	PageIndex int   `json:"pageIndex"` //当前页码
	PageCount int64 `json:"pageCount"` //总页数
	ItemCount int64 `json:"itemCount"` //总条数
	List      []*T  `json:"list"`      //数据列表
}
