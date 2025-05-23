package kafka

import (
	"kafka-connect/internal/infra/config"

	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
	"go.uber.org/zap"
)

type Schema struct {
	c     *schemaregistry.Client
	ser   *jsonschema.Serializer
	deser *jsonschema.Deserializer
}

func (s Schema) Serialize(topic string, msg interface{}) ([]byte, error) {
	return s.ser.Serialize(topic, msg)
}

func (s Schema) DeserializeInto(topic string, payload []byte, msg interface{}) error {
	return s.deser.DeserializeInto(topic, payload, msg)
}

func NewSchema(
	conf config.Config,
	l *zap.SugaredLogger,
) *Schema {
	c, err := schemaregistry.NewClient(schemaregistry.NewConfig(conf.SchemaregistryUrl))
	if err != nil {
		l.Fatal("Failed to create schema registry ", zap.Error(err))
	}

	ser, err := jsonschema.NewSerializer(c, serde.ValueSerde, jsonschema.NewSerializerConfig())
	if err != nil {
		l.Fatal("Failed to create ser ", zap.Error(err))
	}

	deser, err := jsonschema.NewDeserializer(c, serde.ValueSerde, jsonschema.NewDeserializerConfig())
	if err != nil {
		l.Fatal("Failed to create deser ", zap.Error(err))
	}

	return &Schema{&c, ser, deser}
}
