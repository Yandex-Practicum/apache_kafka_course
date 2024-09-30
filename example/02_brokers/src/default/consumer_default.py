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
