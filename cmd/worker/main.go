package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/mathoyer/temporal-reservation/internal/config"
	"github.com/mathoyer/temporal-reservation/internal/reservation"
)

func main() {
	cfg := config.Load()

	c, err := client.Dial(client.Options{
		HostPort:  cfg.TemporalHost,
		Namespace: cfg.TemporalNS,
	})
	if err != nil {
		log.Fatalf("failed to create temporal client: %v", err)
	}
	defer c.Close()

	w := worker.New(c, cfg.TemporalTaskQ, worker.Options{})
	w.RegisterWorkflow(reservation.ReservationWorkflow)
	w.RegisterActivity(&reservation.Activities{})

	log.Printf("worker started, task queue %q", cfg.TemporalTaskQ)
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalf("worker stopped: %v", err)
	}
}
