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
