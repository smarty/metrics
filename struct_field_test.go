package metrics

import (
	"testing"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func TestMetricsAsStructField(t *testing.T) {
	thing := &Thing{}
	thing.metrics = Capture()
	thing.DoStuff()

	assert := assertions.New(t)
	assert.So(thing.metrics.CounterValue(0), should.Equal, 1)
	assert.So(thing.metrics.CounterValue(1), should.Equal, 2)
	assert.So(thing.metrics.CounterValue(2), should.Equal, 3)
	assert.So(thing.metrics.GaugeValue(3), should.Equal, 4)
	assert.So(thing.metrics.HistogramValue(4).Max(), should.Equal, 5)
}

///////////////////////////////////////////////////////////////////////////////

type Thing struct {
	metrics *Metrics
}

func (this *Thing) DoStuff() {
	this.metrics.RawCount(2, 3)
	this.metrics.Measure(3, 4)
	this.metrics.Record(4, 5)
	this.metrics.Count(0)
	this.metrics.CountN(1, 2)
}
