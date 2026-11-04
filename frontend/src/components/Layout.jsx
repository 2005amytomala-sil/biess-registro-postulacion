import Navbar from './Navbar';

const Layout = ({ children }) => {
  return (
    <div className="flex min-h-screen flex-col bg-slate-50">
      <Navbar />
      <main className="flex-1">{children}</main>
      <footer className="border-t border-slate-200 bg-white">
        <div className="mx-auto max-w-7xl px-4 py-4 text-center text-xs text-slate-500 sm:px-6 lg:px-8">
          Banco del Instituto Ecuatoriano de Seguridad Social · Concurso Público de Méritos y Oposición
        </div>
      </footer>
    </div>
  );
};

export default Layout;
