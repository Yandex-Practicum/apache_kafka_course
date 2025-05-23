package emitter

import (
	"log"
	"math/rand"
	"time"

	"github.com/lovoo/goka"
	"project/message"
)

const sendTimeInterval = time.Second

var (
	users = []string{
		"Alex",
		"Dian",
		"Xenia",
	}

	contents = []string{
		"Hi how are you doing",
		"Hello let's have lunch together",
		"i'm sad",
		"i'm happy",
	}
)

func RunEmitter(brokers []string, outputTopic goka.Stream) {
	var emitter, err = goka.NewEmitter(brokers, outputTopic, new(message.Codec))
	if err != nil {
		log.Fatal(err)
	}
	defer emitter.Finish()

	t := time.NewTicker(sendTimeInterval)
	defer t.Stop()

	for range t.C {
		sender := users[rand.Intn(len(users))]
		receiver := users[rand.Intn(len(users))]
		for receiver == sender {
			sender = users[rand.Intn(len(users))]
		}
		content := contents[rand.Intn(len(contents))]

		fakeUserMessage := &message.Message{
			From:    sender,
			To:      receiver,
			Content: content,
		}

		err = emitter.EmitSync(receiver, fakeUserMessage)
		if err != nil {
			log.Fatal(err)
		}
	}
}
