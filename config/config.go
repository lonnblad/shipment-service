package config

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	defaultServiceName     = "shipment-service"
	defaultServiceVersion  = "dev"
	defaultRestPort        = "8080"
	defaultRestURL         = "http://localhost:" + defaultRestPort
	defaultShutdownTimeout = 20 * time.Second

	configKeyEnvironment    = "environment"
	configKeyServiceName    = "service-name"
	configKeyServiceVersion = "service-version"
	configKeyRestPort       = "rest-port"
	configKeyRestURL        = "rest-url"
	configKeyShutdownTimout = "shutdown-timeout"
)

func init() {
	viper.AutomaticEnv()

	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)

	if viper.GetString(configKeyEnvironment) == "" {
		viper.SetDefault(configKeyEnvironment, EnvironmentLocal.String())
	}

	if viper.GetString(configKeyServiceName) == "" {
		viper.SetDefault(configKeyServiceName, defaultServiceName)
	}

	if viper.GetString(configKeyServiceVersion) == "" {
		viper.SetDefault(configKeyServiceVersion, defaultServiceVersion)
	}

	if viper.GetString(configKeyRestPort) == "" {
		viper.SetDefault(configKeyRestPort, defaultRestPort)
	}

	if viper.GetString(configKeyRestURL) == "" {
		viper.SetDefault(configKeyRestURL, defaultRestURL)
	}

	if viper.GetDuration(configKeyShutdownTimout) == 0 {
		viper.SetDefault(configKeyShutdownTimout, defaultShutdownTimeout)
	}
}

func mustGetString(key string) string {
	value := viper.GetString(key)
	if value == "" {
		panic(fmt.Sprintf("%q is not set", key))
	}

	return value
}

func GetEnvironment() Environment {
	return Environment(mustGetString(configKeyEnvironment))
}

func GetServiceName() string {
	return mustGetString(configKeyServiceName)
}

func GetServiceVersion() string {
	return strings.ToLower(viper.GetString(configKeyServiceVersion))
}

func GetRestPort() string {
	return mustGetString(configKeyRestPort)
}

func GetRestURL() url.URL {
	url, err := url.Parse(viper.GetString(configKeyRestURL))
	if err != nil {
		panic(err)
	}

	return *url
}

func GetShutdownTimeout() time.Duration {
	return viper.GetDuration(configKeyShutdownTimout)
}
