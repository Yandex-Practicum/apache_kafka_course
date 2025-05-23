package filter

import (
	"context"
	"log"
	"strings"

	"github.com/lovoo/goka"
	"project/blocker"
	"project/censure"
	"project/message"
)

var (
	filterGroup goka.Group = "filter"
)

func shouldDrop(ctx goka.Context, msg interface{}) bool {
	v := ctx.Join(goka.GroupTable(blocker.Group))
	msgSend, ok := msg.(*message.Message)
	if !ok {
		return false
	}
	return v != nil && v.(*blocker.BlockValue).Blocked[msgSend.From]
}

func censor(ctx goka.Context, m *message.Message) *message.Message {
	words := strings.Split(m.Content, " ")
	for i, w := range words {
		if tw := ctx.Lookup(goka.GroupTable(censure.Group), w); tw != nil {
			words[i] = tw.(string)
		}
	}
	return &message.Message{
		From:    m.From,
		To:      m.To,
		Content: strings.Join(words, " "),
	}
}

func RunFilter(brokers []string, inputTopic goka.Stream, outputTopic goka.Stream) {
	g := goka.DefineGroup(filterGroup,
		goka.Input(inputTopic, new(message.Codec), func(ctx goka.Context, msg interface{}) {
			if shouldDrop(ctx, msg) {
				return
			}
			m := censor(ctx, msg.(*message.Message))
			ctx.Emit(outputTopic, ctx.Key(), m)
		}),
		goka.Output(outputTopic, new(message.Codec)),
		goka.Join(goka.GroupTable(blocker.Group), new(blocker.BlockValueCodec)),
		goka.Lookup(goka.GroupTable(censure.Group), new(censure.ValueCodec)),
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
