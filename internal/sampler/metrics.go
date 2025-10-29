package sampler

type Metric interface {
	GetServiceName() string
	GetValue() float64
}

type ResponseTime struct {
	serviceName string
	value       float64
}

func (r ResponseTime) GetServiceName() string { return r.serviceName }
func (r ResponseTime) GetValue() float64      { return r.value }

type Throughput struct {
	serviceName string
	value       float64
}

func (t Throughput) GetServiceName() string { return t.serviceName }
func (t Throughput) GetValue() float64      { return t.value }
