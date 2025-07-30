#!/usr/bin/env python3
"""
Скрипт для чтения сообщений из Kafka и записи в HDFS
"""

import json
import logging
import time
from datetime import datetime
from typing import Dict, Any

from kafka import KafkaConsumer
from hdfs import InsecureClient
import yaml


class KafkaToHDFSProcessor:
    """Класс для обработки сообщений из Kafka и записи в HDFS"""
    
    def __init__(self, config_path: str = "config.yaml"):
        """
        Инициализация процессора
        
        Args:
            config_path: Путь к файлу конфигурации
        """
        self.config = self._load_config(config_path)
        self.setup_logging()
        self.consumer = None
        self.hdfs_client = None
        
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
                logging.FileHandler('kafka_to_hdfs.log'),
                logging.StreamHandler()
            ]
        )
        self.logger = logging.getLogger(__name__)
    
    def setup_kafka_consumer(self):
        """Настройка Kafka consumer"""
        kafka_config = self.config['kafka']
        
        self.consumer = KafkaConsumer(
            kafka_config['topic'],
            bootstrap_servers=kafka_config['bootstrap_servers'],
            group_id=kafka_config['group_id'],
            auto_offset_reset=kafka_config.get('auto_offset_reset', 'earliest'),
            enable_auto_commit=kafka_config.get('enable_auto_commit', True),
            value_deserializer=lambda x: x.decode('utf-8'),
            key_deserializer=lambda x: x.decode('utf-8') if x else None
        )
        
        self.logger.info(f"Kafka consumer настроен для топика: {kafka_config['topic']}")
    
    def setup_hdfs_client(self):
        """Настройка HDFS клиента"""
        hdfs_config = self.config['hdfs']
        
        self.hdfs_client = InsecureClient(
            url=hdfs_config['url'],
            user=hdfs_config.get('user', 'root')
        )
        
        # Проверяем подключение к HDFS
        try:
            self.hdfs_client.status('/')
            self.logger.info("Подключение к HDFS установлено")
        except Exception as e:
            self.logger.error(f"Ошибка подключения к HDFS: {e}")
            raise
    
    def create_hdfs_directory(self, directory_path: str):
        """Создание директории в HDFS если она не существует"""
        try:
            if not self.hdfs_client.status(directory_path, strict=False):
                self.hdfs_client.makedirs(directory_path)
                self.logger.info(f"Создана директория в HDFS: {directory_path}")
        except Exception as e:
            self.logger.error(f"Ошибка создания директории {directory_path}: {e}")
            raise
    
    def write_message_to_hdfs(self, message: str, topic: str, partition: int, offset: int):
        """Запись сообщения в HDFS"""
        try:
            # Создаем путь для записи на основе даты
            current_date = datetime.now().strftime('%Y-%m-%d')
            hdfs_path = f"/kafka_data/{topic}/{current_date}"
            
            # Создаем директорию если не существует
            self.create_hdfs_directory(hdfs_path)
            
            # Создаем имя файла с информацией о партиции и оффсете
            filename = f"partition_{partition}_offset_{offset}_{int(time.time())}.json"
            full_path = f"{hdfs_path}/{filename}"
            
            # Записываем сообщение в HDFS
            with self.hdfs_client.write(full_path, overwrite=True) as writer:
                writer.write(message.encode('utf-8'))
            
            self.logger.info(f"Сообщение записано в HDFS: {full_path}")
            
        except Exception as e:
            self.logger.error(f"Ошибка записи в HDFS: {e}")
            raise
    
    def process_messages(self):
        """Основной цикл обработки сообщений"""
        self.logger.info("Начинаем обработку сообщений из Kafka...")
        
        try:
            for message in self.consumer:
                try:
                    # Логируем информацию о сообщении
                    self.logger.info(
                        f"Получено сообщение: topic={message.topic}, "
                        f"partition={message.partition}, offset={message.offset}"
                    )
                    
                    # Записываем сообщение в HDFS
                    self.write_message_to_hdfs(
                        message.value,
                        message.topic,
                        message.partition,
                        message.offset
                    )
                    
                except Exception as e:
                    self.logger.error(f"Ошибка обработки сообщения: {e}")
                    continue
                    
        except KeyboardInterrupt:
            self.logger.info("Получен сигнал остановки")
        except Exception as e:
            self.logger.error(f"Критическая ошибка: {e}")
        finally:
            self.cleanup()
    
    def cleanup(self):
        """Очистка ресурсов"""
        if self.consumer:
            self.consumer.close()
            self.logger.info("Kafka consumer закрыт")
        
        if self.hdfs_client:
            self.logger.info("HDFS клиент закрыт")
    
    def run(self):
        """Запуск процессора"""
        try:
            self.setup_kafka_consumer()
            self.setup_hdfs_client()
            self.process_messages()
        except Exception as e:
            self.logger.error(f"Ошибка запуска: {e}")
            raise


def main():
    """Главная функция"""
    try:
        processor = KafkaToHDFSProcessor()
        processor.run()
    except Exception as e:
        logging.error(f"Критическая ошибка: {e}")
        exit(1)


if __name__ == "__main__":
    main() 