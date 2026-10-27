package config

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/frostgrove/vv/app/module"
	"github.com/frostgrove/vv/utils/vvcfg"
	"github.com/frostgrove/vv/vvdb"
)

const DefaultPath = "./config/app.yml"

type Config struct {
	App           App           `yaml:"app"`
	Transports    Transports    `yaml:"transports"`
	Infra         Infra         `yaml:"infra"`
	Clients       Clients       `yaml:"clients"`
	Internal      Internal      `yaml:"internal"`
	Observability Observability `yaml:"observability"`
}

// Load reads the configuration strictly: a key no field declares refuses the
// start instead of silently falling back to a default (D-086). The file is
// --config-path, then CONFIG_PATH, then DefaultPath.
func Load() (*Config, *vvcfg.Report, error) {
	return vvcfg.LoadFrom[Config](vvcfg.Source{
		Arguments:   os.Args[1:],
		DefaultPath: DefaultPath,
		Strict:      true,
	})
}

func MustLoad() (*Config, *vvcfg.Report) {
	configuration, report, err := Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "configuration: %v\n", err)
		os.Exit(2)
	}
	return configuration, report
}

type App struct {
	Name string `yaml:"name" env:"APP_NAME" env-default:"vv-template"`
	// Env selects which seeders run (see app.Seeder.Envs).
	Env string `yaml:"env" env:"APP_ENV" env-default:"dev"`
	// Profile decides which roles of every module this process wires:
	// complete (api + worker + seeder), serving (api), working (worker).
	Profile string `yaml:"profile" env:"APP_PROFILE" env-default:"complete"`
}

var Envs = []string{"dev", "test", "staging", "prod"}

func (this *App) ValidateSelf() error {
	var problems []error
	if strings.TrimSpace(this.Name) == "" {
		problems = append(problems, errors.New("name is required"))
	}
	if !slices.Contains(Envs, this.Env) {
		problems = append(problems, fmt.Errorf("env %q is not one of %v", this.Env, Envs))
	}
	if _, err := this.ModuleProfile(); err != nil {
		problems = append(problems, err)
	}
	return errors.Join(problems...)
}

func (this *App) ModuleProfile() (module.Profile, error) {
	switch this.Profile {
	case "complete":
		return module.Complete, nil
	case "serving":
		return module.Serving, nil
	case "working":
		return module.Working, nil
	}
	return module.Profile{}, fmt.Errorf("profile %q is not one of complete, serving, working", this.Profile)
}

// Transport-level configs: anything that handles communication.
type Transports struct {
	Fiber FiberConfig `yaml:"fiber"`
}

type FiberConfig struct {
	Addr         string   `yaml:"addr" env:"HTTP_ADDR" env-default:":8009"`
	Prefix       string   `yaml:"prefix" env:"HTTP_PREFIX" env-default:"/api"`
	AllowOrigins []string `yaml:"allow_origins" env:"HTTP_ALLOW_ORIGINS"`
	// APIKeys is the development authenticator. Replace it with a key store or
	// JWT (auth/authjwt, auth/access) before anything leaves your machine.
	APIKeys []APIKey `yaml:"api_keys"`
}

type APIKey struct {
	Key         string   `yaml:"key"`
	Subject     string   `yaml:"subject"`
	Permissions []string `yaml:"permissions"`
}

func (this *FiberConfig) ValidateSelf() error {
	var problems []error
	if strings.TrimSpace(this.Addr) == "" {
		problems = append(problems, errors.New("addr is required"))
	}
	if this.Prefix != "" && (!strings.HasPrefix(this.Prefix, "/") || strings.HasSuffix(this.Prefix, "/")) {
		problems = append(problems, fmt.Errorf("prefix %q must start with / and not end with one", this.Prefix))
	}
	for index, key := range this.APIKeys {
		if strings.TrimSpace(key.Key) == "" || strings.TrimSpace(key.Subject) == "" {
			problems = append(problems, fmt.Errorf("api_keys[%d] needs a key and a subject", index))
		}
	}
	return errors.Join(problems...)
}

type Infra struct {
	Db vvdb.Config `yaml:"db"`
}

// Clients to external APIs, SDKs.
type Clients struct {
	// Example: PaymentAPI, NotificationService, etc.
}

// Internal microservices or internal modules.
type Internal struct {
	// Example: UserServiceConfig, EmailServiceConfig, etc.
}

type Observability struct {
	Logger LoggerConfig `yaml:"logger"`
	Otel   OtelConfig   `yaml:"otel"`
}

type LoggerConfig struct {
	Level  string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	Format string `yaml:"format" env:"LOG_FORMAT" env-default:"text"`
}

func (this *LoggerConfig) ValidateSelf() error {
	var problems []error
	if !slices.Contains([]string{"debug", "info", "warn", "error"}, this.Level) {
		problems = append(problems, fmt.Errorf("level %q is not one of debug, info, warn, error", this.Level))
	}
	if !slices.Contains([]string{"text", "json"}, this.Format) {
		problems = append(problems, fmt.Errorf("format %q is not one of text, json", this.Format))
	}
	return errors.Join(problems...)
}

// OtelConfig turns on the OpenTelemetry SDK with OTLP/gRPC export of traces and
// metrics. Off, nothing is exported and no SDK is started. The standard
// OTEL_EXPORTER_OTLP_* variables apply when Endpoint is empty.
type OtelConfig struct {
	Enabled     bool    `yaml:"enabled" env:"OTEL_ENABLED" env-default:"false"`
	Endpoint    string  `yaml:"endpoint" env:"OTEL_ENDPOINT"`
	Insecure    bool    `yaml:"insecure" env:"OTEL_INSECURE" env-default:"false"`
	SampleRatio float64 `yaml:"sample_ratio" env:"OTEL_SAMPLE_RATIO" env-default:"1"`
}

func (this *OtelConfig) ValidateSelf() error {
	if this.SampleRatio < 0 || this.SampleRatio > 1 {
		return fmt.Errorf("sample_ratio %v is not between 0 and 1", this.SampleRatio)
	}
	return nil
}
