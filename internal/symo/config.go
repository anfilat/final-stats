package symo

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func NewConfig(configFile string) (Config, error) {
	config := Config{}

	v := viper.New()

	configure(v)

	if configFile != "" {
		v.SetConfigFile(configFile)
		err := v.ReadInConfig()
		if err != nil {
			return config, fmt.Errorf("failed to read configuration: %w", err)
		}
	}

	if err := v.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	if err := config.Validate(); err != nil {
		return config, fmt.Errorf("failed to validate configuration: %w", err)
	}

	return config, nil
}

func configure(v *viper.Viper) {
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	v.SetDefault("log.level", "INFO")
	v.SetDefault("server.port", "8000")
	v.SetDefault("metric.loadavg", true)
}

type Config struct {
	Logger LoggerConf
	Server ServerConf
	Metric MetricConf
}

func (c Config) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return err
	}

	return nil
}

type LoggerConf struct {
	Level string
}

type ServerConf struct {
	Port string
}

func (c ServerConf) Validate() error {
	if c.Port == "" {
		return errors.New("server port is required")
	}

	return nil
}

type MetricConf struct {
	Loadavg bool
}
