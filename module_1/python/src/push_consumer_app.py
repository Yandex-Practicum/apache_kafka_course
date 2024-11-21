import logging.config

from confluent_kafka import Consumer

from models.message import Message
from settings.logging import LOGGING_CONFIG

logging.config.dictConfig(LOGGING_CONFIG)
logger = logging.getLogger("push_consumer_app")


def push_consume():
    """
    This consumer reacts to messages immediately using a listener.
    """
    consumer_conf = {
        "bootstrap.servers": "localhost:29092,localhost:39092,localhost:49092",
        "group.id": "consumer-group-2",
        "auto.offset.reset": "earliest",
        "enable.auto.commit": True,
    }
    consumer = Consumer(consumer_conf)
    consumer.subscribe(["test-topic"])

    try:
        while True:
            records = consumer.consume(timeout=0.0001)

            for record in records:
                message = Message.deserialize(record.value().decode("utf-8"))
                logger.info(f"Pushed Message: {message.content}")
    except Exception as e:
        logger.exception(f"Failure: {e}")
    finally:
        consumer.close()


if __name__ == "__main__":
    push_consume()
