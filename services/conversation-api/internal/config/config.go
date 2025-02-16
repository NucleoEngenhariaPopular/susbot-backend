package config

import "os"

type Config struct {
	Port            string
	MongoURI        string
	MongoDBName     string
	MongoCollection string
}

func Load() *Config {
	return &Config{
		Port:            getEnv("CONVERSATION_API_PORT", "8082"),
		MongoURI:        getEnv("MONGO_URI", "mongodb://root:example@mongo:27017/"),
		MongoDBName:     getEnv("MONGODB_NAME", "my_database"),
		MongoCollection: getEnv("MONGODB_COLLECTION", "conversations"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
