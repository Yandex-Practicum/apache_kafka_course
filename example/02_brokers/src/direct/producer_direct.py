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
