/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 数据库字符串null值问题处理
 */

package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type NullString sql.NullString

// 实现一个 JSON的MarshalJSON方法，用来解决查询返回null抛异常问题
func (n *NullString) MarshalJSON() ([]byte, error) {
	ns := n.String
	return []byte(fmt.Sprintf("\"%s\"", ns)), nil
}

// 实现一个 JSON的UnmarshalJSON方法，用来解决接收前端参数抛异常的问题
func (n *NullString) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &n.String)
}
