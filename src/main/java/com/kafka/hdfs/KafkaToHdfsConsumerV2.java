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
 * Улучшенная версия основного класса для чтения данных из Kafka и записи в HDFS
 */
public class KafkaToHdfsConsumerV2 {
    private static final Logger logger = LoggerFactory.getLogger(KafkaToHdfsConsumerV2.class);

    private final KafkaConsumer<String, String> consumer;
    private final FileSystem hdfs;
    private final String hdfsOutputPath;
    private final ObjectMapper objectMapper;
    private final AtomicLong messageCounter;
    private final int batchSize;
    private final long pollTimeoutMs;
    private final Config config;

    public KafkaToHdfsConsumerV2(Config config) throws IOException {
        this.config = config;
        this.hdfsOutputPath = config.getString(Config.HDFS_OUTPUT_PATH);
        this.batchSize = config.getInt(Config.BATCH_SIZE, 100);
        this.pollTimeoutMs = config.getLong(Config.POLL_TIMEOUT_MS, 1000);
        this.messageCounter = new AtomicLong(0);
        this.objectMapper = new ObjectMapper();

        // Настройка Kafka Consumer
        Properties props = new Properties();
        props.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG, config.getString(Config.KAFKA_BOOTSTRAP_SERVERS));
        props.put(ConsumerConfig.GROUP_ID_CONFIG, config.getString(Config.KAFKA_GROUP_ID));
        props.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        props.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG,
                config.getString(Config.KAFKA_AUTO_OFFSET_RESET, "earliest"));
        props.put(ConsumerConfig.ENABLE_AUTO_COMMIT_CONFIG, config.getString(Config.KAFKA_ENABLE_AUTO_COMMIT, "false"));

        this.consumer = new KafkaConsumer<>(props);
        this.consumer.subscribe(Collections.singletonList(config.getString(Config.KAFKA_TOPIC)));

        // Настройка HDFS
        Configuration hdfsConf = new Configuration();
        hdfsConf.set("fs.defaultFS", config.getString(Config.HDFS_NAMENODE_URL));
        this.hdfs = FileSystem.get(hdfsConf);

        // Создаем директорию в HDFS если она не существует
        createHdfsDirectoryIfNotExists();

        logger.info("KafkaToHdfsConsumerV2 инициализирован");
        logger.info("Kafka: {} -> {}", config.getString(Config.KAFKA_BOOTSTRAP_SERVERS),
                config.getString(Config.KAFKA_TOPIC));
        logger.info("HDFS: {}", hdfsOutputPath);
        logger.info("Batch size: {}, Poll timeout: {}ms", batchSize, pollTimeoutMs);
    }

    private void createHdfsDirectoryIfNotExists() throws IOException {
        Path outputPath = new Path(hdfsOutputPath);
        if (!hdfs.exists(outputPath)) {
            hdfs.mkdirs(outputPath);
            logger.info("Создана директория в HDFS: {}", hdfsOutputPath);
        }
    }

    public void start() {
        logger.info("Начинаем чтение данных из Kafka и запись в HDFS");

        try {
            while (true) {
                ConsumerRecords<String, String> records = consumer.poll(Duration.ofMillis(pollTimeoutMs));

                if (!records.isEmpty()) {
                    logger.debug("Получено {} записей из Kafka", records.count());
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
        try {
            Config config;

            if (args.length > 0) {
                // Загружаем конфигурацию из файла
                config = new Config(args[0]);
                logger.info("Конфигурация загружена из файла: {}", args[0]);
            } else {
                // Используем конфигурацию по умолчанию
                config = new Config();
                logger.info("Используется конфигурация по умолчанию");
            }

            KafkaToHdfsConsumerV2 consumer = new KafkaToHdfsConsumerV2(config);

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