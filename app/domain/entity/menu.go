/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc menu
 */

package entity

import "bailu/utils/types"

var _ IModel = (*Menu)(nil)

type Menu struct {
	//ID       uint64  `json:"id,string" gorm:"primarykey"`
	ID      uint64  `json:"id" gorm:"primarykey"`
	Pid     *uint64 `json:"pid" gorm:"default:0;comment:父菜单ID" `
	Name    string  `json:"name" gorm:"NOT NULL;default:'';size:30;comment:菜单名称;"` //这里数据库里面取名dept_name,为避免链接查询重复
	I18nKey string  `json:"i18nKey" gorm:"default:'';size:30;comment:国际化名称key;"`
	Path    string  `json:"path" gorm:"default:'';comment:路由路径（链接地址）"`
	//Component types.NullString `json:"component" gorm:"size:255;default:null;comment:组件路径"`
	Component *string `json:"component" gorm:"size:255;default:null;comment:组件路径"`
	Meta      `gorm:"embedded;comment:附加属性"` //`json:"meta" gorm:"embedded;comment:附加属性"`
	//Type      types.NullString `json:"type" gorm:"type:char(1);default:NULL;comment:菜单类型(M目录 C菜单 F按钮）;binding:required"`
	Type     *string                   `json:"type" gorm:"type:char(1);default:NULL;comment:菜单类型(M目录 C菜单 F按钮）;binding:required"`
	Children []*Menu                   `json:"children" gorm:"-"`
	Apis     types.EmptySlice[MenuApi] `json:"apis" gorm:"foreignkey:menu_id;references:id"` //该参数只针对按钮(Type==F)
	//Apis types.EmptySlice[*MenuApi] `json:"apis" gorm:"-"` //该参数只针对按钮(Type==F)
	SortAndStatus
	BaseEntity
}

func (m Menu) GetID() uint64 {
	return m.ID
}

// bool在mysql中对应TINYINT(1) 0对应false
type Meta struct {
	Icon string `json:"icon" gorm:"size:100;default:'#';comment:菜单图标"`
	//Query      types.NullString `json:"query" gorm:"size:255;default:NULL;comment:路由参数"`
	Query *string `json:"query" gorm:"size:255;default:NULL;comment:路由参数"`
	//IgnoreAuth bool    `json:"ignoreAuth" gorm:"default:false;comment:忽略验证; true忽略"`
	//Affix     bool `json:"affix" gorm:"default:false;comment:是否固定在tab上,true固定"`
	KeepAlive bool `json:"keepAlive" gorm:"type:tinyint(1);default:true;comment:是否缓存"` // 是否缓存
	// 是否外链，点击之后，直接在浏览器打开一个新的选项卡
	// 注意，是否外链跟地址是不是超链接没有关系，如果地址是超链接，但是IsFrame是false，那么这个链接会在当前系统的iframe中打开
	IsFrame bool `json:"isFrame" gorm:"type:tinyint(1);default:false;comment:是否外链"`
	Hide    bool `json:"hide" gorm:"type:tinyint(1);default:false;comment:是否隐藏 true隐藏 false可见;不可见的话不会出现在侧边栏，但可以访问"`
	//Permission types.NullString `json:"permission" gorm:"size:100; default:NULL;comment:权限标识，例如：user:add"` //主要是给前端做权限按钮控制(list表示列表查询，query表示详情)，后台则使用casbin做接口权限校验
	Permission *string `json:"permission" gorm:"size:100; default:NULL;comment:权限标识，例如：user:add"` //主要是给前端做权限按钮控制(list表示列表查询，query表示详情)，后台则使用casbin做接口权限校验
}

var MenuTN = "sys_menu"

func (Menu) TableName() string {
	return MenuTN
}

// 排序
type MenuSort []*Menu

//PersonSort 实现sort SDK 中的Interface接口

func (s MenuSort) Len() int {
	//返回传入数据的总数
	return len(s)
}
func (s MenuSort) Swap(i, j int) {
	//两个对象满足Less()则位置对换
	//表示执行交换数组中下标为i的数据和下标为j的数据
	s[i], s[j] = s[j], s[i]
}
func (s MenuSort) Less(i, j int) bool {
	//按字段比较大小,此处是降序排序
	//返回数组中下标为i的数据是否小于下标为j的数据
	//return s[i] > s[j]

	if *s[i].Pid == *s[j].Pid {
		//return *s[i].Sort > *s[j].Sort  //sort有可能为nil，这样写会空指针
		var si, sj uint
		//var si = global.Ternary(s[i].Sort == nil, 0, *s[i].Sort)
		//var sj = global.Ternary(s[j].Sort == nil, 0, *s[j].Sort)
		if s[i].Sort != nil {
			si = *s[i].Sort
		}
		if s[j].Sort != nil {
			sj = *s[j].Sort
		}
		return si > sj
	}
	return *s[i].Pid > *s[j].Pid
}

// whether it has parent node or not
func (s *Menu) HasParentNode(menus []*Menu) (has bool) {
	for _, ele := range menus {
		has = *s.Pid == ele.ID
		if has {
			break
		}
	}
	return
}

// 按钮绑定api
type MenuApi struct {
	MenuId uint64 `json:"menuId" gorm:"not null;comment:菜单id"`
	Method string `json:"method" gorm:"not null;size:20;comment:请求方法" binding:"required"`
	Path   string `json:"path" gorm:"not null;comment:按钮对应的接口路径" binding:"required"`
}

var MenuApiTN = "sys_menu_api"

func (MenuApi) TableName() string {
	return MenuApiTN
}

func (m MenuApi) GetID() uint64 {
	return m.MenuId
}
