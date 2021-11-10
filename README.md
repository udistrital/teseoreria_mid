# tesoreria_mid
Api intermediaria entre el cliente de tesoreria y las apis necesarias para la gestión de la información para estos mismos.
Api mid para el subsistema de tesoreria que hace parte del sistema KRONOS


## Especificaciones Técnicas

### Tecnologías Implementadas y Versiones
* [Golang](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/golang.md)
* [BeeGo](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/beego.md)
* [Docker](https://docs.docker.com/engine/install/ubuntu/)
* [Docker Compose](https://docs.docker.com/compose/)

### Variables de Entorno
```shell
# Ejemplo que se debe actualizar acorde al proyecto
API_NAME = [Nombre del API]
TESORERIA_MID_HTTP_PORT = [Puerto de ejecución del API]
TESORERIA_MID_RUN_MODE = [Modo de ejecución]
AVANCES_CRUD_URL = [URL del despliegue del api de avances_crud]
SOLICITUDES_CRUD_URL = [URL del despliegue del api de solicitudes_crud]
TERCEROS_CRUD_URL = [URL del despliegue del api de terceros_crud]
GIROS_CRUD_URL = [URL del despliegue del api de giros_crud]
```
**NOTA:** Las variables se pueden ver en el fichero conf/app.conf y .env


### Ejecución del Proyecto
```shell
#1. Obtener el repositorio con Go
go get github.com/udistrital/tesoreria_mid

#2. Moverse a la carpeta del repositorio
cd $GOPATH/src/github.com/udistrital/tesoreria_mid

# 3. Moverse a la rama **develop**
git pull origin develop && git checkout develop

# 4. alimentar todas las variables de entorno que utiliza el proyecto.
TESORERIA_MID_HTTP_PORT=8080 AVANCES_CRUD_URL=https://127.0.0.1/v1 TESORERIA_MID_RUN_MODE=dev bee run
```

### Ejecución Dockerfile
```shell
# Implementado para despliegue del Sistema de integración continua CI.
# docker build --tag=avances_crud . --no-cache
# docker run -p 80:80 avances_crud
```

### Ejecución docker-compose
```shell
#1. Clonar el repositorio
git clone -b develop https://github.com/udistrital/tesoreria_mid

#2. Moverse a la carpeta del repositorio
cd solicitudes_crud

#3. Crear un fichero con el nombre **custom.env**
touch .env

#4. Crear la network **back_end** para los contenedores
docker network create back_end

#5. Ejecutar el compose del contenedor
docker-compose up --build

#6. Comprobar que los contenedores estén en ejecución
docker ps
```

### Apis Requeridas
1. [avances_crud](https://github.com/udistrital/avances_crud)

### Ejecución Pruebas

Pruebas unitarias
```shell
# Not Data
```

## Estado CI
| Develop | Release 0.1.1 | Master |
| -- | -- | -- |
| [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/tesoreria_mid/status.svg?ref=refs/heads/develop)](https://hubci.portaloas.udistrital.edu.co/udistrital/tesoreria_mid) | [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/tesoreria_mid/status.svg?ref=refs/heads/release/0.1.1)](https://hubci.portaloas.udistrital.edu.co/udistrital/tesoreria_mid) | [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/tesoreria_mid/status.svg?ref=refs/heads/master)](https://hubci.portaloas.udistrital.edu.co/udistrital/tesoreria_mid) |


## Licencia
This file is part of tesoreria_mid

tesoreria_mid is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

tesoreria_mid is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with tesoreria_mid. If not, see https://www.gnu.org/licenses/.
