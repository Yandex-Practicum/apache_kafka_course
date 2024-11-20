package producer

import (
	"log/slog"

	"github.com/apache_kafka_course/module1/go/avro-example/internal/config"
	"github.com/apache_kafka_course/module1/go/avro-example/internal/dto"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
)

type Broker struct {
	producer   *kafka.Producer
	serializer serde.Serializer
	log        *slog.Logger
}

type Response struct {
	UserUUID string
	Err      error
}

var FlushBrokerTimeMs = 100

// New returns kafka producer with schema registry.
func New(cfg *config.Config, log *slog.Logger) (*Broker, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": cfg.Kafka.KafkaURL})
	if err != nil {
		return nil, err
	}

	client, err := schemaregistry.NewClient(schemaregistry.NewConfig(cfg.Kafka.SchemaRegistryURL))
	if err != nil {
		return nil, err
	}
	ser, err := avro.NewSpecificSerializer(client, serde.ValueSerde, avro.NewSerializerConfig())
	if err != nil {
		return nil, err
	}

	// Delivery report handler for produced messages
	go func() {
		for {
			for e := range p.Events() {
				switch e := e.(type) {
				// https://github.com/confluentinc/confluent-kafka-go/blob/master/examples/producer_example/producer_example.go
				case *kafka.Message:
					// The message delivery report, indicating success or
					// permanent failure after retries have been exhausted.
					// Application level retries won't help since the client
					// is already configured to do that.
					if e.TopicPartition.Error != nil {
						log.Error("sending message finished with failure", "err", e.TopicPartition.Error, "key", string(e.Key))
						continue
					}
					log.Debug("sending message finished with success ", "key", string(e.Key))
				case kafka.Error:
					// Generic client instance-level errors, such as
					// broker connection failures, authentication issues, etc.
					//
					// These errors should generally be considered informational
					// as the underlying client will automatically try to
					// recover from any errors encountered, the application
					// does not need to take action on them.
					log.Error("kafka general error", "err", e.Error())
				}
			}
		}
	}()

	return &Broker{
			producer:   p,
			serializer: ser,
			log:        log,
		},
		nil
}

// Close closes serialization agent and kafka producer.
func (b *Broker) Close() {
	b.log.Info("kafka stops")
	b.serializer.Close()
	// https://docs.confluent.io/platform/current/clients/confluent-kafka-go/index.html#hdr-Producer
	// * When done producing messages it's necessary  to make sure all messages are
	// indeed delivered to the broker (or failed),
	// because this is an asynchronous client so some messages may be
	// lingering in internal channels or transmission queues.
	// Call the convenience function `.Flush()` will block code until all
	// message deliveries are done or the provided timeout elapses.
	b.producer.Flush(FlushBrokerTimeMs)
	b.producer.Close()
}

// Send sends serialized message to kafka using schema registry.
func (b *Broker) Send(msg dto.User, topic string, key string) error {
	b.log.Info("sending message", "msg", msg)
	payload, err := b.serializer.Serialize(topic, &msg)
	if err != nil {
		return err
	}
	err = b.producer.Produce(&kafka.Message{
		Key:            []byte(key),
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          payload,
		Headers:        []kafka.Header{{Key: "Course", Value: []byte("Kafka")}},
	}, nil)
	if err != nil {
		return err
	}
	return nil
}
