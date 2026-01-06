/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc initial Router
 */

package core

import (
	"bailu/app/config"
	"bailu/app/domain/resp"
	"bailu/app/middleware"
	"bailu/app/router"
	"bailu/global/consts"
	"bailu/pkg/translate"
	"github.com/LyricTian/gzip"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// 2.如果是有nginx前置做代理，基本不需要gin框架记录访问日志
func InitRouter(r router.IRouter, logger *zap.SugaredLogger, www string) *gin.Engine {
	gin.SetMode(config.Conf.Server.Mode)
	app := gin.New()
	// 限制表单上传大小 MB，默认为32MB
	app.MaxMultipartMemory = int64(config.Conf.Upload.MaxMultipartMemory << 20)

	app.NoMethod(middleware.NoMethodHandler())
	// Router.Use(middleware.LoadTls())  // 打开就能玩https了

	prefixes := r.Prefixes()

	//recovery
	app.Use(middleware.RecoveryMiddleware(true, translate.Trans))

	// trace id
	app.Use(middleware.TraceMiddleware())

	//log
	if config.Conf.Server.Mode != "release" {
		app.Use(middleware.LoggerMiddleware(logger,
			middleware.AllowPathPrefixNoSkipper(prefixes...),
			middleware.AllowPathPrefixSkipper("/api/upload", "/api/file"),
		))
	}

	//CORS
	if config.Conf.CORS.Enable {
		app.Use(middleware.CORSMiddleware())
	}

	//GZIP
	if config.Conf.GZIP.Enable {
		app.Use(gzip.Gzip(gzip.BestCompression,
			gzip.WithExcludedExtensions(config.Conf.GZIP.ExcludedExt),
			gzip.WithExcludedPaths(config.Conf.GZIP.ExcludedPaths),
		))
	}

	//指定静态文件目录 (如上传的文件,文件下载/预览)
	app.StaticFS("/upload", http.Dir(config.Conf.Local.Dir))

	//注册路由
	_ = r.Register(app)

	// api列表
	//tree, err := route.RouteTree(app, "/api/tree", "/assets", "/swagger", "/stream")
	//if err != nil {
	//	panic(err)
	//}
	tree := router.RouteTree()
	api := app.Group("/api")
	{
		_r := r.(*router.Router)
		api.Use(middleware.LocaleMiddleware(), middleware.AuthMiddleware(_r.TokenProvider), middleware.RateLimiterMiddleware()).
			GET("tree", func(c *gin.Context) {
				resp.OKWithData(c, tree)
			})
	}

	// Swagger, 上产环境屏蔽
	if config.Conf.Swagger && config.Conf.Server.Mode != consts.MODE_RELEASE {
		app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		//app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, func(c *ginSwagger.Config) {
		//	c.DeepLinking = false
		//}))
	}

	//
	// --- 3. 前端托管兼容性处理 ---
	registerWWW(app, www)

	noRoute(app, www)

	return app
}

func registerWWW(app *gin.Engine, www string) {
	//前端托管兼容性处理
	if www != "" {
		// 场景 A：后端托管前端 (内嵌部署)
		// 这里的路径处理要非常小心，避免覆盖 API
		app.Static("assets", filepath.Join(www, "assets"))
		app.StaticFile("/favicon.ico", filepath.Join(www, "favicon.ico"))

		// 根路径处理
		app.GET("/", func(c *gin.Context) {
			c.File(filepath.Join(www, "index.html"))
		})
	}
}

func noRoute(app *gin.Engine, www string) {
	app.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 如果是 API 请求找不到，永远返回 JSON 404
		if strings.HasPrefix(path, "/api/") {
			resp.NotFound(c)
			return
		}

		// 场景 A 的补充：如果是后端托管前端，处理 History 模式
		if www != "" {
			indexPath := filepath.Join(www, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				c.File(indexPath)
				return
			}
		}

		// 场景 B：前端独立部署 (Nginx)
		// 此时后端 NoRoute 只会收到非 /api 开头的异常请求
		// 直接返回标准 404 即可，因为前端路由分发由 Nginx 处理了
		c.Status(http.StatusNotFound)
	})
}
