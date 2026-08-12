'use client';

import { useEffect, useState } from 'react';
import { AlertTriangle, Info, Loader2, ShieldAlert, X } from 'lucide-react';
import { fetchSyncFindings } from '../../infra/repository/sync-findings';
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

    return (
        <div className="absolute inset-0 z-30 flex items-center justify-center p-4" role="dialog" aria-modal="true">
            <button
                aria-label="Cerrar hallazgos"
                onClick={onClose}
                className="absolute inset-0 cursor-default bg-gray-900/30 backdrop-blur-[2px]"
            />

            <div className="relative flex max-h-full w-full max-w-2xl flex-col overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-2xl dark:border-gray-700 dark:bg-gray-900">
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

                <div className="min-h-0 flex-1 space-y-2 overflow-y-auto p-4">
                    {loading && (
                        <div className="flex items-center justify-center gap-2 py-10 text-[12px] text-gray-500">
                            <Loader2 size={14} className="animate-spin" />
                            Cruzando la informacion de tus canales
                        </div>
                    )}

                    {!loading && findings.length === 0 && (
                        <div className="py-10 text-center text-[12px] text-gray-500 dark:text-gray-400">
                            Todo en orden. Vuelve a comparar productos cuando cambies algo en tus canales.
                        </div>
                    )}

                    {!loading && findings.map(finding => <FindingCard key={finding.code} finding={finding} />)}

                    {!loading && (report?.channels?.length ?? 0) > 0 && (
                        <div className="mt-3 rounded-xl border border-gray-100 bg-gray-50 p-3 dark:border-gray-700 dark:bg-gray-800/50">
                            <p className="mb-2 text-[11px] font-bold uppercase tracking-wider text-gray-400">Por canal</p>
                            <div className="space-y-1">
                                {report?.channels.map(channel => {
                                    const problemas = channel.channel_no_sku + channel.sku_changed + channel.sku_typo + channel.not_associated;
                                    return (
                                        <div key={channel.integration_id} className="flex items-center justify-between gap-2 text-[11.5px]">
                                            <span className="truncate font-semibold text-gray-700 dark:text-gray-200">
                                                {channel.integration_name}
                                            </span>
                                            <span className="flex-shrink-0 text-gray-500 dark:text-gray-400">
                                                {channel.matched} asociados
                                                {problemas > 0 && (
                                                    <span className="ml-1.5 font-bold text-amber-600 dark:text-amber-400">
                                                        {problemas} por revisar
                                                    </span>
                                                )}
                                            </span>
                                        </div>
                                    );
                                })}
                            </div>
                        </div>
                    )}
                </div>

                {!loading && graves > 0 && (
                    <div className="border-t border-gray-100 bg-red-50 px-4 py-2 text-[11.5px] font-semibold text-red-700 dark:border-gray-700 dark:bg-red-900/20 dark:text-red-300">
                        {graves === 1 ? 'Hay 1 hallazgo' : `Hay ${graves} hallazgos`} que puede estar costando ventas.
                    </div>
                )}
            </div>
        </div>
    );
}
