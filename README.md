# Sistema de Registro y Postulación - Directorio del BIESS

Módulo inicial (EF-11 y EF-12) del Concurso Público de Méritos y Oposición para el Directorio del BIESS.

**Stack tecnológico:**
- **Frontend:** React 19 + Vite + Tailwind CSS 4
- **Backend:** Go 1.27 + Gin + PostgreSQL 16 + Redis 7

---

## 📋 Prerrequisitos

Antes de comenzar, asegúrese de tener instaladas las siguientes versiones mínimas:

| Herramienta | Versión mínima | Instalación |
|-------------|---------------|-------------|
| **Go** | 1.27+ | [golang.org/dl](https://golang.org/dl/) |
| **Node.js** | 18+ | [nodejs.org](https://nodejs.org/) |
| **PostgreSQL** | 16+ | [postgresql.org/download](https://www.postgresql.org/download/) |
| **Redis** | 7+ | [redis.io/download](https://redis.io/download) |
| **Docker** | 20.10+ | [docker.com/get-docker](https://www.docker.com/get-docker) |
| **Docker Compose** | 2.0+ | Incluido con Docker Desktop |

---

## 🚀 Ejecución Local con Docker Compose

Esta es la forma recomendada de ejecutar el proyecto en desarrollo, ya que levanta PostgreSQL y Redis automáticamente.

### 1. Clonar el repositorio

```bash
git clone https://github.com/2005amytomala-sil/biess-registro-postulacion.git
cd biess-registro-postulacion
```

### 2. Levantar servicios de infraestructura

```bash
docker-compose up -d
```

Esto iniciará:
- **PostgreSQL** en `localhost:5434` (usuario: `postgres`, password: `postgres`, base de datos: `biess`)
- **Redis** en `localhost:6379`

El script `db/init.sql` se ejecutará automáticamente al crear el contenedor de PostgreSQL, creando todas las tablas (DDL) y sembrando los catálogos de países y provincias del Ecuador.

### 3. Configurar Backend

```bash
cd backend
cp .env.example .env
```

El archivo `.env.example` contiene:

```env
# Server
SERVER_PORT=8080
ENV=development

# PostgreSQL (docker-compose defaults)
POSTGRES_HOST=localhost
POSTGRES_PORT=5434
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=biess

# Redis (docker-compose defaults)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000

# JWT
JWT_SECRET_KEY=secret-key-change-in-production
JWT_EXPIRATION_HRS=24
```

### 4. Instalar dependencias y ejecutar Backend

```bash
go mod download
go run cmd/api/main.go
```

El backend estará disponible en `http://localhost:8080`

### 5. Configurar Frontend

En una nueva terminal:

```bash
cd frontend
cp .env.example .env
```

El archivo `.env.example` contiene:

```env
VITE_API_URL=http://localhost:8080/api/v1
```

### 6. Instalar dependencias y ejecutar Frontend

```bash
npm install
npm run dev
```

El frontend estará disponible en `http://localhost:5173`

---

## 🗄️ Base de Datos

### Ubicación del script de inicialización

`db/init.sql` - Contiene:
- **DDL normalizado (3FN)**: Definición de todas las tablas del sistema
- **Semilla de división política**: Catálogos de países y las 24 provincias del Ecuador

### Estructura principal de tablas

- `postulantes` - Usuarios registrados en el sistema
- `postulaciones` - Registro de postulaciones (estado: borrador, enviada, cerrada)
- `datos_personales` - Datos obtenidos del Registro Civil
- `domicilio_contacto` - Dirección y medios de contacto
- `familiares_iess` - Familiares que trabajan en IESS/BIESS
- `respaldo_organizacion` - Organización que respalda la postulación
- `formacion_academica` - Estudios realizados
- `experiencia_profesional` - Trayectoria laboral
- `inhabilidades` - Declaración de inhabilidades legales
- `paises`, `provincias`, `cantones` - Catálogos geográficos

---

## 👤 Usuarios de Prueba (Mock Registro Civil)

El sistema incluye un mock del servicio del Registro Civil para desarrollo. Las siguientes cédulas están preconfiguradas con datos de prueba:

| Cédula | Nombres | Apellidos | Nacionalidad | Estado Civil | Observaciones |
|--------|---------|-----------|-------------|--------------|---------------|
| `0101010101` | Juan Carlos | Pérez López | Ecuatoriana | Soltero | Caso válido básico |
| `0202020202` | María Isabel | García Martínez | Ecuatoriana | Soltera | Caso válido femenino |
| `0303030303` | John | Smith | Estadounidense | Soltero | **Nacionalidad no permitida** (rechazado) |
| `0404040404` | Carlos Eduardo | Rodríguez Torres | Ecuatoriana | Casado | Con datos de pareja |
| `0505050505` | Ana María | Gómez Sánchez | Ecuatoriana | Casada | Con datos de pareja |
| `0606060606` | Luis Fernando | Mendoza Ruiz | Ecuatoriana | Divorciado | Caso válido divorciado |
| `0909090909` | Carlos | Test Prueba | Ecuatoriana | Soltero | Para pruebas de bloqueo |

**Códigos dactilares y fechas de expedición** se encuentran en `backend/internal/platform/mocks/registro_civil.go`

---

## 🔌 API Endpoints

### Autenticación

| Método | Endpoint | Descripción | Autenticación |
|--------|----------|-------------|---------------|
| POST | `/api/v1/auth/verificar-identidad` | Verificar identidad con Registro Civil | No |
| POST | `/api/v1/auth/registro` | Registrar nuevo usuario | No |
| POST | `/api/v1/auth/login` | Iniciar sesión | No |
| POST | `/api/v1/auth/recuperar` | Solicitar recuperación de contraseña | No |
| POST | `/api/v1/auth/restablecer` | Restablecer contraseña con token | No |
| GET | `/api/v1/auth/me` | Obtener usuario autenticado | Sí (JWT) |

### Catálogos

| Método | Endpoint | Descripción | Autenticación |
|--------|----------|-------------|---------------|
| GET | `/api/v1/catalogos/paises` | Obtener lista de países | No |
| GET | `/api/v1/catalogos/provincias` | Obtener lista de provincias | No |
| GET | `/api/v1/catalogos/provincias/:id/cantones` | Obtener cantones por provincia | No |

### Postulaciones

| Método | Endpoint | Descripción | Autenticación |
|--------|----------|-------------|---------------|
| GET | `/api/v1/postulaciones/me` | Obtener postulación del usuario | Sí (JWT) |
| PUT | `/api/v1/postulaciones/me` | Guardar progreso de postulación | Sí (JWT) |
| POST | `/api/v1/postulaciones/me/enviar` | Enviar postulación final | Sí (JWT) |

### Health Check

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/health` | Verificar estado de PostgreSQL y Redis |

---

## 🏗️ Arquitectura

### Diagrama de Capas (Backend)

```
┌─────────────────────────────────────────────────────────────┐
│                      HTTP Requests                           │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                      Handler Layer                          │
│  (Gin Router + HTTP Handlers)                              │
│  - Parse request body                                       │
│  - Validate input                                           │
│  - Call service layer                                       │
│  - Return HTTP response                                     │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                      Service Layer                           │
│  (Business Logic)                                          │
│  - Domain validation                                        │
│  - Call external services (Registro Civil mock)             │
│  - Coordinate repositories                                  │
│  - Manage transactions                                      │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                   Repository Layer                          │
│  (Data Access)                                             │
│  - SQL queries (pgx)                                       │
│  - Database operations                                     │
│  - Cache integration (Redis)                               │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│              PostgreSQL & Redis                              │
└─────────────────────────────────────────────────────────────┘
```

### Redis: Caché y Bloqueo

#### 1. Caché de Catálogos

Los catálogos (países, provincias, cantones) se cachean en Redis para reducir consultas a PostgreSQL:

```
┌──────────────┐      GET      ┌──────────────┐
│   Frontend   │ ────────────> │   Handler    │
└──────────────┘               └──────┬───────┘
                                     │
                                     ▼
                            ┌──────────────┐
                            │   Service    │
                            └──────┬───────┘
                                   │
                ┌──────────────────┼──────────────────┐
                │                  │                  │
                ▼                  ▼                  ▼
        ┌──────────────┐   ┌──────────────┐   ┌──────────────┐
        │   Redis      │   │ Repository   │   │    SQL       │
        │   Cache      │   │              │   │ PostgreSQL   │
        │ (TTL: 24h)   │   │              │   │              │
        └──────────────┘   └──────────────┘   └──────────────┘
```

**Claves Redis:**
- `catalogo:paises` - Lista de países (JSON)
- `catalogo:provincias` - Lista de provincias (JSON)
- `catalogo:cantones:{provincia_id}` - Cantones por provincia (JSON)

#### 2. Bloqueo por Intentos Fallidos (Rate Limiting)

El sistema implementa rate limiting en el flujo de autenticación para prevenir ataques de fuerza bruta:

```
┌──────────────┐     POST /auth/login     ┌──────────────┐
│   Frontend   │ ───────────────────────> │   Handler    │
└──────────────┘                         └──────┬───────┘
                                               │
                                               ▼
                                      ┌──────────────┐
                                      │   Service    │
                                      └──────┬───────┘
                                             │
                     ┌───────────────────────┼───────────────────────┐
                     │                       │                       │
                     ▼                       ▼                       ▼
           ┌──────────────┐         ┌──────────────┐        ┌──────────────┐
           │ Redis Limiter│         │  Repository  │        │  Bloqueo?    │
           │              │         │              │        │              │
           │ auth:intentos│         │              │        │ auth:bloqueo │
           │ :{cedula}    │         │              │        │ :{cedula}    │
           │              │         │              │        │              │
           │ Max: 3       │         │              │        │ TTL: 10min   │
           │ TTL: 15min   │         │              │        │              │
           └──────────────┘         └──────────────┘        └──────────────┘
```

**Lógica de bloqueo:**
1. Cada intento fallido incrementa el contador en `auth:intentos:{cedula}`
2. Después de 3 intentos fallidos, se establece un bloqueo en `auth:bloqueo:{cedula}` por 10 minutos
3. El contador de intentos expira en 15 minutos
4. Si el usuario se autentica correctamente, se reinician los intentos

**Claves Redis:**
- `auth:intentos:{cedula}` - Contador de intentos (TTL: 15 min)
- `auth:bloqueo:{cedula}` - Indicador de bloqueo (TTL: 10 min)

---

## 📁 Estructura del Proyecto

```
biess-registro-postulacion/
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go          # Punto de entrada
│   ├── internal/
│   │   ├── domain/              # Entidades y errores de dominio
│   │   ├── handler/             # HTTP handlers (Gin)
│   │   ├── middleware/          # JWT, CORS, rate limiting
│   │   ├── repository/          # Acceso a datos (PostgreSQL)
│   │   ├── service/             # Lógica de negocio
│   │   └── platform/            # Infraestructura (Redis, mocks)
│   ├── .env.example
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── components/          # Componentes React
│   │   ├── pages/               # Páginas (Dashboard, Postulación)
│   │   ├── services/            # API clients
│   │   └── utils/               # Utilidades
│   ├── .env.example
│   └── package.json
├── db/
│   └── init.sql                 # DDL + semilla de división política
├── docker-compose.yml
└── README.md
```

---

## 🧪 Ejecutar Tests

### Backend

```bash
cd backend
go test ./...
```

### Tests de integración (requiere PostgreSQL corriendo)

```bash
cd backend
go test -tags=integration ./...
```

---

## 🐛 Troubleshooting

### Puerto 8080 ya está en uso

```bash
# En Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# En Linux/Mac
lsof -ti:8080 | xargs kill -9
```

### Puerto 5173 ya está en uso (Vite)

Si el puerto 5173 está ocupado, Vite automáticamente intentará el siguiente puerto disponible (5174, 5175, etc.). Si necesitas liberar el puerto:

```bash
# En Windows
netstat -ano | findstr :5173
taskkill /PID <PID> /F

# En Linux/Mac
lsof -ti:5173 | xargs kill -9
```

### Docker compose no levanta los contenedores

```bash
# Verificar logs
docker-compose logs

# Reconstruir contenedores
docker-compose down -v
docker-compose up -d --build
```

### Error de conexión a PostgreSQL

Verificar que el puerto 5434 esté disponible y que docker-compose esté corriendo:

```bash
docker-compose ps
```

### Dependencias de Node.js no se instalan

Si encuentras errores al instalar dependencias, intenta limpiar la caché:

```bash
cd frontend
rm -rf node_modules package-lock.json
npm cache clean --force
npm install
```

---

## 📝 Licencia

Este proyecto es propiedad del BIESS (Instituto Ecuatoriano de Seguridad Social).

---

## 👥 Autores

- [Amy T. Silvestre]
- [Ismael Anchundia]
