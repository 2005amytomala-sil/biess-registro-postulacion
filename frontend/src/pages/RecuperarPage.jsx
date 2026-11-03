import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import Button from '../components/Button';
import Input from '../components/Input';
import Select from '../components/Select';
import Alert from '../components/Alert';
import Card from '../components/Card';

const RecuperarPage = () => {
  const [paso, setPaso] = useState(1);
  const [cedulaToken, setCedulaToken] = useState('');
  const [formData, setFormData] = useState({
    cedula: '',
    tipoValidacion: 'fecha',
    fechaExpedicion: '',
    codigoDactilar: '',
    token: '',
    password: '',
    confirmPassword: '',
  });
  const [errors, setErrors] = useState({});
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false);
  const { recuperarPassword, restablecerPassword } = useAuth();
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

  const validatePassword = (password) => {
    if (!password) return 'La contraseña es requerida';
    if (password.length < 8) return 'La contraseña debe tener al menos 8 caracteres';
    return '';
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });

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

  const handleRecuperar = async (e) => {
    e.preventDefault();
    setError('');
    setMessage('');

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
    const result = await recuperarPassword(
      formData.cedula,
      formData.tipoValidacion === 'fecha' ? formData.fechaExpedicion : '',
      formData.tipoValidacion === 'codigo' ? formData.codigoDactilar : ''
    );
    setLoading(false);

    if (result.success) {
      setCedulaToken(formData.cedula);
      setMessage('Se ha enviado un token de recuperación. Revise la consola del servidor para obtenerlo.');
      setPaso(2);
    } else {
      setError(result.message);
    }
  };

  const handleRestablecer = async (e) => {
    e.preventDefault();
    setError('');
    setMessage('');

    const passwordError = validatePassword(formData.password);
    const confirmError = formData.password !== formData.confirmPassword ? 'Las contraseñas no coinciden' : '';

    if (passwordError || confirmError) {
      setErrors({
        ...errors,
        password: passwordError,
        confirmPassword: confirmError,
      });
      return;
    }

    setLoading(true);
    const result = await restablecerPassword(
      cedulaToken,
      formData.token,
      formData.password,
      formData.confirmPassword
    );
    setLoading(false);

    if (result.success) {
      setMessage('Contraseña restablecida exitosamente');
      setTimeout(() => navigate('/login'), 2000);
    } else {
      setError(result.message);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full">
        <Card title={paso === 1 ? 'Recuperar Contraseña' : 'Restablecer Contraseña'}>
          {error && <Alert type="error" message={error} />}
          {message && <Alert type="success" message={message} />}
          
          {paso === 1 ? (
            <form onSubmit={handleRecuperar} className="space-y-6">
              <p className="text-sm text-gray-600">
                Ingrese su cédula y el método de validación para recuperar su contraseña.
              </p>
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
                Enviar Token
              </Button>
            </form>
          ) : (
            <form onSubmit={handleRestablecer} className="space-y-6">
              <p className="text-sm text-gray-600">
                Ingrese el token de recuperación y su nueva contraseña.
              </p>
              <Input
                label="Token de recuperación"
                name="token"
                type="text"
                value={formData.token}
                onChange={handleChange}
                placeholder="Ingrese el token"
                required
              />
              <Input
                label="Nueva contraseña"
                name="password"
                type="password"
                value={formData.password}
                onChange={handleChange}
                placeholder="Mínimo 8 caracteres"
                error={errors.password}
                required
              />
              <Input
                label="Confirmar contraseña"
                name="confirmPassword"
                type="password"
                value={formData.confirmPassword}
                onChange={handleChange}
                placeholder="Confirme su contraseña"
                error={errors.confirmPassword}
                required
              />
              <Button type="submit" loading={loading} className="w-full">
                Restablecer Contraseña
              </Button>
              <Button
                type="button"
                onClick={() => setPaso(1)}
                className="w-full"
                variant="secondary"
              >
                Atrás
              </Button>
            </form>
          )}

          <div className="mt-6 text-center">
            <Link to="/login" className="text-sm text-blue-600 hover:text-blue-500">
              Volver a Iniciar Sesión
            </Link>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default RecuperarPage;
