package org.example;

import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.clients.consumer.KafkaConsumer;
import org.apache.kafka.common.serialization.Serdes;
import org.apache.kafka.streams.KafkaStreams;
import org.apache.kafka.streams.KeyValue;
import org.apache.kafka.streams.StreamsBuilder;
import org.apache.kafka.streams.StreamsConfig;
import org.apache.kafka.streams.kstream.KStream;
import org.apache.kafka.streams.kstream.KTable;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashSet;
import java.util.List;
import java.util.Properties;
import java.util.Set;

public class MessageProcessor {

    private static List<String> bannedWords = new ArrayList<>();
    private static final Logger logger = LoggerFactory.getLogger(MessageProcessor.class);

    public static void main(String[] args) {
        loadBannedWords("src/main/resources/banned_words.txt");

        Properties props = new Properties();
        props.put(StreamsConfig.APPLICATION_ID_CONFIG, "message-processor");
        props.put(StreamsConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:29092");
        props.put(StreamsConfig.DEFAULT_KEY_SERDE_CLASS_CONFIG, Serdes.String().getClass());
        props.put(StreamsConfig.DEFAULT_VALUE_SERDE_CLASS_CONFIG, Serdes.String().getClass());
        props.put(ConsumerConfig.GROUP_INSTANCE_ID_CONFIG, "message-processor-instance");
        props.put(StreamsConfig.REQUEST_TIMEOUT_MS_CONFIG, 60000);
        props.put(ConsumerConfig.SESSION_TIMEOUT_MS_CONFIG, 30000);
        props.put(ConsumerConfig.HEARTBEAT_INTERVAL_MS_CONFIG, 10000);
        props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");


        StreamsBuilder builder = new StreamsBuilder();

        // Чтение списка заблокированных пользователей
        KTable<String, String> blockedUsers = builder.table("blocked_users");

        // Чтение входящих сообщений
        KStream<String, String> messages = builder.stream("messages");
        messages.foreach((key, value) -> logger.info("Direction: " + key + ", message: " + value));
        // Установить получатель как key, перенести сообщение в новую структуру
        messages.map((direction, message) ->
                        new KeyValue<>(direction.split("->")[1],
                                        "%s:%s".formatted(direction.split("->")[0], message)))
                .peek((key, value) -> logger.info("key=%s,value=%s".formatted(key, value)))
        // Коннект к таблице заблокированных пользователей и фильтрация
                .leftJoin(
                        blockedUsers,
                        (message, blockedList) -> {
                            String sender = message.split(":")[0];
                            List<String> blocked = Arrays.stream(blockedList.split(","))
                                                            .map(String::trim)
                                                            .toList();
                            logger.info("Check that %s is not in %s".formatted(sender, blocked));
                            if (blocked.contains(sender)) {
                                logger.info("This message will be ignored");
                                return null; // Если получатель заблокировал отправителя, исключаем сообщение
                            }
                            return message;
                        }
                )
                .filter((key, value) -> value != null) // Исключаем null сообщения
        // Возвращаем исходную структуру потока
                .map((receiver, message) ->
                    new KeyValue<>("%s->%s".formatted(message.split(":")[0], receiver),
                                    message.split(":")[1]))
                .peek((key, value) -> logger.info("key=%s,value=%s".formatted(key, value)))
        // Замена запрещенных слов на звездочки
                .mapValues(value -> {
                    for (String bannedWord : bannedWords) {
                        value = value.replaceAll(bannedWord, "***");
                    }
                    return value;
                })
        // Отправка отфильтрованных сообщений в `filtered_messages`
                .to("filtered_messages");

        KafkaStreams streams = new KafkaStreams(builder.build(), props);
        streams.start();

        // Завершение потока
//        Runtime.getRuntime().addShutdownHook(new Thread(streams::close));
    }
    private static void loadBannedWords(String filePath) {
        try {
            bannedWords = Files.readAllLines(Paths.get(filePath));
            logger.info("Loaded banned words: " + bannedWords);
        } catch (IOException e) {
            logger.error("Error loading banned words from file: " + filePath, e);
        }
    }
}
