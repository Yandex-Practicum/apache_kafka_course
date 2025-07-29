package com.kafka.hdfs;

import org.junit.Test;
import static org.junit.Assert.*;

/**
 * Простые тесты для KafkaToHdfsConsumer
 */
public class KafkaToHdfsConsumerTest {

    @Test
    public void testMessageDataCreation() {
        String topic = "test-topic";
        int partition = 0;
        long offset = 123L;
        String key = "test-key";
        String value = "test-value";
        long processedAt = System.currentTimeMillis();

        MessageData messageData = new MessageData(topic, partition, offset, key, value, processedAt);

        assertEquals(topic, messageData.getTopic());
        assertEquals(partition, messageData.getPartition());
        assertEquals(offset, messageData.getOffset());
        assertEquals(key, messageData.getKey());
        assertEquals(value, messageData.getValue());
        assertEquals(processedAt, messageData.getProcessedAt());
        assertTrue(messageData.getTimestamp() > 0);
    }

    @Test
    public void testConfigDefaults() {
        Config config = new Config();

        assertNotNull(config.getString(Config.KAFKA_BOOTSTRAP_SERVERS));
        assertNotNull(config.getString(Config.KAFKA_TOPIC));
        assertNotNull(config.getString(Config.KAFKA_GROUP_ID));
        assertNotNull(config.getString(Config.HDFS_OUTPUT_PATH));
        assertNotNull(config.getString(Config.HDFS_NAMENODE_URL));

        assertTrue(config.getInt(Config.BATCH_SIZE) > 0);
        assertTrue(config.getLong(Config.FLUSH_INTERVAL_MS) > 0);
        assertTrue(config.getLong(Config.POLL_TIMEOUT_MS) > 0);
    }

    @Test
    public void testConfigWithCustomValues() {
        Config config = new Config();

        // Устанавливаем кастомные значения
        config.setProperty(Config.KAFKA_TOPIC, "custom-topic");
        config.setProperty(Config.BATCH_SIZE, 500);

        assertEquals("custom-topic", config.getString(Config.KAFKA_TOPIC));
        assertEquals(500, config.getInt(Config.BATCH_SIZE));
    }
}