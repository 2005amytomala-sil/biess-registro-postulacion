import { useState, useEffect, useCallback, useMemo } from 'react';
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
import Button from '../components/Button';
import { obtenerPostulacion, guardarPostulacion, enviarPostulacion } from '../services/postulacionApi';
import { calcularFaltantes } from '../utils/completitud';

const DOMICILIO_INICIAL = {
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
};

// Errores del backend que significan que la postulación quedó cerrada (motivo exacto en el mensaje)
const CODIGOS_CIERRE = ['NACIONALIDAD_NO_PERMITIDA', 'SECTOR_NO_VALIDO', 'POSTULACION_CERRADA'];

const steps = [
  { id: 1, title: 'Información Personal', subtitle: 'Datos del Registro Civil' },
  { id: 2, title: 'Domicilio y Contacto', subtitle: 'Dirección y medios de contacto' },
  { id: 3, title: 'Información Adicional', subtitle: 'Familiares en IESS/BIESS' },
  { id: 4, title: 'Selección del Cargo', subtitle: 'Sector y respaldo' },
  { id: 5, title: 'Formación Académica', subtitle: 'Estudios realizados' },
  { id: 6, title: 'Experiencia Profesional', subtitle: 'Trayectoria laboral' },
  { id: 7, title: 'Inhabilidades', subtitle: 'Declaración de inhabilidades' },
  { id: 8, title: 'Consentimiento', subtitle: 'Envío de postulación' },
];

// Aviso destacado del estado de la postulación (cerrada o enviada)
const AvisoEstado = ({ tipo, titulo, detalle, nota }) => {
  const cerrada = tipo === 'cerrada';
  return (
    <div
      role={cerrada ? 'alert' : 'status'}
      className={`flex gap-4 rounded-xl border-l-4 p-4 shadow-sm sm:p-5 ${
        cerrada ? 'border-red-600 bg-red-50 text-red-900' : 'border-emerald-600 bg-emerald-50 text-emerald-900'
      }`}
    >
      <span
        className={`flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full ${
          cerrada ? 'bg-red-100 text-red-700' : 'bg-emerald-100 text-emerald-700'
        }`}
        aria-hidden="true"
      >
        <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          {cerrada ? (
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
          ) : (
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          )}
        </svg>
      </span>
      <div className="min-w-0">
        <p className="font-semibold">{titulo}</p>
        {detalle && <p className="mt-1 break-words text-sm">{detalle}</p>}
        {nota && <p className={`mt-1 text-xs ${cerrada ? 'text-red-700' : 'text-emerald-700'}`}>{nota}</p>}
      </div>
    </div>
  );
};

const PostulacionPage = () => {
  const navigate = useNavigate();
  const [currentStep, setCurrentStep] = useState(0);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [enviando, setEnviando] = useState(false);
  const [error, setError] = useState(null);
  const [successMessage, setSuccessMessage] = useState(null);
  const [postulacion, setPostulacion] = useState(null);
  const [motivoCierre, setMotivoCierre] = useState(null);
  const [erroresDomicilio, setErroresDomicilio] = useState({});
  const [faltantesServidor, setFaltantesServidor] = useState([]);

  const [domicilio, setDomicilio] = useState(DOMICILIO_INICIAL);
  const [familiares, setFamiliares] = useState({ tiene_familiares: false, familiares: [] });
  const [cargo, setCargo] = useState({
    sector: '',
    respaldo_organizacion: { tiene_respaldo: false, nombre_organizacion: '' },
  });
  const [formacion, setFormacion] = useState([]);
  const [experiencia, setExperiencia] = useState([]);
  const [inhabilidades, setInhabilidades] = useState([]);

  const enviada = postulacion?.estado === 'enviada';
  const cerrada = motivoCierre !== null || postulacion?.estado === 'cerrada';
  const soloLectura = enviada || cerrada;
  const faltantes = useMemo(() => calcularFaltantes(postulacion), [postulacion]);

  // Vuelca la postulación del servidor en el estado de cada sección
  const aplicarPostulacion = (data) => {
    setPostulacion(data);
    if (data.estado === 'cerrada') setMotivoCierre(data.motivo_cierre || 'Su postulación ha sido cerrada.');
    setDomicilio(data.domicilio ? { ...DOMICILIO_INICIAL, ...data.domicilio } : DOMICILIO_INICIAL);
    setFamiliares({ tiene_familiares: !!data.tiene_familiares, familiares: data.familiares || [] });
    setCargo({
      sector: data.respaldo_organizacion ? data.sector || '' : '',
      respaldo_organizacion: data.respaldo_organizacion || { tiene_respaldo: false, nombre_organizacion: '' },
    });
    setFormacion(data.formacion || []);
    setExperiencia(data.experiencia || []);
    setInhabilidades(data.inhabilidades || []);
  };

  // Trata las respuestas de error comunes: cierre (422/409) y postulación ya enviada (409)
  const manejarError = (result) => {
    if (CODIGOS_CIERRE.includes(result.codigo)) {
      setMotivoCierre(result.data?.motivo_cierre || result.message);
      return true;
    }
    if (result.codigo === 'POSTULACION_ENVIADA') {
      setPostulacion((prev) => ({ ...prev, estado: 'enviada', codigo: result.data?.codigo || prev?.codigo }));
      setError(result.message);
      return true;
    }
    setError(result.message);
    return false;
  };

  const cargarPostulacion = useCallback(async ({ silencioso = false } = {}) => {
    if (!silencioso) setLoading(true);
    const result = await obtenerPostulacion();
    if (result.success) {
      aplicarPostulacion(result.data);
    } else if (!silencioso) {
      manejarError(result);
    }
    if (!silencioso) setLoading(false);
  }, []);

  useEffect(() => {
    cargarPostulacion();
  }, [cargarPostulacion]);

  const irAPaso = (index) => {
    setError(null);
    setCurrentStep(index);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  // Guardado parcial: envía solo la sección indicada y refresca el estado para la completitud
  const guardarSeccion = async (datos, mensajeExito, { onFieldErrors } = {}) => {
    if (soloLectura) return;
    setSaving(true);
    setError(null);
    setSuccessMessage(null);
    onFieldErrors?.({});

    const result = await guardarPostulacion(datos);

    if (result.success) {
      setSuccessMessage(mensajeExito);
      setTimeout(() => setSuccessMessage(null), 3000);
      setFaltantesServidor([]);
      await cargarPostulacion({ silencioso: true });
    } else {
      if (result.data?.campo) onFieldErrors?.({ [result.data.campo]: result.message });
      manejarError(result);
    }
    setSaving(false);
  };

  const handleGuardarDomicilio = () =>
    guardarSeccion({ domicilio }, 'Domicilio guardado exitosamente', { onFieldErrors: setErroresDomicilio });

  const handleGuardarFamiliares = () =>
    guardarSeccion(
      { tiene_familiares: familiares.tiene_familiares, familiares: familiares.familiares },
      'Familiares guardados exitosamente'
    );

  const handleGuardarCargo = () =>
    guardarSeccion(
      { sector: cargo.sector, respaldo_organizacion: cargo.respaldo_organizacion },
      'Selección de cargo guardada exitosamente'
    );

  const handleGuardarFormacion = () => guardarSeccion({ formacion }, 'Formación académica guardada exitosamente');

  const handleGuardarExperiencia = () =>
    guardarSeccion({ experiencia }, 'Experiencia profesional guardada exitosamente');

  // Inhabilidades entrega los datos al guardar porque el estado del padre aún no se actualizó
  const handleGuardarInhabilidades = (datos) =>
    guardarSeccion(
      { inhabilidades: Array.isArray(datos) ? datos : inhabilidades },
      'Inhabilidades guardadas exitosamente'
    );

  // Inhabilidades no está aquí: mantiene sus respuestas internamente y guarda con su propio botón
  const guardadoPorPaso = {
    1: handleGuardarDomicilio,
    2: handleGuardarFamiliares,
    3: handleGuardarCargo,
    4: handleGuardarFormacion,
    5: handleGuardarExperiencia,
  };

  const handleEnviarPostulacion = async (aceptaConsentimiento) => {
    setEnviando(true);
    setError(null);
    setFaltantesServidor([]);

    const result = await enviarPostulacion(aceptaConsentimiento);

    if (result.success) {
      setPostulacion((prev) => ({ ...prev, estado: 'enviada', codigo: result.data.codigo }));
      setSuccessMessage('Postulación enviada exitosamente');
    } else if (result.codigo === 'POSTULACION_INCOMPLETA') {
      setFaltantesServidor(result.data?.faltantes || []);
      setError(result.message);
      await cargarPostulacion({ silencioso: true });
    } else {
      manejarError(result);
    }
    setEnviando(false);
  };

  if (loading) {
    return (
      <div className="flex min-h-[60vh] items-center justify-center px-4" role="status" aria-live="polite">
        <div className="text-center text-blue-700">
          <Spinner size="lg" className="mx-auto" />
          <p className="mt-4 text-sm font-medium text-slate-600">Cargando postulación...</p>
        </div>
      </div>
    );
  }

  // Cerrada sin datos (el backend rechaza la consulta): solo se muestra el motivo exacto
  if (cerrada && !postulacion) {
    return (
      <div className="px-4 py-10 sm:px-6">
        <div className="mx-auto max-w-3xl">
          <Card title="Postulación Cerrada">
            <div className="space-y-5">
              <AvisoEstado tipo="cerrada" titulo="Su postulación ha sido cerrada" detalle={motivoCierre} />
              <Button variant="secondary" onClick={() => navigate('/dashboard')}>
                Volver al inicio
              </Button>
            </div>
          </Card>
        </div>
      </div>
    );
  }

  if (!postulacion) {
    return (
      <div className="px-4 py-10 sm:px-6">
        <div className="mx-auto max-w-3xl space-y-4">
          <Alert type="error" message={error || 'No se pudo cargar la postulación'} />
          <Button onClick={() => cargarPostulacion()}>Reintentar</Button>
        </div>
      </div>
    );
  }

  const esUltimoPaso = currentStep === steps.length - 1;
  const guardarPasoActual = guardadoPorPaso[currentStep];

  return (
    <div className="px-4 py-6 sm:px-6 sm:py-8 lg:px-8">
      <div className="mx-auto max-w-6xl space-y-5">
        <header className="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <h1 className="text-2xl font-bold text-slate-900 sm:text-3xl">Formulario de Postulación</h1>
            <p className="mt-1 text-sm text-slate-600 sm:text-base">
              Complete los pasos para postular al concurso público de méritos del BIESS
            </p>
          </div>
          {postulacion.codigo && (
            <p className="text-xs text-slate-500">
              Código: <span className="font-mono text-slate-700">{postulacion.codigo}</span>
            </p>
          )}
        </header>

        {cerrada && (
          <AvisoEstado
            tipo="cerrada"
            titulo="Postulación cerrada"
            detalle={motivoCierre}
            nota="La postulación está cerrada y no admite cambios."
          />
        )}

        {enviada && (
          <AvisoEstado
            tipo="enviada"
            titulo="Postulación enviada"
            detalle={
              <>
                Código de postulación: <strong className="font-mono text-base">{postulacion.codigo}</strong>
              </>
            }
            nota="El formulario está en modo solo lectura."
          />
        )}

        {error && !cerrada && <Alert type="error" message={error} onClose={() => setError(null)} />}

        {successMessage && (
          <Alert type="success" message={successMessage} onClose={() => setSuccessMessage(null)} />
        )}

        <Stepper steps={steps} currentStep={currentStep} onStepClick={irAPaso} />

        <div>
          {currentStep === 0 && <DatosPersonales datos={postulacion.datos_personales} loading={false} />}

          {/* fieldset deshabilitado bloquea todos los campos y botones de las secciones editables */}
          <fieldset disabled={soloLectura || saving} className="contents">
            {currentStep === 1 && (
              <DomicilioContacto
                datos={domicilio}
                onChange={setDomicilio}
                onSave={handleGuardarDomicilio}
                saving={saving}
                serverErrors={erroresDomicilio}
              />
            )}

            {currentStep === 2 && (
              <FamiliaresIESS
                datos={familiares}
                onChange={setFamiliares}
                onSave={handleGuardarFamiliares}
                saving={saving}
              />
            )}

            {currentStep === 3 && (
              <SeleccionCargo datos={cargo} onChange={setCargo} onSave={handleGuardarCargo} saving={saving} />
            )}

            {currentStep === 4 && (
              <FormacionAcademica
                datos={formacion}
                onChange={setFormacion}
                onSave={handleGuardarFormacion}
                saving={saving}
              />
            )}

            {currentStep === 5 && (
              <ExperienciaProfesional
                datos={experiencia}
                onChange={setExperiencia}
                onSave={handleGuardarExperiencia}
                saving={saving}
              />
            )}

            {currentStep === 6 && (
              <Inhabilidades
                datos={inhabilidades}
                onChange={setInhabilidades}
                onSave={handleGuardarInhabilidades}
                saving={saving}
              />
            )}
          </fieldset>

          {currentStep === 7 && !cerrada && (
            <Consentimiento
              postulacion={postulacion}
              faltantes={faltantes}
              faltantesServidor={faltantesServidor}
              onEnviar={handleEnviarPostulacion}
              enviando={enviando}
              onIrAPaso={irAPaso}
            />
          )}
        </div>

        <div className="flex flex-col-reverse gap-3 border-t border-slate-200 pt-5 sm:flex-row sm:items-center sm:justify-between">
          <Button variant="secondary" onClick={() => irAPaso(currentStep - 1)} disabled={currentStep === 0}>
            <span aria-hidden="true">←</span> Anterior
          </Button>

          <div className="flex flex-col gap-3 sm:flex-row">
            {!soloLectura && guardarPasoActual && (
              <Button variant="outline" onClick={() => guardarPasoActual()} disabled={saving}>
                {saving ? (
                  <>
                    <Spinner size="sm" /> Guardando...
                  </>
                ) : (
                  'Guardar Progreso'
                )}
              </Button>
            )}
            <Button onClick={() => irAPaso(currentStep + 1)} disabled={esUltimoPaso}>
              Siguiente <span aria-hidden="true">→</span>
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default PostulacionPage;
