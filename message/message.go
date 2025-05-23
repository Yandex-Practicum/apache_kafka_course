package message

import (
	"encoding/json"
	"fmt"

	"github.com/lovoo/goka"
)

var (
	Group goka.Group = "user-message"
)

// Message — содержит парамтеры пользователя
type Message struct {
	From    string // отправитель
	To      string // получатель
	Content string // данные
}

// Codec позволяет сериализовать и десериализовать Message в/из групповой таблицы.
type Codec struct{}

// Encode переводит Message в []byte
func (uc *Codec) Encode(value any) ([]byte, error) {
	if _, isMessage := value.(*Message); !isMessage {
		return nil, fmt.Errorf("тип должен быть *Message, получен %T", value)
	}
	return json.Marshal(value)
}

// Decode переводит Message из []byte в структуру.
func (uc *Codec) Decode(data []byte) (any, error) {
	var (
		c   Message
		err error
	)
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("ошибка десериализации: %v", err)
	}
	return &c, nil
}
