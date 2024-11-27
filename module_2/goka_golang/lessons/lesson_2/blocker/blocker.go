package blocker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/lovoo/goka"
)

var (
	Group goka.Group = "blocker"
)

type BlockEvent struct {
	Unblock bool
}

type BlockEventCodec struct{}

func (c *BlockEventCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (c *BlockEventCodec) Decode(data []byte) (interface{}, error) {
	var m BlockEvent
	return &m, json.Unmarshal(data, &m)
}

type BlockValue struct {
	Blocked bool
}
type BlockValueCodec struct{}

func (c *BlockValueCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (c *BlockValueCodec) Decode(data []byte) (interface{}, error) {
	var m BlockValue
	return &m, json.Unmarshal(data, &m)
}

func block(ctx goka.Context, msg interface{}) {
	var s *BlockValue
	if v := ctx.Value(); v == nil {
		s = new(BlockValue)
	} else {
		s = v.(*BlockValue)
	}

	if msg.(*BlockEvent).Unblock {
		s.Blocked = false
	} else {
		s.Blocked = true
	}
	ctx.SetValue(s)
}

func RunBlocker(brokers []string, inputTopic goka.Stream) {
	g := goka.DefineGroup(Group,
		goka.Input(inputTopic, new(BlockEventCodec), block),
		goka.Persist(new(BlockValueCodec)),
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
