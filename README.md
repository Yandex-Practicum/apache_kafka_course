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

   a. Код продюсера на Java:
   
   ```java
   public class ProducerApp {
   private static final Logger logger = LoggerFactory.getLogger(ProducerApp.class);
   public static void main(String[] args) {
   Properties props = new Properties();
   props.put("bootstrap.servers", "localhost:19092,localhost:29092,localhost:39092");
   props.put("group.id", "consumer-ssl-group");
   props.put("security.protocol", "SASL_SSL");
   props.put(SslConfigs.SSL_TRUSTSTORE_LOCATION_CONFIG, "./kafka-0-creds/kafka-0.truststore.jks");
   props.put(SslConfigs.SSL_TRUSTSTORE_PASSWORD_CONFIG, "your-password");
   props.put(SslConfigs.SSL_KEYSTORE_LOCATION_CONFIG, "./kafka-0-creds/kafka-0.keystore.jks");
   props.put(SslConfigs.SSL_KEYSTORE_PASSWORD_CONFIG, "your-password");
   props.put(SslConfigs.SSL_KEY_PASSWORD_CONFIG, "your-password");

        props.put("sasl.mechanism", "PLAIN");
        props.put("sasl.jaas.config",
                "org.apache.kafka.common.security.plain.PlainLoginModule required " +
                        "username=\"admin\" password=\"your-password\";");

        props.put("key.serializer", "org.apache.kafka.common.serialization.StringSerializer");
        props.put("value.serializer", "org.apache.kafka.common.serialization.StringSerializer");

        try (KafkaProducer<String, String> producer = new KafkaProducer<>(props)) {
            for (int i = 1; i <= 10; i++) {
                String key = "key-" + UUID.randomUUID();
                String value = "SSL message " + i;

                logger.info("Sending message {}: key={}, value={}", i, key, value);

                ProducerRecord<String, String> record = new ProducerRecord<>("secure-topic", key, value);
                producer.send(record, (metadata, exception) -> {
                    if (exception == null) {
                        logger.info("Message sent: partition={}, offset={}", metadata.partition(), metadata.offset());
                    } else {
                        logger.error("Error while producing: ", exception);
                    }
                });

                Thread.sleep(600);
            }
            producer.flush();
        } catch (Exception e) {
            logger.error("Error while producing: ", e);
        }
    }
   }
   ```
   
   b. Код консюмера на Java:
   
   ```java
   public class ConsumerApp {
   private static final Logger logger = LoggerFactory.getLogger(ConsumerApp.class);

    public static void main(String[] args) {
        Properties props = new Properties();
        props.put("bootstrap.servers", "localhost:19092,localhost:29092,localhost:39092");
        props.put("group.id", "consumer-ssl-group");
        props.put("security.protocol", "SASL_SSL");
        props.put(SslConfigs.SSL_TRUSTSTORE_LOCATION_CONFIG, "./kafka-0-creds/kafka-0.truststore.jks");
        props.put(SslConfigs.SSL_TRUSTSTORE_PASSWORD_CONFIG, "your-password");
        props.put(SslConfigs.SSL_KEYSTORE_LOCATION_CONFIG, "./kafka-0-creds/kafka-0.keystore.jks");
        props.put(SslConfigs.SSL_KEYSTORE_PASSWORD_CONFIG, "your-password");
        props.put(SslConfigs.SSL_KEY_PASSWORD_CONFIG, "your-password");


        props.put("sasl.mechanism", "PLAIN");
        props.put("sasl.jaas.config",
                "org.apache.kafka.common.security.plain.PlainLoginModule required " +
                        "username=\"admin\" password=\"your-password\";");

        props.put("key.deserializer", "org.apache.kafka.common.serialization.StringDeserializer");
        props.put("value.deserializer", "org.apache.kafka.common.serialization.StringDeserializer");

        Consumer<String, String> consumer = new KafkaConsumer<>(props);
        consumer.subscribe(Collections.singletonList("secure-topic"));

        while (true) {
            ConsumerRecords<String, String> records = consumer.poll(Duration.ofSeconds(1));
            for (ConsumerRecord<String, String> record : records) {
                logger.info("Received record: key={} value={}", record.key(), record.value());
            }
        }
    }
   }
   ```


