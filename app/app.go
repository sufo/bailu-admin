/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/google/gops/agent"
	"github.com/sufo/bailu-admin/app/config"
	"github.com/sufo/bailu-admin/app/core"
	"github.com/sufo/bailu-admin/app/middleware"
	"github.com/sufo/bailu-admin/global"
	respErr "github.com/sufo/bailu-admin/pkg/exception"
	"github.com/sufo/bailu-admin/pkg/i18n"
	"github.com/sufo/bailu-admin/pkg/ip2region"
	"github.com/sufo/bailu-admin/pkg/mq"
	"github.com/sufo/bailu-admin/pkg/store"
	"github.com/sufo/bailu-admin/pkg/translate"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type options struct {
	Version  string
	WWW      string
	Conf     string
	MenuFile string
}
type Opt func(option *options)

func SetConfigFile(s string) Opt {
	return func(opt *options) {
		opt.Conf = s
	}
}
func SetWWWDir(s string) Opt {
	return func(opt *options) {
		opt.WWW = s
	}
}
func SetMenuFile(m string) Opt {
	return func(opt *options) {
		opt.MenuFile = m
	}
}
func SetVersion(s string) Opt {
	return func(opt *options) {
		opt.Version = s
	}
}

func Init(ctx context.Context, opts ...Opt) (*Injector, func(), error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	//加载配置文件
	core.InitViper(o.Conf)
	respErr.Initial() //初始化一些常用响应错误变量

	if v := o.MenuFile; v != "" {
		config.Conf.Menu.Path = v
	}

	//wire inject
	injector, clearFunc, err := BuildInjector(o.WWW)
	if err != nil {
		return nil, nil, err
	}

	injector.Logger.Infof("Start server,#run_mode %s,#version %s,#pid %d", config.Conf.Server.Mode, o.Version, os.Getpid())

	//初始化监视器
	clearMonitor := InitMonitor(injector.Logger)

	//初始化全局ip2region
	ip2region.InitIp2Region()

	//初始化全局validator翻译器
	tranErr := translate.InitTrans(config.Conf.Server.Locale)
	if tranErr != nil {
		injector.Logger.Errorf("init translate error: %s", tranErr.Error())
	}

	//init i18n
	languages := map[string]string{
		"en":    "English",
		"zh-CN": "简体中文",
	}
	// The default instance initialized directly here
	i18n.Init("app/locales/lang", "en", languages)

	//初始化redis服务和mq
	if config.Conf.Server.UseRedis {
		store.RedisClient = store.NewRedisClient(config.Conf.Store.Redis.DB)
		mq.Publisher = mq.NewPublisher(store.RedisClient, "annmq")
		mq.Consumer = mq.NewRedisConsumer(&mq.RedisStreamConsumerConfig{
			Client: store.RedisClient, StreamName: "annmq", ConsumerGroupName: "trmqcg",
		})
	}

	//配置静态文件夹路径 第一个参数是api，第二个是文件夹路径
	if o.WWW != "" {
		injector.Engine.Use(middleware.WWWMiddleware(o.WWW,
			middleware.AllowPathPrefixSkipper(injector.Router.Prefixes()...)))
	}

	//menu init
	if config.Conf.Menu.Enable && config.Conf.Menu.Path != "" {
		err = injector.MenuSrv.InitData(ctx, config.Conf.Menu.Path)
		if err != nil {
			return nil, nil, err
		}
	}

	//处理定时任务
	go func() {
		//cron.InitFuncJobs()
		injector.Injector2Job.InitFuncJobs()
		injector.Job.Start(ctx)
	}()

	return injector, func() {
		clearFunc()
		clearMonitor()
		//关闭redis连接
		if config.Conf.Server.UseRedis {
			err := store.RedisClient.Close()
			if err != nil {
				fmt.Println(err)
			}
		}
		// 程序结束前关闭数据库链接
		//sqlDb, _ := db.DB()
		//_ = sqlDb.Close()
	}, nil
}

func InitMonitor(logger *zap.SugaredLogger) func() {
	if c := config.Conf.Monitor; c.Enable {
		// ShutdownCleanup set false to prevent automatically closes on os.Interrupt
		// and close agent manually before service shutting down
		err := agent.Listen(agent.Options{Addr: c.Addr, ConfigDir: c.ConfigDir, ShutdownCleanup: false})
		if err != nil {
			logger.Errorf("Agent monitor exception: %s", err.Error())
		}
		return func() {
			agent.Close()
		}
	}
	return func() {}
}

func InitHTTPServer(handler http.Handler, logger *zap.SugaredLogger) {
	cfg := config.Conf.Server
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	srv := &http.Server{
		Addr:           addr,
		Handler:        handler,
		WriteTimeout:   time.Duration(cfg.WriterTimeout) * time.Second,
		ReadTimeout:    time.Duration(cfg.ReadTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	go func() {
		logger.Info("Listening and serving HTTP on ", addr)
		var err error
		if cfg.CertFile != "" && cfg.KeyFile != "" {
			srv.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
			err = srv.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			logger.Fatalf("HTTP server ListenAndServe error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: ", err)
	}
	logger.Info("Server exiting")
}

func Run(ctx context.Context, opts ...Opt) error {
	global.StartTime = time.Now() //记录启动时间
	injector, clearFunc, err := Init(ctx, opts...)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer clearFunc()

	InitHTTPServer(injector.Engine, injector.Logger)

	return nil
}
