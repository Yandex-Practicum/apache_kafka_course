# Модуль 3
## Задание 1
| Номер эксперимента | batch.size | linger.ms | compression.type | buffer.memory | Source Record Write Rate (ops/sec) |
|:------------------:|:----------:|:---------:|:----------------:|:-------------:|:----------------------------------:|
|         1          |   16384    |     0     |       none       |   33554432    |                134                 |
|         2          |   16384    |    100    |       none       |   33554432    |                141                 |
|         3          |   16384    |    100    |      snappy      |   33554432    |                143                 |
|         4          |   92000    |    100    |      snappy      |   33554432    |                136                 |
|         5          |   92000    |    300    |      snappy      |   33554432    |                129                 |
|         6          |   92000    |     0     |      snappy      |   33554432    |                133                 |
|         7          |   920000   |     0     |      snappy      |   33554432    |                125                 |
|         8          |   920000   |     0     |      snappy      |   67108864    |                138                 |
|         9          |   128000   |     0     |      snappy      |   33554432    |                140                 |
|         10         |   128000   |     0     |       none       |   33554432    |                133                 |
|         11         |   256000   |     0     |       none       |   33554432    |                130                 |
|         12         |   256000   |     0     |      snappy      |   33554432    |                139                 |
|         13         |   256000   |    100    |      snappy      |   33554432    |                124                 |
|         14         |   32768    |     0     |      snappy      |   33554432    |                132                 |
|         15         |   256000   |    30     |      snappy      |   33554432    |                125                 |
|         16         |   92000    |    200    |      snappy      |   33554432    |                132                 |
|         17         |   92000    |     0     |      snappy      |   33554432    |                132                 |
|         18         |   16384    |    80     |      snappy      |   33554432    |                149                 |


## Задание 2
1. Собрать проект
```bash
mvn clean package
```
2. Переместить собранный артефакт module3-0.0.1-SNAPSHOT.jar в папку infra/confluent-hub-components
3. Запустить докер и проект из папки infra 
```bash
docker compose up
```
4. Перейти в графану по пути localhost:3000 и импортировать infra/grafana/dashboards/Kafka Connect WithCustom.json. 
Первые 4 панели настроены на получение метрик из кастомного коннектора.
5. Запустить коннектор
```bash
curl -X POST "http://localhost:8083/connectors" \
     -H "Content-Type: application/json" \
     --data "@custom_connector_settings.json" | jq
```
6. Забросить в топик prometheus-topic пример сообщения и убедится, что данные появились на дешборде
```bash
cat task2_sample.json | tr -d '\n' | kafka-console-producer.sh --broker-list localhost:9095 --topic prometheus-topic
```
## Задание 3
Лог находится в файле debezium.log