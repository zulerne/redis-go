package config

type Config struct {
	Addr string
}

func DefaultConfig() Config {
	return Config{
		Addr: "0.0.0.0:6379",
	}
}
