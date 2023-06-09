package constants

import (
	"time"

	"github.com/rs/zerolog"
)

const (
	ConfigFileName = ".env"

	// MySQL URL with the following format: HOST:PORT.
	MySQLURL = "MYSQL_URL"

	// MySQL user.
	MySQLUser = "MYSQL_USER"

	// MySQL password.
	MySQLPassword = "MYSQL_PASSWORD"

	// MySQL database name.
	MySQLDatabase = "MYSQL_DATABASE"

	// RabbitMQ address.
	RabbitMQAddress = "RABBITMQ_ADDRESS"

	// Redis address.
	RedisAddress = "REDIS_ADDRESS"

	// Redis cache retention. Duration type.
	RedisCacheRetention = "REDIS_CACHE_RETENTION"

	// Redis cache size, following LFU rules.
	RedisCacheSize = "REDIS_CACHE_SIZE"

	// Timeout to retrieve Dofus data. Duration type.
	DofusDudeTimeout = "HTTP_TIMEOUT"

	// Metric port.
	MetricPort = "METRIC_PORT"

	// Zerolog values from [trace, debug, info, warn, error, fatal, panic].
	LogLevel = "LOG_LEVEL"

	// Boolean; used to register commands at development guild level or globally.
	Production = "PRODUCTION"

	defaultMySQLURL            = "localhost:3306"
	defaultMySQLUser           = ""
	defaultMySQLPassword       = ""
	defaultMySQLDatabase       = "kaellybot"
	defaultRabbitMQAddress     = "amqp://localhost:5672"
	defaultRedisAddress        = "localhost:6379"
	defaultRedisCacheRetention = 60 * time.Minute
	defaultRedisCacheSize      = 1000
	defaultDofusDudeTimeout    = 10 * time.Second
	defaultMetricPort          = 2112
	defaultLogLevel            = zerolog.InfoLevel
	defaultProduction          = false
)

func GetDefaultConfigValues() map[string]any {
	return map[string]any{
		MySQLURL:            defaultMySQLURL,
		MySQLUser:           defaultMySQLUser,
		MySQLPassword:       defaultMySQLPassword,
		MySQLDatabase:       defaultMySQLDatabase,
		RabbitMQAddress:     defaultRabbitMQAddress,
		RedisAddress:        defaultRedisAddress,
		RedisCacheRetention: defaultRedisCacheRetention,
		RedisCacheSize:      defaultRedisCacheSize,
		DofusDudeTimeout:    defaultDofusDudeTimeout,
		MetricPort:          defaultMetricPort,
		LogLevel:            defaultLogLevel.String(),
		Production:          defaultProduction,
	}
}
