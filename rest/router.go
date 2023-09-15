package rest

import (
	"fmt"

	"github.com/nmarsollier/authgo/rest/engine"
	"github.com/nmarsollier/authgo/tools/env"
)

// Start this server
func Start() {
	engine.Router().Run(fmt.Sprintf(":%d", env.Get().Port))
}
