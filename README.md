# Sistema de Registro y Postulación del BIESS

Módulo inicial (EF-11 y EF-12) del Concurso Público de Méritos y Oposición para el Directorio del BIESS. Este repositorio contiene un prototipo académico: la verificación de identidad usa datos ficticios y no está conectada al Registro Civil.

## Stack y versiones

| Componente | Tecnología | Versión requerida o usada |
| --- | --- | --- |
| Frontend | React, Vite, Tailwind CSS | Node.js 20.19+ o 22.12+; npm |
| Backend | Go, Gin | Go 1.27.1 (según `go.mod`) |
| Base de datos | PostgreSQL | 16 (imagen de Docker Compose) |
| Caché, sesiones y limitación de intentos | Redis | 7 (imagen de Docker Compose) |
| Contenedores locales | Docker Compose | Compose v2 |

Vite 8 requiere Node.js 20.19+ o 22.12+. Consulta la [guía oficial de Vite](https://vite.dev/guide/) para los requisitos vigentes.

## Prerrequisitos

- Git.
- Go 1.27.1.
- Node.js 20.19+ o 22.12+ y npm.
- Docker Desktop o Docker Engine con Docker Compose v2.
- Alternativamente, PostgreSQL 16 y Redis 7 instalados localmente.
- `psql` es opcional; se requiere para cargar scripts SQL desde la terminal.

## Instalación y ejecución local

Clona el repositorio:

```bash
git clone https://github.com/2005amytomala-sil/biess-registro-postulacion.git
cd biess-registro-postulacion
```

### 1. Iniciar PostgreSQL y Redis

```bash
docker compose up -d
```

Docker Compose expone PostgreSQL en `localhost:5434` y Redis en `localhost:6379`. En la primera inicialización de PostgreSQL se ejecutan los scripts `.sql` de `db/`, que crean las tablas y cargan los catálogos. Docker solo ejecuta esos scripts cuando el directorio de datos de PostgreSQL está vacío.

Para comprobar los contenedores:

```bash
docker compose ps
```

Para detenerlos sin borrar los datos:

```bash
docker compose down
```

### 2. Configurar y arrancar el backend

Crea `backend/.env` a partir de `backend/.env.example` y revisa los valores si cambiaste los puertos o credenciales.

PowerShell:

```powershell
Copy-Item backend\.env.example backend\.env
Set-Location backend
go run ./cmd/api
```

macOS, Linux o Git Bash:

```bash
cp backend/.env.example backend/.env
cd backend
go run ./cmd/api
```

La API queda disponible en `http://localhost:8080`. Para compilarla desde `backend/`:

```bash
go build -o bin/api ./cmd/api
```

### 3. Configurar y arrancar el frontend

En otra terminal, desde la raíz del repositorio, crea `frontend/.env` desde su ejemplo y ejecuta Vite.

PowerShell:

```powershell
Copy-Item frontend\.env.example frontend\.env
Set-Location frontend
npm ci
npm run dev
```

macOS, Linux o Git Bash:

```bash
cp frontend/.env.example frontend/.env
cd frontend
npm ci
npm run dev
```

Abre `http://localhost:5173`. `VITE_API_URL` debe apuntar a la URL base del backend sin `/api/v1`; el frontend agrega ese prefijo automáticamente.

### Compilación

Para compilar el frontend desde `frontend/`:

```bash
npm ci
npm run build
```

Para compilar el backend desde `backend/`:

```bash
go build -o bin/api ./cmd/api
```

## Variables de entorno

Los ejemplos están separados por servicio:

- [Backend](backend/.env.example): puerto y entorno (`SERVER_PORT`, `ENV`), PostgreSQL (`POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`), Redis (`REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB`), CORS (`CORS_ALLOWED_ORIGINS`) y JWT (`JWT_SECRET_KEY`, `JWT_EXPIRATION_HRS`).
- [Frontend](frontend/.env.example): `VITE_API_URL`, URL base del backend sin `/api/v1`.

Los archivos `.env` están excluidos de Git. No publiques contraseñas ni secretos; sustituye especialmente `JWT_SECRET_KEY` en cualquier entorno compartido.

## Base de datos y datos semilla

- [`db/init.sql`](db/init.sql) contiene el DDL de las tablas y datos semilla de países, provincias y cantones.
- [`db/seed.sql`](db/seed.sql) contiene únicamente los datos semilla; úsalo para completar la carga cuando el esquema ya existe, no para inicializar una base vacía.
- Si los scripts automáticos no se ejecutan porque la base ya tenía datos, puedes cargar `db/init.sql` con `psql` contra una base vacía. En Windows, si aparece un error de codificación, configura primero la sesión de PowerShell:

```powershell
$env:PGCLIENTENCODING = "UTF8"
```

En Render, PostgreSQL no carga los scripts automáticamente. Consulta [`DEPLOY_RENDER.md`](DEPLOY_RENDER.md) para los pasos de inicialización.

## Datos ficticios de prueba

El backend utiliza `RegistroCivilMock`. Para probar el flujo de registro, se pueden usar estas cédulas y sus fechas de expedición de ejemplo:

| Cédula ficticia | Fecha de expedición |
| --- | --- |
| `0101010101` | `01/01/2010` |
| `0202020202` | `15/06/2012` |
| `0404040404` | `10/09/2018` |
| `0505050505` | `05/12/2019` |
| `0606060606` | `25/04/2020` |
| `0909090909` | `01/01/2015` |

La cédula `0303030303` representa un caso de nacionalidad no ecuatoriana y se rechaza. Los datos completos están en [`backend/internal/platform/mocks/registro_civil.go`](backend/internal/platform/mocks/registro_civil.go). No se valida identidad real hasta integrar un proveedor autorizado.

La recuperación de contraseña también es simulada: el token aparece en los logs del backend, no se envía por correo.

## API

El backend publica estas rutas principales:

| Método | Ruta | Descripción |
| --- | --- | --- |
| `GET` | `/health` | Estado de la API, PostgreSQL y Redis |
| `POST` | `/api/v1/auth/verificar-identidad` | Verificación de identidad de prueba |
| `POST` | `/api/v1/auth/registro` | Registro de una cuenta |
| `POST` | `/api/v1/auth/login` | Inicio de sesión |
| `POST` | `/api/v1/auth/recuperar` | Inicio de recuperación de contraseña |
| `POST` | `/api/v1/auth/restablecer` | Cambio de contraseña con token |
| `GET` | `/api/v1/auth/me` | Usuario autenticado (JWT) |
| `GET` | `/api/v1/catalogos/paises` | Catálogo de países |
| `GET` | `/api/v1/catalogos/provincias` | Catálogo de provincias |
| `GET` | `/api/v1/catalogos/provincias/:id/cantones` | Cantones de una provincia |
| `GET`, `PUT` | `/api/v1/postulaciones/me` | Consultar y guardar la postulación autenticada |
| `POST` | `/api/v1/postulaciones/me/enviar` | Enviar la postulación autenticada |

Prueba el estado local desde el navegador o con `curl` en `http://localhost:8080/health`.

## Estructura del repositorio

```text
backend/   API Go, dominio, servicios, handlers y acceso a PostgreSQL/Redis
frontend/  Aplicación React y servicios HTTP
db/        Esquema SQL y catálogos iniciales
render.yaml       Blueprint de Render
docker-compose.yml  PostgreSQL y Redis locales
```

## Despliegue en Render

[`render.yaml`](render.yaml) define PostgreSQL, Key Value, la API Go y el sitio estático. El frontend y la API pueden desplegarse automáticamente al actualizar la rama vinculada. La base PostgreSQL gratuita de Render expira a los 30 días y Key Value gratuito no conserva sus datos al reiniciarse; revisa los límites antes de usar datos reales. Para desplegar e inicializar la base, sigue [`DEPLOY_RENDER.md`](DEPLOY_RENDER.md).

## Autores

- Amy T. Silvestre
- Ismael Anchundia
