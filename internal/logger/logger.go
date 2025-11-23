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

// Debug - инкапсуляциия логгера для метода Debug
func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

// Info - инкапсуляциия логгера для метода Info
func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

// Warn - инкапсуляциия логгера для метода Warn
func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

// Error - инкапсуляциия логгера для метода Error
func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

// Fatal - инкапсуляциия логгера для метода Fatal
func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

// Debugf - инкапсуляциия логгера для метода Debug c форматированием
func Debugf(format string, args ...interface{}) {
	globalLogger.Debug(fmt.Sprintf(format, args...))
}

// Infof - инкапсуляциия логгера для метода Info c форматированием
func Infof(format string, args ...interface{}) {
	globalLogger.Info(fmt.Sprintf(format, args...))
}

// Warnf - инкапсуляциия логгера для метода Warn c форматированием
func Warnf(format string, args ...interface{}) {
	globalLogger.Warn(fmt.Sprintf(format, args...))
}

// Errorf - инкапсуляциия логгера для метода Error c форматированием
func Errorf(format string, args ...interface{}) {
	globalLogger.Error(fmt.Sprintf(format, args...))
}

// Fatalf - инкапсуляциия логгера для метода Fatal c форматированием
func Fatalf(format string, args ...interface{}) {
	globalLogger.Fatal(fmt.Sprintf(format, args...))
}
