package main

import "context"

type PerformanceDataProvider interface {
	GetResponseTime(ctx context.Context) ([]ResponseTime, error)
	GetThroughput(ctx context.Context) ([]Throughput, error)
}
