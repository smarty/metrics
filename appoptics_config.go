package metrics

type AppOpticsConfigLoader func() AppOpticsConfig

type AppOpticsConfig struct {
	Key string
}
