package main

import "context"

type PerformanceSample struct {
	ServiceName     string
	ResponseTime_ms float64
	Throughput_rps  float64
}

type Sampler struct {
	provider PerformanceDataProvider
}

func NewSampler(provider PerformanceDataProvider) *Sampler {
	return &Sampler{provider: provider}
}

func (s *Sampler) SampleClusterData(ctx context.Context) []PerformanceSample {
	responseTimes := s.provider.GetResponseTime(ctx)
	throughputs := s.provider.GetThroughput(ctx)

	throughputMap := make(map[string]float64, len(throughputs))
	for _, t := range throughputs {
		if t.serviceName != "" {
			throughputMap[t.serviceName] = t.value
		}
	}

	responseMap := make(map[string]float64, len(responseTimes))
	for _, r := range responseTimes {
		if r.serviceName != "" {
			responseMap[r.serviceName] = r.value
		}
	}

	samples := make([]PerformanceSample, 0, len(responseTimes))
	for key := range throughputMap {
		if response, found := responseMap[key]; found {
			sample := PerformanceSample{
				ServiceName:     key,
				ResponseTime_ms: response,
				Throughput_rps:  throughputMap[key],
			}
			samples = append(samples, sample)
		}
	}

	return samples
}
