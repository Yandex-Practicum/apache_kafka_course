package com.kafka.hdfs;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.FSDataOutputStream;
import org.apache.hadoop.fs.FileSystem;
import org.apache.hadoop.fs.Path;
import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.apache.kafka.clients.consumer.ConsumerRecords;
import org.apache.kafka.clients.consumer.KafkaConsumer;
import org.apache.kafka.common.serialization.StringDeserializer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.time.Duration;
import java.util.Collections;
import java.util.Properties;
import java.util.concurrent.atomic.AtomicLong;

/**
 * Основной класс для чтения данных из Kafka и записи в HDFS
 */
public class KafkaToHdfsConsumer {
    private static final Logger logger = LoggerFactory.getLogger(KafkaToHdfsConsumer.class);

    private final KafkaConsumer<String, String> consumer;
    private final FileSystem hdfs;
    private final String hdfsOutputPath;
    private final ObjectMapper objectMapper;
    private final AtomicLong messageCounter;
    private final int batchSize;
    private final long flushIntervalMs;

    public KafkaToHdfsConsumer(String bootstrapServers, String topic, String groupId,
            String hdfsOutputPath, int batchSize, long flushIntervalMs) throws IOException {
        this.hdfsOutputPath = hdfsOutputPath;
        this.batchSize = batchSize;
        this.flushIntervalMs = flushIntervalMs;
        this.messageCounter = new AtomicLong(0);
        this.objectMapper = new ObjectMapper();

        // Настройка Kafka Consumer
        Properties props = new Properties();
        props.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG, bootstrapServers);
        props.put(ConsumerConfig.GROUP_ID_CONFIG, groupId);
        props.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        props.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");
        props.put(ConsumerConfig.ENABLE_AUTO_COMMIT_CONFIG, "false");

        this.consumer = new KafkaConsumer<>(props);
        this.consumer.subscribe(Collections.singletonList(topic));

        // Настройка HDFS
        Configuration conf = new Configuration();
        conf.set("fs.defaultFS", "hdfs://namenode:9000");
        this.hdfs = FileSystem.get(conf);

        logger.info("KafkaToHdfsConsumer инициализирован");
        logger.info("Kafka: {} -> {}", bootstrapServers, topic);
        logger.info("HDFS: {}", hdfsOutputPath);
    }

    public void start() {
        logger.info("Начинаем чтение данных из Kafka и запись в HDFS");

        try {
            while (true) {
                ConsumerRecords<String, String> records = consumer.poll(Duration.ofMillis(1000));

                if (!records.isEmpty()) {
                    processBatch(records);
                }
            }
        } catch (Exception e) {
            logger.error("Ошибка при обработке сообщений", e);
        } finally {
            close();
        }
    }

    private void processBatch(ConsumerRecords<String, String> records) {
        try {
            StringBuilder batchData = new StringBuilder();
            int count = 0;

            for (ConsumerRecord<String, String> record : records) {
                try {
                    // Создаем объект для записи
                    MessageData messageData = new MessageData(
                            record.topic(),
                            record.partition(),
                            record.offset(),
                            record.key(),
                            record.value(),
                            System.currentTimeMillis());

                    // Сериализуем в JSON
                    String jsonLine = objectMapper.writeValueAsString(messageData) + "\n";
                    batchData.append(jsonLine);
                    count++;

                    if (count >= batchSize) {
                        writeToHdfs(batchData.toString(), count);
                        batchData.setLength(0);
                        count = 0;
                    }

                } catch (Exception e) {
                    logger.error("Ошибка при обработке записи: {}", record, e);
                }
            }

            // Записываем оставшиеся данные
            if (count > 0) {
                writeToHdfs(batchData.toString(), count);
            }

            // Подтверждаем обработку
            consumer.commitSync();

        } catch (Exception e) {
            logger.error("Ошибка при обработке батча", e);
        }
    }

    private void writeToHdfs(String data, int count) throws IOException {
        long timestamp = System.currentTimeMillis();
        String fileName = String.format("kafka_data_%d_%d.json", timestamp, messageCounter.incrementAndGet());
        Path filePath = new Path(hdfsOutputPath + "/" + fileName);

        try (FSDataOutputStream outputStream = hdfs.create(filePath, true)) {
            outputStream.writeBytes(data);
            outputStream.flush();
            logger.info("Записано {} сообщений в файл: {}", count, filePath);
        } catch (IOException e) {
            logger.error("Ошибка при записи в HDFS: {}", filePath, e);
            throw e;
        }
    }

    public void close() {
        try {
            if (consumer != null) {
                consumer.close();
                logger.info("Kafka consumer закрыт");
            }
            if (hdfs != null) {
                hdfs.close();
                logger.info("HDFS connection закрыт");
            }
        } catch (Exception e) {
            logger.error("Ошибка при закрытии ресурсов", e);
        }
    }

    public static void main(String[] args) {
        if (args.length < 4) {
            System.out.println(
                    "Использование: java KafkaToHdfsConsumer <bootstrap-servers> <topic> <group-id> <hdfs-output-path> [batch-size] [flush-interval-ms]");
            System.exit(1);
        }

        String bootstrapServers = args[0];
        String topic = args[1];
        String groupId = args[2];
        String hdfsOutputPath = args[3];
        int batchSize = args.length > 4 ? Integer.parseInt(args[4]) : 100;
        long flushIntervalMs = args.length > 5 ? Long.parseLong(args[5]) : 5000;

        try {
            KafkaToHdfsConsumer consumer = new KafkaToHdfsConsumer(
                    bootstrapServers, topic, groupId, hdfsOutputPath, batchSize, flushIntervalMs);

            // Добавляем обработчик сигналов для graceful shutdown
            Runtime.getRuntime().addShutdownHook(new Thread(() -> {
                logger.info("Получен сигнал завершения, закрываем consumer...");
                consumer.close();
            }));

            consumer.start();

        } catch (Exception e) {
            logger.error("Ошибка при запуске consumer", e);
            System.exit(1);
        }
    }
}