import Card from './Card';
import Alert from './Alert';
import Spinner from './Spinner';
import Button from './Button';

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
            <Button variant="secondary" onClick={onClose}>
              Cerrar
            </Button>
          )}
        </div>
      </Card>
    );
  }

  return (
    <Card title="1.1 Información Personal">
      <div className="space-y-5">
        <div className="flex flex-wrap items-center gap-2 text-sm text-slate-600">
          <span className="inline-flex items-center gap-1.5 rounded-full bg-slate-100 px-2.5 py-1 text-xs font-medium text-slate-700 ring-1 ring-slate-200">
            <svg className="h-3.5 w-3.5" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
              <path fillRule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clipRule="evenodd" />
            </svg>
            Solo lectura
          </span>
          Datos obtenidos del Registro Civil
        </div>

        <dl className="grid grid-cols-1 gap-4 rounded-lg border border-dashed border-slate-300 bg-slate-50 p-4 sm:grid-cols-2 sm:p-5">
          <Dato etiqueta="Nombres" valor={datos.nombres} />
          <Dato etiqueta="Apellidos" valor={datos.apellidos} />
          <Dato etiqueta="Cédula" valor={datos.cedula} />
          <Dato etiqueta="Género" valor={datos.genero} capitalizar />
          <Dato etiqueta="Fecha de Nacimiento" valor={datos.fecha_nacimiento} />
          <Dato etiqueta="Nacionalidad" valor={datos.nacionalidad} />
          <Dato etiqueta="Estado Civil" valor={datos.estado_civil} capitalizar />
        </dl>

        {(datos.pareja_nombres || datos.pareja_cedula) && (
          <div>
            <h3 className="mb-3 text-sm font-semibold text-slate-900">Datos de Pareja</h3>
            <dl className="grid grid-cols-1 gap-4 rounded-lg border border-dashed border-slate-300 bg-slate-50 p-4 sm:grid-cols-2 sm:p-5">
              <Dato etiqueta="Nombres" valor={datos.pareja_nombres} />
              <Dato etiqueta="Cédula" valor={datos.pareja_cedula} />
            </dl>
          </div>
        )}

        <Alert type="info" message="Esta información se obtiene del Registro Civil y no puede ser modificada." />
      </div>
    </Card>
  );
};

const Dato = ({ etiqueta, valor, capitalizar = false }) => (
  <div>
    <dt className="text-xs font-medium uppercase tracking-wide text-slate-500">{etiqueta}</dt>
    <dd className={`mt-1 text-sm font-medium text-slate-900 ${capitalizar ? 'capitalize' : ''}`}>{valor || '—'}</dd>
  </div>
);

export default DatosPersonales;
