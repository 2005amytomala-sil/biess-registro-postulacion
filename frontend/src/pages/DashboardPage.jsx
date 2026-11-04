import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Card from '../components/Card';
import Button from '../components/Button';
import Alert from '../components/Alert';
import Spinner from '../components/Spinner';
import { obtenerPostulacion } from '../services/postulacionApi';

const DashboardPage = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [postulacion, setPostulacion] = useState(null);
  const [error, setError] = useState(null);
  const [fases, setFases] = useState([
    { id: 1, nombre: 'Información Personal', estado: 'pendiente' },
    { id: 2, nombre: 'Domicilio y Contacto', estado: 'pendiente' },
    { id: 3, nombre: 'Información Adicional', estado: 'pendiente' },
    { id: 4, nombre: 'Selección del Cargo', estado: 'pendiente' },
    { id: 5, nombre: 'Formación Académica', estado: 'pendiente' },
    { id: 6, nombre: 'Experiencia Profesional', estado: 'pendiente' },
    { id: 7, nombre: 'Inhabilidades', estado: 'pendiente' },
    { id: 8, nombre: 'Consentimiento y Envío', estado: 'pendiente' },
  ]);

  useEffect(() => {
    cargarEstado();
  }, []);

  const cargarEstado = async () => {
    setLoading(true);
    setError(null);

    const result = await obtenerPostulacion();

    if (result.success) {
      setPostulacion(result.data);
      
      // Actualizar estado de fases según lo completado
      const nuevasFases = [...fases];
      
      if (result.data.datos_personales) nuevasFases[0].estado = 'completado';
      if (result.data.domicilio) nuevasFases[1].estado = 'completado';
      if (result.data.familiares) nuevasFases[2].estado = 'completado';
      if (result.data.respaldo_organizacion) nuevasFases[3].estado = 'completado';
      if (result.data.formacion && result.data.formacion.length > 0) nuevasFases[4].estado = 'completado';
      if (result.data.experiencia && result.data.experiencia.length > 0) nuevasFases[5].estado = 'completado';
      if (result.data.inhabilidades && result.data.inhabilidades.length > 0) nuevasFases[6].estado = 'completado';
      
      if (result.data.estado === 'enviada') {
        nuevasFases[7].estado = 'completado';
      }
      
      setFases(nuevasFases);
    } else if (!result.message.includes('no existe')) {
      setError(result.message);
    }

    setLoading(false);
  };

  const getFaseColor = (estado) => {
    switch (estado) {
      case 'completado': return 'bg-green-100 text-green-800 border-green-300';
      case 'en_progreso': return 'bg-blue-100 text-blue-800 border-blue-300';
      default: return 'bg-slate-100 text-slate-800 border-slate-300';
    }
  };

  const getFaseIcon = (estado) => {
    switch (estado) {
      case 'completado': return '✓';
      case 'en_progreso': return '→';
      default: return '○';
    }
  };

  const handleIniciar = () => {
    // Encontrar la primera fase pendiente
    const fasePendiente = fases.find(f => f.estado === 'pendiente');
    if (fasePendiente) {
      navigate('/postulacion');
    } else {
      // Si todas están completadas pero no enviada, ir a consentimiento
      if (postulacion && postulacion.estado !== 'enviada') {
        navigate('/postulacion');
      } else {
        navigate('/postulacion');
      }
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-slate-50 flex items-center justify-center">
        <div className="text-center">
          <Spinner />
          <p className="mt-4 text-slate-600">Cargando estado de postulación...</p>
        </div>
      </div>
    );
  }

  const fasesCompletadas = fases.filter(f => f.estado === 'completado').length;
  const progreso = Math.round((fasesCompletadas / fases.length) * 100);

  return (
    <div className="min-h-screen bg-slate-50 py-8 px-4">
      <div className="max-w-4xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-slate-900">Bienvenido al Sistema de Postulación</h1>
          <p className="text-slate-600 mt-2">
            Progreso de su postulación al concurso público de méritos del BIESS
          </p>
        </div>

        {error && <Alert type="error" message={error} onClose={() => setError(null)} className="mb-4" />}

        {/* Resumen de progreso */}
        <Card title="Estado de su Postulación">
          <div className="space-y-6">
            {postulacion && postulacion.estado === 'enviada' ? (
              <div className="bg-green-50 border border-green-200 rounded-lg p-6 text-center">
                <svg className="w-16 h-16 text-green-600 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <h3 className="text-xl font-semibold text-green-900 mb-2">
                  ¡Su postulación ha sido enviada!
                </h3>
                <p className="text-green-800 mb-4">
                  Código de postulación: <strong>{postulacion.codigo}</strong>
                </p>
                <Button onClick={() => navigate('/postulacion')}>
                  Ver Detalles
                </Button>
              </div>
            ) : (
              <>
                <div className="flex items-center justify-between mb-4">
                  <div>
                    <p className="text-sm text-slate-600">Progreso general</p>
                    <p className="text-2xl font-bold text-slate-900">{progreso}%</p>
                  </div>
                  <div className="w-32 bg-slate-200 rounded-full h-4">
                    <div
                      className="bg-blue-600 h-4 rounded-full transition-all"
                      style={{ width: `${progreso}%` }}
                    />
                  </div>
                </div>

                <div className="space-y-3">
                  {fases.map((fase) => (
                    <div
                      key={fase.id}
                      className={`flex items-center justify-between p-4 rounded-lg border ${getFaseColor(fase.estado)}`}
                    >
                      <div className="flex items-center">
                        <span className="text-xl font-bold mr-3">{getFaseIcon(fase.estado)}</span>
                        <span className="font-medium">{fase.nombre}</span>
                      </div>
                      <span className="text-sm capitalize">
                        {fase.estado === 'completado' ? 'Completado' : 'Pendiente'}
                      </span>
                    </div>
                  ))}
                </div>

                <div className="flex justify-end mt-6">
                  <Button onClick={handleIniciar}>
                    {fasesCompletadas === 0 ? 'Iniciar Postulación' : 'Continuar Postulación'}
                  </Button>
                </div>
              </>
            )}
          </div>
        </Card>
      </div>
    </div>
  );
};

export default DashboardPage;
