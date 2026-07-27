'use client';

import { useMemo, useState } from 'react';
import type { ProductApplyActions } from '../providers';
import type { DetailGroup, ProductActionKey, SyncDetailItem } from '../sync-activity-context';

interface SyncDetailPanelProps {
    items: SyncDetailItem[];
    providerLabel: string;
    apply: ProductApplyActions;
    isProducts: boolean;
    busyAction: ProductActionKey | null;
    onApply: (action: ProductActionKey, skus: string[]) => void;
}

interface GroupStyle {
    label: string;
    dot: string;
    text: string;
    chip: string;
}

const GROUP_STYLES: Record<DetailGroup, GroupStyle> = {
    both: {
        label: 'En ambos',
        dot: 'bg-emerald-500',
        text: 'text-emerald-700 dark:text-emerald-300',
        chip: 'border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-500/40 dark:bg-emerald-900/30 dark:text-emerald-300',
    },
    not_associated: {
        label: 'Sin asociar',
        dot: 'bg-blue-500',
        text: 'text-blue-700 dark:text-blue-300',
        chip: 'border-blue-200 bg-blue-50 text-blue-700 dark:border-blue-500/40 dark:bg-blue-900/30 dark:text-blue-300',
    },
    only_probability: {
        label: 'Solo en Probability',
        dot: 'bg-orange-500',
        text: 'text-orange-700 dark:text-orange-300',
        chip: 'border-orange-200 bg-orange-50 text-orange-700 dark:border-orange-500/40 dark:bg-orange-900/30 dark:text-orange-300',
    },
    only_channel: {
        label: 'Solo en el canal',
        dot: 'bg-fuchsia-500',
        text: 'text-fuchsia-700 dark:text-fuchsia-300',
        chip: 'border-fuchsia-200 bg-fuchsia-50 text-fuchsia-700 dark:border-fuchsia-500/40 dark:bg-fuchsia-900/30 dark:text-fuchsia-300',
    },
    updated: {
        label: 'Actualizados',
        dot: 'bg-emerald-500',
        text: 'text-emerald-700 dark:text-emerald-300',
        chip: 'border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-500/40 dark:bg-emerald-900/30 dark:text-emerald-300',
    },
    skipped: {
        label: 'Omitidos',
        dot: 'bg-amber-500',
        text: 'text-amber-700 dark:text-amber-300',
        chip: 'border-amber-200 bg-amber-50 text-amber-700 dark:border-amber-500/40 dark:bg-amber-900/30 dark:text-amber-300',
    },
    failed: {
        label: 'Fallidos',
        dot: 'bg-red-500',
        text: 'text-red-600 dark:text-red-400',
        chip: 'border-red-200 bg-red-50 text-red-600 dark:border-red-500/40 dark:bg-red-900/30 dark:text-red-300',
    },
};

const PRODUCT_ORDER: DetailGroup[] = ['both', 'not_associated', 'only_probability', 'only_channel'];
const INVENTORY_ORDER: DetailGroup[] = ['updated', 'skipped', 'failed'];

const GROUP_ACTION: Partial<Record<DetailGroup, { key: ProductActionKey; label: (channel: string) => string }>> = {
    not_associated: { key: 'associate', label: () => 'Asociar al canal' },
    only_probability: { key: 'createInChannel', label: channel => `Crear en ${channel}` },
    only_channel: { key: 'createInProbability', label: () => 'Crear en Probability' },
    both: { key: 'updateInProbability', label: () => 'Actualizar en Probability' },
};

export function SyncDetailPanel({ items, providerLabel, apply, isProducts, busyAction, onApply }: SyncDetailPanelProps) {
    const [filter, setFilter] = useState<DetailGroup | 'all'>('all');
    const [selected, setSelected] = useState<Set<string>>(new Set());

    const counts = useMemo(() => {
        const map = new Map<DetailGroup, number>();
        for (const item of items) map.set(item.group, (map.get(item.group) ?? 0) + 1);
        return map;
    }, [items]);

    const order = isProducts ? PRODUCT_ORDER : INVENTORY_ORDER;
    const visible = filter === 'all' ? items : items.filter(item => item.group === filter);

    const action = filter === 'all' ? null : GROUP_ACTION[filter] ?? null;
    const canApply = action !== null && (action.key === 'associate' || apply[action.key] !== undefined);

    const pick = (group: DetailGroup | 'all') => {
        setFilter(group);
        setSelected(new Set());
    };

    const toggle = (sku: string) => {
        setSelected(prev => {
            const next = new Set(prev);
            if (next.has(sku)) next.delete(sku);
            else next.add(sku);
            return next;
        });
    };

    const allVisibleSelected = visible.length > 0 && visible.every(item => selected.has(item.sku));

    return (
        <div className="space-y-1.5">
            <div className="flex flex-wrap items-center gap-1">
                <button
                    onClick={() => pick('all')}
                    className={`rounded-full border px-2 py-0.5 text-[10.5px] font-bold transition-colors ${
                        filter === 'all'
                            ? 'border-gray-400 bg-gray-700 text-white dark:border-gray-500 dark:bg-gray-600'
                            : 'border-gray-200 bg-white text-gray-500 hover:bg-gray-50 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-300'
                    }`}
                >
                    Todos {items.length}
                </button>
                {order.map(group => {
                    const count = counts.get(group) ?? 0;
                    if (count === 0) return null;
                    const style = GROUP_STYLES[group];
                    return (
                        <button
                            key={group}
                            onClick={() => pick(group)}
                            className={`rounded-full border px-2 py-0.5 text-[10.5px] font-bold transition-colors ${style.chip} ${
                                filter === group ? 'ring-2 ring-offset-1 ring-current dark:ring-offset-gray-800' : 'opacity-80 hover:opacity-100'
                            }`}
                        >
                            {style.label} {count}
                        </button>
                    );
                })}
            </div>

            <div className="max-h-40 overflow-y-auto rounded-lg bg-gray-50 p-1.5 dark:bg-gray-900/40">
                {visible.map((item, index) => {
                    const style = GROUP_STYLES[item.group];
                    return (
                        <label
                            key={`${item.group}-${item.sku}-${index}`}
                            className="flex cursor-pointer items-start gap-2 border-b border-gray-100 px-1 py-1 text-[11px] last:border-0 hover:bg-white dark:border-gray-700/60 dark:hover:bg-gray-800/60"
                        >
                            {canApply && (
                                <input
                                    type="checkbox"
                                    checked={selected.has(item.sku)}
                                    onChange={() => toggle(item.sku)}
                                    className="mt-0.5 h-3 w-3 flex-shrink-0 accent-blue-600"
                                />
                            )}
                            <span className={`mt-1 h-1.5 w-1.5 flex-shrink-0 rounded-full ${style.dot}`} />
                            <span className="w-28 flex-shrink-0 truncate font-mono font-semibold text-gray-700 dark:text-gray-200">{item.sku}</span>
                            <span className={`min-w-0 flex-1 truncate ${style.text}`}>{item.label}</span>
                        </label>
                    );
                })}
            </div>

            {action && canApply && (
                <div className="flex flex-wrap items-center gap-2">
                    <button
                        onClick={() => setSelected(allVisibleSelected ? new Set() : new Set(visible.map(item => item.sku)))}
                        className="rounded-lg border border-gray-200 px-2 py-1 text-[11px] font-semibold text-gray-600 transition-colors hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700"
                    >
                        {allVisibleSelected ? 'Quitar seleccion' : `Seleccionar los ${visible.length}`}
                    </button>
                    <button
                        onClick={() => onApply(action.key, Array.from(selected))}
                        disabled={selected.size === 0 || busyAction !== null}
                        className="flex items-center gap-1.5 rounded-lg bg-[#0d5c80] px-3 py-1 text-[11px] font-bold text-white transition-colors hover:bg-[#0a4964] disabled:cursor-not-allowed disabled:opacity-40"
                    >
                        {busyAction === action.key && (
                            <span className="h-3 w-3 animate-spin rounded-full border border-transparent border-t-current" />
                        )}
                        {action.label(providerLabel)}
                        {selected.size > 0 && <span className="rounded-full bg-white/25 px-1.5 tabular-nums">{selected.size}</span>}
                    </button>
                    {selected.size === 0 && (
                        <span className="text-[11px] italic text-gray-400 dark:text-gray-500">Elige los productos a aplicar</span>
                    )}
                </div>
            )}
        </div>
    );
}
