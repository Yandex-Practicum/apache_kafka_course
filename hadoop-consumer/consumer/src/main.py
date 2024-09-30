import asyncio
from contextlib import asynccontextmanager

import uvicorn
from fastapi import FastAPI
from aiokafka import AIOKafkaConsumer
from aiokafka.helpers import create_ssl_context
from confluent_kafka.schema_registry import SchemaRegistryClient
from apscheduler.schedulers.asyncio import AsyncIOScheduler
from apscheduler.triggers.interval import IntervalTrigger
from hdfs import InsecureClient

from src.consumer import routers as kafka_routers
from src.consumer.service import consume_batch, consume_single
from src.config import KafkaConfig

config = KafkaConfig()


async def start_consumer():
    """Запуск Kafka Consumer"""

    context = create_ssl_context(
        cafile=config.KAFKA_CA_PATH
    )

    kafka_consumer = AIOKafkaConsumer(
            config.KAFKA_TOPIC,
            bootstrap_servers=config.KAFKA_BOOTSTRAP_SERVERS,
            group_id=config.GROUP_ID,
            enable_auto_commit=config.AUTO_COMMIT,
            max_poll_records=config.MESSAGES_BATCH_MAX,
            security_protocol='SASL_SSL',
            sasl_mechanism='PLAIN',
            sasl_plain_username=config.KAFKA_USER,
            sasl_plain_password=config.KAFKA_PASSWORD,
            ssl_context=context
        )
    await kafka_consumer.start()
    return kafka_consumer


async def shutdown_consumer(kafka_consumer: AIOKafkaConsumer):
    """Остановка Kafka Consumer"""
    await kafka_consumer.stop()


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Добавление kafka consumer в объект FastApi, чтобы не запускать consumer при каждом запросе"""
    app.kafka_consumer = await start_consumer()
    app.scheduler = AsyncIOScheduler()
    app.scheduler.add_job(
        consume_batch if config.MESSAGES_BATCH_MIN > 1 else consume_single, 
        args=(app, config),
        trigger=IntervalTrigger(seconds=config.POLL_INTERVAL_SECONDS),
        id='poll_kafka',
        replace_existing=False
    )
    app.scheduler.start()
    app.schema_client = SchemaRegistryClient({
        'url': config.SCHEMA_REGISTRY_SERVER, 
        'ssl.ca.location': '/ca.crt'})
    app.include_router(kafka_routers.all_routers)
    app.hadoop_client = InsecureClient(config.HADOOP_ENDPOINT, user="root")
    yield
    await shutdown_consumer(app.kafka_consumer)
    app.scheduler.shutdown()

app = FastAPI(
    title="Kafka consumer",
    description="""
            Получает сообщения из Kafka
            """,
    version="1.0.0",
    lifespan=lifespan
)

async def start_fastapi():
    """Запуск FastAPI сервера"""
    config = uvicorn.Config(
        app=app,
        host="0.0.0.0",
        port=9001,
    )
    server = uvicorn.Server(config)
    await server.serve()


async def main():
    """Основная точка входа - запуск веб-сервера Uvicorn
    """
    await asyncio.create_task(
        start_fastapi(),
    )


if __name__ == "__main__":
    asyncio.run(main())
