import asyncio
from contextlib import asynccontextmanager
import json
import logging

from confluent_kafka import Consumer
from fastapi import FastAPI
from jinja2 import Template
import uvicorn
from starlette.responses import PlainTextResponse

TOPIC = 'metrics'
TIMEOUT = 30 * 60
metrics = {}
lock = asyncio.Lock()

logger = logging.getLogger(__name__)
logging.basicConfig(level=logging.INFO)
logger.info('filtered messages consumer started')

conf = {
    "group.id": 'filtered',
    "bootstrap.servers": "kafka-0:9092",
    "auto.offset.reset": "earliest",
    "enable.auto.commit": "true",
    "fetch.min.bytes": 1024,
}

consumer = Consumer(conf)
consumer.subscribe([TOPIC])


def gen_str(metrics: dict) -> str:
    strs = []
    for key, value in metrics.items():
        schema = f"# HELP {value['Name']} {value['Description']}\n# TYPE {value['Name']} {value['Type']}\n{value['Name']} {value['Value']}"
        strs.append(schema)
    return '\n'.join(strs)


async def read_metrics(consumer: Consumer) -> None:
    try:
        while True:
            msg = consumer.poll(timeout=TIMEOUT)

            if msg is None:
                continue
            if msg.error():
                logger.error(f'получена ошибка: {msg.error()}')
                continue

            logger.info(f'полученое сообщение: {msg.value().decode("utf-8")}')

            async with lock:
                global metrics
                metrics = msg.value().decode("utf-8")

            break

    except Exception as e:
        logger.error(f'поймана ошибка: {e}')

    #finally:
    #    consumer.close()


async def run_scheduler(consumer: Consumer):
    while True:
        await read_metrics(consumer)
        await asyncio.sleep(TIMEOUT + 5)


@asynccontextmanager
async def lifespan(app: FastAPI):
    asyncio.create_task(
        run_scheduler(consumer)
    )

    yield


app = FastAPI(lifespan=lifespan)


@app.get("/metrics", response_class=PlainTextResponse)
async def metrics():
    metrics_dict = json.loads(metrics)
    return gen_str(metrics_dict)

if __name__ == '__main__':
    uvicorn.run(app, host='0.0.0.0', port=8000)
