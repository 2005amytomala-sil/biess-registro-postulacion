import api from './api';

// Obtener postulación del usuario actual
export const obtenerPostulacion = async () => {
  try {
    const response = await api.get('/postulaciones/me');
    return { success: true, data: response.data };
  } catch (error) {
    return {
      success: false,
      message: error.response?.data?.error || 'Error al obtener postulación',
    };
  }
};

// Guardar postulación (guardado de progreso por sección)
export const guardarPostulacion = async (datos) => {
  try {
    const response = await api.put('/postulaciones/me', datos);
    return { success: true, data: response.data };
  } catch (error) {
    return {
      success: false,
      message: error.response?.data?.error || 'Error al guardar postulación',
      codigo: error.response?.data?.codigo,
    };
  }
};

// Enviar postulación (valida completitud y cambia estado a "enviada")
export const enviarPostulacion = async () => {
  try {
    const response = await api.post('/postulaciones/me/enviar');
    return { success: true, data: response.data };
  } catch (error) {
    return {
      success: false,
      message: error.response?.data?.error || 'Error al enviar postulación',
    };
  }
};

// Obtener catálogo de países
export const obtenerPaises = async () => {
  try {
    const response = await api.get('/catalogos/paises');
    return { success: true, data: response.data };
  } catch (error) {
    return {
      success: false,
      message: 'Error al obtener países',
    };
  }
};

// Obtener catálogo de provincias
export const obtenerProvincias = async () => {
  try {
    const response = await api.get('/catalogos/provincias');
    return { success: true, data: response.data };
  } catch (error) {
    return {
      success: false,
      message: 'Error al obtener provincias',
    };
  }
};

// Obtener cantones por provincia
export const obtenerCantonesPorProvincia = async (provinciaId) => {
  try {
    const response = await api.get(`/catalogos/provincias/${provinciaId}/cantones`);
    return { success: true, data: response.data };
  } catch (error) {
    return {
      success: false,
      message: 'Error al obtener cantones',
    };
  }
};
