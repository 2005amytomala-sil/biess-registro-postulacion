/**
 * Calcula la duración entre dos fechas en años, meses y días
 * Considera años bisiestos correctamente
 * @param {string} fechaInicio - Fecha de inicio en formato YYYY-MM-DD
 * @param {string} fechaFin - Fecha de fin en formato YYYY-MM-DD (o null si está en curso)
 * @returns {Object} Objeto con { anos, meses, dias, totalDias }
 */
export const calcularDuracion = (fechaInicio, fechaFin) => {
  if (!fechaInicio) {
    return { anos: 0, meses: 0, dias: 0, totalDias: 0 };
  }

  const inicio = new Date(fechaInicio);
  const fin = fechaFin ? new Date(fechaFin) : new Date();

  // Validar fechas
  if (isNaN(inicio.getTime()) || (fechaFin && isNaN(fin.getTime()))) {
    return { anos: 0, meses: 0, dias: 0, totalDias: 0 };
  }

  // Calcular diferencia en milisegundos
  const diffMs = fin - inicio;
  const totalDias = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  // Calcular años, meses y días
  let anos = fin.getFullYear() - inicio.getFullYear();
  let meses = fin.getMonth() - inicio.getMonth();
  let dias = fin.getDate() - inicio.getDate();

  // Ajustar si los días son negativos
  if (dias < 0) {
    meses--;
    const ultimoDiaMesAnterior = new Date(fin.getFullYear(), fin.getMonth(), 0);
    dias += ultimoDiaMesAnterior.getDate();
  }

  // Ajustar si los meses son negativos
  if (meses < 0) {
    anos--;
    meses += 12;
  }

  return {
    anos: Math.max(0, anos),
    meses: Math.max(0, meses),
    dias: Math.max(0, dias),
    totalDias: Math.max(0, totalDias),
  };
};

/**
 * Calcula el total acumulado de experiencia por tipo
 * @param {Array} experiencias - Array de experiencias profesionales
 * @returns {Object} Objeto con totales por tipo: { general: {anos, meses, dias, totalDias}, direccion: {...} }
 */
export const calcularTotalPorTipo = (experiencias) => {
  const totales = {
    general: { anos: 0, meses: 0, dias: 0, totalDias: 0 },
    direccion: { anos: 0, meses: 0, dias: 0, totalDias: 0 },
  };

  if (!experiencias || experiencias.length === 0) {
    return totales;
  }

  experiencias.forEach((exp) => {
    const duracion = calcularDuracion(exp.fecha_inicio, exp.fecha_fin);
    const tipo = exp.tipo === 'direccion' ? 'direccion' : 'general';

    totales[tipo].anos += duracion.anos;
    totales[tipo].meses += duracion.meses;
    totales[tipo].dias += duracion.dias;
    totales[tipo].totalDias += duracion.totalDias;
  });

  // Normalizar (convertir meses excedentes a años, días excedentes a meses)
  Object.keys(totales).forEach((tipo) => {
    const total = totales[tipo];
    
    // Convertir días excedentes a meses (30 días ≈ 1 mes)
    const mesesExtras = Math.floor(total.dias / 30);
    total.meses += mesesExtras;
    total.dias = total.dias % 30;

    // Convertir meses excedentes a años
    const anosExtras = Math.floor(total.meses / 12);
    total.anos += anosExtras;
    total.meses = total.meses % 12;
  });

  return totales;
};

/**
 * Formatea la duración para mostrar en pantalla
 * @param {Object} duracion - Objeto con { anos, meses, dias }
 * @returns {string} String formateado ej: "2 años, 3 meses, 15 días"
 */
export const formatearDuracion = (duracion) => {
  const partes = [];

  if (duracion.anos > 0) {
    partes.push(`${duracion.anos} año${duracion.anos !== 1 ? 's' : ''}`);
  }

  if (duracion.meses > 0) {
    partes.push(`${duracion.meses} mes${duracion.meses !== 1 ? 'es' : ''}`);
  }

  if (duracion.dias > 0 || partes.length === 0) {
    partes.push(`${duracion.dias} día${duracion.dias !== 1 ? 's' : ''}`);
  }

  return partes.join(', ');
};
