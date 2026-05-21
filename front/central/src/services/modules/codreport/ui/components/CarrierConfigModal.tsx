'use client';

import { useCallback, useEffect, useState } from 'react';
import { X, Truck, Save, RefreshCw, Percent, CheckCircle2 } from 'lucide-react';
import { getCarrierConfigsAction, saveCarrierConfigAction } from '../../infra/actions';
import { CarrierConfig } from '../../domain/types';
import { carrierLabel } from './helpers';

interface Props {
    businessId?: number | null;
    onClose: () => void;
    onSaved?: () => void;
}

export default function CarrierConfigModal({ businessId, onClose, onSaved }: Props) {
    const [rows, setRows] = useState<CarrierConfig[]>([]);
    const [loading, setLoading] = useState(true);
    const [savingName, setSavingName] = useState<string | null>(null);
    const [feedback, setFeedback] = useState<{ ok: boolean; msg: string } | null>(null);

    const load = useCallback(async () => {
        setLoading(true);
        const res = await getCarrierConfigsAction(businessId || undefined);
        if (res.success) setRows(res.data || []);
        setLoading(false);
    }, [businessId]);

    useEffect(() => { load(); }, [load]);

    const updateRow = (name: string, patch: Partial<CarrierConfig>) => {
        setRows(rs => rs.map(r => (r.carrier_name === name ? { ...r, ...patch } : r)));
    };

    const save = async (row: CarrierConfig) => {
        setSavingName(row.carrier_name);
        setFeedback(null);
        const res = await saveCarrierConfigAction(
            {
                carrier_name: row.carrier_name,
                discount_percentage: Number(row.discount_percentage) || 0,
                is_active: row.is_active,
            },
            businessId || undefined,
        );
        if (res.success) {
            setFeedback({ ok: true, msg: `${carrierLabel(row.carrier_name)}: descuento guardado` });
            onSaved?.();
        } else {
            setFeedback({ ok: false, msg: (res as any).message || 'Error al guardar' });
        }
        setSavingName(null);
        setTimeout(() => setFeedback(null), 2800);
    };

    return (
        <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
            <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl max-w-2xl w-full max-h-[85vh] flex flex-col">
                <div className="flex items-center justify-between px-5 py-3 border-b border-gray-200 dark:border-gray-700">
                    <div className="flex items-center gap-2">
                        <Percent size={18} className="text-purple-600" />
                        <h3 className="text-base font-bold text-gray-900 dark:text-white">
                            Descuentos por transportadora
                        </h3>
                    </div>
                    <button onClick={onClose} className="p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700">
                        <X size={18} />
                    </button>
                </div>

                <div className="px-5 py-2 text-xs text-gray-500 dark:text-gray-400">
                    El porcentaje se descuenta del recaudo contra entrega de cada transportadora.
                </div>

                {feedback && (
                    <div className={`mx-5 mb-2 px-3 py-2 rounded-md text-sm ${feedback.ok ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' : 'bg-red-50 text-red-700 border border-red-200'}`}>
                        {feedback.msg}
                    </div>
                )}

                <div className="flex-1 overflow-y-auto px-5 pb-4">
                    {loading ? (
                        <div className="flex items-center justify-center py-10 text-gray-400">
                            <RefreshCw size={18} className="animate-spin mr-2" /> Cargando...
                        </div>
                    ) : rows.length === 0 ? (
                        <div className="text-center text-gray-400 py-10 text-sm">
                            <Truck size={32} className="mx-auto mb-2 opacity-50" />
                            No hay transportadoras con ordenes contra entrega.
                        </div>
                    ) : (
                        <div className="space-y-2">
                            {rows.map(row => (
                                <div
                                    key={row.carrier_name}
                                    className="flex items-center gap-3 p-3 rounded-lg border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/40"
                                >
                                    <div className="flex items-center gap-2 flex-1 min-w-0">
                                        <Truck size={15} className="text-gray-400 shrink-0" />
                                        <span className="text-sm font-semibold text-gray-900 dark:text-white truncate">
                                            {carrierLabel(row.carrier_name)}
                                        </span>
                                    </div>
                                    <div className="relative w-28">
                                        <input
                                            type="number"
                                            min={0}
                                            max={100}
                                            step={0.5}
                                            value={row.discount_percentage}
                                            onChange={e => updateRow(row.carrier_name, { discount_percentage: parseFloat(e.target.value) })}
                                            className="w-full pl-3 pr-7 py-1.5 text-sm rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                                        />
                                        <span className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 text-sm">%</span>
                                    </div>
                                    <label className="flex items-center gap-1.5 text-xs text-gray-600 dark:text-gray-300 cursor-pointer select-none">
                                        <input
                                            type="checkbox"
                                            checked={row.is_active}
                                            onChange={e => updateRow(row.carrier_name, { is_active: e.target.checked })}
                                            className="rounded"
                                        />
                                        Activa
                                    </label>
                                    <button
                                        onClick={() => save(row)}
                                        disabled={savingName === row.carrier_name}
                                        className="px-3 py-1.5 bg-purple-600 hover:bg-purple-700 disabled:opacity-50 text-white text-xs font-semibold rounded-md inline-flex items-center gap-1.5"
                                    >
                                        {savingName === row.carrier_name
                                            ? <RefreshCw size={13} className="animate-spin" />
                                            : <Save size={13} />}
                                        Guardar
                                    </button>
                                </div>
                            ))}
                        </div>
                    )}
                </div>

                <div className="px-5 py-3 border-t border-gray-200 dark:border-gray-700 flex justify-end">
                    <button
                        onClick={onClose}
                        className="px-4 py-2 text-sm rounded-md bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-200 font-semibold inline-flex items-center gap-1.5"
                    >
                        <CheckCircle2 size={14} /> Listo
                    </button>
                </div>
            </div>
        </div>
    );
}
