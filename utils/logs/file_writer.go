package logs

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap/zapcore"
	"os"
)

// 文件分割对象
var lumberJackLogger *lumberjack.Logger

// 文件写入对象
func getFileWriter(config LoggerConfig) zapcore.WriteSyncer {

	lumberJackLogger = &lumberjack.Logger{
		Filename:   config.FileInfo.Path,
		MaxSize:    2 * 1024, // 每个文件2G,单位:M
		MaxBackups: 7,        // 最多保存7个文件
		MaxAge:     10,       // 最多保存10天
		Compress:   true,     // 压缩
		LocalTime:  true,     // 文件名本地日期
	}

	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(lumberJackLogger),
		zapcore.AddSync(os.Stdout))
}

// 关闭文件
func closeFile() error {

	var err error
	if lumberJackLogger != nil {
		err = lumberJackLogger.Close()
	}

	return err
}
