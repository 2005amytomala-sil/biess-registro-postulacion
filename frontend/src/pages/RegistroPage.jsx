import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import Button from '../components/Button';
import Input from '../components/Input';
import Alert from '../components/Alert';
import Card from '../components/Card';

const RegistroPage = () => {
  const [formData, setFormData] = useState({
    cedula: '',
    correo: '',
    password: '',
    confirmPassword: '',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { register } = useAuth();
  const navigate = useNavigate();

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    if (formData.password !== formData.confirmPassword) {
      setError('Las contraseñas no coinciden');
      return;
    }

    setLoading(true);

    const result = await register({
      cedula: formData.cedula,
      correo: formData.correo,
      password: formData.password,
    });
    setLoading(false);

    if (result.success) {
      navigate('/login');
    } else {
      setError(result.message);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full">
        <Card title="Registro de Usuario">
          <form onSubmit={handleSubmit} className="space-y-6">
            {error && <Alert type="error" message={error} />}
            <Input
              label="Cédula"
              name="cedula"
              type="text"
              value={formData.cedula}
              onChange={handleChange}
              placeholder="Ingrese su cédula"
              required
            />
            <Input
              label="Correo Electrónico"
              name="correo"
              type="email"
              value={formData.correo}
              onChange={handleChange}
              placeholder="Ingrese su correo"
              required
            />
            <Input
              label="Contraseña"
              name="password"
              type="password"
              value={formData.password}
              onChange={handleChange}
              placeholder="Ingrese su contraseña"
              required
            />
            <Input
              label="Confirmar Contraseña"
              name="confirmPassword"
              type="password"
              value={formData.confirmPassword}
              onChange={handleChange}
              placeholder="Confirme su contraseña"
              required
            />
            <Button type="submit" loading={loading} className="w-full">
              Registrarse
            </Button>
          </form>
          <div className="mt-6 text-center">
            <span className="text-sm text-gray-600">
              ¿Ya tiene cuenta?{' '}
              <Link to="/login" className="text-blue-600 hover:text-blue-500">
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
