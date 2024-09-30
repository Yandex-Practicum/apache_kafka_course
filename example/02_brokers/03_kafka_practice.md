# Немного практики

Проще всего получить работающую Kafka через использование контейнеров Docker.

Для этого нужно:

- Docker:
  - [Установленный и запущенный Docker](https://docs.docker.com/engine/installation/){target="_blank"} версии 1.11 или выше.
  - [Установленный Docker Compose](https://docs.docker.com/compose/install/){target="_blank"}.
  - Как минимум 8 Гб оперативной памяти, выделенной в Docker.
- [Git](https://git-scm.com/downloads){target="_blank"}.
- Доступ в Интернет.

## Запускаем кластер Kafka

Начнём с запуска Kafka на одном узле и последовательно будем масштабировать наш кластер.

1. Создадим `docker-compose`-файл для запуска первого узла:

    ```yml
    version: "3.9"
    
    services:
      kafka-0:
       image: bitnami/kafka:3.4
       ports:
         - "9094:9094"
       environment:
         - KAFKA_ENABLE_KRAFT=yes
         - ALLOW_PLAINTEXT_LISTENER=yes
         - KAFKA_CFG_NODE_ID=0
         - KAFKA_CFG_PROCESS_ROLES=broker,controller
         - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
         - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=0@kafka-0:9093
         - KAFKA_KRAFT_CLUSTER_ID=abcdefghijklmnopqrstuv
         - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093,EXTERNAL://:9094
         - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://kafka-0:9092,EXTERNAL://127.0.0.1:9094
         - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,EXTERNAL:PLAINTEXT,PLAINTEXT:PLAINTEXT
    
       volumes:
         - kafka_0_data:/bitnami/kafka
    
    volumes:
      kafka_0_data:
    ```

   Поясним некоторые переменные окружения:
   - `KAFKA_ENABLE_KRAFT=yes` — разрешить использование протокола KRaft;
   - `KAFKA_CFG_PROCESS_ROLE=broker,controller` — узел может входить в кворум как контроллер, но также как брокер обеспечивает хранение разделов и добавление новых сообщений в разделы;
   - `KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER` — определение типа слушателя для публикации контроллера (используется далее в KAFKA_CFG_LISTENERS и KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP);
   - `KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=0@kafka-0:9093` — обозначаем идентификатор контроллера (как части кворума) и его адрес и порт (здесь нужно перечислить адреса всех известных контроллеров);
   - `KAFKA_KRAFT_CLUSTER_ID=somevalue` — идентификатор кластера (должен быть одинаковым у всех контроллеров и брокеров), его мы получим после первого запуска;
   - `KAFKA_CFG_LISTENERS` —  этот параметр определяет, на каких адресах и портах брокер Kafka будет прослушивать входящие соединения. Это включает в себя все типы соединений, внутри докер-сети или локальные. По сути, здесь мы говорим Kafka, где искать входящие запросы. В вашем конфигурационном файле вы указываете три слушателя: `PLAINTEXT://:9092`, `CONTROLLER://:9093` и `EXTERNAL://:9094`. Первые два предназначены для внутренних коммуникаций внутри Docker, а последний — для внешних соединений с хост-машиной;
   - `KAFKA_CFG_ADVERTISED_LISTENERS` —  этот параметр указывает на те адреса и порты, которые брокер Kafka будет «объявлять» или «рекламировать» своим клиентам. Это важно, потому что клиенты Kafka должны знать, как подключиться к брокеру. Если ваш брокер Kafka работает внутри докер-сети, этот параметр может указывать на адрес в этой сети. Если же вам нужно подключаться к брокеру локально, здесь вы можете указать соответствующий порт и протокол. В вашем конфигурационном файле вы объявляете два слушателя `PLAINTEXT://kafka-0:9092` и `EXTERNAL://127.0.0.1:9094`. Первый слушатель используется для обслуживания внутренних коммуникаций в Docker, второй — для обслуживания внешних коммуникаций с хост-машиной;
   - `KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP` — этот параметр позволяет вам определить, какой протокол безопасности будет использоваться для каждого прослушивателя, которые вы определили в параметрах выше. Это важно для обеспечения безопасности ваших данных при передаче между брокером и клиентами. В вашем конфигурационном файле вы указываете, что все ваши слушатели (`CONTROLLER`, `EXTERNAL`, `PLAINTEXT`) будут использовать протокол `PLAINTEXT`, который не предполагает шифрования данных.
   
   Про последние три параметра более подробно можно прочитать в [этой](https://www.confluent.io/blog/kafka-listeners-explained/){target="_blank"} статье.

2. Добавим ещё два узла, каждый из которых будет выполнять функции контроллера и брокера. Для того чтобы определить voters (избирателей), нам придётся перечислить все контроллеры, которые в настоящее время находятся в системе в переменной `KAFKA_CFG_CONTROLLER_QUORUM_VOTERS`. Два новых узла будут обладать похожими настройками, отличия будут в имени контейнера, идентификаторе брокера и используемого volume для хранения данных. Дополненная конфигурация будет выглядеть так:

   ```yml
   services:
     kafka-0:
       image: bitnami/kafka:3.4
       ports:
         - "9094:9094"
       environment:
         - KAFKA_ENABLE_KRAFT=yes
         - KAFKA_CFG_PROCESS_ROLES=broker,controller
         - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
         - ALLOW_PLAINTEXT_LISTENER=yes
         - KAFKA_CFG_NODE_ID=0
         - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=0@kafka-0:9093,1@kafka-1:9093,2@kafka-2:9093
         - KAFKA_KRAFT_CLUSTER_ID=abcdefghijklmnopqrstuv
         - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093,EXTERNAL://:9094
         - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://kafka-0:9092,EXTERNAL://127.0.0.1:9094
         - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,EXTERNAL:PLAINTEXT,PLAINTEXT:PLAINTEXT
      
       volumes:
         - kafka_0_data:/bitnami/kafka
      
      
     kafka-1:
       image: bitnami/kafka:3.4
       ports:
         - "9095:9095"
       environment:
         - KAFKA_ENABLE_KRAFT=yes
         - ALLOW_PLAINTEXT_LISTENER=yes
         - KAFKA_CFG_NODE_ID=1
         - KAFKA_CFG_PROCESS_ROLES=broker,controller
         - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
         - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=0@kafka-0:9093,1@kafka-1:9093,2@kafka-2:9093
         - KAFKA_KRAFT_CLUSTER_ID=abcdefghijklmnopqrstuv
         - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093,EXTERNAL://:9095
         - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://kafka-1:9092,EXTERNAL://127.0.0.1:9095
         - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,EXTERNAL:PLAINTEXT,PLAINTEXT:PLAINTEXT
      
       volumes:
         - kafka_1_data:/bitnami/kafka
      
     kafka-2:
       image: bitnami/kafka:3.4
       ports:
         - "9096:9096"
       environment:
         - KAFKA_ENABLE_KRAFT=yes
         - ALLOW_PLAINTEXT_LISTENER=yes
         - KAFKA_CFG_NODE_ID=2
         - KAFKA_CFG_PROCESS_ROLES=broker,controller
         - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
         - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=0@kafka-0:9093,1@kafka-1:9093,2@kafka-2:9093
         - KAFKA_KRAFT_CLUSTER_ID=abcdefghijklmnopqrstuv
         - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093,EXTERNAL://:9096
         - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://kafka-2:9092,EXTERNAL://127.0.0.1:9096
         - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,EXTERNAL:PLAINTEXT,PLAINTEXT:PLAINTEXT
       volumes:
         - kafka_2_data:/bitnami/kafka   
   
   volumes:
     kafka_0_data:
     kafka_1_data:
     kafka_2_data:
   ```

3. Запускаем кластер.

    ```bash
    docker-compose up -d
    ```

4. Проверим успешность запуска. Введём в терминале:

   ```bash
   docker ps
   ```
      
   В терминале должно быть что-то вроде:
   
   ```shell
   CONTAINER ID   IMAGE               COMMAND                  CREATED          STATUS          PORTS                     NAMES
   33f26b857624   bitnami/kafka:3.4   "/opt/bitnami/script…"   26 seconds ago   Up 25 seconds   0.0.0.0:9094->9092/tcp   yap-kafka-1-1
   08de103bcada   bitnami/kafka:3.4   "/opt/bitnami/script…"   26 seconds ago   Up 25 seconds   0.0.0.0:9095->9092/tcp   yap-kafka-2-1
   d8153faf348d   bitnami/kafka:3.4   "/opt/bitnami/script…"   26 seconds ago   Up 25 seconds   0.0.0.0:9096->9092/tcp   yap-kafka-0-1
   ```

5. Подключим интерфейс для взаимодействия с Kafka. Добавим новый сервис в `docker-compose`-файл:

   ```yml
   ui:
     image: provectuslabs/kafka-ui:v0.7.0
     ports:
       - "8080:8080"
     environment:
       - KAFKA_CLUSTERS_0_BOOTSTRAP_SERVERS=kafka-0:9092
       - KAFKA_CLUSTERS_0_NAME=kraft
   ```

Теперь по адресу `http://localhost:8080` у нас доступен интерфейс для управления Kafka.

![](https://github.com/Yandex-Practicum/python_p2f_ugc/blob/77942931f003fff70972826778916c38973a325c/8th-sprint-refactor/02_brokers/images/kafka_ui_fin.png)

## Создание топика

Создадим наш первый топик. Перейдем во вкладку Topics и верхнем правом углу кликнем на **Add a topic**. В открывшемся окне поставьте следующие настройки:

![](https://github.com/Yandex-Practicum/python_p2f_ugc/blob/77942931f003fff70972826778916c38973a325c/8th-sprint-refactor/02_brokers/images/topic_fin.png)

Мы создаем топик, состоящий из трёх партиций. Давайте разберём каждый выбраный параметр:

### Replications Factor

Фактор репликации определяет количество копий каждой партиции в Kafka, которые в свою очередь хранят сообщения. Если фактор репликации равен 1, тогда каждая партиция (включая все её сообщения) хранится только на одном брокере и, если этот брокер выйдет из строя, все сообщения в этой партиции будут утеряны. Поэтому, как минимум, фактор репликации должен быть больше 1 для обеспечения отказоустойчивости. Обычно фактор репликации устанавливают равным 3, что обеспечивает хороший баланс между отказоустойчивостью и затратами на хранение. Мы можем себе это позволить, потому что в нашем кластере сейчас ровно 3 узла.

### Min-in-sync-replicas

Определяет минимальное количество реплик, которые должны подтвердить запись, прежде чем запись считается успешной. Например, если у вас фактор репликации равен 3 и min-in-sync-replicas установлено в 2, то запись будет считаться успешной только после того, как две реплики подтвердят запись. Этот параметр важен для обеспечения консистентности данных в Kafka. Если этот параметр слишком высок, это может замедлить производительность, так как Kafka будет ждать подтверждения от большего числа реплик. А если поставить его намного больше, чем число самих реплик, то Kafka вообще сломается. Обычно, если вы установите фактор репликации в 3, хорошей отправной точкой для min-in-sync-replicas будет 2.

### Cleanup policy

Политика очистки определяет, что происходит с записями, которые превышают установленный срок хранения или размер. В нашем случае, мы выбираем delete, что означает, что устаревшие или превысившие установленный размер записи будут удалены. Существует также политика compact. В этом режиме, вместо того чтобы удалять устаревшие сообщения, Kafka сохраняет самую последнюю запись для каждого уникального ключа в партиции. Это особенно полезно для ситуаций, где важно сохранить состояние, представленное последним сообщением каждого ключа, независимо от того, превышает ли он установленные ограничения по срокам хранения или размеру.

### Time to retain

Время хранения данных  определяет период времени (в миллисекундах), в течение которого данные будут храниться на сервере, прежде чем они будут автоматически удалены в соответствии с политикой очистки. В нашем случае мы установили этот параметр на уровне 86400000мс, что соответствует одним суткам. В Kafka время сохранения сообщения может быть основано на времени его создания (англ. create time), что устанавливается производителем сообщения, или на времени добавления в журнал (англ. log-append time), что устанавливается брокером при получении сообщения. Выбор между этими двумя режимами влияет на то, как будет применяться политика срока хранения. Если выбрано время создания, то сообщения удаляются после того, как истекает срок их хранения от момента их создания. Если же выбрано время добавления в журнал, то сообщения удаляются после того, как истекает срок их хранения от момента их добавления в журнал Kafka. 

Также существует параметр `retention.bytes`: он определяет максимальный размер данных, который может быть сохранен в партиции, прежде чем старые данные начнут удаляться. 

При необходимости можно установить и другие ограничения, например, общий размер сообщений на диске или максимальный размер одного сообщения.

После создания топика отправим в него первое сообщение нажатием на копку **Produce Message** в верхнем правом углу.

![](https://github.com/Yandex-Practicum/python_p2f_ugc/blob/77942931f003fff70972826778916c38973a325c/8th-sprint-refactor/02_brokers/images/message_fin.png)

На данный момент у нас один топик с тремя партициями и одним сообщением.

![](https://github.com/Yandex-Practicum/python_p2f_ugc/blob/77942931f003fff70972826778916c38973a325c/8th-sprint-refactor/02_brokers/images/topic-info_fin.png)


## Публикуем и читаем события из топика

Чтобы публиковать события в топик и читать из него события, будем использовать библиотеку `kafka-python`. Установите её, используя свой любимый пакетный менеджер.

```sh
pip install kafka-python==2.0.2
```

Следующий сниппет показывает принципы использования библиотеки `kafka-python` для публикации сообщений в Kafka.

```python
from kafka import KafkaProducer
from time import sleep


producer = KafkaProducer(bootstrap_servers=['localhost:9094'])

producer.send(
    topic='messages',
    value=b'my message from python',
    key=b'python-message',
)

sleep(1)
```

Создаём объект класса `KafkaProducer` и указываем ему в качестве подключения хост и порт, по которому слушает сообщения брокер Kafka.

Используем метод `send` для отправки данных.

- `topic` — имя топика, в который нужно отправить данные.
- `value` — сами данные для отправки. Формат данных может быть любой. В этом примере это просто бинарная строка.
- `key` — это ключ партиционирования данных. Kafka гарантирует, что сообщения с одним и тем же ключом партиционирования будут отправлены в одну и ту же партицию. 

Чтобы убедиться, что сообщение действительно было отправлено в топик, откройте страницу топика `messages` в панели управления.

Считать сообщения из топика можно с помощью такого сниппета:

```python
from kafka import KafkaConsumer


consumer = KafkaConsumer(
    'messages',
    bootstrap_servers=['localhost:9094'],
    auto_offset_reset='earliest',
    group_id='echo-messages-to-stdout',
)

for message in consumer:
    print(message.value)
```

Создадим объект класса `KafkaConsumer` и укажем ему адрес для подключения к брокеру.

- `auto_offset_reset` задаёт, с какой позиции нужно начать считывать сообщения. В этом примере сообщения будут считываться с самого начала.
- `group_id` задаёт имя группы, к которой относится этот `consumer`. В примере консьюмер просто выводит сообщения в стандартный поток вывода, что и отражено в имени.

Запуск сниппета приводит к выводу в консоль следующих данных:

```sh
b'my message from python'
```

Это не что иное, как тело сообщения, которое вы отправили ранее в топик.
