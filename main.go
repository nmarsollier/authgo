package main

import (
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
	routes.Start()
}
