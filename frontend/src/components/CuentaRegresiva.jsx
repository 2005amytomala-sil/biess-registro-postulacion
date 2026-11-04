import { formatearTiempo } from '../hooks/useCuentaRegresiva';

const CuentaRegresiva = ({ mensaje, restante }) => (
  <div className="rounded-lg border border-red-200 bg-red-50 p-4 shadow-sm" role="alert">
    <div className="flex items-start gap-3">
      <svg className="mt-0.5 h-5 w-5 flex-shrink-0 text-red-600" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
        <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clipRule="evenodd" />
      </svg>
      <div className="flex-1 text-sm text-red-900">
        {mensaje && <p className="font-medium">{mensaje}</p>}
        <p className="mt-1 flex flex-wrap items-baseline gap-x-2">
          <span>Intente nuevamente en:</span>
          <strong className="font-mono text-lg tabular-nums text-red-700" aria-live="polite">
            {formatearTiempo(restante)}
          </strong>
        </p>
      </div>
    </div>
  </div>
);

export default CuentaRegresiva;
