package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

// 全局日志对象，主要用于公共库没有调用NewLogs的场景
var gLogger *Logger

// 日志库
type Logger struct {
	zapLogger *zap.Logger
	config    LoggerConfig
}

// 创建一个Logger对象
func NewLogger(config LoggerConfig) *Logger {

	gLogger = &Logger{
		config: config,
	}

	// 通过type获取writer对象，默认为终端
	writer := getConsoleWriter()
	config.Mode = strings.ToLower(config.Mode)
	switch config.Mode {
	case "file":
		writer = getFileWriter(config)
	case "kafka":
		writer = getKafkaWriter(config)
	}

	if writer == nil {
		return gLogger
	}

	// 格式相关的配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05") // 修改时间戳的格式
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder                       // 日志级别使用大写显示
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 日志等级
	level := zapcore.InfoLevel
	config.Level = strings.ToLower(config.Level)
	switch config.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}
	core := zapcore.NewCore(encoder, writer, level)

	gLogger.zapLogger = zap.New(core,
		zap.Fields(zap.String("service", config.ServiceName)), // 添加Service字段
		zap.AddCaller(), // 添加调用位置
		zap.AddCallerSkip(2),
		zap.AddStacktrace(zap.ErrorLevel),
	)

	return gLogger
}

// 关闭日志
func (l *Logger) Close() error {

	if l.zapLogger == nil {
		return nil
	}

	var err error
	switch l.config.Mode {
	case "file":
		err = closeFile()
	case "kafka":
		err = closeKafka()
	}

	return err
}
