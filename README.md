Experimental : Auth Service en GO
==

Requisitos
-

Go 1.10  <https://golang.org/doc/install>

Configuración inicial
-

establecer variables de entorno (consultar documentación de la version instalada)

```bash
export GOPATH="$HOME/go"
export GOROOT=/usr/local/go
export PATH="$PATH:$GOPATH/bin:$GOROOT/bin"
```

Para descargar el proyecto correctamente hay que ejecutar :

```bash
go get github.com/nmarsollier/ms_auth_go
```

Una vez descargado, tendremos el proyecto en la carpeta

```bash
GOPATH="$HOME/go/src/github.com/nmarsollier/ms_auth_go
```

Instalar Librerías requeridas
-

```bash
go get github.com/gin-gonic/gin
go get github.com/mongodb/mongo-go-driver/mongo
go get golang.org/x/crypto/bcrypt
go get github.com/dgrijalva/jwt-go
go get github.com/itsjamie/gin-cors
go get github.com/patrickmn/go-cache
go get github.com/streadway/amqp
go get github.com/gin-contrib/static
```

Build y ejecución
-

```bash
go install github.com/nmarsollier/ms_auth_go
ms_auth_go
```

Apidoc
-

Apidoc es una herramienta para proyectos node, para que funcione correctamente hay que instalarla globalmente con

```bash
npm install apidoc -g
```

La documentación necesita ser generada manualmente ejecutando la siguiente linea en la carpeta ms_auth_go :

```bash
apidoc -s src -o www
```

Esto nos genera una carpeta public con la documentación, esta carpeta debe estar presente desde donde se ejecute auth, auth busca ./www para localizarlo, aunque se puede configurar desde el archivo de properties.
