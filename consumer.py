import logging
from confluent_kafka import Consumer


logger = logging.getLogger(__name__)

if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
    params = {
        'bootstrap.servers': 'rc1a-tvhv0qrfoh26dp4h.mdb.yandexcloud.net:9091',
        'security.protocol': 'SASL_SSL',
        'ssl.ca.location': './secrets/YandexInternalRootCA.crt',
        'sasl.mechanism': 'SCRAM-SHA-512',
        'sasl.username': 'consumer',
        'sasl.password': 'consumer',
        'group.id': 'test-consumer-group',
        'auto.offset.reset': 'earliest',
        'enable.auto.commit': False,
    }
    consumer = Consumer(params)
    consumer.subscribe(['messages'])
    try:
        while True:
            message = consumer.poll(0.1)

            if message is None:
                continue
            if message.error():
                logger.error(f"Ошибка топика {message.topic()}: {message.error()}")
                continue

            key = message.key()#.decode("utf-8")
            value = message.value()#.decode("utf-8")
            logger.info(f"Получено сообщение из топика {message.topic()}: {key=}, {value=}, offset={message.offset()}")
    finally:
        consumer.close()
