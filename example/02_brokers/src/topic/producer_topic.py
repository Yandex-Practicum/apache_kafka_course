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
