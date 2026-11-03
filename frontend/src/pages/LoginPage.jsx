import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import Button from '../components/Button';
import Input from '../components/Input';
import Alert from '../components/Alert';
import Card from '../components/Card';

const LoginPage = () => {
  const [cedula, setCedula] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    const result = await login(cedula, password);
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
              type="text"
              value={cedula}
              onChange={(e) => setCedula(e.target.value)}
              placeholder="Ingrese su cédula"
              required
            />
            <Input
              label="Contraseña"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Ingrese su contraseña"
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
