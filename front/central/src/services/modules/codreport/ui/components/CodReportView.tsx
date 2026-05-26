'use client';

import { useMemo, useState } from 'react';
import { DollarSign, BarChart3, Package, CalendarCheck } from 'lucide-react';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { ReportFilters, RangeKey } from '../../domain/types';
import { RANGE_OPTIONS, carrierLabel } from './helpers';
import CodSummaryTab from './CodSummaryTab';
import CodOrdersTab from './CodOrdersTab';
import CodCutsTab from './CodCutsTab';

interface Props {
    selectedBusinessId?: number | null;
}

type TabKey = 'resumen' | 'ordenes' | 'cortes';

const TABS: { key: TabKey; label: string; icon: React.ReactNode }[] = [
    { key: 'resumen', label: 'Resumen', icon: <BarChart3 size={15} /> },
    { key: 'ordenes', label: 'Ordenes', icon: <Package size={15} /> },
    { key: 'cortes', label: 'Cortes de pago', icon: <CalendarCheck size={15} /> },
];

export default function CodReportView({ selectedBusinessId }: Props) {
    const { isSuperAdmin, permissions } = usePermissions();
    const isAdmin = isSuperAdmin || (permissions?.role_name || '').toLowerCase().includes('admin');

    const [tab, setTab] = useState<TabKey>('resumen');
    const [range, setRange] = useState<RangeKey>('month');
    const [customStart, setCustomStart] = useState('');
    const [customEnd, setCustomEnd] = useState('');
    const [carrier, setCarrier] = useState('');
    const carriers: string[] = [];

    const filters: ReportFilters = useMemo(() => ({
        range,
        start_date: range === 'custom' ? customStart : undefined,
        end_date: range === 'custom' ? customEnd : undefined,
        carrier: carrier || undefined,
        business_id: selectedBusinessId || undefined,
    }), [range, customStart, customEnd, carrier, selectedBusinessId]);

    return (
        <div className="flex flex-col h-full bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700">
            <div className="px-4 py-3 border-b border-gray-200 dark:border-gray-700 flex items-center gap-3 flex-wrap">
                <div className="flex items-center gap-2">
                    <DollarSign className="text-emerald-600" size={20} />
                    <h2 className="text-lg font-bold text-gray-900 dark:text-white">Recaudo contra entrega</h2>
                </div>
                <div className="flex-1" />
                <div className="flex items-center gap-1 bg-gray-100 dark:bg-gray-700 rounded-lg p-1">
                    {TABS.map(t => (
                        <button
                            key={t.key}
                            onClick={() => setTab(t.key)}
                            className={`px-3 py-1.5 rounded-md text-sm font-semibold inline-flex items-center gap-1.5 transition-colors ${
                                tab === t.key
                                    ? 'bg-white dark:bg-gray-800 shadow-sm'
                                    : 'text-gray-500 dark:text-gray-400 hover:text-gray-700'
                            }`}
                            style={tab === t.key ? { color: 'var(--color-primary)' } : undefined}
                        >
                            {t.icon} {t.label}
                        </button>
                    ))}
                </div>
            </div>

            {tab !== 'cortes' && (
                <div className="px-4 py-2.5 border-b border-gray-200 dark:border-gray-700 flex items-center gap-2 flex-wrap">
                    <div className="flex items-center gap-1 bg-gray-100 dark:bg-gray-700 rounded-lg p-1">
                        {RANGE_OPTIONS.map(opt => (
                            <button
                                key={opt.key}
                                onClick={() => setRange(opt.key)}
                                className={`px-2.5 py-1 rounded-md text-xs font-semibold transition-colors ${
                                    range === opt.key
                                        ? 'bg-white dark:bg-gray-800 shadow-sm'
                                        : 'text-gray-500 dark:text-gray-400 hover:text-gray-700'
                                }`}
                                style={range === opt.key ? { color: 'var(--color-primary)' } : undefined}
                            >
                                {opt.label}
                            </button>
                        ))}
                        <button
                            onClick={() => setRange('custom')}
                            className={`px-2.5 py-1 rounded-md text-xs font-semibold transition-colors ${
                                range === 'custom'
                                    ? 'bg-white dark:bg-gray-800 shadow-sm'
                                    : 'text-gray-500 dark:text-gray-400 hover:text-gray-700'
                            }`}
                            style={range === 'custom' ? { color: 'var(--color-primary)' } : undefined}
                        >
                            Personalizado
                        </button>
                    </div>

                    {range === 'custom' && (
                        <div className="flex items-center gap-1.5">
                            <input
                                type="date"
                                value={customStart}
                                onChange={e => setCustomStart(e.target.value)}
                                className="px-2 py-1 text-xs rounded-md border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                            />
                            <span className="text-gray-400 text-xs">a</span>
                            <input
                                type="date"
                                value={customEnd}
                                onChange={e => setCustomEnd(e.target.value)}
                                className="px-2 py-1 text-xs rounded-md border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                            />
                        </div>
                    )}

                    <div className="flex-1" />

                    <select
                        value={carrier}
                        onChange={e => setCarrier(e.target.value)}
                        className="px-2 py-1.5 text-sm rounded-md border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                    >
                        <option value="">Todas las transportadoras</option>
                        {carriers.map(c => (
                            <option key={c} value={c}>{carrierLabel(c)}</option>
                        ))}
                    </select>
                </div>
            )}

            <div className="flex-1 min-h-0 overflow-y-auto p-4">
                {tab === 'resumen' && <CodSummaryTab filters={filters} />}
                {tab === 'ordenes' && <CodOrdersTab filters={filters} />}
                {tab === 'cortes' && <CodCutsTab businessId={selectedBusinessId} isAdmin={isAdmin} />}
            </div>

        </div>
    );
}
