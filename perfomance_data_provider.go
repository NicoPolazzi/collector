package main

import "context"

type PerformanceDataProvider interface {
	GetResponseTime(ctx context.Context) []ResponseTime
	GetThroughput(ctx context.Context) []Throughput
}
