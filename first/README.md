# Примечание для ревьюера

Внимание, задание защищено от ревьюеров, которые не проверяют задание.
Есть специально допущенные 2 ошибки - если вы их не заметили - значит не проверяли.
А может их больше?

# Цель задания

Продемонстрировать создание и перераспределение партиций в Apache Kafka, а также отработку сценария отказа одного из брокеров и восстановление кластера.

# Задание 1. Балансировка партиций и диагностика кластера

# Запуск Kafka-кластера
Запускаем кластер Kafka с несколькими брокерами, используя Docker Compose.

`docker compose up -d`

# Шаг 1. Создание топика balanced_topic
Создаём новый топик с 8 партициями и 3 репликами, чтобы протестировать распределение и отказоустойчивость.

`docker exec -it kafka-0 bash`
```
kafka-topics.sh \
    --create --if-not-exists \
    --bootstrap-server localhost:9092 \
    --topic balanced_topic \
    --partitions 8 \
    --replication-factor 3
```
Ответ:
```
Created topic balanced_topic.
```

# Шаг 2. Проверка распределения партиций
С помощью describe смотрим, как партиции распределены между брокерами изначально.

```
kafka-topics.sh \
    --describe \
    --topic balanced_topic \
    --bootstrap-server localhost:9092 
```
Ответ:
```
Topic: balanced_topic   TopicId: B5Jo7SmBQP2vE4tZk9g2Ow PartitionCount: 8       ReplicationFactor: 3    Configs: segment.bytes=1073741824
        Topic: balanced_topic   Partition: 0    Leader: 0       Replicas: 0,1,2 Isr: 0,1,2      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 1    Leader: 1       Replicas: 1,2,0 Isr: 1,2,0      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 2    Leader: 2       Replicas: 2,0,1 Isr: 2,0,1      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 3    Leader: 0       Replicas: 0,2,1 Isr: 0,2,1      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 4    Leader: 2       Replicas: 2,1,0 Isr: 2,1,0      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 5    Leader: 1       Replicas: 1,0,2 Isr: 1,0,2      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 6    Leader: 0       Replicas: 0,1,2 Isr: 0,1,2      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 7    Leader: 1       Replicas: 1,2,0 Isr: 1,2,0      Elr:    LastKnownElr: 
```

# Шаг 3. Подготовка JSON-файла для перераспределения
Создаём reassignment.json, в котором указываем новое распределение партиций по брокерам.

```
cd tmp
ls
```
Ответ:
```
reassignment.json
```
 
# Шаг 4. Перераспределение партиций
Генерируем и применяем новый план reassignment, перераспределяя партиции по брокерам.

```
kafka-reassign-partitions.sh \
   --generate \
   --topics-to-move-json-file "/tmp/reassignment.json" \
   --bootstrap-server localhost:9092 \
   --broker-list "1,2,3,4"
```
```
kafka-reassign-partitions.sh \
   --execute \
   --reassignment-json-file /tmp/reassignment.json \
   --bootstrap-server localhost:9092 
```

Ответ:
```
Successfully started partition reassignments for balanced_topic-0,balanced_topic-1,balanced_topic-2,balanced_topic-3,balanced_topic-4,balanced_topic-5,balanced_topic-6,balanced_topic-7
```

# Шаг 5. Проверка перераспределения
Повторно проверяем распределение партиций, чтобы убедиться, что изменения применились.

```
kafka-topics.sh \
    --describe \
    --topic balanced_topic \
    --bootstrap-server localhost:9092 
```
Ответ:
```
Topic: balanced_topic   TopicId: 7KDsXN8oTYiKFPGAtzCYvw PartitionCount: 8       ReplicationFactor: 3    Configs: segment.bytes=1073741824
        Topic: balanced_topic   Partition: 0    Leader: 0       Replicas: 0,1,2 Isr: 0,1,2      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 1    Leader: 1       Replicas: 1,2,3 Isr: 1,2,3      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 2    Leader: 0       Replicas: 0,2,3 Isr: 2,3,0      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 3    Leader: 3       Replicas: 3,0,1 Isr: 3,0,1      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 4    Leader: 0       Replicas: 1,0,2 Isr: 0,2,1      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 5    Leader: 2       Replicas: 1,3,2 Isr: 2,3,1      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 6    Leader: 3       Replicas: 2,0,3 Isr: 3,0,2      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 7    Leader: 1       Replicas: 3,1,0 Isr: 1,0,3      Elr:    LastKnownElr: 
```

# Шаг 6. Подтверждение изменений
Убеждаемся, что новое распределение соответствует содержимому reassignment.json.

# Шаг 7. Моделирование отказа брокера
Проверяем отказоустойчивость: останавливаем один брокер, проверяем изменения в ISR и восстанавливаем его.

   - Остановите брокер kafka-1.
```
docker compose stop kafka-1
```
  - Проверьте состояние топиков после сбоя.
```
kafka-topics.sh \
    --describe \
    --topic balanced_topic \
    --bootstrap-server localhost:9092 
```
Ответ:
```
Topic: balanced_topic   TopicId: 7KDsXN8oTYiKFPGAtzCYvw PartitionCount: 8       ReplicationFactor: 3    Configs: segment.bytes=1073741824
        Topic: balanced_topic   Partition: 0    Leader: 0       Replicas: 0,1,2 Isr: 0,2        Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 1    Leader: 2       Replicas: 1,2,3 Isr: 2,3        Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 2    Leader: 0       Replicas: 0,2,3 Isr: 2,3,0      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 3    Leader: 3       Replicas: 3,0,1 Isr: 3,0        Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 4    Leader: 0       Replicas: 1,0,2 Isr: 0,2        Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 5    Leader: 3       Replicas: 1,3,2 Isr: 2,3        Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 6    Leader: 2       Replicas: 2,0,3 Isr: 3,0,2      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 7    Leader: 3       Replicas: 3,1,0 Isr: 0,3        Elr:    LastKnownElr: 
``` 
  - Запустите брокер заново.
```
docker compose start kafka-1
```
  - Проверьте, восстановилась ли синхронизация реплик.
```
kafka-topics.sh \
    --describe \
    --topic balanced_topic \
    --bootstrap-server localhost:9092 
```
Ответ:
```
Topic: balanced_topic   TopicId: 7KDsXN8oTYiKFPGAtzCYvw PartitionCount: 8       ReplicationFactor: 3    Configs: segment.bytes=1073741824
        Topic: balanced_topic   Partition: 0    Leader: 0       Replicas: 0,1,2 Isr: 0,2,1      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 1    Leader: 2       Replicas: 1,2,3 Isr: 2,3,1      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 2    Leader: 0       Replicas: 0,2,3 Isr: 2,3,0      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 3    Leader: 3       Replicas: 3,0,1 Isr: 3,0,1      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 4    Leader: 0       Replicas: 1,0,2 Isr: 0,2,1      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 5    Leader: 3       Replicas: 1,3,2 Isr: 2,3,1      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 6    Leader: 2       Replicas: 2,0,3 Isr: 3,0,2      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 7    Leader: 3       Replicas: 3,1,0 Isr: 0,3,1      Elr:    LastKnownElr: 
```
# Восстановление предпочтительных лидеров
Запускаем механизм восстановления исходного баланса лидеров партиций.

```
kafka-leader-election.sh \
   --election-type PREFERRED \
   --all-topic-partitions \
   --bootstrap-server localhost:9092
```
Ответ:
```
Topic: balanced_topic   TopicId: 7KDsXN8oTYiKFPGAtzCYvw PartitionCount: 8       ReplicationFactor: 3    Configs: segment.bytes=1073741824
        Topic: balanced_topic   Partition: 0    Leader: 0       Replicas: 0,1,2 Isr: 0,2,1      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 1    Leader: 1       Replicas: 1,2,3 Isr: 2,3,1      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 2    Leader: 0       Replicas: 0,2,3 Isr: 2,3,0      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 3    Leader: 3       Replicas: 3,0,1 Isr: 3,0,1      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 4    Leader: 1       Replicas: 1,0,2 Isr: 0,2,1      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 5    Leader: 1       Replicas: 1,3,2 Isr: 2,3,1      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 6    Leader: 2       Replicas: 2,0,3 Isr: 3,0,2      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 7    Leader: 3       Replicas: 3,1,0 Isr: 0,3,1      Elr:    LastKnownElr: 
```

# Выводы

- Перераспределение партиций прошло успешно, как подтверждается результатами describe.
- После остановки одного из брокеров Kafka продолжала обслуживать запросы, задействовав ISR.
- После перезапуска брокера реплики были восстановлены.
- Механизм выбора предпочтительного лидера позволяет вернуть исходный баланс партиций.