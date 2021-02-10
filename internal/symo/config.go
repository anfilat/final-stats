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

	v.SetDefault("app.maxSeconds", 600)
	v.SetDefault("log.level", "INFO")
	v.SetDefault("server.port", "8000")
	v.SetDefault("metric.loadavg", true)
	v.SetDefault("metric.cpu", true)
	v.SetDefault("metric.loaddisks", true)
	v.SetDefault("metric.usedfs", true)
}

type Config struct {
	App    AppConf
	Log    LoggerConf
	Server ServerConf
	Metric MetricConf
}

func (c Config) Validate() error {
	if err := c.App.Validate(); err != nil {
		return err
	}
	if err := c.Server.Validate(); err != nil {
		return err
	}

	return nil
}

type AppConf struct {
	MaxSeconds int
}

func (c AppConf) Validate() error {
	if c.MaxSeconds <= 0 {
		return errors.New("time to keep metrics must be greater than zero")
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
	Loadavg   bool
	CPU       bool
	Loaddisks bool
	UsedFS    bool
}
