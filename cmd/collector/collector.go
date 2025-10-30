// Command collector is the main entrypoint for our application.
// It initializes the provider and starts the data collection loop.
package main

import (
	"context"
	"log"
	"time"

	"github.com/nicopolazzi/collector/internal/sampler"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
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

	dataProvider := sampler.NewPrometheusProvider(v1.NewAPI(client))
	dataSampler := sampler.NewDataSampler(dataProvider)

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var samples []sampler.PerformanceSample
	for {
		select {
		case <-ctx.Done():
			log.Println("Collection period finished. Shutting down agent.")
			return

		case <-ticker.C:
			samples = dataSampler.SampleClusterData(ctx)
			log.Printf("Performance samples: %v\n", samples)
		}
	}
}
