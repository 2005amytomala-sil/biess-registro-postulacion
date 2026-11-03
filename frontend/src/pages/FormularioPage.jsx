import Card from '../components/Card';
import Alert from '../components/Alert';

const FormularioPage = () => {
  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
        <Alert type="info" message="Formulario de postulación en construcción..." />
        <Card title="Formulario de Postulación">
          <div className="space-y-6">
            <p className="text-gray-600">
              Aquí se implementará el formulario de postulación al concurso público
              de méritos del BIESS.
            </p>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="border rounded-lg p-4">
                <h3 className="font-semibold text-gray-900 mb-2">Datos Personales</h3>
                <p className="text-sm text-gray-600">Pendiente de implementación</p>
              </div>
              <div className="border rounded-lg p-4">
                <h3 className="font-semibold text-gray-900 mb-2">Domicilio y Contacto</h3>
                <p className="text-sm text-gray-600">Pendiente de implementación</p>
              </div>
              <div className="border rounded-lg p-4">
                <h3 className="font-semibold text-gray-900 mb-2">Formación Académica</h3>
                <p className="text-sm text-gray-600">Pendiente de implementación</p>
              </div>
              <div className="border rounded-lg p-4">
                <h3 className="font-semibold text-gray-900 mb-2">Experiencia Profesional</h3>
                <p className="text-sm text-gray-600">Pendiente de implementación</p>
              </div>
            </div>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default FormularioPage;
