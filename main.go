package main

import (
	"flag"

	routes "github.com/nmarsollier/authgo/rest"
)

//	@title			AuthGo
//	@version		1.0
//	@description	Microservicio de Autentificación.

//	@contact.name	Nestor Marsollier
//	@contact.email	nmarsollier@gmail.com

// @host		localhost:3000
// @BasePath	/v1
func main() {
	// For logging
	flag.Parse()
	flag.Set("logtostderr", "true")
	flag.Set("v", "2")

	routes.Start()
}
