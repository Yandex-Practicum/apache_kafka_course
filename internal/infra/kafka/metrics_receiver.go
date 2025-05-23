package kafka

import (
	"kafka-connect/internal/domain"
	"kafka-connect/internal/infra/config"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"go.uber.org/zap"
)

type MetricsReceiver struct {
	cons *kafka.Consumer
	conf config.Config
	l    *zap.SugaredLogger
	sch  *Schema
}

func (mr MetricsReceiver) ReceiveMetrics(done <-chan struct{}, commit <-chan bool) <-chan domain.List {
	out := make(chan domain.List)

	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			default:
				ev := mr.cons.Poll(mr.conf.KafkaConsumerPullTimeoutMs)
				if ev == nil {
					continue
				}
				switch e := ev.(type) {
				case *kafka.Message:
					mr.l.Info("Get message")
					data := domain.List{}
					err := mr.sch.DeserializeInto(mr.conf.Topic, e.Value, &data)
					if err != nil {
						mr.l.Warn(err.Error())
					} else {
						out <- data
						isCommit := <-commit
						if isCommit {
							_, _ = mr.cons.Commit()
						}
					}
				case kafka.Error:
					mr.l.Warn("Error: ", e)
				default:
					mr.l.Warn("Some event: ", e)
				}
			}
		}
	}()

	return out
}

func NewMetricsReceiver(
	conf config.Config,
	l *zap.SugaredLogger,
	sch *Schema,
) *MetricsReceiver {
	cons, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  conf.KafkaBootstrapServers,
		"group.id":           conf.KafkaCustomerGroup1,
		"session.timeout.ms": conf.KafkaSessionTimeoutMs,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": "false",
	})
	if err != nil {
		l.Fatal("Error creating consumer: ", err)
	}

	l.Info("Consumer created")
	err = cons.SubscribeTopics([]string{conf.Topic}, nil)
	if err != nil {
		l.Fatal("Subscribe error:", err)
	}

	return &MetricsReceiver{cons, conf, l, sch}
}
