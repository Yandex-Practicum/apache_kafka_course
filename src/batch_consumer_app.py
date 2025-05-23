import logging.config

from confluent_kafka import Consumer

from models.message import Message
from settings.logging import LOGGING_CONFIG

logging.config.dictConfig(LOGGING_CONFIG)
logger = logging.getLogger("batch_consumer_app")


def batch_consume():
    """
    This consumer manually polls Kafka for new messages, implementing a batch model.
    """
    consumer_conf = {
        "bootstrap.servers": "localhost:29092,localhost:39092,localhost:49092",
        "group.id": "consumer-group-1",
        "auto.offset.reset": "earliest",
        "enable.auto.commit": False,          # Отключает автоматическое сохранение смещения
        "fetch.min.bytes": 10 * 60,  # примерный размер пакета сообщений
        "fetch.wait.max.ms": 300_000,
    }
    consumer = Consumer(consumer_conf)
    consumer.subscribe(["test-topic"])

    try:
        while True:
            records = consumer.consume(timeout=300_000, num_messages=10)
            logger.info(f"batch records: {len(records)}")

            for record in records:
                message = Message.deserialize(record.value().decode("utf-8"))
                logger.info(f"batched Message: {message.content}")
            consumer.commit()
    except Exception as e:
        logger.exception(f"Failure: {e}")
        consumer.commit()
    finally:
        consumer.close()


if __name__ == "__main__":
    batch_consume()
