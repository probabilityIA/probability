'use client';

import { useEffect, useState } from 'react';
import { AlertTriangle, Info, Loader2, ShieldAlert, X } from 'lucide-react';
import { fetchSyncFindings } from '../../infra/repository/sync-findings';
import { FindingItemsTable } from './FindingItemsTable';
import type { Finding, FindingSeverity, FindingsReport } from '../../domain/types';

interface FindingsPanelProps {
    businessId: number | null;
    onClose: () => void;
}

const SEVERITY = {
    error: {
        label: 'Requiere atencion',
        icon: ShieldAlert,
        card: 'border-red-200 bg-red-50 dark:border-red-500/40 dark:bg-red-900/20',
        text: 'text-red-700 dark:text-red-300',
        badge: 'bg-red-600 text-white',
    },
    warn: {
        label: 'Revisar',
        icon: AlertTriangle,
        card: 'border-amber-200 bg-amber-50 dark:border-amber-500/40 dark:bg-amber-900/20',
        text: 'text-amber-800 dark:text-amber-300',
        badge: 'bg-amber-500 text-white',
    },
    info: {
        label: 'Informativo',
        icon: Info,
        card: 'border-sky-200 bg-sky-50 dark:border-sky-500/40 dark:bg-sky-900/20',
        text: 'text-sky-800 dark:text-sky-300',
        badge: 'bg-sky-500 text-white',
    },
} satisfies Record<FindingSeverity, unknown> as Record<FindingSeverity, {
    label: string;
    icon: typeof Info;
    card: string;
    text: string;
    badge: string;
}>;

function FindingTab({ finding, active, onClick }: { finding: Finding; active: boolean; onClick: () => void }) {
    const style = SEVERITY[finding.severity] ?? SEVERITY.info;
    const Icon = style.icon;
    return (
        <button
            onClick={onClick}
            className={`flex flex-shrink-0 items-center gap-1.5 rounded-lg border px-2.5 py-1.5 text-[11.5px] font-semibold transition-all ${
                active
                    ? `${style.card} ${style.text} ring-2 ring-current ring-offset-1 dark:ring-offset-gray-900`
                    : 'border-gray-200 bg-white text-gray-500 hover:bg-gray-50 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700'
            }`}
        >
            <Icon size={12} className="flex-shrink-0" />
            <span className="max-w-[11rem] truncate">{finding.title}</span>
            <span className={`rounded-full px-1.5 py-0.5 text-[10px] font-black ${style.badge}`}>{finding.count}</span>
        </button>
    );
}

function FindingCard({ finding }: { finding: Finding }) {
    const style = SEVERITY[finding.severity] ?? SEVERITY.info;
    const Icon = style.icon;
    return (
        <div className={`rounded-xl border p-3 ${style.card}`}>
            <div className="flex items-start gap-2.5">
                <Icon size={16} className={`mt-0.5 flex-shrink-0 ${style.text}`} />
                <div className="min-w-0 flex-1">
                    <div className="flex items-baseline gap-2">
                        <span className={`text-[13px] font-bold ${style.text}`}>{finding.title}</span>
                        <span className={`rounded-full px-2 py-0.5 text-[10.5px] font-bold ${style.badge}`}>
                            {finding.count}
                        </span>
                    </div>
                    <p className="mt-1 text-[11.5px] leading-snug text-gray-600 dark:text-gray-300">{finding.detail}</p>
                    {finding.channels && finding.channels.length > 0 && (
                        <div className="mt-1.5 flex flex-wrap gap-1">
                            {finding.channels.map(channel => (
                                <span
                                    key={channel}
                                    className="rounded-full border border-gray-200 bg-white/70 px-1.5 py-0.5 text-[10px] font-semibold text-gray-600 dark:border-gray-600 dark:bg-gray-800/70 dark:text-gray-300"
                                >
                                    {channel}
                                </span>
                            ))}
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}

export function FindingsPanel({ businessId, onClose }: FindingsPanelProps) {
    const [report, setReport] = useState<FindingsReport | null>(null);
    const [loading, setLoading] = useState(true);
    const [active, setActive] = useState<string | null>(null);

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
    const graves = findings.filter(f => f.severity === 'error').length;
    const activo = findings.find(f => f.code === active) ?? findings[0] ?? null;

    return (
        <div className="absolute inset-0 z-30 flex items-center justify-center p-4" role="dialog" aria-modal="true">
            <button
                aria-label="Cerrar hallazgos"
                onClick={onClose}
                className="absolute inset-0 cursor-default bg-gray-900/30 backdrop-blur-[2px]"
            />

            <div className="relative flex max-h-full h-full w-full max-w-4xl flex-col overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-2xl dark:border-gray-700 dark:bg-gray-900">
                <div className="flex items-start justify-between gap-3 border-b border-gray-100 px-4 py-3 dark:border-gray-700">
                    <div className="min-w-0">
                        <h2 className="text-sm font-bold text-gray-800 dark:text-white">Hallazgos de tus integraciones</h2>
                        <p className="mt-0.5 text-[11.5px] text-gray-500 dark:text-gray-400">
                            {loading
                                ? 'Analizando todos los canales...'
                                : findings.length === 0
                                    ? 'No encontramos cruces ni conflictos entre tus canales.'
                                    : `${findings.length} ${findings.length === 1 ? 'hallazgo' : 'hallazgos'} sobre ${report?.total ?? 0} productos, cruzando todos tus canales.`}
                        </p>
                    </div>
                    <button
                        onClick={onClose}
                        className="flex-shrink-0 rounded-lg p-1 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:hover:bg-gray-800 dark:hover:text-gray-200"
                    >
                        <X size={16} />
                    </button>
                </div>

                {loading && (
                    <div className="flex items-center justify-center gap-2 py-16 text-[12px] text-gray-500">
                        <Loader2 size={14} className="animate-spin" />
                        Cruzando la informacion de tus canales
                    </div>
                )}

                {!loading && findings.length === 0 && (
                    <div className="py-16 text-center text-[12px] text-gray-500 dark:text-gray-400">
                        Todo en orden. Vuelve a comparar productos cuando cambies algo en tus canales.
                    </div>
                )}

                {!loading && findings.length > 0 && (
                    <>
                        <div className="flex flex-wrap gap-1.5 border-b border-gray-100 px-4 py-2.5 dark:border-gray-700">
                            {findings.map(finding => (
                                <FindingTab
                                    key={finding.code}
                                    finding={finding}
                                    active={activo?.code === finding.code}
                                    onClick={() => setActive(finding.code)}
                                />
                            ))}
                        </div>

                        {activo && (
                            <div className="flex min-h-0 flex-1 flex-col gap-2.5 p-4">
                                <FindingCard finding={activo} />
                                <FindingItemsTable code={activo.code} businessId={businessId} total={activo.count} />
                            </div>
                        )}
                    </>
                )}

                {!loading && (report?.channels?.length ?? 0) > 0 && (
                    <div className="flex flex-wrap items-center gap-x-4 gap-y-1 border-t border-gray-100 bg-gray-50 px-4 py-2 dark:border-gray-700 dark:bg-gray-800/50">
                        <span className="text-[10px] font-bold uppercase tracking-wider text-gray-400">Por canal</span>
                        {report?.channels.map(channel => {
                            const problemas = channel.channel_no_sku + channel.sku_changed + channel.sku_typo + channel.not_associated;
                            return (
                                <span key={channel.integration_id} className="text-[11px] text-gray-500 dark:text-gray-400">
                                    <span className="font-semibold text-gray-700 dark:text-gray-200">{channel.integration_name}</span>
                                    {' '}{channel.matched} asociados
                                    {problemas > 0 && (
                                        <span className="ml-1 font-bold text-amber-600 dark:text-amber-400">{problemas} por revisar</span>
                                    )}
                                </span>
                            );
                        })}
                    </div>
                )}

                {!loading && graves > 0 && (
                    <div className="border-t border-gray-100 bg-red-50 px-4 py-2 text-[11.5px] font-semibold text-red-700 dark:border-gray-700 dark:bg-red-900/20 dark:text-red-300">
                        {graves === 1 ? 'Hay 1 hallazgo' : `Hay ${graves} hallazgos`} que puede estar costando ventas.
                    </div>
                )}
            </div>
        </div>
    );
}
