package main

import (
	"time"

	"github.com/smartystreets/metrics"
)

func main() {
	counter1 := metrics.AddCounter("sample-counter1", time.Second*1)
	// counter2 := metrics.AddCounter("sample-counter2", time.Second*1)
	// counter3 := metrics.AddCounter("sample-counter3", time.Second*1)

	// gauge1 := metrics.AddGauge("sample-gauge1", time.Second*1)
	// gauge2 := metrics.AddGauge("sample-gauge2", time.Second*1)
	// gauge3 := metrics.AddGauge("sample-gauge3", time.Second*1)
	queue := make(chan []metrics.Measurement, 64)
	metrics.RegisterChannelDestination(queue)
	metrics.StartMeasuring()
	metrics.Count(counter1)
	// metrics.Count(counter2)
	// metrics.Count(counter3)
	// metrics.Measure(gauge1, 41)
	// metrics.Measure(gauge2, 42)
	// metrics.Measure(gauge3, 43)

	librato := metrics.WithLibrato("librato-operations@smartystreets.com", "417d3fd5dd34621420f1e1871d220bd78f9187652b530b16b4386bcc71ec16d5", "mymachinename", 2)
	librato.Listen(queue)
}
