package env

import (
	"os"
	"strconv"
)

// Configuration properties
type Configuration struct {
	Port           int    `json:"port"`
	GqlPort        int    `json:"gqlPort"`
	RabbitURL      string `json:"rabbitUrl"`
	JWTSecret      string `json:"jwtSecret"`
	FluentUrl      string `json:"fluentUrl"`
	AwsAccessKeyId string `json:"AwsAccessKeyId"`
	AwsSecret      string `json:"AwsSecret"`
	AwsRegion      string `json:"AwsRegion"`
}

var config *Configuration

func new() *Configuration {
	return &Configuration{
		Port:           3000,
		GqlPort:        4000,
		RabbitURL:      "amqp://localhost",
		JWTSecret:      "ecb6d3479ac3823f1da7f314d871989b",
		FluentUrl:      "localhost:24224",
		AwsAccessKeyId: "placeholder_value",
		AwsSecret:      "placeholder_value",
		AwsRegion:      "placeholder_value",
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

	if value := os.Getenv("AWS_ACCESS_KEY_ID"); len(value) > 0 {
		result.AwsAccessKeyId = value
	}

	if value := os.Getenv("AWS_SECRET_ACCESS_KEY"); len(value) > 0 {
		result.AwsSecret = value
	}

	if value := os.Getenv("AWS_REGION"); len(value) > 0 {
		result.AwsRegion = value
	}

	if value := os.Getenv("FLUENT_URL"); len(value) > 0 {
		result.FluentUrl = value
	}

	if value := os.Getenv("PORT"); len(value) > 0 {
		if intVal, err := strconv.Atoi(value); err == nil {
			result.Port = intVal
		}
	}

	if value := os.Getenv("GQL_PORT"); len(value) > 0 {
		if intVal, err := strconv.Atoi(value); err == nil {
			result.GqlPort = intVal
		}
	}

	if value := os.Getenv("JWT_SECRET"); len(value) > 0 {
		result.JWTSecret = value
	}

	return result
}
