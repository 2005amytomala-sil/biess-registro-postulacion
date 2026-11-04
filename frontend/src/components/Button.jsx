import Spinner from './Spinner';

const sizes = {
  sm: 'px-3 py-1.5 text-sm gap-1.5',
  md: 'px-4 py-2 text-sm gap-2',
  lg: 'px-6 py-3 text-base gap-2',
};

const variants = {
  primary: 'bg-blue-600 text-white shadow-sm hover:bg-blue-700 active:bg-blue-800 focus-visible:ring-blue-500',
  secondary:
    'bg-white text-slate-700 border border-slate-300 shadow-sm hover:bg-slate-50 active:bg-slate-100 focus-visible:ring-slate-400',
  danger: 'bg-red-600 text-white shadow-sm hover:bg-red-700 active:bg-red-800 focus-visible:ring-red-500',
  outline: 'border-2 border-blue-600 text-blue-700 hover:bg-blue-50 focus-visible:ring-blue-500',
  success: 'bg-emerald-600 text-white shadow-sm hover:bg-emerald-700 active:bg-emerald-800 focus-visible:ring-emerald-500',
};

const Button = ({ children, type = 'button', variant = 'primary', loading = false, disabled = false, size = 'md', className = '', onClick }) => {
  const classes = [
    'inline-flex items-center justify-center rounded-lg font-medium transition-colors',
    'focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-2',
    'disabled:opacity-60 disabled:cursor-not-allowed disabled:shadow-none',
    sizes[size] || sizes.md,
    variants[variant] || variants.primary,
    className,
  ].join(' ');

  return (
    <button type={type} className={classes} disabled={disabled || loading} onClick={onClick} aria-busy={loading || undefined}>
      {loading ? (
        <>
          <Spinner size="sm" />
          <span>Cargando...</span>
        </>
      ) : (
        children
      )}
    </button>
  );
};

export default Button;
