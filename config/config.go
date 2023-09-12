package config

var (
	CollectorEndpoint = "localhost:4317"
	ServiceName       = "tyke-go-service"
	Insecure          = "true"
)

func SetCollectorHost(host string) {
	CollectorEndpoint = host
}

func SetServiceName(name string) {
	ServiceName = name
}

func WithInsecure(insecure bool) {
	if !insecure {
		Insecure = "false"
	}
}
