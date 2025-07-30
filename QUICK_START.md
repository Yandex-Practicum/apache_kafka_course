# Быстрый старт: Kafka to HDFS

## 🚀 Быстрый запуск

### 1. Запуск инфраструктуры
```bash
# Запуск Kafka и HDFS
docker-compose up -d

# Ждем запуска (30-60 секунд)
sleep 60
```

### 2. Установка зависимостей
```bash
pip install -r requirements.txt
```

### 3. Проверка системы
```bash
python health_check.py
```

### 4. Демонстрация (полный цикл)
```bash
python run_demo.py
```

## 📋 Ручной запуск

### Отправка тестовых сообщений
```bash
python kafka_producer.py --count 10
```

### Запуск процессора
```bash
python kafka_to_hdfs.py
```

## 🔍 Проверка результатов

### Просмотр файлов в HDFS
```bash
docker exec -it hadoop-namenode hdfs dfs -ls -R /kafka_data
```

### Web интерфейсы
- HDFS: http://localhost:9870
- Kafka (если настроен): http://localhost:9094

## 🛠️ Устранение проблем

### Если Kafka не отвечает
```bash
docker logs kafka-0
```

### Если HDFS не отвечает
```bash
docker logs hadoop-namenode
```

### Проверка портов
```bash
netstat -an | grep -E "(9094|9870)"
```

## 📁 Структура файлов

```
├── kafka_to_hdfs.py      # Основной процессор
├── kafka_producer.py     # Тестовый producer
├── health_check.py       # Проверка здоровья
├── run_demo.py          # Демонстрация
├── config.yaml          # Конфигурация
├── requirements.txt     # Зависимости
└── *.log               # Логи
``` 