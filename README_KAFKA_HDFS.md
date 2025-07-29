# Kafka to HDFS Consumer

Проект на Java с Maven для чтения данных из Apache Kafka и записи их в Apache HDFS.

## Описание

Этот проект представляет собой consumer приложение, которое:
- Подключается к Apache Kafka и читает сообщения из указанного топика
- Обрабатывает сообщения батчами для оптимизации производительности
- Записывает данные в Apache HDFS в формате JSON
- Поддерживает настройку через конфигурационные файлы
- Включает логирование и обработку ошибок

## Структура проекта

```
kafka-to-hdfs/
├── src/
│   ├── main/
│   │   ├── java/com/kafka/hdfs/
│   │   │   ├── KafkaToHdfsConsumer.java      # Основной класс consumer
│   │   │   ├── KafkaToHdfsConsumerV2.java    # Улучшенная версия с конфигурацией
│   │   │   ├── MessageData.java              # Модель данных сообщения
│   │   │   └── Config.java                   # Управление конфигурацией
│   │   └── resources/
│   │       ├── application.properties        # Конфигурация приложения
│   │       └── logback.xml                   # Конфигурация логирования
│   └── test/java/com/kafka/hdfs/
│       └── KafkaToHdfsConsumerTest.java      # Тесты
├── pom.xml                                   # Maven конфигурация
├── run_kafka_hdfs_consumer.sh               # Скрипт запуска для Linux/Mac
├── run_kafka_hdfs_consumer.ps1              # Скрипт запуска для Windows
└── README_KAFKA_HDFS.md                     # Этот файл
```

## Требования

- Java 11 или выше
- Apache Maven 3.6+
- Apache Kafka (локально или удаленно)
- Apache Hadoop HDFS (локально или удаленно)

## Установка и сборка

1. Клонируйте репозиторий или скопируйте файлы проекта
2. Перейдите в директорию проекта
3. Выполните сборку:

```bash
mvn clean package
```

## Конфигурация

### Основные параметры

Создайте файл `src/main/resources/application.properties` или используйте параметры командной строки:

```properties
# Kafka Configuration
kafka.bootstrap.servers=localhost:9092
kafka.topic=test-topic
kafka.group.id=kafka-hdfs-consumer
kafka.auto.offset.reset=earliest
kafka.enable.auto.commit=false

# HDFS Configuration
hdfs.output.path=/kafka-data
hdfs.namenode.url=hdfs://namenode:9000

# Performance Configuration
batch.size=100
flush.interval.ms=5000
poll.timeout.ms=1000
```

### Параметры Kafka

- `kafka.bootstrap.servers` - адреса Kafka брокеров
- `kafka.topic` - имя топика для чтения
- `kafka.group.id` - ID группы consumer
- `kafka.auto.offset.reset` - стратегия сброса offset (earliest/latest)
- `kafka.enable.auto.commit` - автоматическое подтверждение offset

### Параметры HDFS

- `hdfs.output.path` - путь в HDFS для записи данных
- `hdfs.namenode.url` - URL NameNode HDFS

### Параметры производительности

- `batch.size` - размер батча для записи в HDFS
- `flush.interval.ms` - интервал принудительной записи (мс)
- `poll.timeout.ms` - таймаут опроса Kafka (мс)

## Запуск

### Способ 1: Использование скриптов

#### Linux/Mac:
```bash
chmod +x run_kafka_hdfs_consumer.sh
./run_kafka_hdfs_consumer.sh [bootstrap-servers] [topic] [group-id] [hdfs-path] [batch-size] [flush-interval]
```

#### Windows PowerShell:
```powershell
.\run_kafka_hdfs_consumer.ps1 [bootstrap-servers] [topic] [group-id] [hdfs-path] [batch-size] [flush-interval]
```

### Способ 2: Прямой запуск через Maven

```bash
mvn exec:java -Dexec.mainClass="com.kafka.hdfs.KafkaToHdfsConsumer" \
    -Dexec.args="localhost:9092 test-topic kafka-hdfs-consumer /kafka-data 100 5000"
```

### Способ 3: Запуск JAR файла

```bash
java -cp target/kafka-to-hdfs-1.0.0.jar \
     -Dlogback.configurationFile=src/main/resources/logback.xml \
     com.kafka.hdfs.KafkaToHdfsConsumer \
     localhost:9092 test-topic kafka-hdfs-consumer /kafka-data 100 5000
```

## Примеры использования

### Базовый запуск с параметрами по умолчанию:
```bash
./run_kafka_hdfs_consumer.sh
```

### Запуск с кастомными параметрами:
```bash
./run_kafka_hdfs_consumer.sh kafka:9092 my-topic my-group /my-data 200 3000
```

### Запуск с конфигурационным файлом (V2):
```bash
java -cp target/kafka-to-hdfs-1.0.0.jar com.kafka.hdfs.KafkaToHdfsConsumerV2 config.properties
```

## Формат выходных данных

Данные записываются в HDFS в формате JSON, по одному сообщению на строку:

```json
{"topic":"test-topic","partition":0,"offset":123,"key":"message-key","value":"message-value","timestamp":1640995200000,"processed_at":1640995201000}
{"topic":"test-topic","partition":0,"offset":124,"key":"message-key-2","value":"message-value-2","timestamp":1640995202000,"processed_at":1640995203000}
```

## Мониторинг и логирование

Приложение создает логи в директории `logs/`:
- `kafka-hdfs-consumer.log` - основной лог файл
- Ротация логов по дням

Уровни логирования настраиваются в `src/main/resources/logback.xml`.

## Тестирование

Запуск тестов:
```bash
mvn test
```

## Устранение неполадок

### Проблемы подключения к Kafka
- Проверьте доступность Kafka брокеров
- Убедитесь в правильности адресов в конфигурации
- Проверьте существование топика

### Проблемы подключения к HDFS
- Проверьте доступность NameNode
- Убедитесь в правильности URL HDFS
- Проверьте права доступа к директории записи

### Проблемы производительности
- Увеличьте размер батча для лучшей производительности
- Уменьшите интервал записи для более частой записи
- Мониторьте использование памяти и CPU

## Расширение функциональности

### Добавление новых форматов данных
Создайте новые классы в пакете `com.kafka.hdfs` и измените логику сериализации в `KafkaToHdfsConsumer`.

### Добавление фильтрации
Добавьте логику фильтрации в метод `processBatch()`.

### Добавление трансформации данных
Создайте классы трансформации и используйте их в обработке сообщений.

## Лицензия

Этот проект распространяется под лицензией Apache 2.0.

## Поддержка

При возникновении проблем создайте issue в репозитории проекта или обратитесь к документации Apache Kafka и Apache Hadoop. 