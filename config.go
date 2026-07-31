package main

import "time"

type Config struct {
	Port                int
	Backends            []string
	HealthCheckInterval time.Duration
}

func NewConfig() *Config {
	return &Config{
		Port:                8080,
		Backends:            []string{"127.0.0.1:8081", "127.0.0.1:8082"},
		HealthCheckInterval: 10 * time.Second,
	}
}
