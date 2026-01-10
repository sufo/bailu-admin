package log

import (
	"fmt"
	zaprotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sufo/bailu-admin/app/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"time"
)

var (
	level zapcore.Level
	L     *zap.SugaredLogger // Global logger for backward compatibility
)

// InitLogger initializes the global logger `L` and also returns the instance
// and cleanup function for use with dependency injection.
func InitLogger() (*zap.SugaredLogger, func(), error) {
	if ok, _ := pathExists(config.Conf.Zap.Director); !ok {
		fmt.Printf("create %v directory\n", config.Conf.Zap.Director)
		_ = os.Mkdir(config.Conf.Zap.Director, os.ModePerm)
	}

	switch config.Conf.Zap.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
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
	if err != nil {
		return nil, nil, err
	}

	var logger *zap.Logger
	if level <= zap.DebugLevel {
		logger = zap.New(zapCore, zap.AddStacktrace(level))
	} else {
		logger = zap.New(zapCore)
	}

	if config.Conf.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}

	// Assign to the global logger for backward compatibility
	L = logger.Sugar()

	// Also return the instance for dependency injection
	return L, clearFunc, err
}

// getEncoderConfig returns a developer-friendly encoder configuration.
func getEncoderConfig() zapcore.EncoderConfig {
	// Use NewDevelopmentEncoderConfig for a human-readable console output
	encoderConfig := zap.NewDevelopmentEncoderConfig()

	// Customize the time format
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006/01/02 - 15:04:05.000"))
	}

	// Customize the level encoding based on config
	switch config.Conf.Zap.EncodeLevel {
	case "LowercaseLevelEncoder":
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	case "LowercaseColorLevelEncoder":
		encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case "CapitalLevelEncoder":
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	case "CapitalColorLevelEncoder":
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder // Default to color
	}

	encoderConfig.ConsoleSeparator = "  "

	return encoderConfig
}

// getEncoder returns the appropriate encoder based on the config.
func getEncoder() zapcore.Encoder {
	if config.Conf.Zap.Format == "json" {
		return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

// getEncoderCore creates the zapcore.Core.
func getEncoderCore() (zapcore.Core, func(), error) {
	writer, clearFunc, err := GetWriterSyncer()
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return nil, nil, err
	}
	return zapcore.NewCore(getEncoder(), writer, level), clearFunc, err
}

// GetWriterSyncer sets up the log writer (file and/or console).
func GetWriterSyncer() (zapcore.WriteSyncer, func(), error) {
	fileWriter, err := zaprotatelogs.New(
		path.Join(config.Conf.Zap.Director, "%Y-%m-%d.log"),
		zaprotatelogs.WithMaxAge(time.Duration(config.Conf.Zap.CleanDaysAgo*24)*time.Hour),
		zaprotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		return nil, nil, err
	}

	clearFunc := func() {
		if fileWriter != nil {
			_ = fileWriter.Close()
		}
	}

	if config.Conf.Zap.LogInConsole {
		// Correctly write to Stdout, not Stdin
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), clearFunc, nil
	}

	return zapcore.AddSync(fileWriter), clearFunc, nil
}

// pathExists checks if a path exists.
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
