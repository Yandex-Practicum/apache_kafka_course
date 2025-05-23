package main

import (
	"flag"
	"log"

	"github.com/lovoo/goka"
	"project/blocker"
)

var (
	user   = flag.String("user", "", "user to block")
	block  = flag.Bool("block", true, "block user")
	name   = flag.String("name", "", "name of user to block")
	broker = flag.String("broker", "localhost:29092", "boostrap Kafka broker")
	stream = flag.String("stream", "", "stream name")
)

func main() {
	flag.Parse()
	if *user == "" {
		log.Fatal("невозможно заблокировать пользователя ''")
	}
	emitter, err := goka.NewEmitter([]string{*broker}, goka.Stream(*stream), new(blocker.BlockEventCodec))
	if err != nil {
		log.Fatal(err)
	}
	defer emitter.Finish()
	err = emitter.EmitSync(*user, &blocker.BlockEvent{Block: *block, Name: *name})
	if err != nil {
		log.Fatal(err)
	}
}
