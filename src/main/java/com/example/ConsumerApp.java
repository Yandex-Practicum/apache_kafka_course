package com.example;

import org.apache.kafka.clients.consumer.*;
import org.apache.kafka.common.config.SslConfigs;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.Collections;
import java.util.Properties;

public class ConsumerApp {
    private static final Logger logger = LoggerFactory.getLogger(ConsumerApp.class);

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

        props.put("key.deserializer", "org.apache.kafka.common.serialization.StringDeserializer");
        props.put("value.deserializer", "org.apache.kafka.common.serialization.StringDeserializer");

        Consumer<String, String> consumer = new KafkaConsumer<>(props);
        consumer.subscribe(Collections.singletonList("secure-topic"));

        while (true) {
            ConsumerRecords<String, String> records = consumer.poll(Duration.ofSeconds(1));
            for (ConsumerRecord<String, String> record : records) {
                logger.info("Received record: key={} value={}", record.key(), record.value());
            }
        }
    }
}
