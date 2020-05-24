define({ "api": [
  {
    "type": "fanout",
    "url": "auth/fanout",
    "title": "Invalidar Token",
    "group": "RabbitMQ_POST",
    "description": "<p>AuthService enviá un broadcast a todos los usuarios cuando un token ha sido invalidado. Los clientes deben eliminar de sus caches las sesiones invalidadas.</p>",
    "success": {
      "examples": [
        {
          "title": "Mensaje",
          "content": "{\n   \"type\": \"logout\",\n   \"message\": \"{Token revocado}\"\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./rabbit/rabbit.go",
    "groupTitle": "RabbitMQ_POST",
    "name": "FanoutAuthFanout"
  },
  {
    "type": "post",
    "url": "/v1/user/password",
    "title": "Cambiar Password",
    "name": "Cambiar_Password",
    "group": "Seguridad",
    "description": "<p>Cambia la contraseña del usuario actual.</p>",
    "examples": [
      {
        "title": "Body",
        "content": "{\n  \"currentPassword\" : \"{Contraseña actual}\",\n  \"newPassword\" : \"{Nueva Contraseña}\",\n}",
        "type": "json"
      },
      {
        "title": "Header Autorización",
        "content": "Authorization=bearer {token}",
        "type": "String"
      }
    ],
    "success": {
      "examples": [
        {
          "title": "Respuesta",
          "content": "HTTP/1.1 200 OK",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./routes/controllers.go",
    "groupTitle": "Seguridad",
    "error": {
      "examples": [
        {
          "title": "401 Unauthorized",
          "content": "HTTP/1.1 401 Unauthorized\n{\n   \"error\" : \"Unauthorized\"\n}",
          "type": "json"
        },
        {
          "title": "400 Bad Request",
          "content": "HTTP/1.1 400 Bad Request\n{\n   \"messages\" : [\n     {\n       \"path\" : \"{Nombre de la propiedad}\",\n       \"message\" : \"{Motivo del error}\"\n     },\n     ...\n  ]\n}",
          "type": "json"
        },
        {
          "title": "500 Server Error",
          "content": "HTTP/1.1 500 Internal Server Error\n{\n   \"error\" : \"Not Found\"\n}",
          "type": "json"
        }
      ]
    }
  },
  {
    "type": "post",
    "url": "/v1/users/:userId/disable",
    "title": "Deshabilitar Usuario",
    "name": "Deshabilitar_Usuario",
    "group": "Seguridad",
    "description": "<p>Deshabilita un usuario en el sistema.   El usuario logueado debe tener permisos &quot;admin&quot;.</p>",
    "success": {
      "examples": [
        {
          "title": "Respuesta",
          "content": "HTTP/1.1 200 OK",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./routes/controllers.go",
    "groupTitle": "Seguridad",
    "examples": [
      {
        "title": "Header Autorización",
        "content": "Authorization=bearer {token}",
        "type": "String"
      }
    ],
    "error": {
      "examples": [
        {
          "title": "401 Unauthorized",
          "content": "HTTP/1.1 401 Unauthorized\n{\n   \"error\" : \"Unauthorized\"\n}",
          "type": "json"
        },
        {
          "title": "400 Bad Request",
          "content": "HTTP/1.1 400 Bad Request\n{\n   \"messages\" : [\n     {\n       \"path\" : \"{Nombre de la propiedad}\",\n       \"message\" : \"{Motivo del error}\"\n     },\n     ...\n  ]\n}",
          "type": "json"
        },
        {
          "title": "500 Server Error",
          "content": "HTTP/1.1 500 Internal Server Error\n{\n   \"error\" : \"Not Found\"\n}",
          "type": "json"
        }
      ]
    }
  },
  {
    "type": "post",
    "url": "/v1/users/:userId/enable",
    "title": "Habilitar Usuario",
    "name": "Habilitar_Usuario",
    "group": "Seguridad",
    "description": "<p>Habilita un usuario en el sistema. El usuario logueado debe tener permisos &quot;admin&quot;.</p>",
    "success": {
      "examples": [
        {
          "title": "Respuesta",
          "content": "HTTP/1.1 200 OK",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./routes/controllers.go",
    "groupTitle": "Seguridad",
    "examples": [
      {
        "title": "Header Autorización",
        "content": "Authorization=bearer {token}",
        "type": "String"
      }
    ],
    "error": {
      "examples": [
        {
          "title": "401 Unauthorized",
          "content": "HTTP/1.1 401 Unauthorized\n{\n   \"error\" : \"Unauthorized\"\n}",
          "type": "json"
        },
        {
          "title": "400 Bad Request",
          "content": "HTTP/1.1 400 Bad Request\n{\n   \"messages\" : [\n     {\n       \"path\" : \"{Nombre de la propiedad}\",\n       \"message\" : \"{Motivo del error}\"\n     },\n     ...\n  ]\n}",
          "type": "json"
        },
        {
          "title": "500 Server Error",
          "content": "HTTP/1.1 500 Internal Server Error\n{\n   \"error\" : \"Not Found\"\n}",
          "type": "json"
        }
      ]
    }
  },
  {
    "type": "get",
    "url": "/v1/users",
    "title": "Listar Usuarios",
    "name": "Listar_Usuarios",
    "group": "Seguridad",
    "description": "<p>Obtiene información de todos los usuarios.</p>",
    "success": {
      "examples": [
        {
          "title": "Respuesta",
          "content": "    HTTP/1.1 200 OK\n    [{\n       \"id\": \"{Id usuario}\",\n       \"name\": \"{Nombre del usuario}\",\n       \"login\": \"{Login de usuario}\",\n       \"permissions\": [\n           \"{Permission}\"\n       ],\n\t      \"enabled\": true|false\n    }, ...]",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./routes/controllers.go",
    "groupTitle": "Seguridad",
    "examples": [
      {
        "title": "Header Autorización",
        "content": "Authorization=bearer {token}",
        "type": "String"
      }
    ],
    "error": {
      "examples": [
        {
          "title": "401 Unauthorized",
          "content": "HTTP/1.1 401 Unauthorized\n{\n   \"error\" : \"Unauthorized\"\n}",
          "type": "json"
        },
        {
          "title": "500 Server Error",
          "content": "HTTP/1.1 500 Internal Server Error\n{\n   \"error\" : \"Not Found\"\n}",
          "type": "json"
        }
      ]
    }
  },
  {
    "type": "post",
    "url": "/v1/user/signin",
    "title": "Login",
    "name": "Login",
    "group": "Seguridad",
    "description": "<p>Loguea un usuario en el sistema.</p>",
    "examples": [
      {
        "title": "Body",
        "content": "{\n  \"login\": \"{Login de usuario}\",\n  \"password\": \"{Contraseña}\"\n}",
        "type": "json"
      }
    ],
    "version": "0.0.0",
    "filename": "./routes/controllers.go",
    "groupTitle": "Seguridad",
    "success": {
      "examples": [
        {
          "title": "Respuesta",
          "content": "HTTP/1.1 200 OK\n{\n  \"token\": \"{Token de autorización}\"\n}",
          "type": "json"
        }
      ]
    },
    "error": {
      "examples": [
        {
          "title": "400 Bad Request",
          "content": "HTTP/1.1 400 Bad Request\n{\n   \"messages\" : [\n     {\n       \"path\" : \"{Nombre de la propiedad}\",\n       \"message\" : \"{Motivo del error}\"\n     },\n     ...\n  ]\n}",
          "type": "json"
        },
        {
          "title": "500 Server Error",
          "content": "HTTP/1.1 500 Internal Server Error\n{\n   \"error\" : \"Not Found\"\n}",
          "type": "json"
        }
      ]
    }
  },
  {
    "type": "get",
    "url": "/v1/user/signout",
    "title": "Logout",
    "name": "Logout",
    "group": "Seguridad",
    "description": "<p>Desloguea un usuario en el sistema, invalida el token.</p>",
    "success": {
      "examples": [
        {
          "title": "Respuesta",
          "content": "HTTP/1.1 200 OK",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./routes/controllers.go",
    "groupTitle": "Seguridad",
    "examples": [
      {
        "title": "Header Autorización",
        "content": "Authorization=bearer {token}",
        "type": "String"
      }
    ],
    "error": {
      "examples": [
        {
          "title": "401 Unauthorized",
          "content": "HTTP/1.1 401 Unauthorized\n{\n   \"error\" : \"Unauthorized\"\n}",
          "type": "json"
        },
        {
          "title": "500 Server Error",
          "content": "HTTP/1.1 500 Internal Server Error\n{\n   \"error\" : \"Not Found\"\n}",
          "type": "json"
        }
      ]
    }
  },
  {
    "type": "post",
    "url": "/v1/users/:userId/grant",
    "title": "Otorga Permisos",
    "name": "Otorga_Permisos",
    "group": "Seguridad",
    "description": "<p>Otorga permisos al usuario indicado, el usuario logueado tiene que tener permiso &quot;admin&quot;.</p>",
    "examples": [
      {
        "title": "Body",
        "content": "{\n  \"permissions\" : [\"permiso\", ...],\n}",
        "type": "json"
      },
      {
        "title": "Header Autorización",
        "content": "Authorization=bearer {token}",
        "type": "String"
      }
    ],
    "success": {
      "examples": [
        {
          "title": "Respuesta",
          "content": "HTTP/1.1 200 OK",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./routes/controllers.go",
    "groupTitle": "Seguridad",
    "error": {
      "examples": [
        {
          "title": "401 Unauthorized",
          "content": "HTTP/1.1 401 Unauthorized\n{\n   \"error\" : \"Unauthorized\"\n}",
          "type": "json"
        },
        {
          "title": "400 Bad Request",
          "content": "HTTP/1.1 400 Bad Request\n{\n   \"messages\" : [\n     {\n       \"path\" : \"{Nombre de la propiedad}\",\n       \"message\" : \"{Motivo del error}\"\n     },\n     ...\n  ]\n}",
          "type": "json"
        },
        {
          "title": "500 Server Error",
          "content": "HTTP/1.1 500 Internal Server Error\n{\n   \"error\" : \"Not Found\"\n}",
          "type": "json"
        }
      ]
    }
  },
  {
    "type": "post",
    "url": "/v1/user",
    "title": "Registrar Usuario",
    "name": "Registrar_Usuario",
    "group": "Seguridad",
    "description": "<p>Registra un nuevo usuario en el sistema.</p>",
    "examples": [
      {
        "title": "Body",
        "content": "{\n  \"name\": \"{Nombre de Usuario}\",\n  \"login\": \"{Login de usuario}\",\n  \"password\": \"{Contraseña}\"\n}",
        "type": "json"
      }
    ],
    "version": "0.0.0",
    "filename": "./routes/controllers.go",
    "groupTitle": "Seguridad",
    "success": {
      "examples": [
        {
          "title": "Respuesta",
          "content": "HTTP/1.1 200 OK\n{\n  \"token\": \"{Token de autorización}\"\n}",
          "type": "json"
        }
      ]
    },
    "error": {
      "examples": [
        {
          "title": "400 Bad Request",
          "content": "HTTP/1.1 400 Bad Request\n{\n   \"messages\" : [\n     {\n       \"path\" : \"{Nombre de la propiedad}\",\n       \"message\" : \"{Motivo del error}\"\n     },\n     ...\n  ]\n}",
          "type": "json"
        },
        {
          "title": "500 Server Error",
          "content": "HTTP/1.1 500 Internal Server Error\n{\n   \"error\" : \"Not Found\"\n}",
          "type": "json"
        }
      ]
    }
  },
  {
    "type": "post",
    "url": "/v1/users/:userId/revoke",
    "title": "Revoca Permisos",
    "name": "Revoca_Permisos",
    "group": "Seguridad",
    "description": "<p>Quita permisos al usuario indicado, el usuario logueado tiene que tener permiso &quot;admin&quot;.</p>",
    "examples": [
      {
        "title": "Body",
        "content": "{\n  \"permissions\" : [\"permiso\", ...],\n}",
        "type": "json"
      },
      {
        "title": "Header Autorización",
        "content": "Authorization=bearer {token}",
        "type": "String"
      }
    ],
    "success": {
      "examples": [
        {
          "title": "Respuesta",
          "content": "HTTP/1.1 200 OK",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./routes/controllers.go",
    "groupTitle": "Seguridad",
    "error": {
      "examples": [
        {
          "title": "401 Unauthorized",
          "content": "HTTP/1.1 401 Unauthorized\n{\n   \"error\" : \"Unauthorized\"\n}",
          "type": "json"
        },
        {
          "title": "400 Bad Request",
          "content": "HTTP/1.1 400 Bad Request\n{\n   \"messages\" : [\n     {\n       \"path\" : \"{Nombre de la propiedad}\",\n       \"message\" : \"{Motivo del error}\"\n     },\n     ...\n  ]\n}",
          "type": "json"
        },
        {
          "title": "500 Server Error",
          "content": "HTTP/1.1 500 Internal Server Error\n{\n   \"error\" : \"Not Found\"\n}",
          "type": "json"
        }
      ]
    }
  },
  {
    "type": "get",
    "url": "/v1/users/current",
    "title": "Usuario Actual",
    "name": "Usuario_Actual",
    "group": "Seguridad",
    "description": "<p>Obtiene información del usuario actual.</p>",
    "success": {
      "examples": [
        {
          "title": "Respuesta",
          "content": "HTTP/1.1 200 OK\n{\n   \"id\": \"{Id usuario}\",\n   \"name\": \"{Nombre del usuario}\",\n   \"login\": \"{Login de usuario}\",\n   \"permissions\": [\n       \"{Permission}\"\n   ]\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./routes/controllers.go",
    "groupTitle": "Seguridad",
    "examples": [
      {
        "title": "Header Autorización",
        "content": "Authorization=bearer {token}",
        "type": "String"
      }
    ],
    "error": {
      "examples": [
        {
          "title": "401 Unauthorized",
          "content": "HTTP/1.1 401 Unauthorized\n{\n   \"error\" : \"Unauthorized\"\n}",
          "type": "json"
        },
        {
          "title": "500 Server Error",
          "content": "HTTP/1.1 500 Internal Server Error\n{\n   \"error\" : \"Not Found\"\n}",
          "type": "json"
        }
      ]
    }
  }
] });
