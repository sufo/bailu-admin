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

// 实现一个 JSON的MarshalNilArray方法，用来解决查询数组返回null问题
type EmptySlice[T any] []T

//func (n EmptySlice[T]) MarshalJSON() ([]byte, error) {
//	if n == nil {
//		return []byte(fmt.Sprintf("[]")), nil
//	} else {
//		return json.Marshal(n)
//	}
//}

func (n EmptySlice[T]) MarshalJSON() ([]byte, error) {
	if n == nil {
		return []byte(fmt.Sprintf("[]")), nil
	} else {
		var temp = make([]T, len(n))
		for k, item := range n {
			temp[k] = item
		}
		return json.Marshal(temp)
	}
}

type _EmptySlice[T any] EmptySlice[T]

// 实现一个 JSON的UnmarshalJSON方法，用来解决接收前端参数抛异常的问题
func (n *EmptySlice[T]) UnmarshalJSON(b []byte) error {
	//这样会出现无限循环
	//return json.Unmarshal(b, n)

	var n2 _EmptySlice[T] = make(_EmptySlice[T], 0)

	err := json.Unmarshal(b, &n2)
	if err != nil {
		return err
	}
	*n = EmptySlice[T](n2)
	return nil
}
