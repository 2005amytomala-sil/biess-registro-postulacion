# Desplegar en Render

El archivo `render.yaml` define tres servicios: API en Go, frontend estático y Redis (Render Key Value), además de PostgreSQL.

## Crear los servicios

1. Sube estos cambios al repositorio de GitHub.
2. En Render, crea un **Blueprint** y conecta `2005amytomala-sil/biess-registro-postulacion`.
3. Cuando Render solicite `CORS_ALLOWED_ORIGINS`, ingresa la URL pública del servicio `biess-frontend` (por ejemplo, `https://biess-frontend.onrender.com`).
4. Render genera `JWT_SECRET_KEY` automáticamente y vincula el resto de servicios mediante variables internas.

## Inicializar PostgreSQL

Render crea la base de datos, pero no ejecuta automáticamente `db/init.sql`. Antes de usar el formulario, conecta `db/init.sql` a la base de Render usando su **External Database URL** y `psql`. Verifica que la URL incluya `sslmode=require`:

```bash
psql "<EXTERNAL_DATABASE_URL>" -v ON_ERROR_STOP=1 -f db/init.sql
```

Usa la URL externa solo desde tu equipo; la API queda conectada por la red interna de Render.

## Límites del plan gratuito

Esta configuración usa recursos gratuitos para facilitar una primera publicación. Render indica que la base PostgreSQL gratuita expira a los 30 días y que Key Value gratuito no persiste datos al reiniciar. Para conservar datos y sesiones, cambia estos recursos a planes pagados antes de usar la aplicación en producción.

Al terminar el despliegue, prueba `https://<URL-DE-BIESS-API>/health` y abre la URL del frontend.
