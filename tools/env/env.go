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
		if ok := Load("config.json"); !ok {
			config = new()
		}
	}

	return config
}

// Load file properties
func Load(fileName string) bool {
	file, err := os.Open(fileName)
	if err != nil {
		log.Output(1, fmt.Sprintf("%s : %s", fileName, err.Error()))
		return false
	}
	defer file.Close()

	loaded := new()
	err = json.NewDecoder(file).Decode(&loaded)
	if err != nil {
		log.Output(1, fmt.Sprintf("%s : %s", fileName, err.Error()))
		return false
	}

	config = loaded
	log.Output(1, fmt.Sprintf("%s cargado correctamente", fileName))
	return true
}
