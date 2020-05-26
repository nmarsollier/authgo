package env

import (
	"os"
	"strconv"
)

// Configuration properties
type Configuration struct {
	Port      int    `json:"port"`
	RabbitURL string `json:"rabbitUrl"`
	MongoURL  string `json:"mongoUrl"`
	WWWWPath  string `json:"wwwPath"`
	JWTSecret string `json:"jwtSecret"`
}

var config *Configuration

func new() *Configuration {
	return &Configuration{
		Port:      3000,
		RabbitURL: "amqp://localhost",
		MongoURL:  "mongodb://localhost:27017",
		WWWWPath:  "www",
		JWTSecret: "ecb6d3479ac3823f1da7f314d871989b",
	}
}

// Get Obtiene las variables de entorno del sistema
func Get() *Configuration {
	if config == nil {
		config = load()
	}

	return config
}

// Load file properties
func load() *Configuration {
	result := new()

	if value := os.Getenv("RABBIT_URL"); len(value) > 0 {
		result.RabbitURL = value
	}

	if value := os.Getenv("MONGO_URL"); len(value) > 0 {
		result.MongoURL = value
	}

	if value := os.Getenv("PORT"); len(value) > 0 {
		if intVal, err := strconv.Atoi(value); err != nil {
			result.Port = intVal
		}
	}

	if value := os.Getenv("WWW_PATH"); len(value) > 0 {
		result.WWWWPath = value
	}

	if value := os.Getenv("JWT_SECRET"); len(value) > 0 {
		result.JWTSecret = value
	}

	return result
}
