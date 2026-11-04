import { useCallback, useEffect, useState } from 'react';

// Cuenta regresiva basada en una hora límite (no se desfasa si la pestaña se suspende)
const useCuentaRegresiva = () => {
  const [hasta, setHasta] = useState(0);
  const [restante, setRestante] = useState(0);

  useEffect(() => {
    if (!hasta) return undefined;
    const id = setInterval(() => {
      const segundos = Math.max(0, Math.ceil((hasta - Date.now()) / 1000));
      setRestante(segundos);
      if (segundos === 0) setHasta(0);
    }, 1000);
    return () => clearInterval(id);
  }, [hasta]);

  const iniciar = useCallback((segundos) => {
    const s = Math.max(0, Math.ceil(Number(segundos) || 0));
    setRestante(s);
    setHasta(s > 0 ? Date.now() + s * 1000 : 0);
  }, []);

  return { bloqueado: restante > 0, restante, iniciar };
};

export const formatearTiempo = (segundos) => {
  const mins = Math.floor(segundos / 60);
  const secs = segundos % 60;
  return `${mins}:${secs.toString().padStart(2, '0')}`;
};

export default useCuentaRegresiva;
