from fastapi import APIRouter

from producer.endpoints.health import router as health_router
from producer.endpoints.kafka import router as kafka_router

all_routers = APIRouter(prefix="/api/v1")

all_routers.include_router(
    kafka_router,
    prefix="/kafka",
    tags=["kafka"],
)
all_routers.include_router(
    health_router,
    prefix="/health",
    tags=["health"],
)
