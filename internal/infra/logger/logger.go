package logger

import "go.uber.org/zap"

func NewLogger() *zap.SugaredLogger {
	var sugar zap.SugaredLogger

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	sugar = *logger.Sugar()

	return &sugar
}
