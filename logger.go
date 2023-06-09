package webhook

import (
	"go.uber.org/zap"
)

// Logger returns a logger for the webhook package.
func Logger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return logger
}
