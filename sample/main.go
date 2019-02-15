package main

import (
	"fmt"
	"os"
	"time"

	"github.com/smartystreets/metrics"
)

func main() {
	metrics.StartAppOptics(key, 1024, 4)

	app := new(Application)

	for x := int64(0); ; x++ {
		fmt.Println(x)
		app.DoStuff(x)
		time.Sleep(time.Millisecond * 100)
	}
}

////////////////////////////////////////////////////////////////////////////

type Application struct {
	metrics *metrics.Metrics
}

func (this *Application) DoStuff(x int64) {
	this.metrics.CountN(Counter, x)
	this.metrics.Measure(Gauge, x)
	this.metrics.Record(Histogram, x%1000)
}

////////////////////////////////////////////////////////////////////////////

var (
	key = os.Getenv("APPOPTICS_KEY")

	Counter   = metrics.AddCounter("sample-counter", time.Second*5)
	Gauge     = metrics.AddGauge("sample-gauge", time.Second*5)
	Histogram = metrics.AddHistogram("sample-histogram", time.Second*5, 0, 1000, 5, 50, 90, 99, 99.9)
)

func init() {
	metrics.TagCounter(Counter, "color", "blue")
}