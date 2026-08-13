'use client';

import { useEffect, useState } from 'react';
import { createPortal } from 'react-dom';
import { Loader2, X } from 'lucide-react';
import { fetchSyncFindings } from '../../infra/repository/sync-findings';
import { FindingItemsTable } from './FindingItemsTable';
import { MatchMatrixTable } from './MatchMatrixTable';
import { ChannelDataTable } from './ChannelDataTable';
import { InventoryMatrixTable } from './InventoryMatrixTable';
import { channelBrand } from '../../domain/types';
import type { Finding, FindingSeverity, FindingsReport } from '../../domain/types';
import type { Integration } from '@/services/integrations/core/domain/types';
import { ACCENT, CARD_BORDER } from '../panel-theme';

interface FindingsPanelProps {
    businessId: number | null;
    integrations: Integration[];
    initialTab?: string;
    onClose: () => void;
}

const COUNT_TONE: Record<FindingSeverity, string> = {
    error: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300',
    warn: 'bg-amber-100 text-amber-800 dark:bg-amber-900/40 dark:text-amber-300',
    info: 'bg-sky-100 text-sky-800 dark:bg-sky-900/40 dark:text-sky-300',
};

const RESUMEN = 'resumen';
const DATOS = 'datos';
const INVENTARIO = 'inventario';

function Pildora({
    label,
    count,
    tone,
    active,
    onClick,
}: {
    label: string;
    count?: number;
    tone?: string;
    active: boolean;
    onClick: () => void;
}) {
    return (
        <button
            onClick={onClick}
            className="inline-flex items-center gap-1.5 rounded-full border px-3 py-1.5 text-[12px] font-semibold transition-colors"
            style={
                active
                    ? { backgroundColor: ACCENT, borderColor: ACCENT, color: 'var(--color-on-primary)' }
                    : { backgroundColor: '#ffffff', borderColor: CARD_BORDER, color: '#4b5563' }
            }
        >
            {label}
            {count !== undefined && (
                <span
                    className={`rounded-full px-1.5 py-0.5 text-[10.5px] font-bold ${active ? 'bg-white/25' : tone}`}
                >
                    {count}
                </span>
            )}
        </button>
    );
}

function KpiCanal({
    name,
    code,
    matched,
    problemas,
}: {
    name: string;
    code: string;
    matched: number;
    problemas: number;
}) {
    const brand = channelBrand(code);
    return (
        <div
            className="min-w-[9.5rem] rounded-xl border px-3 py-2 dark:bg-gray-800/60"
            style={{ borderColor: CARD_BORDER, backgroundColor: '#fafafd' }}
        >
            <span className="flex items-center gap-1.5 text-[10.5px] font-semibold text-gray-500 dark:text-gray-400">
                <span className="h-2 w-2 flex-shrink-0 rounded-full" style={{ backgroundColor: brand.dot }} />
                <span className="truncate">{name}</span>
            </span>
            <div className="mt-0.5 flex items-baseline gap-1">
                <span className="text-[17px] font-black leading-none text-gray-900 dark:text-white">
                    {matched.toLocaleString('es-CO')}
                </span>
                <span className="text-[10.5px] text-gray-500 dark:text-gray-400">asociados</span>
            </div>
            {problemas > 0 ? (
                <span className="mt-0.5 block text-[10.5px] font-bold text-amber-600 dark:text-amber-400">
                    {problemas.toLocaleString('es-CO')} por revisar
                </span>
            ) : (
                <span className="mt-0.5 block text-[10.5px] text-emerald-600 dark:text-emerald-400">sin pendientes</span>
            )}
        </div>
    );
}

export function FindingsPanel({ businessId, integrations, initialTab, onClose }: FindingsPanelProps) {
    const [report, setReport] = useState<FindingsReport | null>(null);
    const [loading, setLoading] = useState(true);
    const [active, setActive] = useState<string>(initialTab ?? RESUMEN);
    const [montado, setMontado] = useState(false);

    useEffect(() => setMontado(true), []);

    useEffect(() => {
        const controller = new AbortController();
        setLoading(true);
        fetchSyncFindings(businessId ?? undefined, controller.signal)
            .then(setReport)
            .catch(() => setReport(null))
            .finally(() => setLoading(false));
        return () => controller.abort();
    }, [businessId]);

    const findings = report?.findings ?? [];
    const codeByChannel = Object.fromEntries(
        (report?.channels ?? []).map(c => [c.integration_name, c.channel_code]),
    );
    const enResumen = active === RESUMEN;
    const enDatos = active === DATOS;
    const enInventario = active === INVENTARIO;
    const activo = findings.find(f => f.code === active) ?? null;

    if (!montado) return null;

    return createPortal(
        <div className="fixed inset-0 z-[200] flex items-center justify-center p-3 sm:p-5" role="dialog" aria-modal="true">
            <button
                aria-label="Cerrar hallazgos"
                onClick={onClose}
                className="absolute inset-0 cursor-default bg-gray-900/40 backdrop-blur-[2px]"
            />

            <div className="relative flex h-[90vh] w-full max-w-[96rem] flex-col overflow-hidden rounded-2xl bg-white shadow-2xl dark:bg-gray-900">
                <div className="flex items-start justify-between gap-4 px-6 pb-3 pt-5">
                    {!loading && (
                        <div className="flex min-w-0 flex-1 flex-wrap gap-2">
                            <Pildora label="Resumen general" active={enResumen} onClick={() => setActive(RESUMEN)} />
                            <Pildora label="Actualizar productos en Probability" active={enDatos} onClick={() => setActive(DATOS)} />
                            <Pildora label="Sincronizar inventario a los canales" active={enInventario} onClick={() => setActive(INVENTARIO)} />
                            {findings.map(finding => (
                                <Pildora
                                    key={finding.code}
                                    label={finding.title}
                                    count={finding.count}
                                    tone={COUNT_TONE[finding.severity] ?? COUNT_TONE.info}
                                    active={activo?.code === finding.code}
                                    onClick={() => setActive(finding.code)}
                                />
                            ))}
                        </div>
                    )}

                    {!loading && (report?.channels?.length ?? 0) > 0 && (
                        <div className="hidden flex-wrap justify-end gap-2 lg:flex">
                            {report?.channels.map(channel => (
                                <KpiCanal
                                    key={channel.integration_id}
                                    name={channel.integration_name}
                                    code={channel.channel_code}
                                    matched={channel.matched}
                                    problemas={
                                        channel.channel_no_sku + channel.sku_changed + channel.sku_typo + channel.not_associated
                                    }
                                />
                            ))}
                        </div>
                    )}

                    <button
                        onClick={onClose}
                        aria-label="Cerrar"
                        className="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg bg-gray-100 text-gray-500 transition-colors hover:bg-gray-200 hover:text-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-gray-700"
                    >
                        <X size={15} />
                    </button>
                </div>

                {loading && (
                    <div className="flex items-center justify-center gap-2 py-20 text-[12px] text-gray-500">
                        <Loader2 size={14} className="animate-spin" />
                        Cruzando la informacion de tus canales
                    </div>
                )}

                {!loading && (
                    <>
                        {enResumen && (
                            <div className="flex min-h-0 flex-1 flex-col px-6 pb-4">
                                <MatchMatrixTable businessId={businessId} />
                            </div>
                        )}

                        {enDatos && (
                            <div className="flex min-h-0 flex-1 flex-col px-6 pb-4">
                                <ChannelDataTable businessId={businessId} />
                            </div>
                        )}

                        {enInventario && (
                            <div className="flex min-h-0 flex-1 flex-col px-6 pb-4">
                                <InventoryMatrixTable businessId={businessId} integrations={integrations} />
                            </div>
                        )}

                        {!enResumen && !enDatos && !enInventario && activo && (
                            <div className="flex min-h-0 flex-1 flex-col px-6 pb-4">
                                <FindingItemsTable
                                    code={activo.code}
                                    detail={activo.detail}
                                    businessId={businessId}
                                    total={activo.count}
                                    channels={activo.channels}
                                    codeByChannel={codeByChannel}
                                />
                            </div>
                        )}
                    </>
                )}

            </div>
        </div>,
        document.body,
    );
}
