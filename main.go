// Package main represents the main module of the collector. It is the only entry point of the application.
package main

import (
	"context"
	"log"
	"time"

	"github.com/prometheus/client_golang/api"
	apiv1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

func main() {
	log.Println("Starting Collector...")
	prometheusURL := "http://prometheus.istio-system.svc.cluster.local:9090"

	client, err := api.NewClient(api.Config{
		Address: prometheusURL,
	})

	if err != nil {
		log.Fatalf("Failed to create Prometheus client: %v", err)
	}

	dataProvider := NewPrometheusProvider(apiv1.NewAPI(client))
	sampler := NewSampler(dataProvider)

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var samples []PerformanceSample
	for {
		select {
		case <-ctx.Done():
			log.Println("Collection period finished. Shutting down agent.")
			return

		case <-ticker.C:
			samples = sampler.SampleClusterData(ctx)
			log.Printf("Performance samples: %v\n", samples)
		}
	}
}
