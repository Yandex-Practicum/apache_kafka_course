package org.example;

import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.apache.kafka.clients.consumer.ConsumerRecords;
import org.apache.kafka.clients.consumer.KafkaConsumer;
import org.apache.kafka.clients.producer.ProducerConfig;
import org.apache.kafka.common.serialization.StringDeserializer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.List;
import java.util.Properties;

public class BatchMessageConsumerApp {
    private static final Logger logger = LoggerFactory.getLogger(BatchMessageConsumerApp.class);

    public static void main(String[] args) {
        Properties props = new Properties();
        props.put(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:29092,localhost:39092,localhost:49092");
        props.put(ConsumerConfig.GROUP_ID_CONFIG, "consumer-group-1");
        props.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        props.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        props.put(ConsumerConfig.ENABLE_AUTO_COMMIT_CONFIG, "false");
        props.put(ConsumerConfig.FETCH_MIN_BYTES_CONFIG, 10 * 120); // Примерный размер 10 сообщений
        props.put(ConsumerConfig.FETCH_MAX_WAIT_MS_CONFIG, 300000); // Максимальное время ожидания сообщений

        KafkaConsumer<String, String> consumer = new KafkaConsumer<>(props);
        consumer.subscribe(List.of("test-topic"));

        try {
            while (true) {
                ConsumerRecords<String, String> records = consumer.poll(Duration.ofMillis(1000));
                logger.info("Received {} messages", records.count());
                for (ConsumerRecord<String, String> record : records) {
                    Message message = Message.deserialize(record.value());
                    logger.info("Received Message: " + message.getContent());
                }
                consumer.commitSync();
            }
        } catch (Exception e) {
            logger.error("Failure", e);
            consumer.commitSync();
        } finally {
            consumer.close();
        }
    }
}
