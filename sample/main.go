package main

import (
	"fmt"
	"os"
	"time"

	"github.com/smartystreets/metrics"
)

func main() {
	metrics.StartLibrato(email, key, 1024, 4)

	// TODO: build struct and invoke method that measures stuff

	fmt.Println("Hello, World!")
}

// TODO: define struct that counts metrics

var (
	email = os.Getenv("LIBRATO_EMAIL")
	key   = os.Getenv("LIBRATO_KEY")

	Counter   = metrics.AddCounter("sample-counter", time.Second*5)
	Gauge     = metrics.AddGauge("sample-gauge", time.Second*5)
	Histogram = metrics.AddHistogram("sample-histogram", time.Second*5, 0, 1000, 5, 50, 90, 99, 99.9)
)
