const CheckIcon = () => (
  <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" />
  </svg>
);

const Stepper = ({ steps, currentStep, onStepClick }) => {
  const actual = steps[currentStep];
  const progreso = steps.length > 1 ? (currentStep / (steps.length - 1)) * 100 : 100;

  return (
    <nav aria-label="Progreso de la postulación" className="mb-2">
      {/* Móvil: paso actual y barra de progreso */}
      <div className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm lg:hidden">
        <p className="text-xs font-medium uppercase tracking-wide text-blue-700">
          Paso {currentStep + 1} de {steps.length}
        </p>
        <p className="mt-0.5 text-base font-semibold text-slate-900">{actual?.title}</p>
        {actual?.subtitle && <p className="text-sm text-slate-500">{actual.subtitle}</p>}
        <div className="mt-3 h-2 overflow-hidden rounded-full bg-slate-200" aria-hidden="true">
          <div className="h-full rounded-full bg-blue-600 transition-all" style={{ width: `${progreso}%` }} />
        </div>
        <ol className="mt-3 flex flex-wrap gap-1.5">
          {steps.map((step, index) => {
            const isCompleted = index < currentStep;
            const isCurrent = index === currentStep;
            const isClickable = step.clickable !== false && onStepClick;
            return (
              <li key={step.id}>
                <button
                  type="button"
                  disabled={!isClickable}
                  onClick={() => isClickable && onStepClick(index)}
                  aria-current={isCurrent ? 'step' : undefined}
                  aria-label={`Paso ${index + 1}: ${step.title}`}
                  className={`flex h-8 w-8 items-center justify-center rounded-full text-xs font-semibold transition-colors ${
                    isCurrent
                      ? 'bg-blue-600 text-white ring-2 ring-blue-200'
                      : isCompleted
                        ? 'bg-emerald-600 text-white'
                        : 'bg-slate-200 text-slate-600 hover:bg-slate-300'
                  }`}
                >
                  {isCompleted ? <CheckIcon /> : index + 1}
                </button>
              </li>
            );
          })}
        </ol>
      </div>

      {/* Escritorio */}
      <ol className="hidden rounded-xl border border-slate-200 bg-white p-4 shadow-sm lg:flex">
        {steps.map((step, index) => {
          const isCompleted = index < currentStep;
          const isCurrent = index === currentStep;
          const isClickable = step.clickable !== false && onStepClick;

          return (
            <li key={step.id} className="flex flex-1 items-start">
              <button
                type="button"
                disabled={!isClickable}
                onClick={() => isClickable && onStepClick(index)}
                aria-current={isCurrent ? 'step' : undefined}
                className="group flex min-w-0 flex-col items-center gap-2 rounded-lg px-1 text-center disabled:cursor-default"
              >
                <span
                  className={`flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full text-sm font-semibold transition-colors ${
                    isCurrent
                      ? 'bg-blue-600 text-white ring-4 ring-blue-100'
                      : isCompleted
                        ? 'bg-emerald-600 text-white group-hover:bg-emerald-700'
                        : 'border-2 border-slate-300 bg-white text-slate-500 group-hover:border-blue-400 group-hover:text-blue-600'
                  }`}
                >
                  {isCompleted ? <CheckIcon /> : index + 1}
                </span>
                <span className="min-w-0">
                  <span
                    className={`block text-xs font-semibold leading-tight ${
                      isCurrent ? 'text-blue-700' : isCompleted ? 'text-emerald-700' : 'text-slate-600'
                    }`}
                  >
                    {step.title}
                  </span>
                  {step.subtitle && <span className="mt-0.5 block text-[11px] leading-tight text-slate-500">{step.subtitle}</span>}
                </span>
              </button>
              {index < steps.length - 1 && (
                <span
                  className={`mx-1 mt-5 h-0.5 flex-1 rounded ${isCompleted ? 'bg-emerald-500' : 'bg-slate-200'}`}
                  aria-hidden="true"
                />
              )}
            </li>
          );
        })}
      </ol>
    </nav>
  );
};

export default Stepper;
