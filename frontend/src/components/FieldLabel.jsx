// Etiqueta común para Input y Select: marca los campos obligatorios con un asterisco visible.
// Acepta etiquetas que ya traen " *" al final (se normalizan para no duplicarlo).
const FieldLabel = ({ htmlFor, label, required }) => {
  const texto = String(label).replace(/\s*\*\s*$/, '');
  const obligatorio = required || /\*\s*$/.test(String(label));

  return (
    <label htmlFor={htmlFor} className="mb-1.5 block text-sm font-medium text-slate-700">
      {texto}
      {obligatorio && (
        <span className="ml-0.5 text-red-600" aria-hidden="true">
          *
        </span>
      )}
    </label>
  );
};

export const fieldClasses = (error) =>
  [
    'block w-full rounded-lg border bg-white px-3 py-2 text-sm text-slate-900 shadow-sm transition-colors',
    'placeholder:text-slate-400',
    'focus:outline-none focus:ring-2 focus:ring-offset-0',
    'disabled:cursor-not-allowed disabled:bg-slate-100 disabled:text-slate-500',
    'read-only:bg-slate-50',
    error
      ? 'border-red-500 focus:border-red-500 focus:ring-red-200'
      : 'border-slate-300 hover:border-slate-400 focus:border-blue-500 focus:ring-blue-200',
  ].join(' ');

export const FieldError = ({ id, error }) =>
  error ? (
    <p id={id} className="mt-1.5 text-sm text-red-600" role="alert">
      {error}
    </p>
  ) : null;

export default FieldLabel;
