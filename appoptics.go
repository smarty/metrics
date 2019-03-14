package metrics

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

type AppOptics struct {
	config         AppOpticsConfigLoader
	hostname       string
	maxRequests    int32
	activeRequests int32
	buffer         map[int]MetricMeasurement
	client         *http.Client
}

func newAppOptics(config AppOpticsConfigLoader, hostname string, maxRequests int32) *AppOptics {
	// TODO: validate inputs

	client := &http.Client{
		Transport: &http.Transport{DisableCompression: true},
		Timeout:   time.Duration(time.Second * 10),
	}

	return &AppOptics{
		config:      config,
		hostname:    hostname,
		maxRequests: maxRequests,
		buffer:      map[int]MetricMeasurement{},
		client:      client,
	}
}

func (this *AppOptics) Listen(queue chan []MetricMeasurement) {
	for measurements := range queue {
		for _, measurement := range measurements {
			this.buffer[measurement.ID] = measurement // last one wins
		}
		if len(queue) == 0 {
			this.publish()
		}
	}
}

func (this *AppOptics) publish() {
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
					if response.StatusCode >= 200 && response.StatusCode < 300 {
						io.Copy(ioutil.Discard, response.Body)
						response.Body.Close()
					} else {
						buf := new(bytes.Buffer)
						buf.ReadFrom(response.Body)
						newStr := buf.String()
						log.Println("AppOptics:", response.StatusCode, newStr)
					}
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

type Measurement struct {
	Name  string            `json:"name"`
	Value int64             `json:"value"`
	Time  int64             `json:"time"`
	Tags  map[string]string `json:"tags"`
}
type Measurements struct {
	Measurements []Measurement `json:"measurements"`
}

func (this *AppOptics) serializeNext() io.Reader {
	written := 0
	var measurements Measurements

	for index, metric := range this.buffer {
		unixTime := metric.Captured.Unix()
		measurements.Measurements = append(measurements.Measurements, Measurement{
			Name:  metric.Name,
			Value: metric.Value,
			Time:  unixTime,
			Tags:  this.buildTags(metric),
		})

		delete(this.buffer, index)
		written++
		if written >= maxMetricsPerBatch {
			break
		}
	}
	jsonBody, err := json.Marshal(measurements)
	if err != nil {
		log.Println("[ERROR] Unable to marshal measurements into json:", err)
	}

	return bytes.NewReader(jsonBody)
}

func (this *AppOptics) buildTags(metric MetricMeasurement) map[string]string {
	tags := make(map[string]string)
	for key, value := range metric.Tags {
		// filter AppOptics invalid tags
		if key != "" && value != "" {
			tags[key] = value
		}
	}
	tags["hostname"] = this.hostname
	tags["metrictype"] = intToWord(metric.MetricType)
	return tags
}
func intToWord(metricType int) string {
	switch metricType {
	case 1:
		return "counter"
	case 2:
		return "gauge"
	}
	return "unknown(" + strconv.Itoa(metricType) + ")"
}
func (this *AppOptics) buildRequest(body io.Reader) *http.Request {
	request, err := http.NewRequest("POST", "https://api.appoptics.com/v1/measurements", body)
	if err != nil {
		log.Println("[ERROR] Unable to creating new http request:", err)
	}

	config := this.config()
	request.SetBasicAuth(config.Key, "") // AppOptics only requires key as username, password blank
	request.Header.Set("Content-Type", "application/json")
	sendBlankUserAgent(request)
	return request
}
func sendBlankUserAgent(request *http.Request) { request.Header.Set("User-Agent", "") }

const (
	maxMetricsPerBatch    = 256
	logRequestInterrupted = "[WARN] (Metrics) Unable to complete HTTP request:"
	logSkippingPublish    = "[INFO] (Metrics) Skipping publish. No open lanes. (Active/Max: %d/%d)\n"
	logTruncatingPublish  = "[INFO] (Metrics) Truncating publish. Not enough lanes to fully publish entire request. (Needed/Available: %d/%d)\n"
)
