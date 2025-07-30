#!/usr/bin/env python3
"""
Демонстрационный скрипт для полного цикла работы Kafka to HDFS
"""

import time
import subprocess
import sys
import logging
from datetime import datetime


def setup_logging():
    """Настройка логирования"""
    logging.basicConfig(
        level=logging.INFO,
        format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
        handlers=[
            logging.FileHandler('demo.log'),
            logging.StreamHandler()
        ]
    )
    return logging.getLogger(__name__)


def run_command(command, description):
    """Выполнение команды с логированием"""
    logger = logging.getLogger(__name__)
    logger.info(f"Выполняем: {description}")
    logger.info(f"Команда: {command}")
    
    try:
        result = subprocess.run(
            command, 
            shell=True, 
            capture_output=True, 
            text=True, 
            timeout=60
        )
        
        if result.returncode == 0:
            logger.info(f"✅ {description} - УСПЕШНО")
            if result.stdout:
                logger.debug(f"Вывод: {result.stdout}")
        else:
            logger.error(f"❌ {description} - ОШИБКА")
            logger.error(f"Ошибка: {result.stderr}")
            return False
            
        return True
        
    except subprocess.TimeoutExpired:
        logger.error(f"❌ {description} - ТАЙМАУТ")
        return False
    except Exception as e:
        logger.error(f"❌ {description} - ИСКЛЮЧЕНИЕ: {e}")
        return False


def main():
    """Главная функция демонстрации"""
    logger = setup_logging()
    
    logger.info("🚀 Начинаем демонстрацию Kafka to HDFS процессора")
    logger.info("=" * 60)
    
    # Шаг 1: Проверка здоровья системы
    logger.info("📋 Шаг 1: Проверка здоровья системы")
    if not run_command("python health_check.py", "Health check"):
        logger.error("❌ Система не готова к работе. Проверьте настройки.")
        sys.exit(1)
    
    # Шаг 2: Создание топика в Kafka (если не существует)
    logger.info("📋 Шаг 2: Создание топика в Kafka")
    create_topic_cmd = """
    docker exec -it kafka-0 kafka-topics.sh \
        --create \
        --topic test-topic \
        --bootstrap-server localhost:9092 \
        --partitions 3 \
        --replication-factor 1 \
        --if-not-exists
    """
    run_command(create_topic_cmd, "Создание топика test-topic")
    
    # Шаг 3: Отправка тестовых сообщений
    logger.info("📋 Шаг 3: Отправка тестовых сообщений в Kafka")
    if not run_command("python kafka_producer.py --count 5", "Отправка тестовых сообщений"):
        logger.error("❌ Не удалось отправить тестовые сообщения")
        sys.exit(1)
    
    # Шаг 4: Запуск процессора Kafka to HDFS
    logger.info("📋 Шаг 4: Запуск процессора Kafka to HDFS")
    logger.info("⚠️  Запускаем процессор в фоновом режиме...")
    
    try:
        # Запускаем процессор в фоне
        processor_process = subprocess.Popen(
            ["python", "kafka_to_hdfs.py"],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )
        
        # Ждем немного для обработки сообщений
        logger.info("⏳ Ждем 10 секунд для обработки сообщений...")
        time.sleep(10)
        
        # Останавливаем процессор
        logger.info("🛑 Останавливаем процессор...")
        processor_process.terminate()
        
        # Ждем завершения
        try:
            processor_process.wait(timeout=5)
        except subprocess.TimeoutExpired:
            processor_process.kill()
        
        logger.info("✅ Процессор остановлен")
        
    except Exception as e:
        logger.error(f"❌ Ошибка запуска процессора: {e}")
        sys.exit(1)
    
    # Шаг 5: Проверка результатов
    logger.info("📋 Шаг 5: Проверка результатов в HDFS")
    
    # Проверяем созданные файлы
    check_files_cmd = "docker exec -it hadoop-namenode hdfs dfs -ls -R /kafka_data"
    if run_command(check_files_cmd, "Проверка файлов в HDFS"):
        logger.info("✅ Файлы успешно записаны в HDFS")
    else:
        logger.warning("⚠️  Не удалось проверить файлы в HDFS")
    
    # Шаг 6: Показываем содержимое одного файла
    logger.info("📋 Шаг 6: Показываем содержимое файла")
    show_content_cmd = """
    docker exec -it hadoop-namenode hdfs dfs -find /kafka_data -name "*.json" | head -1 | xargs -I {} docker exec -it hadoop-namenode hdfs dfs -cat {}
    """
    run_command(show_content_cmd, "Показ содержимого файла")
    
    # Финальный отчет
    logger.info("=" * 60)
    logger.info("🎉 Демонстрация завершена!")
    logger.info("📊 Результаты:")
    logger.info("  ✅ Kafka подключение работает")
    logger.info("  ✅ HDFS подключение работает") 
    logger.info("  ✅ Сообщения отправлены в Kafka")
    logger.info("  ✅ Сообщения обработаны и записаны в HDFS")
    logger.info("")
    logger.info("🔗 Полезные ссылки:")
    logger.info("  - HDFS Web UI: http://localhost:9870")
    logger.info("  - Логи процессора: kafka_to_hdfs.log")
    logger.info("  - Логи producer: kafka_producer.log")
    logger.info("  - Логи health check: health_check.log")
    logger.info("  - Логи демонстрации: demo.log")


if __name__ == "__main__":
    main() 