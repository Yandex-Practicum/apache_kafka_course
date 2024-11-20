package app

import (
	"context"
	"errors"
	"lesson1/avro-example/app/consumer"
	"lesson1/avro-example/app/producer"
	"lesson1/avro-example/internal/config"
	"lesson1/avro-example/internal/logger"
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
