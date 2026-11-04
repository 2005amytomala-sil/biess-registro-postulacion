import { useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import Button from '../components/Button';
import Input from '../components/Input';
import Alert from '../components/Alert';
import Card from '../components/Card';
import AuthLayout from '../components/AuthLayout';
import CuentaRegresiva from '../components/CuentaRegresiva';
import useCuentaRegresiva from '../hooks/useCuentaRegresiva';

const LoginPage = () => {
  const [formData, setFormData] = useState({
    cedula: '',
    password: '',
  });
  const [errors, setErrors] = useState({});
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { bloqueado, restante, iniciar: iniciarBloqueo } = useCuentaRegresiva();
  const { login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [exito, setExito] = useState(location.state?.mensaje || '');

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
    setExito('');

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
      navigate('/postulacion', { replace: true });
    } else {
      setError(result.message);
      if (result.retryAfter) iniciarBloqueo(result.retryAfter);
    }
  };

  return (
    <AuthLayout>
        <Card title="Iniciar Sesión">
          <form onSubmit={handleSubmit} className="space-y-5">
            {exito && <Alert type="success" message={exito} />}
            {bloqueado ? (
              <CuentaRegresiva mensaje={error} restante={restante} />
            ) : (
              error && <Alert type="error" message={error} />
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
            <Link to="/recuperar" className="rounded text-sm font-medium text-blue-700 hover:text-blue-900 hover:underline">
              ¿Olvidó su contraseña?
            </Link>
          </div>
          <div className="mt-4 text-center">
            <span className="text-sm text-slate-600">
              ¿No tiene cuenta?{' '}
              <Link to="/registro" className="rounded font-medium text-blue-700 hover:text-blue-900 hover:underline">
                Regístrese
              </Link>
            </span>
          </div>
        </Card>
    </AuthLayout>
  );
};

export default LoginPage;
