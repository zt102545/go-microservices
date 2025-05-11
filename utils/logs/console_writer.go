package logs

import (
	"go.uber.org/zap/zapcore"
	"os"
)

// 终端写入对象
func getConsoleWriter() zapcore.WriteSyncer {

	return zapcore.AddSync(os.Stdout)
}
