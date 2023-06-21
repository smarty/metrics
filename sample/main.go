package main

import (
	"log"
	"net/http"
	"time"

	"github.com/smarty/metrics/v2"
)

func main() {
	counter := metrics.NewCounter("my_counter",
		metrics.Options.Description("this is a description"),
		metrics.Options.Label("label_key", "label_value"),
	)

	exporter := metrics.NewExporter()
	exporter.Add(counter)

	go func() {
		for {
			counter.Increment()
			time.Sleep(time.Millisecond * 10)
		}
	}()

	server := &http.Server{Addr: "127.0.0.1:8080", Handler: exporter}
	log.Printf("[INFO] Listening for HTTP traffic on [%s]", server.Addr)
	_ = server.ListenAndServe()
}
