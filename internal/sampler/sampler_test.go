package sampler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockPerformanceDataProvider struct {
	responseTimes []responseTime
	throughputs   []throughput
}

func (p *mockPerformanceDataProvider) getResponseTime(_ context.Context) ([]responseTime, error) {
	return p.responseTimes, nil
}

func (p *mockPerformanceDataProvider) getThroughput(_ context.Context) ([]throughput, error) {
	return p.throughputs, nil
}

func TestSampleClusterData(t *testing.T) {

	t.Run("returns the correct sample of the cluster", func(t *testing.T) {
		provider := &mockPerformanceDataProvider{
			responseTimes: []responseTime{
				{
					serviceName: "test-service1",
					value:       10.33,
				},
				{
					serviceName: "test-service2",
					value:       12.33,
				},
			},
			throughputs: []throughput{
				{
					serviceName: "test-service2",
					value:       11.3,
				},
				{
					serviceName: "test-service1",
					value:       40.2,
				},
			},
		}

		sampler := NewDataSampler(provider)
		samples := sampler.SampleClusterData(context.Background())

		expected := []PerformanceSample{
			{
				ServiceName:    "test-service1",
				ResponseTimeMs: 10.33,
				ThroughputRps:  40.2,
			},
			{
				ServiceName:    "test-service2",
				ResponseTimeMs: 12.33,
				ThroughputRps:  11.3,
			},
		}

		assert.Equal(t, expected, samples)
	})

	t.Run("discard incomplete samples", func(t *testing.T) {
		provider := &mockPerformanceDataProvider{
			responseTimes: []responseTime{
				{
					serviceName: "test-service2",
					value:       10.33,
				},
			},
			throughputs: []throughput{
				{
					serviceName: "test-service2",
					value:       11.3,
				},
				{
					serviceName: "test-service1",
					value:       40.2,
				},
			},
		}

		sampler := NewDataSampler(provider)
		samples := sampler.SampleClusterData(context.Background())

		expected := []PerformanceSample{
			{
				ServiceName:    "test-service2",
				ResponseTimeMs: 10.33,
				ThroughputRps:  11.3,
			},
		}

		assert.Equal(t, expected, samples)
	})
}
