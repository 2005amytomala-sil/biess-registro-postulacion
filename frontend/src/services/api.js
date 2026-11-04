import axios from 'axios';

const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';
const api = axios.create({
  baseURL: `${apiUrl.replace(/\/$/, '')}/api/v1`,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor de request: agrega el token JWT si existe
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Interceptor de respuesta: normaliza los errores a { message, status, codigo, data, retryAfter }
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (!error.response) {
      // Error de red (servidor no responde)
      return Promise.reject({
        message: 'No se pudo conectar con el servidor. Verifique su conexión a internet.',
        originalError: error,
      });
    }

    const { status, data } = error.response;
    const serverMessage = data?.error || data?.message;
    const base = { status, codigo: data?.codigo, data: data || {}, originalError: error };

    switch (status) {
      case 401: {
        // No autorizado: solo cerrar sesión y redirigir si el usuario tiene una sesión activa
        if (localStorage.getItem('token')) {
          localStorage.removeItem('token');
          localStorage.removeItem('usuario');
          window.location.href = '/login';
          return Promise.reject({ ...base, message: 'Su sesión ha expirado. Por favor, inicie sesión nuevamente.' });
        }
        // Si no hay sesión (ej. registro), usar mensaje del servidor
        return Promise.reject({ ...base, message: serverMessage || 'Error de autenticación' });
      }

      case 429: {
        // Bloqueo: exponer el mensaje del servidor y los segundos restantes reales
        const retryAfter = Number(data?.retry_after || error.response.headers?.['retry-after']) || 600;
        return Promise.reject({
          ...base,
          message: serverMessage || 'Demasiados intentos. Intente nuevamente más tarde.',
          retryAfter,
        });
      }

      case 500:
        return Promise.reject({ ...base, message: 'Error interno del servidor. Intente nuevamente más tarde.' });

      case 503:
        return Promise.reject({ ...base, message: 'El servicio no está disponible temporalmente.' });

      default:
        // Otros errores (400, 404, 409, 422...): usar mensaje del servidor si existe
        return Promise.reject({ ...base, message: serverMessage || 'Ocurrió un error inesperado.' });
    }
  }
);

export default api;
