/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sufo/bailu-admin/app/config"
	"github.com/sufo/bailu-admin/app/middleware"
	"github.com/sufo/bailu-admin/global"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (r *Router) RegisterAPI(app *gin.Engine) {
	g := app.Group("/api")
	//默认页面
	RegisterWelcome(app)

	//本地化，限流
	g.Use(middleware.LocaleMiddleware(), middleware.RateLimiterMiddleware())
	//无需登录验证，无需 RBAC 权限验证
	{

		g.GET("captcha", r.AuthAPI.GetCaptcha)
		g.POST("login", r.AuthAPI.Login)
		g.GET("sms", r.AuthAPI.SendSMS)
		g.GET("sms-login", r.AuthAPI.SMSLogin)

		g.POST("user/register", r.UserApI.Register)
		g.GET("user/username/:username", r.UserApI.NameExist)
		g.GET("user/phone/:dialCode/:phone", r.UserApI.PhoneExist)
	}

	g.Use(middleware.AuthMiddleware(r.TokenProvider,
		middleware.AllowPathPrefixSkipper("/admin/auth/"),
	))

	//操作日志 //要放在AuthMiddleware路由之后
	if config.Conf.OperLog.Enable {
		g.Use(middleware.OperationMiddleware(
			r.OperApi.OperSrv,
			middleware.AllowPathPrefixSkipper("/admin/login", "/admin/captcha", "/admin/oper"),
		))
	}

	//解锁屏幕
	g.POST("unlock", r.AuthAPI.UnLock)
	//casbin
	g.Use(middleware.CasbinMiddleware(r.Enforcer))

	_gUser := g.Group("/user")
	gUser := NewResourceRouterGroup(_gUser, "用户", &ResourceRoutes)
	{
		gUser.GET("info", "用户信息", "获取当前用户信息", r.UserApI.GetInfo)
		gUser.POST("logout", "退出登录", "", r.UserApI.Logout)
		gUser.GET("", "用户列表", "", r.UserApI.Index)
		gUser.PUT("", "编辑用户", "", r.UserApI.Edit)
		gUser.POST("", "创建用户", "", r.UserApI.Create)
		gUser.PATCH("status", "启用/禁用用户", "", r.UserApI.Status)
		gUser.PATCH("resetPassword", "管理员重置用户密码", "通过管理员来直接重置用户密码", r.UserApI.ResetPwd)
		gUser.PUT("resetPasswordBySmsCode", "短信验证修改用户密码", "", r.UserApI.ResetPwdBySMSCode)
		gUser.DELETE(":userIds", "批量删除用户", "", r.UserApI.Destroy)

		//profile
		gUser.POST("avatar", "上传头像", "上传头像", r.ProfileApi.UploadAvatar)
		gUser.PUT("profile", "修改个人信息", "", r.ProfileApi.Edit)
		gUser.PATCH("changePwd", "修改密码", "", r.ProfileApi.ChangePwd)
	}

	//g.Use(middleware.CasbinMiddleware(a.CasbinEnforcer,
	//	middleware.AllowPathPrefixSkipper("/admin/v1/pub"),
	//))

	_gMenu := g.Group("/menu")
	gMenu := NewResourceRouterGroup(_gMenu, "菜单", &ResourceRoutes)
	{
		gMenu.GET("", "菜单列表", "", r.MenuApi.Index)
		gMenu.GET("menus", "获取菜单树(不包含按钮)", "", r.MenuApi.Menus)
		gMenu.GET("routes", "获取动态路由", "", r.MenuApi.Routes)
		gMenu.GET(":menuId", "获取子菜单树", "", r.MenuApi.SubMenus)
		gMenu.DELETE(":menuId", "删除菜单", "", r.MenuApi.Destroy)
		gMenu.POST("", "创建菜单", "", r.MenuApi.Create)
		gMenu.PATCH("", "编辑菜单", "", r.MenuApi.Edit)
		gMenu.GET("tree", "菜单权限(角色模块使用)", "", r.MenuApi.Tree)
		gMenu.GET("tree/:roleId", "角色菜单树", "", r.MenuApi.TreeSelect)
	}

	_gDept := g.Group("/dept")
	gDept := NewResourceRouterGroup(_gDept, "部门", &ResourceRoutes)
	{
		gDept.GET("", "部门列表", "", r.DeptApi.Index)
		gDept.DELETE(":deptId", "删除部门", "", r.DeptApi.Destroy)
		gDept.POST("", "创建部门", "", r.DeptApi.Create)
		gDept.PUT("", "编辑部门", "", r.DeptApi.Edit)
		gDept.GET("tree", "部门树", "", r.DeptApi.Tree)
		gDept.GET("tree/:roleId", "角色对应的部门树", "", r.DeptApi.TreeSelect)
	}

	_gRole := g.Group("/role")
	gRole := NewResourceRouterGroup(_gRole, "角色", &ResourceRoutes)
	{
		gRole.GET("", "角色列表", "", r.RoleApi.Index)
		gRole.POST("", "创建角色", "", r.RoleApi.Create)
		gRole.PUT("", "编辑角色", "", r.RoleApi.Edit)
		gRole.DELETE(":roleIds", "删除角色", "", r.RoleApi.Destroy)
		gRole.PATCH("status", "启用/禁用角色", "", r.RoleApi.Status)
		gRole.PATCH("dataScope", "数据权限", "", r.RoleApi.DataScope)
		gRole.GET("options", "角色下拉选项列表", "", r.RoleApi.Options)
		gRole.GET(":id", "查询角色", "", r.RoleApi.Detail)
		gRole.DELETE("authUser", "取消用户授权", "", r.RoleApi.CancelAuthUsers)
		gRole.POST("authUser", "用户授权", "", r.RoleApi.AuthUsers)
	}

	_gPost := g.Group("/post")
	gPost := NewResourceRouterGroup(_gPost, "岗位", &ResourceRoutes)
	{
		gPost.GET("", "岗位列表", "", r.PostApi.Index)
		gPost.POST("", "创建岗位", "", r.PostApi.Create)
		gPost.PUT("", "编辑岗位", "", r.PostApi.Edit)
		gPost.DELETE(":ids", "删除岗位", "", r.PostApi.Destroy)
		gPost.GET("options", "岗位下拉列表", "", r.PostApi.Options)
	}

	_gDict := g.Group("/dict")
	gDict := NewResourceRouterGroup(_gDict, "字典", &ResourceRoutes)
	{
		gDict.GET("", "字典列表", "", r.DictApi.Index)
		gDict.POST("", "创建字典", "", r.DictApi.Create)
		gDict.PUT("", "编辑字典", "", r.DictApi.Edit)
		gDict.DELETE(":codes", "删除字典", "", r.DictApi.Destroy)

		gDict.GET(":code", "字典项列表", "", r.DictItemApi.Index)
	}

	_gDictItem := g.Group("/dictItem")
	gDictItem := NewResourceRouterGroup(_gDictItem, "字典项", &ResourceRoutes)
	{
		gDictItem.POST("", "创建字典项", "", r.DictItemApi.Create)
		gDictItem.PUT("", "编辑字典项", "", r.DictItemApi.Edit)
		gDictItem.DELETE(":ids", "删除字典项", "", r.DictItemApi.Destroy)
		gDictItem.GET("options", "字典项下拉列表", "", r.DictItemApi.Options)
		gDictItem.PATCH("status", "编辑字典项状态", "", r.DictItemApi.Status)
	}

	_gSysConfig := g.Group("/config")
	gSysConfig := NewResourceRouterGroup(_gSysConfig, "配置", &ResourceRoutes)
	{
		gSysConfig.GET("", "系统参数配置", "", r.ConfigApi.Index)
		gSysConfig.POST("", "创建系统参数", "", r.ConfigApi.Create)
		gSysConfig.PUT("", "编辑系统参数", "", r.ConfigApi.Edit)
		gSysConfig.DELETE(":ids", "删除系统参数", "", r.ConfigApi.Destroy)
		gSysConfig.PATCH("status", "启用/禁用配置", "", r.ConfigApi.Status)
	}

	//monitor
	_gOper := g.Group("/oper")
	gOper := NewResourceRouterGroup(_gOper, "操作日志", &ResourceRoutes)
	{
		gOper.GET("", "操作日志列表", "", r.OperApi.Index)
		gOper.DELETE(":ids", "删除操作日志", "", r.OperApi.Destroy)
	}

	_gServInfo := g.Group("/server")
	gServInfo := NewResourceRouterGroup(_gServInfo, "系统信息", &ResourceRoutes)
	{
		gServInfo.GET("", "服务器信息", "", r.ServInfo.GetServerInfo)
	}

	_gOnline := g.Group("/online")
	gOnline := NewResourceRouterGroup(_gOnline, "在线用户", &ResourceRoutes)
	{
		gOnline.GET("", "在线用户列表", "", r.OnlineApi.Index)
		gOnline.DELETE(":ids", "剔出用户", "", r.OnlineApi.KickOut)
	}

	//task
	_gTask := g.Group("task")
	gTask := NewResourceRouterGroup(_gTask, "定时任务", &ResourceRoutes)
	{
		gTask.GET("", "任务列表", "", r.TaskApi.Index)
		gTask.POST("", "创建任务", "", r.TaskApi.Create)
		gTask.PUT("", "编辑任务", "", r.TaskApi.Edit)
		gTask.DELETE(":ids", "删除任务", "", r.TaskApi.Destroy)
		gTask.GET(":id", "任务详情", "", r.TaskApi.Detail)
		gTask.GET(":id/logs", "任务执行日志", "", r.TaskApi.Logs)
		gTask.GET("log/:ids", "删除任务日志", "", r.TaskApi.DestroyLogs)
		gTask.POST("invoke/:id", "任务立即执行", "", r.TaskApi.Exec)
		gTask.GET("jobs", "函数任务列表", "", r.TaskApi.FuncJobs)
		gTask.PATCH(":id/:status", "修改任务状态", "", r.TaskApi.Status)
	}
	_gLoginLog := g.Group("/loginLog")
	gLoginLog := NewResourceRouterGroup(_gLoginLog, "登录日志", &ResourceRoutes)
	{
		gLoginLog.GET("", "登录日志列表", "", r.LoginLogApi.Index)
		gLoginLog.GET("findByUsername", "通过用户名登录日志", "", r.LoginLogApi.FindByName)
		gLoginLog.DELETE(":ids", "删除登录日志", "", r.LoginLogApi.Destroy)
		gLoginLog.DELETE("clean", "清空登录日志", "", r.LoginLogApi.Clean)
	}

	_gNotice := g.Group("/notice")
	gNotice := NewResourceRouterGroup(_gNotice, "通知公告", &ResourceRoutes)
	{
		gNotice.GET("", "通知列表", "", r.NoticeApi.Index)
		gNotice.POST("", "创建通知", "", r.NoticeApi.Create)
		gNotice.PUT("", "编辑通知", "", r.NoticeApi.Edit)
		gNotice.DELETE(":ids", "删除通知", "", r.NoticeApi.Destroy)
		gNotice.PATCH("release/:id", "发布通知", "", r.NoticeApi.Release)
		gNotice.PATCH("revoke/:id", "撤销通知", "", r.NoticeApi.Revoke)
	}

	gUpload := g.Group("/upload")
	{
		gUpload.POST("", r.UploadApi.Upload)
	}

	gMine := g.Group("/mine")
	{
		gMine.GET("message/unread", r.Message.Unread)
		gMine.GET(":msgType/unread_count", r.Message.UnreadCount)
		gMine.PUT(":msgType/read_all", r.Message.ReadAll)
		gMine.PUT(":msgType/read/:id", r.Message.Read)
		gMine.DELETE(":msgType/:ids", r.Message.Destroy)
		gMine.DELETE(":msgType/clear", r.Message.Clear)
	}

	_gFile := g.Group("/file")
	gFile := NewResourceRouterGroup(_gFile, "文件", &ResourceRoutes)
	{
		gFile.GET("", "文件列表", "", r.FileApi.Index)
		gFile.POST("", "新增文件", "", r.FileApi.Create)
		gFile.DELETE(":ids", "删除文件", "", r.FileApi.Destroy)

		gFile.GET("category", "文件分类查询", "", r.FileApi.Category)
		gFile.POST("category", "删除文件分类", "", r.FileApi.CategorySave)
		gFile.PUT("category", "编辑文件分类", "", r.FileApi.CategoryEdit)
		gFile.DELETE("category/:ids", "删除文件分类", "", r.FileApi.CategoryDestroy)
	}
}

func RegisterWelcome(app *gin.Engine) {
	//默认页面
	var files []string
	filepath.Walk(global.Root+"/public", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			files = append(files, path)
		}
		return nil
	})
	app.LoadHTMLFiles(files...)
	//如果之加载一个文件
	//app.LoadHTMLFiles("public/index.html")
	app.GET("/api", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"title": strings.ToUpper(config.Conf.Server.Name)})
	})
}
