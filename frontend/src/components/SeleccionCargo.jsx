import Card from './Card';
import Button from './Button';
import Alert from './Alert';

const SeleccionCargo = ({ datos, onChange, onSave, saving, error }) => {

  const sectores = [
    {
      value: 'afiliados',
      label: 'Afiliados',
      description: 'Representación del Sector de Afiliados ante el Directorio del BIESS',
    },
    {
      value: 'jubilados',
      label: 'Jubilados',
      description: 'Representación del Sector de Jubilados ante el Directorio del BIESS',
    },
    {
      value: 'empleadores',
      label: 'Empleadores',
      description: 'Representación del Sector de Empleadores ante el Directorio del BIESS',
    },
  ];

  const handleSectorChange = (sector) => {
    // Actualizar el sector en el estado local
    const nuevosDatos = {
      ...datos,
      sector: sector,
      respaldo_organizacion: {
        tiene_respaldo: false,
        nombre_organizacion: '',
      },
    };
    onChange(nuevosDatos);
  };

  const handleTieneRespaldoChange = (value) => {
    const tieneRespaldo = value === 'si';
    const nuevosDatos = {
      ...datos,
      respaldo_organizacion: {
        tiene_respaldo: tieneRespaldo,
        nombre_organizacion: tieneRespaldo ? datos.respaldo_organizacion?.nombre_organizacion || '' : '',
      },
    };
    onChange(nuevosDatos);
  };

  const handleNombreOrganizacionChange = (value) => {
    const nuevosDatos = {
      ...datos,
      respaldo_organizacion: {
        ...datos.respaldo_organizacion,
        nombre_organizacion: value,
      },
    };
    onChange(nuevosDatos);
  };

  return (
    <Card title="1.4 Selección del Cargo">
      <div className="space-y-6">
        {error && <Alert type="error" message={error} onClose={() => {}} />}

        {/* Selección de sector */}
        <div className="space-y-4">
          <label className="block text-sm font-medium text-gray-700">
            Seleccione el sector al que desea postularse:
          </label>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {sectores.map((sector) => (
              <div
                key={sector.value}
                className={`border-2 rounded-lg p-4 cursor-pointer transition-all ${
                  datos.sector === sector.value
                    ? 'border-blue-500 bg-blue-50'
                    : 'border-gray-200 hover:border-gray-300'
                }`}
                onClick={() => !datos.sector && handleSectorChange(sector.value)}
              >
                <div className="flex items-start gap-3">
                  <input
                    type="radio"
                    name="sector"
                    value={sector.value}
                    checked={datos.sector === sector.value}
                    onChange={() => handleSectorChange(sector.value)}
                    disabled={!!datos.sector}
                    className="mt-1"
                  />
                  <div className="flex-1">
                    <h3 className="font-medium text-gray-900">{sector.label}</h3>
                    <p className="text-sm text-gray-600 mt-1">{sector.description}</p>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Pregunta de respaldo por organización */}
        {datos.sector && (
          <div className="space-y-4 pt-4 border-t">
            <div className="space-y-3">
              <label className="block text-sm font-medium text-gray-700">
                ¿Cuenta con el respaldo de alguna organización?
              </label>
              <div className="flex gap-4">
                <label className="flex items-center">
                  <input
                    type="radio"
                    name="tiene_respaldo"
                    value="si"
                    checked={datos.respaldo_organizacion?.tiene_respaldo === true}
                    onChange={(e) => handleTieneRespaldoChange(e.target.value)}
                    className="mr-2"
                  />
                  <span>Sí</span>
                </label>
                <label className="flex items-center">
                  <input
                    type="radio"
                    name="tiene_respaldo"
                    value="no"
                    checked={datos.respaldo_organizacion?.tiene_respaldo === false}
                    onChange={(e) => handleTieneRespaldoChange(e.target.value)}
                    className="mr-2"
                  />
                  <span>No</span>
                </label>
              </div>
            </div>

            {/* Nombre de organización */}
            {datos.respaldo_organizacion?.tiene_respaldo && (
              <div className="space-y-2">
                <label className="block text-sm font-medium text-gray-700">
                  Nombre de la organización *
                </label>
                <input
                  type="text"
                  value={datos.respaldo_organizacion?.nombre_organizacion || ''}
                  onChange={(e) => handleNombreOrganizacionChange(e.target.value)}
                  placeholder="Ingrese el nombre de la organización"
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                />
                {!datos.respaldo_organizacion?.nombre_organizacion?.trim() && (
                  <p className="text-sm text-red-600">El nombre de la organización es obligatorio</p>
                )}
              </div>
            )}
          </div>
        )}

        {/* Botón de guardar */}
        {datos.sector && (
          <div className="flex justify-end pt-4 border-t">
            <button
              onClick={onSave}
              disabled={
                saving ||
                (datos.respaldo_organizacion?.tiene_respaldo &&
                  !datos.respaldo_organizacion?.nombre_organizacion?.trim())
              }
              className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {saving ? 'Guardando...' : 'Guardar'}
            </button>
          </div>
        )}
      </div>
    </Card>
  );
};

export default SeleccionCargo;
