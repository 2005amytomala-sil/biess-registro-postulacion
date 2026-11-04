import { createContext, useContext, useState, useEffect } from 'react';
import api from '../services/api';

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
  const [usuario, setUsuario] = useState(null);
  const [token, setToken] = useState(null);
  const [loading, setLoading] = useState(true);

  // Restaurar sesión al cargar la aplicación
  useEffect(() => {
    const storedToken = localStorage.getItem('token');
    const storedUsuario = localStorage.getItem('usuario');

    if (storedToken && storedUsuario) {
      try {
        setToken(storedToken);
        setUsuario(JSON.parse(storedUsuario));
      } catch (error) {
        console.error('Error al restaurar sesión:', error);
        localStorage.removeItem('token');
        localStorage.removeItem('usuario');
      }
    }

    setLoading(false);
  }, []);

  const login = async (cedula, password) => {
    try {
      const response = await api.post('/auth/login', { cedula, password });
      const { token: newToken, user: newUsuario } = response.data;

      setToken(newToken);
      setUsuario(newUsuario);

      // Persistir en localStorage
      localStorage.setItem('token', newToken);
      localStorage.setItem('usuario', JSON.stringify(newUsuario));

      return { success: true };
    } catch (error) {
      return {
        success: false,
        message: error.message || 'Error al iniciar sesión',
        retryAfter: error.retryAfter,
      };
    }
  };

  const logout = () => {
    setToken(null);
    setUsuario(null);
    localStorage.removeItem('token');
    localStorage.removeItem('usuario');
  };

  // El registro no inicia sesión: el usuario debe ingresar luego por /login
  const register = async (datos) => {
    try {
      const response = await api.post('/auth/registro', datos);
      return { success: true, data: response.data };
    } catch (error) {
      return {
        success: false,
        message: error.message || 'Error al registrar usuario',
        retryAfter: error.retryAfter,
      };
    }
  };

  const verificarIdentidad = async (cedula, fechaExpedicion, codigoDactilar) => {
    try {
      const response = await api.post('/auth/verificar-identidad', {
        cedula,
        fecha_expedicion: fechaExpedicion,
        codigo_dactilar: codigoDactilar,
      });
      return { success: true, data: response.data };
    } catch (error) {
      return {
        success: false,
        message: error.message || 'Error al verificar identidad',
        retryAfter: error.retryAfter,
      };
    }
  };

  const recuperarPassword = async (cedula, fechaExpedicion, codigoDactilar) => {
    try {
      const response = await api.post('/auth/recuperar', {
        cedula,
        fecha_expedicion: fechaExpedicion,
        codigo_dactilar: codigoDactilar,
      });
      return { success: true, data: response.data };
    } catch (error) {
      return {
        success: false,
        message: error.message || 'Error al recuperar contraseña',
        retryAfter: error.retryAfter,
      };
    }
  };

  const restablecerPassword = async (cedula, token, password, confirmPassword) => {
    try {
      const response = await api.post('/auth/restablecer', {
        cedula,
        token,
        password,
        confirm_password: confirmPassword,
      });
      return { success: true, data: response.data };
    } catch (error) {
      return {
        success: false,
        message: error.message || 'Error al restablecer contraseña',
      };
    }
  };

  const isAuthenticated = !!token;

  const value = {
    usuario,
    token,
    loading,
    isAuthenticated,
    login,
    logout,
    register,
    verificarIdentidad,
    recuperarPassword,
    restablecerPassword,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth debe usarse dentro de un AuthProvider');
  }
  return context;
};
