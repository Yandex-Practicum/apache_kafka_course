## Авторское решение на Go

Этот проект демонстрирует систему обработки потоков сообщений с функциональностью 
блокировки пользователей и цензуры сообщений.


### Установка

1. Установите Go с официального сайта: https://go.dev/doc/install

2. Установите требуемые зависимости:
```bash
go mod tidy
```


### Запуск
1. Запустите инфраструктуру
```bash
cd ./infra && docker compose up -d
```

2. Создайте топики: 

Создание топика для отправки сообщения
```bash
docker exec -it kafka-1 ../../usr/bin/kafka-topics --create --topic emitter2filter-stream --bootstrap-server localhost:9092 --partitions 3 --replication-factor 2
```
Создание топика отслеживающего изменения в групповой таблицы сообщений.
```bash
docker exec -it kafka-1 ../../usr/bin/kafka-topics --create --topic user-message-table --bootstrap-server localhost:9092 --partitions 3 --replication-factor 2 --config cleanup.policy=compact
```

Создание топика для передачи данных от фильтра в обработчик пользовательских сообщений
```bash
docker exec -it kafka-1 ../../usr/bin/kafka-topics --create --topic filter2userprocessor-stream --bootstrap-server localhost:9092 --partitions 3 --replication-factor 2
```
Создание топика отслеживающего изменения в групповой таблицы фильтра.
```bash
docker exec -it kafka-1 ../../usr/bin/kafka-topics --create --topic filter-table --bootstrap-server localhost:9092 --partitions 3 --replication-factor 2 --config cleanup.policy=compact
```


Создание топика для передачи данных заблокированных пользователей.
```bash
docker exec -it kafka-1 ../../usr/bin/kafka-topics --create --topic blocker-stream --bootstrap-server localhost:9092 --partitions 3 --replication-factor 2 --config cleanup.policy=compact
```

Создание топика отслеживающего изменения в групповой таблицы заблокированных пользователей.
```bash
docker exec -it kafka-1 ../../usr/bin/kafka-topics --create --topic blocker-table --bootstrap-server localhost:9092 --partitions 3 --replication-factor 2 --config cleanup.policy=compact
```

Создание топика для передачи данных о зацензуринных словах.
```bash
docker exec -it kafka-1 ../../usr/bin/kafka-topics --create --topic censor-stream --bootstrap-server localhost:9092 --partitions 3 --replication-factor 2 --config cleanup.policy=compact
```

Создание топика отслеживающего изменения в групповой таблицы о зацензуринных словах.
```bash
docker exec -it kafka-1 ../../usr/bin/kafka-topics --create --topic censor-table --bootstrap-server localhost:9092 --partitions 3 --replication-factor 2 --config cleanup.policy=compact
```

3. Запустите программу:
```bash
go run main.go
```

Сообщения будут авторматически отправляться, последние 5 сообщений доступны по адресам:
http://localhost:9095/Alex
http://localhost:9095/Dian
http://localhost:9095/Xenia

Ожидаемый вывод:

![messages.png](docs%2Fmessages.png)

4. Цензура, заменяем "sad" на "not so happy"
```bash
go run cmd/censore/main.go -word "sad" -with "not so happy" 
```

5. Блокировка пользователей. Для пользователя Алекс заблокировать все сообщения от Dian и Xenia
```bash
go run cmd/block-user/main.go -user Alex -stream blocker-stream -name Dian 
go run cmd/block-user/main.go -user Alex -stream blocker-stream -name Xenia 
```

6. Разблокировка пользователей. Для пользователя Алекс разблокировать все сообщения от Dian и Xenia
```bash
go run cmd/block-user/main.go -user Alex -stream blocker-stream -name Dian -block=false
go run cmd/block-user/main.go -user Alex -stream blocker-stream -name Xenia -block=false
```
