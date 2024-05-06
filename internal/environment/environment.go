package environment

import (
	"context"
)

type ApiConfig struct {
	Port       int    `env:"PORT, default=8080"`
	EnvName    string `env:"ENV_NAME, default=development"`
	ApiVersion string `env:"VERSION, default=v1"`
}

func (c *ApiConfig) IsDevelopment() bool {
	return c.EnvName == "development"
}

type DatabaseConfig struct {
	Url string `env:"URL, required"`
}

type CloudConfig struct {
	OrderProductionTopic string `env:"ORDER_PRODUCTION_TOPIC_NAME, required"`
	UpdateOrderTopic     string `env:"UPDATE_ORDER_TOPIC_NAME, required"`
	OrderPaymentQueue    string `env:"ORDER_PAYMENT_QUEUE_NAME, required"`

	BaseEndpoint string `env:"BASE_ENDPOINT"`
}

func (c *CloudConfig) IsBaseEndpointSet() bool {
	return c.BaseEndpoint != ""
}

type Config struct {
	ApiConfig   *ApiConfig      `env:",prefix=API_"`
	DbConfig    *DatabaseConfig `env:",prefix=DB_"`
	CloudConfig *CloudConfig    `env:",prefix=AWS_"`
}

type Environment interface {
	GetEnvironmentFromFile(ctx context.Context, fileName string) (*Config, error)
	GetEnvironment(ctx context.Context) (*Config, error)
}
