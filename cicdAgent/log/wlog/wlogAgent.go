package wlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
)

var (
	Debug func(msg string, fields ...zap.Field)
	Info  func(msg string, fields ...zap.Field)
	Warn  func(msg string, fields ...zap.Field)
	Error func(msg string, fields ...zap.Field)
	Fatal func(msg string, fields ...zap.Field)
)

func InitLogger(logFile string) {
	var (
		_logger *zap.Logger
		once    sync.Once
	)
	once.Do(func() {
		_logger = logConfigDeploy(logFile)
	})

	Debug = _logger.Debug
	Info = _logger.Info
	Warn = _logger.Warn
	Error = _logger.Error
	Fatal = _logger.Fatal

}

// LogConfig 日志配置
func logConfigDeploy(logFile string) *zap.Logger {

	//filePath := os.Getenv("GIN_LOG_PATH")
	//if filePath == "" {
	//	filePath = "./logs/server.log"
	//}
	filePath := logFile

	loggerFileWriter := lumberjack.Logger{
		Filename:   filePath, // 日志文件路径
		MaxSize:    10,       // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 5,        // 日志文件最多保存多少个备份
		MaxAge:     7,        // 文件最多保存多少天
		Compress:   true,     // 是否压缩
	}

	// 日志文件输出配置
	fileEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,                        // 全大写日志等级标识
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"), // 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 终端输出配置
	stdEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	//fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)  // 设置日志文件内容编码格式：json格式
	fileEncoder := zapcore.NewConsoleEncoder(fileEncoderConfig) // 设置日志文件内容编码格式：正常格式
	stdEncoder := zapcore.NewConsoleEncoder(stdEncoderConfig)   // 终端格式

	// 创建写入的目标 writer
	fileWriter := zapcore.NewMultiWriteSyncer(zapcore.AddSync(&loggerFileWriter))
	stdWriter := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))

	// 设置同时写入文件和终端
	// logWriter := zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook), zapcore.AddSync(os.Stdout))

	// 定义日志级别
	debugLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.DebugLevel
	})

	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.InfoLevel
	})

	// 组合所有配置，分别设置写入文件的参数和显示到终端的参数
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, fileWriter, infoLevel),
		zapcore.NewCore(stdEncoder, stdWriter, debugLevel),
	)
	caller := zap.AddCaller()
	development := zap.Development()
	return zap.New(core, caller, development)
	//if os.Getenv("GIN_DEBUG") == "true" {
	//	// 开启开发模式，堆栈跟踪
	//	caller := zap.AddCaller()
	//	// // 开启文件及行号
	//	development := zap.Development()
	//	// 构造日志对象
	//	return zap.New(core, caller, development)
	//} else {
	//	return zap.New(core)
	//}
}
