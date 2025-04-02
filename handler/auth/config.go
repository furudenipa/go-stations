package auth

import "os"

type Config struct {
	UserID   string
	Password string
}

func NewConfigFromEnv() *Config {
	userID := os.Getenv("BASIC_AUTH_USER_ID")
	if userID == "" {
		userID = "test"
	}
	pass := os.Getenv("BASIC_AUTH_PASSWORD")
	if pass == "" {
		pass = "test"
	}
	return &Config{
		UserID:   userID,
		Password: pass,
	}
}
