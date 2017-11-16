package metrics

type LibratoConfigLoader func() LibratoConfig

type LibratoConfig struct {
	Email string
	Key   string
}
