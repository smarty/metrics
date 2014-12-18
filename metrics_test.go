package metrics_test

import (
	"bytes"
	"log"
	"os"
	"testing"
	"time"

	"github.com/smartystreets/metrics"

	. "github.com/smartystreets/goconvey/convey/assertions"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func TestMetrics(t *testing.T) {
	log.SetOutput(tWriter{t})
	defer log.SetOutput(os.Stdout)

	log.Println("We setup a queue for receiving measurements...")
	outbound := make(chan []metrics.Measurement, 10)
	metrics.RegisterChannelDestination(outbound)

	log.Println("We configure two measurements, a and b...")
	a := metrics.Add("a", time.Millisecond)
	b := metrics.Add("b", time.Millisecond*2)

	log.Println("We start capturing metrics...")
	metrics.StartMeasuring()

	log.Println("We count and measure various things...")
	for x := int64(0); x < 5; x++ {
		metrics.Count(a)
		metrics.Measure(b, x*x)
	}

	log.Println("We stop measuring...")
	metrics.StopMeasuring()

	log.Println("We wait for measurements to be aggregated and sent...")
	time.Sleep(time.Millisecond * 2)

	log.Println("We gather the collected measurements...", len(outbound))
	measurements := []metrics.Measurement{}
	for set := range outbound {
		for _, measurement := range set {
			measurements = append(measurements, measurement)
			t.Logf("Measurement: %#v\n", measurement)
		}

		if len(outbound) == 0 {
			break
		}
	}

	log.Println("We assert that we have at least one measurement...")
	if ok, message := So(measurements, ShouldNotBeEmpty); !ok {
		t.Error(message)
		t.FailNow()
	}
}

///////////////////////////////////////////////////////////////////////////////

type tWriter struct{ *testing.T }

func (self tWriter) Write(value []byte) (int, error) {
	self.T.Log(string(bytes.TrimRight(value, "\n")))
	return len(value), nil
}

///////////////////////////////////////////////////////////////////////////////
