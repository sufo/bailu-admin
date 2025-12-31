/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package vo

import "bailu/app/domain/entity"

type Route struct {
	Name      string   `json:"name"` //路由名称
	Path      string   `json:"path"`
	Component string   `json:"component"`
	Meta      *Meta    `json:"meta"`
	Children  []*Route `json:"children"`
	Redirect  string   `json:"redirect"`
}

type Meta struct {
	Title   string `json:"title"`
	I18nKey string `json:"i18nKey"`
	Icon    string `json:"icon"`
	Query   string `json:"query"`
	//IgnoreAuth bool   `json:"ignoreAuth"`
	//Affix      bool   `json:"affix"`
	FrameSrc           string `json:"frameSrc"`
	KeepAlive          bool   `json:"keepAlive"` // 是否缓存
	IsFrame            bool   `json:"isFrame"`   // 是否外链
	Hide               bool   `json:"hide"`
	Permission         string `json:"permission"`         //主要是给前端做权限按钮控制(list表示列表查询，query表示详情)，后台则使用casbin做接口权限校验
	HideChildrenInMenu bool   `json:"hideChildrenInMenu"` //隐藏子菜单, 解决顶级菜单(C,Pid=0)没有子菜单时，此时为true
}

type Menu struct {
	ID        uint64  `json:"id,string"`
	Pid       *uint64 `json:"-"`
	Name      string  `json:"name"`
	I18nKey   string  `json:"i18nKey"`
	Path      string  `json:"path"`
	Component *string `json:"component"`
	Meta
	Type     *string `json:"type"`
	Children []*Menu `json:"children"`
	entity.SortAndStatus
	entity.BaseEntity
}
