package metrics

import (
	"bytes"
	"fmt"
	"net/http"
	"sync/atomic"
)

type Librato struct {
	email          string
	key            string
	hostname       string
	maxRequests    int32
	activeRequests int32
	buffer         map[int]Measurement
	client         *http.Client
}

func WithLibrato(email, key, hostname string, maxRequests int32) *Librato {
	// TODO: validate inputs

	// TODO: all HTTP-related timeouts
	transport := &http.Transport{DisableCompression: true}
	client := &http.Client{Transport: transport}

	return &Librato{
		email:       email,
		key:         key,
		hostname:    hostname,
		maxRequests: maxRequests,
		buffer:      map[int]Measurement{},
		client:      client,
	}
}

func (this *Librato) Listen(queue chan []Measurement) {
	for measurements := range queue {
		for _, measurement := range measurements {
			this.buffer[measurement.ID] = measurement // overwrite the oldest with the newest
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
			break // all lanes are busy
		}

		// how many requests would we need to take care of all of the items
		required := int32(countBatches(len(this.buffer)))
		if required > available {
			required = available // not enough open/available lanes to accomodate all the requests;
		}

		for i := int32(0); i < required; i++ {
			body := this.serializeNext()
			request := this.buildRequest(body)
			atomic.AddInt32(&this.activeRequests, 1)

			go func() {
				this.client.Do(request) // ignore errors
				atomic.AddInt32(&this.activeRequests, -1)
			}()
		}
	}
}

func (this *Librato) buildRequest(body *bytes.Buffer) *http.Request {
	request, _ := http.NewRequest("POST", "https://metrics-api.librato.com/v1/metrics", body)
	request.SetBasicAuth(this.email, this.key)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return request
}

func (this *Librato) serializeNext() *bytes.Buffer {
	written, counterIndex, gaugeIndex := 0, 0, 0
	body := bytes.NewBuffer([]byte{})

	if len(this.hostname) > 0 {
		fmt.Fprintf(body, "source=%s&", this.hostname)
	}

	for index, metric := range this.buffer {
		meta := standard.meta[index]
		unixTime := metric.Captured.Unix()
		if meta.MetricType == CounterMetric {
			fmt.Fprintf(body, counterFormat, counterIndex, meta.Name, counterIndex, metric.Value, counterIndex, unixTime)
			counterIndex++
		} else {
			fmt.Fprintf(body, gaugeFormat, gaugeIndex, meta.Name, gaugeIndex, metric.Value, gaugeIndex, unixTime)
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

func countBatches(itemCount int) int {
	remainder := itemCount % maxMetricsPerBatch
	if remainder == 0 {
		return itemCount / maxMetricsPerBatch
	} else {
		return itemCount/maxMetricsPerBatch + 1
	}
}

// const maxMetricsPerBatch = 256
const maxMetricsPerBatch = 3
const counterFormat = "counters[%d][name]=%s&counters[%d][value]=%d&counters[%d][measure_time]=%d&"
const gaugeFormat = "gauges[%d][name]=%s&gauges[%d][value]=%d&gauges[%d][measure_time]=%d&"
