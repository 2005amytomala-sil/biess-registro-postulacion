import { useState } from 'react';
import Card from './Card';
import Input from './Input';
import Select from './Select';
import Button from './Button';
import Alert from './Alert';

const FamiliaresIESS = ({ datos, onChange, onSave, saving, error }) => {
  const [editandoIndex, setEditandoIndex] = useState(null);
  const [familiarForm, setFamiliarForm] = useState({
    nombres: '',
    parentesco: '',
    institucion: '',
    area: '',
    cargo: '',
  });

  const instituciones = [
    { value: 'IESS', label: 'IESS' },
    { value: 'BIESS', label: 'BIESS' },
  ];

  const parentescos = [
    { value: 'padre', label: 'Padre' },
    { value: 'madre', label: 'Madre' },
    { value: 'hermano', label: 'Hermano/a' },
    { value: 'hijo', label: 'Hijo/a' },
    { value: 'esposo', label: 'Esposo/a' },
    { value: 'abuelo', label: 'Abuelo/a' },
    { value: 'tio', label: 'Tío/a' },
    { value: 'sobrino', label: 'Sobrino/a' },
    { value: 'primo', label: 'Primo/a' },
    { value: 'cunado', label: 'Cuñado/a' },
    { value: 'nuera', label: 'Nuera' },
    { value: 'yerno', label: 'Yerno' },
    { value: 'otro', label: 'Otro' },
  ];

  const handleTieneFamiliaresChange = (value) => {
    const tieneFamiliares = value === 'si';
    const nuevosDatos = {
      ...datos,
      tiene_familiares: tieneFamiliares,
      familiares: tieneFamiliares ? datos.familiares : [],
    };
    onChange(nuevosDatos);
  };

  const handleAgregarFamiliar = () => {
    setEditandoIndex('nuevo');
    setFamiliarForm({
      nombres: '',
      parentesco: '',
      institucion: '',
      area: '',
      cargo: '',
    });
  };

  const handleEditarFamiliar = (index) => {
    setEditandoIndex(index);
    setFamiliarForm({ ...datos.familiares[index] });
  };

  const handleEliminarFamiliar = (index) => {
    const nuevosFamiliares = datos.familiares.filter((_, i) => i !== index);
    onChange({ ...datos, familiares: nuevosFamiliares });
  };

  const handleGuardarFamiliar = () => {
    // Validar
    if (!familiarForm.nombres.trim()) {
      alert('El nombre es obligatorio');
      return;
    }
    if (!familiarForm.parentesco) {
      alert('El parentesco es obligatorio');
      return;
    }
    if (!familiarForm.institucion) {
      alert('La institución es obligatoria');
      return;
    }

    let nuevosFamiliares;
    if (editandoIndex === 'nuevo') {
      nuevosFamiliares = [...datos.familiares, { ...familiarForm }];
    } else {
      nuevosFamiliares = datos.familiares.map((f, i) =>
        i === editandoIndex ? { ...familiarForm } : f
      );
    }

    onChange({ ...datos, familiares: nuevosFamiliares });
    setEditandoIndex(null);
    setFamiliarForm({
      nombres: '',
      parentesco: '',
      institucion: '',
      area: '',
      cargo: '',
    });
  };

  const handleCancelarEdicion = () => {
    setEditandoIndex(null);
    setFamiliarForm({
      nombres: '',
      parentesco: '',
      institucion: '',
      area: '',
      cargo: '',
    });
  };

  const handleFamiliarFormChange = (field, value) => {
    setFamiliarForm({ ...familiarForm, [field]: value });
  };

  return (
    <Card title="1.3 Información Adicional">
      <div className="space-y-6">
        {error && <Alert type="error" message={error} onClose={() => {}} />}

        {/* Pregunta sobre familiares */}
        <div className="space-y-3">
          <label className="block text-sm font-medium text-gray-700">
            ¿Tiene familiares dentro del 2.º grado de consanguinidad y 4.º de afinidad en IESS/BIESS?
          </label>
          <div className="flex gap-4">
            <label className="flex items-center">
              <input
                type="radio"
                name="tiene_familiares"
                value="si"
                checked={datos.tiene_familiares === true}
                onChange={(e) => handleTieneFamiliaresChange(e.target.value)}
                className="mr-2"
              />
              <span>Sí</span>
            </label>
            <label className="flex items-center">
              <input
                type="radio"
                name="tiene_familiares"
                value="no"
                checked={datos.tiene_familiares === false}
                onChange={(e) => handleTieneFamiliaresChange(e.target.value)}
                className="mr-2"
              />
              <span>No</span>
            </label>
          </div>
        </div>

        {/* Lista de familiares */}
        {datos.tiene_familiares && (
          <div className="space-y-4">
            <div className="flex justify-between items-center">
              <h3 className="text-sm font-medium text-gray-900">Lista de Familiares</h3>
              <Button
                type="button"
                variant="outline"
                onClick={handleAgregarFamiliar}
                disabled={editandoIndex !== null}
              >
                + Agregar Familiar
              </Button>
            </div>

            {/* Formulario de edición/creación */}
            {editandoIndex !== null && (
              <div className="bg-gray-50 border rounded-lg p-4 space-y-4">
                <h4 className="text-sm font-medium text-gray-900">
                  {editandoIndex === 'nuevo' ? 'Nuevo Familiar' : 'Editar Familiar'}
                </h4>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <Input
                    label="Nombres *"
                    value={familiarForm.nombres}
                    onChange={(e) => handleFamiliarFormChange('nombres', e.target.value)}
                    required
                  />

                  <Select
                    label="Parentesco *"
                    options={parentescos}
                    value={familiarForm.parentesco}
                    onChange={(e) => handleFamiliarFormChange('parentesco', e.target.value)}
                    required
                  />

                  <Select
                    label="Institución *"
                    options={instituciones}
                    value={familiarForm.institucion}
                    onChange={(e) => handleFamiliarFormChange('institucion', e.target.value)}
                    required
                  />

                  <Input
                    label="Área"
                    value={familiarForm.area}
                    onChange={(e) => handleFamiliarFormChange('area', e.target.value)}
                  />

                  <Input
                    label="Cargo"
                    value={familiarForm.cargo}
                    onChange={(e) => handleFamiliarFormChange('cargo', e.target.value)}
                    className="md:col-span-2"
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
                    onClick={handleGuardarFamiliar}
                  >
                    {editandoIndex === 'nuevo' ? 'Agregar' : 'Actualizar'}
                  </Button>
                </div>
              </div>
            )}

            {/* Lista de familiares existentes */}
            {datos.familiares && datos.familiares.length > 0 ? (
              <div className="space-y-3">
                {datos.familiares.map((familiar, index) => (
                  <div
                    key={index}
                    className="bg-white border rounded-lg p-4 flex justify-between items-start"
                  >
                    <div className="flex-1">
                      <p className="font-medium text-gray-900">{familiar.nombres}</p>
                      <p className="text-sm text-gray-600">
                        {familiar.parentesco} - {familiar.institucion}
                      </p>
                      {familiar.area && (
                        <p className="text-sm text-gray-500">Área: {familiar.area}</p>
                      )}
                      {familiar.cargo && (
                        <p className="text-sm text-gray-500">Cargo: {familiar.cargo}</p>
                      )}
                    </div>
                    <div className="flex gap-2 ml-4">
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        onClick={() => handleEditarFamiliar(index)}
                        disabled={editandoIndex !== null}
                      >
                        Editar
                      </Button>
                      <Button
                        type="button"
                        variant="danger"
                        size="sm"
                        onClick={() => handleEliminarFamiliar(index)}
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
                <div className="text-center py-8 text-gray-500 border-2 border-dashed rounded-lg">
                  No hay familiares agregados. Click en "Agregar Familiar" para comenzar.
                </div>
              )
            )}

            {/* Validación: si tiene_familiares pero lista vacía */}
            {datos.tiene_familiares && datos.familiares.length === 0 && !editandoIndex && (
              <Alert
                type="error"
                message="Debe agregar al menos un familiar si indica que tiene familiares"
              />
            )}
          </div>
        )}

        <div className="flex justify-end">
          <button
            onClick={onSave}
            disabled={saving || (datos.tiene_familiares && datos.familiares.length === 0)}
            className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {saving ? 'Guardando...' : 'Guardar'}
          </button>
        </div>
      </div>
    </Card>
  );
};

export default FamiliaresIESS;
