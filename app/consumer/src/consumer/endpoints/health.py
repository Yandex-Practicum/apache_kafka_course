from fastapi import APIRouter
from fastapi.responses import JSONResponse

router = APIRouter()


@router.get(
    "/ping",
    response_class=JSONResponse,
    summary="Healthcheck",
    description="Выполняет функцию проверки доступности сервиса (healthcheck)",
)

async def ping():
    return JSONResponse(
        content={"message": "pong"},
        status_code=200,
    )
