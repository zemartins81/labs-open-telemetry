package config

import "os"

type Config struct {
	Port        string
	ServiceBURL string
}

func LoadConfig() *Config {
	return &Config{
		Port:        os.Getenv("PORT"),
		ServiceBURL: os.Getenv("SERVICE_B_URL"),
	}
}
