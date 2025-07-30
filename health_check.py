#!/usr/bin/env python3
"""
Скрипт для проверки здоровья системы Kafka и HDFS
"""

import requests
import logging
import yaml
from typing import Dict, Any
from kafka import KafkaConsumer, KafkaProducer
from hdfs import InsecureClient


class HealthChecker:
    """Класс для проверки здоровья системы"""
    
    def __init__(self, config_path: str = "config.yaml"):
        """
        Инициализация проверки здоровья
        
        Args:
            config_path: Путь к файлу конфигурации
        """
        self.config = self._load_config(config_path)
        self.setup_logging()
        
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
                logging.FileHandler('health_check.log'),
                logging.StreamHandler()
            ]
        )
        self.logger = logging.getLogger(__name__)
    
    def check_kafka_connection(self) -> bool:
        """Проверка подключения к Kafka"""
        try:
            kafka_config = self.config['kafka']
            
            # Проверяем подключение через consumer
            consumer = KafkaConsumer(
                bootstrap_servers=kafka_config['bootstrap_servers'],
                group_id='health-check-group',
                auto_offset_reset='earliest'
            )
            
            # Получаем список топиков
            topics = consumer.topics()
            consumer.close()
            
            self.logger.info(f"Kafka подключение успешно. Доступные топики: {list(topics)}")
            return True
            
        except Exception as e:
            self.logger.error(f"Ошибка подключения к Kafka: {e}")
            return False
    
    def check_kafka_producer(self) -> bool:
        """Проверка Kafka producer"""
        try:
            kafka_config = self.config['kafka']
            
            producer = KafkaProducer(
                bootstrap_servers=kafka_config['bootstrap_servers']
            )
            
            # Проверяем подключение
            producer.metrics()
            producer.close()
            
            self.logger.info("Kafka producer работает корректно")
            return True
            
        except Exception as e:
            self.logger.error(f"Ошибка Kafka producer: {e}")
            return False
    
    def check_hdfs_connection(self) -> bool:
        """Проверка подключения к HDFS"""
        try:
            hdfs_config = self.config['hdfs']
            
            client = InsecureClient(
                url=hdfs_config['url'],
                user=hdfs_config.get('user', 'root')
            )
            
            # Проверяем статус корневой директории
            status = client.status('/')
            self.logger.info(f"HDFS подключение успешно. Статус корневой директории: {status}")
            
            return True
            
        except Exception as e:
            self.logger.error(f"Ошибка подключения к HDFS: {e}")
            return False
    
    def check_hdfs_write_permissions(self) -> bool:
        """Проверка прав на запись в HDFS"""
        try:
            hdfs_config = self.config['hdfs']
            
            client = InsecureClient(
                url=hdfs_config['url'],
                user=hdfs_config.get('user', 'root')
            )
            
            # Пытаемся создать тестовую директорию
            test_dir = '/health_check_test'
            client.makedirs(test_dir, permission='755')
            
            # Пытаемся записать тестовый файл
            test_file = f'{test_dir}/test.txt'
            with client.write(test_file, overwrite=True) as writer:
                writer.write(b'health check test')
            
            # Удаляем тестовые файлы
            client.delete(test_file, recursive=False)
            client.delete(test_dir, recursive=True)
            
            self.logger.info("Права на запись в HDFS работают корректно")
            return True
            
        except Exception as e:
            self.logger.error(f"Ошибка прав на запись в HDFS: {e}")
            return False
    
    def check_hdfs_web_ui(self) -> bool:
        """Проверка HDFS Web UI"""
        try:
            hdfs_config = self.config['hdfs']
            
            # Проверяем доступность Web UI
            response = requests.get(f"{hdfs_config['url']}/jmx", timeout=10)
            
            if response.status_code == 200:
                self.logger.info("HDFS Web UI доступен")
                return True
            else:
                self.logger.error(f"HDFS Web UI недоступен. Статус: {response.status_code}")
                return False
                
        except Exception as e:
            self.logger.error(f"Ошибка доступа к HDFS Web UI: {e}")
            return False
    
    def run_all_checks(self) -> Dict[str, bool]:
        """Запуск всех проверок"""
        self.logger.info("Начинаем проверку здоровья системы...")
        
        checks = {
            'kafka_connection': self.check_kafka_connection(),
            'kafka_producer': self.check_kafka_producer(),
            'hdfs_connection': self.check_hdfs_connection(),
            'hdfs_write_permissions': self.check_hdfs_write_permissions(),
            'hdfs_web_ui': self.check_hdfs_web_ui()
        }
        
        # Выводим результаты
        self.logger.info("Результаты проверки здоровья:")
        for check_name, result in checks.items():
            status = "✅ УСПЕШНО" if result else "❌ ОШИБКА"
            self.logger.info(f"  {check_name}: {status}")
        
        # Общий статус
        all_passed = all(checks.values())
        if all_passed:
            self.logger.info("🎉 Все проверки прошли успешно!")
        else:
            self.logger.error("⚠️  Некоторые проверки не прошли")
        
        return checks


def main():
    """Главная функция"""
    import argparse
    
    parser = argparse.ArgumentParser(description='Health Check для Kafka и HDFS')
    parser.add_argument(
        '--config', 
        type=str, 
        default='config.yaml', 
        help='Путь к файлу конфигурации (по умолчанию: config.yaml)'
    )
    
    args = parser.parse_args()
    
    try:
        checker = HealthChecker(args.config)
        results = checker.run_all_checks()
        
        # Возвращаем код выхода
        exit(0 if all(results.values()) else 1)
        
    except Exception as e:
        logging.error(f"Критическая ошибка: {e}")
        exit(1)


if __name__ == "__main__":
    main() 