package org.example;

import org.apache.kafka.clients.producer.KafkaProducer;
import org.apache.kafka.clients.producer.ProducerConfig;
import org.apache.kafka.clients.producer.ProducerRecord;
import org.apache.kafka.common.serialization.StringSerializer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Properties;
import java.util.Random;

/**
 * This producer sends messages to Kafka with "at-least-once" delivery guarantee
 */
public class ProducerApp {
    private static final Logger logger = LoggerFactory.getLogger(ProducerApp.class);

    public static void main(String[] args) {
        Properties props = new Properties();
        props.put(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:29092,localhost:39092,localhost:49092");
        props.put(ProducerConfig.KEY_SERIALIZER_CLASS_CONFIG, StringSerializer.class.getName());
        props.put(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG, StringSerializer.class.getName());

        // Настройка подтверждения доставки
        props.put(ProducerConfig.ACKS_CONFIG, "all"); // гарантия «at-least-once»
        props.put(ProducerConfig.RETRIES_CONFIG, 3);  // поведение при сбое

        for (int i = 1; i < 100; i++) {
            logger.info(String.format("Preparing to send message number %s", i));
            try (KafkaProducer<String, String> producer = new KafkaProducer<>(props)) {

                Message message = new Message("user123" + new Random().nextInt(1, 100), "Hello Kafka!");
                ProducerRecord<String, String> record = new ProducerRecord<>("test-topic", message.getUserId(),
                        Message.serialize(message));

                producer.send(record, (metadata, exception) -> {
                    if (exception == null) {
                        logger.info(String.format("Message sent successfully: %s", metadata.toString()));
                    } else {
                        logger.error("Error while producing: ", exception);
                    }
                });
                Thread.sleep(600);
            } catch (Exception e) {
                logger.error("Error while producing: ", e);
            }
        }
    }
}
