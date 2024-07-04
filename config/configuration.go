package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

var Configuration AppConfig

type AppConfig struct {
	Microservices Microservices `yaml:"microservices"`
	Environment   Environment   `mapstructure:"environment"`
	Log           Log           `mapstructure:"log"`
}

type Logger struct {
	Level string `mapstructure:"level" yaml:"level"`
}

type Microservices struct {
	Name    string `mapstructure:"name" yaml:"name"`
	Port    string `mapstructure:"port" yaml:"port"`
	Version string `mapstructure:"version" yaml:"version"`
}

type Environment struct {
	Value string `mapstructure:"value" yaml:"value"`
}

type Log struct {
	Level string `mapstructure:"level" yaml:"level"`
}

func LoadConfigurationMicroservice(path string) {
	fmt.Println("Loading configuration microservice from file [application.yml]", path)
	viper.SetConfigName("application.yml")
	viper.AddConfigPath(path)
	viper.SetConfigType("yml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading configuration filem, %s", err)
		panic("Error reading configuration file")
	}

	// Set undefined variables:
	hostname, _ := os.Hostname()

	fmt.Println("hostname: ", hostname)
	viper.SetDefault("microserviceServer", hostname)
	viper.SetDefault("microservicePathRoot", "./")

	if err := viper.Unmarshal(&Configuration); err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
		panic(fmt.Sprintf("unable to decode into struct, %v", err))
	}

	fmt.Println("Configuration loaded: ", Configuration.Microservices.Name)
}
