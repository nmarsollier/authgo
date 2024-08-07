// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Nestor Marsollier",
            "email": "nmarsollier@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/rabbit/logout": {
            "put": {
                "description": "SendLogout envía un broadcast a rabbit con logout. Esto no es Rest es RabbitMQ.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Rabbit"
                ],
                "summary": "Mensage Rabbit",
                "parameters": [
                    {
                        "description": "Token deshabilitado",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/rabbit.message"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/v1/user": {
            "post": {
                "description": "Registra un nuevo usuario en el sistema.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Seguridad"
                ],
                "summary": "Registrar Usuario",
                "parameters": [
                    {
                        "description": "Informacion de ususario",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.SignUpRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User Token",
                        "schema": {
                            "$ref": "#/definitions/rest.tokenResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/app_errors.ErrValidation"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    }
                }
            }
        },
        "/v1/user/password": {
            "post": {
                "description": "Cambia la contraseña del usuario actual.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Seguridad"
                ],
                "summary": "Cambiar Password",
                "parameters": [
                    {
                        "description": "Passwords",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/rest.changePasswordBody"
                        }
                    },
                    {
                        "type": "string",
                        "description": "bearer {token}",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/app_errors.ErrValidation"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    }
                }
            }
        },
        "/v1/user/signin": {
            "post": {
                "description": "Loguea un usuario en el sistema.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Seguridad"
                ],
                "summary": "Login",
                "parameters": [
                    {
                        "description": "Sign in information",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.SignInRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User Token",
                        "schema": {
                            "$ref": "#/definitions/rest.tokenResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/app_errors.ErrValidation"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    }
                }
            }
        },
        "/v1/user/signout": {
            "get": {
                "description": "Desloguea un usuario en el sistema, invalida el token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Seguridad"
                ],
                "summary": "Logout",
                "parameters": [
                    {
                        "type": "string",
                        "description": "bearer {token}",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "No Content"
                    },
                    "500": {
                        "description": "Error response",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    }
                }
            }
        },
        "/v1/users": {
            "get": {
                "description": "Obtiene información de todos los usuarios. El usuario logueado debe tener permisos \"admin\".",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Seguridad"
                ],
                "summary": "Listar Usuarios",
                "parameters": [
                    {
                        "type": "string",
                        "description": "bearer {token}",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Users",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/rest.UserDataResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/app_errors.ErrValidation"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    }
                }
            }
        },
        "/v1/users/:userID/grant": {
            "post": {
                "responses": {
                    "200": {
                        "description": "User Token",
                        "schema": {
                            "$ref": "#/definitions/rest.tokenResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/app_errors.ErrValidation"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    }
                }
            }
        },
        "/v1/users/:userID/revoke": {
            "post": {
                "responses": {
                    "200": {
                        "description": "User Token",
                        "schema": {
                            "$ref": "#/definitions/rest.revokePermissionBody"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/app_errors.ErrValidation"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    }
                }
            }
        },
        "/v1/users/:userId/disable": {
            "post": {
                "responses": {
                    "200": {
                        "description": "User Token",
                        "schema": {
                            "$ref": "#/definitions/rest.tokenResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/app_errors.ErrValidation"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    }
                }
            }
        },
        "/v1/users/:userId/enable": {
            "post": {
                "responses": {
                    "200": {
                        "description": "User Token",
                        "schema": {
                            "$ref": "#/definitions/rest.tokenResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/app_errors.ErrValidation"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    }
                }
            }
        },
        "/v1/users/current": {
            "get": {
                "description": "Obtiene información del usuario actual.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Seguridad"
                ],
                "summary": "Usuario Actual",
                "parameters": [
                    {
                        "type": "string",
                        "description": "bearer {token}",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User data",
                        "schema": {
                            "$ref": "#/definitions/rest.UserResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/app_errors.ErrValidation"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app_errors.OtherErrors"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "app_errors.ErrField": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "path": {
                    "type": "string"
                }
            }
        },
        "app_errors.ErrValidation": {
            "type": "object",
            "properties": {
                "messages": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/app_errors.ErrField"
                    }
                }
            }
        },
        "app_errors.OtherErrors": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "rabbit.message": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "rest.UserDataResponse": {
            "type": "object",
            "properties": {
                "enabled": {
                    "type": "boolean"
                },
                "id": {
                    "type": "string"
                },
                "login": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "permissions": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "rest.UserResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "login": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "permissions": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "rest.changePasswordBody": {
            "type": "object",
            "required": [
                "currentPassword",
                "newPassword"
            ],
            "properties": {
                "currentPassword": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 1
                },
                "newPassword": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 1
                }
            }
        },
        "rest.revokePermissionBody": {
            "type": "object",
            "required": [
                "permissions",
                "user"
            ],
            "properties": {
                "permissions": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "user": {
                    "type": "string"
                }
            }
        },
        "rest.tokenResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "user.SignInRequest": {
            "type": "object",
            "required": [
                "login",
                "password"
            ],
            "properties": {
                "login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "user.SignUpRequest": {
            "type": "object",
            "required": [
                "login",
                "name",
                "password"
            ],
            "properties": {
                "login": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:3000",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "AuthGo",
	Description:      "Microservicio de Autentificación.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
