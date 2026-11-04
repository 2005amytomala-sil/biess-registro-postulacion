import { useState, useEffect } from 'react';
import Card from '../Card';
import Input from '../Input';
import Button from '../Button';
import Alert from '../Alert';

const Inhabilidades = ({ datos, onChange, onSave, saving, error }) => {
  const [respuestas, setRespuestas] = useState({});
  const [detalle1_9_7, setDetalle1_9_7] = useState({
    institucion: '',
    cargo: '',
    fecha_finalizacion: '',
    motivo: '',
  });
  const [errors, setErrors] = useState({});

  // Preguntas de inhabilidad (código_pregunta, texto)
  const preguntas = [
    { codigo: '1.9.1', texto: '¿Ha sido condenado por delito contra la administración pública?' },
    { codigo: '1.9.2', texto: '¿Ha sido destituido del servicio público?' },
    { codigo: '1.9.3', texto: '¿Tiene antecedentes penales?' },
    { codigo: '1.9.4', texto: '¿Ha sido sancionado por responsabilidad civil o penal?' },
    { codigo: '1.9.5', texto: '¿Tiene deudas con el IESS?' },
    { codigo: '1.9.6', texto: '¿Ha sido objeto de remoción de cargo?' },
    { codigo: '1.9.7', texto: '¿Ha ocupado cargo público en los últimos 2 años?' },
    { codigo: '1.9.8', texto: '¿Es cónyuge o pariente de autoridad del IESS?' },
    { codigo: '1.9.9', texto: '¿Ha participado en actos de corrupción?' },
    { codigo: '1.9.10', texto: '¿Tiene conflicto de intereses?' },
    { codigo: '1.9.11', texto: '¿Ha sido inhabilitado para ejercer cargos públicos?' },
    { codigo: '1.9.12', texto: '¿Tiene procesos judiciales pendientes?' },
    { codigo: '1.9.13', texto: '¿Otras inhabilidades previstas en la ley?' },
  ];

  // Cargar respuestas existentes al montar
  useEffect(() => {
    if (datos && datos.length > 0) {
      const nuevasRespuestas = {};
      const nuevoDetalle = { institucion: '', cargo: '', fecha_finalizacion: '', motivo: '' };

      datos.forEach((inh) => {
        nuevasRespuestas[inh.codigo_pregunta] = inh.respuesta;
        
        // Si es la pregunta 1.9.7 y tiene detalle, cargarlo
        if (inh.codigo_pregunta === '1.9.7' && inh.detalle) {
          try {
            const detalleParsed = typeof inh.detalle === 'string' 
              ? JSON.parse(inh.detalle) 
              : inh.detalle;
            if (detalleParsed) {
              nuevoDetalle.institucion = detalleParsed.institucion || '';
              nuevoDetalle.cargo = detalleParsed.cargo || '';
              nuevoDetalle.fecha_finalizacion = detalleParsed.fecha_finalizacion || '';
              nuevoDetalle.motivo = detalleParsed.motivo || '';
            }
          } catch (e) {
            console.error('Error al parsear detalle:', e);
          }
        }
      });

      setRespuestas(nuevasRespuestas);
      setDetalle1_9_7(nuevoDetalle);
    }
  }, [datos]);

  const handleRespuestaChange = (codigo, valor) => {
    setRespuestas({ ...respuestas, [codigo]: valor === 'si' });
    
    // Limpiar errores si existe
    if (errors[codigo]) {
      const newErrors = { ...errors };
      delete newErrors[codigo];
      setErrors(newErrors);
    }
  };

  const handleDetalle1_9_7Change = (field, value) => {
    setDetalle1_9_7({ ...detalle1_9_7, [field]: value });
    
    // Validación en vivo
    const newErrors = { ...errors };
    if (field === 'institucion' && value.trim()) {
      delete newErrors.detalle1_9_7_institucion;
    }
    if (field === 'cargo' && value.trim()) {
      delete newErrors.detalle1_9_7_cargo;
    }
    if (field === 'fecha_finalizacion' && value) {
      delete newErrors.detalle1_9_7_fecha_finalizacion;
    }
    if (field === 'motivo' && value.trim()) {
      delete newErrors.detalle1_9_7_motivo;
    }
    setErrors(newErrors);
  };

  const handleGuardar = async () => {
    const newErrors = {};

    // Validar que todas las preguntas tengan respuesta
    preguntas.forEach((pregunta) => {
      if (respuestas[pregunta.codigo] === undefined) {
        newErrors[pregunta.codigo] = 'Debe responder esta pregunta';
      }
    });

    // Si 1.9.7 es Sí, validar campos adicionales
    if (respuestas['1.9.7'] === true) {
      if (!detalle1_9_7.institucion.trim()) {
        newErrors.detalle1_9_7_institucion = 'La institución es obligatoria';
      }
      if (!detalle1_9_7.cargo.trim()) {
        newErrors.detalle1_9_7_cargo = 'El cargo es obligatorio';
      }
      if (!detalle1_9_7.fecha_finalizacion) {
        newErrors.detalle1_9_7_fecha_finalizacion = 'La fecha de finalización es obligatoria';
      }
      if (!detalle1_9_7.motivo.trim()) {
        newErrors.detalle1_9_7_motivo = 'El motivo es obligatorio';
      }
    }

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    // Construir array de inhabilidades
    const inhabilidades = preguntas.map((pregunta) => {
      const inhabilidad = {
        codigo_pregunta: pregunta.codigo,
        respuesta: respuestas[pregunta.codigo] === true,
        detalle: '',
      };

      // Si es 1.9.7 y la respuesta es Sí, agregar detalle como JSON
      if (pregunta.codigo === '1.9.7' && respuestas['1.9.7'] === true) {
        inhabilidad.detalle = JSON.stringify(detalle1_9_7);
      }

      return inhabilidad;
    });

    onChange(inhabilidades);
    // Se pasan los datos directamente: el estado del padre aún no se actualizó en este render
    await onSave(inhabilidades);
  };

  return (
    <Card title="1.9 Declaración de Inhabilidades">
      <div className="space-y-6">
        {error && <Alert type="error" message={error} onClose={() => {}} />}

        <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
          <p className="text-sm text-yellow-800">
            <strong>Importante:</strong> Debe responder todas las preguntas. Si alguna respuesta es afirmativa,
            proporcione la información adicional solicitada.
          </p>
        </div>

        {/* Lista de preguntas */}
        <div className="space-y-4">
          {preguntas.map((pregunta, index) => (
            <div key={pregunta.codigo} className="border-b border-slate-200 pb-4 last:border-0">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <p className="text-sm font-medium text-slate-900 mb-2">
                    {index + 1}. {pregunta.texto}
                  </p>
                  <div className="flex gap-4">
                    <label className="flex items-center">
                      <input
                        type="radio"
                        name={`pregunta-${pregunta.codigo}`}
                        value="si"
                        checked={respuestas[pregunta.codigo] === true}
                        onChange={(e) => handleRespuestaChange(pregunta.codigo, e.target.value)}
                        className="mr-2"
                      />
                      <span className="text-sm">Sí</span>
                    </label>
                    <label className="flex items-center">
                      <input
                        type="radio"
                        name={`pregunta-${pregunta.codigo}`}
                        value="no"
                        checked={respuestas[pregunta.codigo] === false}
                        onChange={(e) => handleRespuestaChange(pregunta.codigo, e.target.value)}
                        className="mr-2"
                      />
                      <span className="text-sm">No</span>
                    </label>
                  </div>
                  {errors[pregunta.codigo] && (
                    <p className="mt-1 text-sm text-red-600">{errors[pregunta.codigo]}</p>
                  )}
                </div>
              </div>

              {/* Campos adicionales para 1.9.7 si la respuesta es Sí */}
              {pregunta.codigo === '1.9.7' && respuestas['1.9.7'] === true && (
                <div className="mt-4 ml-4 bg-slate-50 border rounded-lg p-4 space-y-4">
                  <h4 className="text-sm font-medium text-slate-900">Información adicional requerida:</h4>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <Input
                      label="Institución *"
                      value={detalle1_9_7.institucion}
                      onChange={(e) => handleDetalle1_9_7Change('institucion', e.target.value)}
                      error={errors.detalle1_9_7_institucion}
                      required
                    />
                    <Input
                      label="Cargo *"
                      value={detalle1_9_7.cargo}
                      onChange={(e) => handleDetalle1_9_7Change('cargo', e.target.value)}
                      error={errors.detalle1_9_7_cargo}
                      required
                    />
                    <Input
                      label="Fecha de Finalización *"
                      type="date"
                      value={detalle1_9_7.fecha_finalizacion}
                      onChange={(e) => handleDetalle1_9_7Change('fecha_finalizacion', e.target.value)}
                      error={errors.detalle1_9_7_fecha_finalizacion}
                      required
                    />
                    <div className="md:col-span-2">
                      <label className="block text-sm font-medium text-slate-700 mb-1">
                        Motivo *
                      </label>
                      <textarea
                        value={detalle1_9_7.motivo}
                        onChange={(e) => handleDetalle1_9_7Change('motivo', e.target.value)}
                        rows={3}
                        className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                          errors.detalle1_9_7_motivo ? 'border-red-500' : 'border-slate-300'
                        }`}
                        placeholder="Describa el motivo..."
                      />
                      {errors.detalle1_9_7_motivo && (
                        <p className="mt-1 text-sm text-red-600">{errors.detalle1_9_7_motivo}</p>
                      )}
                    </div>
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>

        <div className="flex justify-end border-t border-slate-100 pt-5">
          <Button
            onClick={handleGuardar}
            disabled={saving}
          >
            {saving ? 'Guardando...' : 'Guardar'}
          </Button>
        </div>
      </div>
    </Card>
  );
};

export default Inhabilidades;
