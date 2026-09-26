package main

import (
	"log"
	"net/http"

	"go.temporal.io/sdk/client"

	"github.com/mathoyer/temporal-reservation/internal/api"
	"github.com/mathoyer/temporal-reservation/internal/config"
	"github.com/mathoyer/temporal-reservation/internal/webassets"
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

	server := api.NewServer(c, cfg.TemporalTaskQ)
	handler := server.Routes(webassets.FS())

	log.Printf("api listening on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, handler); err != nil {
		log.Fatalf("api server stopped: %v", err)
	}
}
