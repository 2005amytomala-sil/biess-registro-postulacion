import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Stepper from '../components/Stepper';
import DatosPersonales from '../components/DatosPersonales';
import Card from '../components/Card';
import Alert from '../components/Alert';
import Spinner from '../components/Spinner';
import { obtenerPostulacion, guardarPostulacion } from '../services/postulacionApi';

const PostulacionPage = () => {
  const navigate = useNavigate();
  const [currentStep, setCurrentStep] = useState(0);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState(null);
  const [successMessage, setSuccessMessage] = useState(null);
  const [postulacion, setPostulacion] = useState(null);
  const [datosPersonales, setDatosPersonales] = useState(null);
  const [postulacionCerrada, setPostulacionCerrada] = useState(false);

  const steps = [
    { id: 1, title: 'Información Personal', subtitle: 'Datos del Registro Civil' },
    { id: 2, title: 'Domicilio y Contacto', subtitle: 'Dirección y medios de contacto' },
    { id: 3, title: 'Información Adicional', subtitle: 'Familiares y formación' },
    { id: 4, title: 'Selección del Cargo', subtitle: 'Sector y respaldo' },
  ];

  // Cargar postulación al montar
  useEffect(() => {
    cargarPostulacion();
  }, []);

  const cargarPostulacion = async () => {
    setLoading(true);
    setError(null);

    const result = await obtenerPostulacion();

    if (result.success) {
      setPostulacion(result.data);
      // Los datos personales se obtienen del backend y se mostrarán en el componente
      setDatosPersonales({
        nombres: 'Juan Carlos', // Estos datos vendrán del backend
        apellidos: 'Pérez López',
        cedula: '0101010101',
        genero: 'masculino',
        fecha_nacimiento: '01/01/1980',
        nacionalidad: 'Ecuatoriana',
        estado_civil: 'soltero',
      });
    } else if (result.message.includes('no existe')) {
      // No existe postulación, crearemos una nueva al guardar
      setPostulacion(null);
    } else {
      setError(result.message);
    }

    setLoading(false);
  };

  const handleStepClick = (stepIndex) => {
    // Solo permitir ir hacia atrás o a pasos completados
    if (stepIndex <= currentStep) {
      setCurrentStep(stepIndex);
    }
  };

  const handleNext = () => {
    if (currentStep < steps.length - 1) {
      setCurrentStep(currentStep + 1);
    }
  };

  const handlePrevious = () => {
    if (currentStep > 0) {
      setCurrentStep(currentStep - 1);
    }
  };

  const handleGuardar = async () => {
    setSaving(true);
    setError(null);
    setSuccessMessage(null);

    // Simular guardado del progreso actual
    // En la implementación real, aquí se enviarían los datos de la sección actual
    const datosAGuardar = {
      sector: postulacion?.sector || 'afiliados',
      // Aquí se agregarían los datos de la sección actual
    };

    const result = await guardarPostulacion(datosAGuardar);

    if (result.success) {
      setSuccessMessage('Progreso guardado exitosamente');
      setTimeout(() => setSuccessMessage(null), 3000);
    } else {
      setError(result.message);
    }

    setSaving(false);
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <Spinner />
          <p className="mt-4 text-gray-600">Cargando postulación...</p>
        </div>
      </div>
    );
  }

  if (postulacionCerrada) {
    return (
      <div className="min-h-screen bg-gray-50 py-8 px-4">
        <div className="max-w-4xl mx-auto">
          <Card title="Postulación Cerrada">
            <div className="space-y-4">
              <Alert
                type="error"
                message="Su postulación ha sido cerrada. No puede continuar con el proceso."
              />
              <button
                onClick={() => navigate('/')}
                className="px-4 py-2 bg-gray-200 text-gray-900 rounded-md hover:bg-gray-300"
              >
                Volver al inicio
              </button>
            </div>
          </Card>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8 px-4">
      <div className="max-w-7xl mx-auto">
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900">Formulario de Postulación</h1>
          <p className="text-gray-600 mt-2">
            Complete los pasos para postular al concurso público de méritos del BIESS
          </p>
        </div>

        {error && (
          <Alert type="error" message={error} onClose={() => setError(null)} className="mb-4" />
        )}

        {successMessage && (
          <Alert
            type="success"
            message={successMessage}
            onClose={() => setSuccessMessage(null)}
            className="mb-4"
          />
        )}

        <Stepper steps={steps} currentStep={currentStep} onStepClick={handleStepClick} />

        <div className="mt-6">
          {currentStep === 0 && (
            <DatosPersonales
              datos={datosPersonales}
              loading={loading}
              error={error}
              onClose={() => setPostulacionCerrada(true)}
            />
          )}

          {currentStep === 1 && (
            <Card title="1.2 Domicilio y Contacto">
              <Alert type="info" message="Sección en construcción..." />
            </Card>
          )}

          {currentStep === 2 && (
            <Card title="1.3 Información Adicional">
              <Alert type="info" message="Sección en construcción..." />
            </Card>
          )}

          {currentStep === 3 && (
            <Card title="1.4 Selección del Cargo">
              <Alert type="info" message="Sección en construcción..." />
            </Card>
          )}
        </div>

        <div className="mt-6 flex justify-between">
          <button
            onClick={handlePrevious}
            disabled={currentStep === 0}
            className="px-4 py-2 bg-gray-200 text-gray-900 rounded-md hover:bg-gray-300 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Anterior
          </button>

          <div className="flex gap-3">
            <button
              onClick={handleGuardar}
              disabled={saving}
              className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {saving ? 'Guardando...' : 'Guardar Progreso'}
            </button>

            <button
              onClick={handleNext}
              disabled={currentStep === steps.length - 1}
              className="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {currentStep === steps.length - 1 ? 'Finalizar' : 'Siguiente'}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default PostulacionPage;
