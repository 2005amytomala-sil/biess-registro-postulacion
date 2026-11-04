# Sistema de Registro y Postulación del BIESS

Aplicación web para el registro de personas postulantes y la gestión de información del Concurso Público de Méritos y Oposición. Es un prototipo académico: la verificación de identidad usa datos simulados y no está conectada al Registro Civil.

## Tecnologías y versiones

| Componente | Tecnología | Versión requerida o usada |
| --- | --- | --- |
| Backend | Go, Gin | Go 1.27.1 (según `go.mod`) |
| Frontend | React, Vite | Node.js 20.19+ o 22.12+; npm |
| Base de datos | PostgreSQL | 16 (imagen de Docker Compose) |
| Caché y sesiones | Redis / Render Key Value | Redis 7 (imagen de Docker Compose) |

Vite 8 requiere Node.js 20.19+ o 22.12+. Consulta la [guía oficial de Vite](https://vite.dev/guide/) si necesitas comprobar requisitos actualizados.

## Prerrequisitos

- Git.
- Go 1.27.1.
- Node.js 20.19+ o 22.12+ y npm.
- Docker Desktop o Docker Engine con Docker Compose v2. También puedes usar PostgreSQL 16 y Redis 7 instalados localmente.
- `psql` es opcional; se necesita para cargar el SQL manualmente desde la terminal.

## Instalación y ejecución local

Clona el repositorio y entra en la carpeta del proyecto:

```bash
git clone https://github.com/2005amytomala-sil/biess-registro-postulacion.git
cd biess-registro-postulacion
```

### 1. Iniciar PostgreSQL y Redis

```bash
docker compose up -d
```

Docker Compose publica PostgreSQL en `localhost:5434` y Redis en `localhost:6379`. En la primera creación de la base, PostgreSQL ejecuta los archivos SQL de `db/` y crea las tablas y catálogos. Los scripts de inicialización automáticos solo se ejecutan cuando el directorio de datos de PostgreSQL está vacío.

### 2. Configurar y arrancar el backend

Las variables de ejemplo están en `backend/.env.example`. Crea `backend/.env` a partir de ese archivo y ajústalo si cambiaste los puertos o credenciales de los servicios.

En PowerShell:

```powershell
Copy-Item backend\.env.example backend\.env
Set-Location backend
go run ./cmd/api
```

En macOS, Linux o Git Bash:

```bash
cp backend/.env.example backend/.env
cd backend
go run ./cmd/api
```

El backend escucha por defecto en `http://localhost:8080`. Para compilarlo desde la carpeta `backend`:

```bash
go build -o bin/api ./cmd/api
```

### 3. Configurar y arrancar el frontend

En otra terminal, desde la raíz del repositorio, crea `frontend/.env` y arranca Vite.

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

Abre `http://localhost:5173`. La variable `VITE_API_URL` del ejemplo apunta al backend local. El frontend agrega `/api/v1` a esa URL.

## Variables de entorno

Los ejemplos están separados por servicio:

- [Backend](backend/.env.example): `SERVER_PORT`, `ENV`, `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB`, `CORS_ALLOWED_ORIGINS`, `JWT_SECRET_KEY` y `JWT_EXPIRATION_HRS`.
- [Frontend](frontend/.env.example): `VITE_API_URL` (URL base del backend, sin `/api/v1`).

No subas archivos `.env` ni contraseñas reales al repositorio. Los valores de `.env.example` son solo para desarrollo local; reemplaza especialmente `JWT_SECRET_KEY` en cualquier entorno compartido.

## Base de datos y datos semilla

- [`db/init.sql`](db/init.sql) contiene el DDL de las tablas y datos semilla de países, provincias y cantones.
- [`db/seed.sql`](db/seed.sql) contiene solo los datos semilla; sirve para completar la carga cuando el esquema ya existe. No lo uses como sustituto de `init.sql` en una base nueva.
- En Render, la base no ejecuta estos archivos automáticamente. Sigue [`DEPLOY_RENDER.md`](DEPLOY_RENDER.md) para inicializarla.

Si cargas un archivo desde Windows con `psql` y aparece un error de codificación, configura primero el cliente como UTF-8 en PowerShell:

```powershell
$env:PGCLIENTENCODING = "UTF8"
```

## Comprobación y compilación

Con backend y servicios iniciados, consulta el estado:

```text
http://localhost:8080/health
```

Compilar el frontend desde `frontend/`:

```bash
npm ci
npm run build
```

Ejecutar las pruebas Go desde la raíz del repositorio:

```bash
go test ./...
```

## API

- `GET /health`: estado del backend, PostgreSQL y Redis.
- Las rutas de autenticación usan `/api/v1/auth`.
- Los catálogos usan `/api/v1/catalogos`.
- Las postulaciones usan `/api/v1/postulaciones`.

## Despliegue en Render

El [Blueprint](render.yaml) define PostgreSQL, Key Value, una API Go y el sitio estático de React. Los servicios están configurados para el plan gratuito, que tiene límites y no es adecuado para datos de producción. Render no carga `db/init.sql` automáticamente; inicializa la base siguiendo [`DEPLOY_RENDER.md`](DEPLOY_RENDER.md).

## Alcance de la verificación de identidad

El backend actual instancia `RegistroCivilMock`, con identidades y fechas ficticias definidas en [`registro_civil.go`](backend/internal/platform/mocks/registro_civil.go). Para probar el flujo, usa uno de esos registros de prueba; por ejemplo, la cédula `0101010101` con fecha de expedición `01/01/2010`. El sistema no puede verificar identidades reales hasta integrar un proveedor autorizado.

La recuperación de contraseña también es una simulación: el token se escribe en los logs del backend en lugar de enviarse por correo.
