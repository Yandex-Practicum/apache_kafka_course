package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	KafkaBootstrapServers      string
	Topic                      string
	KafkaCustomerGroup1        string
	KafkaSessionTimeoutMs      int
	KafkaConsumerPullTimeoutMs int
	SchemaregistryUrl          string
	RedisOpt                   string
}

func NewConfig() Config {
	var err error
	if os.Getenv("IS_TEST_ENV") == "true" {
		err = godotenv.Load("../../.env.test")
	} else {
		err = godotenv.Load(".env")
	}

	if err != nil {
		panic(err)
	}

	return Config{
		KafkaBootstrapServers:      getEnv("KAFKA_BOOTSTRAP_SERVERS", ""),
		Topic:                      getEnv("TOPIC", ""),
		KafkaCustomerGroup1:        getEnv("KAFKA_CUSTOMER_GROUP_1", ""),
		KafkaSessionTimeoutMs:      getEnvInt("KAFKA_SESSION_TIMEOUT_MS", 0),
		KafkaConsumerPullTimeoutMs: getEnvInt("KAFKA_CONSUMER_PULL_TIMEOUT_MS", 0),
		SchemaregistryUrl:          getEnv("SCHEMA_REGISTRY_URL", ""),
		RedisOpt:                   getEnv("REDIS_OPT", ""),
	}
}

func getEnvInt(key string, def int) int {
	v, e := strconv.Atoi(getEnv(key, strconv.Itoa(def)))
	if e != nil {
		return def
	} else {
		return v
	}
}

func getEnv(key string, def string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return def
}
