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
