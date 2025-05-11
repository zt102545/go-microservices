package logs

import (
	"context"
	"fmt"
	"go-microservices/utils/consts"
	"strings"

	"github.com/zeromicro/go-zero/core/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/metadata"
)

// 写入日志
func writeLog(ctx context.Context, logger *Logger, level zapcore.Level, msg string, format ...interface{}) {

	if logger == nil || logger.zapLogger == nil {
		return
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if len(md[consts.CTX_USERID]) > 0 {
			format = append(format, Any(consts.CTX_USERID, md[consts.CTX_USERID][0]))
		}
		if len(md[consts.CTX_TRACEID]) > 0 {
			format = append(format, Any(consts.CTX_TRACEID, md[consts.CTX_TRACEID][0]))
		}
	} else {
		format = append(format, Any(consts.CTX_USERID, ctx.Value(consts.CTX_USERID)))
		format = append(format, Any(consts.CTX_TRACEID, ctx.Value(consts.CTX_TRACEID)))
	}
	msg, fields := formatMsg(msg, format...)
	switch level {
	case zap.DebugLevel:
		logger.zapLogger.Debug(msg, fields...)
	case zap.InfoLevel:
		logger.zapLogger.Info(msg, fields...)
	case zap.WarnLevel:
		logger.zapLogger.Warn(msg, fields...)
	case zap.ErrorLevel:
		if !strings.Contains(msg, "skywalking") {
			logger.zapLogger.Error(msg, fields...)
		}
	}

	if logger.config.Env == service.DevMode || logger.config.Env == service.TestMode {
		fmt.Printf(msg, format...)
	}
}

// 格式化日志数据
func formatMsg(format string, a ...interface{}) (string, []zapcore.Field) {

	var (
		args   []interface{}
		fields []zapcore.Field
	)

	for _, arg := range a {
		if b, ok := arg.(zapcore.Field); ok {
			fields = append(fields, b)
			continue
		}
		args = append(args, arg)
	}

	msg := fmt.Sprintf(format, args...)

	return msg, fields
}

// 设置Arg参数
func Any(key string, data interface{}) zap.Field {

	return zap.Any(key, data)
}

// 设置Flag参数
func Flag(flag string) zap.Field {

	return zap.String("flag", flag)
}

// 全局Debug
func Debug(ctx context.Context, msg string, format ...interface{}) {

	writeLog(ctx, gLogger, zap.DebugLevel, msg, format...)
}

// 全局Info
func Info(ctx context.Context, msg string, format ...interface{}) {

	writeLog(ctx, gLogger, zap.InfoLevel, msg, format...)
}

// 全局Warn
func Warn(ctx context.Context, msg string, format ...interface{}) {
	writeLog(ctx, gLogger, zap.WarnLevel, msg, format...)
}

// 全局Err
func Err(ctx context.Context, msg string, format ...interface{}) {

	writeLog(ctx, gLogger, zap.ErrorLevel, msg, format...)
}
