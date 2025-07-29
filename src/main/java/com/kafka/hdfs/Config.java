package com.kafka.hdfs;

import org.apache.commons.configuration2.Configuration;
import org.apache.commons.configuration2.builder.fluent.Configurations;
import org.apache.commons.configuration2.ex.ConfigurationException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;

/**
 * Класс для управления конфигурацией приложения
 */
public class Config {
    private static final Logger logger = LoggerFactory.getLogger(Config.class);

    // Kafka настройки
    public static final String KAFKA_BOOTSTRAP_SERVERS = "kafka.bootstrap.servers";
    public static final String KAFKA_TOPIC = "kafka.topic";
    public static final String KAFKA_GROUP_ID = "kafka.group.id";
    public static final String KAFKA_AUTO_OFFSET_RESET = "kafka.auto.offset.reset";
    public static final String KAFKA_ENABLE_AUTO_COMMIT = "kafka.enable.auto.commit";

    // HDFS настройки
    public static final String HDFS_OUTPUT_PATH = "hdfs.output.path";
    public static final String HDFS_NAMENODE_URL = "hdfs.namenode.url";

    // Настройки производительности
    public static final String BATCH_SIZE = "batch.size";
    public static final String FLUSH_INTERVAL_MS = "flush.interval.ms";
    public static final String POLL_TIMEOUT_MS = "poll.timeout.ms";

    private Configuration config;

    public Config(String configFile) throws ConfigurationException {
        Configurations configs = new Configurations();
        this.config = configs.properties(new File(configFile));
        logger.info("Конфигурация загружена из файла: {}", configFile);
    }

    public Config() {
        // Конструктор по умолчанию с настройками по умолчанию
        setDefaults();
    }

    private void setDefaults() {
        // Устанавливаем значения по умолчанию
        setDefault(KAFKA_BOOTSTRAP_SERVERS, "localhost:9092");
        setDefault(KAFKA_TOPIC, "test-topic");
        setDefault(KAFKA_GROUP_ID, "kafka-hdfs-consumer");
        setDefault(KAFKA_AUTO_OFFSET_RESET, "earliest");
        setDefault(KAFKA_ENABLE_AUTO_COMMIT, "false");

        setDefault(HDFS_OUTPUT_PATH, "/kafka-data");
        setDefault(HDFS_NAMENODE_URL, "hdfs://namenode:9000");

        setDefault(BATCH_SIZE, "100");
        setDefault(FLUSH_INTERVAL_MS, "5000");
        setDefault(POLL_TIMEOUT_MS, "1000");
    }

    private void setDefault(String key, String value) {
        if (config == null || !config.containsKey(key)) {
            if (config == null) {
                // Создаем простую конфигурацию в памяти
                config = new org.apache.commons.configuration2.PropertiesConfiguration();
            }
            config.setProperty(key, value);
        }
    }

    public String getString(String key) {
        return config.getString(key);
    }

    public String getString(String key, String defaultValue) {
        return config.getString(key, defaultValue);
    }

    public int getInt(String key) {
        return config.getInt(key);
    }

    public int getInt(String key, int defaultValue) {
        return config.getInt(key, defaultValue);
    }

    public long getLong(String key) {
        return config.getLong(key);
    }

    public long getLong(String key, long defaultValue) {
        return config.getLong(key, defaultValue);
    }

    public boolean getBoolean(String key) {
        return config.getBoolean(key);
    }

    public boolean getBoolean(String key, boolean defaultValue) {
        return config.getBoolean(key, defaultValue);
    }

    public boolean containsKey(String key) {
        return config.containsKey(key);
    }

    public void setProperty(String key, Object value) {
        config.setProperty(key, value);
    }

    public Configuration getConfiguration() {
        return config;
    }
}