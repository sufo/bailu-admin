/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package app

import (
	"bailu/app/core"
	"bailu/app/middleware"
	"bailu/global"
	respErr "bailu/pkg/exception"
	"bailu/pkg/i18n"
	"bailu/pkg/ip2region"
	"bailu/pkg/mq"
	"bailu/pkg/store"
	"bailu/pkg/translate"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/fvbock/endless"
	"github.com/google/gops/agent"
	"os"
	//_log "log"
	"bailu/app/config"
	"bailu/pkg/log"
	"net/http"
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

func Init(ctx context.Context, opts ...Opt) (func(), error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	//加载配置文件
	core.InitViper(o.Conf)
	respErr.Initial() //初始化一些常用响应错误变量
	//db := util.Gorm()
	//di.Add("gorm", db)
	if v := o.MenuFile; v != "" {
		config.Conf.Menu.Path = v
	}
	//初始化日志
	clearLogFunc, err := log.InitLogger()
	if err != nil {
		return nil, err
	}
	log.L.Infof("Start server,#run_mode %s,#version %s,#pid %d", config.Conf.Server.Mode, o.Version, os.Getpid())

	//初始化监视器
	clearMonitor := InitMonitor()

	//初始化全局ip2region
	ip2region.InitIp2Region()

	//初始化全局validator翻译器
	tranErr := translate.InitTrans(config.Conf.Server.Locale)
	if tranErr != nil {
		log.L.Errorf("init translate error: %s", tranErr.Error())
	}

	//init i18n
	languages := map[string]string{
		"en":    "English",
		"zh-CN": "简体中文",
		// "zh-TW": "繁体中文",
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

	//wire inject
	injector, clearFunc, err := BuildInjector(o.WWW)
	if err != nil {
		return nil, err
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
			return nil, err
		}
	}

	//处理定时任务
	go func() {
		//cron.InitFuncJobs()
		injector.Injector2Job.InitFuncJobs()
		injector.Job.Start(ctx)
	}()

	//启动服务后会阻塞， 一些服务启动时的动作不要放在它后面
	InitHTTPServer(injector.Engine)

	return func() {
		<-injector.SSE.ClosedClients //关闭sse
		clearFunc()
		clearMonitor()
		clearLogFunc()
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

func InitMonitor() func() {
	if c := config.Conf.Monitor; c.Enable {
		// ShutdownCleanup set false to prevent automatically closes on os.Interrupt
		// and close agent manually before service shutting down
		err := agent.Listen(agent.Options{Addr: c.Addr, ConfigDir: c.ConfigDir, ShutdownCleanup: false})
		if err != nil {
			log.L.Errorf("Agent monitor exception: %s", err.Error())
		}
		return func() {
			agent.Close()
		}
	}
	return func() {}
}

func InitHTTPServer(handler http.Handler) func() {
	cfg := config.Conf.Server
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	s := endless.NewServer(addr, handler)
	s.WriteTimeout = time.Duration(cfg.WriterTimeout) * time.Second
	s.ReadTimeout = time.Duration(cfg.ReadTimeout) * time.Second
	s.MaxHeaderBytes = 1 << 20 //请求头最大1M
	s.BeforeBegin = func(add string) {
		//_log.Printf("Actual pid is %d", syscall.Getpid())
		log.L.Info("addr:", addr)
	}

	var err error
	if cfg.CertFile != "" && cfg.KeyFile != "" {
		s.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
		err = s.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
	} else {
		err = s.ListenAndServe()
	}
	//err := s.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
	return func() {}
}

func Run(ctx context.Context, opts ...Opt) error {
	global.StartTime = time.Now() //记录启动时间
	clearFunc, err := Init(ctx, opts...)
	defer clearFunc()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

//func Run() exception {
//	state := 1
//	sc := make(chan os.Signal, 1)
//	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
//	clearFunc, err := Init()
//	if err != nil {
//		return err
//	}
//EXIT:
//	for {
//		sig := <-sc
//
//	}
//
//}
