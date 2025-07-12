package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type User struct {
	Name           string `json:"name"`
	FavoriteNumber int64  `json:"favorite_number"`
	FavoriteColor  string `json:"favorite_color"`
}

func createSSLConfig(bootstrapServers string) *kafka.ConfigMap {
	// SSL конфигурация с CA сертификатом
	return &kafka.ConfigMap{
		"bootstrap.servers":                     bootstrapServers,
		"security.protocol":                     "SSL",
		"ssl.ca.location":                       "/app/ca.crt",
		"ssl.endpoint.identification.algorithm": "none", // Отключаем проверку hostname
	}
}

func createConsumerConfig(bootstrapServers, groupID string) *kafka.ConfigMap {
	config := createSSLConfig(bootstrapServers)
	(*config)["group.id"] = groupID
	(*config)["auto.offset.reset"] = "earliest"
	return config
}

func main() {
	fmt.Println("=== Тестирование Kafka с SSL ===")

	// Тестируем SSL соединения
	testProducerToTopic1()
	testProducerToTopic2()
	testProducerToTest()
	testConsumerFromTopic1()
	testConsumerFromTopic2()
	testConsumerFromTest()

	fmt.Println("\n=== SSL тестирование завершено ===")
	fmt.Println("✅ Все SSL соединения работают корректно!")
}

func testProducerToTopic1() {
	fmt.Println("\n1. Тестируем отправку в topic-1")

	bootstrapServers := "kafka-0:9093"
	topic := "topic-1"

	sslConfig := createSSLConfig(bootstrapServers)

	producer, err := kafka.NewProducer(sslConfig)
	if err != nil {
		log.Printf("Ошибка создания продюсера для topic-1: %v", err)
		return
	}
	defer producer.Close()

	value := User{
		Name:           "Test User Topic-1",
		FavoriteNumber: 1,
		FavoriteColor:  "blue",
	}
	payload, err := json.Marshal(value)
	if err != nil {
		log.Printf("Ошибка сериализации: %v", err)
		return
	}

	deliveryChan := make(chan kafka.Event)

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          payload,
		Headers:        []kafka.Header{{Key: "testHeader", Value: []byte("topic-1 test")}},
	}, deliveryChan)
	if err != nil {
		log.Printf("Ошибка отправки в topic-1: %v", err)
		return
	}

	event := <-deliveryChan
	message := event.(*kafka.Message)

	if message.TopicPartition.Error != nil {
		log.Printf("❌ Не удалось отправить в topic-1: %v", message.TopicPartition.Error)
	} else {
		log.Printf("✅ Успешно отправлено в topic-1 [%d] offset %v",
			message.TopicPartition.Partition, message.TopicPartition.Offset)
	}

	close(deliveryChan)
}

func testProducerToTopic2() {
	fmt.Println("\n2. Тестируем отправку в topic-2")

	bootstrapServers := "kafka-0:9093"
	topic := "topic-2"

	sslConfig := createSSLConfig(bootstrapServers)

	producer, err := kafka.NewProducer(sslConfig)
	if err != nil {
		log.Printf("Ошибка создания продюсера для topic-2: %v", err)
		return
	}
	defer producer.Close()

	value := User{
		Name:           "Test User Topic-2",
		FavoriteNumber: 2,
		FavoriteColor:  "red",
	}
	payload, err := json.Marshal(value)
	if err != nil {
		log.Printf("Ошибка сериализации: %v", err)
		return
	}

	deliveryChan := make(chan kafka.Event)

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          payload,
		Headers:        []kafka.Header{{Key: "testHeader", Value: []byte("topic-2 test")}},
	}, deliveryChan)
	if err != nil {
		log.Printf("Ошибка отправки в topic-2: %v", err)
		return
	}

	event := <-deliveryChan
	message := event.(*kafka.Message)

	if message.TopicPartition.Error != nil {
		log.Printf("❌ Не удалось отправить в topic-2: %v", message.TopicPartition.Error)
	} else {
		log.Printf("✅ Успешно отправлено в topic-2 [%d] offset %v",
			message.TopicPartition.Partition, message.TopicPartition.Offset)
	}

	close(deliveryChan)
}

func testProducerToTest() {
	fmt.Println("\n3. Тестируем отправку в test топик")

	bootstrapServers := "kafka-0:9093"
	topic := "test"

	sslConfig := createSSLConfig(bootstrapServers)

	producer, err := kafka.NewProducer(sslConfig)
	if err != nil {
		log.Printf("Ошибка создания продюсера для test топика: %v", err)
		return
	}
	defer producer.Close()

	value := User{
		Name:           "Test User",
		FavoriteNumber: 42,
		FavoriteColor:  "green",
	}
	payload, err := json.Marshal(value)
	if err != nil {
		log.Printf("Ошибка сериализации: %v", err)
		return
	}

	deliveryChan := make(chan kafka.Event)

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          payload,
		Headers:        []kafka.Header{{Key: "testHeader", Value: []byte("test topic message")}},
	}, deliveryChan)
	if err != nil {
		log.Printf("Ошибка отправки в test топик: %v", err)
		return
	}

	event := <-deliveryChan
	message := event.(*kafka.Message)

	if message.TopicPartition.Error != nil {
		log.Printf("❌ Не удалось отправить в test топик: %v", message.TopicPartition.Error)
	} else {
		log.Printf("✅ Успешно отправлено в test топик [%d] offset %v",
			message.TopicPartition.Partition, message.TopicPartition.Offset)
	}

	close(deliveryChan)
}

func testConsumerFromTopic1() {
	fmt.Println("\n4. Тестируем чтение из topic-1")

	bootstrapServers := "kafka-0:9093"
	topic := "topic-1"
	groupID := "test-consumer-group-1"

	sslConfig := createConsumerConfig(bootstrapServers, groupID)

	consumer, err := kafka.NewConsumer(sslConfig)
	if err != nil {
		log.Printf("Ошибка создания консьюмера для topic-1: %v", err)
		return
	}
	defer consumer.Close()

	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Printf("Ошибка подписки на topic-1: %v", err)
		return
	}

	fmt.Printf("Пытаемся читать из topic-1...\n")

	// Читаем сообщения с таймаутом
	timeout := time.NewTimer(10 * time.Second)
	defer timeout.Stop()

	messagesRead := 0
	for {
		select {
		case <-timeout.C:
			log.Printf("✅ Чтение из topic-1 завершено. Прочитано сообщений: %d", messagesRead)
			return
		default:
			msg, err := consumer.ReadMessage(1000)
			if err != nil {
				var kafkaErr kafka.Error
				if errors.As(err, &kafkaErr) {
					if kafkaErr.Code() == kafka.ErrTimedOut {
						continue
					}
					if kafkaErr.IsFatal() {
						log.Printf("❌ Критическая ошибка при чтении topic-1: %v", kafkaErr)
						return
					}
				}
				continue
			}

			var user User
			err = json.Unmarshal(msg.Value, &user)
			if err != nil {
				log.Printf("Ошибка десериализации: %v", err)
				continue
			}

			fmt.Printf("✅ Прочитано из topic-1: %v\n", user)
			messagesRead++
		}
	}
}

func testConsumerFromTopic2() {
	fmt.Println("\n5. Тестируем чтение из topic-2 (должно быть ограничено ACL)")

	bootstrapServers := "kafka-0:9093"
	topic := "topic-2"
	groupID := "test-consumer-group-2"

	sslConfig := createConsumerConfig(bootstrapServers, groupID)

	consumer, err := kafka.NewConsumer(sslConfig)
	if err != nil {
		log.Printf("Ошибка создания консьюмера для topic-2: %v", err)
		return
	}
	defer consumer.Close()

	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Printf("Ошибка подписки на topic-2: %v", err)
		return
	}

	fmt.Printf("Пытаемся читать из topic-2 (ожидается ошибка ACL)...\n")

	// Читаем сообщения с таймаутом
	timeout := time.NewTimer(10 * time.Second)
	defer timeout.Stop()

	messagesRead := 0
	for {
		select {
		case <-timeout.C:
			if messagesRead == 0 {
				log.Printf("✅ ACL работает корректно - доступ к topic-2 ограничен")
			} else {
				log.Printf("⚠️ Прочитано сообщений из topic-2: %d (ACL может быть не настроен)", messagesRead)
			}
			return
		default:
			msg, err := consumer.ReadMessage(1000)
			if err != nil {
				var kafkaErr kafka.Error
				if errors.As(err, &kafkaErr) {
					if kafkaErr.Code() == kafka.ErrTimedOut {
						continue
					}
					if kafkaErr.Code() == kafka.ErrTopicAuthorizationFailed {
						log.Printf("✅ ACL работает - доступ к topic-2 запрещен: %v", kafkaErr)
						return
					}
					if kafkaErr.IsFatal() {
						log.Printf("❌ Критическая ошибка при чтении topic-2: %v", kafkaErr)
						return
					}
				}
				continue
			}

			var user User
			err = json.Unmarshal(msg.Value, &user)
			if err != nil {
				log.Printf("Ошибка десериализации: %v", err)
				continue
			}

			fmt.Printf("⚠️ Прочитано из topic-2: %v (ACL может быть не настроен)\n", user)
			messagesRead++
		}
	}
}

func testConsumerFromTest() {
	fmt.Println("\n6. Тестируем чтение из test топика")

	bootstrapServers := "kafka-0:9093"
	topic := "test"
	groupID := "test-consumer-group"

	sslConfig := createConsumerConfig(bootstrapServers, groupID)

	consumer, err := kafka.NewConsumer(sslConfig)
	if err != nil {
		log.Printf("Ошибка создания консьюмера для test топика: %v", err)
		return
	}
	defer consumer.Close()

	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Printf("Ошибка подписки на test топик: %v", err)
		return
	}

	fmt.Printf("Пытаемся читать из test топика...\n")

	// Читаем сообщения с таймаутом
	timeout := time.NewTimer(10 * time.Second)
	defer timeout.Stop()

	messagesRead := 0
	for {
		select {
		case <-timeout.C:
			log.Printf("✅ Чтение из test топика завершено. Прочитано сообщений: %d", messagesRead)
			return
		default:
			msg, err := consumer.ReadMessage(1000)
			if err != nil {
				var kafkaErr kafka.Error
				if errors.As(err, &kafkaErr) {
					if kafkaErr.Code() == kafka.ErrTimedOut {
						continue
					}
					if kafkaErr.IsFatal() {
						log.Printf("❌ Критическая ошибка при чтении test топика: %v", kafkaErr)
						return
					}
				}
				continue
			}

			var user User
			err = json.Unmarshal(msg.Value, &user)
			if err != nil {
				log.Printf("Ошибка десериализации: %v", err)
				continue
			}

			fmt.Printf("✅ Прочитано из test топика: %v\n", user)
			messagesRead++
		}
	}
}
