from fastapi import APIRouter

from consumer.endpoints.health import router as health_router

all_routers = APIRouter(prefix="/api/v1")

all_routers.include_router(
    health_router,
    prefix="/health",
    tags=["health"],
)
