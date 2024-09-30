# Задание 2. Интеграция Kafka с внешними системами (Apache Hadoop)

Проект демонстрирует интеграцию Apache Kafka с Apache Hadoop HDFS для потоковой обработки данных.

## Архитектура

- **Kafka**: Брокер сообщений для потоковой передачи данных
- **Hadoop HDFS**: Распределенная файловая система для хранения данных
- **Producer**: Go приложение для отправки сообщений в Kafka
- **Consumer**: Go приложение для чтения из Kafka и записи в HDFS

## Компоненты

### Kafka кластер
- 1 брокер Kafka (порт 9094)
- Топик: `hadoop-topic`

### Hadoop кластер
- 1 NameNode (порт 9000, веб-интерфейс 9870)
- 3 DataNode (порты 9864, 9865, 9866)

## Запуск системы

### 1. Запуск инфраструктуры

```bash
# Запуск всех контейнеров
docker-compose up -d

# Проверка статуса
docker-compose ps
```

### 2. Запуск приложений

#### Отдельные скрипты:

```bash
# Запуск продюсера (отправка сообщений)
chmod +x run_producer.sh
./run_producer.sh

# Запуск консьюмера (чтение и запись в HDFS)
chmod +x run_consumer.sh
./run_consumer.sh
```

### 3. Полный тест интеграции
## Проверка результатов

### Kafka топики
```bash
# Просмотр сообщений в топике
docker exec kafka-0 kafka-console-consumer.sh \
    --bootstrap-server localhost:9092 \
    --topic hadoop-topic \
    --from-beginning
```

### HDFS данные
```bash
# Список файлов в HDFS
docker exec hadoop-namenode hdfs dfs -ls /data

# Содержимое файла
docker exec hadoop-namenode hdfs dfs -cat /data/message_<uuid>
```

### Веб-интерфейсы
- Hadoop NameNode: http://localhost:9870
- HDFS Browser: http://localhost:9870/explorer.html#/data

## Логи
- `producer_logs.txt` - логи продюсера
- `consumer_logs.txt` - логи консьюмера
- `kafka_consumer_logs.txt` - логи Kafka консьюмера
- `hdfs_ls_logs.txt` - список файлов в HDFS

