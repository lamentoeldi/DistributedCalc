package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Host string
	Port int

	AdditionTime       time.Duration
	SubtractionTime    time.Duration
	MultiplicationTime time.Duration
	DivisionTime       time.Duration
}

func NewConfig() *Config {
	viper.AutomaticEnv()

	host := viper.GetString("host")
	port := viper.GetInt("port")

	additionTime := viper.GetInt("time_addition_ms")
	subtractionTime := viper.GetInt("time_subtraction_ms")
	multiplicationTime := viper.GetInt("time_multiplication_ms")
	divisionTime := viper.GetInt("time_division_ms")

	if host == "" {
		host = "0.0.0.0"
	}

	if port == 0 {
		port = 8080
	}

	if additionTime == 0 {
		additionTime = 1
	}

	if subtractionTime == 0 {
		subtractionTime = 1
	}

	if multiplicationTime == 0 {
		multiplicationTime = 1
	}

	if divisionTime == 0 {
		divisionTime = 1
	}

	return &Config{
		Host:               host,
		Port:               port,
		AdditionTime:       time.Duration(additionTime) * time.Millisecond,
		SubtractionTime:    time.Duration(subtractionTime) * time.Millisecond,
		MultiplicationTime: time.Duration(multiplicationTime) * time.Millisecond,
		DivisionTime:       time.Duration(divisionTime) * time.Millisecond,
	}
}
