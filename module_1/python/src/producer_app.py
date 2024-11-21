import logging.config
import random
import time

from confluent_kafka import Producer

from models.message import Message
from settings.logging import LOGGING_CONFIG

logging.config.dictConfig(LOGGING_CONFIG)
logger = logging.getLogger("producer_app")


def delivery_report(err, msg):
    if err:
        logger.error(f"Error while producing: {err}")
    else:
        logger.info(
            f"Message sent successfully: topic='{msg.topic()}', partition={msg.partition()}, offset={msg.offset()}"
        )


def produce():
    """
    This producer sends messages to Kafka with "at-least-once" delivery guarantee.
    """
    producer_conf = {
        "bootstrap.servers": "localhost:29092,localhost:39092,localhost:49092",
        # Настройка подтверждения доставки:
        "acks": "all",  # Гарантия «at-least-once»
        "retries": 3,   # Поведение при сбоях
    }
    producer = Producer(producer_conf)

    try:
        for i in range(1, 11):
            logger.info(f"Preparing to send message number {i}")
            message = Message(f"user123{random.randint(1, 100)}", "Hello Kafka!")

            try:
                serialized_message = message.serialize()
                producer.produce(
                    topic="test-topic",
                    key=message.user_id.encode("utf-8"),
                    value=serialized_message.encode("utf-8"),
                    callback=delivery_report,
                )
                producer.poll(0)
                time.sleep(0.6)
            except Exception as e:
                logger.error(f"Error while producing: {e}")
    finally:
        producer.flush()


if __name__ == "__main__":
    produce()
