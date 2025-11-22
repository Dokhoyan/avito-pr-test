package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

// Init - инкапсуляциия логгера для метода Init
func Init(core zapcore.Core, options ...zap.Option) {
	globalLogger = zap.New(core, options...)
}

// Debugf - инкапсуляциия логгера для метода Debug
func Debugf(format string, args ...interface{}) {
	globalLogger.Debug(fmt.Sprintf(format, args...))
}

// Infof - инкапсуляциия логгера для метода Info
func Infof(format string, args ...interface{}) {
	globalLogger.Info(fmt.Sprintf(format, args...))
}

// Warnf - инкапсуляциия логгера для метода Warn
func Warnf(format string, args ...interface{}) {
	globalLogger.Warn(fmt.Sprintf(format, args...))
}

// Errorf - инкапсуляциия логгера для метода Error
func Errorf(format string, args ...interface{}) {
	globalLogger.Error(fmt.Sprintf(format, args...))
}

// Fatalf - инкапсуляциия логгера для метода Fatal
func Fatalf(format string, args ...interface{}) {
	globalLogger.Fatal(fmt.Sprintf(format, args...))
}