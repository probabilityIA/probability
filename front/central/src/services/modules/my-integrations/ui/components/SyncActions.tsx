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
        <span className="flex flex-wrap items-center gap-x-1 gap-y-1">
            <span className="flex items-center gap-1 rounded-xl bg-gray-100 p-1 dark:bg-gray-700">
                {VIEWS.map(item => {
                    const Icon = item.icon;
                    const active = view === item.key;
                    return (
                        <button
                            key={item.key}
                            onClick={() => setView(item.key)}
                            title={item.hint}
                            style={active ? { color: 'var(--color-primary)' } : {}}
                            className={`flex items-center gap-1.5 whitespace-nowrap rounded-lg px-2.5 py-1.5 text-xs font-semibold transition-colors ${
                                active
                                    ? 'bg-white shadow-sm dark:bg-gray-800'
                                    : 'text-gray-500 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200'
                            }`}
                        >
                            <Icon size={13} />
                            {item.label}
                        </button>
                    );
                })}
            </span>

            <span className="mx-1 h-4 w-px bg-gray-200 dark:bg-gray-700" />

            <button
                onClick={() => { reset(); setEnvironment(null); }}
                disabled={running}
                title="Resumen de órdenes: cuántas entraron por cada canal y como van"
                style={environment === null ? { color: 'var(--color-primary)' } : {}}
                className={`flex items-center gap-1.5 whitespace-nowrap px-2.5 py-1.5 text-xs font-semibold transition-colors disabled:cursor-not-allowed disabled:opacity-50 ${
                    environment === null ? '' : 'text-gray-500 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200'
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
                        style={active ? { color: 'var(--color-primary)' } : {}}
                        className={`flex items-center gap-1.5 whitespace-nowrap px-2.5 py-1.5 text-xs font-semibold transition-colors disabled:cursor-not-allowed disabled:opacity-50 ${
                            active ? '' : 'text-gray-500 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200'
                        }`}
                    >
                        <Icon size={13} />
                        {action.label}
                    </button>
                );
            })}
            {finished && (
                <>
                    <span className="mx-1 h-4 w-px bg-gray-200 dark:bg-gray-700" />
                    <button
                        onClick={reset}
                        className="whitespace-nowrap px-2.5 py-1.5 text-xs font-semibold text-gray-500 transition-colors hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200"
                    >
                        Reiniciar
                    </button>
                </>
            )}
        </span>
    );
}
