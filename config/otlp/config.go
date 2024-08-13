package otlp

type Config struct {
	ServiceName  string `env:"OTLP_SERVICE_NAME,required"`
	Endpoint     string `env:"OTLP_EXPORTER_ENDPOINT,required"`
	InSecureMode bool   `env:"OTLP_INSECURE_MODE,required"`
}
