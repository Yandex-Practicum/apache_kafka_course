package di

import (
	"kafka-connect/internal/application/metric"
	"kafka-connect/internal/application/repo"
	"kafka-connect/internal/infra/config"
	"kafka-connect/internal/infra/http/handlers"
	"kafka-connect/internal/infra/http/httpserver"
	"kafka-connect/internal/infra/http/mux"
	"kafka-connect/internal/infra/kafka"
	"kafka-connect/internal/infra/logger"
	"kafka-connect/internal/infra/storage"

	"go.uber.org/fx"
)

func InjectApp() fx.Option {
	return fx.Provide(
		kafka.NewSchema,
		kafka.NewMetricsSender,
		kafka.NewMetricsReceiver,
		config.NewConfig,
		handlers.NewMetricsHandler,
		logger.NewLogger,
		mux.NewMux,
		httpserver.NewHTTPServer,

		storage.NewRedis,
		func(redis *storage.Redis) repo.Repo {
			return redis
		},

		metric.NewProcess,
	)
}
