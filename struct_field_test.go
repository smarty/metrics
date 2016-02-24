package metrics

import (
	"testing"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func TestMetricsAsStructField(t *testing.T) {
	thing := &Thing{}
	thing.metrics = Capture()
	thing.Stuff()

	assert := assertions.New(t)
	assert.So(thing.metrics.CounterValue(1), should.Equal, 1)
	assert.So(thing.metrics.CounterValue(2), should.Equal, 2)
	assert.So(thing.metrics.CounterValue(3), should.Equal, 3)
	assert.So(thing.metrics.CounterValue(4), should.Equal, 4)
}

///////////////////////////////////////////////////////////////////////////////

type Thing struct {
	metrics *Metrics
}

func (this *Thing) Stuff() {
	this.metrics.Count(1)
	this.metrics.CountN(2, 2)
	this.metrics.RawCount(3, 3)
	this.metrics.Measure(4, 4)
}