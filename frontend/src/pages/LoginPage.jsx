import { useState } from 'react';
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
      navigate('/formulario');
    } else {
      setError(result.message);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full">
        <Card title="Iniciar Sesión">
          <form onSubmit={handleSubmit} className="space-y-6">
            {error && <Alert type="error" message={error} />}
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
            <Input
              label="Contraseña"
              name="password"
              type="password"
              value={formData.password}
              onChange={handleChange}
              placeholder="Ingrese su contraseña"
              error={errors.password}
              required
            />
            <Button type="submit" loading={loading} className="w-full">
              Iniciar Sesión
            </Button>
          </form>
          <div className="mt-6 text-center">
            <Link to="/recuperar" className="text-sm text-blue-600 hover:text-blue-500">
              ¿Olvidó su contraseña?
            </Link>
          </div>
          <div className="mt-4 text-center">
            <span className="text-sm text-gray-600">
              ¿No tiene cuenta?{' '}
              <Link to="/registro" className="text-blue-600 hover:text-blue-500">
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
