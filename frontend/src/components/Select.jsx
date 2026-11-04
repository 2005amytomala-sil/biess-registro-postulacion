const Select = ({ label, options = [], error = '', className = '', loading = false, disabled = false, ...props }) => {
  const baseClasses = 'w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 transition-colors bg-white';
  const errorClasses = error ? 'border-red-500 focus:ring-red-500' : 'border-gray-300';
  const disabledClasses = disabled || loading ? 'opacity-50 cursor-not-allowed' : '';
  const classes = `${baseClasses} ${errorClasses} ${disabledClasses} ${className}`;

  return (
    <div className="w-full">
      {label && (
        <label className="block text-sm font-medium text-gray-700 mb-1">
          {label}
        </label>
      )}
      <select className={classes} disabled={disabled || loading} {...props}>
        {loading ? (
          <option value="">Cargando...</option>
        ) : (
          options.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))
        )}
      </select>
      {error && (
        <p className="mt-1 text-sm text-red-600">{error}</p>
      )}
    </div>
  );
};

export default Select;
