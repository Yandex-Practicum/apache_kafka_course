import logging.config

from confluent_kafka import Consumer

from models.message import Message
from settings.logging import LOGGING_CONFIG

logging.config.dictConfig(LOGGING_CONFIG)
logger = logging.getLogger("pull_consumer_app")


def pull_consume():
    """
    This consumer manually polls Kafka for new messages, implementing a pull model.
    """
    consumer_conf = {
        "bootstrap.servers": "localhost:29092,localhost:39092,localhost:49092",
        "group.id": "consumer-group-1",
        "auto.offset.reset": "earliest",
        "enable.auto.commit": False,          # Отключает автоматическое сохранение смещения
        "fetch.min.bytes": 10 * 1024 * 1024,  # Для большей наглядности ограничен минимальный размер пакета сообщений
        "fetch.wait.max.ms": 10_000,
    }
    consumer = Consumer(consumer_conf)
    consumer.subscribe(["test-topic"])

    try:
        while True:
            records = consumer.consume(timeout=10)

            for record in records:
                message = Message.deserialize(record.value().decode("utf-8"))
                logger.info(f"Pulled Message: {message.content}")
            consumer.commit()
    except Exception as e:
        logger.exception(f"Failure: {e}")
    finally:
        consumer.close()


if __name__ == "__main__":
    pull_consume()
