package config

import (
	"errors"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Server ServerConfig
		DB     DBConfig
		Auth   AuthConfig
		Logger Logger
	}

	ServerConfig struct {
		Debug            bool
		Addr             string
		JwtSecret        string `mapstructure:"jwt_secret"`
		JwtRefreshSecret string `mapstructure:"jwt_refresh_secret"`
	}

	DBConfig struct {
		Driver   string
		Host     string
		Port     int
		User     string
		Password string
		Name     string
		SSL      string
	}

	TokenConfig struct {
		MaxAge   int `mapstructure:"max_age"`
		Secure   bool
		HttpOnly bool `mapstructure:"http_only"`
	}

	AuthConfig struct {
		AccessToken  TokenConfig `mapstructure:"access_token"`
		RefreshToken TokenConfig `mapstructure:"refresh_token"`
	}

	Logger struct {
		Level string
	}
)

func LoadConfig(path string, name string) (*Config, error) {
	config := new(Config)

	viper.SetConfigName(name)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}

		return nil, err
	}

	return config, viper.Unmarshal(config)
}
