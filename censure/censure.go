package censure

import (
	"context"
	"log"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

var (
	Group goka.Group = "censor"
)

type ValueCodec struct {
	codec.String
}

func replaceWord(ctx goka.Context, msg interface{}) {
	ctx.SetValue(msg.(string))
}

func RunCensore(broker []string, inputStream goka.Stream) {
	g := goka.DefineGroup(Group,
		goka.Input(inputStream, new(ValueCodec), replaceWord),
		goka.Persist(new(ValueCodec)),
	)
	p, err := goka.NewProcessor(broker, g)
	if err != nil {
		log.Fatal(err)
	}
	err = p.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
