import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Layout from '../components/Layout';
import ProtectedRoute from './ProtectedRoute';
import { useAuth } from '../context/AuthContext';
import LoginPage from '../pages/LoginPage';
import RegistroPage from '../pages/RegistroPage';
import RecuperarPage from '../pages/RecuperarPage';
import DashboardPage from '../pages/DashboardPage';
import PostulacionPage from '../pages/PostulacionPage';

// Redirige según la sesión: con sesión a /dashboard, sin sesión a /login
const RedireccionPorSesion = () => {
  const { isAuthenticated, loading } = useAuth();
  if (loading) return null;
  return <Navigate to={isAuthenticated ? '/dashboard' : '/login'} replace />;
};

const AppRoutes = () => {
  return (
    <BrowserRouter>
      <Layout>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/registro" element={<RegistroPage />} />
          <Route path="/recuperar" element={<RecuperarPage />} />
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <DashboardPage />
              </ProtectedRoute>
            }
          />
          <Route
            path="/postulacion"
            element={
              <ProtectedRoute>
                <PostulacionPage />
              </ProtectedRoute>
            }
          />
          <Route path="*" element={<RedireccionPorSesion />} />
        </Routes>
      </Layout>
    </BrowserRouter>
  );
};

export default AppRoutes;
