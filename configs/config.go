package config

import (
	"fmt"
	"os"

	"github.com/itoken417/go-simple-webapp/internal/secret"
)

type Config struct {
	AppName      string
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
	SystemTo     string // サービスメールのデフォルト宛先
	ErrorTo      string // エラー通知の宛先
	MailEncoding string
}

var config *Config

func Get() (*Config, error) {
	if config != nil {
		return config, nil
	}
	smtpPassword, err := decryptSMTPPassword()
	if err != nil {
		return nil, fmt.Errorf("SMTP パスワードの復号に失敗: %w", err)
	}
	config = &Config{
		AppName:      os.Getenv("APP_NAME"),
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPUser:     os.Getenv("SMTP_USER"),
		SMTPPassword: smtpPassword,
		SMTPFrom:     os.Getenv("SMTP_FROM"),
		SystemTo:     os.Getenv("SYSTEM_TO"),
		ErrorTo:      os.Getenv("ERROR_TO"),
		MailEncoding: os.Getenv("MAIL_ENCODING"),
	}
	return config, nil
}

func decryptSMTPPassword() (string, error) {
	encrypted := os.Getenv("SMTP_PASSWORD_ENCRYPTED")
	salt := os.Getenv("SMTP_PASSWORD_SALT")
	if encrypted == "" || salt == "" {
		return "", fmt.Errorf("SMTP_PASSWORD_ENCRYPTED または SMTP_PASSWORD_SALT が未設定です")
	}
	key, err := secret.LoadKey()
	if err != nil {
		return "", err
	}
	return secret.Decrypt(encrypted, salt, key)
}
