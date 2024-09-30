## Подготовка

* В первую очередь, необходимо установить docker, docker-compose-plugin ([по инструкции](https://docs.docker.com/engine/install/ubuntu/)) и git
* Добавить пользователя в группу docker: `sudo usermod -aG docker $USER` и перелогиниться
* Скачать этот репозиторий `git clone https://github.com/sakharovaan/yandex-kafka-production.git`
* Скачать сертификат yandex cloud `wget "https://storage.yandexcloud.net/cloud-certs/CA.pem" --output-document yandex-kafka-production/app/ca.crt`
* Посмотреть в консоли Yandex Cloud хост для подключения к кластеру Kafka и создать файл `yandex-kafka-production/app/.env` со следующим содержимым (заменив имя хоста и пароль на свои):
```
PRODUCER_KAFKA_BOOTSTRAP_SERVERS=rc1a-96eh53qo92ss79ic.mdb.yandexcloud.net:9091
PRODUCER_SCHEMA_REGISTRY_SERVER=https://kafka-user:...@rc1a-96eh53qo92ss79ic.mdb.yandexcloud.net:443
PRODUCER_KAFKA_PASSWORD=...

CONSUMER_KAFKA_BOOTSTRAP_SERVERS=rc1a-96eh53qo92ss79ic.mdb.yandexcloud.net:9091
CONSUMER_SCHEMA_REGISTRY_SERVER=https://kafka-user:...@rc1a-96eh53qo92ss79ic.mdb.yandexcloud.net:443
CONSUMER_KAFKA_PASSWORD=...
```
* Выполнить `cd yandex-kafka-production/app`
* Выполнить `docker compose up -d`

## Тестирование

* Пошлём сообщения через producer: `for i in $(seq 1 100); do curl -vv -XPOST -d '{"text": "message"}' -H "Content-Type: application/json" localhost:9000/api/v1/kafka/send; done`
* Посмотрим логи контейнеров `docker compose logs -f`
* В логах продьюсера должны быть сообщения об успешной отправке
```
producer-1     | 2025-06-01 17:37:05.372 | INFO     | producer.endpoints.kafka:send_kafka:26 - "{\"text\":\"message\"}"
producer-1     | INFO:     172.19.0.1:58920 - "POST /api/v1/kafka/send HTTP/1.1" 200 OK
producer-1     | 2025-06-01 17:37:05.383 | INFO     | producer.endpoints.kafka:send_kafka:26 - "{\"text\":\"message\"}"
producer-1     | INFO:     172.19.0.1:58924 - "POST /api/v1/kafka/send HTTP/1.1" 200 OK
```
* В логах консьюмера должны быть сообщения об успешном получении
```
consumer-1     | 2025-06-01 17:37:11.034 | INFO     | src.consumer.service:process_message:16 - {"topic": "kafka-topic", "partition": 1, "offset": 0, "key": null, "message": {"text": "message"}, "timestamp": 1748799424732}
consumer-1     | 2025-06-01 17:37:21.020 | INFO     | src.consumer.service:process_message:16 - {"topic": "kafka-topic", "partition": 0, "offset": 0, "key": null, "message": {"text": "message"}, "timestamp": 1748799424743}
```
* Скриншот:

![logs](README_img/logs.png)

* Выполним `curl -k https://kafka-user:...@rc1a-96eh53qo92ss79ic.mdb.yandexcloud.net:443/subjects`

```
["kafka-topic-message"]
```

* Выполним `curl -k https://kafka-user:...@rc1a-96eh53qo92ss79ic.mdb.yandexcloud.net:443/subjects/kafka-topic-message/versions`

```
[1]
```

* Так как мы не используем файл схемы (producer автоматически его генерирует при помощи библиотеки pydantic-avro), извлечём схему с Schema Registry при помощи команды `curl -k https://kafka-user:...@rc1a-96eh53qo92ss79ic.mdb.yandexcloud.net:443/schemas`

```
[{"id":1,"schema":"{\"fields\":[{\"name\":\"text\",\"type\":\"string\"}],\"name\":\"Message\",\"namespace\":\"Message\",\"type\":\"record\"}","schemaType":"AVRO","subject":"kafka-topic-message","version":1}]
```

* Выполним по заданию команду `kafka-topics.sh --describe`. Сначала поменяем файл admin-config.conf, добавив туда креды пользователя, а затем выполним

```bash
docker compose run -u0 kafka-tools bash
keytool -import -alias yc -keystore "/opt/bitnami/java/lib/security/cacerts" -file /ca.crt
kafka-topics.sh --describe --bootstrap-server rc1a-96eh53qo92ss79ic.mdb.yandexcloud.net:9091 --command-config /admin-config.conf
```

Полный вывод команды:
```

Topic: kafka-topic      TopicId: YbZTH3a5QriRZNQmM6qd-A PartitionCount: 3       ReplicationFactor: 3    Configs: min.insync.replicas=2,cleanup.policy=delete,segment.bytes=1073741824,retention.ms=86400000,retention.bytes=1073741824
        Topic: kafka-topic      Partition: 0    Leader: 1       Replicas: 1,2,3 Isr: 1,2,3      Elr: N/A        LastKnownElr: N/A
        Topic: kafka-topic      Partition: 1    Leader: 2       Replicas: 2,3,1 Isr: 2,3,1      Elr: N/A        LastKnownElr: N/A
        Topic: kafka-topic      Partition: 2    Leader: 3       Replicas: 3,1,2 Isr: 3,1,2      Elr: N/A        LastKnownElr: N/A
Topic: __consumer_offsets       TopicId: VZ69wROGT7CTALTt6ZM6rQ PartitionCount: 50      ReplicationFactor: 3    Configs: compression.type=producer,min.insync.replicas=1,cleanup.policy=compact,segment.bytes=104857600
        Topic: __consumer_offsets       Partition: 0    Leader: 2       Replicas: 2,3,1 Isr: 2,3,1      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 1    Leader: 3       Replicas: 3,1,2 Isr: 3,1,2      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 2    Leader: 1       Replicas: 1,2,3 Isr: 1,2,3      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 3    Leader: 2       Replicas: 2,1,3 Isr: 2,1,3      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 4    Leader: 1       Replicas: 1,3,2 Isr: 1,3,2      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 5    Leader: 3       Replicas: 3,2,1 Isr: 3,2,1      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 6    Leader: 1       Replicas: 1,3,2 Isr: 1,3,2      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 7    Leader: 3       Replicas: 3,2,1 Isr: 3,2,1      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 8    Leader: 2       Replicas: 2,1,3 Isr: 2,1,3      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 9    Leader: 3       Replicas: 3,1,2 Isr: 3,1,2      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 10   Leader: 1       Replicas: 1,2,3 Isr: 1,2,3      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 11   Leader: 2       Replicas: 2,3,1 Isr: 2,3,1      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 12   Leader: 1       Replicas: 1,3,2 Isr: 1,3,2      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 13   Leader: 3       Replicas: 3,2,1 Isr: 3,2,1      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 14   Leader: 2       Replicas: 2,1,3 Isr: 2,1,3      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 15   Leader: 3       Replicas: 3,1,2 Isr: 3,1,2      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 16   Leader: 1       Replicas: 1,2,3 Isr: 1,2,3      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 17   Leader: 2       Replicas: 2,3,1 Isr: 2,3,1      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 18   Leader: 3       Replicas: 3,2,1 Isr: 3,2,1      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 19   Leader: 2       Replicas: 2,1,3 Isr: 2,1,3      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 20   Leader: 1       Replicas: 1,3,2 Isr: 1,3,2      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 21   Leader: 2       Replicas: 2,1,3 Isr: 2,1,3      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 22   Leader: 1       Replicas: 1,3,2 Isr: 1,3,2      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 23   Leader: 3       Replicas: 3,2,1 Isr: 3,2,1      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 24   Leader: 2       Replicas: 2,3,1 Isr: 2,3,1      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 25   Leader: 3       Replicas: 3,1,2 Isr: 3,1,2      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 26   Leader: 1       Replicas: 1,2,3 Isr: 1,2,3      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 27   Leader: 2       Replicas: 2,1,3 Isr: 2,1,3      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 28   Leader: 1       Replicas: 1,3,2 Isr: 1,3,2      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 29   Leader: 3       Replicas: 3,2,1 Isr: 3,2,1      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 30   Leader: 3       Replicas: 3,1,2 Isr: 3,1,2      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 31   Leader: 1       Replicas: 1,2,3 Isr: 1,2,3      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 32   Leader: 2       Replicas: 2,3,1 Isr: 2,3,1      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 33   Leader: 3       Replicas: 3,2,1 Isr: 3,2,1      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 34   Leader: 2       Replicas: 2,1,3 Isr: 2,1,3      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 35   Leader: 1       Replicas: 1,3,2 Isr: 1,3,2      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 36   Leader: 3       Replicas: 3,2,1 Isr: 3,2,1      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 37   Leader: 2       Replicas: 2,1,3 Isr: 2,1,3      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 38   Leader: 1       Replicas: 1,3,2 Isr: 1,3,2      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 39   Leader: 2       Replicas: 2,3,1 Isr: 2,3,1      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 40   Leader: 3       Replicas: 3,1,2 Isr: 3,1,2      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 41   Leader: 1       Replicas: 1,2,3 Isr: 1,2,3      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 42   Leader: 2       Replicas: 2,1,3 Isr: 2,1,3      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 43   Leader: 1       Replicas: 1,3,2 Isr: 1,3,2      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 44   Leader: 3       Replicas: 3,2,1 Isr: 3,2,1      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 45   Leader: 3       Replicas: 3,1,2 Isr: 3,1,2      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 46   Leader: 1       Replicas: 1,2,3 Isr: 1,2,3      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 47   Leader: 2       Replicas: 2,3,1 Isr: 2,3,1      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 48   Leader: 2       Replicas: 2,3,1 Isr: 2,3,1      Elr: N/A        LastKnownElr: N/A
        Topic: __consumer_offsets       Partition: 49   Leader: 3       Replicas: 3,1,2 Isr: 3,1,2      Elr: N/A        LastKnownElr: N/A
```