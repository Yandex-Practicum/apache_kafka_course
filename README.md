# Анализ и агрегация сообщений с использованием ksqlDB

Пример проекта, как анализировать сообщения в реальном времени с помощью ksqlDB, Kafka и Docker Compose

### Инфраструктура

- Kafka: Обеспечивает обработку сообщений
- ksqlDB: Реализует анализ в реальном времени
- Docker Compose: Автоматизирует развертывание инфраструктуры

### Запуск

   ```bash
   docker-compose up -d
```

### Тестирование

1. Создать топик

```bash
docker exec -it ksql-kafka-1 bash

kafka-topics --create --bootstrap-server kafka:9092 --replication-factor 1 --partitions 1 --topic messages
```

2. Отправить тестовые данные в  "messages"

```bash
kafka-console-producer --broker-list kafka:9092 --topic messages
```

Пример данных, вставляются построчно:

```json
{"user_id": "user1", "recipient_id": "user2", "message": "Hello!", "timestamp": 1691486400000}
{"user_id": "user1","recipient_id": "user3","message": "Hi!","timestamp": 1691486401000}
{"user_id": "user2","recipient_id": "user1","message": "Reply","timestamp": 1691486402000}
{"user_id": "user3", "recipient_id": "user2","message": "Hey there!","timestamp": 1691486403000}
```

3. Отправить запросы из файла ksqldb-queries.sql в ksqlDB CLI

 ```bash
docker exec -it ksql-ksqldb-cli-1 ksql http://ksqldb-server:8088
RUN SCRIPT '/ksqldb-queries.sql';
```

4. Проверить результат:

```sql
Общее количество отправленных сообщений:
SELECT * FROM total_messages_sent EMIT CHANGES;

Уникальные получатели:
SELECT * FROM unique_recipients EMIT CHANGES;

Топ-5 активных пользователей:
SELECT user_id, sent_messages FROM top_active_users DESC LIMIT 5;
```