package app

import (
	"context"
	"errors"

	"github.com/apache_kafka_course/module1/go/avro-example/app/consumer"
	"github.com/apache_kafka_course/module1/go/avro-example/app/producer"
	"github.com/apache_kafka_course/module1/go/avro-example/internal/config"
	"github.com/apache_kafka_course/module1/go/avro-example/internal/logger"
)

var ErrWrongType = errors.New("wrong type")

type StartGetConfigStopper interface {
	Start(ctx context.Context)
	GetConfig() string
	Stop()
}

func Fabric() (StartGetConfigStopper, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}
	log := logger.New(cfg.Env)
	switch cfg.Kafka.Type {
	case "producer":
		return producer.New(cfg, log)
	case "consumer-push":
		return consumer.New(cfg, log)
	case "consumer-pull":
		return consumer.New(cfg, log)
	default:
		return nil, ErrWrongType
	}
}
