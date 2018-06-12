package env

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Configuration properties
type Configuration struct {
	Port      int    `json:"port"`
	RabbitURL string `json:"rabbitUrl"`
	MongoURL  string `json:"mongoUrl"`
	WWWWPath  string `json:"wwwPath"`
	JWTSecret string `json:"jwtSecret"`
}

var config = Configuration{
	Port:      3000,
	RabbitURL: "amqp://localhost",
	MongoURL:  "mongodb://localhost:27017",
	WWWWPath:  "www",
	JWTSecret: "ecb6d3479ac3823f1da7f314d871989b",
}
var initialized = false

// Get Obtiene las variables de entorno del sistema
func Get() *Configuration {
	if !initialized {
		Load("config.json")
		initialized = true
	}

	return &config
}

// Load file properties
func Load(fileName string) bool {
	file, err := os.Open(fileName)
	if err != nil {
		log.Output(1, fmt.Sprintf("%s : %s", fileName, err.Error()))
		return false
	}

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Output(1, fmt.Sprintf("%s : %s", fileName, err.Error()))
		return false
	}

	log.Output(1, fmt.Sprintf("%s cargado correctamente", fileName))
	initialized = true
	return true
}
