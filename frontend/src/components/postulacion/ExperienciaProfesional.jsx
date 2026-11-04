import { useState } from 'react';
import Card from '../Card';
import Input from '../Input';
import Select from '../Select';
import Button from '../Button';
import Alert from '../Alert';
import { calcularDuracion, formatearDuracion, calcularTotalPorTipo } from '../../utils/duracion';

const ExperienciaProfesional = ({ datos, onChange, onSave, saving, error }) => {
  const [editandoIndex, setEditandoIndex] = useState(null);
  const [experienciaForm, setExperienciaForm] = useState({
    tipo: '',
    cargo: '',
    institucion: '',
    fecha_inicio: '',
    fecha_fin: '',
    actividades_relacionadas: '',
    descripcion: '',
  });
  const [errors, setErrors] = useState({});

  const tipos = [
    { value: 'general', label: 'General' },
    { value: 'direccion', label: 'Dirección' },
  ];

  const handleAgregarExperiencia = () => {
    setEditandoIndex('nuevo');
    setExperienciaForm({
      tipo: '',
      cargo: '',
      institucion: '',
      fecha_inicio: '',
      fecha_fin: '',
      actividades_relacionadas: '',
      descripcion: '',
    });
    setErrors({});
  };

  const handleEditarExperiencia = (index) => {
    setEditandoIndex(index);
    const experiencia = datos[index];
    setExperienciaForm({
      tipo: experiencia.tipo || '',
      cargo: experiencia.cargo || '',
      institucion: experiencia.institucion || '',
      fecha_inicio: experiencia.fecha_inicio || '',
      fecha_fin: experiencia.fecha_fin || '',
      actividades_relacionadas: experiencia.actividades_relacionadas || '',
      descripcion: experiencia.descripcion || '',
    });
    setErrors({});
  };

  const handleEliminarExperiencia = (index) => {
    const nuevasExperiencias = datos.filter((_, i) => i !== index);
    onChange(nuevasExperiencias);
  };

  const handleGuardarExperiencia = () => {
    const newErrors = {};

    if (!experienciaForm.tipo) {
      newErrors.tipo = 'El tipo es obligatorio';
    }
    if (!experienciaForm.cargo.trim()) {
      newErrors.cargo = 'El cargo es obligatorio';
    }
    if (!experienciaForm.institucion.trim()) {
      newErrors.institucion = 'La institución es obligatoria';
    }
    if (!experienciaForm.fecha_inicio) {
      newErrors.fecha_inicio = 'La fecha de inicio es obligatoria';
    }

    // Validar fecha fin >= fecha inicio
    if (experienciaForm.fecha_inicio && experienciaForm.fecha_fin) {
      const inicio = new Date(experienciaForm.fecha_inicio);
      const fin = new Date(experienciaForm.fecha_fin);
      if (fin < inicio) {
        newErrors.fecha_fin = 'La fecha de fin debe ser posterior o igual a la fecha de inicio';
      }
    }

    // Validar descripción <= 1000 caracteres
    if (experienciaForm.descripcion.length > 1000) {
      newErrors.descripcion = 'La descripción no puede exceder 1000 caracteres';
    }

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    let nuevasExperiencias;
    if (editandoIndex === 'nuevo') {
      nuevasExperiencias = [...datos, { ...experienciaForm }];
    } else {
      nuevasExperiencias = datos.map((e, i) =>
        i === editandoIndex ? { ...experienciaForm } : e
      );
    }

    onChange(nuevasExperiencias);
    setEditandoIndex(null);
    setExperienciaForm({
      tipo: '',
      cargo: '',
      institucion: '',
      fecha_inicio: '',
      fecha_fin: '',
      actividades_relacionadas: '',
      descripcion: '',
    });
    setErrors({});
  };

  const handleCancelarEdicion = () => {
    setEditandoIndex(null);
    setExperienciaForm({
      tipo: '',
      cargo: '',
      institucion: '',
      fecha_inicio: '',
      fecha_fin: '',
      actividades_relacionadas: '',
      descripcion: '',
    });
    setErrors({});
  };

  const handleExperienciaFormChange = (field, value) => {
    setExperienciaForm({ ...experienciaForm, [field]: value });
    
    // Validación en vivo
    const newErrors = { ...errors };
    if (field === 'cargo' && value.trim()) {
      delete newErrors.cargo;
    }
    if (field === 'institucion' && value.trim()) {
      delete newErrors.institucion;
    }
    if (field === 'tipo' && value) {
      delete newErrors.tipo;
    }
    if (field === 'fecha_inicio' && value) {
      delete newErrors.fecha_inicio;
    }
    if (field === 'fecha_fin' && experienciaForm.fecha_inicio && value) {
      const inicio = new Date(experienciaForm.fecha_inicio);
      const fin = new Date(value);
      if (fin >= inicio) {
        delete newErrors.fecha_fin;
      }
    }
    if (field === 'descripcion' && value.length <= 1000) {
      delete newErrors.descripcion;
    }
    setErrors(newErrors);
  };

  // Calcular totales por tipo
  const totales = calcularTotalPorTipo(datos);

  return (
    <Card title="1.6 Experiencia Profesional">
      <div className="space-y-6">
        {error && <Alert type="error" message={error} onClose={() => {}} />}

        {/* Resumen de totales */}
        {datos && datos.length > 0 && (
          <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <h3 className="text-sm font-medium text-slate-900 mb-3">Total de Experiencia Acumulada</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <p className="text-sm font-medium text-slate-700">General:</p>
                <p className="text-sm text-slate-600">{formatearDuracion(totales.general)}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-slate-700">Dirección:</p>
                <p className="text-sm text-slate-600">{formatearDuracion(totales.direccion)}</p>
              </div>
            </div>
          </div>
        )}

        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-sm font-semibold uppercase tracking-wide text-slate-500">Lista de Experiencia Profesional</h3>
            <Button
              type="button"
              variant="outline"
              onClick={handleAgregarExperiencia}
              disabled={editandoIndex !== null}
            >
              + Agregar Experiencia
            </Button>
          </div>

          {/* Formulario de edición/creación */}
          {editandoIndex !== null && (
            <div className="bg-slate-50 border rounded-lg p-4 space-y-4">
              <h4 className="text-sm font-medium text-slate-900">
                {editandoIndex === 'nuevo' ? 'Nueva Experiencia' : 'Editar Experiencia'}
              </h4>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <Select
                  label="Tipo *"
                  options={tipos}
                  value={experienciaForm.tipo}
                  onChange={(e) => handleExperienciaFormChange('tipo', e.target.value)}
                  error={errors.tipo}
                  required
                />

                <Input
                  label="Cargo *"
                  value={experienciaForm.cargo}
                  onChange={(e) => handleExperienciaFormChange('cargo', e.target.value)}
                  error={errors.cargo}
                  required
                />

                <Input
                  label="Institución *"
                  value={experienciaForm.institucion}
                  onChange={(e) => handleExperienciaFormChange('institucion', e.target.value)}
                  error={errors.institucion}
                  required
                />

                <div></div>

                <Input
                  label="Fecha de Inicio *"
                  type="date"
                  value={experienciaForm.fecha_inicio}
                  onChange={(e) => handleExperienciaFormChange('fecha_inicio', e.target.value)}
                  error={errors.fecha_inicio}
                  required
                />

                <Input
                  label="Fecha de Fin"
                  type="date"
                  value={experienciaForm.fecha_fin}
                  onChange={(e) => handleExperienciaFormChange('fecha_fin', e.target.value)}
                  error={errors.fecha_fin}
                  helperText="Dejar vacío si está en curso"
                />

                <Input
                  label="Actividades Relacionadas"
                  value={experienciaForm.actividades_relacionadas}
                  onChange={(e) => handleExperienciaFormChange('actividades_relacionadas', e.target.value)}
                  className="md:col-span-2"
                />

                <div className="md:col-span-2">
                  <label className="block text-sm font-medium text-slate-700 mb-1">
                    Descripción * (máx. 1000 caracteres)
                  </label>
                  <textarea
                    value={experienciaForm.descripcion}
                    onChange={(e) => handleExperienciaFormChange('descripcion', e.target.value)}
                    rows={4}
                    className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                      errors.descripcion ? 'border-red-500' : 'border-slate-300'
                    }`}
                    placeholder="Describa sus funciones y logros..."
                  />
                  {errors.descripcion && (
                    <p className="mt-1 text-sm text-red-600">{errors.descripcion}</p>
                  )}
                  <p className="mt-1 text-sm text-slate-500">
                    Caracteres: {experienciaForm.descripcion.length}/1000
                  </p>
                </div>
              </div>

              <div className="flex gap-3 justify-end">
                <Button
                  type="button"
                  variant="secondary"
                  onClick={handleCancelarEdicion}
                >
                  Cancelar
                </Button>
                <Button
                  type="button"
                  onClick={handleGuardarExperiencia}
                >
                  {editandoIndex === 'nuevo' ? 'Agregar' : 'Actualizar'}
                </Button>
              </div>
            </div>
          )}

          {/* Lista de experiencias existentes */}
          {datos && datos.length > 0 ? (
            <div className="space-y-3">
              {datos.map((experiencia, index) => {
                const duracion = calcularDuracion(experiencia.fecha_inicio, experiencia.fecha_fin);
                return (
                  <div
                    key={index}
                    className="bg-white border rounded-lg p-4 flex justify-between items-start"
                  >
                    <div className="flex-1">
                      <p className="font-medium text-slate-900">{experiencia.cargo}</p>
                      <p className="text-sm text-slate-600">
                        {tipos.find(t => t.value === experiencia.tipo)?.label} - {experiencia.institucion}
                      </p>
                      <p className="text-sm text-slate-500">
                        {experiencia.fecha_inicio} {experiencia.fecha_fin ? ` - ${experiencia.fecha_fin}` : ' - En curso'}
                      </p>
                      <p className="text-sm text-slate-500">
                        Duración: {formatearDuracion(duracion)}
                      </p>
                      {experiencia.actividades_relacionadas && (
                        <p className="text-sm text-slate-500">
                          Actividades: {experiencia.actividades_relacionadas}
                        </p>
                      )}
                      {experiencia.descripcion && (
                        <p className="text-sm text-slate-500 mt-1">
                          {experiencia.descripcion.substring(0, 100)}
                          {experiencia.descripcion.length > 100 && '...'}
                        </p>
                      )}
                    </div>
                    <div className="flex gap-2 ml-4">
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        onClick={() => handleEditarExperiencia(index)}
                        disabled={editandoIndex !== null}
                      >
                        Editar
                      </Button>
                      <Button
                        type="button"
                        variant="danger"
                        size="sm"
                        onClick={() => handleEliminarExperiencia(index)}
                        disabled={editandoIndex !== null}
                      >
                        Eliminar
                      </Button>
                    </div>
                  </div>
                );
              })}
            </div>
          ) : (
            !editandoIndex && (
              <div className="text-center py-8 text-slate-500 border-2 border-dashed rounded-lg">
                No hay experiencia profesional agregada. Click en "Agregar Experiencia" para comenzar.
              </div>
            )
          )}
        </div>

        <div className="flex justify-end border-t border-slate-100 pt-5">
          <Button
            onClick={onSave}
            disabled={saving}
          >
            {saving ? 'Guardando...' : 'Guardar'}
          </Button>
        </div>
      </div>
    </Card>
  );
};

export default ExperienciaProfesional;
