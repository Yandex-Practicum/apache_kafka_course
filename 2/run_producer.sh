#!/bin/bash

echo "🚀 Запуск Kafka продюсера"
echo "========================"

# Проверка Kafka
if ! nc -z localhost 9094 2>/dev/null; then
    echo "❌ Kafka недоступен на порту 9094"
    echo "Запустите Docker Compose: docker-compose up -d"
    exit 1
fi

echo "✅ Kafka доступен"
echo "📤 Отправка тестовых сообщений в топик hadoop-topic..."
echo ""

go run main.go producer
