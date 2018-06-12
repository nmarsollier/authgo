Experimental : Auth Service en GO
==

Requisitos
-

Go 1.10  <https://golang.org/doc/install>

Dep <https://github.com/golang/dep>

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

Una vez descargado, tendremos el codigo fuente del proyecto en la carpeta

```bash
cd $GOPATH/src/github.com/nmarsollier/ms_auth_go
```

Instalar Librerías requeridas
-

```bash
dep ensure
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
