import logging
from confluent_kafka import Producer


logger = logging.getLogger(__name__)

def delivery_report(err, msg):
    """
    Проверить статус доставки сообщения
    """
    if err is not None:
        logger.info(f"Ошибка доставки сообщения в топик {msg.topic()}: {err}")
    else:
        logger.info(f"Сообщение «{msg.__str__()}» доставлено в {msg.topic()} [{msg.partition()}]")

if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')

    params = {
        'bootstrap.servers': 'rc1a-tvhv0qrfoh26dp4h.mdb.yandexcloud.net:9091',
        'security.protocol': 'SASL_SSL',
        'ssl.ca.location': './secrets/YandexInternalRootCA.crt',
        'sasl.mechanism': 'SCRAM-SHA-512',
        'sasl.username': 'producer',
        'sasl.password': 'producer',
        'error_cb': delivery_report,
    }

    producer = Producer(params)
    producer.produce(
        topic='messages',
        key='msg', 
        value='test message',
        callback=delivery_report 
    )
    producer.flush(10)