package main

import (
	"flag"
	"log"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

var (
	word   = flag.String("word", "", "word censored")
	with   = flag.String("with", "", "new word")
	broker = flag.String("broker", "localhost:29092", "boostrap Kafka broker")
)

func main() {
	flag.Parse()
	if *word == "" {
		log.Fatalln("cannot censor word ''")
	}
	emitter, err := goka.NewEmitter([]string{*broker}, "censor-stream", new(codec.String))
	if err != nil {
		panic(err)
	}
	defer emitter.Finish()

	err = emitter.EmitSync(*word, *with)
	if err != nil {
		panic(err)
	}
}
