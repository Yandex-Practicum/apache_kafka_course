## Локальная разработка

`poetry install`
`export PYTHONPATH=.`
`poetry run python3 src/main.py`

## Реализация

Продьюсер ожидает POST-запросы по пути `/api/v1/kafka/send` в формате json. Отправить запрос можно следующими способами:
* `curl -vv -XPOST -d '{"text": "message"}' -H "Content-Type: application/json"  localhost:9000/api/v1/kafka/send`
* Через swagger ui по адресу `http://localhost:9000/docs#/kafka/send_kafka_api_v1_kafka_send_post`

При успешной обработке сообщения приложение возращает код ответа 200. Отправляемое сообщение выводится в консоль

Реализацию отправки сообщения в kafka можно посмотреть [тут](src/producer/endpoints/kafka.py)

Для гарантии "At least once" включён параметр "acks=all" и включена проверка идемпотентности при отправке ("enable_idempotence=True"). Aiokafka автоматически делает retry доставки при ошибках в течение request_timeout_ms.