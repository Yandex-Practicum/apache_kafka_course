# Проект модуля №5: "Безопасность в Kafka" (Python)

## Задание "Настройка защищённого соединения и управление доступом"

**Цели задания:** настроить защищённое SSL-соединение для кластера Apache 
Kafka из трёх брокеров с использованием Docker Compose, создать новый топик 
и протестировать отправку и получение зашифрованных сообщений.

**Задание:**
1. Создайте сертификаты для каждого брокера. 
2. Создайте Truststore и Keystore для каждого брокера.
3. Настройте дополнительные брокеры в режиме SSL. Ранее в курсе вы уже 
   работали с кластером Kafka, состоящим из трёх брокеров. Используйте
   имеющийся `docker-compose` кластера и настройте для него SSL. 
4. Создайте топики:
   * **topic-1**
   * **topic-2**
5. Настройте права доступа:
   * **topic-1**: доступен как для продюсеров, так и для консьюмеров.
   * **topic-2**: продюсеры могут отправлять сообщения; консьюмеры не имеют 
     доступа к чтению данных.
6. Реализуйте продюсера и консьюмера.
7. Проверьте права доступа.

## Решение

1. **Создайте сертификаты для каждого брокера.**

   a. Создаем файл конфигурации для корневого сертификата (Root CA) `ca.cnf`:
   
   ```
   [ policy_match ]
   countryName = match
   stateOrProvinceName = match
   organizationName = match
   organizationalUnitName = optional
   commonName = supplied
   emailAddress = optional
   
   [ req ]
   prompt = no
   distinguished_name = dn
   default_md = sha256
   default_bits = 4096
   x509_extensions = v3_ca
   
   [ dn ]
   countryName = RU
   organizationName = Yandex
   organizationalUnitName = Practice
   localityName = Moscow
   commonName = yandex-practice-kafka-ca
   
   [ v3_ca ]
   subjectKeyIdentifier = hash
   basicConstraints = critical,CA:true
   authorityKeyIdentifier = keyid:always,issuer:always
   keyUsage = critical,keyCertSign,cRLSign
   ```
   
   b. Создаем корневой сертификат - Root CA (локальный терминал):
   
   ```bash
   openssl req -new -nodes \
      -x509 \
      -days 365 \
      -newkey rsa:2048 \
      -keyout ca.key \
      -out ca.crt \
      -config ca.cnf
   ```
   
   c. Создаем файл для хранения сертификата безопасности `ca.pem` (локальный 
   терминал):
   
   ```bash
   cat ca.crt ca.key > ca.pem
   ```
   
   d. Создаем файлы конфигурации для каждого брокера:
   
      *  Для `kafka-0` создаем файл `kafka-0-creds/kafka-0.cnf`:
      
      ```bash
      [req]
      prompt = no
      distinguished_name = dn
      default_md = sha256
      default_bits = 4096
      req_extensions = v3_req
      
      [ dn ]
      countryName = RU
      organizationName = Yandex
      organizationalUnitName = Practice
      localityName = Moscow
      commonName = kafka-0
      
      [ v3_ca ]
      subjectKeyIdentifier = hash
      basicConstraints = critical,CA:true
      authorityKeyIdentifier = keyid:always,issuer:always
      keyUsage = critical,keyCertSign,cRLSign
      
      [ v3_req ]
      subjectKeyIdentifier = hash
      basicConstraints = CA:FALSE
      nsComment = "OpenSSL Generated Certificate"
      keyUsage = critical, digitalSignature, keyEncipherment
      extendedKeyUsage = serverAuth, clientAuth
      subjectAltName = @alt_names
      
      [ alt_names ]
      DNS.1 = kafka-0
      DNS.2 = kafka-0-external
      DNS.3 = localhost
      ```
      
      * Для `kafka-1` создаем файл `kafka-1-creds/kafka-1.cnf`:
      
      ```bash
      [req]
      prompt = no
      distinguished_name = dn
      default_md = sha256
      default_bits = 4096
      req_extensions = v3_req
      
      [ dn ]
      countryName = RU
      organizationName = Yandex
      organizationalUnitName = Practice
      localityName = Moscow
      commonName = kafka-1
      
      [ v3_ca ]
      subjectKeyIdentifier = hash
      basicConstraints = critical,CA:true
      authorityKeyIdentifier = keyid:always,issuer:always
      keyUsage = critical,keyCertSign,cRLSign
      
      [ v3_req ]
      subjectKeyIdentifier = hash
      basicConstraints = CA:FALSE
      nsComment = "OpenSSL Generated Certificate"
      keyUsage = critical, digitalSignature, keyEncipherment
      extendedKeyUsage = serverAuth, clientAuth
      subjectAltName = @alt_names
      
      [ alt_names ]
      DNS.1 = kafka-1
      DNS.2 = kafka-1-external
      DNS.3 = localhost
      ```
      
      * Для `kafka-2` создаем файл `kafka-2-creds/kafka-2.cnf`:
      
      ```bash
      [req]
      prompt = no
      distinguished_name = dn
      default_md = sha256
      default_bits = 4096
      req_extensions = v3_req
      
      [ dn ]
      countryName = RU
      organizationName = Yandex
      organizationalUnitName = Practice
      localityName = Moscow
      commonName = kafka-2
      
      [ v3_ca ]
      subjectKeyIdentifier = hash
      basicConstraints = critical,CA:true
      authorityKeyIdentifier = keyid:always,issuer:always
      keyUsage = critical,keyCertSign,cRLSign
      
      [ v3_req ]
      subjectKeyIdentifier = hash
      basicConstraints = CA:FALSE
      nsComment = "OpenSSL Generated Certificate"
      keyUsage = critical, digitalSignature, keyEncipherment
      extendedKeyUsage = serverAuth, clientAuth
      subjectAltName = @alt_names
      
      [ alt_names ]
      DNS.1 = kafka-2
      DNS.2 = kafka-2-external
      DNS.3 = localhost
      ```
   
   e. Создаем приватные ключи и запросы на сертификат - CSR (локальный терминал): 
   
   ```bash
   openssl req -new \
       -newkey rsa:2048 \
       -keyout kafka-0-creds/kafka-0.key \
       -out kafka-0-creds/kafka-0.csr \
       -config kafka-0-creds/kafka-0.cnf \
       -nodes
   
   openssl req -new \
       -newkey rsa:2048 \
       -keyout kafka-1-creds/kafka-1.key \
       -out kafka-1-creds/kafka-1.csr \
       -config kafka-1-creds/kafka-1.cnf \
       -nodes
   
   openssl req -new \
       -newkey rsa:2048 \
       -keyout kafka-2-creds/kafka-2.key \
       -out kafka-2-creds/kafka-2.csr \
       -config kafka-2-creds/kafka-2.cnf \
       -nodes
   ```
   
   f. Создаем сертификаты брокеров, подписанный CA (локальный терминал):
   
   ```bash
   openssl x509 -req \
       -days 3650 \
       -in kafka-0-creds/kafka-0.csr \
       -CA ca.crt \
       -CAkey ca.key \
       -CAcreateserial \
       -out kafka-0-creds/kafka-0.crt \
       -extfile kafka-0-creds/kafka-0.cnf \
       -extensions v3_req
   
   openssl x509 -req \
       -days 3650 \
       -in kafka-1-creds/kafka-1.csr \
       -CA ca.crt \
       -CAkey ca.key \
       -CAcreateserial \
       -out kafka-1-creds/kafka-1.crt \
       -extfile kafka-1-creds/kafka-1.cnf \
       -extensions v3_req
   
   openssl x509 -req \
       -days 3650 \
       -in kafka-2-creds/kafka-2.csr \
       -CA ca.crt \
       -CAkey ca.key \
       -CAcreateserial \
       -out kafka-2-creds/kafka-2.crt \
       -extfile kafka-2-creds/kafka-2.cnf \
       -extensions v3_req
   ```
   
   g. Создаем PKCS12-хранилища (локальный терминал):
   
   ```bash
   openssl pkcs12 -export \
       -in kafka-0-creds/kafka-0.crt \
       -inkey kafka-0-creds/kafka-0.key \
       -chain \
       -CAfile ca.pem \
       -name kafka-0 \
       -out kafka-0-creds/kafka-0.p12 \
       -password pass:your-password
   
   openssl pkcs12 -export \
       -in kafka-1-creds/kafka-1.crt \
       -inkey kafka-1-creds/kafka-1.key \
       -chain \
       -CAfile ca.pem \
       -name kafka-1 \
       -out kafka-1-creds/kafka-1.p12 \
       -password pass:your-password
   
   openssl pkcs12 -export \
       -in kafka-2-creds/kafka-2.crt \
       -inkey kafka-2-creds/kafka-2.key \
       -chain \
       -CAfile ca.pem \
       -name kafka-2 \
       -out kafka-2-creds/kafka-2.p12 \
       -password pass:your-password
   ```


2. **Создайте Truststore и Keystore для каждого брокера.**

   a. Начнем с создания Keystore (локальный терминал):
   
   ```bash
   keytool -importkeystore \
       -deststorepass your-password \
       -destkeystore kafka-0-creds/kafka.kafka-0.keystore.pkcs12 \
       -srckeystore kafka-0-creds/kafka-0.p12 \
       -deststoretype PKCS12  \
       -srcstoretype PKCS12 \
       -noprompt \
       -srcstorepass your-password
   
   keytool -importkeystore \
       -deststorepass your-password \
       -destkeystore kafka-1-creds/kafka.kafka-1.keystore.pkcs12 \
       -srckeystore kafka-1-creds/kafka-1.p12 \
       -deststoretype PKCS12  \
       -srcstoretype PKCS12 \
       -noprompt \
       -srcstorepass your-password
   
   keytool -importkeystore \
       -deststorepass your-password \
       -destkeystore kafka-2-creds/kafka.kafka-2.keystore.pkcs12 \
       -srckeystore kafka-2-creds/kafka-2.p12 \
       -deststoretype PKCS12  \
       -srcstoretype PKCS12 \
       -noprompt \
       -srcstorepass your-password
   ```
   
   b. Создаем Truststore для Kafka (локальный терминал):
   
   ```bash
   keytool -import \
       -file ca.crt \
       -alias ca \
       -keystore kafka-0-creds/kafka.kafka-0.truststore.jks \
       -storepass your-password \
       -noprompt
   
   keytool -import \
       -file ca.crt \
       -alias ca \
       -keystore kafka-1-creds/kafka.kafka-1.truststore.jks \
       -storepass your-password \
       -noprompt
   
   keytool -import \
       -file ca.crt \
       -alias ca \
       -keystore kafka-2-creds/kafka.kafka-2.truststore.jks \
       -storepass your-password \
       -noprompt
   ```
   
   c. Создаем файлы с паролями, которые указывали в предыдущих командах (локальный терминал):
   
   ```bash
   echo "your-password" > kafka-0-creds/kafka-0_sslkey_creds
   echo "your-password" > kafka-0-creds/kafka-0_keystore_creds
   echo "your-password" > kafka-0-creds/kafka-0_truststore_creds
   
   echo "your-password" > kafka-1-creds/kafka-1_sslkey_creds
   echo "your-password" > kafka-1-creds/kafka-1_keystore_creds
   echo "your-password" > kafka-1-creds/kafka-1_truststore_creds
   
   echo "your-password" > kafka-2-creds/kafka-2_sslkey_creds
   echo "your-password" > kafka-2-creds/kafka-2_keystore_creds
   echo "your-password" > kafka-2-creds/kafka-2_truststore_creds
   ```
   
   d. Импортируем PKCS12 в JKS (локальный терминал):
   
   ```bash
   keytool -importkeystore \
       -srckeystore kafka-0-creds/kafka-0.p12 \
       -srcstoretype PKCS12 \
       -destkeystore kafka-0-creds/kafka-0.keystore.jks \
       -deststoretype JKS \
       -deststorepass your-password
   
   keytool -importkeystore \
       -srckeystore kafka-1-creds/kafka-1.p12 \
       -srcstoretype PKCS12 \
       -destkeystore kafka-1-creds/kafka-1.keystore.jks \
       -deststoretype JKS \
       -deststorepass your-password
   
   keytool -importkeystore \
       -srckeystore kafka-2-creds/kafka-2.p12 \
       -srcstoretype PKCS12 \
       -destkeystore kafka-2-creds/kafka-2.keystore.jks \
       -deststoretype JKS \
       -deststorepass your-password
   ```
   
   e. Импортируем CA в Truststore (локальный терминал)::
   
   ```bash
   keytool -import -trustcacerts -file ca.crt \
       -keystore kafka-0-creds/kafka-0.truststore.jks \
       -storepass your-password -noprompt -alias ca
   
   keytool -import -trustcacerts -file ca.crt \
       -keystore kafka-1-creds/kafka-1.truststore.jks \
       -storepass your-password -noprompt -alias ca
   
   keytool -import -trustcacerts -file ca.crt \
       -keystore kafka-2-creds/kafka-2.truststore.jks \
       -storepass your-password -noprompt -alias ca
   ```
   
   f. Создаем конфигурацию для ZooKeeper (для аутентификации через SASL/PLAIN) в 
   файле `zookeeper.sasl.jaas.conf`:
   
   ```
   Server {
     org.apache.zookeeper.server.auth.DigestLoginModule required
     user_admin="your-password";
   };
   ```
   
   g. Создаем конфигурацию Kafka для авторизации в ZooKeeper в файле 
   `kafka_server_jaas.conf`:
   
   ```
   KafkaServer {
      org.apache.kafka.common.security.plain.PlainLoginModule required
      username="admin"
      password="your-password"
      user_admin="your-password"
      user_kafka="your-password"
      user_producer="your-password"
      user_consumer="your-password";
   };
   
   Client {
      org.apache.kafka.common.security.plain.PlainLoginModule required
      username="admin"
      password="your-password";
   };
   ```
   
   h. Добавим учетные записи клиента, создав файл `admin.properties`:
   
   ```
   security.protocol=SASL_SSL
   ssl.truststore.location=/etc/kafka/secrets/kafka.kafka-0.truststore.jks
   ssl.truststore.password=your-password
   sasl.mechanism=PLAIN
   sasl.jaas.config=org.apache.kafka.common.security.plain.PlainLoginModule required username="admin" password="your-password";
   ```

3. **Настройте дополнительные брокеры в режиме SSL.**

   Реализуем `docker-compose.yaml` (в нем также реализован запуск будущих 
   producer и consumer, поэтому лучше всего дождаться их реализации). Для 
   запуска используется команда:

   ```bash
   docker compose up -d
   ```

4. **Создайте топики.**

   a. После запуска контейнера проверяем, что топики еще не созданы (локальный 
   терминал):
   
   ```bash
   docker exec -it kafka-0 kafka-topics \
     --bootstrap-server kafka-0:9092 \
     --command-config /etc/kafka/secrets/admin.properties \
     --list
   ```
   
   Должны увидеть пустой вывод.
   
   b. Создаем два новых топика:
   
   ```bash
   docker exec -it kafka-0 kafka-topics \
     --bootstrap-server kafka-0:9092 \
     --command-config /etc/kafka/secrets/admin.properties \
     --create --topic topic-1 \
     --partitions 3 \
     --replication-factor 3
   
   docker exec -it kafka-0 kafka-topics \
     --bootstrap-server kafka-0:9092 \
     --command-config /etc/kafka/secrets/admin.properties \
     --create --topic topic-2 \
     --partitions 3 \
     --replication-factor 3
   ```
   
   В выводе должно быть `Created topic topic-1.` и `Created topic topic-2.` 
   соответственно.
   
   c. Проверяем созданные топики:
   
   ```bash
   docker exec -it kafka-0 kafka-topics \
     --bootstrap-server kafka-0:9092 \
     --command-config /etc/kafka/secrets/admin.properties \
     --list
   ```
   
   Должны увидеть записи о `topic-1` и `topic-2`.

5. **Настройте права доступа.**

   a. Настраиваем права доступа на запись для пользователя `producer` в топик 
   `topic-1` (локальный терминал):
   
   ```bash
   docker exec -it kafka-0 kafka-acls \
     --bootstrap-server kafka-0:9092 \
     --command-config /etc/kafka/secrets/admin.properties \
     --add \
     --allow-principal User:producer \
     --operation Write \
     --topic topic-1
   ```
   
   Должно вывестись что-то подобное:
   
   ```
   Adding ACLs for resource `ResourcePattern(resourceType=TOPIC, name=topic-1, patternType=LITERAL)`: 
           (principal=User:producer, host=*, operation=WRITE, permissionType=ALLOW) 
   ```
   
   b. Настраиваем права доступа на чтение для пользователя `consumer` в топик 
   `topic-1` (локальный терминал):
   
   ```bash
   docker exec -it kafka-0 kafka-acls \
     --bootstrap-server kafka-0:9092 \
     --command-config /etc/kafka/secrets/admin.properties \
     --add \
     --allow-principal User:consumer \
     --operation Read \
     --topic topic-1 \
     --group consumer-ssl-group
   ```
   
   Должно вывестись что-то подобное:
   
   ```
   Adding ACLs for resource `ResourcePattern(resourceType=TOPIC, name=topic-1, patternType=LITERAL)`: 
           (principal=User:consumer, host=*, operation=READ, permissionType=ALLOW) 
   
   Adding ACLs for resource `ResourcePattern(resourceType=GROUP, name=consumer-ssl-group, patternType=LITERAL)`: 
           (principal=User:consumer, host=*, operation=READ, permissionType=ALLOW)
   ```
   
   c. Настраиваем права доступа на запись для пользователя `producer` в топик 
   `topic-2` (локальный терминал):
   
   ```bash
   docker exec -it kafka-0 kafka-acls \
     --bootstrap-server kafka-0:9092 \
     --command-config /etc/kafka/secrets/admin.properties \
     --add \
     --allow-principal User:producer \
     --operation Write \
     --topic topic-2
   ```
   
   Должно вывестись что-то подобное:
   
   ```
   Adding ACLs for resource `ResourcePattern(resourceType=TOPIC, name=topic-2, patternType=LITERAL)`: 
           (principal=User:producer, host=*, operation=WRITE, permissionType=ALLOW) 
   ```

6. **Реализуйте продюсера и консьюмера.**.

   a. Код продюсера на Python:
   
   ```python
   import logging
   import os
   import time
   import uuid
   
   from confluent_kafka import Producer
   
   
   logger = logging.getLogger(__name__)
   
   logging.basicConfig(
       level=logging.INFO,
       filename="kafka_producer.log",
       format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
   )
   
   
   if __name__ == "__main__":
       logger.info("До запуска 'producer': 60 сек.")
       time.sleep(60)
       logger.info("Конфигурация 'producer'...")
   
       producer_conf = {
           "bootstrap.servers": os.getenv("KAFKA_BOOTSTRAP_SERVERS"),
   
           "security.protocol": "SASL_SSL",
           "ssl.ca.location": "ca.crt",
           "ssl.certificate.location": "kafka-0-creds/kafka-0.crt",
           "ssl.key.location": "kafka-0-creds/kafka-0.key",
   
           "ssl.endpoint.identification.algorithm": "none",
           "sasl.mechanism": "PLAIN",
           "sasl.username": os.getenv("SASL_USERNAME"),
           "sasl.password": os.getenv("SASL_PASSWORD"),
       }
   
       try:
           producer = Producer(producer_conf)
   
           while True:
               key = f"key-{uuid.uuid4()}"
               value = "SASL/PLAIN"
               producer.produce(
                   "topic-1",
                   key=key,
                   value=value,
               )
               producer.flush()
               logger.info(f"Отправлено сообщение: {key=}, {value=}")
               time.sleep(1)
       except Exception as e:
           logger.error(e)
   ```
   
   b. Код консюмера на Python:
   
   ```python
   import logging
   import os
   import time
   
   from confluent_kafka import Consumer
   
   
   logger = logging.getLogger(__name__)
   
   logging.basicConfig(
       level=logging.INFO,
       filename="kafka_consumer.log",
       format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
   )
   
   
   if __name__ == "__main__":
       logger.info("До запуска 'consumer': 60 сек.")
       time.sleep(60)
       logger.info("Конфигурация 'consumer'...")
   
       consumer_conf = {
           "bootstrap.servers": os.getenv("KAFKA_BOOTSTRAP_SERVERS"),
           "group.id": os.getenv("KAFKA_GROUP_ID"),
           "auto.offset.reset": "earliest",
   
           "security.protocol": "SASL_SSL",
           "ssl.ca.location": "ca.crt",
           "ssl.certificate.location": "kafka-0-creds/kafka-0.crt",
           "ssl.key.location": "kafka-0-creds/kafka-0.key",
   
           "ssl.endpoint.identification.algorithm": "none",
           "sasl.mechanism": "PLAIN",
           "sasl.username": os.getenv("SASL_USERNAME"),
           "sasl.password": os.getenv("SASL_PASSWORD"),
       }
   
       try:
           consumer = Consumer(consumer_conf)
           consumer.subscribe(["topic-1", "topic-2"])
   
           try:
               while True:
                   message = consumer.poll(1)
   
                   if message is None:
                       continue
                   if message.error():
                       logger.error(f"Ошибка при получении сообщения: {message.error()}")
                       continue
   
                   key = message.key().decode("utf-8")
                   value = message.value().decode("utf-8")
                   offset = message.offset()
                   logger.info(f"Получено сообщение: key='{key}', value='{value}', {offset=}")
           finally:
               consumer.close()
       except Exception as e:
               logger.error(e)
   ```

7. **Проверьте права доступа.**.

   После запуска `docker-compose.yaml` и выполнения команд, рассмотренных выше, 
   в логах продюсера должно быть что-то похожее на:
   
   ```
   2025-04-13 15:07:13,430 - __main__ - INFO - До запуска 'producer': 60 сек.
   2025-04-13 15:08:13,434 - __main__ - INFO - Конфигурация 'producer'...
   2025-04-13 15:08:13,536 - __main__ - INFO - Отправлено сообщение: key='key-f2840e33-6caf-4780-bd1f-973252cf09d6', value='SASL/PLAIN'
   2025-04-13 15:08:14,589 - __main__ - INFO - Отправлено сообщение: key='key-8afe6838-ab99-46a0-8247-83c6b616d063', value='SASL/PLAIN'
   2025-04-13 15:08:15,603 - __main__ - INFO - Отправлено сообщение: key='key-2c9783e1-395a-409b-902e-28785b226d3e', value='SASL/PLAIN'
   ```
   
   в логах консюмера что-то похожее на:
   
   ```
   2025-04-13 15:07:16,910 - __main__ - INFO - До запуска 'consumer': 60 сек.
   2025-04-13 15:08:16,913 - __main__ - INFO - Конфигурация 'consumer'...
   2025-04-13 15:08:21,142 - __main__ - INFO - Получено сообщение: key='key-f2840e33-6caf-4780-bd1f-973252cf09d6', value='SASL/PLAIN', offset=0
   2025-04-13 15:08:21,142 - __main__ - INFO - Получено сообщение: key='key-8afe6838-ab99-46a0-8247-83c6b616d063', value='SASL/PLAIN', offset=1
   2025-04-13 15:08:21,142 - __main__ - INFO - Получено сообщение: key='key-2c9783e1-395a-409b-902e-28785b226d3e', value='SASL/PLAIN', offset=2
   ```

Настройка SASL/SSL выполнена успешно. Задание выполнено.
