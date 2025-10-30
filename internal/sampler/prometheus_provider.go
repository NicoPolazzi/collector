package sampler

import (
	"context"
	"fmt"
	"log"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

const (
	responseTimeQuery = `histogram_quantile(0.95, 
		sum(rate(istio_request_duration_milliseconds_bucket{
			reporter="destination"
		}[5m])) by (le, destination_service_name)
	)`
	throughputQuery = `sum(
	rate(istio_requests_total{reporter="destination", destination_service_name!~"unknown"}[5m])
	) by (destination_service_name)`

	labelDestination = "destination_service_name"
)

// The PrometheusProvider rapresents a provider of metrics contained in a Prometheus server.
type PrometheusProvider struct {
	api v1.API
}

func NewPrometheusProvider(api v1.API) *PrometheusProvider {
	return &PrometheusProvider{api: api}
}

func (p *PrometheusProvider) getResponseTime(ctx context.Context) ([]responseTime, error) {
	result, err := p.queryPrometheus(ctx, responseTimeQuery)

	if err != nil {
		return nil, err
	}

	return collectMetricsData[responseTime](result)
}

// queryPrometheus is responsible to extract metrics from the Prometheus server using query.
// User of this method are in charge of checking the returned value type before extracting the data samples.
func (p *PrometheusProvider) queryPrometheus(ctx context.Context, query string) (model.Value, error) {
	result, warnings, err := p.api.Query(ctx, query, time.Now())

	if len(warnings) > 0 {
		log.Printf("Warnings received from Prometheus, %v", warnings)
	}

	if err != nil {
		return nil, fmt.Errorf("problem with querying Prometheus server, %v", err)
	}

	return result, nil
}

func collectMetricsData[T responseTime | throughput](result model.Value) ([]T, error) {
	vector, ok := result.(model.Vector)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	metrics := make([]T, 0, len(vector))
	for _, sample := range vector {
		serviceName := string(sample.Metric[labelDestination])
		value := float64(sample.Value)

		metrics = append(metrics, T{
			serviceName: serviceName,
			value:       value,
		})
	}

	return metrics, nil
}

func (p *PrometheusProvider) getThroughput(ctx context.Context) ([]throughput, error) {
	result, err := p.queryPrometheus(ctx, throughputQuery)

	if err != nil {
		return nil, err
	}

	return collectMetricsData[throughput](result)
}
