package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/apache_kafka_course/module1/go/avro-example/app"
)

func main() {
	application, err := app.Fabric()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("application starts with cfg -> %s \n", application.GetConfig())
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	application.Start(ctx)
}
