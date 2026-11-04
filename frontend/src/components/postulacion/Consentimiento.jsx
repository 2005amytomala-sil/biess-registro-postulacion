import { useState } from 'react';
import Card from '../Card';
import Button from '../Button';
import { PASO_POR_SECCION } from '../../utils/completitud';

const ListaFaltantes = ({ titulo, faltantes, onIrAPaso }) => (
  <div className="bg-yellow-50 border border-yellow-300 rounded-lg p-4">
    <p className="text-sm font-medium text-yellow-900 mb-2">{titulo}</p>
    <ul className="text-sm text-yellow-800 space-y-1 list-disc list-inside">
      {faltantes.map((f) => (
        <li key={`${f.seccion}-${f.mensaje}`}>
          {PASO_POR_SECCION[f.seccion] !== undefined && f.seccion !== 'consentimiento' ? (
            <button
              type="button"
              onClick={() => onIrAPaso(PASO_POR_SECCION[f.seccion])}
              className="underline hover:text-yellow-900"
            >
              {f.mensaje}
            </button>
          ) : (
            f.mensaje
          )}
        </li>
      ))}
    </ul>
  </div>
);

const Consentimiento = ({ postulacion, faltantes = [], faltantesServidor = [], onEnviar, enviando, onIrAPaso }) => {
  const [aceptaTerminos, setAceptaTerminos] = useState(false);
  const enviado = postulacion?.estado === 'enviada';
  const codigoUnico = postulacion?.codigo;
  const completa = faltantes.length === 0;

  const handleEnviar = () => {
    if (aceptaTerminos && completa) onEnviar(aceptaTerminos);
  };

  if (enviado) {
    return (
      <Card title="Postulación Enviada Exitosamente">
        <div className="space-y-6">
          <div className="bg-green-50 border border-green-200 rounded-lg p-6">
            <div className="flex items-center justify-center mb-4">
              <svg className="w-16 h-16 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <h3 className="text-xl font-semibold text-green-900 text-center mb-2">
              ¡Su postulación ha sido enviada exitosamente!
            </h3>
            <p className="text-green-800 text-center mb-4">
              Gracias por participar en el concurso público de méritos del BIESS.
            </p>
            <div className="bg-white border border-green-300 rounded-lg p-4 text-center">
              <p className="text-sm text-slate-600 mb-1">Su código de postulación es:</p>
              <p className="text-2xl font-bold text-green-700">{codigoUnico}</p>
              <p className="text-xs text-slate-500 mt-2">
                Guarde este código para futuras consultas sobre su postulación.
              </p>
            </div>
          </div>

          <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <h4 className="text-sm font-medium text-blue-900 mb-2">Próximos pasos:</h4>
            <ul className="text-sm text-blue-800 space-y-1 list-disc list-inside">
              <li>Recibirá notificaciones sobre el proceso de evaluación</li>
              <li>Puede consultar el estado de su postulación con su código</li>
              <li>Manténgase atento a las comunicaciones oficiales del BIESS</li>
            </ul>
          </div>

          <div className="flex justify-center">
            <Button
              type="button"
              onClick={() => (window.location.href = '/dashboard')}
              variant="secondary"
            >
              Volver al Inicio
            </Button>
          </div>
        </div>
      </Card>
    );
  }

  return (
    <Card title="1.10 Consentimiento y Envío">
      <div className="space-y-6">
        {faltantesServidor.length > 0 && (
          <ListaFaltantes
            titulo="El servidor rechazó el envío. Complete lo siguiente:"
            faltantes={faltantesServidor}
            onIrAPaso={onIrAPaso}
          />
        )}
        {!completa && (
          <ListaFaltantes
            titulo="Para enviar la postulación primero debe completar:"
            faltantes={faltantes}
            onIrAPaso={onIrAPaso}
          />
        )}

        {/* Texto de consentimiento */}
        <div className="bg-slate-50 border border-slate-200 rounded-lg p-6 space-y-4">
          <h3 className="text-lg font-semibold text-slate-900">Declaración y Consentimiento</h3>
          
          <div className="space-y-3 text-sm text-slate-700">
            <p>
              <strong>Declaro bajo juramento</strong> que la información proporcionada en este formulario 
              de postulación es veraz, completa y correcta en todos sus aspectos.
            </p>
            
            <p>
              Autorizo al BIESS a verificar la veracidad de los datos aquí consignados, incluyendo pero 
              no limitándose a: consultas al Registro Civil, IESS, SENESCYT, y otras entidades pertinentes.
            </p>
            
            <p>
              Me comprometo a presentar la documentación original que respalde la información declarada 
              en caso de ser requerida durante el proceso de selección.
            </p>
            
            <p>
              Entiendo que la presentación de información falsa o inexactura en este formulario puede 
              causar mi descalificación del proceso de selección y dar lugar a las acciones legales 
              correspondientes.
            </p>
            
            <p>
              Acepto que mis datos personales serán procesados de acuerdo con la Ley de Protección de 
              Datos Personales del Ecuador y normativas vigentes del BIESS.
            </p>
            
            <p>
              Conozco y acepto las bases del concurso público de méritos y oposición para el Directorio 
              del BIESS, así como los criterios de evaluación establecidos.
            </p>
          </div>

          <div className="border-t border-slate-300 pt-4 mt-4">
            <p className="text-sm text-slate-600 italic">
              Al hacer clic en "Enviar Postulación", confirmo que he leído, entendido y aceptado los 
              términos y condiciones establecidos en esta declaración.
            </p>
          </div>
        </div>

        {/* Checkbox de aceptación */}
        <div className="flex items-start">
          <input
            type="checkbox"
            id="aceptaTerminos"
            checked={aceptaTerminos}
            onChange={(e) => setAceptaTerminos(e.target.checked)}
            disabled={enviando}
            className="mt-1 mr-3 w-5 h-5 text-blue-600 border-slate-300 rounded focus:ring-blue-500"
          />
          <label htmlFor="aceptaTerminos" className="text-sm text-slate-700">
            <strong>He leído, entendido y acepto</strong> los términos y condiciones de la declaración 
            y autorizo el envío de mi postulación.
          </label>
        </div>

        {/* Botones de acción */}
        <div className="flex gap-4 justify-end">
          <Button
            type="button"
            variant="secondary"
            onClick={() => onIrAPaso(6)}
            disabled={enviando}
          >
            Volver
          </Button>
          <Button
            type="button"
            onClick={handleEnviar}
            loading={enviando}
            disabled={!aceptaTerminos || !completa || enviando}
          >
            {enviando ? 'Enviando...' : 'Enviar Postulación'}
          </Button>
        </div>

        {/* Advertencia */}
        <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
          <p className="text-sm text-yellow-800">
            <strong>Importante:</strong> Una vez enviada la postulación, no podrá modificar los datos 
            proporcionados. Asegúrese de revisar toda la información antes de enviar.
          </p>
        </div>
      </div>
    </Card>
  );
};

export default Consentimiento;
