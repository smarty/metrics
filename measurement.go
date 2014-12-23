package metrics

import "time"

// Measurement is the struct that is sent onto the outbound channel for
// publishing to whatever backend service that happens to be configured.
type Measurement struct {
	ID       int
	Captured time.Time
	Value    int64
}
