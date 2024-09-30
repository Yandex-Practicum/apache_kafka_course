# Выполненное задание 1 должно включать:
# Краткое описание выполненных шагов (текстовый файл или документ).

Результаты всех необходимы шагов фиксирую в этом файле.

# Шаг 1. Разверните Kafka. Для этого:
# Разверните кластер Kafka в Yandex Cloud с 3 брокерами.

Данные из yandex cloud:

Описание	 kafka585
Метки	     c9qhd09k2eiqf6eq5hp5

# Укажите оптимальные параметры аппаратных ресурсов для брокеров (количество дисков, CPU, RAM).

Класс хоста s3-c2-m8 (2 vCPU, 100% vCPU rate, 8 ГБ RAM)

# Шаг 2. Настройте репликацию и хранение данных:
# Создайте топик с 3 партициями и коэффициентом репликации 3.
# Настройте политику очистки логов (log.cleanup.policy)/
# Установите параметры хранения (log.retention.ms, log.segment.bytes).

Обзор топика
Имя topic
Количество разделов 3
Фактор репликации 3
Высокая доступность Да
Настройки топика 
Политика очистки лога CLEANUP_POLICY_DELETE
Тип сжатия COMPRESSION_TYPE_UNSPECIFIED
Максимальный размер раздела, байт 10000
Время жизни сегмента лога, мс 10000
Предварительное выделение файла сегмента Нет

# Шаг 3. Настройте Schema Registry:
# Разверните Schema Registry.
# Зарегистрируйте схему данных.

cd ./task1

jq \
    -n --slurpfile data schema-key.json \
    '{
       "schemaType": "AVRO",
       "schema": "\($data)"
    }' \
| curl \
      --request POST \
      -k \
      --url 'https://rc1a-v868adntv9bpbg1i.mdb.yandexcloud.net:443/subjects/topic-key/versions' \
      --user user:user1234 \
      --header 'Content-Type: application/vnd.schemaregistry.v1+json' \
      --data "@-"

{"id":1}


 jq \
    -n --slurpfile data schema-value.json \
    '{
       "schemaType": "AVRO",
       "schema": "\($data)"
    }' \
| curl \
      --request POST \
      -k \
      --url 'https://rc1a-v868adntv9bpbg1i.mdb.yandexcloud.net:443/subjects/topic-value/versions' \
      --user user:user1234 \
      --header 'Content-Type: application/vnd.schemaregistry.v1+json' \
      --data "@-"
{"id":2}



# Шаг 4. Проверьте работу Kafka:
# Напишите простой продюсер и консьюмер.
# Отправьте тестовые сообщения и убедитесь, что они передаются через Kafka.

# Код продюсера и консьюмера.

mkdir -p /usr/local/share/ca-certificates/Yandex && \
wget "https://storage.yandexcloud.net/cloud-certs/CA.pem" \
   --output-document /usr/local/share/ca-certificates/Yandex/YandexInternalRootCA.crt && \
chmod 0655 /usr/local/share/ca-certificates/Yandex/YandexInternalRootCA.crt

sudo apt update && sudo apt install --yes kafkacat

# Скриншоты, подтверждающие успешную передачу сообщений:
# логи продюсера, которые показывают успешную отправку сообщений;

echo "test message" | kcat -P \
    -b rc1a-v868adntv9bpbg1i.mdb.yandexcloud.net:9091 \
    -t topic \
    -k key \
    -X security.protocol=SASL_SSL \
    -X sasl.mechanism=SCRAM-SHA-512 \
    -X sasl.username="user" \
    -X sasl.password="user1234" \
    -X ssl.ca.location=./YandexInternalRootCA.crt -Z


# логи консьюмера, которые показывают, что сообщения успешно прочитаны.

kcat -C \
         -b rc1a-v868adntv9bpbg1i.mdb.yandexcloud.net:9091 \
         -t topic \
         -X security.protocol=SASL_SSL \
         -X sasl.mechanism=SCRAM-SHA-512 \
         -X sasl.username="user" \
         -X sasl.password="user1234" \
         -X ssl.ca.location=./YandexInternalRootCA.crt -Z -K:
key:test message
% Reached end of topic topic [2] at offset 0
% Reached end of topic topic [1] at offset 1
% Reached end of topic topic [0] at offset 0


# Информацию по аппаратным ресурсам.

Класс хоста s3-c2-m8 (2 vCPU, 100% vCPU rate, 8 ГБ RAM)

# Скрипты конфигурации.

Kластер kafka создан через консоль yandex cloud

Данные из yandex cloud:

Описание	 kafka585
Метки	     c9qhd09k2eiqf6eq5hp5

# Описание параметров кластера.

Имя kafka585
Идентификатор c9qhd09k2eiqf6eq5hp5
Дата создания 19.04.2025, в 20:44
Окружение PRODUCTION
Версия 3.5
Реестр схем данных Да
Kafka Rest API Да
Кластер отказоустойчив Да

# Скриншот ответа вызова curl http://localhost:8081/subjects

curl \
    --request GET -k \
    --url 'https://rc1a-v868adntv9bpbg1i.mdb.yandexcloud.net:443/subjects' \
    --user user:user1234 \
    --header 'Accept: application/vnd.schemaregistry.v1+json'

["topic-key","topic-value"]

# и curl -X GET http://localhost:8081/subjects/<название_схемы>/versions.

curl \
    --request GET -k \
    --url 'https://rc1a-v868adntv9bpbg1i.mdb.yandexcloud.net:443/subjects/topic-key/versions' \
    --user user:user1234 \
    --header 'Accept: application/vnd.schemaregistry.v1+json'

[1]

curl \
    --request GET -k \
    --url 'https://rc1a-v868adntv9bpbg1i.mdb.yandexcloud.net:443/subjects/topic-value/versions' \
    --user user:user1234 \
    --header 'Accept: application/vnd.schemaregistry.v1+json'
    
[1]


curl \
    --request GET -k \
    --url 'https://rc1a-v868adntv9bpbg1i.mdb.yandexcloud.net:443/schemas' \
    --user user:user1234 \
    --header 'Accept: application/vnd.schemaregistry.v1+json'

[{"id":1,"schema":"[{\"fields\":[{\"name\":\"id\",\"type\":\"int\"},{\"name\":\"sid\",\"type\":\"string\"}],\"name\":\"my_key\",\"type\":\"record\"}]","schemaType":"AVRO","subject":"topic-key","version":1},{"id":2,"schema":"[{\"fields\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"city\",\"type\":\"string\"},{\"name\":\"age\",\"type\":\"int\"}],\"name\":\"my_value\",\"type\":\"record\"}]","schemaType":"AVRO","subject":"topic-value","version":1}]

# Файл схемы (.avsc или .json).

./task1/schema-key.json
./task1/schema-value.json

# Вывод команды kafka-topics.sh --describe.

yc managed-kafka topic list --cluster-name kafka585
+-------+------------------+--------------------+
| NAME  | PARTITIONS COUNT | REPLICATION FACTOR |
+-------+------------------+--------------------+
| task2 |                3 |                  3 |
| topic |                3 |                  3 |
+-------+------------------+--------------------+
