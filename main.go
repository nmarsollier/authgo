package main

import (
	_ "github.com/nmarsollier/authgo/docs"
	"github.com/nmarsollier/authgo/internal/graph"
	"github.com/nmarsollier/authgo/internal/rest"
)

//	@title			AuthGo
//	@version		1.0
//	@description	Microservicio de Autentificación.

//	@contact.name	Nestor Marsollier
//	@contact.email	nmarsollier@gmail.com

// @host		localhost:3000
// @BasePath	/v1
func main() {
	go graph.Start()

	rest.Start()
}
