package config

import (
	"os"
)

type Config struct {
	Port                string
	UserAPIHost         string
	UserAPIPort         string
	AddressAPIHost      string
	AddressAPIPort      string
	ConversationAPIHost string
	ConversationAPIPort string
	BOTKIT_URL          string
	TWILIO_SID          string
}

var Env *Config

func Load() *Config {
	return &Config{
		Port:                getEnv("PORT", "8080"),
		UserAPIHost:         getEnv("USER_API_HOST", "user-api"),
		UserAPIPort:         getEnv("USER_API_PORT", "8080"),
		AddressAPIHost:      getEnv("ADDRESS_API_HOST", "address-api"),
		AddressAPIPort:      getEnv("ADDRESS_API_PORT", "8081"),
		ConversationAPIHost: getEnv("CONVERSATION_API_HOST", "conversation-api"),
		ConversationAPIPort: getEnv("CONVERSATION_API_PORT", "8082"),
		BOTKIT_URL:          getEnv("BOTKIT_URL", "http://fluxo:3000/api/messages"),
		TWILIO_SID:          getEnv("TWILIO_SID", "XXXXXXXX"),
	}
}

// Tenta pegar as variaveis de ambiente, se nao encontrar cai no callback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
