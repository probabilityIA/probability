'use client';

import { RefreshCw, ArrowRightLeft, ArrowDownToLine, ReceiptText, LayoutGrid, Share2, Table2, ClipboardList, type LucideIcon } from 'lucide-react';
import { useSyncActivity, type HubView, type SyncEnvironment } from '../sync-activity-context';

interface EnvironmentAction {
    key: SyncEnvironment;
    label: string;
    icon: LucideIcon;
    hint: string;
    disabled?: boolean;
}

const ACTIONS: EnvironmentAction[] = [
    {
        key: 'products',
        label: 'Comparar productos',
        icon: ArrowRightLeft,
        hint: 'Que producto esta en cada canal y cual falta por publicar',
    },
    {
        key: 'data',
        label: 'Actualizar productos',
        icon: ArrowDownToLine,
        hint: 'Que dato del canal (nombre, imagen, categoría) puede entrar a Probability',
    },
    {
        key: 'inventory',
        label: 'Sincronizar inventario',
        icon: RefreshCw,
        hint: 'Cuanto stock tiene cada canal y cual quedaría distinto al de Probability',
    },
    {
        key: 'orders_compare',
        label: 'Comparar órdenes',
        icon: ClipboardList,
        hint: 'Que orden existe en el canal y no en Probability, y crearla acá',
    },
    {
        key: 'invoicing',
        label: 'Facturar',
        icon: ReceiptText,
        hint: 'Facturación desde el hub: próximamente',
        disabled: true,
    },
];

const VIEWS: { key: HubView; label: string; icon: LucideIcon; hint: string }[] = [
    {
        key: 'diagrama',
        label: 'Diagrama',
        icon: Share2,
        hint: 'Ver tus canales conectados al núcleo, con el resumen de cada uno',
    },
    {
        key: 'informe',
        label: 'Informe',
        icon: Table2,
        hint: 'Ver lo mismo como tabla: una fila por producto y una columna por canal',
    },
];

export function SyncActions() {
    const { running, nodes, environment, setEnvironment, view, setView, reset } = useSyncActivity();
    const states = Object.values(nodes);
    const finished = !running && states.length > 0 && states.every(s => s === 'done' || s === 'error');

    return (
        <span className="flex items-center gap-2">
            <span className="mr-1 flex items-center gap-0.5 rounded-xl border border-white/25 bg-white/10 p-0.5">
                {VIEWS.map(item => {
                    const Icon = item.icon;
                    const active = view === item.key;
                    return (
                        <button
                            key={item.key}
                            onClick={() => setView(item.key)}
                            title={item.hint}
                            className={`flex items-center gap-1.5 whitespace-nowrap rounded-lg px-2.5 py-1.5 text-xs font-semibold transition-colors ${
                                active
                                    ? 'bg-white text-[#0d5c80] shadow-sm'
                                    : 'text-white/85 hover:bg-white/15 hover:text-white'
                            }`}
                        >
                            <Icon size={13} />
                            {item.label}
                        </button>
                    );
                })}
            </span>
            <button
                onClick={() => { reset(); setEnvironment(null); }}
                disabled={running}
                title="Resumen de órdenes: cuántas entraron por cada canal y como van"
                className={`flex items-center gap-1.5 whitespace-nowrap rounded-lg border px-3 py-1.5 text-xs font-semibold transition-colors disabled:cursor-not-allowed disabled:opacity-50 ${
                    environment === null
                        ? 'border-white bg-white text-[#0d5c80] shadow-sm'
                        : 'border-white/30 bg-white/15 text-white hover:bg-white/25'
                }`}
            >
                <LayoutGrid size={13} />
                Vista general
            </button>
            {ACTIONS.map(action => {
                const Icon = action.icon;
                const active = environment === action.key;
                return (
                    <button
                        key={action.key}
                        onClick={() => {
                            setEnvironment(active ? null : action.key);
                            if (!active && action.key === 'orders_compare') setView('informe');
                        }}
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
