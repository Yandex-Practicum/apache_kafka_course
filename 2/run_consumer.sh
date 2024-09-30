#!/bin/bash

echo "🚀 Запуск Kafka-HDFS консьюмера"
echo "==============================="

# Проверка Kafka
if ! nc -z localhost 9094 2>/dev/null; then
    echo "❌ Kafka недоступен на порту 9094"
    echo "Запустите Docker Compose: docker-compose up -d"
    exit 1
fi

# Проверка Hadoop
if ! nc -z localhost 9000 2>/dev/null; then
    echo "❌ Hadoop NameNode недоступен на порту 9000"
    echo "Запустите Docker Compose: docker-compose up -d"
    exit 1
fi

echo "✅ Kafka доступен"
echo "✅ Hadoop доступен"
echo "📥 Чтение сообщений из топика hadoop-topic и запись в HDFS..."
echo ""

go run main.go consumer
