package opentelCore

type Config struct {
	OpenTel struct {
		URL         string  `yaml:"url"`
		MetricPort  *string `yaml:"metricPort"`
		MetricPath  *string `yaml:"metricPath"`
		ServiceName string  `yaml:"serviceName"`
	} `yaml:"openTel"`
}
