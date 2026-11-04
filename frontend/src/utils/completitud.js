// Refleja la verificación de completitud del backend (EnviarPostulacion) para habilitar el envío.
// El backend sigue siendo la fuente de verdad: si responde 422 se muestran sus faltantes.

export const CODIGOS_INHABILIDADES = [
  '1.9.1', '1.9.2', '1.9.3', '1.9.4', '1.9.5', '1.9.6', '1.9.7',
  '1.9.8', '1.9.9', '1.9.10', '1.9.11', '1.9.12', '1.9.13',
];

const SECTORES = ['afiliados', 'jubilados', 'empleadores'];

// Paso del formulario donde se corrige cada sección faltante
export const PASO_POR_SECCION = {
  datos_personales: 0,
  domicilio: 1,
  sector: 3,
  respaldo_organizacion: 3,
  inhabilidades: 6,
  consentimiento: 7,
};

export const calcularFaltantes = (postulacion) => {
  const faltantes = [];
  const agregar = (seccion, mensaje) => faltantes.push({ seccion, mensaje });

  if (!postulacion?.datos_personales) agregar('datos_personales', 'Datos personales del Registro Civil');
  if (!postulacion?.domicilio) agregar('domicilio', 'Domicilio y datos de contacto');
  if (!SECTORES.includes(postulacion?.sector)) agregar('sector', 'Sector al que postula');

  const respaldo = postulacion?.respaldo_organizacion;
  if (!respaldo) {
    agregar('respaldo_organizacion', 'Respaldo de organización');
  } else if (respaldo.tiene_respaldo && !respaldo.nombre_organizacion?.trim()) {
    agregar('respaldo_organizacion', 'Nombre de la organización que respalda la postulación');
  }

  const respondidas = new Set((postulacion?.inhabilidades || []).map((i) => i.codigo_pregunta));
  const sinResponder = CODIGOS_INHABILIDADES.filter((c) => !respondidas.has(c));
  if (sinResponder.length > 0) {
    agregar('inhabilidades', `Inhabilidades sin responder: ${sinResponder.join(', ')}`);
  }

  return faltantes;
};
