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

Docker Compose expone PostgreSQL en `localhost:5434` (usuario `postgres`, contraseña `postgres`, base `biess`) y Redis en `localhost:6379`. En la primera inicialización de PostgreSQL se ejecutan en orden alfabético los scripts de `db/`: `init.sql` (tablas y catálogos) y luego `seed.sql` (catálogos otra vez; es idempotente y no duplica datos). Docker solo ejecuta esos scripts cuando el volumen de datos de PostgreSQL está vacío.

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

- [`db/init.sql`](db/init.sql) contiene el DDL de las tablas y los datos semilla de países (194), provincias (24) y cantones (199). Debe ejecutarse sobre una base vacía.
- [`db/seed.sql`](db/seed.sql) contiene solo los datos semilla. Es idempotente (`ON CONFLICT DO NOTHING`), así que puede ejecutarse varias veces o sobre una base que ya tiene el esquema.
- Ambos scripts sincronizan las secuencias de `paises` y `provincias` después de insertar ids explícitos.
- Ecuador es el país con id `1`; la restricción `CHECK` de `domicilio_contacto` depende de ese id.

Si los scripts automáticos no se ejecutaron porque el volumen ya tenía datos, puedes recrear la base desde cero (borra los datos locales):

```powershell
docker compose down -v
docker compose up -d
```

O cargarlos manualmente con `psql` contra una base vacía:

```powershell
$env:PGCLIENTENCODING = "UTF8"
psql -h localhost -p 5434 -U postgres -d biess -v ON_ERROR_STOP=1 -f db/init.sql
psql -h localhost -p 5434 -U postgres -d biess -v ON_ERROR_STOP=1 -f db/seed.sql
```

En Render, PostgreSQL no carga los scripts automáticamente. Consulta [`DEPLOY_RENDER.md`](DEPLOY_RENDER.md) para los pasos de inicialización.

## Datos ficticios de prueba

El backend usa `RegistroCivilMock` y `HistorialLaboralMock` (datos ficticios, sin conexión a servicios reales). Para registrarse se ingresa la cédula y **uno** de los dos datos de identidad: fecha de expedición o código dactilar (10 dígitos).

| Cédula | Fecha de expedición | Código dactilar | Caso que prueba |
| --- | --- | --- | --- |
| `0101010101` | `01/01/2010` | `1234567890` | Válida: ecuatoriano, soltero. Flujo completo de registro, login, postulación y envío. |
| `0202020202` | `15/06/2012` | `0987654321` | Válida: ecuatoriana, soltera. Segundo usuario para el flujo completo. |
| `0303030303` | `20/03/2015` | `1122334455` | Nacionalidad no ecuatoriana (estadounidense): el registro se rechaza con "La información ingresada es incorrecta". Si llegara a postular, la postulación se cierra por nacionalidad. |
| `0404040404` | `10/09/2018` | `5566778899` | Casado: muestra los datos de la pareja (`0505050505`). Al elegir el sector **Jubilados** la postulación se cierra con el mensaje del Historial Laboral. |
| `0505050505` | `05/12/2019` | `9988776655` | Casada: muestra los datos de la pareja (`0404040404`). Al elegir el sector **Empleadores** la postulación se cierra con el mensaje del Historial Laboral. |
| `0606060606` | `25/04/2020` | `4455667788` | Válida: ecuatoriano, divorciado. |
| `0909090909` | `01/01/2015` | `1111222233` | Bloqueo: tres intentos fallidos de verificación (por ejemplo, fecha `01/01/1999`) bloquean el registro 10 minutos ("Registro fallido, inténtelo nuevamente en 10 minutos"), con cuenta regresiva en pantalla. El login tiene su propio contador: tres contraseñas erróneas lo bloquean 10 minutos. |

Los datos completos están en [`registro_civil.go`](backend/internal/platform/mocks/registro_civil.go) y [`historial_laboral.go`](backend/internal/platform/mocks/historial_laboral.go). No se valida identidad real hasta integrar un proveedor autorizado.

## Reglas de la postulación (EF-12)

- Antes de consultar, guardar o enviar, el backend verifica la nacionalidad (Registro Civil, con la cédula del token) y el sector (Historial Laboral). Si alguno no es válido, la postulación pasa a `cerrada` con su `motivo_cierre` y se responde 422.
- Una postulación `cerrada` o `enviada` no admite cambios: guardar o enviar responde 409 con `motivo_cierre`.
- El guardado es parcial: cada sección enviada en `PUT /postulaciones/me` modifica solo su tabla; las secciones ausentes no se tocan.
- Para enviar (`POST /postulaciones/me/enviar` con `{"acepta_consentimiento": true}`) deben estar completos los datos personales, el domicilio, el sector y respaldo, y las 13 inhabilidades. Si falta algo, se responde 422 con la lista `faltantes`.

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

Prueba el estado local desde el navegador en `http://localhost:8080/health` o con PowerShell:

```powershell
Invoke-RestMethod http://localhost:8080/health | ConvertTo-Json
```

## Pruebas

Desde `backend/`:

```powershell
go build ./...
go test ./...
```

Las pruebas de integración usan una base PostgreSQL real con el esquema cargado. Por ejemplo, con un contenedor temporal:

```powershell
docker run --rm -d --name biess-test-db -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=biess_test -p 55432:5432 postgres:16
# Esperar unos segundos y cargar el esquema
docker cp ..\db\init.sql biess-test-db:/tmp/init.sql
docker exec biess-test-db psql -q -v ON_ERROR_STOP=1 -U postgres -d biess_test -f /tmp/init.sql
$env:TEST_DB_PORT = "55432"; $env:TEST_DB_NAME = "biess_test"
go test -tags=integration ./internal/service/
docker stop biess-test-db
```

Desde `frontend/`:

```powershell
npm run build
```

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
