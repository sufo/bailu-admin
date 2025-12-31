/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package types

import (
	"encoding/json"
	"fmt"
)

type Uint64EmptySlice []uint64

// 实现一个 JSON的MarshalNilArray方法，用来解决查询数组返回null问题
func (n Uint64EmptySlice) MarshalJSON() ([]byte, error) {
	if n == nil {
		//return []byte{}, nil
		//return []byte(fmt.Sprintf("[%v]", strings.Join(values, ","))), nil
		return []byte(fmt.Sprintf("[]")), nil
	} else {
		return json.Marshal(n)
	}
}

// 实现一个 JSON的UnmarshalJSON方法，用来解决接收前端参数抛异常的问题
func (n *Uint64EmptySlice) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, n)
}
