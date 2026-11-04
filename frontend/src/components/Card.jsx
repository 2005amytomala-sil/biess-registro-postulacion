const Card = ({ children, className = '', title }) => {
  return (
    <section className={`rounded-xl border border-slate-200 bg-white p-5 shadow-sm sm:p-6 ${className}`}>
      {title && (
        <h2 className="mb-5 border-b border-slate-100 pb-3 text-lg font-semibold text-slate-900">{title}</h2>
      )}
      {children}
    </section>
  );
};

export default Card;
