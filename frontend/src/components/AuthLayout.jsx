// Contenedor visual de las pantallas de acceso (login, registro, recuperación)
const AuthLayout = ({ children }) => (
  <div className="flex min-h-[calc(100vh-8rem)] items-center justify-center bg-gradient-to-b from-blue-50 via-slate-50 to-slate-100 px-4 py-10 sm:px-6">
    <div className="w-full max-w-md">
      <header className="mb-6 text-center">
        <span className="mx-auto mb-3 flex h-12 w-12 items-center justify-center rounded-xl bg-blue-700 text-lg font-bold text-white shadow-md">
          B
        </span>
        <h1 className="text-xl font-bold text-slate-900 sm:text-2xl">BIESS – Concurso de Méritos y Oposición</h1>
        <p className="mt-1 text-sm text-slate-600">Sistema de Registro y Postulación al Directorio</p>
      </header>
      {children}
    </div>
  </div>
);

export default AuthLayout;
