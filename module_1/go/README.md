## Проект Kafka Avro на Go

Этот проект демонстрирует работу с Kafka и Avro на языке Go.

### Описание

Код проекта можно скомпилировать в два независимо запускаемых приложений:

* Продьюсер (Producer): Записывает данные в формате Avro в топик Kafka.
* Консьюмер (Consumer): Считывает данные из топика Kafka и выводит их в консоль. 

### Установка

1. Установите Go с официального сайта: https://go.dev/doc/install

2. Установите требуемые зависимости:
```bash
go mod tidy
```

3. Настройте конфигурационные файлы (можно оставить без изменений):
    * `config/имя.yaml`: файл конфигурации для подключения к Kafka и установки параметров.

### Запуск
1. Запустите инфраструктуру
```bash
cd ./infra && docker compose up -d
```

2. Создайте топик и проверьте его конфигурацию:

```bash
 docker exec -it kafka-1 ../../usr/bin/kafka-topics --create --topic users --bootstrap-server kafka-1:9092 --partitions 3 --replication-factor 2
```
Параметры `localhost:9092` - адрес Kafka-брокера, `users` - имя топика, `3` - количество партиций, `2` - количество реплик.

```bash
docker exec -it kafka-1 ../../usr/bin/kafka-topics  --describe --topic users  --bootstrap-server localhost:9092
```
Ожидаемый результат:
```
Topic: users    TopicId: YEE8ywCDR2mWWht0TNlZsw PartitionCount: 3       ReplicationFactor: 2    Configs: 
        Topic: users    Partition: 0    Leader: 2       Replicas: 2,1   Isr: 2,1
        Topic: users    Partition: 1    Leader: 1       Replicas: 1,2   Isr: 1,2
        Topic: users    Partition: 2    Leader: 2       Replicas: 2,1   Isr: 2,1
```

3. Запустите Producer:
```bash
    go run ./avro-example/cmd/main.go -c ./avro-example/config/producer.yaml
```
    а. Программа у вас запросит требуемое действие (Command):   
    введите `send`, чтобы отправить сообщение в брокер
    введите `exit`, чтобы выйти

    б. Программа у вас запросит требуемые данные для отправки - имя (Enter name:), 
    любимое число (Enter favorite number:), любимый цвет (Enter favorite color:). Пожалуйста,
    заполняйте правильно вводимые значения, так как основная задача показать возможность передачи
    данных в формате Avro, используя Kafka, валидация данных не сделана. 
    
    Пример:
    ```
    Command:send
    Enter name: alex
    Enter favorite number: 55
    Enter favorite color:black
    ```


4. Запустите consumer-push:
```bash
   go run ./avro-example/cmd/main.go -c ./avro-example/config/consumer1.yaml 
```
5. Запустите consumer-pull:
```bash
   go run ./avro-example/cmd/main.go -c ./avro-example/config/consumer2.yaml 
```
