package app

import (
	"kafka-connect/internal/application/metric"
	"kafka-connect/internal/domain"
	"kafka-connect/internal/infra/kafka"
	"math/rand"
	"net/http"

	"github.com/google/uuid"
)

func Start(
	_ *http.Server,
	rec *kafka.MetricsReceiver,
	snd *kafka.MetricsSender,
	prc metric.Processor,
) {
	done := make(chan struct{})
	commit := make(chan bool)
	send := make(chan kafka.ListMsg)

	rch := rec.ReceiveMetrics(done, commit)
	snd.SendMetrics(done, send)

	msgCnt := 100

	// Sending metric to kafka
	go func() {
		for i := 0; i < msgCnt; i++ {
			l := domain.List{
				Alloc: domain.Metric{
					Type:        domain.Gauge,
					Name:        "Alloc",
					Description: "",
					Value:       rand.Int63n(100),
				},
				FreeMemory: domain.Metric{
					Type:        domain.Gauge,
					Name:        "FreeMemory",
					Description: "",
					Value:       rand.Int63n(100),
				},
				PollCount: domain.Metric{
					Type:        domain.Counter,
					Name:        "PollCount",
					Description: "",
					Value:       rand.Int63n(100),
				},
				TotalMemory: domain.Metric{
					Type:        domain.Gauge,
					Name:        "TotalMemory",
					Description: "",
					Value:       rand.Int63n(100),
				},
			}
			send <- kafka.ListMsg{
				Key: uuid.New().String(),
				Msg: l,
			}
		}
	}()

	// Receiving metric from kafka
	go func() {
		defer close(done)
		defer close(commit)
		defer close(send)
		for i := 0; i < msgCnt; i++ {
			select {
			case <-done:
				return
			case l := <-rch:
				commit <- prc.Process(l) == nil
			}
		}
	}()
}
