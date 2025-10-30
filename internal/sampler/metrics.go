package sampler

type metric interface {
	GetServiceName() string
	GetValue() float64
}

type responseTime struct {
	serviceName string
	value       float64
}

func (r responseTime) GetServiceName() string { return r.serviceName }
func (r responseTime) GetValue() float64      { return r.value }

type throughput struct {
	serviceName string
	value       float64
}

func (t throughput) GetServiceName() string { return t.serviceName }
func (t throughput) GetValue() float64      { return t.value }
