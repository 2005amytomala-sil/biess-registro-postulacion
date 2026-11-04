import { useState, useEffect } from 'react';
import Button from './Button';
import Card from './Card';
import Input from './Input';
import Select from './Select';
import Spinner from './Spinner';
import Alert from './Alert';
import { obtenerPaises, obtenerProvincias, obtenerCantonesPorProvincia } from '../services/postulacionApi';

const REGEX_CORREO = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

const validarCampo = (field, value) => {
  switch (field) {
    case 'celular':
      return /^\d{10}$/.test(value || '') ? '' : 'El celular debe tener exactamente 10 dígitos';
    case 'correo_1':
      return value && REGEX_CORREO.test(value) ? '' : 'Ingrese un correo válido';
    case 'correo_2':
      return !value || REGEX_CORREO.test(value) ? '' : 'Ingrese un correo válido';
    default:
      return '';
  }
};

const SIN_ERRORES = {};

const DomicilioContacto = ({ datos, onChange, onSave, saving, serverErrors = SIN_ERRORES }) => {
  const [paises, setPaises] = useState([]);
  const [provincias, setProvincias] = useState([]);
  const [cantones, setCantones] = useState([]);
  const [loadingPaises, setLoadingPaises] = useState(false);
  const [loadingProvincias, setLoadingProvincias] = useState(false);
  const [loadingCantones, setLoadingCantones] = useState(false);
  const [errorCatalogo, setErrorCatalogo] = useState('');
  const [errors, setErrors] = useState({});
  const [erroresServidor, setErroresServidor] = useState(serverErrors);

  // Los errores del servidor llegan por campo tras cada guardado
  useEffect(() => {
    setErroresServidor(serverErrors);
  }, [serverErrors]);

  // Ecuador se identifica por nombre en el catálogo (no por un id fijo)
  const ecuadorId = paises.find((p) => p.nombre?.trim().toLowerCase() === 'ecuador')?.id ?? 1;
  const esEcuador = datos.pais_id === ecuadorId;

  // Cargar países al montar
  useEffect(() => {
    cargarPaises();
  }, []);

  // Las cargas de catálogo solo traen opciones; nunca modifican los datos del formulario
  useEffect(() => {
    if (esEcuador) {
      cargarProvincias();
    } else {
      setProvincias([]);
    }
  }, [esEcuador]);

  useEffect(() => {
    if (esEcuador && datos.provincia_id) {
      cargarCantones(datos.provincia_id);
    } else {
      setCantones([]);
    }
  }, [esEcuador, datos.provincia_id]);

  const cargarCatalogo = async (cargar, setLoading, setDatos) => {
    setLoading(true);
    const result = await cargar();
    if (result.success) {
      setDatos(result.data || []);
    } else {
      setErrorCatalogo(result.message);
    }
    setLoading(false);
  };

  const cargarPaises = () => cargarCatalogo(obtenerPaises, setLoadingPaises, setPaises);
  const cargarProvincias = () => cargarCatalogo(obtenerProvincias, setLoadingProvincias, setProvincias);
  const cargarCantones = (provinciaId) =>
    cargarCatalogo(() => obtenerCantonesPorProvincia(provinciaId), setLoadingCantones, setCantones);

  // Aplica cambios al formulario, valida en vivo y limpia el error del servidor de los campos tocados
  const actualizar = (cambios) => {
    const newErrors = { ...errors };
    Object.entries(cambios).forEach(([field, value]) => {
      const msg = validarCampo(field, value);
      if (msg) newErrors[field] = msg;
      else delete newErrors[field];
    });
    setErrors(newErrors);

    setErroresServidor((prev) => {
      const next = { ...prev };
      Object.keys(cambios).forEach((field) => delete next[field]);
      return next;
    });

    onChange({ ...datos, ...cambios });
  };

  const handleChange = (field, value) => actualizar({ [field]: value });

  // Al cambiar de país se reinician provincia, cantón y los campos de texto en el mismo cambio
  const handlePaisChange = (value) => {
    actualizar({
      pais_id: value ? parseInt(value, 10) : null,
      provincia_id: null,
      canton_id: null,
      provincia_texto: '',
      ciudad_texto: '',
    });
  };

  // Al cambiar de provincia se reinicia el cantón en el mismo cambio
  const handleProvinciaChange = (value) => {
    actualizar({ provincia_id: value ? parseInt(value, 10) : null, canton_id: null });
  };

  const handleCantonChange = (value) => {
    actualizar({ canton_id: value ? parseInt(value, 10) : null });
  };

  const errorDe = (field) => errors[field] || erroresServidor[field];

  return (
    <Card title="1.2 Domicilio y Contacto">
      <div className="space-y-6">
        {errorCatalogo && <Alert type="error" message={errorCatalogo} onClose={() => setErrorCatalogo('')} />}

        {/* Cascada geográfica */}
        <div className="space-y-4">
          <h3 className="text-sm font-semibold uppercase tracking-wide text-slate-500">Ubicación</h3>
          
          <Select
            label="País *"
            options={paises.map(p => ({ value: p.id, label: p.nombre }))}
            value={datos.pais_id || ''}
            onChange={(e) => handlePaisChange(e.target.value)}
            loading={loadingPaises}
            error={errorDe('pais_id')}
            required
          />

          {esEcuador ? (
            <>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <Select
                    label="Provincia *"
                    options={provincias.map(p => ({ value: p.id, label: p.nombre }))}
                    value={datos.provincia_id || ''}
                    onChange={(e) => handleProvinciaChange(e.target.value)}
                    loading={loadingProvincias}
                    error={errorDe('provincia_id')}
                    required
                  />
                  {loadingProvincias && (
                    <div className="mt-2 flex items-center text-sm text-slate-500">
                      <Spinner size="small" />
                      <span className="ml-2">Cargando provincias...</span>
                    </div>
                  )}
                </div>

                <div>
                  <Select
                    label="Cantón *"
                    options={cantones.map(c => ({ value: c.id, label: c.nombre }))}
                    value={datos.canton_id || ''}
                    onChange={(e) => handleCantonChange(e.target.value)}
                    loading={loadingCantones}
                    error={errorDe('canton_id')}
                    disabled={!datos.provincia_id}
                    required
                  />
                  {loadingCantones && (
                    <div className="mt-2 flex items-center text-sm text-slate-500">
                      <Spinner size="small" />
                      <span className="ml-2">Cargando cantones...</span>
                    </div>
                  )}
                </div>
              </div>
            </>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <Input
                label="Provincia *"
                value={datos.provincia_texto || ''}
                onChange={(e) => handleChange('provincia_texto', e.target.value)}
                error={errorDe('provincia_texto')}
                required
              />
              <Input
                label="Ciudad *"
                value={datos.ciudad_texto || ''}
                onChange={(e) => handleChange('ciudad_texto', e.target.value)}
                error={errorDe('ciudad_texto')}
                required
              />
            </div>
          )}
        </div>

        {/* Contacto */}
        <div className="space-y-4">
          <h3 className="text-sm font-semibold uppercase tracking-wide text-slate-500">Contacto</h3>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Input
              label="Celular *"
              type="tel"
              value={datos.celular || ''}
              onChange={(e) => handleChange('celular', e.target.value.replace(/\D/g, ''))}
              inputMode="numeric"
              maxLength={10}
              error={errorDe('celular')}
              placeholder="0991234567"
              required
            />
            <Input
              label="Correo Electrónico 1 *"
              type="email"
              value={datos.correo_1 || ''}
              onChange={(e) => handleChange('correo_1', e.target.value)}
              error={errorDe('correo_1')}
              placeholder="ejemplo@correo.com"
              required
            />
            <Input
              label="Correo Electrónico 2"
              type="email"
              value={datos.correo_2 || ''}
              onChange={(e) => handleChange('correo_2', e.target.value)}
              error={errorDe('correo_2')}
              placeholder="ejemplo2@correo.com"
            />
          </div>
        </div>

        {/* Dirección detallada */}
        <div className="space-y-4">
          <h3 className="text-sm font-semibold uppercase tracking-wide text-slate-500">Dirección Detallada</h3>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Input
              label="Calle Principal"
              value={datos.calle_principal || ''}
              onChange={(e) => handleChange('calle_principal', e.target.value)}
            />
            <Input
              label="Número"
              value={datos.numero || ''}
              onChange={(e) => handleChange('numero', e.target.value)}
            />
            <Input
              label="Calle Secundaria"
              value={datos.calle_secundaria || ''}
              onChange={(e) => handleChange('calle_secundaria', e.target.value)}
            />
            <Input
              label="Sector"
              value={datos.sector || ''}
              onChange={(e) => handleChange('sector', e.target.value)}
            />
          </div>

          <Input
            label="Referencia"
            value={datos.referencia || ''}
            onChange={(e) => handleChange('referencia', e.target.value)}
            placeholder="Ej: Casa color azul, cerca del parque"
          />

          <Input
            label="Teléfono de Domicilio"
            type="tel"
            value={datos.telefono_domicilio || ''}
            onChange={(e) => handleChange('telefono_domicilio', e.target.value)}
            placeholder="0221234567"
          />
        </div>

        <div className="flex justify-end border-t border-slate-100 pt-5">
          <Button
            onClick={onSave}
            disabled={saving || Object.keys(errors).length > 0}
          >
            {saving ? 'Guardando...' : 'Guardar'}
          </Button>
        </div>
      </div>
    </Card>
  );
};

export default DomicilioContacto;
