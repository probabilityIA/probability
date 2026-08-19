'use client';

import type { SiigoCatalogItem } from '@/services/integrations/invoicing/siigo/infra/actions';

interface CampoCatalogoProps {
  etiqueta: string;
  opciones?: SiigoCatalogItem[];
  valor: number | string;
  onChange: (valor: string) => void;
  disabled?: boolean;
  requerido?: boolean;
}

export function CampoCatalogo({
  etiqueta,
  opciones,
  valor,
  onChange,
  disabled,
  requerido,
}: CampoCatalogoProps) {
  const clases =
    'w-full px-3 py-2 border rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50';
  const faltante = requerido && !valor;
  const borde = faltante ? 'border-amber-400' : 'border-gray-300';

  const sinCatalogo = !opciones || opciones.length === 0;

  return (
    <div>
      <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">
        {etiqueta}
        {requerido ? <span className="text-amber-600"> *</span> : <span className="text-gray-400"> — opcional</span>}
      </label>

      {sinCatalogo ? (
        <input
          type="number"
          value={valor}
          onChange={(e) => onChange(e.target.value)}
          placeholder="ID"
          disabled={disabled}
          className={`${clases} ${borde}`}
        />
      ) : (
        <select
          value={valor}
          onChange={(e) => onChange(e.target.value)}
          disabled={disabled}
          className={`${clases} ${borde}`}
        >
          <option value="">Selecciona una opcion</option>
          {opciones.map((o) => (
            <option key={o.id} value={o.id}>
              {o.name}
              {o.percent ? ` (${o.percent}%)` : ''}
              {o.detail && !o.percent ? ` — ${o.detail}` : ''}
            </option>
          ))}
        </select>
      )}
    </div>
  );
}
