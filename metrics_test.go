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
		tracker := New()

		outbound := make(chan []Measurement, 10)
		tracker.RegisterChannelDestination(outbound)
		tracker.StartMeasuring()

		Convey("Calls to Add should not result in successful registration", func() {
			a := tracker.Add("a", time.Millisecond)
			wasCounted := tracker.Count(a)

			Convey("Which means that the calls should return a 'negative' responses", func() {
				So(a, ShouldBeLessThan, 0)
				So(wasCounted, ShouldBeFalse)
			})

			Convey("And the metric should not be tracked or published", func() {
				for x := int64(0); x < 5; x++ {
					tracker.Count(a)
					tracker.Measure(a, x)
				}

				tracker.StopMeasuring()
				time.Sleep(time.Millisecond * 2)

				So(len(outbound), ShouldEqual, 0)
			})
		})

		Reset(func() {
			tracker.StopMeasuring()
		})
	})

	Convey("When tracking has been stopped", t, func() {
		tracker := New()

		outbound := make(chan []Measurement, 10)
		tracker.RegisterChannelDestination(outbound)
		a := tracker.Add("a", time.Millisecond)
		tracker.StartMeasuring()
		tracker.StopMeasuring()

		for x := int64(0); x < 5; x++ {
			tracker.Count(a)
		}

		time.Sleep(time.Millisecond * 3)

		Convey("Counts and Measurements should not be accepted", func() {
			So(len(outbound), ShouldBeGreaterThan, 0)

			for measurements := range outbound {
				for _, measurement := range measurements {
					So(measurement.Value, ShouldEqual, 0) // no counting
				}
			}
		})
	})
}

func TestMetrics(t *testing.T) {
	log.SetOutput(tWriter{t})

	Convey("Metrics should be tracked accurately", t, func() {

		// Setup...

		tracker := New()

		outbound := make(chan []Measurement, 10)
		tracker.RegisterChannelDestination(outbound)

		a := tracker.Add("a", time.Millisecond)
		b := tracker.Add("b", time.Millisecond*2)

		// Action...

		before := time.Now()

		tracker.StartMeasuring()

		for x := int64(0); x < 5; x++ {
			tracker.Count(a)
			tracker.Measure(b, x*x)
		}

		tracker.StopMeasuring()
		time.Sleep(time.Millisecond * 2)

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

		Convey("We should have at least 2 measurements", func() {
			So(len(measurements), ShouldBeGreaterThanOrEqualTo, 2)
		})

		Convey("The first measurement should reflect the counted value", func() {
			So([]time.Time{before, measurements[0].Captured, after}, ShouldBeChronological)
			So(measurements[0].Index, ShouldEqual, 0)
			So(measurements[0].Value, ShouldEqual, 5)
		})

		Convey("The second measurement should reflect the measured value", func() {
			So([]time.Time{before, measurements[1].Captured, after}, ShouldBeChronological)
			So(measurements[1].Index, ShouldEqual, 1)
			So(measurements[1].Value, ShouldEqual, 16)
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
