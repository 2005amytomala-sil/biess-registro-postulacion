import { useState } from 'react';
import Card from '../Card';
import Input from '../Input';
import Select from '../Select';
import Button from '../Button';
import Alert from '../Alert';

const FormacionAcademica = ({ datos, onChange, onSave, saving, error }) => {
  const [editandoIndex, setEditandoIndex] = useState(null);
  const [formacionForm, setFormacionForm] = useState({
    nivel: '',
    institucion: '',
    tipo: '',
    titulo: '',
    area_estudios: '',
    registro_senescyt: '',
    fecha_registro_senescyt: '',
  });
  const [errors, setErrors] = useState({});

  const niveles = [
    { value: 'tercer_nivel', label: 'Tercer Nivel de Grado' },
    { value: 'cuarto_nivel_especialidad', label: 'Cuarto Nivel Especialidad' },
    { value: 'maestria', label: 'Maestría' },
    { value: 'phd', label: 'PhD' },
  ];

  const tipos = [
    { value: 'nacional', label: 'Nacional' },
    { value: 'extranjera', label: 'Extranjera' },
  ];

  const handleAgregarFormacion = () => {
    setEditandoIndex('nuevo');
    setFormacionForm({
      nivel: '',
      institucion: '',
      tipo: '',
      titulo: '',
      area_estudios: '',
      registro_senescyt: '',
      fecha_registro_senescyt: '',
    });
    setErrors({});
  };

  const handleEditarFormacion = (index) => {
    setEditandoIndex(index);
    const formacion = datos[index];
    setFormacionForm({
      nivel: formacion.nivel || '',
      institucion: formacion.institucion || '',
      tipo: formacion.tipo || '',
      titulo: formacion.titulo || '',
      area_estudios: formacion.area_estudios || '',
      registro_senescyt: formacion.registro_senescyt || '',
      fecha_registro_senescyt: formacion.fecha_registro_senescyt || '',
    });
    setErrors({});
  };

  const handleEliminarFormacion = (index) => {
    const nuevasFormaciones = datos.filter((_, i) => i !== index);
    onChange(nuevasFormaciones);
  };

  const handleGuardarFormacion = () => {
    const newErrors = {};

    if (!formacionForm.nivel) {
      newErrors.nivel = 'El nivel es obligatorio';
    }
    if (!formacionForm.institucion.trim()) {
      newErrors.institucion = 'La institución es obligatoria';
    }
    if (!formacionForm.tipo) {
      newErrors.tipo = 'El tipo es obligatorio';
    }
    if (!formacionForm.titulo.trim()) {
      newErrors.titulo = 'El título es obligatorio';
    }

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    let nuevasFormaciones;
    if (editandoIndex === 'nuevo') {
      nuevasFormaciones = [...datos, { ...formacionForm }];
    } else {
      nuevasFormaciones = datos.map((f, i) =>
        i === editandoIndex ? { ...formacionForm } : f
      );
    }

    onChange(nuevasFormaciones);
    setEditandoIndex(null);
    setFormacionForm({
      nivel: '',
      institucion: '',
      tipo: '',
      titulo: '',
      area_estudios: '',
      registro_senescyt: '',
      fecha_registro_senescyt: '',
    });
    setErrors({});
  };

  const handleCancelarEdicion = () => {
    setEditandoIndex(null);
    setFormacionForm({
      nivel: '',
      institucion: '',
      tipo: '',
      titulo: '',
      area_estudios: '',
      registro_senescyt: '',
      fecha_registro_senescyt: '',
    });
    setErrors({});
  };

  const handleFormacionFormChange = (field, value) => {
    setFormacionForm({ ...formacionForm, [field]: value });
    
    // Validación en vivo
    const newErrors = { ...errors };
    if (field === 'institucion' && value.trim()) {
      delete newErrors.institucion;
    }
    if (field === 'titulo' && value.trim()) {
      delete newErrors.titulo;
    }
    if (field === 'nivel' && value) {
      delete newErrors.nivel;
    }
    if (field === 'tipo' && value) {
      delete newErrors.tipo;
    }
    setErrors(newErrors);
  };

  return (
    <Card title="1.5 Formación Académica">
      <div className="space-y-6">
        {error && <Alert type="error" message={error} onClose={() => {}} />}

        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-sm font-semibold uppercase tracking-wide text-slate-500">Lista de Formación Académica</h3>
            <Button
              type="button"
              variant="outline"
              onClick={handleAgregarFormacion}
              disabled={editandoIndex !== null}
            >
              + Agregar Formación
            </Button>
          </div>

          {/* Formulario de edición/creación */}
          {editandoIndex !== null && (
            <div className="bg-slate-50 border rounded-lg p-4 space-y-4">
              <h4 className="text-sm font-medium text-slate-900">
                {editandoIndex === 'nuevo' ? 'Nueva Formación' : 'Editar Formación'}
              </h4>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <Select
                  label="Nivel *"
                  options={niveles}
                  value={formacionForm.nivel}
                  onChange={(e) => handleFormacionFormChange('nivel', e.target.value)}
                  error={errors.nivel}
                  required
                />

                <Input
                  label="Institución *"
                  value={formacionForm.institucion}
                  onChange={(e) => handleFormacionFormChange('institucion', e.target.value)}
                  error={errors.institucion}
                  required
                />

                <Select
                  label="Tipo *"
                  options={tipos}
                  value={formacionForm.tipo}
                  onChange={(e) => handleFormacionFormChange('tipo', e.target.value)}
                  error={errors.tipo}
                  required
                />

                <Input
                  label="Título *"
                  value={formacionForm.titulo}
                  onChange={(e) => handleFormacionFormChange('titulo', e.target.value)}
                  error={errors.titulo}
                  required
                />

                <Input
                  label="Área de Estudios"
                  value={formacionForm.area_estudios}
                  onChange={(e) => handleFormacionFormChange('area_estudios', e.target.value)}
                  className="md:col-span-2"
                />

                <Input
                  label="Número de Registro SENESCYT"
                  value={formacionForm.registro_senescyt}
                  onChange={(e) => handleFormacionFormChange('registro_senescyt', e.target.value)}
                  placeholder="Ej: 12345678"
                />

                <Input
                  label="Fecha de Registro SENESCYT"
                  type="date"
                  value={formacionForm.fecha_registro_senescyt}
                  onChange={(e) => handleFormacionFormChange('fecha_registro_senescyt', e.target.value)}
                />
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
                  onClick={handleGuardarFormacion}
                >
                  {editandoIndex === 'nuevo' ? 'Agregar' : 'Actualizar'}
                </Button>
              </div>
            </div>
          )}

          {/* Lista de formaciones existentes */}
          {datos && datos.length > 0 ? (
            <div className="space-y-3">
              {datos.map((formacion, index) => (
                <div
                  key={index}
                  className="bg-white border rounded-lg p-4 flex justify-between items-start"
                >
                  <div className="flex-1">
                    <p className="font-medium text-slate-900">{formacion.titulo}</p>
                    <p className="text-sm text-slate-600">
                      {niveles.find(n => n.value === formacion.nivel)?.label} - {formacion.institucion}
                    </p>
                    <p className="text-sm text-slate-500">
                      Tipo: {tipos.find(t => t.value === formacion.tipo)?.label}
                    </p>
                    {formacion.area_estudios && (
                      <p className="text-sm text-slate-500">Área: {formacion.area_estudios}</p>
                    )}
                    {formacion.registro_senescyt && (
                      <p className="text-sm text-slate-500">Registro SENESCYT: {formacion.registro_senescyt}</p>
                    )}
                    {formacion.fecha_registro_senescyt && (
                      <p className="text-sm text-slate-500">Fecha: {formacion.fecha_registro_senescyt}</p>
                    )}
                  </div>
                  <div className="flex gap-2 ml-4">
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={() => handleEditarFormacion(index)}
                      disabled={editandoIndex !== null}
                    >
                      Editar
                    </Button>
                    <Button
                      type="button"
                      variant="danger"
                      size="sm"
                      onClick={() => handleEliminarFormacion(index)}
                      disabled={editandoIndex !== null}
                    >
                      Eliminar
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            !editandoIndex && (
              <div className="text-center py-8 text-slate-500 border-2 border-dashed rounded-lg">
                No hay formación académica agregada. Click en "Agregar Formación" para comenzar.
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

export default FormacionAcademica;
