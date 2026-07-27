'use client';

import { RefreshCw, ArrowRightLeft } from 'lucide-react';
import { useSyncActivity } from '../sync-activity-context';

export function SyncActions() {
    const { running, nodes, environment, setEnvironment, reset } = useSyncActivity();
    const states = Object.values(nodes);
    const finished = !running && states.length > 0 && states.every(s => s === 'done' || s === 'error');

    return (
        <span className="flex items-center gap-2">
            <button
                onClick={() => setEnvironment('inventory')}
                disabled={running}
                title="Modo inventario: las tarjetas muestran el stock y su sincronizacion"
                className={`flex items-center gap-1.5 whitespace-nowrap rounded-lg border px-3 py-1.5 text-xs font-semibold transition-colors disabled:cursor-not-allowed disabled:opacity-50 ${
                    environment === 'inventory'
                        ? 'border-white bg-white text-[#0d5c80] shadow-sm'
                        : 'border-white/30 bg-white/15 text-white hover:bg-white/25'
                }`}
            >
                <RefreshCw size={13} />
                Sincronizar inventario
            </button>
            <button
                onClick={() => setEnvironment('products')}
                disabled={running}
                title="Modo productos: las tarjetas muestran la comparacion de catalogo"
                className={`flex items-center gap-1.5 whitespace-nowrap rounded-lg border px-3 py-1.5 text-xs font-semibold transition-colors disabled:cursor-not-allowed disabled:opacity-50 ${
                    environment === 'products'
                        ? 'border-white bg-white text-[#0d5c80] shadow-sm'
                        : 'border-white/30 bg-white/15 text-white hover:bg-white/25'
                }`}
            >
                <ArrowRightLeft size={13} />
                Comparar productos
            </button>
            {finished && (
                <button
                    onClick={reset}
                    className="whitespace-nowrap rounded-lg border border-white/20 px-2.5 py-1.5 text-xs font-semibold text-white/80 transition-colors hover:bg-white/15"
                >
                    Reiniciar
                </button>
            )}
        </span>
    );
}
