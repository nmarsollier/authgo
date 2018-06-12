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
		if file, err := os.Open("config.json"); err == nil {
			err = json.NewDecoder(file).Decode(&config)
			if err != nil {
				log.Output(1, fmt.Sprintf("Error al leer archivo config.xml : %s", err.Error()))
			}
		} else {
			log.Output(1, "No se encontró el archivo de configuraion config.json")
		}

		initialized = true
	}

	return &config
}
