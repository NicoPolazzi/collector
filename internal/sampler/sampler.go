package sampler

import (
	"cmp"
	"context"
	"slices"
)

type PerformanceSample struct {
	ServiceName    string
	ResponseTimeMs float64
	ThroughputRps  float64
}

type DataSampler struct {
	provider PerformanceDataProvider
}

func NewDataSampler(provider PerformanceDataProvider) *DataSampler {
	return &DataSampler{provider: provider}
}

func (s *DataSampler) SampleClusterData(ctx context.Context) []PerformanceSample {
	const initialSamplesCount int = 0
	responseTimes, _ := s.provider.GetResponseTime(ctx)
	throughputs, _ := s.provider.GetThroughput(ctx)
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

func convertToMap[T Metric](slice []T) map[string]float64 {
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
