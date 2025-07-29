package com.kafka.hdfs;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * Класс для хранения данных сообщений из Kafka
 */
public class MessageData {
    @JsonProperty("topic")
    private String topic;

    @JsonProperty("partition")
    private int partition;

    @JsonProperty("offset")
    private long offset;

    @JsonProperty("key")
    private String key;

    @JsonProperty("value")
    private String value;

    @JsonProperty("timestamp")
    private long timestamp;

    @JsonProperty("processed_at")
    private long processedAt;

    public MessageData() {
        // Конструктор по умолчанию для Jackson
    }

    public MessageData(String topic, int partition, long offset, String key, String value, long processedAt) {
        this.topic = topic;
        this.partition = partition;
        this.offset = offset;
        this.key = key;
        this.value = value;
        this.timestamp = System.currentTimeMillis();
        this.processedAt = processedAt;
    }

    // Геттеры
    public String getTopic() {
        return topic;
    }

    public int getPartition() {
        return partition;
    }

    public long getOffset() {
        return offset;
    }

    public String getKey() {
        return key;
    }

    public String getValue() {
        return value;
    }

    public long getTimestamp() {
        return timestamp;
    }

    public long getProcessedAt() {
        return processedAt;
    }

    // Сеттеры
    public void setTopic(String topic) {
        this.topic = topic;
    }

    public void setPartition(int partition) {
        this.partition = partition;
    }

    public void setOffset(long offset) {
        this.offset = offset;
    }

    public void setKey(String key) {
        this.key = key;
    }

    public void setValue(String value) {
        this.value = value;
    }

    public void setTimestamp(long timestamp) {
        this.timestamp = timestamp;
    }

    public void setProcessedAt(long processedAt) {
        this.processedAt = processedAt;
    }

    @Override
    public String toString() {
        return "MessageData{" +
                "topic='" + topic + '\'' +
                ", partition=" + partition +
                ", offset=" + offset +
                ", key='" + key + '\'' +
                ", value='" + value + '\'' +
                ", timestamp=" + timestamp +
                ", processedAt=" + processedAt +
                '}';
    }
}