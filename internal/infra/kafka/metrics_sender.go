package kafka

import (
	"context"
	"kafka-connect/internal/infra/config"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MetricsSender struct {
	p   *kafka.Producer
	c   config.Config
	l   *zap.SugaredLogger
	sch *Schema
}

type ListMsg struct {
	Key string
	Msg interface{}
}

func (ms MetricsSender) SendMetrics(done <-chan struct{}, in <-chan ListMsg) {
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				select {
				case <-done:
				case e := <-ms.p.Events():
					switch ev := e.(type) {
					case *kafka.Message:
						m := ev
						if m.TopicPartition.Error != nil {
							ms.l.Warn("Delivery failed ", m.TopicPartition.Error)
						} else {
							ms.l.Info("Delivered message ", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
						}
					case kafka.Error:
						ms.l.Fatal("Error delivered ", ev)
					default:
						ms.l.Info("Ignored event ", ev)
					}
				case lm := <-in:
					val, err := ms.sch.Serialize(ms.c.Topic, lm.Msg)
					if err != nil {
						ms.l.Fatal("Error Serialize ", err)
					}
					err = ms.p.Produce(&kafka.Message{
						TopicPartition: kafka.TopicPartition{Topic: &ms.c.Topic, Partition: kafka.PartitionAny},
						Key:            []byte(lm.Key),
						Value:          val,
						Headers:        []kafka.Header{{Key: "TEST_HEADER", Value: []byte("TEST_HEADER_VAL")}},
					}, nil)
					if err != nil {
						ms.l.Fatal("Error delivered ", err)
					}
				}
			}
		}
	}()
}

func NewMetricsSender(
	lc fx.Lifecycle,
	c config.Config,
	l *zap.SugaredLogger,
	sch *Schema,
) *MetricsSender {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": c.KafkaBootstrapServers,
		"acks":              "all",
	})
	if err != nil {
		l.Fatal("Can't create producer: ", err.Error())
	}
	l.Info("Producer successfully created")

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if !p.IsClosed() {
				p.Flush(1000)
				p.Close()
			}
			return nil
		},
	})

	return &MetricsSender{p, c, l, sch}
}
