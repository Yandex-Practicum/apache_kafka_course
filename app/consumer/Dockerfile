FROM python:3.13.3-slim-bookworm

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV PYTHONPATH=/app

WORKDIR /app
COPY . .

RUN pip install poetry && \
    poetry config virtualenvs.create false && \
    poetry install --no-root

HEALTHCHECK --interval=5s --timeout=10s --retries=3 CMD curl -sS http://localhost:9001/api/v1/health/ping | grep -q pong || exit 1

EXPOSE 9001 

CMD ["python3", "/app/src/main.py"]