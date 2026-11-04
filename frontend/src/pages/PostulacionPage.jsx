import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Stepper from '../components/Stepper';
import DatosPersonales from '../components/DatosPersonales';
import DomicilioContacto from '../components/DomicilioContacto';
import FamiliaresIESS from '../components/FamiliaresIESS';
import SeleccionCargo from '../components/SeleccionCargo';
import FormacionAcademica from '../components/postulacion/FormacionAcademica';
import ExperienciaProfesional from '../components/postulacion/ExperienciaProfesional';
import Inhabilidades from '../components/postulacion/Inhabilidades';
import Consentimiento from '../components/postulacion/Consentimiento';
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
  const [domicilio, setDomicilio] = useState({
    pais_id: 1,
    provincia_id: null,
    canton_id: null,
    provincia_texto: '',
    ciudad_texto: '',
    celular: '',
    correo_1: '',
    correo_2: '',
    calle_principal: '',
    numero: '',
    calle_secundaria: '',
    sector: '',
    referencia: '',
    telefono_domicilio: '',
  });
  const [familiares, setFamiliares] = useState({
    tiene_familiares: false,
    familiares: [],
  });
  const [cargo, setCargo] = useState({
    sector: '',
    respaldo_organizacion: {
      tiene_respaldo: false,
      nombre_organizacion: '',
    },
  });
  const [formacion, setFormacion] = useState([]);
  const [experiencia, setExperiencia] = useState([]);
  const [inhabilidades, setInhabilidades] = useState([]);
  const [postulacionEnviada, setPostulacionEnviada] = useState(false);
  const [postulacionCerrada, setPostulacionCerrada] = useState(false);

  const steps = [
    { id: 1, title: 'Información Personal', subtitle: 'Datos del Registro Civil' },
    { id: 2, title: 'Domicilio y Contacto', subtitle: 'Dirección y medios de contacto' },
    { id: 3, title: 'Información Adicional', subtitle: 'Familiares y formación' },
    { id: 4, title: 'Selección del Cargo', subtitle: 'Sector y respaldo' },
    { id: 5, title: 'Formación Académica', subtitle: 'Estudios realizados' },
    { id: 6, title: 'Experiencia Profesional', subtitle: 'Trayectoria laboral' },
    { id: 7, title: 'Inhabilidades', subtitle: 'Declaración de inhabilidades' },
    { id: 8, title: 'Consentimiento', subtitle: 'Envío de postulación' },
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
      // Los datos personales se obtienen del backend (datos_personales del Registro Civil)
      if (result.data.datos_personales) {
        const dp = result.data.datos_personales;
        setDatosPersonales({
          nombres: dp.nombres,
          apellidos: dp.apellidos,
          cedula: dp.cedula,
          genero: dp.genero,
          fecha_nacimiento: dp.fecha_nacimiento,
          nacionalidad: dp.nacionalidad,
          estado_civil: dp.estado_civil,
          pareja_nombres: dp.pareja_nombres || '',
          pareja_cedula: dp.pareja_cedula || '',
        });
      }
      // Cargar domicilio si existe
      if (result.data.domicilio) {
        setDomicilio(result.data.domicilio);
      }
      // Cargar familiares si existe
      if (result.data.familiares) {
        setFamiliares({
          tiene_familiares: result.data.tiene_familiares,
          familiares: result.data.familiares,
        });
      }
      // Cargar cargo si existe
      if (result.data.respaldo_organizacion) {
        setCargo({
          sector: result.data.sector || 'afiliados',
          respaldo_organizacion: result.data.respaldo_organizacion,
        });
      }
      // Cargar formación si existe
      if (result.data.formacion) {
        setFormacion(result.data.formacion);
      }
      // Cargar experiencia si existe
      if (result.data.experiencia) {
        setExperiencia(result.data.experiencia);
      }
      // Cargar inhabilidades si existen
      if (result.data.inhabilidades) {
        setInhabilidades(result.data.inhabilidades);
      }
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

    // Guardar según el paso actual
    let datosAGuardar = {
      sector: postulacion?.sector || cargo.sector || 'afiliados',
    };

    if (currentStep === 1) {
      // Guardar domicilio
      datosAGuardar = {
        ...datosAGuardar,
        domicilio: domicilio,
      };
    } else if (currentStep === 2) {
      // Guardar familiares
      datosAGuardar = {
        ...datosAGuardar,
        tiene_familiares: familiares.tiene_familiares,
        familiares: familiares.familiares,
      };
    } else if (currentStep === 3) {
      // Guardar cargo
      datosAGuardar = {
        sector: cargo.sector,
        respaldo_organizacion: cargo.respaldo_organizacion,
      };
    }

    const result = await guardarPostulacion(datosAGuardar);

    if (result.success) {
      setSuccessMessage('Progreso guardado exitosamente');
      setTimeout(() => setSuccessMessage(null), 3000);
      // Actualizar el sector en la postulación si se guardó en paso 3
      if (currentStep === 3) {
        setPostulacion({ ...postulacion, sector: cargo.sector });
      }
    } else {
      setError(result.message);
      // Si el error es de sector no válido, cerrar la postulación
      if (result.codigo === 'SECTOR_NO_VALIDO') {
        setPostulacionCerrada(true);
      }
    }

    setSaving(false);
  };

  const handleDomicilioChange = (newDomicilio) => {
    setDomicilio(newDomicilio);
  };

  const handleGuardarDomicilio = async () => {
    setSaving(true);
    setError(null);
    setSuccessMessage(null);

    const datosAGuardar = {
      sector: postulacion?.sector || 'afiliados',
      domicilio: domicilio,
    };

    const result = await guardarPostulacion(datosAGuardar);

    if (result.success) {
      setSuccessMessage('Domicilio guardado exitosamente');
      setTimeout(() => setSuccessMessage(null), 3000);
    } else {
      setError(result.message);
    }

    setSaving(false);
  };

  const handleFamiliaresChange = (newFamiliares) => {
    setFamiliares(newFamiliares);
  };

  const handleGuardarFamiliares = async () => {
    setSaving(true);
    setError(null);
    setSuccessMessage(null);

    const datosAGuardar = {
      sector: postulacion?.sector || 'afiliados',
      tiene_familiares: familiares.tiene_familiares,
      familiares: familiares.familiares,
    };

    const result = await guardarPostulacion(datosAGuardar);

    if (result.success) {
      setSuccessMessage('Familiares guardados exitosamente');
      setTimeout(() => setSuccessMessage(null), 3000);
    } else {
      setError(result.message);
    }

    setSaving(false);
  };

  const handleCargoChange = (newCargo) => {
    setCargo(newCargo);
  };

  const handleGuardarCargo = async () => {
    setSaving(true);
    setError(null);
    setSuccessMessage(null);

    const datosAGuardar = {
      sector: cargo.sector,
      respaldo_organizacion: cargo.respaldo_organizacion,
    };

    const result = await guardarPostulacion(datosAGuardar);

    if (result.success) {
      setSuccessMessage('Selección de cargo guardada exitosamente');
      setTimeout(() => setSuccessMessage(null), 3000);
      // Actualizar el sector en la postulación
      setPostulacion({ ...postulacion, sector: cargo.sector });
    } else {
      setError(result.message);
      // Si el error es de sector no válido, cerrar la postulación
      if (result.codigo === 'SECTOR_NO_VALIDO') {
        setPostulacionCerrada(true);
      }
    }

    setSaving(false);
  };

  const handleFormacionChange = (newFormacion) => {
    setFormacion(newFormacion);
  };

  const handleGuardarFormacion = async () => {
    setSaving(true);
    setError(null);
    setSuccessMessage(null);

    const datosAGuardar = {
      sector: postulacion?.sector || 'afiliados',
      formacion: formacion,
    };

    const result = await guardarPostulacion(datosAGuardar);

    if (result.success) {
      setSuccessMessage('Formación académica guardada exitosamente');
      setTimeout(() => setSuccessMessage(null), 3000);
    } else {
      setError(result.message);
    }

    setSaving(false);
  };

  const handleExperienciaChange = (newExperiencia) => {
    setExperiencia(newExperiencia);
  };

  const handleGuardarExperiencia = async () => {
    setSaving(true);
    setError(null);
    setSuccessMessage(null);

    const datosAGuardar = {
      sector: postulacion?.sector || 'afiliados',
      experiencia: experiencia,
    };

    const result = await guardarPostulacion(datosAGuardar);

    if (result.success) {
      setSuccessMessage('Experiencia profesional guardada exitosamente');
      setTimeout(() => setSuccessMessage(null), 3000);
    } else {
      setError(result.message);
    }

    setSaving(false);
  };

  const handleInhabilidadesChange = (newInhabilidades) => {
    setInhabilidades(newInhabilidades);
  };

  const handleGuardarInhabilidades = async () => {
    setSaving(true);
    setError(null);
    setSuccessMessage(null);

    const datosAGuardar = {
      sector: postulacion?.sector || 'afiliados',
      inhabilidades: inhabilidades,
    };

    const result = await guardarPostulacion(datosAGuardar);

    if (result.success) {
      setSuccessMessage('Inhabilidades guardadas exitosamente');
      setTimeout(() => setSuccessMessage(null), 3000);
    } else {
      setError(result.message);
    }

    setSaving(false);
  };

  const handleEnviarPostulacion = () => {
    setPostulacionEnviada(true);
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
            <DomicilioContacto
              datos={domicilio}
              onChange={handleDomicilioChange}
              onSave={handleGuardarDomicilio}
              saving={saving}
              error={error}
            />
          )}

          {currentStep === 2 && (
            <FamiliaresIESS
              datos={familiares}
              onChange={handleFamiliaresChange}
              onSave={handleGuardarFamiliares}
              saving={saving}
              error={error}
            />
          )}

          {currentStep === 3 && (
            <SeleccionCargo
              datos={cargo}
              onChange={handleCargoChange}
              onSave={handleGuardarCargo}
              saving={saving}
              error={error}
            />
          )}

          {currentStep === 4 && (
            <FormacionAcademica
              datos={formacion}
              onChange={handleFormacionChange}
              onSave={handleGuardarFormacion}
              saving={saving}
              error={error}
            />
          )}

          {currentStep === 5 && (
            <ExperienciaProfesional
              datos={experiencia}
              onChange={handleExperienciaChange}
              onSave={handleGuardarExperiencia}
              saving={saving}
              error={error}
            />
          )}

          {currentStep === 6 && (
            <Inhabilidades
              datos={inhabilidades}
              onChange={handleInhabilidadesChange}
              onSave={handleGuardarInhabilidades}
              saving={saving}
              error={error}
            />
          )}

          {currentStep === 7 && (
            <Consentimiento
              postulacion={postulacion}
              onEnviar={handleEnviarPostulacion}
              saving={saving}
              error={error}
            />
          )}
        </div>

        <div className="mt-6 flex justify-between">
          <button
            onClick={handlePrevious}
            disabled={currentStep === 0 || postulacionEnviada}
            className="px-4 py-2 bg-gray-200 text-gray-900 rounded-md hover:bg-gray-300 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Anterior
          </button>

          {!postulacionEnviada && (
            <div className="flex gap-3">
              <button
                onClick={handleGuardar}
                disabled={saving || currentStep === 7}
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
          )}
        </div>
      </div>
    </div>
  );
};

export default PostulacionPage;
