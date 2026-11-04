import Card from './Card';
import Alert from './Alert';
import Spinner from './Spinner';

const DatosPersonales = ({ datos, loading, error, onClose }) => {
  if (loading) {
    return (
      <Card title="1.1 Información Personal">
        <div className="flex items-center justify-center py-8">
          <Spinner />
        </div>
      </Card>
    );
  }

  if (error) {
    return (
      <Card title="1.1 Información Personal">
        <Alert type="error" message={error} />
      </Card>
    );
  }

  if (!datos) {
    return (
      <Card title="1.1 Información Personal">
        <Alert type="info" message="No hay datos personales disponibles" />
      </Card>
    );
  }

  // Si la nacionalidad no es ecuatoriana, mostrar pantalla de postulación cerrada
  if (datos.nacionalidad && datos.nacionalidad.toLowerCase() !== 'ecuatoriana') {
    return (
      <Card title="1.1 Información Personal">
        <div className="space-y-4">
          <Alert
            type="error"
            message="Su postulación ha sido cerrada. Solo se permite nacionalidad ecuatoriana."
          />
          <div className="bg-red-50 border border-red-200 rounded-lg p-6">
            <h3 className="text-lg font-semibold text-red-900 mb-2">Postulación Cerrada</h3>
            <p className="text-red-700">
              De acuerdo con las normas del concurso público de méritos del BIESS, solo
              pueden postular personas de nacionalidad ecuatoriana.
            </p>
            <p className="text-red-700 mt-2">
              Su nacionalidad registrada es: <strong>{datos.nacionalidad}</strong>
            </p>
          </div>
          {onClose && (
            <button
              onClick={onClose}
              className="mt-4 px-4 py-2 bg-gray-200 text-gray-900 rounded-md hover:bg-gray-300"
            >
              Cerrar
            </button>
          )}
        </div>
      </Card>
    );
  }

  return (
    <Card title="1.1 Información Personal">
      <div className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Nombres</label>
            <p className="px-3 py-2 bg-gray-100 rounded-md text-gray-900">{datos.nombres}</p>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Apellidos</label>
            <p className="px-3 py-2 bg-gray-100 rounded-md text-gray-900">{datos.apellidos}</p>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Cédula</label>
            <p className="px-3 py-2 bg-gray-100 rounded-md text-gray-900">{datos.cedula}</p>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Género</label>
            <p className="px-3 py-2 bg-gray-100 rounded-md text-gray-900 capitalize">{datos.genero}</p>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Fecha de Nacimiento</label>
            <p className="px-3 py-2 bg-gray-100 rounded-md text-gray-900">{datos.fecha_nacimiento}</p>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Nacionalidad</label>
            <p className="px-3 py-2 bg-gray-100 rounded-md text-gray-900">{datos.nacionalidad}</p>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Estado Civil</label>
            <p className="px-3 py-2 bg-gray-100 rounded-md text-gray-900 capitalize">{datos.estado_civil}</p>
          </div>
        </div>

        {(datos.pareja_nombres || datos.pareja_cedula) && (
          <div className="mt-4 pt-4 border-t border-gray-200">
            <h4 className="text-sm font-medium text-gray-900 mb-3">Datos de Pareja</h4>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Nombres</label>
                <p className="px-3 py-2 bg-gray-100 rounded-md text-gray-900">{datos.pareja_nombres}</p>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Cédula</label>
                <p className="px-3 py-2 bg-gray-100 rounded-md text-gray-900">{datos.pareja_cedula}</p>
              </div>
            </div>
          </div>
        )}

        <div className="mt-4 p-4 bg-blue-50 border border-blue-200 rounded-lg">
          <p className="text-sm text-blue-800">
            <strong>Nota:</strong> Esta información se obtiene del Registro Civil y no puede ser modificada.
          </p>
        </div>
      </div>
    </Card>
  );
};

export default DatosPersonales;
