'use client';

import { useMemo, useState } from 'react';
import { Link2, Loader2 } from 'lucide-react';
import type { Integration } from '@/services/integrations/core/domain/types';
import { getSyncProvider } from '../providers';

export interface ProductBrief {
    sku: string;
    name: string;
}

export interface ChannelDiffData {
    integration: Integration;
    notAssociated: ProductBrief[];
    onlyInProbability: ProductBrief[];
    onlyInChannel: ProductBrief[];
}

interface GlobalProductsMatrixProps {
    channels: ChannelDiffData[];
    businessId: number | null;
    onAssociated: () => void;
}

type CellState = 'not_associated' | 'missing_in_channel' | 'only_in_channel' | null;

interface MatrixRow {
    sku: string;
    name: string;
    cells: Record<number, CellState>;
}

const MAX_ROWS = 100;

export function GlobalProductsMatrix({ channels, businessId, onAssociated }: GlobalProductsMatrixProps) {
    const [selected, setSelected] = useState<Record<number, Set<string>>>({});
    const [associating, setAssociating] = useState(false);
    const [resultMsg, setResultMsg] = useState<string | null>(null);

    const rows = useMemo<MatrixRow[]>(() => {
        const map = new Map<string, MatrixRow>();
        const ensure = (sku: string, name: string) => {
            let row = map.get(sku);
            if (!row) {
                row = { sku, name, cells: {} };
                map.set(sku, row);
            }
            if (!row.name && name) row.name = name;
            return row;
        };
        for (const ch of channels) {
            for (const b of ch.notAssociated) ensure(b.sku, b.name).cells[ch.integration.id] = 'not_associated';
            for (const b of ch.onlyInProbability) ensure(b.sku, b.name).cells[ch.integration.id] = 'missing_in_channel';
            for (const b of ch.onlyInChannel) ensure(b.sku, b.name).cells[ch.integration.id] = 'only_in_channel';
        }
        return [...map.values()].sort((a, b) => a.sku.localeCompare(b.sku));
    }, [channels]);

    const totalSelected = Object.values(selected).reduce((acc, set) => acc + set.size, 0);

    const toggleCell = (integrationId: number, sku: string) => {
        if (associating) return;
        setSelected(prev => {
            const next = { ...prev };
            const set = new Set(next[integrationId] || []);
            if (set.has(sku)) set.delete(sku);
            else set.add(sku);
            next[integrationId] = set;
            return next;
        });
    };

    const associate = async () => {
        if (associating || totalSelected === 0) return;
        setAssociating(true);
        setResultMsg(null);
        let okCount = 0;
        let failCount = 0;
        for (const ch of channels) {
            const skus = [...(selected[ch.integration.id] || [])];
            if (skus.length === 0) continue;
            const provider = getSyncProvider(ch.integration.integration_type_id);
            if (!provider) continue;
            try {
                const res = await provider.associateProducts(ch.integration.id, businessId ?? undefined, skus) as { success?: boolean };
                if (res?.success === false) failCount += skus.length;
                else okCount += skus.length;
            } catch {
                failCount += skus.length;
            }
        }
        setAssociating(false);
        setSelected({});
        setResultMsg(failCount === 0 ? `${okCount} productos asociados` : `${okCount} asociados, ${failCount} fallidos`);
        onAssociated();
    };

    if (rows.length === 0) {
        return (
            <p className="py-4 text-center text-xs italic text-gray-400 dark:text-gray-500">
                Todos los productos estan alineados entre los canales comparados
            </p>
        );
    }

    const visible = rows.slice(0, MAX_ROWS);

    return (
        <div className="flex flex-col gap-2">
            <div className="overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700">
                <table className="w-full min-w-[36rem] text-left text-[11px]">
                    <thead>
                        <tr className="border-b border-gray-200 bg-gray-50 dark:border-gray-700 dark:bg-gray-900/40">
                            <th className="sticky left-0 bg-gray-50 px-3 py-2 font-semibold text-gray-600 dark:bg-gray-900/40 dark:text-gray-300">
                                Producto
                            </th>
                            {channels.map(ch => {
                                const typeName = ch.integration.integration_type?.name || ch.integration.name;
                                return (
                                    <th key={ch.integration.id} className="px-3 py-2 font-semibold text-gray-600 dark:text-gray-300">
                                        <span className="flex items-center gap-1.5 whitespace-nowrap">
                                            {ch.integration.integration_type?.image_url && (
                                                <img
                                                    src={ch.integration.integration_type.image_url}
                                                    alt={typeName}
                                                    className="h-4 w-4 rounded-full object-contain"
                                                />
                                            )}
                                            {typeName}
                                        </span>
                                    </th>
                                );
                            })}
                        </tr>
                    </thead>
                    <tbody>
                        {visible.map(row => (
                            <tr key={row.sku} className="border-b border-gray-100 last:border-0 hover:bg-gray-50/60 dark:border-gray-700/60 dark:hover:bg-gray-700/30">
                                <td className="sticky left-0 max-w-[14rem] bg-white px-3 py-1.5 dark:bg-gray-800">
                                    <span className="block truncate font-semibold text-gray-800 dark:text-gray-100">{row.sku}</span>
                                    <span className="block truncate text-[10px] text-gray-500 dark:text-gray-400">{row.name}</span>
                                </td>
                                {channels.map(ch => {
                                    const state = row.cells[ch.integration.id] ?? null;
                                    return (
                                        <td key={ch.integration.id} className="px-3 py-1.5">
                                            {state === null && <span className="text-gray-300 dark:text-gray-600">-</span>}
                                            {state === 'not_associated' && (
                                                <label className="flex w-fit cursor-pointer items-center gap-1.5 rounded-full bg-indigo-50 px-2 py-0.5 font-semibold text-indigo-700 dark:bg-indigo-900/40 dark:text-indigo-300">
                                                    <input
                                                        type="checkbox"
                                                        checked={selected[ch.integration.id]?.has(row.sku) || false}
                                                        onChange={() => toggleCell(ch.integration.id, row.sku)}
                                                        disabled={associating}
                                                        className="h-3 w-3 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
                                                    />
                                                    Asociar
                                                </label>
                                            )}
                                            {state === 'missing_in_channel' && (
                                                <span className="whitespace-nowrap rounded-full bg-amber-50 px-2 py-0.5 font-semibold text-amber-700 dark:bg-amber-900/40 dark:text-amber-300">
                                                    Falta en canal
                                                </span>
                                            )}
                                            {state === 'only_in_channel' && (
                                                <span className="whitespace-nowrap rounded-full bg-purple-50 px-2 py-0.5 font-semibold text-purple-700 dark:bg-purple-900/40 dark:text-purple-300">
                                                    Solo en canal
                                                </span>
                                            )}
                                        </td>
                                    );
                                })}
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>

            <div className="flex items-center justify-between gap-2">
                <span className="text-[11px] text-gray-400 dark:text-gray-500">
                    {rows.length > MAX_ROWS
                        ? `Mostrando ${MAX_ROWS} de ${rows.length} productos con diferencias`
                        : `${rows.length} productos con diferencias`}
                    {resultMsg ? ` - ${resultMsg}` : ''}
                </span>
                <button
                    onClick={associate}
                    disabled={associating || totalSelected === 0}
                    className="flex items-center gap-1.5 rounded-lg bg-indigo-600 px-3 py-1.5 text-[11px] font-semibold text-white transition-colors hover:bg-indigo-700 disabled:cursor-not-allowed disabled:opacity-50"
                >
                    {associating ? <Loader2 size={12} className="animate-spin" /> : <Link2 size={12} />}
                    Asociar seleccionados ({totalSelected})
                </button>
            </div>

            <p className="text-[10px] text-gray-400 dark:text-gray-500">
                Para crear productos faltantes usa Revisar en el canal correspondiente.
            </p>
        </div>
    );
}
