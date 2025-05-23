package main

import (
	"kafka-connect/internal/application/app"
	"kafka-connect/internal/infra/di"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		di.InjectApp(),
		fx.Invoke(app.Start),
	).Run()
}
