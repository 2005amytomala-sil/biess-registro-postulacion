Lee CONTEXTO.md antes de hacer cualquier cosa.
 
PROYECTO: Sistema de Registro y Postulación al Concurso Público de Méritos y Oposición del
Directorio del BIESS (módulo inicial: EF-11 y EF-12 del archivo requerimiento.pdf, que debes
tener adjunto). Sustento: Resolución C.D. 701 del IESS.
 
STACK OBLIGATORIO:
- Frontend: SPA con React + Vite + Tailwind CSS.
- Backend: Go con Gin. Contraseñas con bcrypt (golang.org/x/crypto/bcrypt).
- Persistencia: PostgreSQL (driver pgx). Caché y control de intentos: Redis (go-redis).
- Idioma de la interfaz y de los mensajes: español.
 
ARQUITECTURA BACKEND (OBLIGATORIA): capas Handler -> Service -> Repository con INYECCIÓN DE
DEPENDENCIAS por constructores e interfaces (definidas en la capa que las consume).
- Sin variables globales ni singletons: todo el cableado se hace en cmd/api/main.go.
- Los handlers NO contienen lógica de negocio. Los repositories NO contienen reglas de negocio.
- Errores de dominio tipados (errors.Is / errors.As, wrapping con %w), mapeados a códigos HTTP
  en un solo lugar.
- context.Context en todas las llamadas, con timeouts. Apagado ordenado (graceful shutdown).
- Estructura:
    backend/cmd/api/main.go
    backend/internal/config
    backend/internal/domain        (entidades, errores, interfaces)
    backend/internal/handler
    backend/internal/service
    backend/internal/repository
    backend/internal/platform/postgres y platform/redis
    backend/internal/middleware    (jwt, cors, logging, recover)
    db/init.sql                    (DDL + datos semilla)
 
ARQUITECTURA FRONTEND: src/{components,pages,context,hooks,services,routes}. Estado global con
Context API (AuthContext) y estado local en formularios. Consumo de API con async/await
(axios) manejando siempre loading / success / error. Rutas protegidas por sesión.
Validaciones inmediatas, loaders y mensajes de error visibles. Diseño responsivo con Tailwind.
 
REGLAS DE NEGOCIO CLAVE:
- Registro Civil e Historial Laboral IESS son servicios MOCK detrás de una interfaz Go
  (se reemplazarán por los reales en el futuro).
- Identidad: cédula + (fecha de expedición DD/MM/AAAA O código dactilar de 10 dígitos).
- Bloqueo: al 3.er intento fallido consecutivo, bloqueo de 10 minutos controlado con Redis
  (TTL). Mensaje: "Registro fallido, inténtelo nuevamente en 10 minutos".
- Solo nacionalidad ecuatoriana; si no lo es, se cierra la postulación.
- Cascada País -> Provincia -> Cantón (provincia y cantón por lista solo si el país es Ecuador).
- Los catálogos geográficos se cachean en Redis con TTL.
- El guardado de la postulación es transaccional (BEGIN / COMMIT / ROLLBACK).
 
REGLAS DE TRABAJO (OBLIGATORIAS):
1. Trabaja solo en el alcance de la fase indicada. No modifiques archivos de otras áreas
   salvo que se indique.
2. NO HAGAS COMMITS, NO HAGAS PUSH, NO CREES NI CAMBIES DE RAMA, NO HAGAS MERGE. El control
   de git es exclusivamente mío. No ejecutes git add / commit / push / checkout / merge.
3. Divide cada fase en BLOQUES. Al terminar cada bloque, DETENTE y entrégame un REPORTE.
   No avances al siguiente bloque hasta que yo te diga "continúa".
4. Formato del REPORTE:
   a) Qué implementé (lista breve, en lenguaje claro).
   b) Archivos creados y modificados (rutas completas).
   c) Cómo probarlo (comandos y pasos exactos, datos de prueba).
   d) Decisiones o supuestos que tomé.
   e) Qué quedó pendiente o con dudas.
   f) Mensaje de commit sugerido (formato: tipo(alcance): descripción).
5. Si algo del requerimiento es ambiguo, pregúntame antes de inventar.
6. El código debe compilar y ejecutarse al final de cada bloque (go build ./... y npm run build).

