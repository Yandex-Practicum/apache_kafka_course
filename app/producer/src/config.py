from pydantic_settings import BaseSettings, SettingsConfigDict


class ConfigBase(BaseSettings):
    model_config = SettingsConfigDict(
        env_file=".env", env_file_encoding="utf-8", extra="ignore"
    )


class KafkaConfig(ConfigBase):
    model_config = SettingsConfigDict(env_prefix="PRODUCER_") 

    KAFKA_BOOTSTRAP_SERVERS: str
    SCHEMA_REGISTRY_SERVER: str
    KAFKA_TOPIC: str
    KAFKA_USER: str
    KAFKA_PASSWORD: str
    KAFKA_CA_PATH: str