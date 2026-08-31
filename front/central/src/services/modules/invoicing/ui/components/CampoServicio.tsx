'use client';

import { useEffect, useRef, useState } from 'react';
import { searchSiigoProductsAction, type SiigoCatalogItem } from '@/services/integrations/invoicing/siigo/infra/actions';

interface CampoServicioProps {
  etiqueta: string;
  valor: string;
  onChange: (valor: string) => void;
  disabled?: boolean;
  placeholder?: string;
  ayuda?: string;
  advertencia?: string;
  integrationId?: number;
  buscable?: boolean;
}

export function CampoServicio({
  etiqueta,
  valor,
  onChange,
  disabled,
  placeholder,
  ayuda,
  advertencia,
  integrationId,
  buscable,
}: CampoServicioProps) {
  const [termino, setTermino] = useState('');
  const [resultados, setResultados] = useState<SiigoCatalogItem[] | null>(null);
  const [buscando, setBuscando] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [abierto, setAbierto] = useState(false);
  const contenedor = useRef<HTMLDivElement>(null);

  const clases =
    'w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] disabled:opacity-50';

  useEffect(() => {
    function fuera(e: MouseEvent) {
      if (contenedor.current && !contenedor.current.contains(e.target as Node)) setAbierto(false);
    }
    document.addEventListener('mousedown', fuera);
    return () => document.removeEventListener('mousedown', fuera);
  }, []);

  useEffect(() => {
    if (!buscable || !integrationId || !abierto || valor) return;
    const limpio = termino.trim();

    let cancelado = false;
    setBuscando(true);
    setError(null);

    const id = setTimeout(() => {
      searchSiigoProductsAction(integrationId, limpio)
        .then((r) => {
          if (cancelado) return;
          if (!r.success) {
            setError(r.message ?? 'No se pudo buscar en el facturador');
            setResultados([]);
            return;
          }
          setResultados(r.data ?? []);
        })
        .catch(() => {
          if (!cancelado) setError('No se pudo buscar en el facturador');
        })
        .finally(() => {
          if (!cancelado) setBuscando(false);
        });
    }, 450);

    return () => {
      cancelado = true;
      clearTimeout(id);
    };
  }, [termino, buscable, integrationId, abierto, valor]);

  if (!buscable || !integrationId) {
    return (
      <div>
        <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">{etiqueta}</label>
        <input
          type="text"
          value={valor}
          onChange={(e) => onChange(e.target.value)}
          placeholder={placeholder}
          disabled={disabled}
          className={clases}
        />
        {ayuda && <p className="text-xs text-gray-400 mt-1">{ayuda}</p>}
        {advertencia && !valor && (
          <p className="mt-1 text-[11px] leading-snug text-amber-600 dark:text-amber-400">{advertencia}</p>
        )}
      </div>
    );
  }

  return (
    <div ref={contenedor}>
      <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">{etiqueta}</label>

      {valor ? (
        <div className="flex items-center gap-2 rounded-md border border-gray-300 bg-white dark:bg-gray-800 px-3 py-2">
          <span className="font-mono text-sm text-gray-900 dark:text-white">{valor}</span>
          <button
            type="button"
            onClick={() => {
              onChange('');
              setTermino('');
              setAbierto(true);
            }}
            disabled={disabled}
            className="ml-auto text-xs font-medium hover:underline disabled:opacity-50"
            style={{ color: 'var(--color-primary)' }}
          >
            Cambiar
          </button>
        </div>
      ) : (
        <input
          type="text"
          value={termino}
          onChange={(e) => {
            setTermino(e.target.value);
            setAbierto(true);
          }}
          onFocus={() => setAbierto(true)}
          placeholder="Escribe el nombre o el código para buscarlo"
          disabled={disabled}
          className={clases}
        />
      )}

      {abierto && !valor && (
        <div className="mt-1 max-h-56 overflow-y-auto rounded-md border border-gray-200 bg-white dark:bg-gray-800 dark:border-gray-700">
          {buscando ? (
            <p className="px-3 py-2 text-[11px] text-gray-400">Buscando en el facturador...</p>
          ) : error ? (
            <p className="px-3 py-2 text-[11px] text-amber-600">{error}</p>
          ) : resultados && resultados.length === 0 ? (
            <p className="px-3 py-2 text-[11px] text-gray-400">
              {termino.trim()
                ? 'Nada con ese nombre o código. Si aún no lo creaste en el facturador, créalo y vuelve a buscar.'
                : 'El facturador no devolvió cuentas.'}
            </p>
          ) : (
            <>
            {!termino.trim() && (
              <p className="border-b border-gray-100 px-3 py-1.5 text-[10px] text-gray-400 dark:border-gray-700">
                Primeras cuentas del facturador. Escribe para buscar en todo el catalogo.
              </p>
            )}
            {(resultados ?? []).map((o) => (
              <button
                key={o.code}
                type="button"
                onClick={() => {
                  onChange(o.code ?? '');
                  setAbierto(false);
                }}
                className="flex w-full items-center gap-2 px-3 py-2 text-left text-xs hover:bg-[color-mix(in_srgb,var(--color-primary)_8%,white)] dark:hover:bg-gray-700"
              >
                <span className="font-mono text-gray-900 dark:text-white">{o.code}</span>
                <span className="truncate text-gray-600 dark:text-gray-300">{o.name}</span>
                {o.detail && (
                  <span
                    className={`ml-auto shrink-0 rounded px-1.5 py-0.5 text-[10px] ${
                      o.detail === 'Service' ? 'bg-[color-mix(in_srgb,var(--color-primary)_15%,white)] text-[var(--color-primary)]' : 'bg-gray-100 text-gray-500'
                    }`}
                  >
                    {o.detail === 'Service' ? 'servicio' : o.detail}
                  </span>
                )}
              </button>
            ))}
            </>
          )}
        </div>
      )}

      {ayuda && <p className="text-xs text-gray-400 mt-1">{ayuda}</p>}
      {advertencia && !valor && (
        <p className="mt-1 text-[11px] leading-snug text-amber-600 dark:text-amber-400">{advertencia}</p>
      )}
    </div>
  );
}
