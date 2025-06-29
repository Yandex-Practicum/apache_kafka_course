# Задание 2. Настройка защищённого соединения и управление доступом

Проект реализует полный Kafka кластер с SSL шифрованием и настроенным управлением доступом (ACL).

## Что реализовано

✅ **Кластер Kafka из 3 брокеров с SSL**
✅ **Настроенное управление доступом (ACL)**
✅ **Топики с разными уровнями доступа:**
- `topic-1`: полный доступ (запись + чтение)
- `topic-2`: только запись (консьюмеры не могут читать)
- `test`: тестовый топик с полным доступом

✅ **Автоматическая настройка топиков и ACL**
✅ **Go клиенты для тестирования**
✅ **Java утилиты для тестирования**

## Быстрый запуск

```bash
# Запуск кластера
docker-compose up -d

# Проверка статуса
docker-compose ps

# Просмотр логов настройки
docker-compose logs kafka-setup

# Тестирование Go клиентов
docker-compose logs producer-consumer
```

## Ручное тестирование

```bash
# Вход в контейнер для ручного тестирования
docker-compose exec kafka-0 bash

# Запуск скрипта тестирования ACL
docker-compose run --rm kafka-setup /setup/test-acl.sh
```

1. Генерация CA
- это корневой сертификат, который используется для подписания сертификатов всех участников, и брокеров, и клиентов.
```
openssl req -new -nodes -x509 -days 365 -newkey rsa:2048 -keyout ca.key -out ca.crt -config ca.cnf
cat ca.crt ca.key > ca.pem
```

2. Создание приватного ключа и запроса на сертификат (CSR)
```
openssl req -new -newkey rsa:2048 -keyout kafka-0-creds/kafka-0.key -out kafka-0-creds/kafka-0.csr -config kafka-0-creds/kafka-0.cnf -nodes

openssl req -new -newkey rsa:2048 -keyout kafka-1-creds/kafka-1.key -out kafka-1-creds/kafka-1.csr -config kafka-1-creds/kafka-1.cnf -nodes

openssl req -new -newkey rsa:2048 -keyout kafka-2-creds/kafka-2.key -out kafka-2-creds/kafka-2.csr -config kafka-2-creds/kafka-2.cnf -nodes
```

3. Создание сертификат брокера, подписанного CA
```
openssl x509 -req -days 3650 -in kafka-0-creds/kafka-0.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out kafka-0-creds/kafka-0.crt -extfile kafka-0-creds/kafka-0.cnf -extensions v3_req

openssl x509 -req -days 3650 -in kafka-1-creds/kafka-1.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out kafka-1-creds/kafka-1.crt -extfile kafka-1-creds/kafka-1.cnf -extensions v3_req

openssl x509 -req -days 3650 -in kafka-2-creds/kafka-2.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out kafka-2-creds/kafka-2.crt -extfile kafka-2-creds/kafka-2.cnf -extensions v3_req
```

4. Создание PKCS12-хранилища
```
openssl pkcs12 -export -in kafka-0-creds/kafka-0.crt -inkey kafka-0-creds/kafka-0.key -chain -CAfile ca.pem -name kafka-0 -out kafka-0-creds/kafka-0.p12 -password pass:kafka-password

openssl pkcs12 -export -in kafka-1-creds/kafka-1.crt -inkey kafka-1-creds/kafka-1.key -chain -CAfile ca.pem -name kafka-1 -out kafka-1-creds/kafka-1.p12 -password pass:kafka-password

openssl pkcs12 -export -in kafka-2-creds/kafka-2.crt -inkey kafka-2-creds/kafka-2.key -chain -CAfile ca.pem -name kafka-2 -out kafka-2-creds/kafka-2.p12 -password pass:kafka-password
```

5. Создание keystore для Kafka
```
keytool -importkeystore -deststorepass kafka-password -destkeystore kafka-0-creds/kafka.kafka-0.keystore.pkcs12 -srckeystore kafka-0-creds/kafka-0.p12 -deststoretype PKCS12 -srcstoretype PKCS12 -noprompt -srcstorepass kafka-password

keytool -importkeystore -deststorepass kafka-password -destkeystore kafka-1-creds/kafka.kafka-1.keystore.pkcs12 -srckeystore kafka-1-creds/kafka-1.p12 -deststoretype PKCS12 -srcstoretype PKCS12 -noprompt -srcstorepass kafka-password

keytool -importkeystore -deststorepass kafka-password -destkeystore kafka-2-creds/kafka.kafka-2.keystore.pkcs12 -srckeystore kafka-2-creds/kafka-2.p12 -deststoretype PKCS12 -srcstoretype PKCS12 -noprompt -srcstorepass kafka-password
```

6. Конвертируем keystore для Kafka
```
keytool -importkeystore -srckeystore kafka-0-creds/kafka.kafka-0.keystore.pkcs12 -srcstoretype PKCS12 -srcstorepass kafka-password -destkeystore kafka-0-creds/kafka.keystore.jks -deststoretype JKS -deststorepass kafka-password

keytool -importkeystore -srckeystore kafka-1-creds/kafka.kafka-1.keystore.pkcs12 -srcstoretype PKCS12 -srcstorepass kafka-password -destkeystore kafka-1-creds/kafka.keystore.jks -deststoretype JKS -deststorepass kafka-password

keytool -importkeystore -srckeystore kafka-2-creds/kafka.kafka-2.keystore.pkcs12 -srcstoretype PKCS12 -srcstorepass kafka-password -destkeystore kafka-2-creds/kafka.keystore.jks -deststoretype JKS -deststorepass kafka-password
```

6. Создание truststore для Kafka
```
keytool -import -file ca.crt -alias ca -keystore kafka-0-creds/kafka.truststore.jks -storepass kafka-password -noprompt

keytool -import -file ca.crt -alias ca -keystore kafka-1-creds/kafka.truststore.jks -storepass kafka-password -noprompt

keytool -import -file ca.crt -alias ca -keystore kafka-2-creds/kafka.truststore.jks -storepass kafka-password -noprompt
```

7. Сохранение паролей
```
echo "kafka-password" > kafka-0-creds/kafka-0_sslkey_creds
echo "kafka-password" > kafka-0-creds/kafka-0_keystore_creds
echo "kafka-password" > kafka-0-creds/kafka-0_truststore_creds

echo "kafka-password" > kafka-1-creds/kafka-1_sslkey_creds
echo "kafka-password" > kafka-1-creds/kafka-1_keystore_creds
echo "kafka-password" > kafka-1-creds/kafka-1_truststore_creds

echo "kafka-password" > kafka-2-creds/kafka-2_sslkey_creds
echo "kafka-password" > kafka-2-creds/kafka-2_keystore_creds
echo "kafka-password" > kafka-2-creds/kafka-2_truststore_creds
```

## Структура проекта

```
├── compose.yaml                 # Docker Compose конфигурация
├── main.go                     # Go клиент для тестирования
├── go.mod / go.sum            # Go модули
├── setup/
│   ├── setup-kafka.sh         # Скрипт автонастройки топиков и ACL
│   └── test-acl.sh           # Скрипт тестирования ACL
├── clients-creds/
│   ├── admin.properties       # Конфигурация админа
│   ├── producer.properties    # Конфигурация продюсера
│   └── consumer.properties    # Конфигурация консьюмера
├── kafka-0-creds/            # Сертификаты и настройки для kafka-0
├── kafka-1-creds/            # Сертификаты и настройки для kafka-1
├── kafka-2-creds/            # Сертификаты и настройки для kafka-2
├── ca.crt / ca.key           # Корневой сертификат
└── ca.cnf                    # Конфигурация CA
```

## Настройки безопасности

- **SSL шифрование:** Все соединения между брокерами и клиентами
- **Взаимная аутентификация:** Требуются сертификаты клиентов
- **ACL авторизация:** Настроенные права доступа к топикам
- **Keystore/Truststore:** JKS формат для совместимости с Kafka

## Пользователи и права

| Пользователь | topic-1 | topic-2 | Описание |
|-------------|---------|---------|----------|
| admin       | ✅ R/W   | ✅ R/W   | Полные права |
| producer    | ✅ W     | ✅ W     | Только запись |
| consumer    | ✅ R     | ❌ R     | Чтение только topic-1 |

## Проверка результата

1. **SSL соединения:** Все логи показывают SSL подключения
2. **ACL работает:** Консьюмер не может читать topic-2
3. **Топики созданы:** topic-1, topic-2, test
4. **Сертификаты:** Правильно настроены для всех брокеров
