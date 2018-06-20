package env

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {

	jsonErr, _ := json.Marshal(Get())
	assert.Equal(t, "{\"port\":3000,\"rabbitUrl\":\"amqp://localhost\",\"mongoUrl\":\"mongodb://localhost:27017\",\"wwwPath\":\"www\",\"jwtSecret\":\"ecb6d3479ac3823f1da7f314d871989b\"}", string(jsonErr))

}

func TestLoad(t *testing.T) {
	loaded := Load("./test/mocks/missing.json")
	assert.Equal(t, loaded, false)
	jsonErr, _ := json.Marshal(Get())
	assert.Equal(t, "{\"port\":3000,\"rabbitUrl\":\"amqp://localhost\",\"mongoUrl\":\"mongodb://localhost:27017\",\"wwwPath\":\"www\",\"jwtSecret\":\"ecb6d3479ac3823f1da7f314d871989b\"}", string(jsonErr))

	loaded = Load("missing.json")
	assert.Equal(t, loaded, false)
	jsonErr, _ = json.Marshal(Get())
	assert.Equal(t, "{\"port\":3000,\"rabbitUrl\":\"amqp://localhost\",\"mongoUrl\":\"mongodb://localhost:27017\",\"wwwPath\":\"www\",\"jwtSecret\":\"ecb6d3479ac3823f1da7f314d871989b\"}", string(jsonErr))

	loaded = Load("env_test_config.json")
	assert.Equal(t, loaded, true)
	jsonErr, _ = json.Marshal(Get())
	assert.Equal(t, "{\"port\":80,\"rabbitUrl\":\"otroUrl\",\"mongoUrl\":\"mongodb://localhost:27017\",\"wwwPath\":\"www\",\"jwtSecret\":\"ecb6d3479ac3823f1da7f314d871989b\"}", string(jsonErr))
}
