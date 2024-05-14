package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	LogLevel string        `env:"LOG_LEVEL" env-default:"DEBUG"`
	Timeout  time.Duration `env:"TIMEOUT" env-required:"true"`
	Postgres PostgresConfig
	Rpc      RpcConfig
	Web      WebConfig
}

func InitConfig() (*AppConfig, error) {
	var cfg AppConfig

	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

type PostgresConfig struct {
	URL string `env:"POSTGRES_URL" env-required:"true"`
}

type RpcConfig struct {
	Host string `env:"RPC_HOST" env-required:"true"`
	Port string `env:"RPC_PORT" env-required:"true"`
}

func (rpc *RpcConfig) Address() string {
	return fmt.Sprintf("%s:%s", rpc.Host, rpc.Port)
}

type WebConfig struct {
	Host string `env:"WEB_HOST" env-required:"true"`
	Port string `env:"WEB_PORT" env-required:"true"`
}

func (web *WebConfig) Address() string {
	return fmt.Sprintf("%s:%s", web.Host, web.Port)
}
