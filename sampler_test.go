package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockPerformanceDataProvider struct {
	responseTimes []ResponseTime
	throughputs   []Throughput
}

func (p *mockPerformanceDataProvider) GetResponseTime(ctx context.Context) []ResponseTime {
	return p.responseTimes
}

func (p *mockPerformanceDataProvider) GetThroughput(ctx context.Context) []Throughput {
	return p.throughputs
}

func TestSampleClusterData(t *testing.T) {
	provider := &mockPerformanceDataProvider{
		responseTimes: []ResponseTime{
			{
				serviceName: "test-service1",
				value:       10.33,
			},
			{
				serviceName: "test-service2",
				value:       12.33,
			},
		},
		throughputs: []Throughput{
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

	sampler := NewSampler(provider)

	samples := sampler.SampleClusterData(context.Background())
	expected := []PerformanceSample{
		{
			ServiceName:     "test-service1",
			ResponseTime_ms: 10.33,
			Throughput_rps:  40.2,
		},
		{
			ServiceName:     "test-service2",
			ResponseTime_ms: 12.33,
			Throughput_rps:  11.3,
		},
	}

	assert.ElementsMatch(t, expected, samples)

}
