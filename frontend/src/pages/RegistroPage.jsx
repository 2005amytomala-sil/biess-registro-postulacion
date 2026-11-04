import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import Button from '../components/Button';
import Input from '../components/Input';
import Select from '../components/Select';
import Alert from '../components/Alert';
import Card from '../components/Card';
import AuthLayout from '../components/AuthLayout';
import CuentaRegresiva from '../components/CuentaRegresiva';
import useCuentaRegresiva from '../hooks/useCuentaRegresiva';

const RegistroPage = () => {
  const [paso, setPaso] = useState(1);
  const [identidadData, setIdentidadData] = useState(null);
  const [formData, setFormData] = useState({
    cedula: '',
    tipoValidacion: 'fecha',
    fechaExpedicion: '',
    codigoDactilar: '',
    correo: '',
    password: '',
    confirmPassword: '',
  });
  const [errors, setErrors] = useState({});
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { verificarIdentidad, register } = useAuth();
  const { bloqueado, restante, iniciar: iniciarBloqueo } = useCuentaRegresiva();
  const navigate = useNavigate();

  const manejarError = (result) => {
    setError(result.message);
    if (result.retryAfter) {
      iniciarBloqueo(result.retryAfter);
      setPaso(1);
    }
  };

  const validateCedula = (cedula) => {
    if (!cedula) return 'La cédula es requerida';
    if (!/^\d{10}$/.test(cedula)) return 'La cédula debe tener 10 dígitos';
    return '';
  };

  const validateFechaExpedicion = (fecha) => {
    if (!fecha) return 'La fecha de expedición es requerida';
    if (!/^\d{2}\/\d{2}\/\d{4}$/.test(fecha)) return 'Formato inválido. Use DD/MM/AAAA';
    return '';
  };

  const validateCodigoDactilar = (codigo) => {
    if (!codigo) return 'El código dactilar es requerido';
    if (!/^\d{10}$/.test(codigo)) return 'El código dactilar debe tener exactamente 10 dígitos';
    return '';
  };

  const validateCorreo = (correo) => {
    if (!correo) return 'El correo es requerido';
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(correo)) return 'Correo inválido';
    return '';
  };

  const validatePassword = (password) => {
    if (!password) return 'La contraseña es requerida';
    if (password.length < 8) return 'La contraseña debe tener al menos 8 caracteres';
    return '';
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });

    // Validaciones inmediatas
    let error = '';
    switch (name) {
      case 'cedula':
        error = validateCedula(value);
        break;
      case 'fechaExpedicion':
        error = validateFechaExpedicion(value);
        break;
      case 'codigoDactilar':
        error = validateCodigoDactilar(value);
        break;
      case 'correo':
        error = validateCorreo(value);
        break;
      case 'password':
        error = validatePassword(value);
        break;
      case 'confirmPassword':
        error = value !== formData.password ? 'Las contraseñas no coinciden' : '';
        break;
      default:
        break;
    }
    setErrors({ ...errors, [name]: error });
  };

  const handleVerificarIdentidad = async (e) => {
    e.preventDefault();
    setError('');

    const cedulaError = validateCedula(formData.cedula);
    if (cedulaError) {
      setErrors({ ...errors, cedula: cedulaError });
      return;
    }

    if (formData.tipoValidacion === 'fecha') {
      const fechaError = validateFechaExpedicion(formData.fechaExpedicion);
      if (fechaError) {
        setErrors({ ...errors, fechaExpedicion: fechaError });
        return;
      }
    } else {
      const codigoError = validateCodigoDactilar(formData.codigoDactilar);
      if (codigoError) {
        setErrors({ ...errors, codigoDactilar: codigoError });
        return;
      }
    }

    setLoading(true);
    const result = await verificarIdentidad(
      formData.cedula,
      formData.tipoValidacion === 'fecha' ? formData.fechaExpedicion : '',
      formData.tipoValidacion === 'codigo' ? formData.codigoDactilar : ''
    );
    setLoading(false);

    if (result.success) {
      setIdentidadData(result.data);
      setPaso(2);
    } else {
      manejarError(result);
    }
  };

  const handleRegistro = async (e) => {
    e.preventDefault();
    setError('');

    const correoError = validateCorreo(formData.correo);
    const passwordError = validatePassword(formData.password);
    const confirmError = formData.password !== formData.confirmPassword ? 'Las contraseñas no coinciden' : '';

    if (correoError || passwordError || confirmError) {
      setErrors({
        ...errors,
        correo: correoError,
        password: passwordError,
        confirmPassword: confirmError,
      });
      return;
    }

    setLoading(true);
    const result = await register({
      cedula: formData.cedula,
      fecha_expedicion: formData.tipoValidacion === 'fecha' ? formData.fechaExpedicion : '',
      codigo_dactilar: formData.tipoValidacion === 'codigo' ? formData.codigoDactilar : '',
      correo: formData.correo,
      password: formData.password,
      confirm_password: formData.confirmPassword,
    });
    setLoading(false);

    if (result.success) {
      navigate('/login', {
        replace: true,
        state: { mensaje: 'Registro exitoso. Inicie sesión con su cédula y contraseña.' },
      });
    } else {
      manejarError(result);
    }
  };

  return (
    <AuthLayout>
        <Card title={paso === 1 ? 'Verificar Identidad' : 'Registro de Usuario'}>
          <div className="mb-5 space-y-3 empty:hidden">
            {bloqueado ? (
              <CuentaRegresiva mensaje={error} restante={restante} />
            ) : (
              error && <Alert type="error" message={error} />
            )}
          </div>
          
          {paso === 1 ? (
            <form onSubmit={handleVerificarIdentidad} className="space-y-5">
              <Input
                label="Cédula"
                name="cedula"
                type="text"
                value={formData.cedula}
                onChange={handleChange}
                placeholder="Ingrese su cédula (10 dígitos)"
                error={errors.cedula}
                required
              />
              <Select
                label="Método de validación"
                name="tipoValidacion"
                value={formData.tipoValidacion}
                onChange={handleChange}
                options={[
                  { value: 'fecha', label: 'Fecha de expedición (DD/MM/AAAA)' },
                  { value: 'codigo', label: 'Código dactilar (10 dígitos)' },
                ]}
              />
              {formData.tipoValidacion === 'fecha' ? (
                <Input
                  label="Fecha de expedición"
                  name="fechaExpedicion"
                  type="text"
                  value={formData.fechaExpedicion}
                  onChange={handleChange}
                  placeholder="DD/MM/AAAA"
                  error={errors.fechaExpedicion}
                  required
                />
              ) : (
                <Input
                  label="Código dactilar"
                  name="codigoDactilar"
                  type="text"
                  value={formData.codigoDactilar}
                  onChange={handleChange}
                  placeholder="Ingrese código dactilar (10 dígitos)"
                  inputMode="numeric"
                  maxLength={10}
                  error={errors.codigoDactilar}
                  required
                />
              )}
              <Button type="submit" loading={loading} disabled={bloqueado} className="w-full">
                Verificar Identidad
              </Button>
            </form>
          ) : (
            <form onSubmit={handleRegistro} className="space-y-5">
              {identidadData && (
                <div className="rounded-lg border border-emerald-200 bg-emerald-50 p-4">
                  <p className="mb-2 flex items-center gap-2 text-sm font-semibold text-emerald-900">
                    <svg className="h-4 w-4 text-emerald-600" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
                      <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                    </svg>
                    Identidad verificada en el Registro Civil
                  </p>
                  <dl className="grid grid-cols-1 gap-x-4 gap-y-1 text-sm sm:grid-cols-[auto_1fr]">
                    <dt className="text-emerald-700">Nombre</dt>
                    <dd className="font-medium text-emerald-950">{identidadData.nombres} {identidadData.apellidos}</dd>
                    <dt className="text-emerald-700">Cédula</dt>
                    <dd className="font-medium text-emerald-950">{identidadData.cedula}</dd>
                    <dt className="text-emerald-700">Nacionalidad</dt>
                    <dd className="font-medium text-emerald-950">{identidadData.nacionalidad}</dd>
                  </dl>
                </div>
              )}
              <Input
                label="Correo Electrónico"
                name="correo"
                type="email"
                value={formData.correo}
                onChange={handleChange}
                placeholder="Ingrese su correo"
                error={errors.correo}
                required
              />
              <Input
                label="Contraseña"
                name="password"
                type="password"
                value={formData.password}
                onChange={handleChange}
                placeholder="Mínimo 8 caracteres"
                error={errors.password}
                required
              />
              <Input
                label="Confirmar Contraseña"
                name="confirmPassword"
                type="password"
                value={formData.confirmPassword}
                onChange={handleChange}
                placeholder="Confirme su contraseña"
                error={errors.confirmPassword}
                required
              />
              <Button type="submit" loading={loading} disabled={bloqueado} className="w-full">
                Registrarse
              </Button>
              <Button type="button" onClick={() => setPaso(1)} className="w-full" variant="secondary">
                Atrás
              </Button>
            </form>
          )}
          <div className="mt-6 text-center">
            <span className="text-sm text-slate-600">
              ¿Ya tiene cuenta?{' '}
              <Link to="/login" className="rounded font-medium text-blue-700 hover:text-blue-900 hover:underline">
                Inicie Sesión
              </Link>
            </span>
          </div>
        </Card>
    </AuthLayout>
  );
};

export default RegistroPage;
