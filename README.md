### Si queres sabes mas sobre mi:
[Nestor Marsollier](https://github.com/nmarsollier/profile)

# Auth Service en GO

Este Microservicio de seguridad reemplaza al del proyecto

[Microservicios Auth](https://github.com/nmarsollier/ecommerce)

Se encarga de registrar y autenticar usuarios en el sistema.

Utiliza el esquema JWT con un header Authorization "bearer" estándar.

Cada usuario tiene asociado una lista de permisos, existen 2 permisos genéricos "user" y "admin". Los usuarios que se registran son todos "user",  muchos procesos necesitan un usuario "admin" para poder funcionar, por lo tanto hay que editar el esquema en mongodb para asociarle el permiso admin a algún usuario inicialmente.

[Documentación de API](./README-API.md)

La documentación de las api también se pueden consultar desde el home del microservicio
que una vez levantado el servidor se puede navegar en [localhost:3000](http://localhost:3000/)

## Requisitos

Go 1.14  [golang.org](https://golang.org/doc/install)


## Configuración inicial

establecer variables de entorno (consultar documentación de la version instalada)

```bash
export GO111MODULE=on
export GOFLAGS=-mod=vendor
```

Para descargar el proyecto correctamente hay que ejecutar :

```bash
go get github.com/nmarsollier/authgo
```

Una vez descargado, tendremos el codigo fuente del proyecto en la carpeta

```bash
cd $GOPATH/src/github.com/nmarsollier/authgo
```

## Instalar Librerías requeridas


```bash
go mod download
go mod vendor
```

Build y ejecución
-

```bash
go install
authgo
```

## MongoDB

La base de datos se almacena en MongoDb.

Seguir las guías de instalación de mongo desde el sitio oficial [mongodb.com](https://www.mongodb.com/download-center#community)

No se requiere ninguna configuración adicional, solo levantarlo luego de instalarlo.

## RabbitMQ

Este microservicio notifica los logouts de usuarios con Rabbit.

Seguir los pasos de instalación en la pagina oficial [rabbitmq.com](https://www.rabbitmq.com/)

No se requiere ninguna configuración adicional, solo levantarlo luego de instalarlo.

## Apidoc

Apidoc es una herramienta que genera documentación de apis para proyectos node (ver [Apidoc](http://apidocjs.com/)).

El microservicio muestra la documentación como archivos estáticos si se abre en un browser la raíz del servidor [localhost:3000](http://localhost:3000/)

Ademas se genera la documentación en formato markdown.

Para que funcione correctamente hay que instalarla globalmente con

```bash
npm install apidoc -g
npm install -g apidoc-markdown2
```

La documentación necesita ser generada manualmente ejecutando la siguiente linea en la carpeta auth :

```bash
apidoc -o www
apidoc-markdown2 -p www -o README-API.md
```

Esto nos genera una carpeta con la documentación, esta carpeta debe estar presente desde donde se ejecute auth, auth busca ./www para localizarlo, aunque se puede configurar desde el archivo de properties.

## Archivo config.json

Este archivo permite configurar los parámetros del servidor, ver ejemplos en config-example.json.
El servidor busca el archivo "./config.json". Podemos definir el archivo su ruta completa ejecutando

```bash
authgo [path_to_config.json]
```

Para mas detalles ver el archivo tools/env/env.go


## Docker

Tambien podemos usar docker en este repositorio, ejecutamos :

```bash
docker build -t dev-auth-go -f Dockerfile.dev .
docker run -d --name dev-auth-go --network host dev-auth-go
```

El contenedor se puede parar usando :

```bash
docker stop dev-auth-go
```
Se vuelve a levantar usando 

```bash
docker start dev-auth-go 
```
