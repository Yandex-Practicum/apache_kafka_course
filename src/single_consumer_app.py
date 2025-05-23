import logging.config

from confluent_kafka import Consumer

from models.message import Message
from settings.logging import LOGGING_CONFIG

logging.config.dictConfig(LOGGING_CONFIG)
logger = logging.getLogger("single_consumer_app")


def single_consume():
    """
    This consumer reacts to messages immediately using a listener.
    """
    consumer_conf = {
        "bootstrap.servers": "localhost:29092,localhost:39092,localhost:49092",
        "group.id": "consumer-group-2",
        "auto.offset.reset": "earliest",
        "enable.auto.commit": True,
        "fetch.wait.max.ms": 1,  # Длительность ожидания при получении пакета сообщений
        "fetch.max.bytes": 1000,
        "message.max.bytes": 1000
    }
    consumer = Consumer(consumer_conf)
    consumer.subscribe(["test-topic"])

    try:
        while True:
            records = consumer.consume(timeout=1000)
            logger.info(f"Pull records: {len(records)}")

            for record in records:
                message = Message.deserialize(record.value().decode("utf-8"))
                logger.info(f"singleed Message: {message.content}")
    except Exception as e:
        logger.exception(f"Failure: {e}")
    finally:
        consumer.close()


if __name__ == "__main__":
    single_consume()
