import { useId } from 'react';
import FieldLabel, { FieldError, fieldClasses } from './FieldLabel';

const Select = ({ label, options = [], error = '', className = '', loading = false, disabled = false, ...props }) => {
  const autoId = useId();
  const id = props.id || autoId;
  const errorId = `${id}-error`;

  return (
    <div className="w-full">
      {label && <FieldLabel htmlFor={id} label={label} required={props.required} />}
      <select
        className={`${fieldClasses(error)} pr-8 ${loading ? 'animate-pulse' : ''} ${className}`}
        disabled={disabled || loading}
        aria-invalid={error ? true : undefined}
        aria-describedby={error ? errorId : undefined}
        aria-busy={loading || undefined}
        {...props}
        id={id}
      >
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
      <FieldError id={errorId} error={error} />
    </div>
  );
};

export default Select;
