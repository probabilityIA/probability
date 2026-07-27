'use client';

import { RefreshCw, ArrowRightLeft, ReceiptText, type LucideIcon } from 'lucide-react';
import { useSyncActivity, type SyncEnvironment } from '../sync-activity-context';

interface EnvironmentAction {
    key: SyncEnvironment;
    label: string;
    icon: LucideIcon;
    hint: string;
    disabled?: boolean;
}

const ACTIONS: EnvironmentAction[] = [
    {
        key: 'inventory',
        label: 'Sincronizar inventario',
        icon: RefreshCw,
        hint: 'Las tarjetas muestran el estado del inventario y su ultima sincronizacion',
    },
    {
        key: 'products',
        label: 'Comparar productos',
        icon: ArrowRightLeft,
        hint: 'Las tarjetas muestran la comparacion de catalogo con cada canal',
    },
    {
        key: 'invoicing',
        label: 'Facturar',
        icon: ReceiptText,
        hint: 'Facturacion desde el hub: proximamente',
        disabled: true,
    },
];

export function SyncActions() {
    const { running, nodes, environment, setEnvironment, reset } = useSyncActivity();
    const states = Object.values(nodes);
    const finished = !running && states.length > 0 && states.every(s => s === 'done' || s === 'error');

    return (
        <span className="flex items-center gap-2">
            {ACTIONS.map(action => {
                const Icon = action.icon;
                const active = environment === action.key;
                return (
                    <button
                        key={action.key}
                        onClick={() => setEnvironment(active ? null : action.key)}
                        disabled={running || action.disabled}
                        title={action.hint}
                        className={`flex items-center gap-1.5 whitespace-nowrap rounded-lg border px-3 py-1.5 text-xs font-semibold transition-colors disabled:cursor-not-allowed disabled:opacity-50 ${
                            active
                                ? 'border-white bg-white text-[#0d5c80] shadow-sm'
                                : 'border-white/30 bg-white/15 text-white hover:bg-white/25'
                        }`}
                    >
                        <Icon size={13} />
                        {action.label}
                    </button>
                );
            })}
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
