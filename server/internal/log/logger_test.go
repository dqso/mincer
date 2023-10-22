package log

import (
	"log/slog"
	"testing"
)

func TestDeferLogger(t *testing.T) {
	logger := New()

	defer func() {
		logger.Info("graceful shutdown")
		logger.Debug("with debug")
	}()

	logger.Debug("no debug")

	logger = NewWithConfig(stdConfig{
		level: slog.LevelDebug,
	})

	logger.Info("application started")
}
