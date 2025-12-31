/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package router

import (
	"bailu/app/api/admin/content"
	"bailu/app/api/admin/mine"
	"bailu/app/api/admin/monitor"
	"bailu/app/api/admin/system"
	"bailu/pkg/jwt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var RouterSet = wire.NewSet(wire.Struct(new(Router), "*"), wire.Bind(new(IRouter), new(*Router)))

// 判断某个结构体是否实现了接口
var _ IRouter = (*Router)(nil)

type IRouter interface {
	Register(app *gin.Engine) error
	Prefixes() []string
}

type Router struct {
	TokenProvider *jwt.JwtProvider
	Enforcer      *casbin.SyncedEnforcer
	DictItemApi   *system.DictItemApi
	DictApi       *system.DictApi
	AuthAPI       *system.AuthApi
	UserApI       *system.UserApi
	ProfileApi    *system.ProfileApi
	MenuApi       *system.MenuApi
	DeptApi       *system.DeptApi
	RoleApi       *system.RoleApi
	PostApi       *system.PostApi
	OperApi       *monitor.OperationApi
	ServInfo      *monitor.ServerInfo
	OnlineApi     *monitor.OnlineUserApi
	TaskApi       *monitor.TaskApi
	LoginLogApi   *monitor.LoginLogApi
	ConfigApi     *system.SysConfigApi
	NoticeApi     *system.NoticeApi
	UploadApi     *system.UploadApi
	Event         *monitor.Event
	Message       *mine.Message
	FileApi       *content.FileApi
}

func (r *Router) Prefixes() []string {
	return []string{
		"/api/",
	}
}

func (r *Router) Register(app *gin.Engine) error {
	r.RegisterAPI(app)
	r.RegisterStream(app)
	return nil
}

/**
 * set name for router
 * If the client does not need the API drop-down list,
 * the following code is not required
 */
type RouteMeta struct {
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Path   string `json:"path"`
	Method string `json:"method"`
	Group  string `json:"group"`
}

// all of routes
var ResourceRoutes = make([]RouteMeta, 0)

// RouterGroup 包装了 gin.RouterGroup，并提供了自动记录元数据的方法
type ResourceRouterGroup struct {
	*gin.RouterGroup              // 组合 gin.RouterGroup，继承其所有方法
	groupName        string       // 权限分组名，用于元信息
	metaList         *[]RouteMeta // 指向全局元信息列表的指针
}

// NewRouterGroup 是一个构造函数，用于创建一个新的 RouterGroup
func NewResourceRouterGroup(group *gin.RouterGroup, name string, metaList *[]RouteMeta) *ResourceRouterGroup {
	return &ResourceRouterGroup{
		RouterGroup: group,
		groupName:   name,
		metaList:    metaList,
	}
}

// POST 是对 gin.RouterGroup.POST 的封装
func (p *ResourceRouterGroup) POST(relativePath, name, desc string, handlers ...gin.HandlerFunc) gin.IRoutes {
	// 构造元信息
	meta := RouteMeta{
		Name:   name,
		Desc:   desc,
		Path:   p.BasePath() + relativePath, // 自动拼接路由组的基础路径
		Method: "POST",
		Group:  p.groupName,
	}
	// 添加到全局列表
	*p.metaList = append(*p.metaList, meta)

	// 调用原始的 POST 方法注册路由
	return p.RouterGroup.POST(relativePath, handlers...)
}

// GET 是对 gin.RouterGroup.GET 的封装
func (p *ResourceRouterGroup) GET(relativePath, name, desc string, handlers ...gin.HandlerFunc) gin.IRoutes {
	meta := RouteMeta{
		Name:   name,
		Desc:   desc,
		Path:   p.BasePath() + relativePath,
		Method: "GET",
		Group:  p.groupName,
	}
	*p.metaList = append(*p.metaList, meta)
	return p.RouterGroup.GET(relativePath, handlers...)
}

func (p *ResourceRouterGroup) PUT(relativePath, name, desc string, handlers ...gin.HandlerFunc) gin.IRoutes {
	meta := RouteMeta{
		Name:   name,
		Desc:   desc,
		Path:   p.BasePath() + relativePath,
		Method: "PUT",
		Group:  p.groupName,
	}
	*p.metaList = append(*p.metaList, meta)
	return p.RouterGroup.PUT(relativePath, handlers...)
}
func (p *ResourceRouterGroup) DELETE(relativePath, name, desc string, handlers ...gin.HandlerFunc) gin.IRoutes {
	meta := RouteMeta{
		Name:   name,
		Desc:   desc,
		Path:   p.BasePath() + relativePath,
		Method: "DELETE",
		Group:  p.groupName,
	}
	*p.metaList = append(*p.metaList, meta)
	return p.RouterGroup.DELETE(relativePath, handlers...)
}
func (p *ResourceRouterGroup) PATCH(relativePath, name, desc string, handlers ...gin.HandlerFunc) gin.IRoutes {
	meta := RouteMeta{
		Name:   name,
		Desc:   desc,
		Path:   p.BasePath() + relativePath,
		Method: "PATCH",
		Group:  p.groupName,
	}
	*p.metaList = append(*p.metaList, meta)
	return p.RouterGroup.PATCH(relativePath, handlers...)
}

// Api tree
func RouteTree() []map[string]any {
	var result = make([]map[string]any, 0)
	groupMap := make(map[string][]RouteMeta)
	for _, item := range ResourceRoutes {
		groupMap[item.Group] = append(groupMap[item.Group], item)
	}

	for group, metas := range groupMap {
		var options = make([]map[string]any, 0)
		for _, meta := range metas {
			option := map[string]any{
				"key":    meta.Method + "_" + meta.Path,
				"label":  meta.Name,
				"method": meta.Method,
				"path":   meta.Path,
				"isLeaf": true,
			}
			options = append(options, option)
		}
		result = append(result, map[string]any{"key": group, "label": group, "children": options})
	}
	return result
}
