'use client';

import { useCallback, useEffect, useMemo, useState } from 'react';
import { ClientGroup, ClientSummary, CatalogPriceRow, CatalogPriceTarget } from '../../domain/types';
import {
    getCatalogPricesAction,
    saveCatalogPricesAction,
    listAvailableClientsAction,
} from '../../infra/actions';

interface CatalogPriceTableProps {
    businessId?: number;
    groups: ClientGroup[];
}

type TargetMode = 'group' | 'client';

function formatCurrency(value: number, currency: string): string {
    try {
        return new Intl.NumberFormat('es-CO', {
            style: 'currency',
            currency: currency || 'COP',
            maximumFractionDigits: 0,
        }).format(value);
    } catch {
        return String(value);
    }
}

function nudgeStep(base: number): number {
    return Math.max(1, Math.round(base * 0.01));
}

export function CatalogPriceTable({ businessId, groups }: CatalogPriceTableProps) {
    const [mode, setMode] = useState<TargetMode>('group');
    const [groupId, setGroupId] = useState<number | null>(null);
    const [client, setClient] = useState<ClientSummary | null>(null);

    const [clientSearch, setClientSearch] = useState('');
    const [clientResults, setClientResults] = useState<ClientSummary[]>([]);

    const [search, setSearch] = useState('');
    const [rows, setRows] = useState<CatalogPriceRow[]>([]);
    const [loading, setLoading] = useState(false);
    const [draft, setDraft] = useState<Record<string, number>>({});
    const [dirty, setDirty] = useState<Set<string>>(new Set());
    const [saving, setSaving] = useState(false);
    const [message, setMessage] = useState('');
    const [selected, setSelected] = useState<Set<string>>(new Set());
    const [bulkAmount, setBulkAmount] = useState('');

    const target: CatalogPriceTarget | null = useMemo(() => {
        if (mode === 'group' && groupId) return { client_group_id: groupId };
        if (mode === 'client' && client) return { client_id: client.id };
        return null;
    }, [mode, groupId, client]);

    const baseByProduct = useMemo(() => {
        const map: Record<string, number> = {};
        rows.forEach((row) => { map[row.product_id] = row.base_price; });
        return map;
    }, [rows]);

    const loadPrices = useCallback(async () => {
        if (!target) return;
        setLoading(true);
        setMessage('');
        const result = await getCatalogPricesAction(businessId, target, search, 1);
        setRows(result.data);
        const initial: Record<string, number> = {};
        result.data.forEach((row) => {
            initial[row.product_id] = row.custom_price ?? row.base_price;
        });
        setDraft(initial);
        setDirty(new Set());
        setSelected(new Set());
        setLoading(false);
    }, [businessId, target, search]);

    useEffect(() => {
        if (target) loadPrices();
        else setRows([]);
    }, [target, loadPrices]);

    useEffect(() => {
        if (mode !== 'client') return;
        const handle = setTimeout(async () => {
            const result = await listAvailableClientsAction(businessId, clientSearch, false, 1);
            setClientResults(result.data);
        }, 300);
        return () => clearTimeout(handle);
    }, [mode, clientSearch, businessId]);

    const setPrice = (productId: string, value: number) => {
        const safe = Number.isFinite(value) && value >= 0 ? value : 0;
        setDraft((prev) => ({ ...prev, [productId]: safe }));
        setDirty((prev) => new Set(prev).add(productId));
    };

    const allSelected = rows.length > 0 && selected.size === rows.length;

    const toggleSelected = (productId: string) => {
        setSelected((prev) => {
            const next = new Set(prev);
            if (next.has(productId)) next.delete(productId);
            else next.add(productId);
            return next;
        });
    };

    const toggleSelectAll = () => {
        setSelected(allSelected ? new Set() : new Set(rows.map((r) => r.product_id)));
    };

    const applyBulk = (sign: 1 | -1) => {
        const amount = Number(bulkAmount);
        if (!Number.isFinite(amount) || amount <= 0) return;
        const targetIds = selected.size > 0 ? [...selected] : rows.map((r) => r.product_id);
        setDraft((prev) => {
            const next = { ...prev };
            targetIds.forEach((pid) => {
                const current = next[pid] ?? baseByProduct[pid] ?? 0;
                const updated = current + sign * amount;
                next[pid] = updated >= 0 ? updated : 0;
            });
            return next;
        });
        setDirty((prev) => {
            const next = new Set(prev);
            targetIds.forEach((pid) => next.add(pid));
            return next;
        });
    };

    const handleSave = async () => {
        if (!target || dirty.size === 0) return;
        setSaving(true);
        setMessage('');
        const items = [...dirty].map((productId) => {
            const value = draft[productId];
            const base = baseByProduct[productId];
            return { product_id: productId, price: value === base ? null : value };
        });
        const result = await saveCatalogPricesAction(businessId, target, items);
        setSaving(false);
        if (!result.success) {
            setMessage(result.message || 'No se pudieron guardar los precios');
            return;
        }
        setMessage('Precios guardados');
        await loadPrices();
    };

    return (
        <div className="space-y-3">
            <div className="flex flex-wrap items-end gap-3">
                <div>
                    <label className="block text-xs font-bold text-gray-700 dark:text-gray-200 mb-1">Aplicar a</label>
                    <div className="flex gap-1 bg-gray-100 dark:bg-gray-800 p-1 rounded-lg">
                        <button
                            onClick={() => { setMode('group'); setClient(null); }}
                            className={`px-3 py-1.5 text-sm font-bold rounded-md ${mode === 'group' ? 'btn-business-primary text-white' : 'text-gray-600 dark:text-gray-300'}`}
                        >
                            Grupo
                        </button>
                        <button
                            onClick={() => { setMode('client'); setGroupId(null); }}
                            className={`px-3 py-1.5 text-sm font-bold rounded-md ${mode === 'client' ? 'btn-business-primary text-white' : 'text-gray-600 dark:text-gray-300'}`}
                        >
                            Cliente individual
                        </button>
                    </div>
                </div>

                {mode === 'group' ? (
                    <div>
                        <label className="block text-xs font-bold text-gray-700 dark:text-gray-200 mb-1">Grupo</label>
                        <div className="flex items-center gap-2">
                            {groupId && (
                                <span
                                    className="w-4 h-4 rounded-full flex-shrink-0"
                                    style={{ backgroundColor: groups.find((g) => g.id === groupId)?.color || '#6b7280' }}
                                />
                            )}
                            <select
                                value={groupId ?? ''}
                                onChange={(e) => setGroupId(e.target.value ? Number(e.target.value) : null)}
                                className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                            >
                                <option value="">Selecciona un grupo</option>
                                {groups.map((group) => (
                                    <option key={group.id} value={group.id}>{group.name}</option>
                                ))}
                            </select>
                        </div>
                    </div>
                ) : (
                    <div>
                        <label className="block text-xs font-bold text-gray-700 dark:text-gray-200 mb-1">Cliente</label>
                        {client ? (
                            <div className="flex items-center gap-2">
                                <span className="px-3 py-2 border border-business-primary rounded-lg text-sm bg-business-primary/5 text-gray-900 dark:text-white">
                                    {client.name}
                                </span>
                                <button
                                    onClick={() => setClient(null)}
                                    className="text-xs text-business-primary underline"
                                >
                                    Cambiar
                                </button>
                            </div>
                        ) : (
                            <span className="inline-block px-3 py-2 text-sm text-gray-400">
                                Selecciona un cliente abajo
                            </span>
                        )}
                    </div>
                )}

                <div className="flex-1 min-w-48">
                    <label className="block text-xs font-bold text-gray-700 dark:text-gray-200 mb-1">Buscar producto</label>
                    <input
                        type="text"
                        placeholder="Nombre o SKU"
                        value={search}
                        onChange={(e) => setSearch(e.target.value)}
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                    />
                </div>

                <button
                    onClick={handleSave}
                    disabled={saving || dirty.size === 0}
                    className="px-5 py-2 btn-business-primary text-white text-sm font-bold rounded-lg disabled:opacity-50"
                >
                    {saving ? 'Guardando...' : `Guardar cambios${dirty.size > 0 ? ` (${dirty.size})` : ''}`}
                </button>
            </div>

            {mode === 'client' && !client && (
                <div className="border border-gray-200 dark:border-gray-700 rounded-lg p-3 space-y-2">
                    <input
                        type="text"
                        placeholder="Buscar cliente por nombre, email o documento"
                        value={clientSearch}
                        onChange={(e) => setClientSearch(e.target.value)}
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                    />
                    <div className="max-h-60 overflow-y-auto space-y-1">
                        {clientResults.length === 0 ? (
                            <p className="text-sm text-gray-500">Escribe para buscar clientes.</p>
                        ) : (
                            clientResults.map((c) => (
                                <button
                                    key={c.id}
                                    onClick={() => { setClient(c); setClientSearch(''); setClientResults([]); }}
                                    className="block w-full text-left px-3 py-2 rounded text-sm text-gray-900 dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700"
                                >
                                    {c.name}
                                    {c.group_name && <span className="text-xs text-gray-400"> - {c.group_name}</span>}
                                </button>
                            ))
                        )}
                    </div>
                </div>
            )}

            {message && <p className="text-sm text-business-primary">{message}</p>}

            {target && !loading && rows.length > 0 && (
                <div className="flex flex-wrap items-center gap-2 bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-3">
                    <span className="text-sm font-bold text-gray-700 dark:text-gray-200">Ajuste masivo</span>
                    <span className="text-xs text-gray-500">
                        {selected.size > 0 ? `${selected.size} productos seleccionados` : 'aplica a todos los productos'}
                    </span>
                    <input
                        type="number"
                        placeholder="Monto"
                        value={bulkAmount}
                        onChange={(e) => setBulkAmount(e.target.value)}
                        className="w-28 px-2 py-1.5 border border-gray-300 dark:border-gray-600 rounded-md text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                    />
                    <button
                        onClick={() => applyBulk(1)}
                        className="px-3 py-1.5 rounded-md border border-green-300 text-green-700 hover:bg-green-50 dark:border-green-800 dark:text-green-400 dark:hover:bg-green-900/30 text-sm font-bold"
                    >
                        + Aumentar
                    </button>
                    <button
                        onClick={() => applyBulk(-1)}
                        className="px-3 py-1.5 rounded-md border border-red-300 text-red-700 hover:bg-red-50 dark:border-red-800 dark:text-red-400 dark:hover:bg-red-900/30 text-sm font-bold"
                    >
                        - Disminuir
                    </button>
                </div>
            )}

            {!target ? (
                <p className="text-sm text-gray-500">
                    {mode === 'group' ? 'Selecciona un grupo para ver y editar sus precios.' : 'Busca y selecciona un cliente.'}
                </p>
            ) : loading ? (
                <p className="text-sm text-gray-500">Cargando precios...</p>
            ) : rows.length === 0 ? (
                <p className="text-sm text-gray-500">No hay productos activos para este negocio.</p>
            ) : (
                <div className="border border-gray-200 dark:border-gray-700 rounded-xl overflow-hidden">
                    <table className="w-full text-sm">
                        <thead className="bg-gray-50 dark:bg-gray-800">
                            <tr className="text-left text-gray-600 dark:text-gray-300">
                                <th className="px-3 py-3">
                                    <input
                                        type="checkbox"
                                        checked={allSelected}
                                        onChange={toggleSelectAll}
                                        aria-label="Seleccionar todos"
                                    />
                                </th>
                                <th className="px-4 py-3 font-bold">Producto</th>
                                <th className="px-4 py-3 font-bold">Precio base</th>
                                <th className="px-4 py-3 font-bold">Precio para este grupo/cliente</th>
                                <th className="px-4 py-3 font-bold">Diferencia</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100 dark:divide-gray-700 max-h-96">
                            {rows.map((row) => {
                                const value = draft[row.product_id] ?? row.base_price;
                                const diff = value - row.base_price;
                                const step = nudgeStep(row.base_price);
                                return (
                                    <tr key={row.product_id} className="text-gray-900 dark:text-white">
                                        <td className="px-3 py-3">
                                            <input
                                                type="checkbox"
                                                checked={selected.has(row.product_id)}
                                                onChange={() => toggleSelected(row.product_id)}
                                                aria-label={`Seleccionar ${row.product_name}`}
                                            />
                                        </td>
                                        <td className="px-4 py-3">
                                            <div className="flex items-center gap-3">
                                                {row.image_url || (row as any).family_image_url ? (
                                                    <img
                                                        src={row.image_url || (row as any).family_image_url}
                                                        alt={row.product_name}
                                                        className="w-10 h-10 rounded-md object-cover border border-gray-200 dark:border-gray-700 flex-shrink-0"
                                                    />
                                                ) : (
                                                    <div className="w-10 h-10 rounded-md bg-gray-100 dark:bg-gray-700 border border-gray-200 dark:border-gray-700 flex items-center justify-center flex-shrink-0 text-gray-400 text-xs">
                                                        sin foto
                                                    </div>
                                                )}
                                                <div>
                                                    <p className="font-medium">{row.product_name}</p>
                                                    <p className="text-xs text-gray-500">{row.product_sku}</p>
                                                </div>
                                            </div>
                                        </td>
                                        <td className="px-4 py-3 text-gray-600 dark:text-gray-300">
                                            {formatCurrency(row.base_price, row.currency)}
                                        </td>
                                        <td className="px-4 py-3">
                                            <div className="flex items-center gap-1">
                                                <button
                                                    onClick={() => setPrice(row.product_id, value - step)}
                                                    className="w-7 h-7 rounded-md border border-red-300 text-red-600 hover:bg-red-50 dark:border-red-800 dark:text-red-400 dark:hover:bg-red-900/30 font-bold"
                                                >
                                                    -
                                                </button>
                                                <input
                                                    type="number"
                                                    value={value}
                                                    onChange={(e) => setPrice(row.product_id, Number(e.target.value))}
                                                    className="w-32 px-2 py-1.5 border border-gray-300 dark:border-gray-600 rounded-md text-sm text-center bg-white dark:bg-gray-700"
                                                />
                                                <button
                                                    onClick={() => setPrice(row.product_id, value + step)}
                                                    className="w-7 h-7 rounded-md border border-green-300 text-green-600 hover:bg-green-50 dark:border-green-800 dark:text-green-400 dark:hover:bg-green-900/30 font-bold"
                                                >
                                                    +
                                                </button>
                                                {diff !== 0 && (
                                                    <button
                                                        onClick={() => setPrice(row.product_id, row.base_price)}
                                                        className="ml-1 text-xs text-gray-400 underline"
                                                    >
                                                        base
                                                    </button>
                                                )}
                                            </div>
                                        </td>
                                        <td className="px-4 py-3">
                                            {diff === 0 ? (
                                                <span className="text-xs text-gray-400">Igual al base</span>
                                            ) : (
                                                <span
                                                    className={`px-2 py-1 rounded-md text-xs font-bold ${
                                                        diff > 0
                                                            ? 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-300'
                                                            : 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300'
                                                    }`}
                                                >
                                                    {diff > 0 ? '+' : '-'}{formatCurrency(Math.abs(diff), row.currency)}
                                                </span>
                                            )}
                                        </td>
                                    </tr>
                                );
                            })}
                        </tbody>
                    </table>
                </div>
            )}
        </div>
    );
}
