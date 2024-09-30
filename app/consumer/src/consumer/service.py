import json

from loguru import logger
from fastapi import FastAPI
from aiokafka.structs import ConsumerRecord
from confluent_kafka.schema_registry.avro import AvroDeserializer, SerializationContext

from src.config import KafkaConfig
from src.consumer.schemas import Message


async def process_message(app: FastAPI, config: KafkaConfig, msg: ConsumerRecord):
    deserializer = AvroDeserializer(app.schema_client, json.dumps(Message.avro_schema()))
    context = SerializationContext(config.KAFKA_TOPIC, "message")

    logger.info(json.dumps(dict(
        topic=msg.topic,
        partition=msg.partition,
        offset=msg.offset,
        key=msg.key,
        message=deserializer(msg.value, context),
        timestamp=msg.timestamp
    )))


async def consume_single(app: FastAPI, config: KafkaConfig):
    result = await app.kafka_consumer.getone()

    await process_message(app, config, result)

    if not config.AUTO_COMMIT:
        # Коммитим только для этой партиции
        await app.kafka_consumer.commit({result.partition: result.offset + 1})


async def consume_batch(app: FastAPI, config: KafkaConfig):
    # Получаем список партиций, которые нам назначены
    partitions = app.kafka_consumer.assignment()
    offsets = await app.kafka_consumer.end_offsets(partitions)

    # Для каждой партиции проверяем по последнему офсету и закомиченному офсету, набралось ли уже n записей (по лагу)
    for part in partitions:
        cur_offset = offsets[part]
        cur_committed = await app.kafka_consumer.committed(part) or 0

        logger.debug(f"Partition {part} offset {cur_offset}, commited {cur_committed}, lag {cur_offset - cur_committed}")

        if cur_offset - cur_committed >= config.MESSAGES_BATCH_MIN:
            result = await app.kafka_consumer.getmany(part, timeout_ms=config.POLL_INTERVAL_SECONDS - 1 * 1000)

            for tp, messages in result.items():
                if messages:
                    logger.debug(f"Polled {len(messages)} messages from {part}")
                    for msg in messages:
                        await process_message(app, config, msg)

                    if not config.AUTO_COMMIT:
                        # Коммитим только для этой партиции
                        await app.kafka_consumer.commit({tp: messages[-1].offset + 1})
