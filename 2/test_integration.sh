#!/bin/bash

echo "🚀 Тестирование интеграции Kafka с Hadoop"
echo "=========================================="

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Функция для проверки статуса службы
check_service() {
    local service_name=$1
    local port=$2
    echo -e "${BLUE}Проверка $service_name на порту $port...${NC}"

    if nc -z localhost $port 2>/dev/null; then
        echo -e "${GREEN}✅ $service_name работает${NC}"
        return 0
    else
        echo -e "${RED}❌ $service_name недоступен${NC}"
        return 1
    fi
}

# Проверка всех сервисов
echo -e "${YELLOW}Шаг 1: Проверка статуса сервисов${NC}"
echo "================================"

check_service "Kafka" 9094
kafka_status=$?

check_service "Hadoop NameNode" 9000
hadoop_status=$?

check_service "Hadoop Web UI" 9870
hadoop_web_status=$?

if [ $kafka_status -ne 0 ] || [ $hadoop_status -ne 0 ]; then
    echo -e "${RED}⚠️  Некоторые сервисы недоступны. Убедитесь, что Docker Compose запущен:${NC}"
    echo "docker-compose up -d"
    exit 1
fi

echo ""
echo -e "${YELLOW}Шаг 2: Запуск продюсера для отправки тестовых данных${NC}"
echo "=================================================="

# Запуск продюсера в фоне и сохранение логов
echo "Отправка сообщений в Kafka..."
go run main.go producer > producer_logs.txt 2>&1 &
PRODUCER_PID=$!

# Ждем завершения продюсера
wait $PRODUCER_PID

echo -e "${GREEN}✅ Продюсер завершил работу${NC}"
echo ""

echo -e "${YELLOW}Шаг 3: Проверка топика Kafka${NC}"
echo "============================="

# Проверяем топик через console consumer
echo "Чтение сообщений из топика hadoop-topic:"
timeout 10s docker exec kafka-0 kafka-console-consumer.sh \
    --bootstrap-server localhost:9092 \
    --topic hadoop-topic \
    --from-beginning \
    --max-messages 10 > kafka_consumer_logs.txt 2>&1

if [ -s kafka_consumer_logs.txt ]; then
    echo -e "${GREEN}✅ Сообщения найдены в топике:${NC}"
    cat kafka_consumer_logs.txt
else
    echo -e "${RED}❌ Сообщения не найдены в топике${NC}"
fi

echo ""
echo -e "${YELLOW}Шаг 4: Запуск консьюмера для записи в HDFS${NC}"
echo "=========================================="

# Запуск консьюмера на 30 секунд
echo "Запуск консьюмера на 30 секунд..."
timeout 30s go run main.go consumer > consumer_logs.txt 2>&1 &
CONSUMER_PID=$!

# Ждем завершения консьюмера
wait $CONSUMER_PID

echo -e "${GREEN}✅ Консьюмер завершил работу${NC}"
echo ""

echo -e "${YELLOW}Шаг 5: Проверка данных в HDFS${NC}"
echo "=============================="

# Проверяем содержимое HDFS
echo "Проверка файлов в HDFS директории /data:"
docker exec hadoop-namenode hdfs dfs -ls /data > hdfs_ls_logs.txt 2>&1

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Файлы найдены в HDFS:${NC}"
    cat hdfs_ls_logs.txt

    # Показываем содержимое нескольких файлов
    echo ""
    echo "Содержимое файлов:"
    docker exec hadoop-namenode hdfs dfs -ls /data | grep "message_" | head -3 | while read line; do
        filename=$(echo $line | awk '{print $8}')
        echo -e "${BLUE}Файл: $filename${NC}"
        docker exec hadoop-namenode hdfs dfs -cat $filename
        echo "---"
    done
else
    echo -e "${RED}❌ Ошибка при чтении HDFS${NC}"
    cat hdfs_ls_logs.txt
fi

echo ""
echo -e "${YELLOW}Шаг 6: Сводка результатов${NC}"
echo "========================="

echo -e "${GREEN}✅ Логи продюсера:${NC} producer_logs.txt"
echo -e "${GREEN}✅ Логи Kafka консьюмера:${NC} kafka_consumer_logs.txt"
echo -e "${GREEN}✅ Логи HDFS консьюмера:${NC} consumer_logs.txt"
echo -e "${GREEN}✅ Список файлов HDFS:${NC} hdfs_ls_logs.txt"

echo ""
echo -e "${GREEN}🎉 Тестирование интеграции завершено!${NC}"
echo "Все логи сохранены в соответствующих файлах."
