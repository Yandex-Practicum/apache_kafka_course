package config

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	ErrAbsentConfigFile = errors.New("config file does not exists")
	ErrReadConfigFailed = errors.New("reading config file failed")
)

type KafkaConfig struct {
	KafkaURL          string `yaml:"kafkaUrl" env-required:"true"`
	SchemaRegistryURL string `yaml:"schemaRegistryUrl" env-required:"true"`
	Type              string `yaml:"type" env-required:"true"`
	GroupID           string `yaml:"groupId"`
	Topic             string `yaml:"topic" env-required:"true"`
	Timeout           int    `yaml:"timeout"`
}

type Config struct {
	// without this param will be used "local" as param value
	Env   string      `yaml:"env" env-default:"local"`
	Kafka KafkaConfig `yaml:"kafka"`
}

func (c *Config) String() string {
	return fmt.Sprintf(
		"type: %s, env: %s, kafka url %s, schema registry url %s",
		c.Kafka.Type, c.Env, c.Kafka.KafkaURL, c.Kafka.SchemaRegistryURL,
	)
}

// New loads config.
func New() (*Config, error) {
	cfg := &Config{}
	var err error
	var configPath string
	// path to config yaml file
	flag.StringVar(&configPath, "c", "", "path to config file")
	flag.Parse()
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}
	if configPath != "" {
		cfg, err = LoadByPath(configPath)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}
	return cfg, nil
}

// LoadByPath loads config by path.
func LoadByPath(configPath string) (*Config, error) {
	_, err := os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrAbsentConfigFile
		}
		return nil, fmt.Errorf("LoadByPath stat error: %w", err)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, ErrReadConfigFailed
	}
	return &cfg, nil
}
