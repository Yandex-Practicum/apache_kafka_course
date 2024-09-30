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
