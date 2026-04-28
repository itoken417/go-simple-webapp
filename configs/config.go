package config

import "os"

type Config struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
	SMTPTo       string
}

var config *Config

func Get() *Config {
	if config == nil {
		config = &Config{
			SMTPHost:     os.Getenv("SMTP_HOST"),
			SMTPPort:     os.Getenv("SMTP_PORT"),
			SMTPUser:     os.Getenv("SMTP_USER"),
			SMTPPassword: os.Getenv("SMTP_PASSWORD"),
			SMTPFrom:     os.Getenv("SMTP_FROM"),
			SMTPTo:       os.Getenv("SMTP_TO"),
		}
	}
	return config
}
