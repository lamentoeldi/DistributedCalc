package config

import "github.com/spf13/viper"

type Config struct {
	Host string
	Port int
}

func NewConfig() *Config {
	viper.AutomaticEnv()

	host := viper.GetString("host")
	port := viper.GetInt("port")

	if host == "" {
		host = "0.0.0.0"
	}

	if port == 0 {
		port = 8080
	}

	return &Config{
		Host: host,
		Port: port,
	}
}
