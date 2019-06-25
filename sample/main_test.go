package main

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/metrics"
)

func TestApplicationFixture(t *testing.T) {
	gunit.Run(new(ApplicationFixture), t)
}

type ApplicationFixture struct {
	*gunit.Fixture

	app *Application
}

func (this *ApplicationFixture) Setup() {
	this.app = new(Application)
	this.app.metrics = metrics.Capture()
}

func (this *ApplicationFixture) TestMetricsAreMeasured() {
	this.app.DoStuff(42)

	this.So(this.app.metrics.CounterValue(Counter), should.Equal, 42)
	this.So(this.app.metrics.GaugeValue(Gauge), should.Equal, 42)
	this.So(this.app.metrics.HistogramValue(Histogram).Max(), should.Equal, 42)
}
