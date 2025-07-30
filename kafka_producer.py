#!/usr/bin/env python3
"""
Скрипт для отправки тестовых сообщений в Kafka
"""

import json
import time
import logging
from datetime import datetime
from typing import Dict, Any

from kafka import KafkaProducer
import yaml


class KafkaTestProducer:
    """Класс для отправки тестовых сообщений в Kafka"""
    
    def __init__(self, config_path: str = "config.yaml"):
        """
        Инициализация producer
        
        Args:
            config_path: Путь к файлу конфигурации
        """
        self.config = self._load_config(config_path)
        self.setup_logging()
        self.producer = None
        
    def _load_config(self, config_path: str) -> Dict[str, Any]:
        """Загрузка конфигурации из YAML файла"""
        try:
            with open(config_path, 'r', encoding='utf-8') as file:
                return yaml.safe_load(file)
        except FileNotFoundError:
            logging.error(f"Файл конфигурации {config_path} не найден")
            raise
        except yaml.YAMLError as e:
            logging.error(f"Ошибка парсинга YAML: {e}")
            raise
    
    def setup_logging(self):
        """Настройка логирования"""
        logging.basicConfig(
            level=logging.INFO,
            format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
            handlers=[
                logging.FileHandler('kafka_producer.log'),
                logging.StreamHandler()
            ]
        )
        self.logger = logging.getLogger(__name__)
    
    def setup_kafka_producer(self):
        """Настройка Kafka producer"""
        kafka_config = self.config['kafka']
        
        self.producer = KafkaProducer(
            bootstrap_servers=kafka_config['bootstrap_servers'],
            value_serializer=lambda x: json.dumps(x).encode('utf-8'),
            key_serializer=lambda x: x.encode('utf-8') if x else None
        )
        
        self.logger.info(f"Kafka producer настроен для топика: {kafka_config['topic']}")
    
    def send_test_message(self, message_data: Dict[str, Any], key: str = None):
        """Отправка тестового сообщения"""
        try:
            kafka_config = self.config['kafka']
            
            # Добавляем временную метку к сообщению
            message_data['timestamp'] = datetime.now().isoformat()
            
            # Отправляем сообщение
            future = self.producer.send(
                kafka_config['topic'],
                value=message_data,
                key=key
            )
            
            # Ждем подтверждения
            record_metadata = future.get(timeout=10)
            
            self.logger.info(
                f"Сообщение отправлено: topic={record_metadata.topic}, "
                f"partition={record_metadata.partition}, offset={record_metadata.offset}"
            )
            
        except Exception as e:
            self.logger.error(f"Ошибка отправки сообщения: {e}")
            raise
    
    def generate_test_messages(self, count: int = 10):
        """Генерация и отправка тестовых сообщений"""
        self.logger.info(f"Начинаем отправку {count} тестовых сообщений...")
        
        for i in range(count):
            # Создаем тестовое сообщение
            test_message = {
                "id": i + 1,
                "message": f"Тестовое сообщение #{i + 1}",
                "data": {
                    "value": f"value_{i}",
                    "random_number": i * 42,
                    "description": f"Описание сообщения номер {i + 1}"
                },
                "source": "kafka_producer_test"
            }
            
            # Отправляем сообщение
            self.send_test_message(test_message, key=f"key_{i}")
            
            # Небольшая пауза между сообщениями
            time.sleep(0.5)
        
        self.logger.info("Все тестовые сообщения отправлены")
    
    def cleanup(self):
        """Очистка ресурсов"""
        if self.producer:
            self.producer.flush()
            self.producer.close()
            self.logger.info("Kafka producer закрыт")
    
    def run(self, message_count: int = 10):
        """Запуск producer"""
        try:
            self.setup_kafka_producer()
            self.generate_test_messages(message_count)
        except Exception as e:
            self.logger.error(f"Ошибка запуска: {e}")
            raise
        finally:
            self.cleanup()


def main():
    """Главная функция"""
    import argparse
    
    parser = argparse.ArgumentParser(description='Kafka Test Producer')
    parser.add_argument(
        '--count', 
        type=int, 
        default=10, 
        help='Количество тестовых сообщений для отправки (по умолчанию: 10)'
    )
    parser.add_argument(
        '--config', 
        type=str, 
        default='config.yaml', 
        help='Путь к файлу конфигурации (по умолчанию: config.yaml)'
    )
    
    args = parser.parse_args()
    
    try:
        producer = KafkaTestProducer(args.config)
        producer.run(args.count)
    except Exception as e:
        logging.error(f"Критическая ошибка: {e}")
        exit(1)


if __name__ == "__main__":
    main() 