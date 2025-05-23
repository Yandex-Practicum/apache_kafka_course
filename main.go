package main

import (
	"github.com/lovoo/goka"
	"project/blocker"
	"project/censure"
	"project/emitter"
	"project/filter"
	"project/user"
)

var brokers = []string{"127.0.0.1:29092"}

var (
	Emitter2FilterTopic       goka.Stream = "emitter2filter-stream"
	Filter2UserProcessorTopic goka.Stream = "filter2userprocessor-stream"
	BlockerTopic              goka.Stream = "blocker-stream"
	CensorTopic               goka.Stream = "censor-stream"
)

func main() {
	go blocker.RunBlocker(brokers, BlockerTopic)
	go censure.RunCensore(brokers, CensorTopic)
	go emitter.RunEmitter(brokers, Emitter2FilterTopic)
	go filter.RunFilter(brokers, Emitter2FilterTopic, Filter2UserProcessorTopic)
	go user.RunUserProcessor(brokers, Filter2UserProcessorTopic)
	user.RunUserView(brokers)
}
