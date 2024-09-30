#!/bin/bash

function create_certs() {
    mkdir "kafka-1"
    mkdir "kafka-2"
    mkdir "kafka-3"
    mkdir "client_certs"
    echo "Creating self-signed CA cert..."
    openssl genpkey -algorithm RSA -out ca.key
    openssl req -new -x509 -key ca.key -out ca.crt -days 365 -subj "/CN=Kafka-CA"

    echo "Creating server certificates..."
    openssl genpkey -algorithm RSA -out kafka-server.key

    echo "Creating CSR for the brokers..."
    openssl req -new -key kafka-server.key -out kafka-1.csr -subj "/CN=kafka-1"
    openssl req -new -key kafka-server.key -out kafka-2.csr -subj "/CN=kafka-2"
    openssl req -new -key kafka-server.key -out kafka-3.csr -subj "/CN=kafka-3"

    echo "Signing the broker certificates..."
    openssl x509 -req -in kafka-1.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out kafka-1.crt -days 365
    openssl x509 -req -in kafka-2.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out kafka-2.crt -days 365
    openssl x509 -req -in kafka-3.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out kafka-3.crt -days 365

    echo "Converting broker certificates to PKCS12 format..."
    openssl pkcs12 -export -in kafka-1.crt -inkey kafka-server.key -out kafka-1.p12 -name kafka-1 -password pass:password
    openssl pkcs12 -export -in kafka-2.crt -inkey kafka-server.key -out kafka-2.p12 -name kafka-2 -password pass:password
    openssl pkcs12 -export -in kafka-3.crt -inkey kafka-server.key -out kafka-3.p12 -name kafka-3 -password pass:password

    echo "Importing brokers certificates into a brokers Java keystore..."
    keytool -importkeystore -destkeystore ./kafka-1/kafka.keystore.jks -srckeystore kafka-1.p12 -srcstoretype PKCS12 -alias kafka-1 -storepass password -srcstorepass password
    keytool -importkeystore -destkeystore ./kafka-2/kafka.keystore.jks -srckeystore kafka-2.p12 -srcstoretype PKCS12 -alias kafka-2 -storepass password -srcstorepass password
    keytool -importkeystore -destkeystore ./kafka-3/kafka.keystore.jks -srckeystore kafka-3.p12 -srcstoretype PKCS12 -alias kafka-3 -storepass password -srcstorepass password

    echo "Creating client certificate..."
    openssl genpkey -algorithm RSA -out kafka-client.key
    openssl req -new -key kafka-client.key -out kafka-client.csr -subj "/CN=kafka-client"
    openssl x509 -req -in kafka-client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out kafka-client.crt -days 365

    echo "Converting client certificate to PKCS12 format..."
    openssl pkcs12 -export -in kafka-client.crt -inkey kafka-client.key -out kafka-client.p12 -name kafka-client -password pass:password

    echo "Import client certificate into a Java keystore..."
    keytool -importkeystore -destkeystore ./client_certs/kafka-client.jks -srckeystore kafka-client.p12 -srcstoretype PKCS12 -alias kafka-client -storepass password -srcstorepass password

    echo "Creating truststore..."
    keytool -keystore ./client_certs/client.truststore.jks -alias CARoot -import -file ca.crt -storepass password -noprompt
    keytool -keystore ./kafka-1/kafka.truststore.jks -alias CARoot -import -file ca.crt -storepass password -noprompt
    keytool -keystore ./kafka-2/kafka.truststore.jks -alias CARoot -import -file ca.crt -storepass password -noprompt
    keytool -keystore ./kafka-3/kafka.truststore.jks -alias CARoot -import -file ca.crt -storepass password -noprompt
    rm kafka*.*
    echo "Certificates created successfully!"
}

create_certs