# RabbitMQ на практике

Развернём RabbitMQ локально и посмотрим, что он из себя представляет. Разворачивать будем с помощью уже знакомого Docker.

## Запускаем RabbitMQ

1. Выполните следующую команду в терминале:

    ```bash
    docker run -d --name rabbitmq -p 15672:15672 -p 5672:5672 -e RABBITMQ_DEFAULT_USER=user -e RABBITMQ_DEFAULT_PASS=password rabbitmq:3-management

    ```

   Вывод в консоль должен быть примерно таким:
   
    ```shell
    Unable to find image 'rabbitmq:3-management' locally
    3-management: Pulling from library/rabbitmq
    f4bb4e8dca02: Pull complete 
    4ee9fe00ca90: Pull complete 
    16a6499274d8: Pull complete 
    b969ac944c5e: Pull complete 
    99b654b9971e: Pull complete 
    cd9b45c35c7f: Pull complete 
    556aa769760d: Pull complete 
    3cc127251176: Pull complete 
    6aea37dc4e08: Pull complete 
    63b80c91e123: Pull complete 
    b0f4d036f4a8: Pull complete 
    cdda86baa465: Pull complete 
    Digest: sha256:18d7104751b66c882c109349f537108c7cd979d87fe9020ef4dc4d773d37691e
    Status: Downloaded newer image for rabbitmq:3-management
    7c7f6c12b077bcd470239a9a82c3210f9b0538c5c63fa407afb9042d6e0e6db3
    ```

2. Убедитесь, что «кролик» успешно запустился и работает — запустите команду:

    ```bash
    docker ps
    ```
   
    Вывод будет похож на этот:

    ```shell
    CONTAINER ID   IMAGE                   COMMAND                  CREATED          STATUS          PORTS                                                                                                         NAMES
    7c7f6c12b077   rabbitmq:3-management   "docker-entrypoint.s…"   15 minutes ago   Up 15 minutes   4369/tcp, 5671/tcp, 0.0.0.0:5672->5672/tcp, 15671/tcp, 15691-15692/tcp, 25672/tcp, 0.0.0.0:15672->15672/tcp   rabbitmq
    ```

Теперь RabbitMQ готов к работе. Интерфейс будет доступен в браузере по адресу [http://localhost:15672/](http://localhost:15672/){target="_blank"}. 

{quiz-task} background: | #82bbf2 content: | Обратите внимание, что при запуске контейнера мы определили пользователя и пароль для работы (**user**/**password**) — их мы и будем использовать в коде приложений. {/quiz-task}

Нам понадобится инструментарий для работы с RabbitMQ. Будем использовать популярную библиотеку [pika](https://pika.readthedocs.io/en/stable/){target="_blank"}. Установите её.

```bash
pip install pika==1.3.2
```

## Создание очереди, публикация сообщения и маршрутизация

Сообщения не публикуются непосредственно в очередь. Вместо этого паблишер отправляет сообщения в распределитель (Exchange). 

**Exchange** — агенты маршрутизации сообщений, определённые для виртуального хоста RabbitMQ. Exchange отвечает за маршрутизацию сообщений в разные очереди. 
Exchange принимает сообщения от приложения-публикатора и направляет их в очереди сообщений с помощью атрибутов заголовка, биндинга и ключей маршрутизации. 

**Связывание** (биндинг) — это «ссылка», которую настраивают для привязки очереди к Exchange.

**Ключ маршрутизации** — это атрибут сообщения (может не использоваться), нужен распределителю (Exchange) для маршрутизации (выбора одной или нескольких очередей) сообщений.

Поработаем на практике с различными типами распределителей.

### Default Exchange

**Default Exchange** или Обмен по умолчанию — это предварительно объявленный прямой обмен без имени; определен пустой строкой `«»`. 

Когда вы используете обмен по умолчанию, ваше сообщение доставляется в очередь с именем, равным ключу маршрутизации сообщения. Каждая очередь автоматически привязывается к обмену по умолчанию с ключом маршрутизации, который совпадает с именем очереди.

Создадим два файла с продюсером сообщений и получателем:

```python
# default/producer_default.py

from time import sleep

from pika import BlockingConnection, ConnectionParameters, PlainCredentials

if __name__ == '__main__':
    credentials = PlainCredentials('user', 'password')
    connection = BlockingConnection(ConnectionParameters(host='localhost', credentials=credentials))
    channel = connection.channel()

    channel.queue_declare(queue='hello_default_exchange')

    counter = 0
    while True:
        counter += 1
        message_body = f'Привет, Практикум! {counter}'
        # exchange не указываем, а значит будет использоваться Default Exchange,
        # сообщение будет отправлено в очередь с именем равным routing_key
        channel.basic_publish(exchange='', routing_key='hello_default', body=message_body)
        print(f'[✅] {message_body}')
        sleep(1)
```

В коде создаётся соединение с брокером сообщений с определёнными ранее именем пользователя и паролем, и описывается очередь `hello_default_exchange`. 
Далее в вечном цикле в эту очередь отправляются сообщения с текстом `Привет, Практикум!` и счётчиком.

```python
# default/consumer_default.py

from pika import BlockingConnection, ConnectionParameters, PlainCredentials
from pika.adapters.blocking_connection import BlockingChannel
from pika.spec import BasicProperties, Basic


def handler(ch: BlockingChannel, method: Basic.Deliver, properties: BasicProperties, body: bytes) -> None:
    print(f'[🎉] Получено: {body.decode()}')


if __name__ == '__main__':
    credentials = PlainCredentials('user', 'password')
    connection = BlockingConnection(ConnectionParameters(host='localhost', credentials=credentials))
    channel = connection.channel()
    channel.queue_declare(queue='hello_default')

    channel.basic_consume(queue='hello_default', on_message_callback=handler, auto_ack=True)
    channel.start_consuming()
```

В коде консьюмера также определяем подключение к брокеру и описываем очередь. 

Создание очереди с помощью `queue_declare` является [**идемпотентным**](https://ru.wikipedia.org/wiki/%D0%98%D0%B4%D0%B5%D0%BC%D0%BF%D0%BE%D1%82%D0%B5%D0%BD%D1%82%D0%BD%D0%BE%D1%81%D1%82%D1%8C){target="_blank"} — мы можем запускать команду столько раз, сколько захотим, и будет создана только **одна** очередь. В нашем случае мы это делаем и в консьюмере, и в продюсере, так как не можем знать наверняка, кто из них начнёт работать раньше.

Запускаем продюсер и два инстанса консьюмеров (на практике их бывает и больше):

![](https://github.com/Yandex-Practicum/python_p2f_ugc/blob/0ef4e078f807861eefc547e2bbd605ab16702c23/8th-sprint-refactor/02_brokers/images/images_new/default1_new.png)

Сообщения распределяются между несколькими консьюмерами по алгоритму [round-robin](https://ru.wikipedia.org/wiki/Round-robin_(%D0%B0%D0%BB%D0%B3%D0%BE%D1%80%D0%B8%D1%82%D0%BC)){target="_blank"}.

![](https://github.com/Yandex-Practicum/python_p2f_ugc/blob/0ef4e078f807861eefc547e2bbd605ab16702c23/8th-sprint-refactor/02_brokers/images/images_new/defaultExchange_new.mp4)

### Direct Exchange

**Direct exchange** — используется, когда нужно доставить сообщение в конкретные очереди. Сообщение публикуется в Exchange с определённым ключом маршрутизации и попадает во все очереди, которые связаны с этим обменником аналогичным ключом маршрутизации. 

**Ключ маршрутизации** — это атрибут сообщения, добавленный в заголовок сообщения паблишера. Ключ маршрутизации можно рассматривать как «адрес», который использует Exchange, чтобы решить, как маршрутизировать сообщение. Сообщение отправляется в очередь (или очереди), чей ключ точно соответствует ключу маршрутизации сообщения.

Создадим систему логирования с распределением логов по очередям.

```python
# direct/producer_direct.py

import random
from time import sleep

from pika import BlockingConnection, ConnectionParameters, PlainCredentials
from pika.exchange_type import ExchangeType

EXCHANGE_NAME = 'direct_logs'


if __name__ == '__main__':
    credentials = PlainCredentials('user', 'password')
    connection = BlockingConnection(ConnectionParameters(host='localhost', credentials=credentials))
    channel = connection.channel()

    channel.exchange_declare(exchange=EXCHANGE_NAME, exchange_type=ExchangeType.direct)

    counter = 0
    while True:
        counter += 1
        message_type = random.choice(['error', 'warning', 'info'])
        message_body = f'Сообщение {message_type}! {counter}'
        channel.basic_publish(exchange=EXCHANGE_NAME, routing_key=message_type, body=message_body)
        print(f'[✅] {message_body}')
        sleep(1)
```

В цикле публикуем сообщения с различными ключами маршрутизации.

```python
# direct/consumer_direct.py

import sys

from pika import BlockingConnection, ConnectionParameters, PlainCredentials
from pika.adapters.blocking_connection import BlockingChannel
from pika.exchange_type import ExchangeType
from pika.spec import BasicProperties, Basic

EXCHANGE_NAME = 'direct_logs'


def handler(ch: BlockingChannel, method: Basic.Deliver, properties: BasicProperties, body: bytes) -> None:
    print(f'[🎉] Получено: {body.decode()}')


if __name__ == '__main__':
    credentials = PlainCredentials('user', 'password')
    connection = BlockingConnection(ConnectionParameters(host='localhost', credentials=credentials))
    channel = connection.channel()

    channel.exchange_declare(exchange=EXCHANGE_NAME, exchange_type=ExchangeType.direct)
    queue = channel.queue_declare(queue='', exclusive=True)
    # обратите внимание: выше мы не указали имя очереди, в этом случае имя сгенерируется автоматически
    # и будет иметь вид amq.gen-fxDXeobOVBDZTpWfVM6iDA
    queue_name = queue.method.queue

    # получим нужные ключи маршрутизации из аргументов запуска приложения
    routing_keys = sys.argv[1:] or ['error', 'warning', 'info']
    for key in routing_keys:
        channel.queue_bind(exchange=EXCHANGE_NAME, queue=queue_name, routing_key=key)

    channel.basic_consume(queue=queue_name, on_message_callback=handler, auto_ack=True)
    channel.start_consuming()
```

В консьюмере ключи маршрутизации, с которыми он будет работать, получаем из аргументов запуска приложения.

Запустим продюсер и четыре консьюмера для ключей `info`, `warning`, `error` и для всех трёх ключей одновременно, указав их через пробел при запуске `consumer_direct.py`.

![](https://github.com/Yandex-Practicum/python_p2f_ugc/blob/0ef4e078f807861eefc547e2bbd605ab16702c23/8th-sprint-refactor/02_brokers/images/images_new/direct1_new.png)

Порядок маршрутизации сообщений можно представить так:

![](https://github.com/Yandex-Practicum/python_p2f_ugc/blob/0ef4e078f807861eefc547e2bbd605ab16702c23/8th-sprint-refactor/02_brokers/images/images_new/directExchange_new.mp4)

### Topic exchange

**Topic exchange** — аналогично Direct exchange помогает осуществить выборочную маршрутизацию путём сравнения ключа маршрутизации. Но в этом случае ключ задаётся по шаблону. 

При создании шаблона используются 0 или более слов (буквы A-Z и a-z и цифры 0-9), разделённых точкой, а также символы `*` и `#`.

* `*` — может быть заменен на 1 слово;
* `#` — может быть заменен на 0 или более слов.

Консьюмеры указывают на то, какие темы (топики) им интересны. Консьюмер создаёт очередь и устанавливает привязку с заданным шаблоном маршрутизации к обмену. 
Все сообщения с ключом маршрутизации, которые соответствуют шаблону маршрутизации, будут перенаправляться в эту очередь.

Применение этого типа обменника может стать хорошим выбором для возможного будущего развития приложения, потому что шаблоны всегда можно настроить так, чтобы сообщение публиковалось аналогично Direct Exchange или Fanout Exchange, который рассмотрим далее.

```python
# topic/producer_topic.py

import random
from time import sleep

from pika import BlockingConnection, ConnectionParameters, PlainCredentials
from pika.exchange_type import ExchangeType

EXCHANGE_NAME = 'topic_notifications'


if __name__ == '__main__':
    credentials = PlainCredentials('user', 'password')
    connection = BlockingConnection(ConnectionParameters(host='localhost', credentials=credentials))
    channel = connection.channel()

    channel.exchange_declare(exchange=EXCHANGE_NAME, exchange_type=ExchangeType.topic)

    counter = 0
    while True:
        counter += 1
        routing_key = random.choice(
            ['notification.instant.telegram', 'notification.instant.email', 'notification.delayed.email']
        )

        message_body = f'Сообщение {routing_key=:<30}: {counter}'
        channel.basic_publish(exchange=EXCHANGE_NAME, routing_key=routing_key, body=message_body)
        print(f'[✅] {message_body}')
        sleep(1)
```

```python
# topic/consumer_topic.py

import sys

from pika import BlockingConnection, ConnectionParameters, PlainCredentials
from pika.adapters.blocking_connection import BlockingChannel
from pika.exchange_type import ExchangeType
from pika.spec import BasicProperties, Basic

EXCHANGE_NAME = 'topic_notifications'


def handler(ch: BlockingChannel, method: Basic.Deliver, properties: BasicProperties, body: bytes) -> None:
    print(f'[🎉] Получено: {body.decode()}')


if __name__ == '__main__':
    credentials = PlainCredentials('user', 'password')
    connection = BlockingConnection(ConnectionParameters(host='localhost', credentials=credentials))
    channel = connection.channel()

    channel.exchange_declare(exchange=EXCHANGE_NAME, exchange_type=ExchangeType.topic)
    queue = channel.queue_declare(queue='', exclusive=True)

    queue_name = queue.method.queue

    routing_key = sys.argv[1:][0] if sys.argv[1:] else 'notification.*.*'
    channel.queue_bind(exchange=EXCHANGE_NAME, queue=queue_name, routing_key=routing_key)

    channel.basic_consume(queue=queue_name, on_message_callback=handler, auto_ack=True)
    channel.start_consuming()
```

Запустим продюсер и два консьюмера с ключами маршрутизации в виде шаблонов `"*.instant.*"` и `"#.telegram"`. В первый будут приходить все срочные сообщения, а во второй — все сообщения для мессенджера Telegram, вне зависимости от срочности.

![](https://github.com/Yandex-Practicum/python_p2f_ugc/blob/0ef4e078f807861eefc547e2bbd605ab16702c23/8th-sprint-refactor/02_brokers/images/images_new/topic1_new.png)

![](https://github.com/Yandex-Practicum/python_p2f_ugc/blob/0ef4e078f807861eefc547e2bbd605ab16702c23/8th-sprint-refactor/02_brokers/images/images_new/topicExchange_new.mp4)

### Fanout Exchange

**Fanout** копирует и отправляет полученное сообщение во все очереди, привязанные к нему, независимо от того, какие ключи маршрутизации или шаблоны определены. 
Предоставленные ключи будут просто проигнорированы.

Fanout Exchange может быть полезен, когда одно и то же сообщение необходимо отправить в одну или несколько очередей, с консьюмерами, которые могут обрабатывать одно и то же сообщение по-разному. Идеально подходит для реализации паттерна [Publish/Subscribe(PubSub)](https://ru.wikipedia.org/wiki/%D0%98%D0%B7%D0%B4%D0%B0%D1%82%D0%B5%D0%BB%D1%8C_%E2%80%94_%D0%BF%D0%BE%D0%B4%D0%BF%D0%B8%D1%81%D1%87%D0%B8%D0%BA){target="_blank"}.


```python
# fanOut/producer_fanout.py

from time import sleep

from pika import BlockingConnection, ConnectionParameters, PlainCredentials
from pika.exchange_type import ExchangeType

EXCHANGE_NAME = 'events'


if __name__ == '__main__':
    credentials = PlainCredentials('user', 'password')
    connection = BlockingConnection(ConnectionParameters(host='localhost', credentials=credentials))
    channel = connection.channel()

    channel.exchange_declare(exchange=EXCHANGE_NAME, exchange_type=ExchangeType.fanout)

    counter = 0
    while True:
        counter += 1
        message_body = f'Произошло событие №{counter}'
        channel.basic_publish(exchange=EXCHANGE_NAME, routing_key='', body=message_body)
        print(f'[✅] {message_body}')
        sleep(1)
```

```python
# fanOut/consumer_fanout.py

from pika import BlockingConnection, ConnectionParameters, PlainCredentials
from pika.adapters.blocking_connection import BlockingChannel
from pika.exchange_type import ExchangeType
from pika.spec import BasicProperties, Basic

EXCHANGE_NAME = 'events'


def handler(ch: BlockingChannel, method: Basic.Deliver, properties: BasicProperties, body: bytes) -> None:
    print(f'[🎉] Получено: {body.decode()}')


if __name__ == '__main__':
    credentials = PlainCredentials('user', 'password')
    connection = BlockingConnection(ConnectionParameters(host='localhost', credentials=credentials))
    channel = connection.channel()

    channel.exchange_declare(exchange=EXCHANGE_NAME, exchange_type=ExchangeType.fanout)
    queue = channel.queue_declare(queue='', exclusive=True)

    queue_name = queue.method.queue

    channel.queue_bind(exchange=EXCHANGE_NAME, queue=queue_name, routing_key='')
    channel.basic_consume(queue=queue_name, on_message_callback=handler, auto_ack=True)
    channel.start_consuming()
```

Запустим продюсер и два консьюмера. Консьюмеров может быть и больше, и работать они могут по-разному, но каждый из них получает все сообщения, отправляемые продюсером.

![](https://github.com/Yandex-Practicum/python_p2f_ugc/blob/0ef4e078f807861eefc547e2bbd605ab16702c23/8th-sprint-refactor/02_brokers/images/images_new/fanout1_new.png)

Порядок работы такой:

![](https://github.com/Yandex-Practicum/python_p2f_ugc/blob/0ef4e078f807861eefc547e2bbd605ab16702c23/8th-sprint-refactor/02_brokers/images/images_new/fanoutExchange_new.mp4)


### Headers Exchange

**Headers Exchange** — самый гибкий способ маршрутизации сообщений в RabbitMQ. 
Этот тип распределителя направляет сообщения в очереди в соответствии с совпадением аргументов привязки этих очередей к Exchange с заголовками сообщений. 
Заголовки представляют собой словарь с парами «ключ-значение».

- При использовании этого типа распределителя учитывайте дополнительные накладные расходы на вычисление и тот факт, что он работает медленнее, чем прочие типы Exchange. 
- Все пары (ключ, значение) атрибута headers должны сортироваться по имени ключа перед вычислением значений маршрутизации сообщения. 
- Почти всегда  способу Headers Exchange предпочитают Topic Exchange, однако, знание того, что такой инструмент есть, может выручить, когда не хватает базовых приёмов.


В следующем примере будем публиковать вакансии для разработчиков:

```python
import random
from time import sleep

from pika import BlockingConnection, ConnectionParameters, PlainCredentials
from pika.exchange_type import ExchangeType
from pika.spec import BasicProperties

EXCHANGE_NAME = 'headers'


if __name__ == '__main__':
    credentials = PlainCredentials('user', 'password')
    connection = BlockingConnection(ConnectionParameters(host='localhost', credentials=credentials))
    channel = connection.channel()

    channel.exchange_declare(exchange=EXCHANGE_NAME, exchange_type=ExchangeType.headers)

    vacancies = [
        {'position': 'middle', 'language': 'python'},
        {'position': 'senior', 'language': 'python'},
        {'position': 'junior', 'language': 'C++'},
        {'position': 'middle', 'language': 'java'},
        {'position': 'lead', 'language': 'java'},
        {'position': 'junior', 'language': 'C#'},
        {'position': 'senior', 'language': 'C#'},
    ]

    counter = 0
    while True:
        counter += 1
        vacancy = random.choice(vacancies)
        message_body = f'Опубликована вакансия №{counter} {vacancy}'
        channel.basic_publish(
            exchange=EXCHANGE_NAME,
            routing_key='',
            body=message_body,
            properties=BasicProperties(headers=vacancy),
        )
        print(f'[✅] {message_body}')
        sleep(1)
```

А в консьюмере обрабатываем только лидов и всех питонистов, вне зависимости от опыта:

```python
from pika import BlockingConnection, ConnectionParameters, PlainCredentials
from pika.adapters.blocking_connection import BlockingChannel
from pika.exchange_type import ExchangeType
from pika.spec import BasicProperties, Basic

EXCHANGE_NAME = 'headers'


def handler(ch: BlockingChannel, method: Basic.Deliver, properties: BasicProperties, body: bytes) -> None:
    print(f'[🎉] Получено: {body.decode()}')


if __name__ == '__main__':
    credentials = PlainCredentials('user', 'password')
    connection = BlockingConnection(ConnectionParameters(host='localhost', credentials=credentials))
    channel = connection.channel()

    channel.exchange_declare(exchange=EXCHANGE_NAME, exchange_type=ExchangeType.headers)
    queue = channel.queue_declare(queue='', exclusive=True)

    queue_name = queue.method.queue

    channel.queue_bind(
        exchange=EXCHANGE_NAME,
        queue=queue_name,
        routing_key='',
        arguments={'x-match': 'any', 'position': 'lead', 'language': 'python'},
    )
    channel.basic_consume(queue=queue_name, on_message_callback=handler, auto_ack=True)
    channel.start_consuming()
```

Запустим продюсера и консьюмера. Результат будет примерно следующий:

![](https://github.com/Yandex-Practicum/python_p2f_ugc/blob/0ef4e078f807861eefc547e2bbd605ab16702c23/8th-sprint-refactor/02_brokers/images/images_new/headers1_new.png)

Более подробно с этим типом Exchange можно ознакомиться на странице официальной документации [https://www.rabbitmq.com/tutorials/amqp-concepts#exchange-headers](https://www.rabbitmq.com/tutorials/amqp-concepts#exchange-headers){target="_blank"}.


## RPC (Remote procedure call) 

Рассмотренные выше способы работы с RabbitMQ позволяют реализовывать многие виды задач: от распределённых вычислений и отложенных заданий до реализации паттерна публикатор/подписчик.

Обсудим ещё один тип использования брокера сообщений: вызов удалённых процедур ([Remote procedure call](https://ru.wikipedia.org/wiki/%D0%A3%D0%B4%D0%B0%D0%BB%D1%91%D0%BD%D0%BD%D1%8B%D0%B9_%D0%B2%D1%8B%D0%B7%D0%BE%D0%B2_%D0%BF%D1%80%D0%BE%D1%86%D0%B5%D0%B4%D1%83%D1%80){target="_blank"}).

Сразу окунёмся в код: напишем банковское приложение, которое обновляет баланс пользователя по запросу клиентского приложения.

Серверная часть:

```python
# rpc/rpc_server.py

import json

import pika
from pika.adapters.blocking_connection import BlockingChannel, BlockingConnection
from pika.connection import ConnectionParameters
from pika.credentials import PlainCredentials
from pika.spec import BasicProperties, Basic

credentials = PlainCredentials('user', 'password')
connection = BlockingConnection(ConnectionParameters(host='localhost', credentials=credentials))

channel = connection.channel()

channel.queue_declare(queue='rpc_queue')

USERS = {
    0: {'username': 'Александр Стравинский', 'balance': 0},
    1: {'username': 'Степан Лиходеев', 'balance': 0},
    2: {'username': 'Тимофей Квасцов', 'balance': 0},
}


def update_balance_for_user(user_id: int, amount: int) -> dict:
    USERS[user_id]['balance'] += amount
    return USERS[user_id]


def handler(ch: BlockingChannel, method: Basic.Deliver, properties: BasicProperties, body: bytes):
    message = json.loads(body)
    print(f'[₽] Увеличение баланса пользователя {message["id"]} на {message["amount"]}')
    updated_user = update_balance_for_user(user_id=message['id'], amount=message['amount'])

    ch.basic_publish(
        exchange='',
        routing_key=properties.reply_to,
        properties=pika.BasicProperties(correlation_id=properties.correlation_id),
        body=json.dumps(updated_user),
    )
    ch.basic_ack(delivery_tag=method.delivery_tag)


if __name__ == '__main__':
    channel.basic_qos(prefetch_count=1)
    channel.basic_consume(queue='rpc_queue', on_message_callback=handler)

    print('[🎉] Ожидание RPC запроса')
    channel.start_consuming()
```

Код сервера довольно прост:

- Как обычно, начинаем с установления соединения и объявления очереди `rpc_queue`.
- Объявляем нашу функцию `update_balance_for_user`, которая обновляет баланс пользователя в базе данных.
- Объявляем обратный вызов `handler` — это ядро нашего RPC-сервера. Он выполняется при получении запроса, производит вычисления и отправляет ответ обратно.
- Возможно, нам захочется запустить более одного серверного процесса. Чтобы равномерно распределить нагрузку на несколько серверов, нужно установить параметр `prefetch_count`. При запуске нескольких серверов, они стартуют не всегда одновременно, а RabbitMQ по умолчанию отправляет сообщения из очереди первому же доступному консьюмеру. При использовании опции `prefetch_count=n` [в примере n=1, но может быть и 2, 3, 4, ...] консьюмер не получает следующие `n` сообщений, пока не подтвердит предыдущие.

```python
# rpc/rpc_client.py

import json
import random
import uuid
from time import sleep

from pika.adapters.blocking_connection import BlockingConnection, BlockingChannel
from pika.connection import ConnectionParameters
from pika.credentials import PlainCredentials
from pika.spec import Basic, BasicProperties


class BankRpcClient:
    def __init__(self):
        credentials = PlainCredentials('user', 'password')
        self.connection = BlockingConnection(ConnectionParameters(host='localhost', credentials=credentials))

        self.channel = self.connection.channel()

        result = self.channel.queue_declare(queue='', exclusive=True)
        self.callback_queue = result.method.queue

        self.channel.basic_consume(queue=self.callback_queue, on_message_callback=self.on_response, auto_ack=True)

        self.response = None
        self.corr_id = None

    def on_response(self, ch: BlockingChannel, method: Basic.Deliver, properties: BasicProperties, body: bytes) -> None:
        if self.corr_id == properties.correlation_id:
            self.response = body

    def process(self, data: dict) -> str:
        self.response = None
        self.corr_id = str(uuid.uuid4())
        self.channel.basic_publish(
            exchange='',
            routing_key='rpc_queue',
            properties=BasicProperties(
                reply_to=self.callback_queue,
                correlation_id=self.corr_id,
            ),
            body=json.dumps(data),
        )
        while self.response is None:
            self.connection.process_data_events(time_limit=0)
        return json.loads(self.response)


if __name__ == '__main__':
    rpc_client = BankRpcClient()

    while True:
        user_id = random.randint(0, 2)
        amount = random.randint(1, 10)
        print('-' * 20)
        print(f'[➡️] Увеличение баланса пользователя {user_id=} на {amount}')
        response = rpc_client.process({'id': user_id, 'amount': amount})
        print(f'[✅] Пользователь: {response}')
        sleep(1)
```

Код клиентского приложение немного более сложный:
- Устанавливаем соединение, канал и объявляем [эксклюзивную](https://www.rabbitmq.com/docs/queues#exclusive-queues){target="_blank"} очередь `callback_queue` для ответов.
- Подписываемся на `callback_queue`, чтобы получать ответы RPC.
- Обратный вызов `on_response`, который выполняется при каждом ответе, выполняет очень простую работу: для каждого ответного сообщения он проверяет, является ли `correlation_id` тем, что мы ищем. Если да, он сохраняет ответ в `self.response` и разрывает цикл получения сообщений.
- Определяем наш основной метод `process` — он выполняет фактический запрос RPC.
- В методе `process` генерируем уникальный номер `correlation_id` и сохраняем его, функция обратного вызова `on_response` будет использовать это значение для перехвата соответствующего ответа и однозначного сопоставления ответа с запросом.
- Также в методе `process` мы публикуем сообщение запроса с двумя свойствами: `reply_to` и `correlation_id`.
- В конце ждём ответ и возвращаем его пользователю.


```shell
❯ python rpc/rpc_server.py
[🎉] Ожидание RPC запроса
[₽] Увеличение баланса пользователя 0 на 10
[₽] Увеличение баланса пользователя 2 на 9
[₽] Увеличение баланса пользователя 1 на 2
[₽] Увеличение баланса пользователя 0 на 5
[₽] Увеличение баланса пользователя 1 на 5
[₽] Увеличение баланса пользователя 0 на 6
[₽] Увеличение баланса пользователя 2 на 3
[₽] Увеличение баланса пользователя 0 на 4
```

```shell
❯ python rpc/rpc_client.py
--------------------
[➡️] Увеличение баланса пользователя user_id=0 на 10
[✅] Пользователь: {'username': 'Александр Стравинский', 'balance': 10}
--------------------
[➡️] Увеличение баланса пользователя user_id=2 на 9
[✅] Пользователь: {'username': 'Тимофей Квасцов', 'balance': 9}
--------------------
[➡️] Увеличение баланса пользователя user_id=1 на 2
[✅] Пользователь: {'username': 'Степан Лиходеев', 'balance': 2}
--------------------
[➡️] Увеличение баланса пользователя user_id=0 на 5
[✅] Пользователь: {'username': 'Александр Стравинский', 'balance': 15}
--------------------
[➡️] Увеличение баланса пользователя user_id=1 на 5
[✅] Пользователь: {'username': 'Степан Лиходеев', 'balance': 7}
--------------------
[➡️] Увеличение баланса пользователя user_id=0 на 6
[✅] Пользователь: {'username': 'Александр Стравинский', 'balance': 21}
--------------------
[➡️] Увеличение баланса пользователя user_id=2 на 3
[✅] Пользователь: {'username': 'Тимофей Квасцов', 'balance': 12}
--------------------
[➡️] Увеличение баланса пользователя user_id=0 на 4
[✅] Пользователь: {'username': 'Александр Стравинский', 'balance': 25}
```

Все работает, как и ожидалось 🎉
