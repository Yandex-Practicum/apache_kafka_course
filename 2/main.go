package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/colinmarc/hdfs/v2"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
)

func runConsumer() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  "localhost:9094",
		"group.id":           "hadoop-consumer-group-3",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": true,
		"session.timeout.ms": 6000,
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}

	err = consumer.Subscribe("hadoop-topic", nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %s", err)
	}
	defer consumer.Close()

	dial := &net.Dialer{
		Timeout:   60 * time.Second,
		KeepAlive: 60 * time.Second,
	}

	clientOptions := hdfs.ClientOptions{
		Addresses:           []string{"localhost:9000"},
		User:                "root",
		NamenodeDialFunc:    dial.DialContext,
		DatanodeDialFunc:    dial.DialContext,
		UseDatanodeHostname: true,
	}

	hdfsClient, err := hdfs.NewClient(clientOptions)
	if err != nil {
		log.Fatalf("Failed to create HDFS client: %s", err)
	}

	dataDir := "/data"
	if err = hdfsClient.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create directory %s: %s", dataDir, err)
	}

	for {
		msg := &kafka.Message{}
		msg, err = consumer.ReadMessage(100 * time.Millisecond)
		if err == nil {
			value := string(msg.Value)
			fmt.Printf("Received message: value=%s, partition=%v\n", value, msg.TopicPartition)

			var writer *hdfs.FileWriter
			hdfsFile := fmt.Sprintf("/data/message_%s", uuid.New().String())

			writer, err = hdfsClient.Create(hdfsFile)
			if err != nil {
				log.Printf("Failed to create HDFS file: %s", err)
				continue
			}

			_, err = writer.Write([]byte(value + "\n"))
			if err != nil {
				log.Printf("Failed to write to HDFS file: %s", err)
				writer.Close()
				continue
			}

			err = writer.Close()
			if err != nil {
				log.Printf("Failed to close HDFS file: %s", err)
			}

			fmt.Printf("Message '%s' written to HDFS at path: '%s' \n", value, hdfsFile)
		} else {
			var kafkaErr kafka.Error
			if errors.As(err, &kafkaErr) && kafkaErr.IsFatal() {
				log.Fatalf("Fatal error: %s", kafkaErr)
			}
		}
	}
}

func runProducer() {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9094",
		"client.id":         "hadoop-producer",
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()

	topic := "hadoop-topic"

	// Канал для получения событий доставки
	deliveryChan := make(chan kafka.Event)

	// Отправляем тестовые сообщения
	messages := []string{
		"Привет из Kafka! Это первое сообщение",
		"Второе сообщение с данными для HDFS",
		"JSON данные: {\"user\": \"Ivan\", \"action\": \"login\", \"timestamp\": \"2024-12-31T12:00:00Z\"}",
		"CSV данные: user1,action1,2024-12-31 12:01:00",
		"Лог сообщение: [INFO] Пользователь успешно авторизован",
		"Метрики: CPU=75%, Memory=2.3GB, Disk=45%",
		"Событие: Новый заказ #12345 на сумму 1500 рублей",
		"Последнее тестовое сообщение для Hadoop интеграции",
	}

	fmt.Println("Начинаю отправку сообщений в Kafka...")

	for i, message := range messages {
		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(message),
			Key:            []byte(fmt.Sprintf("key-%d", i)),
		}, deliveryChan)

		if err != nil {
			log.Printf("Failed to produce message: %s", err)
			continue
		}

		// Ждем подтверждения доставки
		e := <-deliveryChan
		m := e.(*kafka.Message)

		if m.TopicPartition.Error != nil {
			log.Printf("Delivery failed: %v", m.TopicPartition.Error)
		} else {
			fmt.Printf("✅ Сообщение доставлено: topic=%s, partition=%d, offset=%v, value=%s\n",
				*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset, string(m.Value))
		}

		// Небольшая пауза между сообщениями
		time.Sleep(2 * time.Second)
	}

	// Ждем отправки всех сообщений
	fmt.Println("Ожидание завершения отправки всех сообщений...")
	producer.Flush(15 * 1000)

	fmt.Println("✅ Все сообщения отправлены успешно!")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование:")
		fmt.Println("  go run main.go producer  - запустить продюсер")
		fmt.Println("  go run main.go consumer  - запустить консьюмер")
		return
	}

	mode := os.Args[1]
	switch mode {
	case "producer":
		fmt.Println("🚀 Запускаю продюсер...")
		runProducer()
	case "consumer":
		fmt.Println("🚀 Запускаю консьюмер...")
		runConsumer()
	default:
		fmt.Printf("Неизвестный режим: %s\n", mode)
		fmt.Println("Используйте 'producer' или 'consumer'")
	}
}
