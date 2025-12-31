/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package vo

import "bailu/utils/types"

// 下拉菜单
type Option[T uint64 | string] struct {
	Value     T      `json:"value"` //
	Label     string `json:"label"`
	IsDefault bool   `json:"isDefault"` //是否默认选中
}

// 下拉菜单树
type Tree[T uint64 | string] struct {
	Option[T]
	Children []Tree[T] `json:"children"`
}

type KV struct {
	Value string `json:"value"` //
	Label string `json:"label"`
}

type Message struct {
	ID      string         `json:"id"`
	Icon    string         `json:"icon"`
	Avatar  string         `json:"avatar"`
	Title   string         `json:"title"`
	Content string         `json:"content"`
	Date    types.JSONTime `json:"date"`
	IsRead  bool           `json:"isRead"`
	ToId    uint64         `json:"toId"`   //chat
	ToName  string         `json:"toName"` //chat
}
