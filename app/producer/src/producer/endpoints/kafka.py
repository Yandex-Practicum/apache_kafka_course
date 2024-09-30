import asyncio
import json

from fastapi import APIRouter
from fastapi import HTTPException
from fastapi import Request
from fastapi.responses import JSONResponse
from aiokafka import AIOKafkaProducer
from confluent_kafka.schema_registry.avro import AvroSerializer, SerializationContext
from loguru import logger

from src.producer.schemas import Message
from src.main import config

event_loop = asyncio.get_event_loop()

router = APIRouter()

@router.post(
    "/send",
    description="Send message to kafka",
)
async def send_kafka(message: Message, request: Request) -> None:
    
    serializer = AvroSerializer(request.app.schema_client, json.dumps(message.avro_schema()))
    logger.info(json.dumps(message.model_dump_json()))

    await request.app.kafka_producer.send(value=serializer(message.model_dump(),
                                                           SerializationContext(config.KAFKA_TOPIC, "message")), 
                                                           topic=config.KAFKA_TOPIC)