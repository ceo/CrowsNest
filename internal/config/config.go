package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// CliBinOpts holds information about user provided binaries to be executed.
type CliBinOpts struct {
	BinaryPath string `validate:"omitempty,file"`
	Flags      []string
	// NOTE: It seems it's not possible to validate as dir when using required_with.
	WorkingDirectory string `validate:"required_with=BinaryPath"`
}

// RepositoryConfig holds information about a users local Git repo.
type RepositoryConfig struct {
	Directory    string `validate:"required,dir"`
	GitPullFlags []string
	Interval     int `validate:"gt=0"`
	PrePullCmd   CliBinOpts
	PostPullCmd  CliBinOpts
}

// Config holds crowsnest's configuration.
type Config struct {
	Repositories map[string]*RepositoryConfig
}

// Get get's the current Config.
func Get(path string) (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")

	if path != "" {
		viper.AddConfigPath(path)
	}

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("fatal error config file: %w", err)
	}

	var c Config

	if err := viper.Unmarshal(&c); err != nil {
		return Config{}, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	if len(c.Repositories) == 0 {
		return Config{}, errors.New("no repositoires found")
	}

	validate := validator.New()
	for cfgName, cfg := range c.Repositories {
		if cfg.Interval == 0 {
			cfg.Interval = 60
		}

		errs := validate.Struct(cfg)

		if errs != nil {
			// NOTE: Basically just return the first err, the user can use some trial and error if need be.
			for _, err := range errs.(validator.ValidationErrors) {
				return Config{}, fmt.Errorf(
					"invalid configuration for %s: %s",
					cfgName,
					strings.ToLower(err.Namespace()),
				)
			}
		}

	}

	return c, nil
}
