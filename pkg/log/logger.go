/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package log

import (
	"bailu/app/config"
	"fmt"
	zaprotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"time"
)

var level zapcore.Level

//var instance *zap.SugaredLogger
//var once sync.Once
//var Logger = Instance()
//单例模式
//func Instance() *zap.SugaredLogger {
//	println("xxxxx")
//	once.Do(func() {
//		instance = newLogger()
//	})
//	return instance
//}

var L *zap.SugaredLogger

// func newLogger() (logger *zap.Logger) {
func InitLogger() (func(), error) {
	if ok, _ := pathExists(config.Conf.Zap.Director); !ok {
		fmt.Printf("create %v directory\n", config.Conf.Zap.Director)
		_ = os.Mkdir(config.Conf.Zap.Director, os.ModePerm)
	}
	// 初始化配置文件的Level
	//switch global.Config.Zap.Level {
	switch config.Conf.Zap.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "exception":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
	zapCore, clearFunc, err := getEncoderCore()

	var logger *zap.Logger
	if level == zap.DebugLevel || level == zap.ErrorLevel {
		logger = zap.New(zapCore, zap.AddStacktrace(level))
	} else {
		logger = zap.New(zapCore)
	}
	//if global.Config.Zap.ShowLine {
	if config.Conf.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	//zap.ReplaceGlobals(logg) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	//zap.ReplaceGlobals(logger)
	L = logger.Sugar()
	return clearFunc, err
}

func getEncoderConfig() (encoderConf zapcore.EncoderConfig) {
	//if err := global.VIPER.Unmarshal(&config); err != nil {
	//	fmt.Println(err)
	//	return
	//}
	encoderConf = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "log",
		CallerKey:      "caller",
		StacktraceKey:  config.Conf.Zap.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     CustomerTimerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
	}

	//switch global.Config.Zap.EncodeLevel {
	switch config.Conf.Zap.EncodeLevel {
	case "LowercaseLevelEncoder": // 小写编码器(默认)
		encoderConf.EncodeLevel = zapcore.LowercaseLevelEncoder
	case "LowercaseColorLevelEncoder": // 小写编码器带颜色
		encoderConf.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case "CapitalLevelEncoder": // 大写编码器
		encoderConf.EncodeLevel = zapcore.CapitalLevelEncoder
	case "CapitalColorLevelEncoder": // 大写编码器带颜色
		encoderConf.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		encoderConf.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return encoderConf
}

// 编码器 (写入日志格式)
func getEncoder() zapcore.Encoder {
	if config.Conf.Zap.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

/**
 * zapcore.Core需要三个配置——Encoder，WriteSyncer，LogLevel。
 */
func getEncoderCore() (zapcore.Core, func(), error) {
	// 使用file-rotatelogs进行日志分割
	writer, clearFunc, err := GetWriterSyncer()
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return nil, nil, err
	}
	return zapcore.NewCore(getEncoder(), writer, level), clearFunc, err
}

// 自定义日志输出时间格式
func CustomerTimerEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(config.Conf.Zap.Prefix + "2006/01/02 - 15:04:05.000"))
}

/**
 * 指定日志写到哪里去。我们使用zapcore.AddSync()函数并且将打开的文件句柄传进去。
 * file-rotate
 */
func GetWriterSyncer() (zapcore.WriteSyncer, func(), error) {
	fileWriter, err := zaprotatelogs.New(
		path.Join(config.Conf.Zap.Director, "%Y-%m-%d.log"),
		//zaprotatelogs.WithLinkName(global.Config.Zap.LinkName),
		zaprotatelogs.WithMaxAge(7*24*time.Hour),
		zaprotatelogs.WithRotationTime(24*time.Hour),
	)
	clearFunc := func() {
		if fileWriter != nil {
			fileWriter.Close()
		}
	}
	if config.Conf.Zap.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdin), zapcore.AddSync(fileWriter)), clearFunc, err
	}
	return zapcore.AddSync(fileWriter), clearFunc, err
}

// 路径是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
