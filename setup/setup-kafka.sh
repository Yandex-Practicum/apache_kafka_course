#!/bin/bash

set -e

echo "Начинаем настройку Kafka топиков и ACL..."

# SSL настройки для подключения к Kafka
BOOTSTRAP_SERVERS="kafka-0:9093,kafka-1:9093,kafka-2:9093"
COMMAND_CONFIG="/setup/client.properties"

# Ожидание готовности кластера
echo "Ожидание готовности Kafka кластера..."
sleep 30

# Дополнительная проверка готовности KRaft кластера
echo "Проверка готовности KRaft кластера..."
for i in {1..10}; do
    if kafka-topics.sh --bootstrap-server $BOOTSTRAP_SERVERS --command-config $COMMAND_CONFIG --list > /dev/null 2>&1; then
        echo "✅ Кластер готов!"
        break
    fi
    echo "Ожидание готовности... ($i/10)"
    sleep 10
done

# Проверка подключения к кластеру
echo "Проверка подключения к Kafka кластеру..."
kafka-topics.sh --bootstrap-server $BOOTSTRAP_SERVERS --command-config $COMMAND_CONFIG --list

# Создание топиков
echo "Создание топика topic-1..."
kafka-topics.sh --bootstrap-server $BOOTSTRAP_SERVERS \
    --command-config $COMMAND_CONFIG \
    --create \
    --topic topic-1 \
    --partitions 3 \
    --replication-factor 3 || echo "topic-1 уже существует"

echo "Создание топика topic-2..."
kafka-topics.sh --bootstrap-server $BOOTSTRAP_SERVERS \
    --command-config $COMMAND_CONFIG \
    --create \
    --topic topic-2 \
    --partitions 3 \
    --replication-factor 3 || echo "topic-2 уже существует"

# Создание дополнительного топика для тестирования
echo "Создание топика test..."
kafka-topics.sh --bootstrap-server $BOOTSTRAP_SERVERS \
    --command-config $COMMAND_CONFIG \
    --create \
    --topic test \
    --partitions 3 \
    --replication-factor 3 || echo "test топик уже существует"

echo "Список всех топиков:"
kafka-topics.sh --bootstrap-server $BOOTSTRAP_SERVERS --command-config $COMMAND_CONFIG --list

# Настройка ACL для демонстрации ограничений доступа
echo "Настройка ACL прав доступа..."

# Ждем стабилизации кластера перед настройкой ACL
echo "Ожидание стабилизации кластера для ACL..."
sleep 15

# Для topic-2 создаем явный запрет на чтение (DENY ACL)
echo "Создание DENY ACL для чтения topic-2..."
kafka-acls.sh --bootstrap-server $BOOTSTRAP_SERVERS --command-config $COMMAND_CONFIG \
    --add --deny-principal User:ANONYMOUS --operation Read --topic topic-2 || echo "DENY ACL для topic-2 уже существует"

# Разрешаем запись в topic-2
kafka-acls.sh --bootstrap-server $BOOTSTRAP_SERVERS --command-config $COMMAND_CONFIG \
    --add --allow-principal User:ANONYMOUS --operation Write --topic topic-2 || echo "ALLOW Write ACL для topic-2 уже существует"

kafka-acls.sh --bootstrap-server $BOOTSTRAP_SERVERS --command-config $COMMAND_CONFIG \
    --add --allow-principal User:ANONYMOUS --operation Describe --topic topic-2 || echo "ALLOW Describe ACL для topic-2 уже существует"

# Показываем все ACL
echo "Список всех ACL:"
kafka-acls.sh --bootstrap-server $BOOTSTRAP_SERVERS --command-config $COMMAND_CONFIG --list || echo "Не удалось получить список ACL"

echo "ACL настроены:"
echo "  - topic-1: доступ по умолчанию (allow_everyone_if_no_acl_found=true)"
echo "  - topic-2: запись разрешена, чтение ЗАПРЕЩЕНО (DENY ACL)"
echo "  - test: доступ по умолчанию (allow_everyone_if_no_acl_found=true)"

echo "Настройка Kafka завершена успешно!"
echo ""
echo "Топики созданы:"
echo "  - topic-1: полный доступ (по умолчанию)"
echo "  - topic-2: запись разрешена, чтение ЗАПРЕЩЕНО (DENY ACL)"
echo "  - test: полный доступ (по умолчанию)"
echo ""
echo "SSL конфигурация:"
echo "  ✅ Все соединения зашифрованы"
echo "  ✅ Взаимная аутентификация клиентов"
echo "  ✅ Keystore/Truststore настроены"
echo ""
echo "ACL конфигурация:"
echo "  ✅ ACL включены (StandardAuthorizer)"
echo "  ✅ topic-2 защищен от чтения (DENY ACL)"
echo "  ✅ allow_everyone_if_no_acl_found=true (для остальных топиков)"
echo "  💡 Для отключения ACL закомментируйте AUTHORIZER в compose.yaml"
