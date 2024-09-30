# KRAFT5
## Итоговый проект пятого модуля

## Подготовка 
1. Клонировать репозиторий
2. Создать виртуальное окружение
```bash
python -m venv venv
```
3. Активировать виртуальное окружение
```bash
venv\Scripts\Activate.ps1
```
4. Установить зависимости
```bash
pip install -r requirements.txt
```
5. Получить сертификат
```bash
mkdir -p ./secrets/ && \
curl -o ./secrets/YandexInternalRootCA.crt "https://storage.yandexcloud.net/cloud-certs/CA.pem" && \
chmod 0655 ./secrets/YandexInternalRootCA.crt
```

## Задание 1.Развёртывание и настройка Kafka-кластера в Yandex Cloud

### Шаг 1. Разверните Kafka. Для этого:

Укажите оптимальные параметры аппаратных ресурсов для брокеров (количество дисков, CPU, RAM):

![alt text](resources/kafka_info1.png)

Разверните кластер Kafka в Yandex Cloud с 3 брокерами:

![alt text](resources/kafka_info2.png)

### Шаг 2. Настройте репликацию и хранение данных:

Создайте топик с 3 партициями и коэффициентом репликации 3:

![alt text](resources/topic_info.png)

Настройте политику очистки логов:

![alt text](resources/clear_logs.png)

Установите параметры хранения (log.retention.ms, log.segment.bytes):

![alt text](resources/storage_parameters.png)

### Шаг 3. Настройте Schema Registry:

Разверните Schema Registry:

![alt text](resources/schema_reg.png)

Зарегистрируйте схему данных:

![alt text](resources/schema_data.png)

### Шаг 4. Проверьте работу Kafka:

Напишите простой продюсер и консьюмер:

В корне имеются producer.py и consumer.py.

Отправьте тестовые сообщения и убедитесь, что они передаются через Kafka:

При запуске продюсера видим в терминале следующее:

![alt text](resources/producer_terminal.png)

При запуске консьюмера видим в терминале следующее:

![alt text](resources/consumer_terminal.png)

## Задание 2. Интеграция Kafka с внешними системами (Apache NiFi / Hadoop)

### Шаг 1. Преобразовать сертификат в формат PKCS12:
```bash
keytool -importcert -alias YandexRootCA -file ./secrets/YandexInternalRootCA.crt -keystore truststore.p12 -storetype PKCS12
```

### Шаг 2. Поднять контейнер с Apache NiFi:
```bash
docker compose up -d
```

### Шаг 3. Перейти по адресу http://localhost:8080/nifi/:

![alt text](resources/nifi_start.png)

### Шаг 4. Сконфигурировать процессоры:

![alt text](resources/processors.png)

### Шаг 5. Проверить вывод консьюмера:

![alt text](resources/nifi_messages.png)