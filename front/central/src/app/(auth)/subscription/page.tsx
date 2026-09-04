'use client';

import { useState, useEffect, useCallback, useRef, useImperativeHandle, forwardRef } from 'react';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { Spinner, Button, Modal, Alert, Input, ConfirmModal } from '@/shared/ui';
import {
    getMySubscriptionAction,
    getMySubscriptionUsageAction,
    registerSubscriptionPaymentAction,
    editSubscriptionDatesAction,
    disableSubscriptionAction,
    reactivateSubscriptionAction,
    extendCourtesyAction,
    revertPaymentAction,
    listPaymentHistoryAction,
    listAuditLogsAction,
    listAdminBusinessesAction,
    getAdminKPIsAction,
    listSubscriptionTypesAction,
    createSubscriptionTypeAction,
    updateSubscriptionTypeAction,
    deleteSubscriptionTypeAction,
    getModuleCatalogAction,
    purchaseSubscriptionAction,
    setAutoPaymentAction,
    listOverridesAction,
    grantOverrideAction,
    revokeOverrideAction,
    listCustomPlansAction,
    createCustomPlanAction,
    updateCustomPlanAction,
    deleteCustomPlanAction,
    BusinessSubscription,
    SubscriptionUsage,
    SubscriptionType,
    BusinessModuleOverride,
    SubscriptionAuditLog,
    AdminBusinessRow,
    AdminKPIs,
    ModuleInfo,
} from '@/services/modules/wallet/infra/subscription-actions';
import { getWalletBalanceAction } from '@/services/modules/wallet/infra/actions';
import { getBoldSignatureAction } from '@/services/modules/pay/infra/actions';
import { BoldPaymentProcessingModal } from '@/app/(auth)/wallet/bold-payment-processing-modal';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';

const formatCurrency = (amount: number) =>
    new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP' }).format(amount);

const formatDate = (dateStr?: string) => {
    if (!dateStr) return '—';
    return new Date(dateStr).toLocaleDateString('es-CO', {
        day: '2-digit', month: 'long', year: 'numeric',
    });
};

function StatusBadge({ status }: { status?: string }) {
    const map: Record<string, { label: string; cls: string }> = {
        active: { label: 'Activo', cls: 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-400' },
        paid: { label: 'Activo', cls: 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-400' },
        expired: { label: 'Vencido', cls: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400' },
        cancelled: { label: 'Suspendido', cls: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400' },
        pending: { label: 'Pendiente', cls: 'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-400' },
    };
    const entry = map[status ?? ''] ?? { label: 'Sin suscripción', cls: 'bg-gray-100 text-gray-500 dark:bg-gray-700 dark:text-gray-400' };
    return (
        <span className={`text-xs px-2 py-1 rounded-full font-medium ${entry.cls}`}>
            {entry.label}
        </span>
    );
}

function PaymentStatusBadge({ status }: { status?: string }) {
    const map: Record<string, { label: string; cls: string }> = {
        paid: { label: 'Pagado', cls: 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-400' },
        reverted: { label: 'Revertido', cls: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400' },
        pending: { label: 'Pendiente', cls: 'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-400' },
        rejected: { label: 'Rechazado', cls: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400' },
    };
    const entry = map[status ?? ''] ?? { label: status || '—', cls: 'bg-gray-100 text-gray-500 dark:bg-gray-700 dark:text-gray-400' };
    return <span className={`text-xs px-2 py-1 rounded-full font-medium ${entry.cls}`}>{entry.label}</span>;
}

function paymentMethodLabel(method?: string) {
    const map: Record<string, string> = {
        WALLET: 'Wallet',
        MANUAL: 'Manual',
        TRANSFER: 'Transferencia',
        CASH: 'Efectivo',
        COURTESY: 'Cortesía',
    };
    return map[method ?? ''] || method || '—';
}

function auditDotColor(action: string) {
    const map: Record<string, string> = {
        payment_registered: '#16a34a',
        payment_reverted: '#dc2626',
        dates_edited: '#f59e0b',
        courtesy_extended: '#f59e0b',
        override_granted: '#7c3aed',
        override_revoked: '#7c3aed',
        subscription_suspended: '#dc2626',
        subscription_reactivated: '#16a34a',
    };
    return map[action] || '#9ca3af';
}

export default function SubscriptionPage() {
    const { isSuperAdmin } = usePermissions();
    const { businesses } = useBusinessesSimple();
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
    const [adminTab, setAdminTab] = useState<'businesses' | 'types' | 'custom'>('businesses');

    const selectedBusiness = businesses.find((b) => b.id === selectedBusinessId);

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
                <div>
                    <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white">
                        Suscripción
                    </h1>
                    <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                        Gestiona el acceso a la plataforma
                    </p>
                </div>

                {isSuperAdmin && businesses.length > 0 && (
                    <div className="flex items-center gap-3 bg-blue-50 dark:bg-blue-900/30 border border-blue-200 dark:border-blue-700 rounded-lg px-4 py-2">
                        <label className="text-sm font-medium text-blue-800 dark:text-blue-300 whitespace-nowrap">
                            Negocio:
                        </label>
                        <select
                            value={selectedBusinessId?.toString() ?? ''}
                            onChange={(e) => {
                                const val = e.target.value;
                                setSelectedBusinessId(val ? Number(val) : null);
                            }}
                            className="px-3 py-1.5 border border-blue-300 rounded-md text-sm bg-white dark:bg-gray-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500 min-w-[200px]"
                        >
                            <option value="">Vista Global (Admin)</option>
                            {businesses.map((b) => (
                                <option key={b.id} value={b.id}>{b.name}</option>
                            ))}
                        </select>
                    </div>
                )}
            </div>

            <div className="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
                <div className="p-6">
                    {isSuperAdmin ? (
                        selectedBusinessId ? (
                            <div className="space-y-8">
                                <BusinessSubscriptionView
                                    businessId={selectedBusinessId}
                                    businessName={selectedBusiness?.name}
                                    isSuperAdminView
                                />
                                <OverridesPanel businessId={selectedBusinessId} businessName={selectedBusiness?.name} />
                            </div>
                        ) : (
                            <div className="space-y-6">
                                <div className="flex gap-2 border-b border-gray-200 dark:border-gray-700">
                                    <button
                                        onClick={() => setAdminTab('businesses')}
                                        className={`px-4 py-2 text-sm font-medium border-b-2 ${adminTab === 'businesses' ? 'border-violet-600 text-violet-600' : 'border-transparent text-gray-500 dark:text-gray-400'}`}
                                    >
                                        Negocios
                                    </button>
                                    <button
                                        onClick={() => setAdminTab('types')}
                                        className={`px-4 py-2 text-sm font-medium border-b-2 ${adminTab === 'types' ? 'border-violet-600 text-violet-600' : 'border-transparent text-gray-500 dark:text-gray-400'}`}
                                    >
                                        Tipos de Suscripción
                                    </button>
                                    <button
                                        onClick={() => setAdminTab('custom')}
                                        className={`px-4 py-2 text-sm font-medium border-b-2 ${adminTab === 'custom' ? 'border-violet-600 text-violet-600' : 'border-transparent text-gray-500 dark:text-gray-400'}`}
                                    >
                                        Planes Personalizados
                                    </button>
                                </div>
                                {adminTab === 'businesses' && <AdminSubscriptionsView />}
                                {adminTab === 'types' && <SubscriptionTypesAdminPanel />}
                                {adminTab === 'custom' && <CustomPlansAdminPanel businesses={businesses} />}
                            </div>
                        )
                    ) : (
                        <BusinessSubscriptionView />
                    )}
                </div>
            </div>
        </div>
    );
}

function computeMembershipProgress(startDate?: string, endDate?: string) {
    if (!startDate || !endDate) return null;
    const start = new Date(startDate).getTime();
    const end = new Date(endDate).getTime();
    const now = Date.now();
    const totalMs = end - start;
    if (totalMs <= 0) return null;

    const elapsedMs = Math.min(Math.max(now - start, 0), totalMs);
    const percent = Math.round((elapsedMs / totalMs) * 100);
    const totalDays = Math.max(1, Math.round(totalMs / 86400000));
    const daysElapsed = Math.min(totalDays, Math.floor(elapsedMs / 86400000));
    const daysRemaining = Math.max(0, totalDays - daysElapsed);

    return { percent, totalDays, daysElapsed, daysRemaining, isExpired: now > end };
}

function MembershipProgress({ startDate, endDate }: { startDate?: string; endDate?: string }) {
    const progress = computeMembershipProgress(startDate, endDate);
    if (!progress) return null;

    const barColor = progress.isExpired
        ? 'bg-red-500'
        : progress.daysRemaining <= 5
            ? 'bg-amber-500'
            : 'bg-green-500';

    return (
        <div className="pt-1">
            <div className="flex items-center justify-between text-[11px] text-gray-500 dark:text-gray-400 mb-1">
                <span>Día {progress.daysElapsed} de {progress.totalDays}</span>
                <span>{progress.isExpired ? 'Vencida' : `Quedan ${progress.daysRemaining} días`}</span>
            </div>
            <div className="w-full h-2 rounded-full bg-gray-100 dark:bg-gray-600 overflow-hidden">
                <div className={`h-full rounded-full transition-all ${barColor}`} style={{ width: `${progress.percent}%` }} />
            </div>
            <div className="flex items-center justify-between text-[10px] text-gray-400 mt-1">
                <span>{formatDate(startDate)}</span>
                <span>{formatDate(endDate)}</span>
            </div>
        </div>
    );
}

const ADMIN_STATUS_FILTERS: Array<{ value: string; label: string }> = [
    { value: '', label: 'Todos' },
    { value: 'active', label: 'Activos' },
    { value: 'expiring_soon', label: 'Por vencer' },
    { value: 'expired', label: 'Vencidos' },
    { value: 'cancelled', label: 'Suspendidos' },
    { value: 'no_plan', label: 'Sin plan' },
];

const ROW_TONE: Record<string, { bg: string; fg: string; bd: string; dot: string; label: string }> = {
    active: { bg: '#f0fdf4', fg: '#15803d', bd: '#bbf7d0', dot: '#16a34a', label: 'Activo' },
    soon: { bg: '#fffbeb', fg: '#b45309', bd: '#fde68a', dot: '#f59e0b', label: 'Por vencer' },
    expired: { bg: '#fef2f2', fg: '#b91c1c', bd: '#fecaca', dot: '#dc2626', label: 'Vencido' },
    suspended: { bg: '#f4f4f5', fg: '#3f3f46', bd: '#dcdce0', dot: '#71717a', label: 'Suspendido' },
    none: { bg: '#f8f9fa', fg: '#6b7280', bd: '#e4e6eb', dot: '#9ca3af', label: 'Sin plan' },
};

function getRowTone(row: AdminBusinessRow): keyof typeof ROW_TONE {
    if (row.status === 'cancelled') return 'suspended';
    if (!row.plan_name) return 'none';
    if (row.status === 'expired') return 'expired';
    if (row.status === 'active' && row.cycle_end_date) {
        const daysLeft = Math.ceil((new Date(row.cycle_end_date).getTime() - Date.now()) / 86400000);
        if (daysLeft <= 5) return 'soon';
    }
    return 'active';
}

const PLAN_TONE: Record<string, { bg: string; fg: string; bd: string }> = {
    premium: { bg: '#f3f0ff', fg: '#6d28d9', bd: '#e6dfff' },
    pro: { bg: '#f5f2ff', fg: '#7c3aed', bd: '#ebe4ff' },
    basico: { bg: '#f4f5f7', fg: '#4b5563', bd: '#e4e6eb' },
};

function getPlanTone(planName?: string) {
    if (!planName) return { bg: '#f8f9fa', fg: '#9ca3af', bd: '#e4e6eb' };
    const key = planName.toLowerCase().normalize('NFD').replace(/[^a-z]/g, '');
    return PLAN_TONE[key] || { bg: '#f3f0ff', fg: '#6d28d9', bd: '#e6dfff' };
}

function formatDueRel(row: AdminBusinessRow): string {
    if (!row.cycle_end_date) return 'sin plan asignado';
    const days = Math.round((new Date(row.cycle_end_date).getTime() - Date.now()) / 86400000);
    if (row.status === 'cancelled') return `suspendida hace ${Math.abs(days)} días`;
    if (days < 0) return `vencida ${Math.abs(days)} días`;
    if (days === 0) return 'vence hoy';
    return `en ${days} días`;
}

const ADMIN_ROW_GRID = '36px minmax(190px,240px) 160px 130px 170px 140px minmax(140px,1fr) minmax(140px,1fr) 44px';

function KpiCard({ label, value, note }: { label: string; value: string; note?: string }) {
    return (
        <div className="bg-white dark:bg-gray-700 rounded-xl border border-gray-200 dark:border-gray-600 p-4">
            <p className="text-xs text-gray-500 dark:text-gray-400">{label}</p>
            <p className="text-xl font-semibold text-gray-900 dark:text-white mt-1">{value}</p>
            {note && <p className="text-xs text-gray-400 mt-0.5">{note}</p>}
        </div>
    );
}

function AdminSubscriptionsView() {
    const [kpis, setKpis] = useState<AdminKPIs | null>(null);
    const [rows, setRows] = useState<AdminBusinessRow[]>([]);
    const [total, setTotal] = useState(0);
    const [page, setPage] = useState(1);
    const [pageSize] = useState(10);
    const [search, setSearch] = useState('');
    const [status, setStatus] = useState('');
    const [loading, setLoading] = useState(false);
    const [activeRow, setActiveRow] = useState<AdminBusinessRow | null>(null);
    const [selectedIds, setSelectedIds] = useState<Set<number>>(new Set());
    const [subscriptionTypes, setSubscriptionTypes] = useState<SubscriptionType[]>([]);
    const [bulkBusy, setBulkBusy] = useState(false);
    const [bulkMessage, setBulkMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

    const [bulkPagoOpen, setBulkPagoOpen] = useState(false);
    const [bulkTypeId, setBulkTypeId] = useState('');
    const [bulkMonths, setBulkMonths] = useState('1');
    const [bulkMethod, setBulkMethod] = useState('TRANSFER');
    const [bulkRef, setBulkRef] = useState('');
    const [bulkNotes, setBulkNotes] = useState('');
    const [bulkExtendOpen, setBulkExtendOpen] = useState(false);
    const [bulkDays, setBulkDays] = useState(7);
    const [bulkReason, setBulkReason] = useState('');
    const [bulkSuspendOpen, setBulkSuspendOpen] = useState(false);

    const refreshKpis = useCallback(() => {
        getAdminKPIsAction().then((res) => {
            if (res.success && res.data) setKpis(res.data);
        });
    }, []);

    const refreshRows = useCallback(async () => {
        setLoading(true);
        const res = await listAdminBusinessesAction({ page, pageSize, search: search.trim() || undefined, status: status || undefined });
        setLoading(false);
        if (res.success) {
            setRows(res.data || []);
            setTotal(res.total || 0);
        }
    }, [page, pageSize, search, status]);

    useEffect(() => { refreshKpis(); }, [refreshKpis]);
    useEffect(() => { refreshRows(); }, [refreshRows]);
    useEffect(() => {
        listSubscriptionTypesAction(true).then((res) => { if (res.success && res.data) setSubscriptionTypes(res.data); });
    }, []);

    const totalPages = Math.max(1, Math.ceil(total / pageSize));
    const allChecked = rows.length > 0 && rows.every((r) => selectedIds.has(r.id));

    const toggleAll = () => {
        setSelectedIds(allChecked ? new Set() : new Set(rows.map((r) => r.id)));
    };
    const toggleOne = (id: number) => {
        setSelectedIds((prev) => {
            const next = new Set(prev);
            if (next.has(id)) next.delete(id); else next.add(id);
            return next;
        });
    };

    const runBulk = async (fn: (id: number) => Promise<{ success: boolean; error?: string }>, successText: string) => {
        setBulkBusy(true);
        let okCount = 0;
        let failCount = 0;
        for (const id of selectedIds) {
            const res = await fn(id);
            if (res.success) okCount++; else failCount++;
        }
        setBulkBusy(false);
        setSelectedIds(new Set());
        await refreshRows();
        await refreshKpis();
        setBulkMessage(failCount
            ? { type: 'error', text: `${okCount} aplicados, ${failCount} con error.` }
            : { type: 'success', text: successText });
    };

    const handleBulkRegisterPayment = async () => {
        if (!bulkTypeId) return;
        await runBulk((id) => registerSubscriptionPaymentAction({
            businessId: id,
            subscriptionTypeId: Number(bulkTypeId),
            monthsToAdd: Number(bulkMonths),
            paymentMethod: bulkMethod,
            paymentReference: bulkRef || undefined,
            notes: bulkNotes || undefined,
        }), 'Pagos registrados.');
        setBulkPagoOpen(false); setBulkTypeId(''); setBulkMonths('1'); setBulkRef(''); setBulkNotes('');
    };

    const handleBulkExtend = async () => {
        if (bulkDays <= 0 || !bulkReason.trim()) return;
        await runBulk((id) => extendCourtesyAction({ businessId: id, days: bulkDays, reason: bulkReason.trim() }), 'Días de cortesía aplicados.');
        setBulkExtendOpen(false); setBulkReason(''); setBulkDays(7);
    };

    const handleBulkSuspend = async () => {
        await runBulk((id) => disableSubscriptionAction(id), 'Negocios suspendidos.');
        setBulkSuspendOpen(false);
    };

    const pageNumbers = Array.from({ length: totalPages }, (_, i) => i + 1).filter((n) => Math.abs(n - page) <= 1 || n === 1 || n === totalPages);

    return (
        <div className="space-y-6">
            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
                <KpiCard label="Activos" value={String(kpis?.active_count ?? '—')} />
                <KpiCard label="Por vencer (5 días)" value={String(kpis?.expiring_soon_count ?? '—')} />
                <KpiCard label="Vencidos / suspendidos" value={String(kpis?.expired_or_suspended_count ?? '—')} />
                <KpiCard label="MRR" value={kpis ? formatCurrency(kpis.mrr) : '—'} note="por mes" />
            </div>

            {bulkMessage && <Alert type={bulkMessage.type} onClose={() => setBulkMessage(null)}>{bulkMessage.text}</Alert>}

            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl overflow-hidden">
                <div className="px-4 py-3 border-b border-gray-100 dark:border-gray-700 flex items-center gap-3 flex-wrap">
                    <input
                        type="text"
                        value={search}
                        onChange={(e) => { setSearch(e.target.value); setPage(1); }}
                        placeholder="Buscar negocio, ID o plan..."
                        className="flex-none w-[280px] px-3 py-1.5 border border-gray-200 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white dark:border-gray-600"
                    />
                    <div className="flex items-center gap-1 p-0.5 bg-gray-100 dark:bg-gray-900/40 rounded-lg">
                        {ADMIN_STATUS_FILTERS.map((f) => (
                            <button
                                key={f.value}
                                onClick={() => { setStatus(f.value); setPage(1); }}
                                className={`h-7 px-2.5 rounded-md text-xs font-medium ${status === f.value ? 'bg-white dark:bg-gray-700 text-gray-900 dark:text-white shadow-sm' : 'text-gray-500 dark:text-gray-400'}`}
                            >
                                {f.label}
                            </button>
                        ))}
                    </div>
                    <span className="ml-auto text-xs text-gray-500 dark:text-gray-400">{loading ? 'Cargando...' : `${total} negocios`}</span>
                </div>

                {selectedIds.size > 0 && (
                    <div className="px-4 py-2.5 bg-violet-50 dark:bg-violet-900/20 border-b border-violet-100 dark:border-violet-800 flex items-center gap-3">
                        <span className="text-sm font-semibold text-violet-700 dark:text-violet-300">
                            {selectedIds.size} {selectedIds.size === 1 ? 'negocio seleccionado' : 'negocios seleccionados'}
                        </span>
                        <div className="w-px h-4 bg-violet-200 dark:bg-violet-700" />
                        <button onClick={() => setBulkPagoOpen(true)} className="h-7 px-2.5 rounded-md border border-green-200 bg-white dark:bg-gray-800 text-green-700 dark:text-green-400 text-xs font-semibold hover:bg-green-50">Registrar pago</button>
                        <button onClick={() => setBulkExtendOpen(true)} className="h-7 px-2.5 rounded-md border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 text-xs font-medium hover:bg-gray-50 dark:hover:bg-gray-700">Extender días</button>
                        <button onClick={() => setBulkSuspendOpen(true)} className="h-7 px-2.5 rounded-md border border-red-200 bg-white dark:bg-gray-800 text-red-700 dark:text-red-400 text-xs font-semibold hover:bg-red-50">Suspender</button>
                        <button onClick={() => setSelectedIds(new Set())} className="ml-auto text-xs text-gray-500 dark:text-gray-400 hover:underline">Limpiar</button>
                    </div>
                )}

                <div className="overflow-x-auto">
                <div className="min-w-[1150px]">
                <div style={{ display: 'grid', gridTemplateColumns: ADMIN_ROW_GRID }} className="items-center px-4 h-9 bg-gray-50 dark:bg-gray-900/40 border-b border-gray-100 dark:border-gray-700 text-[11px] font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">
                    <div><input type="checkbox" checked={allChecked} onChange={toggleAll} className="w-3.5 h-3.5 accent-violet-600 cursor-pointer" /></div>
                    <div>Negocio</div>
                    <div>Plan</div>
                    <div>Estado</div>
                    <div>Ciclo</div>
                    <div>Vence</div>
                    <div className="text-right">Último pago</div>
                    <div className="text-right">Pago pronosticado</div>
                    <div />
                </div>

                {loading ? (
                    <div className="px-4 py-10 text-center text-sm text-gray-400">Cargando...</div>
                ) : rows.length === 0 ? (
                    <div className="py-14 px-6 flex flex-col items-center gap-2 text-center">
                        <p className="text-sm font-semibold text-gray-900 dark:text-white">Ningún negocio coincide con el filtro</p>
                        <button onClick={() => { setSearch(''); setStatus(''); setPage(1); }} className="mt-1 h-8 px-3 rounded-lg border border-gray-200 dark:border-gray-600 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700">Limpiar filtros</button>
                    </div>
                ) : (
                    rows.map((row) => {
                        const tone = ROW_TONE[getRowTone(row)];
                        const planTone = getPlanTone(row.plan_name);
                        const progress = computeMembershipProgress(row.cycle_start_date, row.cycle_end_date);
                        const dueFg = row.status === 'expired' || row.status === 'cancelled' ? '#b91c1c' : (progress && !progress.isExpired && progress.daysRemaining <= 5 ? '#b45309' : '#9ca3af');
                        return (
                            <div
                                key={row.id}
                                onClick={() => toggleOne(row.id)}
                                style={{ display: 'grid', gridTemplateColumns: ADMIN_ROW_GRID, background: selectedIds.has(row.id) ? '#faf8ff' : undefined }}
                                className="items-center px-4 h-16 border-b border-gray-100 dark:border-gray-700 cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700/40"
                            >
                                <div onClick={(e) => e.stopPropagation()}>
                                    <input type="checkbox" checked={selectedIds.has(row.id)} onChange={() => toggleOne(row.id)} className="w-3.5 h-3.5 accent-violet-600 cursor-pointer" />
                                </div>
                                <div className="min-w-0 pr-3">
                                    <p className="text-[13.5px] font-semibold text-gray-900 dark:text-white truncate">{row.name}</p>
                                    <p className="text-[11px] font-mono text-gray-400 truncate">{row.code}</p>
                                </div>
                                <div>
                                    <span className="inline-flex items-center h-[22px] px-2.5 rounded-full text-xs font-semibold" style={{ background: planTone.bg, color: planTone.fg, border: `1px solid ${planTone.bd}` }}>
                                        {row.plan_name || '—'}
                                    </span>
                                </div>
                                <div>
                                    <span className="inline-flex items-center gap-1.5 h-[22px] pl-2 pr-2.5 rounded-full text-xs font-semibold" style={{ background: tone.bg, color: tone.fg, border: `1px solid ${tone.bd}` }}>
                                        <span className="w-1.5 h-1.5 rounded-full" style={{ background: tone.dot }} />
                                        {tone.label}
                                    </span>
                                </div>
                                <div className="flex flex-col gap-1.5 pr-3">
                                    {row.cycle_start_date && row.cycle_end_date ? (
                                        <>
                                            <div className="w-[130px] h-[5px] rounded-full bg-gray-100 dark:bg-gray-600 overflow-hidden">
                                                <div className="h-full rounded-full" style={{ width: `${progress?.percent ?? 0}%`, background: tone.dot }} />
                                            </div>
                                            <span className="text-[11.5px] text-gray-500 dark:text-gray-400">{progress ? `Día ${progress.daysElapsed} de ${progress.totalDays}` : 'Sin ciclo'}</span>
                                        </>
                                    ) : (
                                        <span className="text-[11.5px] text-gray-400">Sin ciclo activo</span>
                                    )}
                                </div>
                                <div className="flex flex-col gap-0.5 pr-3">
                                    <span className="text-[12.5px] text-gray-700 dark:text-gray-300">{row.cycle_end_date ? formatDate(row.cycle_end_date) : '—'}</span>
                                    <span className="text-[11.5px] font-medium" style={{ color: dueFg }}>{formatDueRel(row)}</span>
                                </div>
                                <div className="flex flex-col gap-0.5 text-right">
                                    {row.last_payment_amount ? (
                                        <>
                                            <span className="font-mono text-[12.5px] text-gray-900 dark:text-white">{formatCurrency(row.last_payment_amount)}</span>
                                            <span className="text-[11.5px] text-gray-400">{formatDate(row.last_payment_date)}</span>
                                        </>
                                    ) : <span className="text-[11.5px] text-gray-400">—</span>}
                                </div>
                                <div className="flex flex-col gap-0.5 text-right">
                                    {row.forecasted_payment ? (
                                        (() => {
                                            const hasOverage = !!row.last_payment_amount && row.forecasted_payment > row.last_payment_amount;
                                            return (
                                                <>
                                                    <span className={`font-mono text-[12.5px] ${hasOverage ? 'text-amber-600 dark:text-amber-400 font-semibold' : 'text-gray-900 dark:text-white'}`}>
                                                        {formatCurrency(row.forecasted_payment)}
                                                    </span>
                                                    {hasOverage && <span className="text-[11.5px] text-amber-500 dark:text-amber-400">incluye excedentes</span>}
                                                </>
                                            );
                                        })()
                                    ) : <span className="text-[11.5px] text-gray-400">—</span>}
                                </div>
                                <div className="flex justify-end">
                                    <button
                                        onClick={(e) => { e.stopPropagation(); setActiveRow(row); }}
                                        title="Ver detalle"
                                        className="w-7 h-7 rounded-md flex items-center justify-center text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-600 hover:text-gray-700 dark:hover:text-gray-200"
                                    >
                                        ›
                                    </button>
                                </div>
                            </div>
                        );
                    })
                )}
                </div>
                </div>

                {rows.length > 0 && (
                    <div className="px-4 py-3 flex items-center justify-between">
                        <span className="text-[12.5px] text-gray-500 dark:text-gray-400">
                            Mostrando {(page - 1) * pageSize + 1}–{Math.min(page * pageSize, total)} de {total}
                        </span>
                        <div className="flex items-center gap-1.5">
                            <button onClick={() => setPage((p) => Math.max(1, p - 1))} disabled={page === 1} className="h-7.5 px-2.5 rounded-md border border-gray-200 dark:border-gray-600 text-xs text-gray-700 dark:text-gray-300 disabled:text-gray-300 disabled:cursor-not-allowed">← Anterior</button>
                            {pageNumbers.map((n, i) => (
                                <span key={n} className="flex items-center">
                                    {i > 0 && pageNumbers[i - 1] !== n - 1 && <span className="text-xs text-gray-300 px-1">…</span>}
                                    <button
                                        onClick={() => setPage(n)}
                                        className={`h-7.5 w-7.5 rounded-md text-xs font-medium ${n === page ? 'bg-violet-600 text-white' : 'border border-gray-200 dark:border-gray-600 text-gray-700 dark:text-gray-300'}`}
                                    >
                                        {n}
                                    </button>
                                </span>
                            ))}
                            <button onClick={() => setPage((p) => Math.min(totalPages, p + 1))} disabled={page === totalPages} className="h-7.5 px-2.5 rounded-md border border-gray-200 dark:border-gray-600 text-xs text-gray-700 dark:text-gray-300 disabled:text-gray-300 disabled:cursor-not-allowed">Siguiente →</button>
                        </div>
                    </div>
                )}
            </div>

            {activeRow && (
                <BusinessDetailDrawer
                    business={activeRow}
                    onClose={() => setActiveRow(null)}
                    onChanged={() => { refreshRows(); refreshKpis(); }}
                />
            )}

            <Modal isOpen={bulkPagoOpen} onClose={() => setBulkPagoOpen(false)} title={`Registrar pago (${selectedIds.size} negocios)`} size="md">
                <div className="space-y-4 p-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Plan</label>
                        <select value={bulkTypeId} onChange={(e) => setBulkTypeId(e.target.value)} className="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white">
                            <option value="">Selecciona un tipo</option>
                            {subscriptionTypes.map((t) => (
                                <option key={t.id} value={t.id}>{t.name} — {formatCurrency(t.price)}/{t.billing_period === 'monthly' ? 'mes' : 'año'}</option>
                            ))}
                        </select>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Meses</label>
                        <select value={bulkMonths} onChange={(e) => setBulkMonths(e.target.value)} className="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white">
                            <option value="1">1 mes</option>
                            <option value="3">3 meses</option>
                            <option value="6">6 meses</option>
                            <option value="12">12 meses</option>
                        </select>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Método</label>
                        <select value={bulkMethod} onChange={(e) => setBulkMethod(e.target.value)} className="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white">
                            <option value="TRANSFER">Transferencia</option>
                            <option value="CASH">Efectivo</option>
                            <option value="WALLET">Wallet</option>
                        </select>
                    </div>
                    <Input label="Referencia (opcional)" value={bulkRef} onChange={(e) => setBulkRef(e.target.value)} />
                    <Input label="Notas internas (opcional)" value={bulkNotes} onChange={(e) => setBulkNotes(e.target.value)} />
                    <div className="flex justify-end gap-2 pt-2">
                        <Button variant="secondary" onClick={() => setBulkPagoOpen(false)}>Cancelar</Button>
                        <Button variant="success" onClick={handleBulkRegisterPayment} loading={bulkBusy} disabled={!bulkTypeId}>Registrar pago</Button>
                    </div>
                </div>
            </Modal>

            <Modal isOpen={bulkExtendOpen} onClose={() => setBulkExtendOpen(false)} title={`Extender días (${selectedIds.size} negocios)`} size="sm">
                <div className="space-y-4 p-4">
                    <p className="text-sm text-gray-500 dark:text-gray-400">Sin registrar pago. Se marca como cortesía en cada negocio.</p>
                    <div className="flex gap-2">
                        {[3, 7, 15, 30].map((d) => (
                            <button
                                key={d}
                                onClick={() => setBulkDays(d)}
                                className={`px-3 py-1.5 rounded-lg text-sm border ${bulkDays === d ? 'bg-amber-500 text-white border-amber-500' : 'border-gray-200 dark:border-gray-600 text-gray-700 dark:text-gray-300'}`}
                            >
                                +{d}
                            </button>
                        ))}
                    </div>
                    <Input label="Motivo (obligatorio)" value={bulkReason} onChange={(e) => setBulkReason(e.target.value)} />
                    <div className="flex justify-end gap-2 pt-2">
                        <Button variant="secondary" onClick={() => setBulkExtendOpen(false)}>Cancelar</Button>
                        <Button variant="quaternary" onClick={handleBulkExtend} loading={bulkBusy} disabled={!bulkReason.trim()}>Extender</Button>
                    </div>
                </div>
            </Modal>

            <ConfirmModal
                isOpen={bulkSuspendOpen}
                onClose={() => setBulkSuspendOpen(false)}
                onConfirm={handleBulkSuspend}
                title={`¿Suspender ${selectedIds.size} negocios?`}
                message="Los usuarios de estos negocios perderán acceso al panel de inmediato. Se pueden reactivar después sin perder datos."
                confirmText="Sí, suspender"
                type="danger"
            />
        </div>
    );
}

function buildMonthGrid(year: number, month: number) {
    const first = new Date(year, month, 1);
    const startOffset = (first.getDay() + 6) % 7;
    const daysInMonth = new Date(year, month + 1, 0).getDate();
    const cells: (number | null)[] = [];
    for (let i = 0; i < startOffset; i++) cells.push(null);
    for (let d = 1; d <= daysInMonth; d++) cells.push(d);
    while (cells.length % 7 !== 0) cells.push(null);
    return cells;
}

const sameDay = (a: Date, b: Date) => a.getFullYear() === b.getFullYear() && a.getMonth() === b.getMonth() && a.getDate() === b.getDate();

const fmtDateInput = (d: Date) => `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;

function SubscriptionCalendar({ cycleStartDate, cycleEndDate, cutoffDay, onSave, saving }: {
    cycleStartDate?: string;
    cycleEndDate?: string;
    cutoffDay?: number;
    onSave: (startDate: string, endDate: string, cutoffDay?: number) => Promise<void>;
    saving?: boolean;
}) {
    const start = cycleStartDate ? new Date(cycleStartDate) : null;
    const end = cycleEndDate ? new Date(cycleEndDate) : null;
    const [viewDate, setViewDate] = useState(() => end ? new Date(end.getFullYear(), end.getMonth(), 1) : new Date(new Date().getFullYear(), new Date().getMonth(), 1));
    const [editing, setEditing] = useState(false);
    const [pendingStart, setPendingStart] = useState<Date | null>(null);
    const [pendingEnd, setPendingEnd] = useState<Date | null>(null);
    const [pendingCutoffDay, setPendingCutoffDay] = useState<string>(cutoffDay ? String(cutoffDay) : '');

    const year = viewDate.getFullYear();
    const month = viewDate.getMonth();
    const cells = buildMonthGrid(year, month);
    const today = new Date();

    const inCycle = (day: number) => {
        if (!start || !end) return false;
        const d = new Date(year, month, day);
        return d >= new Date(start.getFullYear(), start.getMonth(), start.getDate()) && d <= new Date(end.getFullYear(), end.getMonth(), end.getDate());
    };

    const inPending = (day: number) => {
        if (!pendingStart || !pendingEnd) return false;
        const d = new Date(year, month, day);
        return d >= pendingStart && d <= pendingEnd;
    };

    const startEditing = () => {
        setPendingStart(start);
        setPendingEnd(end);
        setPendingCutoffDay(cutoffDay ? String(cutoffDay) : '');
        setEditing(true);
    };

    const cancelEditing = () => {
        setEditing(false);
        setPendingStart(null);
        setPendingEnd(null);
        setPendingCutoffDay(cutoffDay ? String(cutoffDay) : '');
    };

    const handleDayClick = (day: number) => {
        if (!editing) return;
        const clicked = new Date(year, month, day);
        if (!pendingStart || pendingEnd) {
            setPendingStart(clicked);
            setPendingEnd(null);
        } else if (clicked < pendingStart) {
            setPendingStart(clicked);
            setPendingEnd(null);
        } else {
            setPendingEnd(clicked);
        }
    };

    const handleSave = async () => {
        if (!pendingStart || !pendingEnd) return;
        const parsedCutoffDay = pendingCutoffDay ? Number(pendingCutoffDay) : undefined;
        await onSave(fmtDateInput(pendingStart), fmtDateInput(pendingEnd), parsedCutoffDay);
        cancelEditing();
    };

    return (
        <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-3.5">
            <div className="flex items-center justify-between mb-2.5">
                <button onClick={() => setViewDate(new Date(year, month - 1, 1))} className="w-6 h-6 flex items-center justify-center rounded text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700">‹</button>
                <span className="text-xs font-semibold text-gray-700 dark:text-gray-200 capitalize">{viewDate.toLocaleDateString('es-CO', { month: 'long', year: 'numeric' })}</span>
                <button onClick={() => setViewDate(new Date(year, month + 1, 1))} className="w-6 h-6 flex items-center justify-center rounded text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700">›</button>
            </div>
            <div className="grid grid-cols-7 gap-1 text-center text-[10px] text-gray-400 mb-1">
                {['L', 'M', 'X', 'J', 'V', 'S', 'D'].map((d, i) => <span key={i}>{d}</span>)}
            </div>
            <div className="grid grid-cols-7 gap-1">
                {cells.map((day, i) => {
                    if (!day) return <div key={i} className="aspect-square" />;
                    const isToday = sameDay(new Date(year, month, day), today);
                    const isEnd = end && sameDay(new Date(year, month, day), end);
                    const isStart = start && sameDay(new Date(year, month, day), start);
                    const highlighted = inCycle(day);
                    const isPendingBoundary = Boolean(pendingStart && sameDay(new Date(year, month, day), pendingStart)) || Boolean(pendingEnd && sameDay(new Date(year, month, day), pendingEnd));
                    const pendingHighlighted = inPending(day);
                    return (
                        <button
                            key={i}
                            type="button"
                            onClick={() => handleDayClick(day)}
                            disabled={!editing}
                            className={`aspect-square flex items-center justify-center rounded-md text-[11px] ${editing ? 'cursor-pointer' : 'cursor-default'} ${isToday ? 'ring-1 ring-violet-600' : ''} ${
                                editing
                                    ? (isPendingBoundary
                                        ? 'bg-amber-500 text-white font-semibold'
                                        : pendingHighlighted
                                            ? 'bg-amber-100 dark:bg-amber-900/40 text-amber-700 dark:text-amber-300'
                                            : 'text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700')
                                    : (isEnd || isStart
                                        ? 'bg-violet-600 text-white font-semibold'
                                        : highlighted
                                            ? 'bg-violet-100 dark:bg-violet-900/40 text-violet-700 dark:text-violet-300'
                                            : 'text-gray-600 dark:text-gray-300')
                            }`}
                        >
                            {day}
                        </button>
                    );
                })}
            </div>

            {editing ? (
                <div className="mt-3 space-y-2">
                    <p className="text-[11px] text-gray-500 dark:text-gray-400">
                        {!pendingStart ? 'Selecciona la fecha de inicio.' : !pendingEnd ? 'Ahora selecciona la fecha de vencimiento.' : `${fmtDateInput(pendingStart)} → ${fmtDateInput(pendingEnd)}`}
                    </p>
                    <div>
                        <label className="block text-[10px] text-gray-400 mb-1">
                            {'D\u00eda de corte (1-31, opcional) \u2014 suspende la cuenta ese d\u00eda del mes si no ha pagado, en vez de al vencer el ciclo'}
                        </label>
                        <input
                            type="number"
                            min={1}
                            max={31}
                            value={pendingCutoffDay}
                            onChange={(e) => setPendingCutoffDay(e.target.value)}
                            placeholder={'Sin d\u00eda de corte'}
                            className="w-full text-[11px] px-2 py-1.5 rounded-md border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-200"
                        />
                    </div>
                    <div className="flex gap-2">
                        <Button size="sm" variant="secondary" className="flex-1 !text-xs" onClick={cancelEditing}>Cancelar</Button>
                        <Button size="sm" variant="purple" className="flex-1 !text-xs" onClick={handleSave} loading={saving} disabled={!pendingStart || !pendingEnd}>Guardar</Button>
                    </div>
                </div>
            ) : (
                <div className="mt-2.5 flex items-center justify-between">
                    {(!start || !end) ? <p className="text-[11px] text-gray-400">Sin ciclo asignado</p> : (
                        <p className="text-[11px] text-gray-400">
                            {cutoffDay ? `Corte: d\u00eda ${cutoffDay}` : 'Sin d\u00eda de corte'}
                        </p>
                    )}
                    <button onClick={startEditing} className="text-[11px] font-semibold text-violet-600 hover:underline">Editar fechas</button>
                </div>
            )}
        </div>
    );
}

function BusinessDetailDrawer({ business, onClose, onChanged }: {
    business: AdminBusinessRow;
    onClose: () => void;
    onChanged: () => void;
}) {
    const [subscription, setSubscription] = useState<BusinessSubscription | null>(null);
    const [history, setHistory] = useState<BusinessSubscription[]>([]);
    const [overrides, setOverrides] = useState<BusinessModuleOverride[]>([]);
    const [auditLogs, setAuditLogs] = useState<SubscriptionAuditLog[]>([]);
    const [subscriptionTypes, setSubscriptionTypes] = useState<SubscriptionType[]>([]);
    const [customPlans, setCustomPlans] = useState<SubscriptionType[]>([]);
    const [moduleCatalog, setModuleCatalog] = useState<ModuleInfo[]>([]);
    const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
    const [busy, setBusy] = useState(false);

    const [registerOpen, setRegisterOpen] = useState(false);
    const [selectedTypeId, setSelectedTypeId] = useState('');
    const [months, setMonths] = useState('1');
    const [paymentMethod, setPaymentMethod] = useState('TRANSFER');
    const [payRef, setPayRef] = useState('');
    const [notes, setNotes] = useState('');

    const [extendOpen, setExtendOpen] = useState(false);
    const [extendDays, setExtendDays] = useState(7);
    const [extendReason, setExtendReason] = useState('');

    const [suspendConfirmOpen, setSuspendConfirmOpen] = useState(false);

    const [grantOpen, setGrantOpen] = useState(false);
    const [grantModule, setGrantModule] = useState('');
    const [grantExpiresAt, setGrantExpiresAt] = useState('');
    const [grantNotes, setGrantNotes] = useState('');

    const load = useCallback(async () => {
        const [subRes, historyRes, overridesRes, auditRes] = await Promise.all([
            getMySubscriptionAction(business.id),
            listPaymentHistoryAction(business.id),
            listOverridesAction(business.id),
            listAuditLogsAction(business.id),
        ]);
        if (subRes.success) setSubscription(subRes.data || null);
        if (historyRes.success) setHistory(historyRes.data || []);
        if (overridesRes.success) setOverrides(overridesRes.data || []);
        if (auditRes.success) setAuditLogs(auditRes.data || []);
    }, [business.id]);

    useEffect(() => { load(); }, [load]);
    useEffect(() => {
        listSubscriptionTypesAction(true).then((res) => { if (res.success && res.data) setSubscriptionTypes(res.data); });
        listCustomPlansAction(business.id).then((res) => { if (res.success && res.data) setCustomPlans(res.data); });
        getModuleCatalogAction().then((res) => { if (res.success && res.data) setModuleCatalog(res.data); });
    }, []);

    const afterChange = async (successText: string) => {
        setMessage({ type: 'success', text: successText });
        await load();
        onChanged();
    };

    const handleRegisterPayment = async () => {
        if (!selectedTypeId) return;
        setBusy(true);
        const res = await registerSubscriptionPaymentAction({
            businessId: business.id,
            subscriptionTypeId: Number(selectedTypeId),
            monthsToAdd: Number(months),
            paymentMethod,
            paymentReference: payRef || undefined,
            notes: notes || undefined,
        });
        setBusy(false);
        if (res.success) {
            setRegisterOpen(false);
            setSelectedTypeId(''); setMonths('1'); setPayRef(''); setNotes('');
            await afterChange('Pago registrado.');
        } else {
            setMessage({ type: 'error', text: res.error || 'Error al registrar el pago' });
        }
    };

    const handleExtend = async () => {
        if (extendDays <= 0 || !extendReason.trim()) return;
        setBusy(true);
        const res = await extendCourtesyAction({ businessId: business.id, days: extendDays, reason: extendReason.trim() });
        setBusy(false);
        if (res.success) {
            setExtendOpen(false);
            setExtendReason(''); setExtendDays(7);
            await afterChange(`Se extendieron ${extendDays} días de cortesía.`);
        } else {
            setMessage({ type: 'error', text: res.error || 'Error al extender la vigencia' });
        }
    };

    const isSuspended = business.status === 'cancelled';

    const handleSuspendOrReactivate = async () => {
        setBusy(true);
        const res = isSuspended
            ? await reactivateSubscriptionAction(business.id)
            : await disableSubscriptionAction(business.id);
        setBusy(false);
        if (res.success) {
            await afterChange(isSuspended ? 'Suscripción reactivada.' : 'Suscripción suspendida.');
        } else {
            setMessage({ type: 'error', text: res.error || 'No se pudo completar la acción' });
        }
    };

    const handleRevert = async (subscriptionId?: number) => {
        if (!subscriptionId) return;
        if (!confirm('¿Revertir este pago? La suscripción se recalculará con el pago anterior.')) return;
        setBusy(true);
        const res = await revertPaymentAction(subscriptionId);
        setBusy(false);
        if (res.success) {
            await afterChange('Pago revertido.');
        } else {
            setMessage({ type: 'error', text: res.error || 'No se pudo revertir el pago' });
        }
    };

    const handleGrantOverride = async () => {
        if (!grantModule) return;
        setBusy(true);
        const res = await grantOverrideAction({
            businessId: business.id,
            moduleCode: grantModule,
            notes: grantNotes || undefined,
            expiresAt: grantExpiresAt || undefined,
        });
        setBusy(false);
        if (res.success) {
            setGrantOpen(false);
            setGrantModule(''); setGrantExpiresAt(''); setGrantNotes('');
            await afterChange('Acceso otorgado.');
        } else {
            setMessage({ type: 'error', text: res.error || 'No se pudo otorgar el acceso' });
        }
    };

    const handleRevoke = async (moduleCode: string) => {
        if (!confirm('¿Revocar este acceso adicional?')) return;
        setBusy(true);
        const res = await revokeOverrideAction(business.id, moduleCode);
        setBusy(false);
        if (res.success) {
            await afterChange('Acceso revocado.');
        } else {
            setMessage({ type: 'error', text: res.error || 'No se pudo revocar el acceso' });
        }
    };

    const cycleStart = subscription?.start_date ?? business.cycle_start_date;
    const cycleEnd = subscription?.end_date ?? business.cycle_end_date;

    const handleSaveDates = async (startDate: string, endDate: string, cutoffDay?: number) => {
        setBusy(true);
        const res = await editSubscriptionDatesAction({ businessId: business.id, startDate, endDate, cutoffDay });
        setBusy(false);
        if (res.success) {
            await afterChange('Fechas del ciclo actualizadas.');
        } else {
            setMessage({ type: 'error', text: res.error || 'Error al actualizar las fechas' });
        }
    };

    const planName = subscription?.subscription_type_name || business.plan_name;
    const isExpired = business.status === 'expired';
    const heroTone = isExpired || isSuspended ? 'from-red-700 to-red-900' : 'from-violet-700 to-violet-900';
    const heroTitle = isSuspended
        ? 'Suscripción suspendida'
        : isExpired
            ? 'Suscripción vencida'
            : planName
                ? `${planName} · activa`
                : 'Sin plan asignado';
    const heroSub = cycleEnd
        ? `${isExpired || isSuspended ? 'Venció el' : 'Vence el'} ${formatDate(cycleEnd)}`
        : 'Sin ciclo asignado';
    const progress = computeMembershipProgress(cycleStart, cycleEnd);
    const dueRel = progress
        ? (progress.isExpired ? `vencida ${Math.abs(progress.daysRemaining)} días` : progress.daysRemaining === 0 ? 'vence hoy' : `en ${progress.daysRemaining} días`)
        : '—';
    const paidTotal = history.filter((h) => h.status === 'paid').reduce((sum, h) => sum + h.amount, 0);
    const initials = business.name.split(' ').filter(Boolean).slice(0, 2).map((w) => w[0]).join('').toUpperCase();

    return (
        <Modal isOpen onClose={onClose} showCloseButton={false} size="4xl" noPadding>
            <div className="flex flex-col min-h-full">
                <div className="sticky top-0 z-10 flex items-start gap-3.5 px-6 py-4 border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 flex-shrink-0">
                    <div className="w-10 h-10 rounded-xl bg-violet-50 dark:bg-violet-900/40 text-violet-700 dark:text-violet-300 flex items-center justify-center font-semibold text-sm flex-shrink-0">
                        {initials || '—'}
                    </div>
                    <div className="min-w-0 flex-1 flex flex-col gap-1">
                        <div className="flex items-center gap-2 flex-wrap">
                            <span className="text-[17px] font-semibold text-gray-900 dark:text-white">{business.name}</span>
                            <StatusBadge status={business.status} />
                            {planName && (
                                <span className="inline-flex items-center h-[22px] px-2.5 rounded-full text-xs font-semibold bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 border border-gray-200 dark:border-gray-600">
                                    {planName}
                                </span>
                            )}
                        </div>
                        <span className="text-xs text-gray-500 dark:text-gray-400 font-mono truncate">{business.code}</span>
                    </div>
                    <button onClick={onClose} className="ml-auto w-8 h-8 rounded-lg border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-600 flex-shrink-0">✕</button>
                </div>

                <div className="p-5 grid gap-4 lg:grid-cols-[1fr_290px]">
                    <div className="flex flex-col gap-3.5 min-w-0">
                        {message && <Alert type={message.type} onClose={() => setMessage(null)}>{message.text}</Alert>}

                        <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl overflow-hidden">
                            <div className="px-4 py-3 border-b border-gray-100 dark:border-gray-700 flex items-center justify-between gap-2">
                                <div className="flex items-baseline gap-2">
                                    <span className="text-sm font-semibold text-gray-900 dark:text-white">Historial de pagos</span>
                                    <span className="text-xs text-gray-400">{history.length} movimiento{history.length !== 1 ? 's' : ''} · <span className="font-mono text-gray-600 dark:text-gray-300">{formatCurrency(paidTotal)}</span> acumulado</span>
                                </div>
                                <button onClick={() => setRegisterOpen(true)} className="h-7 px-2.5 rounded-md border border-green-200 bg-green-50 text-green-700 text-xs font-semibold hover:bg-green-100 dark:bg-green-900/30 dark:border-green-800 dark:text-green-400">
                                    + Registrar pago
                                </button>
                            </div>
                            <table className="w-full text-sm">
                                <thead>
                                    <tr className="bg-gray-50 dark:bg-gray-900/40 text-[11px] font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">
                                        <th className="text-left px-4 py-2">Fecha</th>
                                        <th className="text-left px-4 py-2">Concepto</th>
                                        <th className="text-right px-4 py-2">Monto</th>
                                        <th className="text-left px-4 py-2">Método</th>
                                        <th className="text-left px-4 py-2">Estado</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {history.length === 0 && (
                                        <tr><td colSpan={5} className="px-4 py-6 text-center text-gray-400 text-xs">Sin pagos registrados</td></tr>
                                    )}
                                    {history.map((h) => (
                                        <tr key={h.id} className="border-t border-gray-100 dark:border-gray-700 align-top">
                                            <td className="px-4 py-2.5 text-gray-600 dark:text-gray-300">{formatDate(h.created_at)}</td>
                                            <td className="px-4 py-2.5">
                                                <p className="font-medium text-gray-900 dark:text-white">{h.subscription_type_name} · {h.months} {h.months === 1 ? 'mes' : 'meses'}</p>
                                                {h.notes && <p className="text-xs text-gray-400 mt-0.5">{h.notes}</p>}
                                            </td>
                                            <td className={`px-4 py-2.5 text-right font-mono ${h.status === 'reverted' ? 'line-through text-gray-400' : 'text-gray-900 dark:text-white'}`}>{formatCurrency(h.amount)}</td>
                                            <td className="px-4 py-2.5">
                                                <p className="text-gray-600 dark:text-gray-300">{paymentMethodLabel(h.payment_method)}</p>
                                                {h.payment_reference && <p className="text-[10px] font-mono text-gray-400 mt-0.5">{h.payment_reference}</p>}
                                            </td>
                                            <td className="px-4 py-2.5">
                                                <div className="flex flex-col items-start gap-1">
                                                    <PaymentStatusBadge status={h.status} />
                                                    {h.status === 'paid' && (
                                                        <button className="text-[11px] text-red-600 hover:underline" onClick={() => handleRevert(h.id)} disabled={busy}>
                                                            Revertir
                                                        </button>
                                                    )}
                                                </div>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>

                        <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl overflow-hidden">
                            <div className="px-4 py-3 border-b border-gray-100 dark:border-gray-700 flex items-center justify-between gap-2">
                                <div>
                                    <p className="text-sm font-semibold text-gray-900 dark:text-white">Módulos adicionales</p>
                                    <p className="text-xs text-gray-500 dark:text-gray-400">Accesos otorgados por fuera del plan{planName ? ` ${planName}` : ''}.</p>
                                </div>
                                <button onClick={() => setGrantOpen((v) => !v)} className="h-7 px-2.5 rounded-md border border-violet-200 bg-violet-50 text-violet-700 text-xs font-semibold hover:bg-violet-100 dark:bg-violet-900/30 dark:border-violet-800 dark:text-violet-300 flex-shrink-0">
                                    {grantOpen ? 'Cerrar' : '+ Otorgar módulo'}
                                </button>
                            </div>

                            {grantOpen && (
                                <div className="px-4 py-3 bg-gray-50 dark:bg-gray-900/30 border-b border-gray-100 dark:border-gray-700 grid gap-2 sm:grid-cols-[1fr_130px_1.2fr_auto] items-end">
                                    <label className="flex flex-col gap-1 text-xs font-medium text-gray-600 dark:text-gray-300">
                                        Módulo
                                        <select value={grantModule} onChange={(e) => setGrantModule(e.target.value)} className="h-9 px-2 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-800 dark:text-white">
                                            <option value="">Selecciona</option>
                                            {moduleCatalog.map((m) => <option key={m.code} value={m.code}>{m.name}</option>)}
                                        </select>
                                    </label>
                                    <label className="flex flex-col gap-1 text-xs font-medium text-gray-600 dark:text-gray-300">
                                        Vence
                                        <input type="date" value={grantExpiresAt} onChange={(e) => setGrantExpiresAt(e.target.value)} className="h-9 px-2 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-800 dark:text-white" />
                                    </label>
                                    <label className="flex flex-col gap-1 text-xs font-medium text-gray-600 dark:text-gray-300">
                                        Notas internas
                                        <input type="text" value={grantNotes} onChange={(e) => setGrantNotes(e.target.value)} placeholder="Motivo del acceso extra…" className="h-9 px-2 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-800 dark:text-white" />
                                    </label>
                                    <Button size="sm" variant="purple" onClick={handleGrantOverride} loading={busy} disabled={!grantModule}>Otorgar</Button>
                                </div>
                            )}

                            {overrides.length === 0 ? (
                                <div className="py-6 px-5 flex flex-col items-center gap-1.5 text-center">
                                    <div className="w-8 h-8 rounded-lg border border-dashed border-violet-200 dark:border-violet-800 bg-violet-50 dark:bg-violet-900/20 text-violet-600 dark:text-violet-300 flex items-center justify-center text-sm">＋</div>
                                    <p className="text-sm font-semibold text-gray-900 dark:text-white">Sin módulos extra</p>
                                    <p className="text-xs text-gray-500 dark:text-gray-400 max-w-xs">
                                        Este negocio solo tiene lo incluido en su plan. Puedes otorgar acceso puntual a un módulo.
                                    </p>
                                </div>
                            ) : (
                                <table className="w-full text-sm">
                                    <thead>
                                        <tr className="bg-gray-50 dark:bg-gray-900/40 text-[11px] font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">
                                            <th className="text-left px-4 py-2">Módulo</th>
                                            <th className="text-left px-4 py-2">Otorgado</th>
                                            <th className="text-left px-4 py-2">Notas</th>
                                            <th className="text-left px-4 py-2"></th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {overrides.map((o) => (
                                            <tr key={o.id} className="border-t border-gray-100 dark:border-gray-700 align-top">
                                                <td className="px-4 py-2.5">
                                                    <div className="flex items-start gap-2">
                                                        <span className="w-1.5 h-1.5 mt-1.5 rounded-full bg-violet-500 flex-shrink-0" />
                                                        <div>
                                                            <p className="font-semibold text-gray-900 dark:text-white">{moduleCatalog.find((m) => m.code === o.module_code)?.name || o.module_code}</p>
                                                            <p className={`text-xs mt-0.5 ${o.expires_at ? 'text-gray-400' : 'text-green-600'}`}>{o.expires_at ? `Vence ${formatDate(o.expires_at)}` : 'Sin vencimiento'}</p>
                                                        </div>
                                                    </div>
                                                </td>
                                                <td className="px-4 py-2.5 text-gray-600 dark:text-gray-300">{formatDate(o.created_at)}</td>
                                                <td className="px-4 py-2.5 text-gray-500 dark:text-gray-400">{o.notes || '—'}</td>
                                                <td className="px-4 py-2.5 text-right">
                                                    <button className="h-[26px] px-2.5 rounded-md border border-red-200 bg-white dark:bg-gray-800 text-red-700 text-xs font-semibold hover:bg-red-50 dark:hover:bg-red-900/20" onClick={() => handleRevoke(o.module_code)} disabled={busy}>
                                                        Revocar
                                                    </button>
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            )}
                        </div>

                        <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl overflow-hidden">
                            <div className="px-4 py-3 border-b border-gray-100 dark:border-gray-700 flex items-center justify-between">
                                <span className="text-sm font-semibold text-gray-900 dark:text-white">Auditoría</span>
                            </div>
                            <div className="px-4 py-3.5">
                                {auditLogs.length === 0 && <p className="text-xs text-gray-400">Sin actividad registrada</p>}
                                {auditLogs.map((log, i) => (
                                    <div key={log.id} className="grid grid-cols-[14px_1fr] gap-2.5">
                                        <div className="flex flex-col items-center">
                                            <span className="w-2 h-2 rounded-full mt-1.5 flex-shrink-0" style={{ background: auditDotColor(log.action) }} />
                                            {i < auditLogs.length - 1 && <span className="flex-1 w-px bg-gray-100 dark:bg-gray-700" />}
                                        </div>
                                        <div className="pb-3.5">
                                            <p className="text-[12.5px] text-gray-900 dark:text-gray-100 leading-snug"><span className="font-semibold">{log.actor_label}</span> {log.description}</p>
                                            <p className="text-[11.5px] text-gray-400 mt-0.5">{formatDate(log.created_at)}</p>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </div>

                        <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-3 flex items-center gap-2">
                            <Button variant="success" size="sm" className="flex-1 !text-xs" onClick={() => setRegisterOpen(true)}>Registrar pago</Button>
                            <Button variant="quaternary" size="sm" className="flex-1 !text-xs" onClick={() => setExtendOpen(true)}>Días cortesía</Button>
                            <Button variant="danger" size="sm" className="flex-1 !text-xs" onClick={() => setSuspendConfirmOpen(true)} loading={busy}>
                                {isSuspended ? 'Reactivar' : 'Suspender'}
                            </Button>
                        </div>
                    </div>

                    <div className="flex flex-col gap-3 lg:sticky lg:top-0 lg:self-start lg:max-h-[85vh] lg:overflow-y-auto lg:pr-1">
                        <div className={`rounded-xl p-4 text-white bg-gradient-to-br ${heroTone}`}>
                            <p className="text-[11px] font-semibold uppercase tracking-wider opacity-80">Estado de suscripción</p>
                            <p className="text-lg font-semibold mt-1.5">{heroTitle}</p>
                            <p className="text-xs opacity-85 mt-0.5">{heroSub}</p>
                            <div className="h-1.5 rounded-full bg-white/25 mt-3 overflow-hidden">
                                <div className="h-full rounded-full bg-white" style={{ width: `${progress?.percent ?? (isExpired ? 100 : 0)}%` }} />
                            </div>
                            <div className="flex items-center justify-between text-[11.5px] opacity-85 mt-1.5">
                                <span>{progress ? `Día ${progress.daysElapsed} de ${progress.totalDays}` : 'Sin ciclo'}</span>
                                <span>{dueRel}</span>
                            </div>
                        </div>

                        <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-3.5 flex flex-col gap-2 text-[12.5px]">
                            <div className="grid grid-cols-[96px_1fr] gap-2 items-baseline">
                                <span className="text-gray-500 dark:text-gray-400">Plan actual</span>
                                <span className="font-semibold text-right text-gray-900 dark:text-white">{planName || '—'}</span>
                            </div>
                            <div className="grid grid-cols-[96px_1fr] gap-2 items-baseline">
                                <span className="text-gray-500 dark:text-gray-400">Vencimiento</span>
                                <span className={`font-semibold text-right ${isExpired ? 'text-red-600' : 'text-gray-900 dark:text-white'}`}>{formatDate(cycleEnd)}</span>
                            </div>
                            <div className="grid grid-cols-[96px_1fr] gap-2 items-baseline">
                                <span className="text-gray-500 dark:text-gray-400">Último pago</span>
                                <span className="font-semibold text-right font-mono text-gray-900 dark:text-white">{business.last_payment_amount ? formatCurrency(business.last_payment_amount) : '—'}</span>
                            </div>
                            <div className="grid grid-cols-[96px_1fr] gap-2 items-baseline">
                                <span className="text-gray-500 dark:text-gray-400">Overrides</span>
                                <span className={`font-semibold text-right ${overrides.length ? 'text-violet-600' : 'text-gray-400'}`}>{overrides.length} módulo{overrides.length !== 1 ? 's' : ''}</span>
                            </div>
                        </div>

                        <SubscriptionCalendar
                            cycleStartDate={cycleStart}
                            cycleEndDate={cycleEnd}
                            cutoffDay={business.cutoff_day}
                            onSave={handleSaveDates}
                            saving={busy}
                        />
                    </div>
                </div>
            </div>

            <Modal isOpen={registerOpen} onClose={() => setRegisterOpen(false)} title="Registrar pago manual" size="md">
                <div className="space-y-4 p-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Plan</label>
                        <select value={selectedTypeId} onChange={(e) => setSelectedTypeId(e.target.value)} className="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white">
                            <option value="">Selecciona un tipo</option>
                            {subscriptionTypes.length > 0 && (
                                <optgroup label="Catálogo general">
                                    {subscriptionTypes.map((t) => (
                                        <option key={t.id} value={t.id}>{t.name} — {formatCurrency(t.price)}/{t.billing_period === 'monthly' ? 'mes' : 'año'}</option>
                                    ))}
                                </optgroup>
                            )}
                            {customPlans.length > 0 && (
                                <optgroup label="Planes personalizados de este negocio">
                                    {customPlans.map((t) => (
                                        <option key={t.id} value={t.id}>{t.name} — {formatCurrency(t.price)}/{t.billing_period === 'monthly' ? 'mes' : 'año'}</option>
                                    ))}
                                </optgroup>
                            )}
                        </select>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Meses</label>
                        <select value={months} onChange={(e) => setMonths(e.target.value)} className="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white">
                            <option value="1">1 mes</option>
                            <option value="3">3 meses</option>
                            <option value="6">6 meses</option>
                            <option value="12">12 meses</option>
                        </select>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Método</label>
                        <select value={paymentMethod} onChange={(e) => setPaymentMethod(e.target.value)} className="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white">
                            <option value="TRANSFER">Transferencia</option>
                            <option value="CASH">Efectivo</option>
                            <option value="WALLET">Wallet</option>
                        </select>
                    </div>
                    <Input label="Referencia (folio o SPEI)" value={payRef} onChange={(e) => setPayRef(e.target.value)} />
                    <Input label="Notas internas (opcional)" value={notes} onChange={(e) => setNotes(e.target.value)} />
                    <div className="flex justify-end gap-2 pt-2">
                        <Button variant="secondary" onClick={() => setRegisterOpen(false)}>Cancelar</Button>
                        <Button variant="success" onClick={handleRegisterPayment} loading={busy} disabled={!selectedTypeId}>Registrar pago</Button>
                    </div>
                </div>
            </Modal>

            <Modal isOpen={extendOpen} onClose={() => setExtendOpen(false)} title="Extender vigencia (cortesía)" size="sm">
                <div className="space-y-4 p-4">
                    <p className="text-sm text-gray-500 dark:text-gray-400">Sin registrar pago. Se marca como cortesía.</p>
                    <div className="flex gap-2">
                        {[3, 7, 15, 30].map((d) => (
                            <button
                                key={d}
                                onClick={() => setExtendDays(d)}
                                className={`px-3 py-1.5 rounded-lg text-sm border ${extendDays === d ? 'bg-amber-500 text-white border-amber-500' : 'border-gray-200 dark:border-gray-600 text-gray-700 dark:text-gray-300'}`}
                            >
                                +{d}
                            </button>
                        ))}
                    </div>
                    <Input label="Motivo (obligatorio)" value={extendReason} onChange={(e) => setExtendReason(e.target.value)} />
                    <div className="flex justify-end gap-2 pt-2">
                        <Button variant="secondary" onClick={() => setExtendOpen(false)}>Cancelar</Button>
                        <Button variant="quaternary" onClick={handleExtend} loading={busy} disabled={!extendReason.trim()}>Extender</Button>
                    </div>
                </div>
            </Modal>

            <ConfirmModal
                isOpen={suspendConfirmOpen}
                onClose={() => setSuspendConfirmOpen(false)}
                onConfirm={handleSuspendOrReactivate}
                title={isSuspended ? `¿Reactivar la suscripción de ${business.name}?` : `¿Suspender la suscripción de ${business.name}?`}
                message={isSuspended
                    ? 'El negocio recupera acceso inmediato a la plataforma con los datos que tenía.'
                    : 'Los usuarios del negocio perderán acceso al panel de inmediato. La suscripción se puede reactivar después sin perder datos.'}
                confirmText={isSuspended ? 'Sí, reactivar' : 'Sí, suspender'}
                type={isSuspended ? 'info' : 'danger'}
            />
        </Modal>
    );
}

function SubscriptionTypesAdminPanel() {
    const [types, setTypes] = useState<SubscriptionType[]>([]);
    const [moduleCatalog, setModuleCatalog] = useState<ModuleInfo[]>([]);
    const [loading, setLoading] = useState(true);
    const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
    const [editModal, setEditModal] = useState<{ open: boolean; type?: SubscriptionType }>({ open: false });

    const emptyTypeForm = {
        name: '', code: '', description: '', price: '', billing_period: 'monthly', active: true,
        module_codes: [] as string[], max_ecommerce_channels: '0',
        included_shipments: '', shipment_overage_price: '', included_invoices: '', invoice_overage_price: '',
        included_orders: '', order_overage_price: '',
    };
    const [form, setForm] = useState(emptyTypeForm);

    const load = useCallback(async () => {
        setLoading(true);
        const [typesRes, catalogRes] = await Promise.all([listSubscriptionTypesAction(false), getModuleCatalogAction()]);
        if (typesRes.success && typesRes.data) setTypes(typesRes.data);
        if (catalogRes.success && catalogRes.data) setModuleCatalog(catalogRes.data);
        setLoading(false);
    }, []);

    useEffect(() => { load(); }, [load]);

    const openCreate = () => {
        setForm(emptyTypeForm);
        setEditModal({ open: true });
    };

    const openEdit = (t: SubscriptionType) => {
        setForm({
            name: t.name, code: t.code, description: t.description, price: String(t.price), billing_period: t.billing_period, active: t.active,
            module_codes: t.module_codes ?? [], max_ecommerce_channels: String(t.max_ecommerce_channels ?? 0),
            included_shipments: t.included_shipments != null ? String(t.included_shipments) : '',
            shipment_overage_price: t.shipment_overage_price != null ? String(t.shipment_overage_price) : '',
            included_invoices: t.included_invoices != null ? String(t.included_invoices) : '',
            invoice_overage_price: t.invoice_overage_price != null ? String(t.invoice_overage_price) : '',
            included_orders: t.included_orders != null ? String(t.included_orders) : '',
            order_overage_price: t.order_overage_price != null ? String(t.order_overage_price) : '',
        });
        setEditModal({ open: true, type: t });
    };

    const toggleModule = (code: string) => {
        setForm((prev) => ({
            ...prev,
            module_codes: prev.module_codes.includes(code)
                ? prev.module_codes.filter((c) => c !== code)
                : [...prev.module_codes, code],
        }));
    };

    const handleSave = async () => {
        if (!form.name.trim()) {
            setMessage({ type: 'error', text: 'El nombre es obligatorio' });
            return;
        }
        if (!editModal.type && !form.code.trim()) {
            setMessage({ type: 'error', text: 'El código (único) es obligatorio' });
            return;
        }
        if (form.price.trim() === '' || Number.isNaN(Number(form.price)) || Number(form.price) < 0) {
            setMessage({ type: 'error', text: 'El precio debe ser un número mayor o igual a 0' });
            return;
        }

        const includedShipments = form.included_shipments.trim() ? Number(form.included_shipments) : undefined;
        const shipmentOveragePrice = form.shipment_overage_price.trim() ? Number(form.shipment_overage_price) : undefined;
        const includedInvoices = form.included_invoices.trim() ? Number(form.included_invoices) : undefined;
        const invoiceOveragePrice = form.invoice_overage_price.trim() ? Number(form.invoice_overage_price) : undefined;
        const includedOrders = form.included_orders.trim() ? Number(form.included_orders) : undefined;
        const orderOveragePrice = form.order_overage_price.trim() ? Number(form.order_overage_price) : undefined;

        const res = editModal.type
            ? await updateSubscriptionTypeAction(editModal.type.id, {
                name: form.name,
                description: form.description,
                price: Number(form.price),
                billing_period: form.billing_period,
                active: form.active,
                module_codes: form.module_codes,
                max_ecommerce_channels: Number(form.max_ecommerce_channels) || 0,
                included_shipments: includedShipments,
                shipment_overage_price: shipmentOveragePrice,
                included_invoices: includedInvoices,
                invoice_overage_price: invoiceOveragePrice,
                included_orders: includedOrders,
                order_overage_price: orderOveragePrice,
            })
            : await createSubscriptionTypeAction({
                name: form.name,
                code: form.code,
                description: form.description,
                price: Number(form.price),
                billing_period: form.billing_period,
                module_codes: form.module_codes,
                max_ecommerce_channels: Number(form.max_ecommerce_channels) || 0,
                included_shipments: includedShipments,
                shipment_overage_price: shipmentOveragePrice,
                included_invoices: includedInvoices,
                invoice_overage_price: invoiceOveragePrice,
                included_orders: includedOrders,
                order_overage_price: orderOveragePrice,
            });

        if (res.success) {
            setMessage({ type: 'success', text: 'Tipo de suscripción guardado' });
            setEditModal({ open: false });
            load();
        } else {
            setMessage({ type: 'error', text: res.error || 'Error al guardar' });
        }
    };

    const handleDelete = async (t: SubscriptionType) => {
        if (!confirm(`¿Eliminar el tipo de suscripción "${t.name}"?`)) return;
        const res = await deleteSubscriptionTypeAction(t.id);
        if (res.success) {
            setMessage({ type: 'success', text: 'Tipo de suscripción eliminado' });
            load();
        } else {
            setMessage({ type: 'error', text: res.error || 'Error al eliminar' });
        }
    };

    const moduleName = (code: string) => moduleCatalog.find((m) => m.code === code)?.name ?? code;

    if (loading) return <div className="flex justify-center py-12"><Spinner /></div>;

    return (
        <div className="space-y-4">
            {message && <Alert type={message.type} onClose={() => setMessage(null)}>{message.text}</Alert>}

            <div className="flex justify-end">
                <Button variant="primary" onClick={openCreate}>+ Nuevo tipo de suscripción</Button>
            </div>

            <div className="grid gap-5 lg:grid-cols-3">
                {types.map((t) => (
                    <div
                        key={t.id}
                        className={`relative flex flex-col rounded-2xl border bg-white dark:bg-gray-800 overflow-hidden transition-opacity ${
                            t.active ? 'border-violet-200 dark:border-violet-800/60 shadow-sm shadow-violet-500/5' : 'border-gray-200 dark:border-gray-700 opacity-60'
                        }`}
                    >
                        <div className={`flex items-center justify-center gap-1.5 text-[11px] font-bold uppercase tracking-wide py-1.5 ${
                            t.active
                                ? 'bg-gradient-to-r from-violet-700 to-purple-500 text-white'
                                : 'bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400'
                        }`}>
                            {t.active ? 'Activo' : 'Inactivo'}
                        </div>

                        <div className="pt-5 px-6 pb-5">
                            <div className="flex items-center justify-between mb-3">
                                <h4 className="text-lg font-bold text-gray-900 dark:text-white">{t.name}</h4>
                                <span className="text-[11px] font-mono text-gray-400">{t.code}</span>
                            </div>

                            <div className="flex items-end gap-1 mb-3">
                                <span className="text-3xl font-extrabold tracking-tight text-gray-900 dark:text-white">{formatCurrency(t.price)}</span>
                                <span className="text-sm font-semibold text-gray-400 pb-1">/{t.billing_period === 'monthly' ? 'mes' : 'año'}</span>
                            </div>

                            {t.description && (
                                <p className="text-sm text-gray-500 dark:text-gray-400 leading-relaxed mb-4 min-h-[2.6rem]">{t.description}</p>
                            )}

                            <div className="flex items-center gap-3 rounded-xl p-3 bg-violet-50 dark:bg-violet-900/20">
                                <div className="w-[34px] h-[34px] rounded-lg flex items-center justify-center flex-shrink-0 bg-violet-100 dark:bg-violet-800/40 text-violet-600 dark:text-violet-300">
                                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                                        <path d="M3 3H21V8H3V3Z" stroke="currentColor" strokeWidth="1.8" strokeLinejoin="round" />
                                        <path d="M4 8V20C4 20.5523 4.44772 21 5 21H19C19.5523 21 20 20.5523 20 20V8" stroke="currentColor" strokeWidth="1.8" strokeLinejoin="round" />
                                        <path d="M9 12H15" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" />
                                    </svg>
                                </div>
                                <div>
                                    <div className="text-sm font-bold text-gray-900 dark:text-white leading-tight">
                                        {t.max_ecommerce_channels > 0 ? `Hasta ${t.max_ecommerce_channels}` : 'Ilimitados'}
                                    </div>
                                    <div className="text-xs text-gray-400 leading-tight">canales de ecommerce conectados</div>
                                </div>
                            </div>

                            {(t.included_shipments != null || t.included_invoices != null || t.included_orders != null) && (
                                <div className="mt-3 space-y-1.5 text-xs text-gray-500 dark:text-gray-400">
                                    {t.included_shipments != null && (
                                        <p>
                                            Hasta <span className="font-semibold text-gray-700 dark:text-gray-200">{t.included_shipments}</span> envíos/mes
                                            {t.shipment_overage_price != null && <> · adicional: <span className="font-semibold text-gray-700 dark:text-gray-200">{formatCurrency(t.shipment_overage_price)}</span>/guía</>}
                                        </p>
                                    )}
                                    {t.included_invoices != null && (
                                        <p>
                                            Hasta <span className="font-semibold text-gray-700 dark:text-gray-200">{t.included_invoices}</span> facturas/mes
                                            {t.invoice_overage_price != null && <> · adicional: <span className="font-semibold text-gray-700 dark:text-gray-200">{formatCurrency(t.invoice_overage_price)}</span>/factura</>}
                                        </p>
                                    )}
                                    {t.included_orders != null && (
                                        <p>
                                            Hasta <span className="font-semibold text-gray-700 dark:text-gray-200">{t.included_orders}</span> ordenes/mes
                                            {t.order_overage_price != null && <> · adicional: <span className="font-semibold text-gray-700 dark:text-gray-200">{formatCurrency(t.order_overage_price)}</span>/orden</>}
                                        </p>
                                    )}
                                </div>
                            )}
                        </div>

                        <div className="border-t border-violet-100 dark:border-violet-900/40 mx-6" />

                        <div className="px-6 pt-4 pb-2 flex-1">
                            <div className="text-[11px] font-bold uppercase tracking-wide text-gray-400 mb-3">Módulos incluidos</div>
                            <div className="flex flex-col gap-2.5">
                                {(t.module_codes ?? []).length === 0 && (
                                    <span className="text-xs text-gray-400 italic">Sin módulos asignados</span>
                                )}
                                {(t.module_codes ?? []).map((m) => (
                                    <div key={m} className="flex items-center gap-2.5">
                                        <span className="w-[18px] h-[18px] rounded-full flex items-center justify-center flex-shrink-0 bg-violet-100 dark:bg-violet-800/40 text-violet-600 dark:text-violet-300">
                                            <svg width="11" height="11" viewBox="0 0 24 24" fill="none">
                                                <path d="M4 12.5L9.5 18L20 6.5" stroke="currentColor" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round" />
                                            </svg>
                                        </span>
                                        <span className="text-[13px] font-medium text-gray-700 dark:text-gray-200">{moduleName(m)}</span>
                                    </div>
                                ))}
                            </div>
                        </div>

                        <div className="px-6 pt-5 pb-6 flex gap-2">
                            <Button size="sm" variant="outline-purple" onClick={() => openEdit(t)} className="flex-1">Editar</Button>
                            <Button size="sm" variant="danger" onClick={() => handleDelete(t)} className="flex-1">Eliminar</Button>
                        </div>
                    </div>
                ))}
            </div>

            <Modal isOpen={editModal.open} onClose={() => setEditModal({ open: false })} size="4xl">
                <div className="flex flex-col">
                    <div className="flex items-start justify-between gap-4 pb-4 mb-4 border-b border-gray-100 dark:border-gray-700">
                        <div>
                            <div className="flex items-center gap-2.5">
                                <h3 className="text-lg font-bold text-gray-900 dark:text-white">
                                    {editModal.type ? `Editar ${editModal.type.name}` : 'Nuevo tipo de suscripción'}
                                </h3>
                                {editModal.type && (
                                    <span className="font-mono text-[11.5px] font-semibold text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2.5 py-0.5 rounded-full">
                                        {editModal.type.code}
                                    </span>
                                )}
                            </div>
                            <p className="text-[13px] text-gray-500 dark:text-gray-400 mt-1">
                                {editModal.type ? 'Los cambios se aplican a nuevos clientes y renovaciones.' : 'Define los datos y el alcance del nuevo plan.'}
                            </p>
                        </div>
                        <button
                            onClick={() => setEditModal({ open: false })}
                            className="w-8 h-8 rounded-lg flex items-center justify-center text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors flex-shrink-0"
                        >
                            <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><path d="M5 5L19 19M19 5L5 19" stroke="currentColor" strokeWidth="2" strokeLinecap="round" /></svg>
                        </button>
                    </div>

                    {message && <div className="mb-4"><Alert type={message.type} onClose={() => setMessage(null)}>{message.text}</Alert></div>}

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-x-8 gap-y-4">
                        <div className="space-y-4">
                            <div className="flex items-baseline justify-between">
                                <span className="text-xs font-bold uppercase tracking-wide text-violet-600 dark:text-violet-400">Datos básicos</span>
                                <span className="text-xs text-gray-400">Cómo se cobra</span>
                            </div>

                            <Input label="Nombre *" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
                            {!editModal.type && (
                                <Input label="Código (único) *" value={form.code} onChange={(e) => setForm({ ...form, code: e.target.value })} placeholder="ej: pro-mensual" className="font-mono" />
                            )}
                            <Input label="Descripción" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} />

                            <div>
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">Precio *</label>
                                <div className="relative">
                                    <span className="absolute left-3 top-1/2 -translate-y-1/2 text-sm font-semibold text-gray-400">$</span>
                                    <input
                                        type="number"
                                        value={form.price}
                                        onChange={(e) => setForm({ ...form, price: e.target.value })}
                                        className="w-full pl-7 pr-3 py-2.5 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-2 focus:ring-violet-500/30 focus:border-violet-500"
                                    />
                                </div>
                                <div className="text-xs font-semibold text-violet-600 dark:text-violet-400 mt-1.5">
                                    ${(Number(form.price) || 0).toLocaleString('es-CO')} COP / {form.billing_period === 'monthly' ? 'mes' : 'año'}
                                </div>
                            </div>

                            <div>
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">Periodo de facturación</label>
                                <select value={form.billing_period} onChange={(e) => setForm({ ...form, billing_period: e.target.value })} className="w-full px-3 py-2.5 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-2 focus:ring-violet-500/30 focus:border-violet-500">
                                    <option value="monthly">Mensual</option>
                                    <option value="annual">Anual</option>
                                </select>
                            </div>

                            {editModal.type && (
                                <div
                                    onClick={() => setForm({ ...form, active: !form.active })}
                                    className="flex items-start gap-2.5 bg-violet-50 dark:bg-violet-900/20 rounded-lg px-3.5 py-3 cursor-pointer"
                                >
                                    <span className={`w-5 h-5 rounded-md flex items-center justify-center flex-shrink-0 mt-0.5 border-[1.5px] ${form.active ? 'bg-violet-600 border-violet-600 text-white' : 'border-gray-300 dark:border-gray-600 text-transparent'}`}>
                                        <svg width="12" height="12" viewBox="0 0 24 24" fill="none"><path d="M4 12.5L9.5 18L20 6.5" stroke="currentColor" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round" /></svg>
                                    </span>
                                    <div>
                                        <div className="text-[13.5px] font-semibold text-gray-900 dark:text-white">Plan activo</div>
                                        <div className="text-xs text-gray-500 dark:text-gray-400 mt-0.5 leading-relaxed">Los planes inactivos no pueden ser contratados por nuevos clientes.</div>
                                    </div>
                                </div>
                            )}
                        </div>

                        <div className="space-y-4 sm:border-l sm:border-gray-100 dark:sm:border-gray-700 sm:pl-8">
                            <div className="flex items-baseline justify-between">
                                <span className="text-xs font-bold uppercase tracking-wide text-violet-600 dark:text-violet-400">Alcance del plan</span>
                                <span className="text-xs text-gray-400">Qué incluye</span>
                            </div>

                            <div>
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">
                                    Límite de canales e-commerce <span className="font-normal text-gray-400">— 0 sin límite</span>
                                </label>
                                <input
                                    type="number"
                                    value={form.max_ecommerce_channels}
                                    onChange={(e) => setForm({ ...form, max_ecommerce_channels: e.target.value })}
                                    className="w-32 px-3 py-2.5 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-2 focus:ring-violet-500/30 focus:border-violet-500"
                                />
                                <p className="text-xs text-gray-400 mt-1.5 leading-relaxed">Al superar el límite, el cliente no podrá conectar canales adicionales hasta cambiar de plan.</p>
                            </div>

                            <div className="grid grid-cols-2 gap-3">
                                <div>
                                    <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">Envíos incluidos <span className="font-normal text-gray-400">/ mes</span></label>
                                    <input
                                        type="number"
                                        value={form.included_shipments}
                                        onChange={(e) => setForm({ ...form, included_shipments: e.target.value })}
                                        placeholder="ej: 100"
                                        className="w-full px-3 py-2.5 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-2 focus:ring-violet-500/30 focus:border-violet-500"
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">Costo envío extra</label>
                                    <div className="relative">
                                        <span className="absolute left-3 top-1/2 -translate-y-1/2 text-sm font-semibold text-gray-400">$</span>
                                        <input
                                            type="number"
                                            value={form.shipment_overage_price}
                                            onChange={(e) => setForm({ ...form, shipment_overage_price: e.target.value })}
                                            placeholder="ej: 600"
                                            className="w-full pl-7 pr-3 py-2.5 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-2 focus:ring-violet-500/30 focus:border-violet-500"
                                        />
                                    </div>
                                </div>
                            </div>

                            <div className="grid grid-cols-2 gap-3">
                                <div>
                                    <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">Facturas incluidas <span className="font-normal text-gray-400">/ mes</span></label>
                                    <input
                                        type="number"
                                        value={form.included_invoices}
                                        onChange={(e) => setForm({ ...form, included_invoices: e.target.value })}
                                        placeholder="ej: 6000"
                                        className="w-full px-3 py-2.5 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-2 focus:ring-violet-500/30 focus:border-violet-500"
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">Costo factura extra</label>
                                    <div className="relative">
                                        <span className="absolute left-3 top-1/2 -translate-y-1/2 text-sm font-semibold text-gray-400">$</span>
                                        <input
                                            type="number"
                                            value={form.invoice_overage_price}
                                            onChange={(e) => setForm({ ...form, invoice_overage_price: e.target.value })}
                                            placeholder="ej: 550"
                                            className="w-full pl-7 pr-3 py-2.5 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-2 focus:ring-violet-500/30 focus:border-violet-500"
                                        />
                                    </div>
                                </div>
                            </div>

                            <div className="grid grid-cols-2 gap-3">
                                <div>
                                    <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">{'\u00d3rdenes'} incluidas <span className="font-normal text-gray-400">/ mes</span></label>
                                    <input
                                        type="number"
                                        value={form.included_orders}
                                        onChange={(e) => setForm({ ...form, included_orders: e.target.value })}
                                        placeholder="ej: 6000"
                                        className="w-full px-3 py-2.5 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-2 focus:ring-violet-500/30 focus:border-violet-500"
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">Costo orden extra</label>
                                    <div className="relative">
                                        <span className="absolute left-3 top-1/2 -translate-y-1/2 text-sm font-semibold text-gray-400">$</span>
                                        <input
                                            type="number"
                                            value={form.order_overage_price}
                                            onChange={(e) => setForm({ ...form, order_overage_price: e.target.value })}
                                            placeholder="ej: 550"
                                            className="w-full pl-7 pr-3 py-2.5 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-2 focus:ring-violet-500/30 focus:border-violet-500"
                                        />
                                    </div>
                                </div>
                            </div>

                            <div>
                                <div className="flex items-center justify-between mb-2.5">
                                    <label className="text-sm font-semibold text-gray-700 dark:text-gray-300">Módulos incluidos</label>
                                    <span className="text-xs font-semibold text-violet-600 dark:text-violet-400">{form.module_codes.length} de {moduleCatalog.length}</span>
                                </div>
                                <div className="grid grid-cols-2 lg:grid-cols-3 gap-2">
                                    {moduleCatalog.map(({ code, name }) => {
                                        const checked = form.module_codes.includes(code);
                                        return (
                                            <div
                                                key={code}
                                                onClick={() => toggleModule(code)}
                                                className={`flex items-center gap-2 px-2.5 py-2 rounded-lg border-[1.5px] cursor-pointer transition-colors ${
                                                    checked
                                                        ? 'bg-violet-50 dark:bg-violet-900/20 border-violet-500'
                                                        : 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-600'
                                                }`}
                                            >
                                                <span className={`w-4 h-4 rounded flex items-center justify-center flex-shrink-0 border-[1.5px] ${checked ? 'bg-violet-600 border-violet-600 text-white' : 'border-gray-300 dark:border-gray-600 text-transparent'}`}>
                                                    <svg width="10" height="10" viewBox="0 0 24 24" fill="none"><path d="M4 12.5L9.5 18L20 6.5" stroke="currentColor" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round" /></svg>
                                                </span>
                                                <span className="text-[13px] font-medium text-gray-700 dark:text-gray-200 truncate">{name}</span>
                                            </div>
                                        );
                                    })}
                                </div>
                            </div>
                        </div>
                    </div>

                    <div className="flex justify-end gap-2.5 pt-4 mt-4 border-t border-gray-100 dark:border-gray-700">
                        <Button variant="secondary" onClick={() => setEditModal({ open: false })}>Cancelar</Button>
                        <Button variant="purple" onClick={handleSave}>Guardar</Button>
                    </div>
                </div>
            </Modal>
        </div>
    );
}

function CustomPlansAdminPanel({ businesses }: { businesses: Array<{ id: number; name: string }> }) {
    const [plans, setPlans] = useState<SubscriptionType[]>([]);
    const [moduleCatalog, setModuleCatalog] = useState<ModuleInfo[]>([]);
    const [loading, setLoading] = useState(true);
    const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
    const [editModal, setEditModal] = useState<{ open: boolean; plan?: SubscriptionType }>({ open: false });

    const emptyForm = {
        name: '', code: '', description: '', price: '', billing_period: 'monthly', active: true,
        module_codes: [] as string[], max_ecommerce_channels: '0', business_id: '', months: '1', notes: '',
        included_shipments: '', shipment_overage_price: '', included_invoices: '', invoice_overage_price: '',
        included_orders: '', order_overage_price: '',
    };
    const [form, setForm] = useState(emptyForm);

    const businessName = (id?: number) => businesses.find((b) => b.id === id)?.name ?? `Negocio ${id}`;

    const load = useCallback(async () => {
        setLoading(true);
        const [plansRes, catalogRes] = await Promise.all([listCustomPlansAction(), getModuleCatalogAction()]);
        if (plansRes.success && plansRes.data) setPlans(plansRes.data);
        if (catalogRes.success && catalogRes.data) setModuleCatalog(catalogRes.data);
        setLoading(false);
    }, []);

    useEffect(() => { load(); }, [load]);

    const openCreate = () => {
        setForm(emptyForm);
        setEditModal({ open: true });
    };

    const openEdit = (p: SubscriptionType) => {
        setForm({
            name: p.name, code: p.code, description: p.description, price: String(p.price), billing_period: p.billing_period, active: p.active,
            module_codes: p.module_codes ?? [], max_ecommerce_channels: String(p.max_ecommerce_channels ?? 0), business_id: String(p.business_id ?? ''), months: '1', notes: '',
            included_shipments: p.included_shipments != null ? String(p.included_shipments) : '',
            shipment_overage_price: p.shipment_overage_price != null ? String(p.shipment_overage_price) : '',
            included_invoices: p.included_invoices != null ? String(p.included_invoices) : '',
            invoice_overage_price: p.invoice_overage_price != null ? String(p.invoice_overage_price) : '',
            included_orders: p.included_orders != null ? String(p.included_orders) : '',
            order_overage_price: p.order_overage_price != null ? String(p.order_overage_price) : '',
        });
        setEditModal({ open: true, plan: p });
    };

    const toggleModule = (code: string) => {
        setForm((prev) => ({
            ...prev,
            module_codes: prev.module_codes.includes(code)
                ? prev.module_codes.filter((c) => c !== code)
                : [...prev.module_codes, code],
        }));
    };

    const handleSave = async () => {
        if (!form.name.trim()) {
            setMessage({ type: 'error', text: 'El nombre es obligatorio' });
            return;
        }
        if (form.price.trim() === '' || Number.isNaN(Number(form.price)) || Number(form.price) < 0) {
            setMessage({ type: 'error', text: 'El precio debe ser un número mayor o igual a 0' });
            return;
        }
        if (!editModal.plan) {
            if (!form.business_id) {
                setMessage({ type: 'error', text: 'Selecciona el negocio al que se atará este plan' });
                return;
            }
            if (!form.code.trim()) {
                setMessage({ type: 'error', text: 'El código es obligatorio y debe ser único (no puede repetir uno ya usado por otro plan)' });
                return;
            }
        }

        const includedShipments = form.included_shipments.trim() ? Number(form.included_shipments) : undefined;
        const shipmentOveragePrice = form.shipment_overage_price.trim() ? Number(form.shipment_overage_price) : undefined;
        const includedInvoices = form.included_invoices.trim() ? Number(form.included_invoices) : undefined;
        const invoiceOveragePrice = form.invoice_overage_price.trim() ? Number(form.invoice_overage_price) : undefined;
        const includedOrders = form.included_orders.trim() ? Number(form.included_orders) : undefined;
        const orderOveragePrice = form.order_overage_price.trim() ? Number(form.order_overage_price) : undefined;

        const res = editModal.plan
            ? await updateCustomPlanAction(editModal.plan.id, {
                name: form.name,
                description: form.description,
                price: Number(form.price),
                billing_period: form.billing_period,
                active: form.active,
                module_codes: form.module_codes,
                max_ecommerce_channels: Number(form.max_ecommerce_channels) || 0,
                included_shipments: includedShipments,
                shipment_overage_price: shipmentOveragePrice,
                included_invoices: includedInvoices,
                invoice_overage_price: invoiceOveragePrice,
                included_orders: includedOrders,
                order_overage_price: orderOveragePrice,
            })
            : await createCustomPlanAction({
                name: form.name,
                code: form.code,
                description: form.description,
                price: Number(form.price),
                billing_period: form.billing_period,
                module_codes: form.module_codes,
                max_ecommerce_channels: Number(form.max_ecommerce_channels) || 0,
                business_id: Number(form.business_id),
                months: Number(form.months) || 1,
                notes: form.notes || undefined,
                included_shipments: includedShipments,
                shipment_overage_price: shipmentOveragePrice,
                included_invoices: includedInvoices,
                invoice_overage_price: invoiceOveragePrice,
                included_orders: includedOrders,
                order_overage_price: orderOveragePrice,
            });

        if (res.success) {
            setMessage({ type: 'success', text: editModal.plan ? 'Plan personalizado actualizado' : 'Plan personalizado creado y asociado al negocio' });
            setEditModal({ open: false });
            load();
        } else {
            const isDuplicateCode = res.error?.includes('duplicate key') && res.error?.includes('idx_subscription_types_code');
            setMessage({
                type: 'error',
                text: isDuplicateCode
                    ? 'Ese código ya está en uso por otro plan (el código debe ser único en todo el sistema, no solo para este negocio). Prueba con otro.'
                    : res.error || 'Error al guardar',
            });
        }
    };

    const handleDelete = async (p: SubscriptionType) => {
        if (!confirm(`¿Eliminar el plan personalizado "${p.name}" de ${businessName(p.business_id)}?`)) return;
        const res = await deleteCustomPlanAction(p.id);
        if (res.success) {
            setMessage({ type: 'success', text: 'Plan personalizado eliminado' });
            load();
        } else {
            setMessage({ type: 'error', text: res.error || 'Error al eliminar' });
        }
    };

    const moduleName = (code: string) => moduleCatalog.find((m) => m.code === code)?.name ?? code;

    if (loading) return <div className="flex justify-center py-12"><Spinner /></div>;

    return (
        <div className="space-y-4">
            <p className="text-sm text-gray-500 dark:text-gray-400">
                Planes a la medida, atados a un negocio específico. Al crearlos quedan activados de inmediato para ese negocio.
            </p>

            {message && <Alert type={message.type} onClose={() => setMessage(null)}>{message.text}</Alert>}

            <div className="flex justify-end">
                <Button variant="primary" onClick={openCreate}>+ Nuevo plan personalizado</Button>
            </div>

            {plans.length === 0 && (
                <div className="text-center py-10 text-sm text-gray-400">Aún no hay planes personalizados creados.</div>
            )}

            <div className="grid gap-5 lg:grid-cols-3">
                {plans.map((p) => (
                    <div
                        key={p.id}
                        className={`relative flex flex-col rounded-2xl border bg-white dark:bg-gray-800 overflow-hidden transition-opacity ${
                            p.active ? 'border-violet-200 dark:border-violet-800/60 shadow-sm shadow-violet-500/5' : 'border-gray-200 dark:border-gray-700 opacity-60'
                        }`}
                    >
                        <div className="flex items-center justify-center gap-1.5 text-[11px] font-bold uppercase tracking-wide py-1.5 bg-gradient-to-r from-indigo-700 to-blue-500 text-white">
                            {businessName(p.business_id)}
                        </div>

                        <div className="pt-5 px-6 pb-5">
                            <div className="flex items-center justify-between mb-3">
                                <h4 className="text-lg font-bold text-gray-900 dark:text-white">{p.name}</h4>
                                <span className="text-[11px] font-mono text-gray-400">{p.code}</span>
                            </div>

                            <div className="flex items-end gap-1 mb-3">
                                <span className="text-3xl font-extrabold tracking-tight text-gray-900 dark:text-white">{formatCurrency(p.price)}</span>
                                <span className="text-sm font-semibold text-gray-400 pb-1">/{p.billing_period === 'monthly' ? 'mes' : 'año'}</span>
                            </div>

                            {p.description && (
                                <p className="text-sm text-gray-500 dark:text-gray-400 leading-relaxed mb-2">{p.description}</p>
                            )}

                            <div className="flex items-center gap-3 rounded-xl p-3 bg-indigo-50 dark:bg-indigo-900/20">
                                <div className="w-[34px] h-[34px] rounded-lg flex items-center justify-center flex-shrink-0 bg-indigo-100 dark:bg-indigo-800/40 text-indigo-600 dark:text-indigo-300">
                                    <ChannelsIcon />
                                </div>
                                <div>
                                    <div className="text-sm font-bold text-gray-900 dark:text-white leading-tight">
                                        {p.max_ecommerce_channels > 0 ? `Hasta ${p.max_ecommerce_channels}` : 'Ilimitados'}
                                    </div>
                                    <div className="text-xs text-gray-400 leading-tight">canales de ecommerce conectados</div>
                                </div>
                            </div>

                            {(p.included_shipments != null || p.included_invoices != null || p.included_orders != null) && (
                                <div className="mt-3 space-y-1.5 text-xs text-gray-500 dark:text-gray-400">
                                    {p.included_shipments != null && (
                                        <p>
                                            Hasta <span className="font-semibold text-gray-700 dark:text-gray-200">{p.included_shipments}</span> envíos/mes
                                            {p.shipment_overage_price != null && <> · adicional: <span className="font-semibold text-gray-700 dark:text-gray-200">{formatCurrency(p.shipment_overage_price)}</span>/guía</>}
                                        </p>
                                    )}
                                    {p.included_invoices != null && (
                                        <p>
                                            Hasta <span className="font-semibold text-gray-700 dark:text-gray-200">{p.included_invoices}</span> facturas/mes
                                            {p.invoice_overage_price != null && <> · adicional: <span className="font-semibold text-gray-700 dark:text-gray-200">{formatCurrency(p.invoice_overage_price)}</span>/factura</>}
                                        </p>
                                    )}
                                    {p.included_orders != null && (
                                        <p>
                                            Hasta <span className="font-semibold text-gray-700 dark:text-gray-200">{p.included_orders}</span> ordenes/mes
                                            {p.order_overage_price != null && <> · adicional: <span className="font-semibold text-gray-700 dark:text-gray-200">{formatCurrency(p.order_overage_price)}</span>/orden</>}
                                        </p>
                                    )}
                                </div>
                            )}
                        </div>

                        <div className="border-t border-indigo-100 dark:border-indigo-900/40 mx-6" />

                        <div className="px-6 pt-4 pb-2 flex-1">
                            <div className="text-[11px] font-bold uppercase tracking-wide text-gray-400 mb-3">Módulos incluidos</div>
                            <div className="flex flex-col gap-2.5">
                                {(p.module_codes ?? []).length === 0 && (
                                    <span className="text-xs text-gray-400 italic">Sin módulos asignados</span>
                                )}
                                {(p.module_codes ?? []).map((m) => (
                                    <div key={m} className="flex items-center gap-2.5">
                                        <span className="w-[18px] h-[18px] rounded-full flex items-center justify-center flex-shrink-0 bg-indigo-100 dark:bg-indigo-800/40 text-indigo-600 dark:text-indigo-300">
                                            <CheckIcon />
                                        </span>
                                        <span className="text-[13px] font-medium text-gray-700 dark:text-gray-200">{moduleName(m)}</span>
                                    </div>
                                ))}
                            </div>
                        </div>

                        <div className="px-6 pt-5 pb-6 flex gap-2">
                            <Button size="sm" variant="outline-purple" onClick={() => openEdit(p)} className="flex-1">Editar</Button>
                            <Button size="sm" variant="danger" onClick={() => handleDelete(p)} className="flex-1">Eliminar</Button>
                        </div>
                    </div>
                ))}
            </div>

            <Modal isOpen={editModal.open} onClose={() => setEditModal({ open: false })} size="4xl">
                <div className="flex flex-col">
                    <div className="flex items-start justify-between gap-4 pb-4 mb-4 border-b border-gray-100 dark:border-gray-700">
                        <div>
                            <div className="flex items-center gap-2.5">
                                <h3 className="text-lg font-bold text-gray-900 dark:text-white">
                                    {editModal.plan ? `Editar ${editModal.plan.name}` : 'Nuevo plan personalizado'}
                                </h3>
                                {editModal.plan && (
                                    <span className="font-mono text-[11.5px] font-semibold text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2.5 py-0.5 rounded-full">
                                        {editModal.plan.code}
                                    </span>
                                )}
                            </div>
                            <p className="text-[13px] text-gray-500 dark:text-gray-400 mt-1">
                                {editModal.plan ? 'El negocio asociado no se puede cambiar una vez creado el plan.' : 'Elige el negocio, define permisos y el valor a pagar.'}
                            </p>
                        </div>
                        <button
                            onClick={() => setEditModal({ open: false })}
                            className="w-8 h-8 rounded-lg flex items-center justify-center text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors flex-shrink-0"
                        >
                            <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><path d="M5 5L19 19M19 5L5 19" stroke="currentColor" strokeWidth="2" strokeLinecap="round" /></svg>
                        </button>
                    </div>

                    {message && <div className="mb-4"><Alert type={message.type} onClose={() => setMessage(null)}>{message.text}</Alert></div>}

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-x-8 gap-y-4">
                        <div className="space-y-4">
                            <div className="flex items-baseline justify-between">
                                <span className="text-xs font-bold uppercase tracking-wide text-indigo-600 dark:text-indigo-400">Datos básicos</span>
                                <span className="text-xs text-gray-400">Cómo se cobra</span>
                            </div>

                            <div>
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">Negocio *</label>
                                <select
                                    value={form.business_id}
                                    onChange={(e) => setForm({ ...form, business_id: e.target.value })}
                                    disabled={!!editModal.plan}
                                    className="w-full px-3 py-2.5 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white disabled:opacity-60 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-indigo-500/30 focus:border-indigo-500"
                                >
                                    <option value="">Selecciona un negocio</option>
                                    {businesses.map((b) => (
                                        <option key={b.id} value={b.id}>{b.name}</option>
                                    ))}
                                </select>
                            </div>

                            <Input label="Nombre *" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
                            {!editModal.plan && (
                                <Input label="Código (único) *" value={form.code} onChange={(e) => setForm({ ...form, code: e.target.value })} placeholder="ej: vip-negocio-x" className="font-mono" />
                            )}
                            <Input label="Descripción" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} />

                            <div>
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">Precio *</label>
                                <div className="relative">
                                    <span className="absolute left-3 top-1/2 -translate-y-1/2 text-sm font-semibold text-gray-400">$</span>
                                    <input
                                        type="number"
                                        value={form.price}
                                        onChange={(e) => setForm({ ...form, price: e.target.value })}
                                        className="w-full pl-7 pr-3 py-2.5 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-2 focus:ring-indigo-500/30 focus:border-indigo-500"
                                    />
                                </div>
                                <div className="text-xs font-semibold text-indigo-600 dark:text-indigo-400 mt-1.5">
                                    ${(Number(form.price) || 0).toLocaleString('es-CO')} COP / {form.billing_period === 'monthly' ? 'mes' : 'año'}
                                </div>
                            </div>

                            <div>
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">Periodo de facturación</label>
                                <select value={form.billing_period} onChange={(e) => setForm({ ...form, billing_period: e.target.value })} className="w-full px-3 py-2.5 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-2 focus:ring-indigo-500/30 focus:border-indigo-500">
                                    <option value="monthly">Mensual</option>
                                    <option value="annual">Anual</option>
                                </select>
                            </div>

                            {!editModal.plan && (
                                <div>
                                    <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">Meses a activar de una vez</label>
                                    <select value={form.months} onChange={(e) => setForm({ ...form, months: e.target.value })} className="w-full px-3 py-2.5 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-2 focus:ring-indigo-500/30 focus:border-indigo-500">
                                        <option value="1">1 mes</option>
                                        <option value="3">3 meses</option>
                                        <option value="6">6 meses</option>
                                        <option value="12">12 meses (anual)</option>
                                    </select>
                                    <p className="text-xs text-gray-400 mt-1.5 leading-relaxed">Al guardar, el plan queda activo de inmediato para el negocio seleccionado por esta cantidad de meses.</p>
                                </div>
                            )}

                            {editModal.plan && (
                                <div
                                    onClick={() => setForm({ ...form, active: !form.active })}
                                    className="flex items-start gap-2.5 bg-indigo-50 dark:bg-indigo-900/20 rounded-lg px-3.5 py-3 cursor-pointer"
                                >
                                    <span className={`w-5 h-5 rounded-md flex items-center justify-center flex-shrink-0 mt-0.5 border-[1.5px] ${form.active ? 'bg-indigo-600 border-indigo-600 text-white' : 'border-gray-300 dark:border-gray-600 text-transparent'}`}>
                                        <CheckIcon />
                                    </span>
                                    <div>
                                        <div className="text-[13.5px] font-semibold text-gray-900 dark:text-white">Plan activo</div>
                                        <div className="text-xs text-gray-500 dark:text-gray-400 mt-0.5 leading-relaxed">Si lo desactivas, no afecta la suscripción ya asociada, solo impide reutilizarlo.</div>
                                    </div>
                                </div>
                            )}
                        </div>

                        <div className="space-y-4 sm:border-l sm:border-gray-100 dark:sm:border-gray-700 sm:pl-8">
                            <div className="flex items-baseline justify-between">
                                <span className="text-xs font-bold uppercase tracking-wide text-indigo-600 dark:text-indigo-400">Alcance del plan</span>
                                <span className="text-xs text-gray-400">Qué incluye</span>
                            </div>

                            <div>
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">
                                    Límite de canales e-commerce <span className="font-normal text-gray-400">— 0 sin límite</span>
                                </label>
                                <input
                                    type="number"
                                    value={form.max_ecommerce_channels}
                                    onChange={(e) => setForm({ ...form, max_ecommerce_channels: e.target.value })}
                                    className="w-32 px-3 py-2.5 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-2 focus:ring-indigo-500/30 focus:border-indigo-500"
                                />
                            </div>

                            <div className="grid grid-cols-2 gap-3">
                                <div>
                                    <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">Envíos incluidos <span className="font-normal text-gray-400">/ mes</span></label>
                                    <input
                                        type="number"
                                        value={form.included_shipments}
                                        onChange={(e) => setForm({ ...form, included_shipments: e.target.value })}
                                        placeholder="ej: 100"
                                        className="w-full px-3 py-2.5 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-2 focus:ring-indigo-500/30 focus:border-indigo-500"
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">Costo envío extra</label>
                                    <div className="relative">
                                        <span className="absolute left-3 top-1/2 -translate-y-1/2 text-sm font-semibold text-gray-400">$</span>
                                        <input
                                            type="number"
                                            value={form.shipment_overage_price}
                                            onChange={(e) => setForm({ ...form, shipment_overage_price: e.target.value })}
                                            placeholder="ej: 600"
                                            className="w-full pl-7 pr-3 py-2.5 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-2 focus:ring-indigo-500/30 focus:border-indigo-500"
                                        />
                                    </div>
                                </div>
                            </div>

                            <div className="grid grid-cols-2 gap-3">
                                <div>
                                    <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">Facturas incluidas <span className="font-normal text-gray-400">/ mes</span></label>
                                    <input
                                        type="number"
                                        value={form.included_invoices}
                                        onChange={(e) => setForm({ ...form, included_invoices: e.target.value })}
                                        placeholder="ej: 6000"
                                        className="w-full px-3 py-2.5 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-2 focus:ring-indigo-500/30 focus:border-indigo-500"
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">Costo factura extra</label>
                                    <div className="relative">
                                        <span className="absolute left-3 top-1/2 -translate-y-1/2 text-sm font-semibold text-gray-400">$</span>
                                        <input
                                            type="number"
                                            value={form.invoice_overage_price}
                                            onChange={(e) => setForm({ ...form, invoice_overage_price: e.target.value })}
                                            placeholder="ej: 550"
                                            className="w-full pl-7 pr-3 py-2.5 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-2 focus:ring-indigo-500/30 focus:border-indigo-500"
                                        />
                                    </div>
                                </div>
                            </div>

                            <div className="grid grid-cols-2 gap-3">
                                <div>
                                    <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">{'\u00d3rdenes'} incluidas <span className="font-normal text-gray-400">/ mes</span></label>
                                    <input
                                        type="number"
                                        value={form.included_orders}
                                        onChange={(e) => setForm({ ...form, included_orders: e.target.value })}
                                        placeholder="ej: 6000"
                                        className="w-full px-3 py-2.5 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-2 focus:ring-indigo-500/30 focus:border-indigo-500"
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">Costo orden extra</label>
                                    <div className="relative">
                                        <span className="absolute left-3 top-1/2 -translate-y-1/2 text-sm font-semibold text-gray-400">$</span>
                                        <input
                                            type="number"
                                            value={form.order_overage_price}
                                            onChange={(e) => setForm({ ...form, order_overage_price: e.target.value })}
                                            placeholder="ej: 550"
                                            className="w-full pl-7 pr-3 py-2.5 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-2 focus:ring-indigo-500/30 focus:border-indigo-500"
                                        />
                                    </div>
                                </div>
                            </div>

                            <div>
                                <div className="flex items-center justify-between mb-2.5">
                                    <label className="text-sm font-semibold text-gray-700 dark:text-gray-300">Módulos incluidos</label>
                                    <span className="text-xs font-semibold text-indigo-600 dark:text-indigo-400">{form.module_codes.length} de {moduleCatalog.length}</span>
                                </div>
                                <div className="grid grid-cols-2 lg:grid-cols-3 gap-2">
                                    {moduleCatalog.map(({ code, name }) => {
                                        const checked = form.module_codes.includes(code);
                                        return (
                                            <div
                                                key={code}
                                                onClick={() => toggleModule(code)}
                                                className={`flex items-center gap-2 px-2.5 py-2 rounded-lg border-[1.5px] cursor-pointer transition-colors ${
                                                    checked
                                                        ? 'bg-indigo-50 dark:bg-indigo-900/20 border-indigo-500'
                                                        : 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-600'
                                                }`}
                                            >
                                                <span className={`w-4 h-4 rounded flex items-center justify-center flex-shrink-0 border-[1.5px] ${checked ? 'bg-indigo-600 border-indigo-600 text-white' : 'border-gray-300 dark:border-gray-600 text-transparent'}`}>
                                                    <CheckIcon />
                                                </span>
                                                <span className="text-[13px] font-medium text-gray-700 dark:text-gray-200 truncate">{name}</span>
                                            </div>
                                        );
                                    })}
                                </div>
                            </div>

                            {!editModal.plan && (
                                <Input label="Notas (opcional)" value={form.notes} onChange={(e) => setForm({ ...form, notes: e.target.value })} placeholder="Observaciones internas sobre este plan..." />
                            )}
                        </div>
                    </div>

                    <div className="flex justify-end gap-2.5 pt-4 mt-4 border-t border-gray-100 dark:border-gray-700">
                        <Button variant="secondary" onClick={() => setEditModal({ open: false })}>Cancelar</Button>
                        <Button variant="purple" onClick={handleSave}>Guardar</Button>
                    </div>
                </div>
            </Modal>
        </div>
    );
}

function OverridesPanel({ businessId, businessName }: { businessId: number; businessName?: string }) {
    const [overrides, setOverrides] = useState<BusinessModuleOverride[]>([]);
    const [moduleCatalog, setModuleCatalog] = useState<ModuleInfo[]>([]);
    const [loading, setLoading] = useState(true);
    const [selectedCode, setSelectedCode] = useState('');
    const [notes, setNotes] = useState('');
    const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

    const moduleName = (code: string) => moduleCatalog.find((m) => m.code === code)?.name ?? code;

    const load = useCallback(async () => {
        setLoading(true);
        const [overridesRes, catalogRes] = await Promise.all([listOverridesAction(businessId), getModuleCatalogAction()]);
        if (overridesRes.success && overridesRes.data) setOverrides(overridesRes.data);
        if (catalogRes.success && catalogRes.data) setModuleCatalog(catalogRes.data);
        setLoading(false);
    }, [businessId]);

    useEffect(() => { load(); }, [load]);

    const handleGrant = async () => {
        if (!selectedCode) return;
        const res = await grantOverrideAction({ businessId, moduleCode: selectedCode, notes: notes || undefined });
        if (res.success) {
            setMessage({ type: 'success', text: `Módulo "${selectedCode}" habilitado para ${businessName}` });
            setSelectedCode(''); setNotes('');
            load();
        } else {
            setMessage({ type: 'error', text: res.error || 'Error al otorgar el módulo' });
        }
    };

    const handleRevoke = async (moduleCode: string) => {
        const res = await revokeOverrideAction(businessId, moduleCode);
        if (res.success) {
            setMessage({ type: 'success', text: `Módulo "${moduleCode}" revocado` });
            load();
        } else {
            setMessage({ type: 'error', text: res.error || 'Error al revocar el módulo' });
        }
    };

    if (loading) return <div className="flex justify-center py-6"><Spinner /></div>;

    return (
        <div className="bg-white dark:bg-gray-700 rounded-xl border border-gray-200 dark:border-gray-600 shadow-sm p-6 space-y-4">
            <h3 className="text-lg font-bold text-gray-900 dark:text-white">Módulos adicionales — {businessName}</h3>
            <p className="text-sm text-gray-500 dark:text-gray-400">
                Otorga acceso a un módulo puntual para este negocio, independiente de su plan actual.
            </p>

            {message && <Alert type={message.type} onClose={() => setMessage(null)}>{message.text}</Alert>}

            <div className="flex flex-wrap gap-2">
                {overrides.length === 0 && <span className="text-sm text-gray-400">Sin módulos adicionales otorgados</span>}
                {overrides.map((o) => (
                    <span key={o.id} className="inline-flex items-center gap-2 text-xs bg-violet-50 dark:bg-violet-900/30 text-violet-700 dark:text-violet-300 px-3 py-1.5 rounded-full">
                        {moduleName(o.module_code)}
                        <button onClick={() => handleRevoke(o.module_code)} className="text-violet-400 hover:text-violet-700">×</button>
                    </span>
                ))}
            </div>

            <div className="flex flex-col sm:flex-row gap-2 items-start sm:items-end pt-3 border-t border-gray-100 dark:border-gray-600">
                <div className="flex-1 w-full">
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Módulo</label>
                    <select value={selectedCode} onChange={(e) => setSelectedCode(e.target.value)} className="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white">
                        <option value="">Selecciona un módulo</option>
                        {moduleCatalog.filter((m) => !overrides.some((o) => o.module_code === m.code)).map(({ code, name }) => (
                            <option key={code} value={code}>{name}</option>
                        ))}
                    </select>
                </div>
                <div className="flex-1 w-full">
                    <Input label="Notas (opcional)" value={notes} onChange={(e) => setNotes(e.target.value)} />
                </div>
                <Button variant="primary" onClick={handleGrant} disabled={!selectedCode}>Otorgar</Button>
            </div>
        </div>
    );
}

interface BusinessSubscriptionViewProps {
    businessId?: number;
    businessName?: string;
    isSuperAdminView?: boolean;
}

function useCountUp(value: number, active: boolean, durationMs = 700) {
    const [display, setDisplay] = useState(active ? 0 : value);

    useEffect(() => {
        if (!active) { setDisplay(value); return; }
        let raf = 0;
        const start = performance.now();
        const from = 0;
        const tick = (now: number) => {
            const t = Math.min(1, (now - start) / durationMs);
            const eased = 1 - Math.pow(1 - t, 3);
            setDisplay(Math.round(from + (value - from) * eased));
            if (t < 1) raf = requestAnimationFrame(tick);
        };
        raf = requestAnimationFrame(tick);
        return () => cancelAnimationFrame(raf);
    }, [value, active, durationMs]);

    return display;
}

function UsageLimitRow({ label, used, included, overagePrice, unit, animate }: {
    label: string; used: number; included: number; overagePrice?: number; unit: string; animate: boolean;
}) {
    const displayUsed = useCountUp(used, animate);
    const pct = included > 0 ? Math.min(100, Math.round((used / included) * 100)) : 0;
    const over = used > included;
    const extra = Math.max(0, used - included);
    const barColor = over ? 'bg-amber-500' : pct >= 80 ? 'bg-amber-400' : 'bg-violet-500';

    return (
        <div>
            <div className="flex items-center justify-between text-xs mb-1">
                <span className="font-medium text-gray-700 dark:text-gray-200">{label}</span>
                <span className={over ? 'font-semibold text-amber-600 dark:text-amber-400' : 'text-gray-400'}>
                    {displayUsed} / {included}
                </span>
            </div>
            <div className="w-full h-2 rounded-full bg-gray-100 dark:bg-gray-700 overflow-hidden">
                <div
                    className={`h-full rounded-full transition-all duration-700 ease-out ${barColor}`}
                    style={{ width: `${animate ? pct : 0}%` }}
                />
            </div>
            {over && overagePrice != null && (
                <p className="text-[11px] text-amber-600 dark:text-amber-400 mt-1">
                    {extra} {unit}{extra === 1 ? '' : 's'} de más · {formatCurrency(overagePrice)} c/u
                </p>
            )}
        </div>
    );
}

function BusinessSubscriptionView({ businessId, businessName, isSuperAdminView }: BusinessSubscriptionViewProps = {}) {
    const [subscription, setSubscription] = useState<BusinessSubscription | null>(null);
    const [usage, setUsage] = useState<SubscriptionUsage | null>(null);
    const [moduleCatalog, setModuleCatalog] = useState<ModuleInfo[]>([]);
    const [loading, setLoading] = useState(true);
    const [mounted, setMounted] = useState(false);
    const [savingAutoPayment, setSavingAutoPayment] = useState(false);
    const planCatalogRef = useRef<PlanCatalogHandle>(null);

    const fetchSub = useCallback(async () => {
        setLoading(true);
        setMounted(false);
        const [subRes, usageRes, catalogRes] = await Promise.all([
            getMySubscriptionAction(businessId),
            getMySubscriptionUsageAction(businessId),
            getModuleCatalogAction(),
        ]);
        if (subRes.success && subRes.data) setSubscription(subRes.data);
        if (usageRes.success && usageRes.data) setUsage(usageRes.data);
        if (catalogRes.success && catalogRes.data) setModuleCatalog(catalogRes.data);
        setLoading(false);
    }, [businessId]);

    useEffect(() => { fetchSub(); }, [fetchSub]);
    useEffect(() => {
        if (loading) return;
        const t = setTimeout(() => setMounted(true), 30);
        return () => clearTimeout(t);
    }, [loading]);

    if (loading) return <div className="flex justify-center py-12"><Spinner /></div>;

    const isExpired = !subscription
        || ['pending', 'expired', 'cancelled'].includes(subscription.business_subscription_status || subscription.status);
    const moduleName = (code: string) => moduleCatalog.find((m) => m.code === code)?.name ?? code;

    const limitRows = usage
        ? [
            { label: 'Envíos', used: usage.shipments_used, included: usage.included_shipments, overagePrice: usage.shipment_overage_price, unit: 'envío' },
            { label: 'Facturas', used: usage.invoices_used, included: usage.included_invoices, overagePrice: usage.invoice_overage_price, unit: 'factura' },
            { label: 'Órdenes', used: usage.orders_used, included: usage.included_orders, overagePrice: usage.order_overage_price, unit: 'orden' },
        ].filter((r): r is typeof r & { included: number } => r.included != null)
        : [];

    const hasOverage = !!(usage?.forecasted_payment != null && subscription && usage.forecasted_payment > subscription.amount);

    const canPayNow = isExpired
        || !subscription?.payment_window_start
        || new Date(subscription.payment_window_start) <= new Date();

    const handleToggleAutoPayment = async () => {
        if (!subscription || savingAutoPayment) return;
        setSavingAutoPayment(true);
        const next = !subscription.auto_payment_enabled;
        const res = await setAutoPaymentAction(next, businessId);
        setSavingAutoPayment(false);
        if (res.success) {
            setSubscription({ ...subscription, auto_payment_enabled: next });
        }
    };

    return (
        <div className="space-y-6">
            {isSuperAdminView && (
                <div className="flex items-center gap-3 bg-blue-50 dark:bg-blue-900/30 border border-blue-200 dark:border-blue-700 rounded-lg px-4 py-3">
                    <p className="text-sm text-blue-800 dark:text-blue-300">
                        Vista de suscripción de <strong>{businessName}</strong> — modo super admin
                    </p>
                </div>
            )}

            <div className={`rounded-2xl p-6 text-white ${isExpired ? 'bg-gradient-to-br from-red-500 to-red-700' : 'bg-gradient-to-br from-violet-600 to-purple-800'} shadow-lg transition-all duration-500 ${mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-2'}`}>
                <div className="flex items-center justify-between mb-4">
                    <div>
                        <p className="text-white/70 text-sm font-medium uppercase tracking-wider">Estado de Suscripción</p>
                        <h2 className="text-2xl font-bold mt-1">{isExpired ? 'Suspendida' : 'Activa'}</h2>
                        {subscription?.subscription_type_name && (
                            <p className="text-white/80 text-sm mt-1">Plan: {subscription.subscription_type_name}</p>
                        )}
                    </div>
                    {subscription && (
                        <Button
                            id="subscription-pay-button"
                            variant="secondary"
                            className={`flex-shrink-0 ${canPayNow ? '!bg-white !text-violet-700 hover:!bg-white/90' : '!bg-white/40 !text-violet-700/60 !cursor-not-allowed hover:!bg-white/40'}`}
                            disabled={!canPayNow}
                            onClick={() => planCatalogRef.current?.openCurrentPlanPurchase()}
                        >
                            {isExpired ? 'Renovar ahora' : 'Pagar / Extender'}
                        </Button>
                    )}
                </div>
                {subscription && (
                    <div className="grid grid-cols-2 gap-4 mt-4 pt-4 border-t border-white/20">
                        <div>
                            <p className="text-white/60 text-xs">Último pago</p>
                            <p className="font-semibold">{formatCurrency(subscription.amount)}</p>
                        </div>
                        <div>
                            <p className="text-white/60 text-xs">Válida hasta</p>
                            <p className="font-semibold">{formatDate(subscription.end_date)}</p>
                        </div>
                    </div>
                )}
                {subscription && subscription.payment_window_start && subscription.payment_window_end && (
                    <p id="subscription-payment-window" className="text-white/60 text-[11px] mt-2">
                        {canPayNow
                            ? `Podrás pagar la membresía hasta el ${formatDate(subscription.payment_window_end)} antes de que se suspenda la cuenta.`
                            : `Periodo facturado hasta el ${formatDate(subscription.payment_window_start)}. Podrás pagar la membresía del ${formatDate(subscription.payment_window_start)} al ${formatDate(subscription.payment_window_end)}.`}
                    </p>
                )}
                {subscription && !isExpired && (
                    <div className="mt-4 pt-4 border-t border-white/20">
                        <div className="flex items-center justify-between text-[11px] text-white/70 mb-1">
                            <span>Ciclo actual</span>
                        </div>
                        <div className="w-full h-1.5 rounded-full bg-white/20 overflow-hidden">
                            <div
                                className="h-full rounded-full bg-white transition-all duration-700 ease-out"
                                style={{ width: `${mounted ? (computeMembershipProgress(subscription.start_date, subscription.end_date)?.percent ?? 0) : 0}%` }}
                            />
                        </div>
                    </div>
                )}
            </div>

            {subscription && (
                <div id="subscription-auto-payment-toggle" className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4 flex items-center justify-between gap-4">
                    <div>
                        <p className="text-sm font-semibold text-gray-900 dark:text-white">Pago automático de la suscripción</p>
                        <p className="text-xs text-gray-500 dark:text-gray-400 mt-0.5">
                            Si hay saldo suficiente en tu billetera, la suscripción se paga sola el día que vence, sin necesidad de que la pagues manualmente.
                            Se revisa todos los días a las 8:00 a. m.
                        </p>
                    </div>
                    <button
                        type="button"
                        role="switch"
                        aria-checked={!!subscription.auto_payment_enabled}
                        disabled={savingAutoPayment}
                        onClick={handleToggleAutoPayment}
                        className={`relative inline-flex h-6 w-11 flex-shrink-0 items-center rounded-full transition-colors disabled:opacity-60 ${subscription.auto_payment_enabled ? 'bg-violet-600' : 'bg-gray-300 dark:bg-gray-600'}`}
                    >
                        <span className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${subscription.auto_payment_enabled ? 'translate-x-6' : 'translate-x-1'}`} />
                    </button>
                </div>
            )}

            {usage && !isExpired && (limitRows.length > 0 || usage.module_codes.length > 0) && (
                <div className={`grid gap-4 md:grid-cols-2 transition-all duration-500 delay-100 ${mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-2'}`}>
                    {usage.module_codes.length > 0 && (
                        <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 shadow-sm p-5">
                            <h3 className="text-sm font-bold text-gray-900 dark:text-white mb-3">Módulos incluidos</h3>
                            <div className="grid grid-cols-2 gap-2.5">
                                {usage.module_codes.map((m) => (
                                    <div key={m} className="flex items-center gap-2">
                                        <span className="w-5 h-5 rounded-full bg-violet-100 dark:bg-violet-800/40 text-violet-600 dark:text-violet-300 flex items-center justify-center flex-shrink-0">
                                            <CheckIcon />
                                        </span>
                                        <span className="text-xs font-medium text-gray-700 dark:text-gray-200">{moduleName(m)}</span>
                                    </div>
                                ))}
                            </div>
                        </div>
                    )}

                    {limitRows.length > 0 && (
                        <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 shadow-sm p-5">
                            <h3 className="text-sm font-bold text-gray-900 dark:text-white mb-4">Uso de este ciclo</h3>
                            <div className="space-y-4">
                                {limitRows.map((r) => (
                                    <UsageLimitRow key={r.label} {...r} animate={mounted} />
                                ))}
                            </div>

                            {hasOverage && usage.forecasted_payment != null && (
                                <div className="mt-4 pt-4 border-t border-gray-100 dark:border-gray-700 flex items-center justify-between">
                                    <span className="text-xs text-gray-500 dark:text-gray-400">Pago pronosticado</span>
                                    <span className="text-sm font-bold text-amber-600 dark:text-amber-400">{formatCurrency(usage.forecasted_payment)}</span>
                                </div>
                            )}
                        </div>
                    )}
                </div>
            )}

            {usage && !isExpired && (
                <div className={`bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 shadow-sm p-5 transition-all duration-500 delay-150 ${mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-2'}`}>
                    <h3 className="text-sm font-bold text-gray-900 dark:text-white mb-4">Detalle del plan</h3>
                    <div className="flex flex-wrap items-center gap-x-8 gap-y-4">
                        <div>
                            <p className="text-[11px] text-gray-400 uppercase tracking-wide mb-0.5">Plan</p>
                            <p className="text-base font-bold text-gray-900 dark:text-white">{usage.plan_name}</p>
                        </div>
                        <div>
                            <p className="text-[11px] text-gray-400 uppercase tracking-wide mb-0.5">Precio</p>
                            <p className="text-base font-bold text-gray-900 dark:text-white">
                                {formatCurrency(usage.plan_price)}
                                <span className="text-xs font-medium text-gray-400"> /{usage.billing_period === 'monthly' ? 'mes' : 'año'}</span>
                            </p>
                        </div>
                        <div className="flex items-center gap-2.5">
                            <span className="w-9 h-9 rounded-lg bg-violet-100 dark:bg-violet-800/40 text-violet-600 dark:text-violet-300 flex items-center justify-center flex-shrink-0">
                                <ChannelsIcon />
                            </span>
                            <div>
                                <p className="text-sm font-bold text-gray-900 dark:text-white leading-tight">
                                    {usage.max_ecommerce_channels > 0 ? `Hasta ${usage.max_ecommerce_channels}` : 'Ilimitados'}
                                </p>
                                <p className="text-[11px] text-gray-400 leading-tight">canales de ecommerce conectados</p>
                            </div>
                        </div>
                    </div>
                </div>
            )}

            <PlanCatalog ref={planCatalogRef} businessId={businessId} onPurchased={fetchSub} currentSubscription={subscription} isCurrentActive={!isExpired} usage={usage} />
        </div>
    );
}

const CheckIcon = () => (
    <svg width="11" height="11" viewBox="0 0 24 24" fill="none">
        <path d="M4 12.5L9.5 18L20 6.5" stroke="currentColor" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
);

const StarIcon = () => (
    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" className="flex-shrink-0">
        <path d="M12 2L14.6 8.4L21.5 9L16.2 13.5L18 20.5L12 16.7L6 20.5L7.8 13.5L2.5 9L9.4 8.4L12 2Z" fill="currentColor" />
    </svg>
);

const ChannelsIcon = () => (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
        <path d="M3 3H21V8H3V3Z" stroke="currentColor" strokeWidth="1.8" strokeLinejoin="round" />
        <path d="M4 8V20C4 20.5523 4.44772 21 5 21H19C19.5523 21 20 20.5523 20 20V8" stroke="currentColor" strokeWidth="1.8" strokeLinejoin="round" />
        <path d="M9 12H15" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" />
    </svg>
);

interface PlanCatalogProps {
    businessId?: number;
    onPurchased: () => void;
    currentSubscription: BusinessSubscription | null;
    isCurrentActive: boolean;
    usage: SubscriptionUsage | null;
}

export interface PlanCatalogHandle {
    openCurrentPlanPurchase: () => void | Promise<void>;
}

const PlanCatalog = forwardRef<PlanCatalogHandle, PlanCatalogProps>(function PlanCatalog({ businessId, onPurchased, currentSubscription, isCurrentActive, usage }, ref) {
    const [types, setTypes] = useState<SubscriptionType[]>([]);
    const [moduleCatalog, setModuleCatalog] = useState<ModuleInfo[]>([]);
    const [loading, setLoading] = useState(true);
    const [purchaseModal, setPurchaseModal] = useState<{ open: boolean; type?: SubscriptionType }>({ open: false });
    const [months, setMonths] = useState('1');
    const [buying, setBuying] = useState(false);
    const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
    const [walletBalance, setWalletBalance] = useState<number | null>(null);
    const [loadingBalance, setLoadingBalance] = useState(false);
    const [rechargingBold, setRechargingBold] = useState(false);
    const [boldProcessing, setBoldProcessing] = useState<{ orderId: string; amount: number; pollingEnabled: boolean } | null>(null);

    const moduleName = (code: string) => moduleCatalog.find((m) => m.code === code)?.name ?? code;

    useEffect(() => {
        Promise.all([listSubscriptionTypesAction(true), getModuleCatalogAction()]).then(([typesRes, catalogRes]) => {
            if (typesRes.success && typesRes.data) setTypes(typesRes.data);
            if (catalogRes.success && catalogRes.data) setModuleCatalog(catalogRes.data);
            setLoading(false);
        });
    }, []);

    const fetchBalance = useCallback(async () => {
        setLoadingBalance(true);
        const res = await getWalletBalanceAction(businessId);
        if (res.success && res.data) setWalletBalance(res.data.Balance);
        setLoadingBalance(false);
    }, [businessId]);

    const openPurchaseModal = (t: SubscriptionType) => {
        setMonths('1');
        setPurchaseModal({ open: true, type: t });
        fetchBalance();
    };

    useImperativeHandle(ref, () => ({
        openCurrentPlanPurchase: async () => {
            const targetId = currentSubscription?.subscription_type_id;
            if (!targetId) return;
            const current = currentSubscription?.subscription_type ?? types.find((t) => t.id === targetId);
            if (current) {
                openPurchaseModal(current);
            } else {
                setMessage({ type: 'error', text: 'No se pudo cargar el plan actual. Intenta de nuevo.' });
            }
        },
    }));

    const handleBuy = async () => {
        if (!purchaseModal.type) return;
        setBuying(true);
        const res = await purchaseSubscriptionAction({ subscriptionTypeId: purchaseModal.type.id, months: Number(months) }, businessId);
        setBuying(false);
        if (res.success) {
            setMessage({ type: 'success', text: `Suscripción "${purchaseModal.type.name}" activada correctamente.` });
            setPurchaseModal({ open: false });
            onPurchased();
        } else {
            setMessage({ type: 'error', text: res.error?.includes('insufficient') ? 'Saldo insuficiente en tu billetera. Recárgala e intenta de nuevo.' : (res.error || 'Error al procesar la compra') });
        }
    };

    const handleRechargeBold = async (amount: number) => {
        setRechargingBold(true);
        try {
            const res = await getBoldSignatureAction(Math.ceil(amount), businessId);
            if (!res?.success) {
                throw new Error(res?.message || 'Error al obtener firma de Bold');
            }

            const { order_id, currency, amount: boldAmount, hash, public_key, redirection_url, polling_enabled } = res.data;

            if (!window.hasOwnProperty('BoldCheckout')) {
                await new Promise<void>((resolve, reject) => {
                    const script = document.createElement('script');
                    script.src = 'https://checkout.bold.co/library/boldPaymentButton.js';
                    script.async = true;
                    script.onload = () => resolve();
                    script.onerror = () => reject(new Error('No se pudo cargar el script de Bold'));
                    document.body.appendChild(script);
                    setTimeout(() => reject(new Error('Timeout cargando script de Bold')), 10000);
                });
            }

            const checkoutConfig: Record<string, unknown> = {
                orderId: order_id,
                currency,
                amount: boldAmount,
                apiKey: public_key,
                integritySignature: hash,
                description: `Recarga para suscripción - Orden ${order_id}`,
            };
            if (redirection_url) checkoutConfig.redirectionUrl = redirection_url;

            // @ts-expect-error BoldCheckout is loaded from external script
            const checkout = new BoldCheckout(checkoutConfig);
            checkout.open();
            setBoldProcessing({ orderId: order_id, amount: boldAmount, pollingEnabled: !!polling_enabled });
        } catch (err: any) {
            setMessage({ type: 'error', text: err.message || 'Error al iniciar el pago con Bold' });
        } finally {
            setRechargingBold(false);
        }
    };

    if (loading) return <div className="flex justify-center py-6"><Spinner /></div>;

    const sorted = [...types].sort((a, b) => a.price - b.price);
    const currentIndex = isCurrentActive ? sorted.findIndex((t) => t.id === currentSubscription?.subscription_type_id) : -1;
    const featuredIndex = sorted.length >= 3 ? Math.floor((sorted.length - 1) / 2) : -1;

    return (
        <div className="bg-white dark:bg-gray-700 rounded-xl border border-gray-200 dark:border-gray-600 shadow-sm p-6 space-y-5">
            <div>
                <h3 className="text-lg font-bold text-gray-900 dark:text-white">Planes disponibles</h3>
                <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
                    Compra o cambia de plan pagando con el saldo de tu billetera. El cambio se aplica de inmediato.
                </p>
            </div>

            {message && <Alert type={message.type} onClose={() => setMessage(null)}>{message.text}</Alert>}

            <div className="grid gap-5 lg:grid-cols-3">
                {sorted.map((t, i) => {
                    const featured = i === featuredIndex;
                    const isCurrent = i === currentIndex;
                    const isUpgrade = currentIndex >= 0 && i > currentIndex;
                    const isDowngrade = currentIndex >= 0 && i < currentIndex;

                    let ctaLabel = 'Cambiar a este plan';
                    let footnote: string | null = null;
                    if (currentIndex < 0) { ctaLabel = 'Comprar'; }
                    else if (isUpgrade) { footnote = 'Se prorratea en tu próxima factura'; }
                    else if (isDowngrade) { footnote = 'Aplica al final de tu ciclo actual'; }

                    return (
                        <div
                            key={t.id}
                            className={`relative flex flex-col rounded-2xl border bg-white dark:bg-gray-800 overflow-hidden transition-transform ${
                                featured
                                    ? 'border-violet-500 shadow-lg shadow-violet-500/10 lg:-translate-y-1.5 z-10'
                                    : 'border-gray-200 dark:border-gray-600 shadow-sm'
                            }`}
                        >
                            {featured && (
                                <div className="flex items-center justify-center gap-1.5 bg-gradient-to-r from-violet-600 to-purple-500 text-white text-[11px] font-bold uppercase tracking-wide py-1.5">
                                    <StarIcon /> Más elegido
                                </div>
                            )}

                            <div className={featured ? 'pt-5 px-6 pb-5' : 'pt-6 px-6 pb-5'}>
                                <div className="flex items-center justify-between mb-3">
                                    <h4 className="text-lg font-bold text-gray-900 dark:text-white">{t.name}</h4>
                                    {isCurrent && (
                                        <span className="text-[11px] font-bold text-violet-600 dark:text-violet-300 bg-violet-50 dark:bg-violet-900/30 px-2 py-0.5 rounded-full">
                                            Tu plan
                                        </span>
                                    )}
                                </div>

                                <div className="flex items-end gap-1 mb-3">
                                    <span className="text-3xl font-extrabold tracking-tight text-gray-900 dark:text-white">{formatCurrency(t.price)}</span>
                                    <span className="text-sm font-semibold text-gray-400 pb-1">/{t.billing_period === 'monthly' ? 'mes' : 'año'}</span>
                                </div>

                                {t.description && (
                                    <p className="text-sm text-gray-500 dark:text-gray-400 leading-relaxed mb-4 min-h-[2.6rem]">{t.description}</p>
                                )}

                                <div className={`flex items-center gap-3 rounded-xl p-3 ${featured ? 'bg-violet-50 dark:bg-violet-900/20' : 'bg-gray-50 dark:bg-gray-700/40'}`}>
                                    <div className={`w-[34px] h-[34px] rounded-lg flex items-center justify-center flex-shrink-0 ${
                                        featured ? 'bg-violet-100 dark:bg-violet-800/40 text-violet-600 dark:text-violet-300' : 'bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400'
                                    }`}>
                                        <ChannelsIcon />
                                    </div>
                                    <div>
                                        <div className="text-sm font-bold text-gray-900 dark:text-white leading-tight">
                                            {t.max_ecommerce_channels > 0 ? `Hasta ${t.max_ecommerce_channels}` : 'Ilimitados'}
                                        </div>
                                        <div className="text-xs text-gray-400 leading-tight">canales de ecommerce conectados</div>
                                    </div>
                                </div>
                            </div>

                            <div className="border-t border-gray-100 dark:border-gray-700 mx-6" />

                            <div className="px-6 pt-4 pb-2 flex-1">
                                <div className="text-[11px] font-bold uppercase tracking-wide text-gray-400 mb-3">Módulos incluidos</div>
                                <div className="flex flex-col gap-2.5">
                                    {(t.module_codes ?? []).map((m) => (
                                        <div key={m} className="flex items-center gap-2.5">
                                            <span className={`w-[18px] h-[18px] rounded-full flex items-center justify-center flex-shrink-0 ${
                                                featured ? 'bg-violet-100 dark:bg-violet-800/40 text-violet-600 dark:text-violet-300' : 'bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400'
                                            }`}>
                                                <CheckIcon />
                                            </span>
                                            <span className="text-[13px] font-medium text-gray-700 dark:text-gray-200">{moduleName(m)}</span>
                                        </div>
                                    ))}
                                </div>
                            </div>

                            <div className="px-6 pt-5 pb-6">
                                {isCurrent ? (
                                    <button disabled className="w-full py-3 rounded-lg border border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 text-gray-400 dark:text-gray-500 font-semibold text-sm cursor-default">
                                        Plan actual
                                    </button>
                                ) : (
                                    <Button
                                        variant={featured ? 'purple' : 'outline-purple'}
                                        className="w-full"
                                        onClick={() => openPurchaseModal(t)}
                                    >
                                        {ctaLabel}
                                    </Button>
                                )}
                                {footnote && (
                                    <div className="text-[11px] text-gray-400 text-center mt-2">{footnote}</div>
                                )}
                            </div>
                        </div>
                    );
                })}
            </div>

            <div className="flex items-center justify-center gap-2 text-xs text-gray-400 pt-1">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" className="flex-shrink-0">
                    <circle cx="12" cy="12" r="9" stroke="currentColor" strokeWidth="1.8" />
                    <path d="M12 8V13" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" />
                    <circle cx="12" cy="16.2" r="0.9" fill="currentColor" />
                </svg>
                El cambio de plan se aplica de forma inmediata y se prorratea en tu próxima factura.
            </div>

            <Modal isOpen={purchaseModal.open} onClose={() => setPurchaseModal({ open: false })} size="sm">
                {purchaseModal.type && (() => {
                    const pendingOverage = usage?.forecasted_payment != null
                        ? Math.max(0, usage.forecasted_payment - usage.plan_price)
                        : 0;
                    const total = purchaseModal.type!.price * Number(months) + pendingOverage;
                    const hasBalance = walletBalance !== null;
                    const sufficient = hasBalance && walletBalance! >= total;
                    const missing = hasBalance ? Math.max(0, total - walletBalance!) : 0;

                    return (
                        <div className="flex flex-col">
                            <div className="bg-gradient-to-br from-violet-600 to-purple-800 rounded-2xl px-5 py-4 text-white shadow-lg shadow-violet-900/25 mb-5">
                                <div className="flex items-center gap-3">
                                    <div className="w-11 h-11 rounded-xl bg-white/15 flex items-center justify-center flex-shrink-0">
                                        <svg width="22" height="22" viewBox="0 0 24 24" fill="none">
                                            <path d="M3 3H21V8H3V3Z" stroke="currentColor" strokeWidth="1.8" strokeLinejoin="round" />
                                            <path d="M4 8V20C4 20.5523 4.44772 21 5 21H19C19.5523 21 20 20.5523 20 20V8" stroke="currentColor" strokeWidth="1.8" strokeLinejoin="round" />
                                            <path d="M9 12H15" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" />
                                        </svg>
                                    </div>
                                    <div>
                                        <p className="text-white/70 text-xs font-medium uppercase tracking-wide">Contratar plan</p>
                                        <h3 className="text-lg font-bold leading-tight">{purchaseModal.type!.name}</h3>
                                    </div>
                                </div>
                            </div>

                            <div className="space-y-4">
                                <div>
                                    <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">Duración</label>
                                    <select value={months} onChange={(e) => setMonths(e.target.value)} className="w-full px-3 py-2.5 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-2 focus:ring-violet-500/30 focus:border-violet-500">
                                        <option value="1">1 mes</option>
                                        <option value="3">3 meses</option>
                                        <option value="6">6 meses</option>
                                        <option value="12">12 meses (anual)</option>
                                    </select>
                                </div>

                                <div className="rounded-xl border border-gray-100 dark:border-gray-700 overflow-hidden">
                                    <div className="flex items-center justify-between px-4 py-2.5 text-sm">
                                        <span className="text-gray-500 dark:text-gray-400">{formatCurrency(purchaseModal.type!.price)} × {months} {Number(months) === 1 ? 'mes' : 'meses'}</span>
                                        <span className="font-medium text-gray-700 dark:text-gray-200">{formatCurrency(purchaseModal.type!.price * Number(months))}</span>
                                    </div>
                                    {pendingOverage > 0 && (
                                        <div className="flex items-center justify-between px-4 py-2.5 text-sm border-t border-gray-100 dark:border-gray-700">
                                            <span className="text-amber-600 dark:text-amber-400">{'Excedente del ciclo actual'}</span>
                                            <span className="font-medium text-amber-600 dark:text-amber-400">{formatCurrency(pendingOverage)}</span>
                                        </div>
                                    )}
                                    <div className="flex items-center justify-between px-4 py-3 bg-violet-50 dark:bg-violet-900/20 border-t border-gray-100 dark:border-gray-700">
                                        <span className="text-sm font-bold text-gray-900 dark:text-white">Total a pagar</span>
                                        <span className="text-lg font-extrabold text-violet-700 dark:text-violet-300">{formatCurrency(total)}</span>
                                    </div>
                                </div>

                                <div className={`rounded-xl px-4 py-3.5 flex items-center gap-3 ${
                                    loadingBalance
                                        ? 'bg-gray-50 dark:bg-gray-700/40'
                                        : sufficient
                                            ? 'bg-green-50 dark:bg-green-900/20'
                                            : 'bg-amber-50 dark:bg-amber-900/20'
                                }`}>
                                    {loadingBalance ? (
                                        <>
                                            <Spinner size="sm" />
                                            <span className="text-sm text-gray-500 dark:text-gray-400">Consultando saldo de tu billetera...</span>
                                        </>
                                    ) : sufficient ? (
                                        <>
                                            <span className="w-8 h-8 rounded-full bg-green-100 dark:bg-green-800/40 text-green-600 dark:text-green-300 flex items-center justify-center flex-shrink-0">
                                                <CheckIcon />
                                            </span>
                                            <div className="text-sm">
                                                <div className="font-semibold text-green-800 dark:text-green-300">Saldo suficiente</div>
                                                <div className="text-green-700/80 dark:text-green-400/80 text-xs mt-0.5">
                                                    Saldo actual {formatCurrency(walletBalance!)} — te quedarán {formatCurrency(walletBalance! - total)} después de la compra.
                                                </div>
                                            </div>
                                        </>
                                    ) : (
                                        <>
                                            <span className="w-8 h-8 rounded-full bg-amber-100 dark:bg-amber-800/40 text-amber-600 dark:text-amber-300 flex items-center justify-center flex-shrink-0">
                                                <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M12 9v4m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" /></svg>
                                            </span>
                                            <div className="text-sm">
                                                <div className="font-semibold text-amber-800 dark:text-amber-300">Saldo insuficiente</div>
                                                <div className="text-amber-700/80 dark:text-amber-400/80 text-xs mt-0.5">
                                                    Tienes {formatCurrency(walletBalance ?? 0)}, te faltan {formatCurrency(missing)} para completar la compra.
                                                </div>
                                            </div>
                                        </>
                                    )}
                                </div>

                                <div className="flex justify-end gap-2.5 pt-1">
                                    <Button variant="outline-purple" onClick={() => setPurchaseModal({ open: false })}>Cancelar</Button>
                                    {!loadingBalance && !sufficient ? (
                                        <Button variant="purple" onClick={() => handleRechargeBold(missing)} loading={rechargingBold}>
                                            Recargar con Bold
                                        </Button>
                                    ) : (
                                        <Button variant="purple" onClick={handleBuy} loading={buying} disabled={loadingBalance}>Confirmar Compra</Button>
                                    )}
                                </div>
                            </div>
                        </div>
                    );
                })()}
            </Modal>

            <BoldPaymentProcessingModal
                open={boldProcessing !== null}
                orderId={boldProcessing?.orderId || ''}
                amount={boldProcessing?.amount || 0}
                businessId={businessId}
                pollingEnabled={boldProcessing?.pollingEnabled ?? false}
                onClose={() => {
                    setBoldProcessing(null);
                    fetchBalance();
                }}
                onResolved={(status) => {
                    if (status === 'success') {
                        fetchBalance();
                        setMessage({ type: 'success', text: 'Recarga confirmada por Bold. Ya puedes completar la compra.' });
                    } else if (status === 'failed') {
                        setMessage({ type: 'error', text: 'Bold rechazó el pago de recarga.' });
                    } else if (status === 'timeout') {
                        setMessage({ type: 'error', text: 'La recarga sigue en proceso. Revisa tu billetera en unos minutos.' });
                    }
                }}
            />
        </div>
    );
});
PlanCatalog.displayName = 'PlanCatalog';
