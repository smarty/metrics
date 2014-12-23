package metrics

type Backend interface {
	Listen(chan []Measurement)
}
