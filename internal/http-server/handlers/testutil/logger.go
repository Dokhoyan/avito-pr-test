package testutil

import (
	"io"

	"github.com/Dokhoyan/avito-pr-test/internal/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitTestLogger initializes the global logger for tests
func InitTestLogger() {
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(io.Discard), // Discard output in tests
		zapcore.DebugLevel,
	)
	logger.Init(core)
}
