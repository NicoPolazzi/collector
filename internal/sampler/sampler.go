package sampler

import (
	"cmp"
	"context"
	"slices"
)

// PerformanceSample represents a sample of the performance metric of a k8s cluster.
type PerformanceSample struct {
	ServiceName    string
	ResponseTimeMs float64
	ThroughputRps  float64
}

// DataSampler is responsible to produce performance metrics' samples.
type DataSampler struct {
	provider performanceDataProvider
}

// NewDataSampler returns a pointer to the DataSampler. The sampler will use the provider to query for data.
func NewDataSampler(provider performanceDataProvider) *DataSampler {
	return &DataSampler{provider: provider}
}

// SampleClusterData produces an ordered set of PerfomanceSample. In this way we exclude not complete samples.
func (s *DataSampler) SampleClusterData(ctx context.Context) []PerformanceSample {
	const initialSamplesCount int = 0
	responseTimes, _ := s.provider.getResponseTime(ctx)
	throughputs, _ := s.provider.getThroughput(ctx)
	responseMap := convertToMap(responseTimes)
	throughputMap := convertToMap(throughputs)

	// This logic maybe can be extracted. Here we produce the sample for each service contained in the cluster, that is
	// the merging of the performance metrics.
	samples := make([]PerformanceSample, initialSamplesCount, len(responseTimes))
	for key := range throughputMap {
		if response, found := responseMap[key]; found {
			sample := PerformanceSample{
				ServiceName:    key,
				ResponseTimeMs: response,
				ThroughputRps:  throughputMap[key],
			}
			samples = append(samples, sample)
		}
	}

	sortInAscendingOrder(samples)
	return samples
}

func convertToMap[T metric](slice []T) map[string]float64 {
	lookupMap := make(map[string]float64, len(slice))

	for _, x := range slice {
		if x.GetServiceName() != "" {
			lookupMap[x.GetServiceName()] = x.GetValue()
		}
	}

	return lookupMap
}

func sortInAscendingOrder(s []PerformanceSample) {
	slices.SortFunc(s, func(a, b PerformanceSample) int {
		return cmp.Compare(a.ServiceName, b.ServiceName)
	})
}
