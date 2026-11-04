import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Layout from '../components/Layout';
import ProtectedRoute from './ProtectedRoute';
import LoginPage from '../pages/LoginPage';
import RegistroPage from '../pages/RegistroPage';
import RecuperarPage from '../pages/RecuperarPage';
import PostulacionPage from '../pages/PostulacionPage';

const AppRoutes = () => {
  return (
    <BrowserRouter>
      <Layout>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/registro" element={<RegistroPage />} />
          <Route path="/recuperar" element={<RecuperarPage />} />
          <Route
            path="/postulacion"
            element={
              <ProtectedRoute>
                <PostulacionPage />
              </ProtectedRoute>
            }
          />
          <Route path="/" element={<LoginPage />} />
        </Routes>
      </Layout>
    </BrowserRouter>
  );
};

export default AppRoutes;
