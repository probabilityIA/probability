'use client';

import { useEffect, useState } from 'react';
import { createPortal } from 'react-dom';
import { Loader2 } from 'lucide-react';
import { fetchSyncFindings } from '../../infra/repository/sync-findings';
import { FindingItemsTable } from './FindingItemsTable';
import { MatchMatrixTable } from './MatchMatrixTable';
import { ChannelDataTable } from './ChannelDataTable';
import { InventoryMatrixTable } from './InventoryMatrixTable';
import { OrdersCompareTable } from './OrdersCompareTable';
import { OrdersReport } from './OrdersReport';
import { ChannelLogo } from './ChannelLogo';
import type { FindingSeverity, FindingsReport } from '../../domain/types';
import type { Integration } from '@/services/integrations/core/domain/types';
import type { IntegrationStatsItem } from '@/services/integrations/core/infra/actions/stats';
import { useSyncActivity } from '../sync-activity-context';
import { ACCENT, CARD_BORDER } from '../panel-theme';
import { INTEGRATIONS_SUBHEADER_SLOT_ID } from './PanelToolbar';

interface ReportViewProps {
    businessId: number | null;
    integrations: Integration[];
    orderSources: Integration[];
    stats: Record<number, IntegrationStatsItem>;
    statsLoaded: boolean;
}

const COUNT_TONE: Record<FindingSeverity, string> = {
    error: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300',
    warn: 'bg-amber-100 text-amber-800 dark:bg-amber-900/40 dark:text-amber-300',
    info: 'bg-sky-100 text-sky-800 dark:bg-sky-900/40 dark:text-sky-300',
};

const MATRIZ = 'matriz';

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
    logo,
    matched,
    problemas,
}: {
    name: string;
    code: string;
    logo?: string;
    matched: number;
    problemas: number;
}) {
    return (
        <div className="px-2.5 leading-tight">
            <span className="flex items-center gap-1.5 text-[10px] font-semibold text-gray-500 dark:text-gray-400">
                <ChannelLogo url={logo} code={code} size={13} />
                <span className="max-w-[9rem] truncate">{name}</span>
            </span>
            <span className="flex items-baseline gap-1">
                <span className="text-[14px] font-bold leading-none text-gray-900 dark:text-white">
                    {matched.toLocaleString('es-CO')}
                </span>
                {problemas > 0 ? (
                    <span className="text-[10px] font-bold text-amber-600 dark:text-amber-400">
                        {problemas.toLocaleString('es-CO')} por revisar
                    </span>
                ) : (
                    <span className="text-[10px] text-emerald-600 dark:text-emerald-400">sin pendientes</span>
                )}
            </span>
        </div>
    );
}

export function ReportView({ businessId, integrations, orderSources, stats, statsLoaded }: ReportViewProps) {
    const { environment } = useSyncActivity();
    const [report, setReport] = useState<FindingsReport | null>(null);
    const [loading, setLoading] = useState(true);
    const [hallazgo, setHallazgo] = useState<string>(MATRIZ);
    const [subheaderSlot, setSubheaderSlot] = useState<HTMLElement | null>(null);

    useEffect(() => {
        setSubheaderSlot(document.getElementById(INTEGRATIONS_SUBHEADER_SLOT_ID));
    }, []);

    useEffect(() => {
        const controller = new AbortController();
        setLoading(true);
        fetchSyncFindings(businessId ?? undefined, controller.signal)
            .then(setReport)
            .catch(() => setReport(null))
            .finally(() => setLoading(false));
        return () => controller.abort();
    }, [businessId]);

    useEffect(() => setHallazgo(MATRIZ), [environment]);

    const findings = report?.findings ?? [];
    const codeByChannel = Object.fromEntries(
        (report?.channels ?? []).map(c => [c.integration_name, c.channel_code]),
    );
    const logoPorCanal: Record<number, string | undefined> = {};
    for (const integracion of integrations) {
        logoPorCanal[integracion.id] = integracion.integration_type?.image_url;
    }
    const enOrdenes = environment === null;
    const enProductos = environment === 'products';
    const activo = findings.find(f => f.code === hallazgo) ?? null;

    return (
        <div className="flex h-[78vh] min-h-0 w-full flex-1 flex-col">
            {subheaderSlot && createPortal(
            <>
                {!loading && enProductos && findings.length > 0 && (
                    <div className="order-1 flex flex-wrap items-center gap-2">
                        <Pildora label="Matriz de productos" active={activo === null} onClick={() => setHallazgo(MATRIZ)} />
                        {findings.map(finding => (
                            <Pildora
                                key={finding.code}
                                label={finding.title}
                                count={finding.count}
                                tone={COUNT_TONE[finding.severity] ?? COUNT_TONE.info}
                                active={activo?.code === finding.code}
                                onClick={() => setHallazgo(finding.code)}
                            />
                        ))}
                    </div>
                )}

                {!loading && enProductos && (report?.channels?.length ?? 0) > 0 && (
                    <div className="order-3 ml-auto hidden flex-wrap items-center divide-x divide-gray-200 lg:flex dark:divide-gray-700">
                        {report?.channels.map(channel => (
                            <KpiCanal
                                key={channel.integration_id}
                                name={channel.integration_name}
                                code={channel.channel_code}
                                logo={logoPorCanal[channel.integration_id]}
                                matched={channel.matched}
                                problemas={
                                    channel.channel_no_sku + channel.sku_changed + channel.sku_typo + channel.not_associated
                                }
                            />
                        ))}
                    </div>
                )}
            </>,
            subheaderSlot,
            )}

            {loading && !enOrdenes && (
                <div className="flex items-center justify-center gap-2 py-20 text-[12px] text-gray-500">
                    <Loader2 size={14} className="animate-spin" />
                    Cruzando la informacion de tus canales
                </div>
            )}

            {(!loading || enOrdenes) && (
                <div className="flex min-h-0 flex-1 flex-col">
                    {enOrdenes && (
                        <OrdersReport integrations={orderSources} stats={stats} statsLoaded={statsLoaded} />
                    )}

                    {enProductos && activo === null && <MatchMatrixTable businessId={businessId} />}

                    {enProductos && activo !== null && (
                        <FindingItemsTable
                            code={activo.code}
                            detail={activo.detail}
                            businessId={businessId}
                            total={activo.count}
                            channels={activo.channels}
                            codeByChannel={codeByChannel}
                        />
                    )}

                    {environment === 'data' && (
                        <ChannelDataTable businessId={businessId} integrations={integrations} />
                    )}

                    {environment === 'inventory' && (
                        <InventoryMatrixTable businessId={businessId} integrations={integrations} />
                    )}

                    {environment === 'orders_compare' && (
                        <OrdersCompareTable businessId={businessId} integrations={orderSources} />
                    )}

                    {environment === 'invoicing' && (
                        <p className="py-20 text-center text-[12px] italic text-gray-400 dark:text-gray-500">
                            Facturacion desde el hub: proximamente
                        </p>
                    )}
                </div>
            )}
        </div>
    );
}
