package metrics

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type Librato struct {
	email          string
	key            string
	hostname       string
	maxRequests    int32
	activeRequests int32
	buffer         map[int]MetricMeasurement
	client         *http.Client
}

func newLibrato(email, key, hostname string, maxRequests int32) *Librato {
	// TODO: validate inputs

	client := &http.Client{
		Transport: &http.Transport{DisableCompression: true},
		Timeout:   time.Duration(time.Second * 10),
	}

	return &Librato{
		email:       email,
		key:         key,
		hostname:    hostname,
		maxRequests: maxRequests,
		buffer:      map[int]MetricMeasurement{},
		client:      client,
	}
}

func (this *Librato) Listen(queue chan []Measurement) {
	for measurements := range queue {
		for _, measurement := range measurements {
			this.buffer[measurement.ID] = measurement // last one wins
		}
		if len(queue) == 0 {
			this.publish()
		}
	}
}

func (this *Librato) publish() {
	for {
		if len(this.buffer) == 0 {
			break // no more work to do
		}

		active := atomic.LoadInt32(&this.activeRequests)
		available := this.maxRequests - active
		if available == 0 {
			log.Printf(logSkippingPublish, active, this.maxRequests)
			break // all lanes are busy
		}

		// how many requests would we need to take care of all of the items
		needed := int32(countBatches(len(this.buffer)))
		if needed > available {
			log.Printf(logTruncatingPublish, needed, available)
			needed = available // not enough open/available lanes to accommodate all the requests
		}

		for i := int32(0); i < needed; i++ {
			body := this.serializeNext()
			request := this.buildRequest(body)
			atomic.AddInt32(&this.activeRequests, 1)

			go func() {
				response, err := this.client.Do(request)
				if response != nil && response.Body != nil {
					io.Copy(ioutil.Discard, response.Body)
					response.Body.Close()
				} else if err != nil {
					log.Println(logRequestInterrupted, err)
				}
				atomic.AddInt32(&this.activeRequests, -1)
			}()
		}
	}
}
func countBatches(itemCount int) int {
	remainder := itemCount % maxMetricsPerBatch
	if remainder == 0 {
		return itemCount / maxMetricsPerBatch
	} else {
		return itemCount/maxMetricsPerBatch + 1
	}
}
func (this *Librato) serializeNext() io.Reader {
	written, counterIndex, gaugeIndex := 0, 0, 0
	body := bytes.NewBuffer([]byte{})

	if len(this.hostname) > 0 {
		fmt.Fprintf(body, "source=%s&", this.hostname)
	}

	for index, metric := range this.buffer {
		unixTime := metric.Captured.Unix()
		if metric.MetricType == counterMetricType {
			fmt.Fprintf(body, counterFormat,
				counterIndex, metric.Name, counterIndex, metric.Value, counterIndex, unixTime)
			counterIndex++
		} else {
			fmt.Fprintf(body, gaugeFormat,
				gaugeIndex, metric.Name, gaugeIndex, metric.Value, gaugeIndex, unixTime)
			gaugeIndex++
		}

		delete(this.buffer, index)
		written++
		if written >= maxMetricsPerBatch {
			break
		}
	}

	return body
}
func (this *Librato) buildRequest(body io.Reader) *http.Request {
	request, _ := http.NewRequest("POST", "https://metrics-api.librato.com/v1/metrics", body)
	request.SetBasicAuth(this.email, this.key)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	sendBlankUserAgent(request)
	return request
}
func sendBlankUserAgent(request *http.Request) { request.Header.Set("User-Agent", "") }

const (
	maxMetricsPerBatch    = 256
	counterFormat         = "counters[%d][name]=%s&counters[%d][value]=%d&counters[%d][measure_time]=%d&"
	gaugeFormat           = "gauges[%d][name]=%s&gauges[%d][value]=%d&gauges[%d][measure_time]=%d&"
	logRequestInterrupted = "[WARN] (Metrics) Unable to complete HTTP request:"
	logSkippingPublish    = "[INFO] (Metrics) Skipping publish. No open lanes. (Active/Max: %d/%d)\n"
	logTruncatingPublish  = "[INFO] (Metrics) Truncating publish. Not enough lanes to fully publish entire request. (Needed/Available: %d/%d)\n"
)
