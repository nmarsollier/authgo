package env

import (
	"cmp"
	"os"

	"github.com/nmarsollier/authgo/internal/common/strs"
)

// Configuration properties
type Configuration struct {
	ServerName string `json:"serverName"`
	Port       int    `json:"port"`
	GqlPort    int    `json:"gqlPort"`
	RabbitURL  string `json:"rabbitUrl"`
	MongoURL   string `json:"mongoUrl"`
	JWTSecret  string `json:"jwtSecret"`
	FluentURL  string `json:"fluentUrl"`
}

var config *Configuration

// Get Obtiene las variables de entorno del sistema
func Get() *Configuration {
	if config == nil {
		config = load()
	}

	return config
}

// Load file properties
func load() *Configuration {
	return &Configuration{
		ServerName: cmp.Or(os.Getenv("SERVER_NAME"), "authgo"),
		Port:       cmp.Or(strs.AtoiZero(os.Getenv("PORT")), 3000),
		GqlPort:    cmp.Or(strs.AtoiZero(os.Getenv("GQL_PORT")), 4000),
		RabbitURL:  cmp.Or(os.Getenv("RABBIT_URL"), "amqp://localhost"),
		MongoURL:   cmp.Or(os.Getenv("MONGO_URL"), "mongodb://localhost:27017"),
		JWTSecret:  cmp.Or(os.Getenv("JWT_SECRET"), "ecb6d3479ac3823f1da7f314d871989b"),
		FluentURL:  cmp.Or(os.Getenv("FLUENT_URL"), "localhost:24224"),
	}
}
