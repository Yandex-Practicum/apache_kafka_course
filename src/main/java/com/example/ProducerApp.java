package com.example;

import org.apache.kafka.clients.producer.KafkaProducer;
import org.apache.kafka.clients.producer.ProducerRecord;
import org.apache.kafka.common.config.SslConfigs;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Properties;
import java.util.UUID;

public class ProducerApp {
    private static final Logger logger = LoggerFactory.getLogger(ProducerApp.class);
    public static void main(String[] args) {
        Properties props = new Properties();
        props.put("bootstrap.servers", "localhost:19092,localhost:29092,localhost:39092");
        props.put("group.id", "consumer-ssl-group");
        props.put("security.protocol", "SASL_SSL");
        props.put(SslConfigs.SSL_TRUSTSTORE_LOCATION_CONFIG, "./kafka-0-creds/kafka-0.truststore.jks");
        props.put(SslConfigs.SSL_TRUSTSTORE_PASSWORD_CONFIG, "your-password");
        props.put(SslConfigs.SSL_KEYSTORE_LOCATION_CONFIG, "./kafka-0-creds/kafka-0.keystore.jks");
        props.put(SslConfigs.SSL_KEYSTORE_PASSWORD_CONFIG, "your-password");
        props.put(SslConfigs.SSL_KEY_PASSWORD_CONFIG, "your-password");

        props.put("sasl.mechanism", "PLAIN");
        props.put("sasl.jaas.config",
                "org.apache.kafka.common.security.plain.PlainLoginModule required " +
                        "username=\"admin\" password=\"your-password\";");

        props.put("key.serializer", "org.apache.kafka.common.serialization.StringSerializer");
        props.put("value.serializer", "org.apache.kafka.common.serialization.StringSerializer");

        try (KafkaProducer<String, String> producer = new KafkaProducer<>(props)) {
            for (int i = 1; i <= 10; i++) {
                String key = "key-" + UUID.randomUUID();
                String value = "SSL message " + i;

                logger.info("Sending message {}: key={}, value={}", i, key, value);

                ProducerRecord<String, String> record = new ProducerRecord<>("secure-topic", key, value);
                producer.send(record, (metadata, exception) -> {
                    if (exception == null) {
                        logger.info("Message sent: partition={}, offset={}", metadata.partition(), metadata.offset());
                    } else {
                        logger.error("Error while producing: ", exception);
                    }
                });

                Thread.sleep(600);
            }
            producer.flush();
        } catch (Exception e) {
            logger.error("Error while producing: ", e);
        }
    }
}
