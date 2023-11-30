package config

import (
	"errors"
	"fmt"

	"getBlock/internal/environment"
	"github.com/jessevdk/go-flags"
)

//nolint:lll
type (
	// AppConfig contains full configuration of the service.
	AppConfig struct {
		Env       environment.Env `default:"local" description:"Environment application is running in" env:"ENV" long:"env"`
		Positions int             `description:"Positions in top" default:"5" env:"POSITIONS" long:"positions"`
		Logger    Logger          `env-namespace:"LOGGER"         group:"Logger options"  namespace:"logger"`
		Blockio   Blockio         `env-namespace:"BLOCKIO"         group:"Blockio options"  namespace:"blockio"`
	}

	// Logger contains logger configuration.
	Logger struct {
		Level string `description:"Log level to use; environment-base level is used when empty" env:"LEVEL" long:"level"`
	}

	// Blockio contains blockio configuration.
	Blockio struct {
		Token  string `description:"getblock.io token" env:"TOKEN" long:"token" required:"true"`
		Blocks int64  `default:"100" description:"Show activity count after blocks" env:"BLOCKS" long:"blocks"`
	}
)

// ErrHelp is returned when --help flag is
// used and application should not launch.
var ErrHelp = errors.New("help")

// New reads flags and envs and returns AppConfig
// that corresponds to the values read.
func New() (*AppConfig, error) {
	var config AppConfig
	if _, err := flags.Parse(&config); err != nil {
		var flagsErr *flags.Error
		if errors.As(err, &flagsErr) && flagsErr.Type == flags.ErrHelp {
			return nil, ErrHelp
		}
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}
