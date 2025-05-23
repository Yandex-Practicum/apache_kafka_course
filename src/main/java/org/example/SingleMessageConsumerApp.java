package org.example;

import org.apache.kafka.clients.consumer.Consumer;
import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.clients.consumer.ConsumerRecords;
import org.apache.kafka.clients.consumer.KafkaConsumer;
import org.apache.kafka.clients.producer.ProducerConfig;
import org.apache.kafka.common.serialization.StringDeserializer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.List;
import java.util.Properties;

/**
 * This consumer reacts to messages immediately using a listener
 */
public class SingleMessageConsumerApp {
    private static final Logger logger = LoggerFactory.getLogger(SingleMessageConsumerApp.class);

    public static void main(String[] args) {
        Properties props = new Properties();
        props.put(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:29092,localhost:39092,localhost:49092");
        props.put(ConsumerConfig.GROUP_ID_CONFIG, "consumer-group-2");
        props.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        props.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        props.put(ConsumerConfig.FETCH_MIN_BYTES_CONFIG, 1); // минимальный размер сообщения
        props.put(ConsumerConfig.FETCH_MAX_WAIT_MS_CONFIG, 1); // минимальное время ожидания
        props.put(ConsumerConfig.ENABLE_AUTO_COMMIT_CONFIG, "true");

        Consumer<String, String> consumer = new KafkaConsumer<>(props);
        consumer.subscribe(List.of("test-topic"));

        while (true) {
            ConsumerRecords<String, String> records = consumer.poll(Duration.ofMillis(1000));
            logger.info("Received {} messages", records.count());
            records.forEach(record -> {
                try {
                    Message message = Message.deserialize(record.value());
                    logger.info("Received Message: " + message.getContent());
                } catch (Exception e) {
                    logger.error("Failure", e);
                }
            });
        }
    }
}
