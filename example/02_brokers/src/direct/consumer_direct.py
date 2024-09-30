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
    # обратите внимание, выше мы не указали имя очереди, в этом случае имя сгенерируется автоматически
    # и будет иметь вид amq.gen-fxDXeobOVBDZTpWfVM6iDA
    queue_name = queue.method.queue

    # получим нужные ключи маршрутизации из аргументов запуска приложения
    routing_keys = sys.argv[1:] or ['error', 'warning', 'info']
    for key in routing_keys:
        channel.queue_bind(exchange=EXCHANGE_NAME, queue=queue_name, routing_key=key)

    channel.basic_consume(queue=queue_name, on_message_callback=handler, auto_ack=True)
    channel.start_consuming()
