package config

var (
	CollectorEndpoint = "localhost:4317"
	ServiceName       = "tyke-go-service"
	Insecure          = "true"
	ServiceIdentifier = "tyke-go-service"
	AppName           = "tyke-app"
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

func SetServiceIdentifier(identifier string) {
	ServiceIdentifier = identifier
}

func SetAppName(name string) {
	AppName = name
}
