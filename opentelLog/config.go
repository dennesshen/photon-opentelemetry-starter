package opentelLog

type Config struct {
	Log struct {
		LogLevel string `yaml:"level"`
	} `yaml:"log"`
}
