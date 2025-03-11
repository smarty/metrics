package alerting

import (
	"fmt"

	"github.com/smarty/metrics/v2"
)

type Monitor interface {
	Monitor(event severity)
}

type severity int

const (
	Anomaly  severity = 1
	Failure  severity = 2
	Disaster severity = 3
)

func (this severity) String() string {
	switch this {
	case Anomaly:
		return "Anomaly"
	case Failure:
		return "Failure"
	case Disaster:
		return "Disaster"
	default:
		return fmt.Sprintf("Unknown(%d)", this)
	}
}

//////////////////////////////////////////////////////////////////////////////////

type metricsMonitor struct {
	anomalies metrics.Counter
	failures  metrics.Counter
	disasters metrics.Counter
}

func NewMetricsMonitor(anomalies, failures, disasters metrics.Counter) Monitor {
	return &metricsMonitor{
		anomalies: anomalies,
		failures:  failures,
		disasters: disasters,
	}
}

func (this *metricsMonitor) Monitor(event severity) {
	switch event {
	case Anomaly:
		this.anomalies.Increment()
	case Failure:
		this.failures.Increment()
	case Disaster:
		this.disasters.Increment()
	}
}
