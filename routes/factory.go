package routes

import (
	"github.com/nmarsollier/authgo/security"
	"github.com/nmarsollier/authgo/user"
)

var userServiceInstance user.Service

func getUserService() user.Service {
	if userServiceInstance != nil {
		return userServiceInstance
	}

	userServiceInstance = user.NewService()
	return userServiceInstance
}

var securityServiceInstance security.Service

func getSecurityService() security.Service {
	if securityServiceInstance != nil {
		return securityServiceInstance
	}

	securityServiceInstance = security.NewService()
	return securityServiceInstance
}
