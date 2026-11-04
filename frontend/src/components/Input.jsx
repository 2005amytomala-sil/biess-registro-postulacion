import { useId } from 'react';
import FieldLabel, { FieldError, fieldClasses } from './FieldLabel';

const Input = ({ label, type = 'text', error = '', className = '', ...props }) => {
  const autoId = useId();
  const id = props.id || autoId;
  const errorId = `${id}-error`;

  return (
    <div className="w-full">
      {label && <FieldLabel htmlFor={id} label={label} required={props.required} />}
      <input
        type={type}
        className={`${fieldClasses(error)} ${className}`}
        aria-invalid={error ? true : undefined}
        aria-describedby={error ? errorId : undefined}
        {...props}
        id={id}
      />
      <FieldError id={errorId} error={error} />
    </div>
  );
};

export default Input;
