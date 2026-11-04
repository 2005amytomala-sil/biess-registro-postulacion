import api from './api';

// Normaliza el error del interceptor: { success:false, message, status, codigo, data }
const fallo = (error, mensajePorDefecto) => ({
  success: false,
  message: error.message || mensajePorDefecto,
  status: error.status,
  codigo: error.codigo,
  data: error.data || {},
});

// Obtener postulación del usuario actual
export const obtenerPostulacion = async () => {
  try {
    const response = await api.get('/postulaciones/me');
    return { success: true, data: response.data };
  } catch (error) {
    return fallo(error, 'Error al obtener postulación');
  }
};

// Guardar postulación (guardado parcial: solo las secciones enviadas)
export const guardarPostulacion = async (datos) => {
  try {
    const response = await api.put('/postulaciones/me', datos);
    return { success: true, data: response.data };
  } catch (error) {
    return fallo(error, 'Error al guardar postulación');
  }
};

// Enviar postulación (el backend valida completitud y consentimiento)
export const enviarPostulacion = async (aceptaConsentimiento) => {
  try {
    const response = await api.post('/postulaciones/me/enviar', { acepta_consentimiento: aceptaConsentimiento });
    return { success: true, data: response.data };
  } catch (error) {
    return fallo(error, 'Error al enviar postulación');
  }
};

// Obtener catálogo de países
export const obtenerPaises = async () => {
  try {
    const response = await api.get('/catalogos/paises');
    return { success: true, data: response.data };
  } catch (error) {
    return fallo(error, 'Error al obtener países');
  }
};

// Obtener catálogo de provincias
export const obtenerProvincias = async () => {
  try {
    const response = await api.get('/catalogos/provincias');
    return { success: true, data: response.data };
  } catch (error) {
    return fallo(error, 'Error al obtener provincias');
  }
};

// Obtener cantones por provincia
export const obtenerCantonesPorProvincia = async (provinciaId) => {
  try {
    const response = await api.get(`/catalogos/provincias/${provinciaId}/cantones`);
    return { success: true, data: response.data };
  } catch (error) {
    return fallo(error, 'Error al obtener cantones');
  }
};
