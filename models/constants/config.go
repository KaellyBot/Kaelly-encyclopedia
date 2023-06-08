package constants

import "github.com/rs/zerolog"

const (
	ConfigFileName = ".env"

	// RabbitMQ address.
	RabbitMQAddress = "RABBITMQ_ADDRESS"

	// Timeout to retrieve Dofus data in seconds.
	DofusDudeTimeout = "HTTP_TIMEOUT"

	// Metric port.
	MetricPort = "METRIC_PORT"

	// Zerolog values from [trace, debug, info, warn, error, fatal, panic].
	LogLevel = "LOG_LEVEL"

	// Boolean; used to register commands at development guild level or globally.
	Production = "PRODUCTION"

	defaultRabbitMQAddress  = "amqp://localhost:5672"
	defaultDofusDudeTimeout = 60
	defaultMetricPort       = 2112
	defaultLogLevel         = zerolog.InfoLevel
	defaultProduction       = false
)

func GetDefaultConfigValues() map[string]any {
	return map[string]any{
		RabbitMQAddress:  defaultRabbitMQAddress,
		DofusDudeTimeout: defaultDofusDudeTimeout,
		MetricPort:       defaultMetricPort,
		LogLevel:         defaultLogLevel.String(),
		Production:       defaultProduction,
	}
}
