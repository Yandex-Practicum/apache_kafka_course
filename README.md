# Упрощённый сервис обмена сообщениями

## Описание
Этот проект реализует потоковую обработку сообщений с функциональностью:
- Блокировка пользователей: Сообщения от заблокированных пользователей фильтруются.
- Цензура сообщений: Сообщения фильтруются на наличие запрещенных слов.

## Инфраструктура
Для работы используются:
- Kafka: Для хранения и обработки сообщений.
- Kafka Streams: Для потоковой обработки данных.

## Инструкция по запуску
1. Разверните инфраструктуру с помощью Docker Compose:
   ```bash
   docker-compose up -d
   ```

2. Создайте необходимые топики:
```bash
docker exec -it kafka-1 kafka-topics --create --topic messages --bootstrap-server kafka-1:9092 --partitions 1 --replication-factor 1
docker exec -it kafka-1 kafka-topics --create --topic filtered_messages --bootstrap-server kafka-1:9092 --partitions 1 --replication-factor 1
docker exec -it kafka-1 kafka-topics --create --topic blocked_users --bootstrap-server kafka-1:9092 --partitions 1 --replication-factor 1
```
## Тестирование
1. Запустите приложение
2. Установите список заблокированных пользователей
```bash
 docker exec -it kafka-1 kafka-console-producer --broker-list kafka-1:9092 --topic blocked_users --property "parse.key=true" --property "key.separator=:"
 ```
Пример данных:
```text
user1: user3
user2: user1,user3
```
3. Отправьте тестовые данные в топик messages
```bash
 docker exec -it kafka-1 kafka-console-producer --broker-list kafka-1:9092 --topic messages --property "parse.key=true" --property "key.separator=:"
 ```
Пример данных:
```text
user1->user2: Hello
user2->user1: 123 bad
```

4. Проверьте топик filtered_messages
```bash
docker exec -it kafka-1 kafka-console-consumer --bootstrap-server kafka-1:9092 --topic filtered_messages --from-beginning
```