package config

import "github.com/frostgrove/vv/utils/vvdb"

type Config struct {
	App           App           `yaml:"app"`
	Transports    Transports    `yaml:"transports"`
	Infra         Infra         `yaml:"infra"`
	Clients       Clients       `yaml:"clients"`
	Internal      Internal      `yaml:"internal"`
	Observability Observability `yaml:"observability"`
}

// App
type App struct {
	Name string `yaml:"name" env-default:"RX-TEMPLATE"`
}

// Transport-level configs: anything that handles communication
type Transports struct {
	Fiber FiberConfig `yaml:"fiber"`
}

type FiberConfig struct {
	Port         int      `yaml:"port" env-required:"true"`
	AllowOrigins []string `yaml:"allow_origins" env-required:"true"`
	Prefork      bool     `yaml:"prefork" default:"false"`
}

// Infra

type Infra struct {
	Db vvdb.Config `yaml:"db"`
}

// Clients to external APIs, SDKs
type Clients struct {
	// Example: PaymentAPI, NotificationService, etc.
}

// Internal microservices or internal modules
type Internal struct {
	// Example: UserServiceConfig, EmailServiceConfig, etc.
}

type Observability struct {
	Logger LoggerConfig `yaml:"logger"`
}

type LoggerConfig struct {
	Name string `yaml:"name" env-default:"app-logger"`
}
