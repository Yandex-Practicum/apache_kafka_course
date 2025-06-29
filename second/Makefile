.PHONY: help up down logs status test-acl clean

help: ## Показать справку
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

up: ## Запустить Kafka кластер
	@echo "🚀 Запуск Kafka кластера с SSL и ACL..."
	docker-compose up -d
	@echo "✅ Кластер запущен!"
	@echo "💡 Используйте 'make logs' для просмотра логов"

down: ## Остановить кластер
	@echo "🛑 Остановка Kafka кластера..."
	docker-compose down
	@echo "✅ Кластер остановлен!"

logs: ## Показать логи всех сервисов
	docker-compose logs -f

logs-setup: ## Показать логи настройки
	docker-compose logs kafka-setup

logs-app: ## Показать логи Go приложения
	docker-compose logs producer-consumer

status: ## Показать статус сервисов
	@echo "📊 Статус сервисов:"
	docker-compose ps

test-acl: ## Запустить тестирование ACL
	@echo "🧪 Запуск тестирования ACL..."
	docker-compose run --rm kafka-setup /setup/test-acl.sh

test-topics: ## Показать список топиков
	@echo "📋 Список топиков:"
	docker-compose exec kafka-0 kafka-topics.sh \
		--bootstrap-server kafka-0:9093 \
		--command-config /bitnami/kafka/config/certs/client.properties \
		--list

test-acl-list: ## Показать настроенные ACL
	@echo "🔒 Настроенные ACL:"
	docker-compose exec kafka-0 kafka-acls.sh \
		--bootstrap-server kafka-0:9093 \
		--command-config /bitnami/kafka/config/certs/client.properties \
		--list

restart: down up ## Перезапустить кластер

clean: down ## Полная очистка (удаление volumes)
	@echo "🧹 Очистка данных..."
	docker-compose down -v
	docker system prune -f
	@echo "✅ Очистка завершена!"

setup-only: ## Запустить только настройку топиков и ACL
	docker-compose up kafka-setup

build: ## Собрать Go приложение локально
	@echo "🔨 Сборка Go приложения..."
	go build -o main .
	@echo "✅ Приложение собрано!"

help-ssl: ## Показать информацию о SSL настройках
	@echo "🔐 SSL конфигурация:"
	@echo "  - Все соединения зашифрованы"
	@echo "  - Keystore пароль: kafka-password"
	@echo "  - Truststore пароль: kafka-password"
	@echo "  - Взаимная аутентификация включена"

help-topics: ## Показать информацию о топиках
	@echo "📚 Настроенные топики:"
	@echo "  - topic-1: полный доступ (R/W)"
	@echo "  - topic-2: только запись (W)"
	@echo "  - test: полный доступ для тестирования"
