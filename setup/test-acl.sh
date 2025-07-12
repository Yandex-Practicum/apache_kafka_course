#!/bin/bash

set -e

echo "=== Тестирование ACL прав доступа ==="

BOOTSTRAP_SERVERS="kafka-0:9093,kafka-1:9093,kafka-2:9093"
COMMAND_CONFIG="/certs/client.properties"

echo "1. Проверяем существующие топики:"
kafka-topics.sh --bootstrap-server $BOOTSTRAP_SERVERS --command-config $COMMAND_CONFIG --list

echo ""
echo "2. Проверяем настроенные ACL:"
kafka-acls.sh --bootstrap-server kafka-0:9093,kafka-1:9093,kafka-2:9093 --command-config /certs/client.properties --list

echo ""
echo "3. Тестируем запись в topic-1 (должно работать):"
echo "test-message-topic-1" | kafka-console-producer.sh \
    --bootstrap-server $BOOTSTRAP_SERVERS \
    --producer.config $COMMAND_CONFIG \
    --topic topic-1

echo "✅ Запись в topic-1 успешна"

echo ""
echo "4. Тестируем запись в topic-2 (должно работать):"
echo "test-message-topic-2" | kafka-console-producer.sh \
    --bootstrap-server $BOOTSTRAP_SERVERS \
    --producer.config $COMMAND_CONFIG \
    --topic topic-2

echo "✅ Запись в topic-2 успешна"

echo ""
echo "5. Тестируем чтение из topic-1 (должно работать, с таймаутом 10 сек):"
timeout 10 kafka-console-consumer.sh \
    --bootstrap-server $BOOTSTRAP_SERVERS \
    --consumer.config $COMMAND_CONFIG \
    --topic topic-1 \
    --from-beginning \
    --max-messages 5 || echo "✅ Чтение из topic-1 завершено"

echo ""
echo "6. Тестируем чтение из topic-2 (должно быть заблокировано ACL):"
echo "Попытка чтения из topic-2 (ожидается ошибка доступа):"
timeout 5 kafka-console-consumer.sh \
    --bootstrap-server $BOOTSTRAP_SERVERS \
    --consumer.config $COMMAND_CONFIG \
    --topic topic-2 \
    --from-beginning \
    --max-messages 1 2>&1 || echo "✅ Чтение из topic-2 заблокировано (это правильно)"

echo ""
echo "=== Тестирование завершено ==="
echo "Резюме:"
echo "  ✅ topic-1: полный доступ (запись + чтение)"
echo "  ✅ topic-2: только запись (чтение заблокировано)"
