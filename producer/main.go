package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/
	"io/ioutil"
	"os"
	"strings"

	"github.com/IBM/sarama"
)

func main() {
	brokers := "rc1a-c2pm1g503ot132k8.mdb.yandexcloud.net:9091,rc1a-dbqf9j9ebna59bqh.mdb.yandexcloud.net:9091,rc1a-irgrbvr02vvg114o.mdb.yandexcloud.net:9091,rc1b-jajilbpgfs9k3cjf.mdb.yandexcloud.net:9091,rc1b-rkg2o4ahfggn38pu.mdb.yandexcloud.net:9091,rc1b-u1mk8vh088lf55b1.mdb.yandexcloud.net:9091,rc1d-6oepo6kelrc5u2rl.mdb.yandexcloud.net:9091,rc1d-i6jkmqbm87a923cg.mdb.yandexcloud.net:9091,rc1d-r3vq3q1t9h4tn7fn.mdb.yandexcloud.net:9091"
	splitBrokers := strings.Split(brokers, ",")
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Return.Successes = true
	conf.Version = sarama.V2_8_0_0
	conf.ClientID = "sasl_scram_client"
	conf.Net.SASL.Enable = true
	conf.Net.SASL.Handshake = true
	conf.Net.SASL.User = "kafka-user"
	conf.Net.SASL.Password = "kafka-password"
	conf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
	conf.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512

	certs := x509.NewCertPool()
	pemPath := "/usr/local/share/ca-certificates/Yandex/YandexInternalRootCA.crt"
	pemData, err := ioutil.ReadFile(pemPath)
	if err != nil {
		fmt.Println("Couldn't load cert: ", err.Error())
	}
	certs.AppendCertsFromPEM(pemData)

	conf.Net.TLS.Enable = true
	conf.Net.TLS.Config = &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            certs,
	}

	syncProducer, err := sarama.NewSyncProducer(splitBrokers, conf)
	if err != nil {
		fmt.Println("Couldn't create producer: ", err.Error())
		os.Exit(0)
	}

	jsonMessage := `{"id": "abc123", "name": "Sergеy", "age": 30}`

	publish(jsonMessage, syncProducer)
}

func publish(message string, producer sarama.SyncProducer) {
	msg := &sarama.ProducerMessage{
		Topic: "test",
		Value: sarama.StringEncoder(message),
	}
	p, o, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("Error publish: ", err.Error())
		return
	}

	fmt.Println("Sending message:", message)
	fmt.Println("Partition: ", p)
	fmt.Println("Offset: ", o)
}
