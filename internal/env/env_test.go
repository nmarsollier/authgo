package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	config := load()

	// Assert the values
	assert.Equal(t, "amqp://localhost", config.RabbitURL)
	assert.Equal(t, "mongodb://localhost:27017", config.MongoURL)
	assert.Equal(t, "localhost:24224", config.FluentURL)
	assert.Equal(t, 3000, config.Port)
	assert.Equal(t, "ecb6d3479ac3823f1da7f314d871989b", config.JWTSecret)
}

func TestLoad(t *testing.T) {
	// Set environment variables
	os.Setenv("RABBIT_URL", "amqp://rabbitmq:5672")
	os.Setenv("MONGO_URL", "mongodb://mongourl:27017")
	os.Setenv("FLUENT_URL", "fluentd:24224")
	os.Setenv("PORT", "8080")
	os.Setenv("JWT_SECRET", "mysecret")

	// Call the load function
	config := load()

	// Assert the values
	assert.Equal(t, "amqp://rabbitmq:5672", config.RabbitURL)
	assert.Equal(t, "mongodb://mongourl:27017", config.MongoURL)
	assert.Equal(t, "fluentd:24224", config.FluentURL)
	assert.Equal(t, 8080, config.Port)
	assert.Equal(t, "mysecret", config.JWTSecret)

	// Clear environment variables
	os.Unsetenv("RABBIT_URL")
	os.Unsetenv("MONGO_URL")
	os.Unsetenv("FLUENT_URL")
	os.Unsetenv("PORT")
	os.Unsetenv("JWT_SECRET")
}
