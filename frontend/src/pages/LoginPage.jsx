import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import Button from '../components/Button';
import Input from '../components/Input';
import Alert from '../components/Alert';
import Card from '../components/Card';

const LoginPage = () => {
  const [formData, setFormData] = useState({
    cedula: '',
    password: '',
  });
  const [errors, setErrors] = useState({});
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [bloqueado, setBloqueado] = useState(false);
  const [retryAfter, setRetryAfter] = useState(0);
  const { login } = useAuth();
  const navigate = useNavigate();

  const validateCedula = (cedula) => {
    if (!cedula) return 'La cédula es requerida';
    if (!/^\d{10}$/.test(cedula)) return 'La cédula debe tener 10 dígitos';
    return '';
  };

  const validatePassword = (password) => {
    if (!password) return 'La contraseña es requerida';
    return '';
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });

    let error = '';
    if (name === 'cedula') {
      error = validateCedula(value);
    } else if (name === 'password') {
      error = validatePassword(value);
    }
    setErrors({ ...errors, [name]: error });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    const cedulaError = validateCedula(formData.cedula);
    const passwordError = validatePassword(formData.password);

    if (cedulaError || passwordError) {
      setErrors({ cedula: cedulaError, password: passwordError });
      return;
    }

    setLoading(true);
    const result = await login(formData.cedula, formData.password);
    setLoading(false);

    if (result.success) {
      navigate('/dashboard');
    } else {
      if (result.retryAfter) {
        setBloqueado(true);
        setRetryAfter(result.retryAfter);
        setError('Registro fallido, inténtelo nuevamente en 10 minutos');
      } else {
        setError(result.message);
      }
    }
  };

  useEffect(() => {
    if (bloqueado && retryAfter > 0) {
      const timer = setInterval(() => {
        setRetryAfter((prev) => {
          if (prev <= 1) {
            clearInterval(timer);
            setBloqueado(false);
            setError('');
            return 0;
          }
          return prev - 1;
        });
      }, 1000);

      return () => clearInterval(timer);
    }
  }, [bloqueado, retryAfter]);

  const formatTime = (seconds) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900">BIESS</h1>
          <p className="text-gray-600 mt-2">Sistema de Registro y Postulación</p>
        </div>
        <Card title="Iniciar Sesión">
          <form onSubmit={handleSubmit} className="space-y-6">
            {error && <Alert type="error" message={error} />}
            {bloqueado && (
              <div className="bg-red-50 border border-red-200 rounded-lg p-4">
                <div className="flex items-center">
                  <svg className="w-5 h-5 text-red-500 mr-2" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-9a1 1 0 10-2 0v4a1 1 0 102 0v-4zm0 6a1 1 0 10-2 0 1 1 0 102 0z" clipRule="evenodd" />
                  </svg>
                  <p className="text-sm text-red-800">
                    Cuenta bloqueada. Intente nuevamente en: <strong className="text-red-900">{formatTime(retryAfter)}</strong>
                  </p>
                </div>
              </div>
            )}
            <Input
              label="Cédula"
              name="cedula"
              type="text"
              value={formData.cedula}
              onChange={handleChange}
              placeholder="Ingrese su cédula (10 dígitos)"
              error={errors.cedula}
              disabled={bloqueado}
              required
            />
            <Input
              label="Contraseña"
              name="password"
              type="password"
              value={formData.password}
              onChange={handleChange}
              placeholder="Ingrese su contraseña"
              error={errors.password}
              disabled={bloqueado}
              required
            />
            <Button type="submit" loading={loading} disabled={bloqueado} className="w-full">
              Iniciar Sesión
            </Button>
          </form>
          <div className="mt-6 flex justify-center">
            <Link to="/recuperar" className="text-sm text-blue-600 hover:text-blue-500 font-medium">
              ¿Olvidó su contraseña?
            </Link>
          </div>
          <div className="mt-4 text-center">
            <span className="text-sm text-gray-600">
              ¿No tiene cuenta?{' '}
              <Link to="/registro" className="text-blue-600 hover:text-blue-500 font-medium">
                Regístrese
              </Link>
            </span>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default LoginPage;
