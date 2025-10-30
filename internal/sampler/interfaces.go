package sampler

import "context"

type performanceDataProvider interface {
	getResponseTime(ctx context.Context) ([]responseTime, error)
	getThroughput(ctx context.Context) ([]throughput, error)
}
