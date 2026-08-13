'use client';

import { useCallback, useEffect, useState } from 'react';
import { createPortal } from 'react-dom';
import { AlertTriangle, Loader2, RotateCcw } from 'lucide-react';
import {
    fetchDataSummary,
    undoChannelData,
    type DataApplyResult,
    type DataMode,
    type DataSummaryCell,
    type DataSummaryRow,
} from '../../infra/repository/channel-data';
import { ChannelDataApplyModal, fechaCorta, type ObjetivoDato } from './ChannelDataApplyModal';
import { channelBrand } from '../../domain/types';
import { CARD_BORDER } from '../panel-theme';

interface ChannelDataTableProps {
    businessId: number | null;
}

function BotonCanal({ cell, onPick }: { cell: DataSummaryCell; onPick: (mode: DataMode) => void }) {
    const brand = channelBrand(cell.channel_code);

    if (cell.can_fill === 0 && cell.can_overwrite === 0) {
        return <span className="text-[11px] text-gray-300 dark:text-gray-600">—</span>;
    }

    return (
        <div className="flex flex-col items-start gap-1">
            {cell.can_fill > 0 && (
                <button
                    onClick={() => onPick('fill_empty')}
                    className={`inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1 text-[11.5px] font-semibold transition-all hover:opacity-80 ${brand.chip}`}
                >
                    <span className="h-2 w-2 flex-shrink-0 rounded-full" style={{ backgroundColor: brand.dot }} />
                    Llenar {cell.can_fill.toLocaleString('es-CO')} vacios
                </button>
            )}
            {cell.can_overwrite > 0 && (
                <button
                    onClick={() => onPick('overwrite')}
                    className="inline-flex items-center gap-1 rounded-full border border-amber-200 bg-amber-50 px-2 py-0.5 text-[10.5px] font-semibold text-amber-700 transition-colors hover:bg-amber-100 dark:border-amber-500/40 dark:bg-amber-900/20 dark:text-amber-300"
                >
                    <AlertTriangle size={9} />
                    Reemplazar {cell.can_overwrite.toLocaleString('es-CO')}
                </button>
            )}
        </div>
    );
}

export function ChannelDataTable({ businessId }: ChannelDataTableProps) {
    const [rows, setRows] = useState<DataSummaryRow[]>([]);
    const [snapshotAt, setSnapshotAt] = useState<string | undefined>(undefined);
    const [loading, setLoading] = useState(true);
    const [objetivo, setObjetivo] = useState<ObjetivoDato | null>(null);
    const [ultimo, setUltimo] = useState<DataApplyResult | null>(null);
    const [deshaciendo, setDeshaciendo] = useState(false);

    const cargar = useCallback(async (signal?: AbortSignal) => {
        setLoading(true);
        const summary = await fetchDataSummary(businessId ?? undefined, signal);
        if (signal?.aborted) return;
        setRows(summary.rows);
        setSnapshotAt(summary.snapshot_at);
        setLoading(false);
    }, [businessId]);

    useEffect(() => {
        const controller = new AbortController();
        void cargar(controller.signal);
        return () => controller.abort();
    }, [cargar]);

    const canales = rows[0]?.cells ?? [];

    const deshacer = async () => {
        if (!ultimo) return;
        setDeshaciendo(true);
        await undoChannelData(businessId ?? undefined, ultimo.batch_id);
        setDeshaciendo(false);
        setUltimo(null);
        void cargar();
    };

    return (
        <div className="flex min-h-0 flex-1 flex-col gap-2">
            <div className="flex flex-wrap items-center justify-between gap-2">
                <p className="text-[12px] text-gray-500 dark:text-gray-400">
                    Cada boton <span className="font-semibold text-gray-600 dark:text-gray-300">actualiza ese dato en Probability</span> con lo que dice el canal.
                    El numero es cuantos productos cambiarian.{' '}
                    {snapshotAt && <span className="text-gray-400">Foto de los canales tomada el {fechaCorta(snapshotAt)}.</span>}
                </p>
                {ultimo && (
                    <div className="inline-flex items-center gap-2 rounded-full border border-emerald-200 bg-emerald-50 px-3 py-1 text-[11.5px] font-semibold text-emerald-800 dark:border-emerald-500/40 dark:bg-emerald-900/20 dark:text-emerald-300">
                        Se actualizaron {ultimo.applied.toLocaleString('es-CO')} productos en Probability
                        <button
                            onClick={deshacer}
                            disabled={deshaciendo}
                            className="inline-flex items-center gap-1 rounded-full bg-white px-2 py-0.5 text-[10.5px] font-bold text-emerald-700 hover:bg-emerald-100 disabled:opacity-50 dark:bg-gray-800"
                        >
                            {deshaciendo ? <Loader2 size={10} className="animate-spin" /> : <RotateCcw size={10} />}
                            Deshacer
                        </button>
                    </div>
                )}
            </div>

            <div className="min-h-0 flex-1 overflow-auto rounded-xl border" style={{ borderColor: CARD_BORDER }}>
                <table className="w-full border-collapse">
                    <thead className="sticky top-0 z-10 bg-gray-50 dark:bg-gray-800">
                        <tr className="text-left">
                            <th className="min-w-[14rem] px-3 py-2 text-[10px] font-bold uppercase tracking-wider text-gray-400">Dato en Probability</th>
                            {canales.map(cell => (
                                <th key={cell.integration_id} className="min-w-[11rem] px-3 py-2">
                                    <span className="inline-flex items-center gap-1.5 whitespace-nowrap text-[11.5px] font-semibold text-gray-700 dark:text-gray-200">
                                        <span className="h-2 w-2 rounded-full" style={{ backgroundColor: channelBrand(cell.channel_code).dot }} />
                                        {cell.integration_name}
                                    </span>
                                </th>
                            ))}
                        </tr>
                    </thead>
                    <tbody>
                        {rows.map(row => (
                            <tr key={row.field} className="border-t" style={{ borderColor: CARD_BORDER }}>
                                <td className="px-3 py-2.5 align-top">
                                    <span className="block text-[12.5px] font-bold text-gray-800 dark:text-gray-100">{row.label}</span>
                                    {row.note && <span className="mt-0.5 block text-[10.5px] italic text-gray-400">{row.note}</span>}
                                </td>
                                {row.cells.map(cell => (
                                    <td key={cell.integration_id} className="px-3 py-2.5 align-top">
                                        <BotonCanal
                                            cell={cell}
                                            onPick={mode => setObjetivo({ field: row.field, label: row.label, cell, mode })}
                                        />
                                    </td>
                                ))}
                            </tr>
                        ))}
                    </tbody>
                </table>

                {loading && (
                    <p className="flex items-center justify-center gap-2 py-8 text-[12px] text-gray-500">
                        <Loader2 size={14} className="animate-spin" />
                        Revisando que trae cada canal
                    </p>
                )}

                {!loading && rows.length === 0 && (
                    <p className="px-3 py-10 text-center text-[12px] italic text-gray-400">
                        No hay nada por actualizar en Probability. Corre una comparacion de catalogo para refrescar la foto de cada canal.
                    </p>
                )}
            </div>

            <p className="text-[11px] text-gray-400">
                Los datos entran a Probability y nunca salen hacia los canales. El stock y el SKU no se tocan.
            </p>

            {objetivo && (
                <ChannelDataApplyModal
                    objetivo={objetivo}
                    businessId={businessId}
                    onClose={() => setObjetivo(null)}
                    onApplied={result => {
                        setObjetivo(null);
                        setUltimo(result);
                        void cargar();
                    }}
                />
            )}
        </div>
    );
}
