import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import Button from '../components/Button';
import Input from '../components/Input';
import Select from '../components/Select';
import Alert from '../components/Alert';
import Card from '../components/Card';

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
  const navigate = useNavigate();

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
    if (!/^\d{10}$/.test(codigo)) return 'El código dactilar debe tener 10 dígitos';
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
      setError(result.message);
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
      correo: formData.correo,
      password: formData.password,
      confirm_password: formData.confirmPassword,
    });
    setLoading(false);

    if (result.success) {
      navigate('/formulario');
    } else {
      setError(result.message);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900">BIESS</h1>
          <p className="text-gray-600 mt-2">Sistema de Registro y Postulación</p>
        </div>
        <Card title={paso === 1 ? 'Verificar Identidad' : 'Registro de Usuario'}>
          {error && <Alert type="error" message={error} />}
          
          {paso === 1 ? (
            <form onSubmit={handleVerificarIdentidad} className="space-y-6">
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
                  error={errors.codigoDactilar}
                  required
                />
              )}
              <Button type="submit" loading={loading} className="w-full">
                Verificar Identidad
              </Button>
            </form>
          ) : (
            <form onSubmit={handleRegistro} className="space-y-6">
              {identidadData && (
                <div className="bg-green-50 border border-green-200 rounded-lg p-4 mb-4">
                  <p className="text-sm text-green-800 font-medium">Identidad verificada:</p>
                  <p className="text-sm text-green-700"><strong>Nombre:</strong> {identidadData.nombres} {identidadData.apellidos}</p>
                  <p className="text-sm text-green-700"><strong>Cédula:</strong> {identidadData.cedula}</p>
                  <p className="text-sm text-green-700"><strong>Nacionalidad:</strong> {identidadData.nacionalidad}</p>
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
              <Button type="submit" loading={loading} className="w-full">
                Registrarse
              </Button>
              <Button type="button" onClick={() => setPaso(1)} className="w-full" variant="secondary">
                Atrás
              </Button>
            </form>
          )}
          <div className="mt-6 text-center">
            <span className="text-sm text-gray-600">
              ¿Ya tiene cuenta?{' '}
              <Link to="/login" className="text-blue-600 hover:text-blue-500 font-medium">
                Inicie Sesión
              </Link>
            </span>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default RegistroPage;
