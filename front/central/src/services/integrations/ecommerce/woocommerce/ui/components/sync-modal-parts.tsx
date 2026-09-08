'use client';

import { useState } from 'react';
import { ChevronDown, ChevronRight } from 'lucide-react';

export interface Brief {
    sku: string;
    name: string;
    family_ref?: string;
    family_name?: string;
    variant_label?: string;
    image_url?: string;
}

export interface ChannelFamilyGroup {
    key: string;
    name: string;
    image_url?: string;
    items: Brief[];
}

export function groupByFamily(items: Brief[]): ChannelFamilyGroup[] {
    const order: string[] = [];
    const byKey = new Map<string, ChannelFamilyGroup>();
    for (const item of items) {
        const key = item.family_ref || `sku:${item.sku}`;
        if (!byKey.has(key)) {
            byKey.set(key, { key, name: item.family_name || item.name, image_url: item.image_url, items: [] });
            order.push(key);
        }
        const group = byKey.get(key)!;
        group.items.push(item);
        if (!group.image_url && item.image_url) group.image_url = item.image_url;
    }
    return order.map(key => byKey.get(key)!);
}

export type SyncAction = 'created' | 'updated' | 'failed' | 'skipped';

export function ActionBadge({ action }: { action: SyncAction }) {
    const map = {
        created: { label: 'Creado', cls: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300' },
        updated: { label: 'Actualizado', cls: 'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-300' },
        failed: { label: 'Fallido', cls: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300' },
        skipped: { label: 'Omitido', cls: 'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300' },
    };
    const { label, cls } = map[action];
    return <span className={`px-1.5 py-0.5 rounded font-semibold ${cls}`}>{label}</span>;
}

export function ProductList({ items }: { items: Brief[] }) {
    if (items.length === 0) return null;
    return (
        <div className="mt-2 max-h-32 overflow-y-auto rounded-md bg-gray-50 dark:bg-gray-800/60 divide-y divide-gray-100 dark:divide-gray-700">
            {items.slice(0, 100).map((p, i) => (
                <div key={i} className="flex items-center justify-between px-2.5 py-1.5 text-[11px]">
                    <span className="text-gray-700 dark:text-gray-200 truncate">{p.name || '(sin nombre)'}</span>
                    <span className="text-gray-400 font-mono ml-2 flex-shrink-0">{p.sku}</span>
                </div>
            ))}
            {items.length > 100 && <div className="px-2.5 py-1.5 text-[11px] text-gray-400">y {items.length - 100} más...</div>}
        </div>
    );
}

export function ChannelFamilyList({ groups, onImport }: { groups: ChannelFamilyGroup[]; onImport: (group: ChannelFamilyGroup) => void }) {
    if (groups.length === 0) return null;
    return (
        <div className="mt-2 max-h-56 overflow-y-auto rounded-md bg-gray-50 dark:bg-gray-800/60 divide-y divide-gray-100 dark:divide-gray-700">
            {groups.slice(0, 100).map((g) => (
                <div key={g.key} className="flex items-center justify-between gap-2 px-2.5 py-2 text-[11px]">
                    <div className="flex items-center gap-2 min-w-0">
                        {g.image_url ? (
                            <img src={g.image_url} alt={g.name} className="w-8 h-8 rounded-full object-cover flex-shrink-0 border border-gray-200 dark:border-gray-700" />
                        ) : (
                            <span className="w-8 h-8 rounded-full flex-shrink-0 bg-gray-200 dark:bg-gray-700" />
                        )}
                        <div className="min-w-0">
                            <p className="text-gray-700 dark:text-gray-200 truncate font-medium">{g.name || '(sin nombre)'}</p>
                            <p className="text-gray-400">{g.items.length} variante{g.items.length !== 1 ? 's' : ''}</p>
                        </div>
                    </div>
                    <button
                        onClick={() => onImport(g)}
                        className="inline-flex items-center gap-1 whitespace-nowrap rounded-md bg-blue-600 hover:bg-blue-700 px-2.5 py-1 text-[11px] font-semibold text-white transition-colors flex-shrink-0"
                    >
                        Traer a Probability
                    </button>
                </div>
            ))}
            {groups.length > 100 && <div className="px-2.5 py-1.5 text-[11px] text-gray-400">y {groups.length - 100} más...</div>}
        </div>
    );
}

export function SelectableFamilyList({ groups, selected, onToggleGroup, onToggleItem }: {
    groups: ChannelFamilyGroup[];
    selected: Set<string>;
    onToggleGroup: (group: ChannelFamilyGroup, checked: boolean) => void;
    onToggleItem: (sku: string) => void;
}) {
    const [expanded, setExpanded] = useState<Set<string>>(new Set());
    if (groups.length === 0) return null;

    const toggleExpanded = (key: string) => setExpanded((prev) => {
        const next = new Set(prev);
        if (next.has(key)) next.delete(key); else next.add(key);
        return next;
    });

    return (
        <div className="mt-2 max-h-64 overflow-y-auto rounded-md bg-gray-50 dark:bg-gray-800/60 divide-y divide-gray-100 dark:divide-gray-700">
            {groups.slice(0, 200).map((g) => {
                const isFamily = g.items.length > 1;
                const selectedCount = g.items.filter((i) => selected.has(i.sku)).length;
                const allSelected = selectedCount === g.items.length;
                const someSelected = selectedCount > 0 && !allSelected;
                const isExpanded = expanded.has(g.key);
                return (
                    <div key={g.key}>
                        <div className="flex items-center gap-2 px-2.5 py-2 text-[11px]">
                            <input
                                type="checkbox"
                                checked={allSelected}
                                ref={(el) => { if (el) el.indeterminate = someSelected; }}
                                onChange={() => onToggleGroup(g, !allSelected)}
                                className="h-3.5 w-3.5 rounded border-gray-300 text-violet-600 focus:ring-violet-500 flex-shrink-0"
                            />
                            {g.image_url ? (
                                <img src={g.image_url} alt={g.name} className="w-8 h-8 rounded-full object-cover flex-shrink-0 border border-gray-200 dark:border-gray-700" />
                            ) : (
                                <span className="w-8 h-8 rounded-full flex-shrink-0 bg-gray-200 dark:bg-gray-700" />
                            )}
                            <div className="min-w-0 flex-1 cursor-pointer" onClick={() => onToggleGroup(g, !allSelected)}>
                                <p className="text-gray-700 dark:text-gray-200 truncate font-medium flex items-center gap-1.5">
                                    {g.name || '(sin nombre)'}
                                    {isFamily && (
                                        <span className="px-1.5 py-0.5 rounded-full text-[9px] font-bold uppercase bg-violet-100 dark:bg-violet-900/40 text-violet-700 dark:text-violet-300">
                                            Familia madre
                                        </span>
                                    )}
                                </p>
                                <p className="text-gray-400">
                                    {isFamily ? `${selectedCount}/${g.items.length} variantes seleccionadas` : g.items[0]?.sku}
                                </p>
                            </div>
                            {isFamily && (
                                <button
                                    type="button"
                                    onClick={() => toggleExpanded(g.key)}
                                    className="p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 flex-shrink-0"
                                    title={isExpanded ? 'Ocultar variantes' : 'Ver variantes'}
                                >
                                    {isExpanded ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
                                </button>
                            )}
                        </div>
                        {isFamily && isExpanded && (
                            <div className="pl-8 pb-1.5 pr-2.5 space-y-1">
                                {g.items.map((item) => (
                                    <label key={item.sku} className="flex items-center gap-2 py-1 text-[11px] cursor-pointer hover:bg-violet-50/50 dark:hover:bg-violet-900/10 rounded px-1">
                                        <input
                                            type="checkbox"
                                            checked={selected.has(item.sku)}
                                            onChange={() => onToggleItem(item.sku)}
                                            className="h-3 w-3 rounded border-gray-300 text-violet-600 focus:ring-violet-500"
                                        />
                                        <span className="text-gray-600 dark:text-gray-300 truncate flex-1">{item.variant_label || item.name}</span>
                                        <span className="text-gray-400 font-mono flex-shrink-0">{item.sku}</span>
                                    </label>
                                ))}
                            </div>
                        )}
                    </div>
                );
            })}
            {groups.length > 200 && <div className="px-2.5 py-1.5 text-[11px] text-gray-400">y {groups.length - 200} más...</div>}
        </div>
    );
}

export function SelectableProductList({ items, selected, onToggle }: { items: Brief[]; selected: Set<string>; onToggle: (sku: string) => void }) {
    if (items.length === 0) return null;
    return (
        <div className="mt-2 max-h-40 overflow-y-auto rounded-md bg-white dark:bg-gray-800/60 border border-amber-100 dark:border-amber-900/40 divide-y divide-gray-100 dark:divide-gray-700">
            {items.slice(0, 200).map((p, i) => (
                <label key={i} className="flex items-center gap-2 px-2.5 py-1.5 text-[11px] cursor-pointer hover:bg-amber-50/50 dark:hover:bg-amber-900/10">
                    <input
                        type="checkbox"
                        checked={selected.has(p.sku)}
                        onChange={() => onToggle(p.sku)}
                        className="h-3.5 w-3.5 rounded border-gray-300 text-amber-600 focus:ring-amber-500"
                    />
                    <span className="text-gray-700 dark:text-gray-200 truncate flex-1">{p.name || '(sin nombre)'}</span>
                    <span className="text-gray-400 font-mono ml-2 flex-shrink-0">{p.sku}</span>
                </label>
            ))}
            {items.length > 200 && <div className="px-2.5 py-1.5 text-[11px] text-gray-400">y {items.length - 200} mas (usa "Asociar todos")...</div>}
        </div>
    );
}
