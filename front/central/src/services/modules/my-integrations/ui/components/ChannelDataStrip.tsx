'use client';

import { useCallback, useEffect, useState } from 'react';
import { AlertTriangle, ArrowDownToLine, Loader2, RotateCcw } from 'lucide-react';
import {
    fetchDataSummary,
    undoChannelData,
    type DataApplyResult,
    type DataMode,
    type DataSummaryCell,
    type DataSummaryRow,
} from '../../infra/repository/channel-data';
import { ChannelDataApplyModal, type ObjetivoDato } from './ChannelDataApplyModal';
import { CARD_BG, CARD_BORDER } from '../panel-theme';

interface ChannelDataStripProps {
    businessId?: number;
    integrationId: number;
}

interface Disponible {
    row: DataSummaryRow;
    cell: DataSummaryCell;
}

export function ChannelDataStrip({ businessId, integrationId }: ChannelDataStripProps) {
    const [rows, setRows] = useState<DataSummaryRow[]>([]);
    const [loading, setLoading] = useState(true);
    const [objetivo, setObjetivo] = useState<ObjetivoDato | null>(null);
    const [ultimo, setUltimo] = useState<DataApplyResult | null>(null);
    const [deshaciendo, setDeshaciendo] = useState(false);

    const cargar = useCallback(async (signal?: AbortSignal) => {
        setLoading(true);
        const summary = await fetchDataSummary(businessId, signal);
        if (signal?.aborted) return;
        setRows(summary.rows);
        setLoading(false);
    }, [businessId]);

    useEffect(() => {
        const controller = new AbortController();
        void cargar(controller.signal);
        return () => controller.abort();
    }, [cargar]);

    const deshacer = async () => {
        if (!ultimo) return;
        setDeshaciendo(true);
        await undoChannelData(businessId, ultimo.batch_id);
        setDeshaciendo(false);
        setUltimo(null);
        void cargar();
    };

    const disponibles: Disponible[] = [];
    for (const row of rows) {
        const cell = row.cells.find(c => c.integration_id === integrationId);
        if (cell && (cell.can_fill > 0 || cell.can_overwrite > 0)) disponibles.push({ row, cell });
    }

    if (loading) {
        return (
            <p className="flex items-center gap-1.5 px-1 py-1 text-[11px] text-gray-400">
                <Loader2 size={11} className="animate-spin" />
                Revisando que datos puede actualizar este canal
            </p>
        );
    }

    if (disponibles.length === 0) return null;

    return (
        <div className="rounded-lg border px-2.5 py-2" style={{ borderColor: CARD_BORDER, backgroundColor: CARD_BG }}>
            <div className="mb-1.5 flex flex-wrap items-center justify-between gap-2">
                <span className="inline-flex items-center gap-1.5 text-[10.5px] font-bold uppercase tracking-wider text-gray-500">
                    <ArrowDownToLine size={11} />
                    Actualizar productos en Probability con este canal
                </span>
                {ultimo && (
                    <span className="inline-flex items-center gap-1.5 rounded-full border border-emerald-200 bg-emerald-50 px-2 py-0.5 text-[10.5px] font-semibold text-emerald-800 dark:border-emerald-500/40 dark:bg-emerald-900/20 dark:text-emerald-300">
                        Se actualizaron {ultimo.applied.toLocaleString('es-CO')} en Probability
                        <button
                            onClick={deshacer}
                            disabled={deshaciendo}
                            className="inline-flex items-center gap-1 rounded-full bg-white px-1.5 py-0.5 text-[10px] font-bold text-emerald-700 hover:bg-emerald-100 disabled:opacity-50 dark:bg-gray-800"
                        >
                            {deshaciendo ? <Loader2 size={9} className="animate-spin" /> : <RotateCcw size={9} />}
                            Deshacer
                        </button>
                    </span>
                )}
            </div>

            <div className="flex flex-wrap items-center gap-1.5">
                {disponibles.map(({ row, cell }) => (
                    <span key={row.field} className="inline-flex items-center gap-1">
                        {cell.can_fill > 0 && (
                            <button
                                onClick={() => setObjetivo({ field: row.field, label: row.label, cell, mode: 'fill_empty' as DataMode })}
                                title={`Llenar ${row.label.toLowerCase()} en Probability donde esta vacio${row.note ? `. ${row.note}` : ''}`}
                                className="inline-flex items-center gap-1 rounded-full border border-emerald-200 bg-emerald-50 px-2 py-0.5 text-[10.5px] font-bold text-emerald-700 transition-colors hover:bg-emerald-100 dark:border-emerald-500/40 dark:bg-emerald-900/20 dark:text-emerald-300"
                            >
                                {row.label} {cell.can_fill.toLocaleString('es-CO')}
                            </button>
                        )}
                        {cell.can_fill === 0 && cell.can_overwrite > 0 && (
                            <button
                                onClick={() => setObjetivo({ field: row.field, label: row.label, cell, mode: 'overwrite' as DataMode })}
                                title={`Reemplazar ${row.label.toLowerCase()} en Probability con lo que dice el canal`}
                                className="inline-flex items-center gap-1 rounded-full border border-amber-200 bg-amber-50 px-2 py-0.5 text-[10.5px] font-bold text-amber-700 transition-colors hover:bg-amber-100 dark:border-amber-500/40 dark:bg-amber-900/20 dark:text-amber-300"
                            >
                                <AlertTriangle size={9} />
                                {row.label} {cell.can_overwrite.toLocaleString('es-CO')}
                            </button>
                        )}
                    </span>
                ))}
            </div>

            {objetivo && (
                <ChannelDataApplyModal
                    objetivo={objetivo}
                    businessId={businessId ?? null}
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
