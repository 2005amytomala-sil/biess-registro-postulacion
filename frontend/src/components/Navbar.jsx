import { useState } from 'react';
import { Link, NavLink } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const enlaceClases = ({ isActive }) =>
  `block rounded-md px-3 py-2 text-sm font-medium transition-colors ${
    isActive ? 'bg-blue-800 text-white' : 'text-blue-100 hover:bg-blue-600 hover:text-white'
  }`;

const Navbar = () => {
  const { usuario, isAuthenticated, logout } = useAuth();
  const [menuAbierto, setMenuAbierto] = useState(false);
  const cerrarMenu = () => setMenuAbierto(false);

  const enlaces = isAuthenticated
    ? [
        { to: '/dashboard', label: 'Dashboard' },
        { to: '/postulacion', label: 'Postulación' },
      ]
    : [
        { to: '/login', label: 'Iniciar Sesión' },
        { to: '/registro', label: 'Registrarse' },
      ];

  return (
    <nav className="bg-blue-700 text-white shadow-md" aria-label="Navegación principal">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between gap-4">
          <Link to="/" className="flex min-w-0 items-center gap-3" onClick={cerrarMenu}>
            <span className="flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-lg bg-white text-sm font-bold tracking-tight text-blue-700">
              B
            </span>
            <span className="min-w-0 leading-tight">
              <span className="block text-base font-semibold">BIESS</span>
              <span className="hidden truncate text-xs text-blue-200 sm:block">Concurso de Méritos y Oposición</span>
            </span>
          </Link>

          {/* Escritorio */}
          <div className="hidden items-center gap-1 md:flex">
            {enlaces.map((e) => (
              <NavLink key={e.to} to={e.to} className={enlaceClases}>
                {e.label}
              </NavLink>
            ))}
            {isAuthenticated && (
              <>
                <div className="mx-3 h-8 w-px bg-blue-500" aria-hidden="true" />
                <div className="mr-2 text-right leading-tight">
                  <p className="text-sm font-medium">{usuario?.cedula || 'Usuario'}</p>
                  {usuario?.correo && <p className="max-w-[14rem] truncate text-xs text-blue-200">{usuario.correo}</p>}
                </div>
                <button
                  type="button"
                  onClick={logout}
                  className="rounded-md border border-blue-400 px-3 py-1.5 text-sm font-medium text-white transition-colors hover:bg-blue-600"
                >
                  Cerrar Sesión
                </button>
              </>
            )}
          </div>

          {/* Botón de menú móvil */}
          <button
            type="button"
            className="inline-flex items-center justify-center rounded-md p-2 text-blue-100 hover:bg-blue-600 hover:text-white md:hidden"
            aria-expanded={menuAbierto}
            aria-controls="menu-movil"
            aria-label={menuAbierto ? 'Cerrar menú' : 'Abrir menú'}
            onClick={() => setMenuAbierto((v) => !v)}
          >
            <svg className="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
              {menuAbierto ? (
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              ) : (
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
              )}
            </svg>
          </button>
        </div>
      </div>

      {/* Menú móvil */}
      {menuAbierto && (
        <div id="menu-movil" className="border-t border-blue-600 md:hidden">
          <div className="space-y-1 px-4 py-3">
            {isAuthenticated && (
              <div className="mb-2 rounded-md bg-blue-800 px-3 py-2">
                <p className="text-sm font-medium">{usuario?.cedula || 'Usuario'}</p>
                {usuario?.correo && <p className="truncate text-xs text-blue-200">{usuario.correo}</p>}
              </div>
            )}
            {enlaces.map((e) => (
              <NavLink key={e.to} to={e.to} className={enlaceClases} onClick={cerrarMenu}>
                {e.label}
              </NavLink>
            ))}
            {isAuthenticated && (
              <button
                type="button"
                onClick={() => {
                  cerrarMenu();
                  logout();
                }}
                className="block w-full rounded-md px-3 py-2 text-left text-sm font-medium text-blue-100 hover:bg-blue-600 hover:text-white"
              >
                Cerrar Sesión
              </button>
            )}
          </div>
        </div>
      )}
    </nav>
  );
};

export default Navbar;
