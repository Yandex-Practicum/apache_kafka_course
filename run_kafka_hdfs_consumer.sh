#!/bin/bash

# Скрипт для запуска Kafka to HDFS Consumer

# Проверяем наличие Maven
if ! command -v mvn &> /dev/null; then
    echo "Ошибка: Maven не установлен"
    exit 1
fi

# Проверяем наличие Java
if ! command -v java &> /dev/null; then
    echo "Ошибка: Java не установлена"
    exit 1
fi

echo "Сборка проекта..."
mvn clean package -DskipTests

if [ $? -ne 0 ]; then
    echo "Ошибка при сборке проекта"
    exit 1
fi

echo "Запуск Kafka to HDFS Consumer..."

# Параметры по умолчанию
BOOTSTRAP_SERVERS=${1:-"localhost:9092"}
TOPIC=${2:-"test-topic"}
GROUP_ID=${3:-"kafka-hdfs-consumer"}
HDFS_OUTPUT_PATH=${4:-"/kafka-data"}
BATCH_SIZE=${5:-100}
FLUSH_INTERVAL_MS=${6:-5000}

echo "Параметры запуска:"
echo "  Bootstrap Servers: $BOOTSTRAP_SERVERS"
echo "  Topic: $TOPIC"
echo "  Group ID: $GROUP_ID"
echo "  HDFS Output Path: $HDFS_OUTPUT_PATH"
echo "  Batch Size: $BATCH_SIZE"
echo "  Flush Interval: ${FLUSH_INTERVAL_MS}ms"

# Создаем директорию для логов если её нет
mkdir -p logs

# Запускаем приложение
java -cp target/kafka-to-hdfs-1.0.0.jar \
     -Dlogback.configurationFile=src/main/resources/logback.xml \
     com.kafka.hdfs.KafkaToHdfsConsumer \
     "$BOOTSTRAP_SERVERS" \
     "$TOPIC" \
     "$GROUP_ID" \
     "$HDFS_OUTPUT_PATH" \
     "$BATCH_SIZE" \
     "$FLUSH_INTERVAL_MS" 