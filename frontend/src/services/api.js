import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api',
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

// Interceptor de respuesta: maneja errores específicos
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

    switch (status) {
      case 401:
        // No autorizado: cerrar sesión
        localStorage.removeItem('token');
        localStorage.removeItem('usuario');
        window.location.href = '/login';
        return Promise.reject({
          message: 'Su sesión ha expirado. Por favor, inicie sesión nuevamente.',
          originalError: error,
        });

      case 429:
        // Demasiadas solicitudes: exponer retry_after
        const retryAfter = data?.retry_after || 60;
        return Promise.reject({
          message: `Demasiados intentos. Intente nuevamente en ${retryAfter} segundos.`,
          retryAfter,
          originalError: error,
        });

      case 500:
        return Promise.reject({
          message: 'Error interno del servidor. Intente nuevamente más tarde.',
          originalError: error,
        });

      case 503:
        return Promise.reject({
          message: 'El servicio no está disponible temporalmente.',
          originalError: error,
        });

      default:
        // Otros errores: usar mensaje del servidor si existe
        const serverMessage = data?.message || data?.error;
        return Promise.reject({
          message: serverMessage || 'Ocurrió un error inesperado.',
          originalError: error,
        });
    }
  }
);

export default api;
