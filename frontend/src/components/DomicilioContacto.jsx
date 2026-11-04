import { useState, useEffect } from 'react';
import Card from './Card';
import Input from './Input';
import Select from './Select';
import Spinner from './Spinner';
import Alert from './Alert';
import { obtenerPaises, obtenerProvincias, obtenerCantonesPorProvincia } from '../services/postulacionApi';

const DomicilioContacto = ({ datos, onChange, onSave, saving, error }) => {
  const [paises, setPaises] = useState([]);
  const [provincias, setProvincias] = useState([]);
  const [cantones, setCantones] = useState([]);
  const [loadingPaises, setLoadingPaises] = useState(false);
  const [loadingProvincias, setLoadingProvincias] = useState(false);
  const [loadingCantones, setLoadingCantones] = useState(false);
  const [errors, setErrors] = useState({});

  // Cargar países al montar
  useEffect(() => {
    cargarPaises();
  }, []);

  // Cargar provincias cuando se selecciona Ecuador
  useEffect(() => {
    if (datos.pais_id === 1) {
      cargarProvincias();
    } else {
      setProvincias([]);
      setCantones([]);
      // Limpiar provincia y cantón si no es Ecuador
      if (onChange) {
        onChange({ ...datos, provincia_id: null, canton_id: null });
      }
    }
  }, [datos.pais_id]);

  // Cargar cantones cuando se selecciona provincia
  useEffect(() => {
    if (datos.provincia_id) {
      cargarCantones(datos.provincia_id);
    } else {
      setCantones([]);
      // Limpiar cantón si no hay provincia
      if (onChange) {
        onChange({ ...datos, canton_id: null });
      }
    }
  }, [datos.provincia_id]);

  const cargarPaises = async () => {
    setLoadingPaises(true);
    const result = await obtenerPaises();
    if (result.success) {
      setPaises(result.data);
    }
    setLoadingPaises(false);
  };

  const cargarProvincias = async () => {
    setLoadingProvincias(true);
    const result = await obtenerProvincias();
    if (result.success) {
      setProvincias(result.data);
    }
    setLoadingProvincias(false);
  };

  const cargarCantones = async (provinciaId) => {
    setLoadingCantones(true);
    const result = await obtenerCantonesPorProvincia(provinciaId);
    if (result.success) {
      setCantones(result.data);
    }
    setLoadingCantones(false);
  };

  const handleChange = (field, value) => {
    const newDatos = { ...datos, [field]: value };
    
    // Validaciones en vivo
    const newErrors = { ...errors };
    
    if (field === 'celular') {
      if (!value || value.length < 10) {
        newErrors.celular = 'El celular debe tener al menos 10 dígitos';
      } else {
        delete newErrors.celular;
      }
    }
    
    if (field === 'correo_1') {
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      if (!value || !emailRegex.test(value)) {
        newErrors.correo_1 = 'Ingrese un correo válido';
      } else {
        delete newErrors.correo_1;
      }
    }
    
    if (field === 'correo_2' && value) {
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      if (!emailRegex.test(value)) {
        newErrors.correo_2 = 'Ingrese un correo válido';
      } else {
        delete newErrors.correo_2;
      }
    }
    
    setErrors(newErrors);
    
    if (onChange) {
      onChange(newDatos);
    }
  };

  const handlePaisChange = (value) => {
    const paisId = parseInt(value);
    handleChange('pais_id', paisId);
    
    // Si no es Ecuador, limpiar provincia y cantón
    if (paisId !== 1) {
      handleChange('provincia_texto', '');
      handleChange('ciudad_texto', '');
    }
  };

  const handleProvinciaChange = (value) => {
    const provinciaId = value ? parseInt(value) : null;
    handleChange('provincia_id', provinciaId);
  };

  const handleCantonChange = (value) => {
    const cantonId = value ? parseInt(value) : null;
    handleChange('canton_id', cantonId);
  };

  const esEcuador = datos.pais_id === 1;

  return (
    <Card title="1.2 Domicilio y Contacto">
      <div className="space-y-6">
        {error && <Alert type="error" message={error} onClose={() => {}} />}

        {/* Cascada geográfica */}
        <div className="space-y-4">
          <h3 className="text-sm font-medium text-gray-900">Ubicación</h3>
          
          <Select
            label="País *"
            options={paises.map(p => ({ value: p.id, label: p.nombre }))}
            value={datos.pais_id || ''}
            onChange={(e) => handlePaisChange(e.target.value)}
            loading={loadingPaises}
            error={errors.pais_id}
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
                    error={errors.provincia_id}
                    disabled={!datos.pais_id}
                    required
                  />
                  {loadingProvincias && (
                    <div className="mt-2 flex items-center text-sm text-gray-500">
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
                    error={errors.canton_id}
                    disabled={!datos.provincia_id}
                    required
                  />
                  {loadingCantones && (
                    <div className="mt-2 flex items-center text-sm text-gray-500">
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
                error={errors.provincia_texto}
                required
              />
              <Input
                label="Ciudad *"
                value={datos.ciudad_texto || ''}
                onChange={(e) => handleChange('ciudad_texto', e.target.value)}
                error={errors.ciudad_texto}
                required
              />
            </div>
          )}
        </div>

        {/* Contacto */}
        <div className="space-y-4">
          <h3 className="text-sm font-medium text-gray-900">Contacto</h3>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Input
              label="Celular *"
              type="tel"
              value={datos.celular || ''}
              onChange={(e) => handleChange('celular', e.target.value)}
              error={errors.celular}
              placeholder="0991234567"
              required
            />
            <Input
              label="Correo Electrónico 1 *"
              type="email"
              value={datos.correo_1 || ''}
              onChange={(e) => handleChange('correo_1', e.target.value)}
              error={errors.correo_1}
              placeholder="ejemplo@correo.com"
              required
            />
            <Input
              label="Correo Electrónico 2"
              type="email"
              value={datos.correo_2 || ''}
              onChange={(e) => handleChange('correo_2', e.target.value)}
              error={errors.correo_2}
              placeholder="ejemplo2@correo.com"
            />
          </div>
        </div>

        {/* Dirección detallada */}
        <div className="space-y-4">
          <h3 className="text-sm font-medium text-gray-900">Dirección Detallada</h3>
          
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

        <div className="flex justify-end">
          <button
            onClick={onSave}
            disabled={saving || Object.keys(errors).length > 0}
            className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {saving ? 'Guardando...' : 'Guardar'}
          </button>
        </div>
      </div>
    </Card>
  );
};

export default DomicilioContacto;
