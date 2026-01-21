/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc menus
 */

package sys

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/sufo/bailu-admin/app/domain/dto"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/repo"
	"github.com/sufo/bailu-admin/app/domain/vo"
	"github.com/sufo/bailu-admin/global/consts"
	respErr "github.com/sufo/bailu-admin/pkg/exception"
	"github.com/sufo/bailu-admin/utils"
	"regexp"
	"sort"
	"strings"
)

var MenuSet = wire.NewSet(wire.Struct(new(MenuService), "*"))

type MenuService struct {
	MenuRepo    *repo.MenuRepo
	MenuApiRepo *repo.MenuApiRepo
	RoleRepo    *repo.RoleRepo
	TranRepo    *repo.Trans
}

// init menu file insert into menu table
func (m *MenuService) InitData(ctx context.Context, menuPath string) error {
	isExit, err := m.MenuRepo.IsExist(ctx, nil)
	if err != nil {
		return err
	}
	if isExit {
		return nil
	}
	data, err := utils.LoadYaml2Struct[[]*entity.Menu](menuPath)
	if err != nil {
		return err
	}
	//定义接受结果的切片
	var result = make([]*entity.Menu, 0)
	//
	m.flatTree(&result, uint64(0), nil, data)
	//排序
	sort.Sort(entity.MenuSort(result))

	//批量插入数据库
	return m.MenuRepo.WithModel(ctx).CreateInBatches(result, len(result)).Error
}

// 这里的result之所以传指针，是因为append可能导致result扩容，扩容本质是产生了一个新的数组
// 如果在函数内对切片添加元素导致扩容，会导致元素内的切片指向一个新的数组，但是函数外的切片仍然指向原来旧的数组，
// 则将会导致影响无法传递到函数外。如果希望函数内对切片扩容作用于函数外，就需要以指针形式传递切片。
func (m *MenuService) flatTree(result *[]*entity.Menu, pId uint64, index *uint64, list []*entity.Menu) {
	//这里采用整数自增实现ID自增且唯一
	if index == nil {
		var i uint64 = 0
		index = &i
	}
	for _, item := range list {
		*index++ //自增只能后置，而且要独占一行
		menu := &entity.Menu{
			//ID:         uint64(index + 1 + lastIndex), //因为index从0开始，所以这里加1
			ID:         *index,
			Pid:        &pId,
			Name:       item.Name,
			Component:  item.Component,
			Path:       item.Path,
			Meta:       item.Meta,
			Type:       item.Type,
			BaseEntity: item.BaseEntity,
		}

		*result = append(*result, menu)
		if item.Children != nil && len(item.Children) > 0 {
			m.flatTree(result, menu.ID, index, item.Children)
		}
	}
}

// 查询菜单列表
func (m *MenuService) FindMenuList(c *gin.Context, params dto.MenuParams) ([]*entity.Menu, error) {
	user, _ := c.Get(consts.REQUEST_USER)
	userDto := user.(*entity.OnlineUserDto)
	ctx := c.Request.Context()

	if userDto.IsSuper() {
		return m.MenuRepo.FindMenus(ctx, params)
	} else {
		return m.MenuRepo.FindMenusByUserId(ctx, userDto.ID, params)
	}
}

// 获取role对应的菜单id集合
func (m *MenuService) FindCheckedMenusByRoleId(c *gin.Context, roleId uint64) ([]uint64, error) {
	//user, _ := c.Get(consts.REQUEST_USER)
	//userDto := user.(*entity.OnlineUserDto)
	ctx := c.Request.Context()
	//如果是超级管理员
	if roleId == consts.SUPER_ROLE_ID {
		panic(respErr.ForbiddenError)
	} else {
		return m.MenuRepo.FindByRoleId(ctx, roleId)
	}
}

/**
 * @menus 把菜单列表变为层级结构
 * @pid pId
 */
func (m *MenuService) BuildTree(menus []*entity.Menu, pid uint64) []*entity.Menu {
	var tree = make([]*entity.Menu, 0)
	for _, ele := range menus {
		if *ele.Pid == pid {
			tree = append(tree, ele)
		}
	}
	for _, t := range tree {
		t.Children = m.BuildTree(menus, t.ID)
	}
	return tree
}

// 构建任意菜单数，并不一定要从pid为0开始
func (m *MenuService) BuildAnyTree(menus []*entity.Menu) []*entity.Menu {
	var tree = make([]*entity.Menu, 0)
	for _, ele := range menus {
		if ele.ID == 0 || !ele.HasParentNode(menus) {
			tree = append(tree, ele)
		}
	}
	for _, t := range tree {
		t.Children = m.BuildTree(menus, t.ID)
	}
	return tree
}

// 数据适配
func (m *MenuService) _adapter(tree []*entity.Menu) []*vo.Route {
	var routers = make([]*vo.Route, len(tree))
	for i, menu := range tree {
		var route = &vo.Route{}
		//meta
		var meta = vo.Meta{
			menu.Name,
			menu.I18nKey,
			menu.Icon,
			"", "",
			menu.KeepAlive,
			menu.IsFrame,
			menu.Hide,
			"",
			false,
		}
		if menu.Query != nil {
			meta.Query = *menu.Query
		}
		if menu.Permission != nil {
			meta.Permission = *menu.Permission
		}

		//外链外部打开
		if menu.IsFrame {
			route.Name = ""
			meta.FrameSrc = menu.Path
		} else {
			var routePath = menu.Path
			var routeName = ""
			var component = ""
			if menu.Component != nil {
				component = *menu.Component
			}

			//如果是内部跳转链接
			if isInnerLink(menu) {
				meta.FrameSrc = menu.Path
				routePath = link2RoutePath(menu.Path)
				component = "IFRAME"
			}
			//顶级目录／菜单
			if *menu.Pid == 0 || menu.Pid == nil {
				if menu.Component == nil || *menu.Component == "" {
					component = "Layout"
				}
			} else if consts.TYPE_DIR == *menu.Type {
				// 如果不是一级菜单，并且菜单类型为目录，则代表是多级菜单
				if menu.Component == nil || *menu.Component == "" {
					component = "ParentView"
				}
			}

			route.Path = routePath
			route.Component = component
			routeName = utils.ToUpperForFirstCharAtSymbolBehind(routePath, "/")
			route.Name = routeName

			//有子节点
			if menu.Children != nil && len(menu.Children) > 0 {
				route.Redirect = route.Path + "/" + menu.Children[0].Path
				route.Children = m._adapter(menu.Children)
			} else if (menu.Pid == nil || *menu.Pid == 0) && consts.TYPE_MENU == *menu.Type {
				// 处理是一级菜单并且没有子菜单的情况
				//这种情况路由增加一层父级layout。这是为了符合前端router的结构，把菜单path="",则父级不需要Redirect，
				route.Redirect = route.Path + "/index"
				route.Component = "LAYOUT"
				route.Name = routeName + "Parent" //给父级设置route name，可以随意定义，不重复即可
				var children []*vo.Route
				var _meta = meta
				children = append(children, &vo.Route{
					Name: routeName,
					//Path: global.Ternary(strings.HasPrefix(menu.Path, "/"), menu.Path[1:], menu.Path),
					Path: "index",
					//Path:      "",
					Component: *menu.Component,
					Meta:      &_meta,
				})
				route.Children = children
				//这种情况需要隐藏（因为动态将一级变成了两级）
				meta.HideChildrenInMenu = true
			}
		}

		//顶级目录需要添加/
		if (*menu.Pid == 0 || menu.Pid == nil) && !strings.HasPrefix(menu.Path, "/") {
			route.Path = "/" + route.Path
		}
		route.Meta = &meta

		routers[i] = route
	}
	return routers
}

// 前端路由菜单
func (m *MenuService) FindRoutes(c *gin.Context) ([]*vo.Route, error) {
	var status = 1
	var params = dto.MenuParams{Status: &status, ExcludeType: "F"}
	menus, err := m.FindMenuList(c, params)
	if err != nil {
		return nil, err
	}
	return m._adapter(m.BuildTree(menus, 0)), nil
}

// 根据id查询子菜单
func (m *MenuService) FindByPid(c *gin.Context, id uint64) (menu []*entity.Menu, err error) {
	user, _ := c.Get(consts.REQUEST_USER)
	userDto := user.(entity.OnlineUserDto)
	ctx := c.Request.Context()
	if userDto.IsSuper() {
		err = m.MenuRepo.Where(ctx, "pid=?", id).Find(&menu).Error
		return
	} else {
		roleIds := make([]uint64, len(userDto.Roles))
		for _, role := range userDto.Roles {
			roleIds = append(roleIds, role.ID)
		}
		if roleIds == nil || len(roleIds) == 0 {
			return nil, errors.New("未知角色")
		}
		var pid uint64 = 0
		return m.MenuRepo.FindByRolesAndTypeNot(ctx, roleIds, "F", &pid)
	}

}

// 创建
func (m *MenuService) Save(ctx context.Context, menu entity.Menu) error {
	//如果不是按钮
	if *menu.Type != "F" {
		if menu.Path == "" {
			panic(respErr.BadRequestErrorWithMsg("path不能为空"))
		}
	}
	//相同父级下不允许有相同名称
	m.checkName(ctx, menu.Name, menu.Pid, menu.ID)
	//父级停用不允许新增子菜单
	if menu.Pid != nil && *menu.Pid != 0 {
		p_menu, err := m.MenuRepo.FindById(ctx, *menu.Pid)
		if err != nil {
			panic(respErr.InternalServerErrorWithError(err))
		}
		if p_menu.Status != 1 {
			panic(respErr.WrapLogicResp("父级菜单停用，不允许新增"))
		}
	}
	checkIsFrame(&menu)
	if *menu.Type == "F" { //是按钮
		if menu.Apis == nil {
			return m.MenuRepo.Create(ctx, &menu)
		} else {
			return m.TranRepo.Exec(ctx, func(ctx context.Context) error {
				err := m.MenuRepo.Create(ctx, &menu)
				if err != nil {
					return err
				}
				//插入按钮和api映射关系
				for _, api := range menu.Apis {
					api.MenuId = menu.ID
				}
				return m.MenuApiRepo.Create(ctx, menu.Apis)
			})
		}

	} else { //不是按钮
		return m.MenuRepo.Create(ctx, &menu)
	}
}

// 修改
func (m *MenuService) Update(ctx context.Context, menu entity.Menu) error {
	//直接通过数据库字段设置唯一性来避免重复插入,但是这里不同层级下的名称是可以相同的
	m.checkName(ctx, menu.Name, menu.Pid, menu.ID)
	//如果不是按钮
	if *menu.Type != "F" {
		if menu.Path == "" {
			panic(respErr.BadRequestErrorWithMsg("path is not null"))
		}
	}
	checkIsFrame(&menu)

	//按钮 单独处理
	if *menu.Type == "F" {
		return m.TranRepo.Exec(ctx, func(ctx context.Context) error {
			err := m.MenuRepo.Update(ctx, &menu)
			if err != nil {
				return err
			}

			//处理menu api表，如果Apis为nil，则表示没有修改此部分
			if menu.Apis != nil {
				err := m.MenuApiRepo.Delete(ctx, menu.ID)
				if err != nil {
					return err
				}
				return m.MenuApiRepo.CreateInBatch(ctx, menu.Apis)
			}
			return nil
		})
	} else {
		return m.MenuRepo.Update(ctx, &menu)
	}
}

func (m *MenuService) Get(ctx context.Context, id uint64) (*entity.Menu, error) {
	//builder := m.menuRepo.NewQueryBuilder().Where("id=?", 20)
	//return m.menuRepo.QueryWithBuilder(builder).First()
	return m.MenuRepo.FindById(ctx, id)
}

func (m *MenuService) Delete(ctx context.Context, id uint64) error {
	return m.TranRepo.Exec(ctx, func(ctx context.Context) error {
		err := m.MenuRepo.Delete(ctx, id)
		if err != nil {
			return err
		}
		//解除菜单和角色关系
		err = m.RoleRepo.UntiedMenu(id)
		if err != nil {
			return err
		}

		//如果是按钮，接触按钮和对应api
		return m.MenuApiRepo.Delete(ctx, id)
	})
}

func (m *MenuService) FindPermissions(ctx context.Context, user *entity.User) ([]string, error) {
	var perms = make([]string, 0)
	if user.IsSuper() {
		perms = append(perms, "*:*:*")
	} else {
		var roleIds = make([]uint64, 0)
		for _, role := range user.Roles {
			roleIds = append(roleIds, role.ID)
		}
		if len(roleIds) > 0 {
			menus, err := m.MenuRepo.FindByRoleIds(ctx, roleIds)
			if err != nil {
				return nil, nil
			}
			for _, menu := range menus {
				p := menu.Permission
				if p != nil {
					perms = append(perms, *p)
				}
			}
		}
	}
	return perms, nil
}

func (m *MenuService) HasChildByID(ctx context.Context, id uint64) (bool, error) {
	return m.MenuRepo.IsExist(ctx, "pid=?", id)
}

// 获取按钮对应的api信息
func (m *MenuService) MenuApiById(ctx context.Context, menuId uint64) ([]*entity.MenuApi, error) {
	return m.MenuApiRepo.FindBy(ctx, "menu_id=?", menuId)
}

// name check
// 如果是修改则需要传id，不同id如何name相同才是重名
func (m *MenuService) checkName(ctx context.Context, name string, pid *uint64, id uint64) {
	var _pid uint64 = 0
	if pid != nil {
		_pid = *pid
	}
	var query = "name=? AND pid=?"
	var args = []interface{}{name, _pid}
	if id != 0 {
		query += " AND id <> ?"
		args = append(args, id)
	}
	isExist, err := m.MenuRepo.IsExist(ctx, query, args...)
	if err != nil {
		panic(respErr.InternalServerErrorWithError(err))
	} else if isExist {
		panic(respErr.BadRequestErrorWithMsg("菜单名重复！"))
	}
}

// 外链检查
func checkIsFrame(menu *entity.Menu) {
	//如果是外链
	if menu.Meta.IsFrame {
		path := strings.ToLower(menu.Path)
		if !(strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")) {
			panic(respErr.BadRequestErrorWithMsg("外链必须以http://或者https://开头"))
		}
	}
}

/**
 * 是否为内链组件(iframe打开链接)
 *
 * @param menu 菜单信息
 * @return 结果
 */
func isInnerLink(menu *entity.Menu) bool {
	path := strings.ToLower(menu.Path)
	return !menu.IsFrame && (strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://"))
}

/**
 * 获取路由地址
 *
 * @param menu 菜单信息
 * @return 路由地址
 */
func getRouterPath(menu *entity.Menu) string {
	routerPath := menu.Path
	// 内链打开外网方式
	if *menu.Pid != 0 && isInnerLink(menu) {
		routerPath = link2RoutePath(routerPath)
	}

	// 非外链并且是一级目录（类型为目录）
	if 0 == *menu.Pid && !menu.IsFrame {
		if consts.TYPE_DIR == *menu.Type {
			if !strings.HasPrefix(menu.Path, "/") {
				routerPath = "/" + menu.Path
			}
		} else if consts.TYPE_MENU == *menu.Type {
			// 非外链并且是一级目录（类型为菜单）
			routerPath = "/"
		}
	}
	return routerPath
}

/**
 * 获取路由名称
 *
 * @param menu 菜单信息
 * @return 路由名称
 */
func getRouteName(menu *entity.Menu) string {
	var routeName = ""
	if menu.IsFrame {
		routeName = menu.Path
	} else if *menu.Pid == 0 && *menu.Type == consts.TYPE_MENU {
		//顶级菜单
		routeName = ""
	} else {
		path := menu.Path
		if isInnerLink(menu) {
			path = link2RoutePath(menu.Path)
		}
		routeName = utils.ToUpperForFirstCharAtSymbolBehind(path, "/")
	}
	return routeName
}

func link2RoutePath(path string) string {
	if path == "" {
		return ""
	}
	//去掉query参数
	url := strings.Split(path, "?")[0]
	//替换 (http:// https:// www) 为""
	regx := regexp.MustCompile(`^(http|https):\/\/[www.]?`)
	part := regx.ReplaceAllString(url, "")
	//端口
	regx2 := regexp.MustCompile(`:\d+`)
	noPort := regx2.ReplaceAllString(part, "")
	//去掉host最后一部分
	domain := strings.Replace(noPort, "#/", "", -1)
	x1 := strings.LastIndex(domain, ".")
	x2 := strings.Index(domain, "/")
	// .替换为/
	return strings.ReplaceAll(domain[0:x1]+domain[x2:], ".", "/")
}
