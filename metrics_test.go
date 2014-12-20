package metrics

import (
	"bytes"
	"log"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConventions(t *testing.T) {
	log.SetOutput(tWriter{t})

	Convey("When tracking has already been started", t, func() {
		metrics := New()

		outbound := make(chan []Measurement, 10)
		metrics.RegisterChannelDestination(outbound)
		metrics.StartMeasuring()

		Convey("Calls to Add should not result in successful registration", func() {
			a := metrics.Add("a", time.Millisecond)
			wasCounted := metrics.Count(a)

			Convey("Which means that the calls should return a 'negative' responses", func() {
				So(a, ShouldBeLessThan, 0)
				So(wasCounted, ShouldBeFalse)
			})

			Convey("And the metric should not be tracked or published", func() {
				for x := int64(0); x < 5; x++ {
					metrics.Count(a)
					metrics.Measure(a, x)
				}

				metrics.StopMeasuring()
				time.Sleep(time.Millisecond * 2)

				So(len(outbound), ShouldEqual, 0)
			})
		})

		Reset(func() {
			metrics.StopMeasuring()
		})
	})

	Convey("When tracking has been stopped", t, func() {
		metrics := New()

		outbound := make(chan []Measurement, 10)
		metrics.RegisterChannelDestination(outbound)
		a := metrics.Add("a", time.Millisecond)
		metrics.StartMeasuring()
		metrics.StopMeasuring()

		for x := int64(0); x < 5; x++ {
			metrics.Count(a)
		}

		time.Sleep(time.Millisecond * 5)

		Convey("Counts and Measurements should no longer be accepted, so the value of the metric should remain at 0", func() {
			So(len(outbound), ShouldEqual, 1)
			measurements := <-outbound
			So(len(measurements), ShouldEqual, 1)
			So(measurements[0].Value, ShouldEqual, 0)
		})
	})
}

func TestMetrics(t *testing.T) {
	log.SetOutput(tWriter{t})

	Convey("Metrics should be tracked accurately", t, func() {

		// Setup...

		metrics := New()

		outbound := make(chan []Measurement, 10)
		metrics.RegisterChannelDestination(outbound)

		a := metrics.Add("a", time.Millisecond)
		b := metrics.Add("b", time.Millisecond*2)

		// Action...

		before := time.Now()

		metrics.StartMeasuring()

		// first two measurements
		for x := int64(0); x < 5; x++ {
			metrics.Count(a)
			metrics.Measure(b, x*x)
		}
		time.Sleep(time.Millisecond * 2)

		// last two measurements
		for x := int64(0); x < 5; x++ {
			metrics.Count(a)
			metrics.Measure(b, x*x+1)
		}
		time.Sleep(time.Millisecond * 2)

		metrics.StopMeasuring()

		after := time.Now()

		// Gather...

		measurements := []Measurement{}
		for set := range outbound {
			for _, measurement := range set {
				measurements = append(measurements, measurement)
			}

			if len(outbound) == 0 {
				break
			}
		}

		// Assert...

		Convey("We should have at least 5 measurements, and they should be in chronological order", func() {
			So(len(measurements), ShouldBeGreaterThanOrEqualTo, 5)
			var (
				first  = measurements[0].Captured
				second = measurements[1].Captured
				third  = measurements[2].Captured
				fourth = measurements[3].Captured
				fifth  = measurements[4].Captured
			)
			So([]time.Time{before, first, second, third, fourth, fifth, after}, ShouldBeChronological)
		})

		Convey("The first measurement should reflect the _counted_ value", func() {
			So(measurements[0].ID, ShouldEqual, 0)
			So(measurements[0].Value, ShouldEqual, 5)
		})

		Convey("The second measurement should reflect the _measured_ value", func() {
			So(measurements[1].ID, ShouldEqual, 1)
			So(measurements[1].Value, ShouldEqual, 16)
		})

		Convey("The third measurement should reflect the _counted_ value", func() {
			So(measurements[2].ID, ShouldEqual, 0)
			So(measurements[2].Value, ShouldEqual, 10)
		})

		Convey("The fourth measurement should reflect the _counted_ value", func() {
			So(measurements[3].ID, ShouldEqual, 0)
			So(measurements[3].Value, ShouldEqual, 10)
		})

		Convey("The fifth measurement should reflect the _measured_ value", func() {
			So(measurements[4].ID, ShouldEqual, 1)
			So(measurements[4].Value, ShouldEqual, 17)
		})
	})
}

///////////////////////////////////////////////////////////////////////////////

type tWriter struct{ *testing.T }

func (self tWriter) Write(value []byte) (int, error) {
	self.T.Log(string(bytes.TrimRight(value, "\n")))
	return len(value), nil
}

///////////////////////////////////////////////////////////////////////////////

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

///////////////////////////////////////////////////////////////////////////////
