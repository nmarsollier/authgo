package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Set environment variables
	os.Setenv("RABBIT_URL", "amqp://rabbitmq:5672")
	os.Setenv("MONGO_URL", "mongodb://localhost:27017")
	os.Setenv("FLUENT_URL", "fluentd:24224")
	os.Setenv("PORT", "8080")
	os.Setenv("JWT_SECRET", "mysecret")

	// Call the load function
	config := load()

	// Assert the values
	assert.Equal(t, "amqp://rabbitmq:5672", config.RabbitURL)
	assert.Equal(t, "mongodb://localhost:27017", config.MongoURL)
	assert.Equal(t, "fluentd:24224", config.FluentUrl)
	assert.Equal(t, 8080, config.Port)
	assert.Equal(t, "mysecret", config.JWTSecret)

	// Clear environment variables
	os.Unsetenv("RABBIT_URL")
	os.Unsetenv("MONGO_URL")
	os.Unsetenv("FLUENT_URL")
	os.Unsetenv("PORT")
	os.Unsetenv("JWT_SECRET")
}
