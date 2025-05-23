package user

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lovoo/goka"
	"project/message"
)

var (
	Group goka.Group = "user-message"
)

type MessagesCollected struct {
	Messages map[string][]string
}

type Codec struct{}

// Encode переводит MessagesCollected в []byte
func (up *Codec) Encode(value any) ([]byte, error) {
	if _, isMessageCollected := value.(*MessagesCollected); !isMessageCollected {
		return nil, fmt.Errorf("тип должен быть *MessagesCollected, получен %T", value)
	}
	return json.Marshal(value)
}

// Decode переводит MessagesCollected из []byte в структуру.
func (up *Codec) Decode(data []byte) (any, error) {
	var (
		p   MessagesCollected
		err error
	)
	err = json.Unmarshal(data, &p)
	if err != nil {
		return nil, fmt.Errorf("ошибка десериализации: %v", err)
	}
	return &p, nil
}

func process(ctx goka.Context, msg any) {
	var msgRecived *message.Message
	var ok bool
	var msgCollected *MessagesCollected

	if msgRecived, ok = msg.(*message.Message); !ok {
		return
	}

	if val := ctx.Value(); val != nil {
		msgCollected = val.(*MessagesCollected)
	} else {
		msgCollected = &MessagesCollected{Messages: make(map[string][]string)}
	}

	if len(msgCollected.Messages[msgRecived.From]) < 5 {
		msgCollected.Messages[msgRecived.From] = append(msgCollected.Messages[msgRecived.From], msgRecived.Content)
	} else {
		msgCollected.Messages[msgRecived.From] = append(msgCollected.Messages[msgRecived.From][1:5], msgRecived.Content)
	}
	ctx.SetValue(msgCollected)
}

func RunUserProcessor(brokers []string, inputStream goka.Stream) {
	g := goka.DefineGroup(Group,
		goka.Input(inputStream, new(message.Codec), process),
		goka.Persist(new(Codec)),
	)
	p, err := goka.NewProcessor(brokers, g)
	if err != nil {
		log.Fatal(err)
	}
	err = p.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func RunUserView(brokers []string) {
	view, err := goka.NewView(brokers,
		goka.GroupTable(Group),
		new(Codec),
	)
	if err != nil {
		log.Fatal(err)
	}

	root := mux.NewRouter()
	root.HandleFunc("/{key}", func(w http.ResponseWriter, r *http.Request) {
		value, _ := view.Get(mux.Vars(r)["key"])
		data, _ := json.Marshal(value)
		w.Write(data)
	})
	log.Println("View opened at http://localhost:9095/")
	go func() {
		err = http.ListenAndServe(":9095", root)
		if err != nil {
			log.Fatal(err)
		}
	}()
	err = view.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
