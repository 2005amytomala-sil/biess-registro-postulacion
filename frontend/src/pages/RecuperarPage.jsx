import { useState } from 'react';
import { Link } from 'react-router-dom';
import Button from '../components/Button';
import Input from '../components/Input';
import Alert from '../components/Alert';
import Card from '../components/Card';

const RecuperarPage = () => {
  const [correo, setCorreo] = useState('');
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setMessage('');
    setLoading(true);

    // Simulación - implementar API real cuando esté disponible
    setTimeout(() => {
      setLoading(false);
      setMessage('Se ha enviado un correo de recuperación a su dirección.');
    }, 1000);
  };

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full">
        <Card title="Recuperar Contraseña">
          <form onSubmit={handleSubmit} className="space-y-6">
            {error && <Alert type="error" message={error} />}
            {message && <Alert type="success" message={message} />}
            <p className="text-sm text-gray-600">
              Ingrese su correo electrónico y le enviaremos instrucciones para
              recuperar su contraseña.
            </p>
            <Input
              label="Correo Electrónico"
              type="email"
              value={correo}
              onChange={(e) => setCorreo(e.target.value)}
              placeholder="Ingrese su correo"
              required
            />
            <Button type="submit" loading={loading} className="w-full">
              Enviar Instrucciones
            </Button>
          </form>
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
