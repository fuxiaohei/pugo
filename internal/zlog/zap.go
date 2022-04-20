package zlog

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	_sugaredLogger *zap.SugaredLogger
	logTmFmtWithMS = "2006/01/02 15:04:05"
	_debug         bool
)

func Init(debug bool) {
	_debug = debug
	writer := zapcore.AddSync(os.Stdout)

	encoderConfig := zap.NewProductionEncoderConfig()

	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(logTmFmtWithMS))
	}

	encoderConfig.EncodeTime = customTimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	level := zapcore.InfoLevel
	if debug {
		level = zapcore.DebugLevel
	}
	core := zapcore.NewCore(encoder, writer, level)
	logger := zap.New(core)
	_sugaredLogger = logger.Sugar()
}

func loadZap() {
	if _sugaredLogger == nil {
		Init(false)
	}
}

func Debug(template string, args ...interface{}) {
	loadZap()
	_sugaredLogger.Debugw(template, args...)
}

func Info(template string, args ...interface{}) {
	loadZap()
	_sugaredLogger.Infow(template, args...)
}
func Warn(template string, args ...interface{}) {
	loadZap()
	_sugaredLogger.Warnw(template, args...)
}
func Error(template string, args ...interface{}) {
	loadZap()
	_sugaredLogger.Errorw(template, args...)
}

func Fatal(template string, args ...interface{}) {
	loadZap()
	_sugaredLogger.Fatalw(template, args...)
}
