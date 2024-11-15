### Si queres sabes mas sobre mi:

[Nestor Marsollier](https://github.com/nmarsollier/profile)

# Auth Service en GO

Este Microservicio de seguridad reemplaza al del proyecto

[Microservicios Auth](https://github.com/nmarsollier/ecommerce)

Se encarga de registrar y autenticar usuarios en el sistema.

Utiliza el esquema JWT con un header Authorization "Bearer" estándar.

Cada usuario tiene asociado una lista de permisos, existen 2 permisos genéricos "user" y "admin". Los usuarios que se registran son todos "user", muchos procesos necesitan un usuario "admin" para poder funcionar, por lo tanto hay que editar el esquema en mongodb para asociarle el permiso admin a algún usuario inicialmente.

[Documentación de API](./README-API.md)

La documentación de las api también se pueden consultar desde el home del microservicio
que una vez levantado el servidor se puede navegar en [localhost:3000](http://localhost:3000/docs/index.html)

La interfaz GraphQL en [localhost:4000](http://localhost:4000/docs/index.html)

## Requisitos

Go [golang.org](https://golang.org/doc/install)

## Configuración inicial

Para descargar el proyecto correctamente hay que ejecutar :

```bash
git clone https://github.com/nmarsollier/authgo $GOPATH/src/github.com/nmarsollier/authgo
```

Una vez descargado, tendremos el codigo fuente del proyecto en la carpeta

```bash
cd $GOPATH/src/github.com/nmarsollier/authgo
```

## Instalar Librerías requeridas

```bash
git config core.hooksPath .githooks
go install github.com/swaggo/gin-swagger/swaggerFiles
go install github.com/swaggo/gin-swagger
go install github.com/swaggo/swag/cmd/swag
go install github.com/golang/mock/mockgen@v1.6.0
go install github.com/99designs/gqlgen@v0.17.56
```

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

## Swagger

Usamos [swaggo](https://github.com/swaggo/swag)

Requisitos

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

La documentacion la generamos con el comando

```bash
swag init
```

Para generar el archivo README-API.md

Requisito

```bash
sudo npm install -g swagger-markdown
```

y ejecutamos

```bash
npx swagger-markdown -i ./docs/swagger.yaml -o README-API.md
```

## GraphQL

El esquema se genera con

```bash
./generate_gql.sh
```

## Tests

```bash
go install github.com/golang/mock/mockgen@v1.6.0
```

Y ejecutamos el script generate_mock.sh

## Configuración del servidor

Este servidor usa las siguientes variables de entorno para configuración :

```
RABBIT_URL : Url de rabbit (default amqp://localhost)
MONGO_URL : Url de mongo (default mongodb://localhost:27017)
PORT : Puerto (default 3000)
GQL_PORT : Puerto de graphql (default 4000)
JWT_SECRET : Secret para password (default ecb6d3479ac3823f1da7f314d871989b)
```

## Docker

Estos comandos son para dockerizar el microservicio desde el codigo descargado localmente.

### Build

Hacer un build local para ejecutar en docker :

```bash
docker build -t dev-auth-go .
```

### El contenedor

Mac | Windows

```bash
docker run -it --name dev-auth-go -p 3000:3000 -p 4000:4000 -v $PWD:/go/src/github.com/nmarsollier/authgo dev-auth-go
```

Linux

```bash
docker run -it --add-host host.docker.internal:172.17.0.1 --name dev-auth-go -p 3000:3000 -p 4000:4000 -v $PWD:/go/src/github.com/nmarsollier/authgo dev-auth-go
```
