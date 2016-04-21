package main

import (
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/metrics"
	"github.com/smartystreets/assertions/should"
)

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

//go:generate go install github.com/smartystreets/gunit/gunit
//go:generate gunit
