const Stepper = ({ steps, currentStep, onStepClick }) => {
  return (
    <div className="mb-8">
      <div className="flex items-center justify-between">
        {steps.map((step, index) => {
          const isCompleted = index < currentStep;
          const isCurrent = index === currentStep;
          const isClickable = step.clickable !== false && onStepClick;

          return (
            <div key={step.id} className="flex items-center flex-1">
              <div className="flex items-center">
                <div
                  className={`
                    flex items-center justify-center w-10 h-10 rounded-full font-medium transition-colors
                    ${isCompleted ? 'bg-green-600 text-white' : ''}
                    ${isCurrent ? 'bg-blue-600 text-white' : ''}
                    ${!isCompleted && !isCurrent ? 'bg-gray-200 text-gray-600' : ''}
                    ${isClickable && !isCurrent ? 'cursor-pointer hover:bg-blue-500' : ''}
                  `}
                  onClick={() => isClickable && onStepClick(index)}
                >
                  {isCompleted ? (
                    <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                  ) : (
                    index + 1
                  )}
                </div>
                <div className="ml-3">
                  <p
                    className={`text-sm font-medium ${
                      isCurrent ? 'text-blue-600' : isCompleted ? 'text-green-600' : 'text-gray-600'
                    }`}
                  >
                    {step.title}
                  </p>
                  {step.subtitle && (
                    <p className="text-xs text-gray-500">{step.subtitle}</p>
                  )}
                </div>
              </div>
              {index < steps.length - 1 && (
                <div className="flex-1 mx-4 h-1 bg-gray-200 rounded">
                  <div
                    className={`h-1 rounded transition-all ${
                      isCompleted ? 'bg-green-600' : 'bg-gray-200'
                    }`}
                    style={{ width: isCompleted ? '100%' : '0%' }}
                  />
                </div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default Stepper;
