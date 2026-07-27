'use client';

import { useState, useEffect, useCallback } from 'react';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { Spinner, Button, Modal, Alert, Input } from '@/shared/ui';
import {
    getMySubscriptionAction,
    registerSubscriptionPaymentAction,
    editSubscriptionDatesAction,
    disableSubscriptionAction,
    listSubscriptionTypesAction,
    createSubscriptionTypeAction,
    updateSubscriptionTypeAction,
    deleteSubscriptionTypeAction,
    getModuleCatalogAction,
    purchaseSubscriptionAction,
    listOverridesAction,
    grantOverrideAction,
    revokeOverrideAction,
    listCustomPlansAction,
    createCustomPlanAction,
    updateCustomPlanAction,
    deleteCustomPlanAction,
    BusinessSubscription,
    SubscriptionType,
    BusinessModuleOverride,
    ModuleInfo,
} from '@/services/modules/wallet/infra/subscription-actions';
import { getWalletBalanceAction } from '@/services/modules/wallet/infra/actions';
import { getBoldSignatureAction } from '@/services/modules/pay/infra/actions';
import { BoldPaymentProcessingModal } from '@/app/(auth)/wallet/bold-payment-processing-modal';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';

const formatCurrency = (amount: number) =>
    new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP' }).format(amount);

const toDateInputValue = (dateStr?: string) => {
    if (!dateStr) return '';
    return dateStr.slice(0, 10);
};

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
                                {adminTab === 'businesses' && <AdminSubscriptionsView businesses={businesses} />}
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

function AdminSubscriptionsView({ businesses }: { businesses: Array<{ id: number; name: string }> }) {
    const [filter, setFilter] = useState('all');
    const [search, setSearch] = useState('');
    const [registerModal, setRegisterModal] = useState<{ open: boolean; business?: { id: number; name: string } }>({ open: false });
    const [editDatesModal, setEditDatesModal] = useState<{ open: boolean; business?: { id: number; name: string } }>({ open: false });
    const [editStartDate, setEditStartDate] = useState('');
    const [editEndDate, setEditEndDate] = useState('');
    const [editingDates, setEditingDates] = useState(false);
    const [subscriptionTypes, setSubscriptionTypes] = useState<SubscriptionType[]>([]);
    const [selectedTypeId, setSelectedTypeId] = useState('');
    const [months, setMonths] = useState('1');
    const [payRef, setPayRef] = useState('');
    const [notes, setNotes] = useState('');
    const [startDate, setStartDate] = useState('');
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
    const [disablingId, setDisablingId] = useState<number | null>(null);

    const [subStatuses, setSubStatuses] = useState<Record<number, { status: string; startDate?: string; endDate?: string; typeName?: string }>>({});

    useEffect(() => {
        listSubscriptionTypesAction(true).then((res) => {
            if (res.success && res.data) setSubscriptionTypes(res.data);
        });
    }, []);

    useEffect(() => {
        if (!businesses.length) return;
        businesses.forEach(async (biz) => {
            const res = await getMySubscriptionAction(biz.id);
            if (res.success) {
                setSubStatuses((prev) => ({
                    ...prev,
                    [biz.id]: {
                        status: res.data?.status ?? 'pending',
                        startDate: res.data?.start_date,
                        endDate: res.data?.end_date,
                        typeName: res.data?.subscription_type_name,
                    },
                }));
            }
        });
    }, [businesses]);

    const openEditDates = (biz: { id: number; name: string }) => {
        const info = subStatuses[biz.id];
        setEditStartDate(toDateInputValue(info?.startDate));
        setEditEndDate(toDateInputValue(info?.endDate));
        setEditDatesModal({ open: true, business: biz });
    };

    const handleEditDates = async () => {
        if (!editDatesModal.business || !editStartDate || !editEndDate) return;
        setEditingDates(true);
        const res = await editSubscriptionDatesAction({
            businessId: editDatesModal.business.id,
            startDate: editStartDate,
            endDate: editEndDate,
        });
        setEditingDates(false);
        if (res.success) {
            setMessage({ type: 'success', text: `Fechas actualizadas para ${editDatesModal.business.name}.` });
            setSubStatuses((prev) => ({
                ...prev,
                [editDatesModal.business!.id]: {
                    ...prev[editDatesModal.business!.id],
                    status: prev[editDatesModal.business!.id]?.status ?? 'active',
                    startDate: `${editStartDate}T00:00:00Z`,
                    endDate: `${editEndDate}T00:00:00Z`,
                },
            }));
            setEditDatesModal({ open: false });
        } else {
            setMessage({ type: 'error', text: res.error || 'Error al actualizar las fechas' });
        }
    };

    const handleRegisterPayment = async () => {
        if (!registerModal.business || !selectedTypeId) return;
        setLoading(true);
        const res = await registerSubscriptionPaymentAction({
            businessId: registerModal.business.id,
            subscriptionTypeId: Number(selectedTypeId),
            monthsToAdd: Number(months),
            paymentReference: payRef || undefined,
            notes: notes || undefined,
            startDate: startDate || undefined,
        });
        setLoading(false);
        if (res.success) {
            setMessage({ type: 'success', text: `Pago registrado para ${registerModal.business.name}. Ahora puede usar la plataforma.` });
            setRegisterModal({ open: false });
            setSelectedTypeId(''); setMonths('1'); setPayRef(''); setNotes(''); setStartDate('');
            setSubStatuses((prev) => ({
                ...prev,
                [registerModal.business!.id]: { status: 'paid' },
            }));
        } else {
            setMessage({ type: 'error', text: res.error || 'Error al registrar pago' });
        }
    };

    const handleDisable = async (biz: { id: number; name: string }) => {
        if (!confirm(`¿Deseas suspender la cuenta de ${biz.name}?`)) return;
        setDisablingId(biz.id);
        const res = await disableSubscriptionAction(biz.id);
        setDisablingId(null);
        if (res.success) {
            setMessage({ type: 'success', text: `Cuenta de ${biz.name} suspendida.` });
            setSubStatuses((prev) => ({ ...prev, [biz.id]: { status: 'cancelled' } }));
        } else {
            setMessage({ type: 'error', text: res.error || 'Error al suspender' });
        }
    };

    const filteredBusinesses = businesses.filter((biz) => {
        const s = subStatuses[biz.id]?.status;
        if (filter === 'active' && !(s === 'active' || s === 'paid')) return false;
        if (filter === 'expired' && !(s === 'expired' || s === 'cancelled' || s === 'pending' || !s)) return false;
        if (search.trim() && !biz.name.toLowerCase().includes(search.trim().toLowerCase())) return false;
        return true;
    });

    return (
        <div className="space-y-6">
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Clientes y Suscripciones</h2>
                <div className="flex flex-col sm:flex-row gap-2">
                    <input
                        type="text"
                        value={search}
                        onChange={(e) => setSearch(e.target.value)}
                        placeholder="Buscar por nombre de negocio..."
                        className="px-3 py-1.5 border border-gray-200 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white dark:border-gray-600 min-w-[220px]"
                    />
                    <select
                        value={filter}
                        onChange={(e) => setFilter(e.target.value)}
                        className="px-3 py-1.5 border border-gray-200 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white"
                    >
                        <option value="all">Todos</option>
                        <option value="active">Activos</option>
                        <option value="expired">Vencidos / Suspendidos</option>
                    </select>
                </div>
            </div>

            {search.trim() && (
                <p className="text-xs text-gray-400 dark:text-gray-500">
                    {filteredBusinesses.length} resultado{filteredBusinesses.length !== 1 ? 's' : ''} para "{search.trim()}"
                </p>
            )}

            {message && (
                <Alert type={message.type} onClose={() => setMessage(null)}>
                    {message.text}
                </Alert>
            )}

            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                {filteredBusinesses.map((biz) => {
                    const subInfo = subStatuses[biz.id];
                    return (
                        <div key={biz.id} className="bg-white dark:bg-gray-700 rounded-xl border border-gray-200 dark:border-gray-600 shadow-sm p-5 space-y-3">
                            <div className="flex items-start justify-between">
                                <div>
                                    <h3 className="font-semibold text-gray-900 dark:text-white">{biz.name}</h3>
                                    <span className="text-xs text-gray-400">ID: {biz.id}</span>
                                    {subInfo?.typeName && (
                                        <p className="text-xs text-gray-400 mt-0.5">Plan: {subInfo.typeName}</p>
                                    )}
                                </div>
                                {subInfo
                                    ? <StatusBadge status={subInfo.status} />
                                    : <span className="text-xs text-gray-400 animate-pulse">Cargando...</span>
                                }
                            </div>

                            {subInfo?.startDate && subInfo?.endDate && (
                                <MembershipProgress startDate={subInfo.startDate} endDate={subInfo.endDate} />
                            )}

                            <div className="flex gap-2 pt-2 border-t border-gray-100 dark:border-gray-600">
                                {subInfo?.startDate && subInfo?.endDate && (
                                    <Button size="sm" variant="outline-purple" onClick={() => openEditDates(biz)} className="flex-1 text-xs">
                                        Editar
                                    </Button>
                                )}
                                <Button size="sm" variant="success" onClick={() => setRegisterModal({ open: true, business: biz })} className="flex-1 text-xs">
                                    Registrar Pago
                                </Button>
                                <Button size="sm" variant="danger" onClick={() => handleDisable(biz)} loading={disablingId === biz.id} className="flex-1 text-xs">
                                    Suspender
                                </Button>
                            </div>
                        </div>
                    );
                })}
            </div>

            <Modal isOpen={registerModal.open} onClose={() => setRegisterModal({ open: false })} title={`Registrar Pago — ${registerModal.business?.name}`} size="md">
                <div className="space-y-4 p-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Tipo de suscripción</label>
                        <select value={selectedTypeId} onChange={(e) => setSelectedTypeId(e.target.value)} className="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white">
                            <option value="">Selecciona un tipo</option>
                            {subscriptionTypes.map((t) => (
                                <option key={t.id} value={t.id}>{t.name} — {formatCurrency(t.price)}/{t.billing_period === 'monthly' ? 'mes' : 'año'}</option>
                            ))}
                        </select>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Meses a habilitar</label>
                        <select value={months} onChange={(e) => setMonths(e.target.value)} className="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white">
                            <option value="1">1 mes</option>
                            <option value="3">3 meses</option>
                            <option value="6">6 meses</option>
                            <option value="12">12 meses (anual)</option>
                        </select>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                            Fecha de inicio <span className="font-normal text-gray-400">(opcional — por defecto hoy, o el fin de la suscripción vigente)</span>
                        </label>
                        <input
                            type="date"
                            value={startDate}
                            onChange={(e) => setStartDate(e.target.value)}
                            className="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white dark:border-gray-600"
                        />
                    </div>
                    <Input label="Referencia de pago (opcional)" value={payRef} onChange={(e) => setPayRef(e.target.value)} placeholder="Nro. de transferencia, comprobante..." />
                    <Input label="Notas (opcional)" value={notes} onChange={(e) => setNotes(e.target.value)} placeholder="Observaciones internas..." />
                    <div className="flex justify-end gap-2 pt-2">
                        <Button variant="secondary" onClick={() => setRegisterModal({ open: false })}>Cancelar</Button>
                        <Button variant="success" onClick={handleRegisterPayment} loading={loading} disabled={!selectedTypeId}>Confirmar Pago</Button>
                    </div>
                </div>
            </Modal>

            <Modal isOpen={editDatesModal.open} onClose={() => setEditDatesModal({ open: false })} title={`Editar fechas — ${editDatesModal.business?.name}`} size="sm">
                <div className="space-y-4 p-4">
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                        Corrige la fecha de inicio o de vencimiento de la suscripción vigente, sin registrar un pago nuevo.
                    </p>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Fecha de inicio</label>
                        <input
                            type="date"
                            value={editStartDate}
                            onChange={(e) => setEditStartDate(e.target.value)}
                            className="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white dark:border-gray-600"
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Fecha de vencimiento</label>
                        <input
                            type="date"
                            value={editEndDate}
                            onChange={(e) => setEditEndDate(e.target.value)}
                            className="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white dark:border-gray-600"
                        />
                    </div>
                    <div className="flex justify-end gap-2 pt-2">
                        <Button variant="secondary" onClick={() => setEditDatesModal({ open: false })}>Cancelar</Button>
                        <Button variant="purple" onClick={handleEditDates} loading={editingDates} disabled={!editStartDate || !editEndDate}>Guardar cambios</Button>
                    </div>
                </div>
            </Modal>
        </div>
    );
}

function SubscriptionTypesAdminPanel() {
    const [types, setTypes] = useState<SubscriptionType[]>([]);
    const [moduleCatalog, setModuleCatalog] = useState<ModuleInfo[]>([]);
    const [loading, setLoading] = useState(true);
    const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
    const [editModal, setEditModal] = useState<{ open: boolean; type?: SubscriptionType }>({ open: false });

    const [form, setForm] = useState({ name: '', code: '', description: '', price: '', billing_period: 'monthly', active: true, module_codes: [] as string[], max_ecommerce_channels: '0' });

    const load = useCallback(async () => {
        setLoading(true);
        const [typesRes, catalogRes] = await Promise.all([listSubscriptionTypesAction(false), getModuleCatalogAction()]);
        if (typesRes.success && typesRes.data) setTypes(typesRes.data);
        if (catalogRes.success && catalogRes.data) setModuleCatalog(catalogRes.data);
        setLoading(false);
    }, []);

    useEffect(() => { load(); }, [load]);

    const openCreate = () => {
        setForm({ name: '', code: '', description: '', price: '', billing_period: 'monthly', active: true, module_codes: [], max_ecommerce_channels: '0' });
        setEditModal({ open: true });
    };

    const openEdit = (t: SubscriptionType) => {
        setForm({ name: t.name, code: t.code, description: t.description, price: String(t.price), billing_period: t.billing_period, active: t.active, module_codes: t.module_codes ?? [], max_ecommerce_channels: String(t.max_ecommerce_channels ?? 0) });
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
        if (!form.name || !form.price) return;
        const res = editModal.type
            ? await updateSubscriptionTypeAction(editModal.type.id, {
                name: form.name,
                description: form.description,
                price: Number(form.price),
                billing_period: form.billing_period,
                active: form.active,
                module_codes: form.module_codes,
                max_ecommerce_channels: Number(form.max_ecommerce_channels) || 0,
            })
            : await createSubscriptionTypeAction({
                name: form.name,
                code: form.code,
                description: form.description,
                price: Number(form.price),
                billing_period: form.billing_period,
                module_codes: form.module_codes,
                max_ecommerce_channels: Number(form.max_ecommerce_channels) || 0,
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

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-x-8 gap-y-4">
                        <div className="space-y-4">
                            <div className="flex items-baseline justify-between">
                                <span className="text-xs font-bold uppercase tracking-wide text-violet-600 dark:text-violet-400">Datos básicos</span>
                                <span className="text-xs text-gray-400">Cómo se cobra</span>
                            </div>

                            <Input label="Nombre" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
                            {!editModal.type && (
                                <Input label="Código (único)" value={form.code} onChange={(e) => setForm({ ...form, code: e.target.value })} placeholder="ej: pro-mensual" className="font-mono" />
                            )}
                            <Input label="Descripción" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} />

                            <div>
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">Precio</label>
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

    const [form, setForm] = useState({
        name: '', code: '', description: '', price: '', billing_period: 'monthly', active: true,
        module_codes: [] as string[], max_ecommerce_channels: '0', business_id: '', months: '1', notes: '',
    });

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
        setForm({ name: '', code: '', description: '', price: '', billing_period: 'monthly', active: true, module_codes: [], max_ecommerce_channels: '0', business_id: '', months: '1', notes: '' });
        setEditModal({ open: true });
    };

    const openEdit = (p: SubscriptionType) => {
        setForm({ name: p.name, code: p.code, description: p.description, price: String(p.price), billing_period: p.billing_period, active: p.active, module_codes: p.module_codes ?? [], max_ecommerce_channels: String(p.max_ecommerce_channels ?? 0), business_id: String(p.business_id ?? ''), months: '1', notes: '' });
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
        if (!form.name || !form.price) return;
        if (!editModal.plan && !form.business_id) {
            setMessage({ type: 'error', text: 'Selecciona el negocio al que se atará este plan' });
            return;
        }

        const res = editModal.plan
            ? await updateCustomPlanAction(editModal.plan.id, {
                name: form.name,
                description: form.description,
                price: Number(form.price),
                billing_period: form.billing_period,
                active: form.active,
                module_codes: form.module_codes,
                max_ecommerce_channels: Number(form.max_ecommerce_channels) || 0,
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
            });

        if (res.success) {
            setMessage({ type: 'success', text: editModal.plan ? 'Plan personalizado actualizado' : 'Plan personalizado creado y asociado al negocio' });
            setEditModal({ open: false });
            load();
        } else {
            setMessage({ type: 'error', text: res.error || 'Error al guardar' });
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

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-x-8 gap-y-4">
                        <div className="space-y-4">
                            <div className="flex items-baseline justify-between">
                                <span className="text-xs font-bold uppercase tracking-wide text-indigo-600 dark:text-indigo-400">Datos básicos</span>
                                <span className="text-xs text-gray-400">Cómo se cobra</span>
                            </div>

                            <div>
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">Negocio</label>
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

                            <Input label="Nombre" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
                            {!editModal.plan && (
                                <Input label="Código (único)" value={form.code} onChange={(e) => setForm({ ...form, code: e.target.value })} placeholder="ej: vip-negocio-x" className="font-mono" />
                            )}
                            <Input label="Descripción" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} />

                            <div>
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">Precio</label>
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

function BusinessSubscriptionView({ businessId, businessName, isSuperAdminView }: BusinessSubscriptionViewProps = {}) {
    const [subscription, setSubscription] = useState<BusinessSubscription | null>(null);
    const [loading, setLoading] = useState(true);

    const fetchSub = useCallback(async () => {
        setLoading(true);
        const res = await getMySubscriptionAction(businessId);
        if (res.success && res.data) setSubscription(res.data);
        setLoading(false);
    }, [businessId]);

    useEffect(() => { fetchSub(); }, [fetchSub]);

    if (loading) return <div className="flex justify-center py-12"><Spinner /></div>;

    const isExpired = !subscription || ['pending', 'expired', 'cancelled'].includes(subscription.status);

    return (
        <div className="space-y-6">
            {isSuperAdminView && (
                <div className="flex items-center gap-3 bg-blue-50 dark:bg-blue-900/30 border border-blue-200 dark:border-blue-700 rounded-lg px-4 py-3">
                    <p className="text-sm text-blue-800 dark:text-blue-300">
                        Vista de suscripción de <strong>{businessName}</strong> — modo super admin
                    </p>
                </div>
            )}

            <div className={`rounded-2xl p-6 text-white ${isExpired ? 'bg-gradient-to-br from-red-500 to-red-700' : 'bg-gradient-to-br from-violet-600 to-purple-800'} shadow-lg`}>
                <div className="flex items-center justify-between mb-4">
                    <div>
                        <p className="text-white/70 text-sm font-medium uppercase tracking-wider">Estado de Suscripción</p>
                        <h2 className="text-2xl font-bold mt-1">{isExpired ? 'Suspendida' : 'Activa'}</h2>
                        {subscription?.subscription_type_name && (
                            <p className="text-white/80 text-sm mt-1">Plan: {subscription.subscription_type_name}</p>
                        )}
                    </div>
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
            </div>

            <PlanCatalog businessId={businessId} onPurchased={fetchSub} currentSubscription={subscription} isCurrentActive={!isExpired} />
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
}

function PlanCatalog({ businessId, onPurchased, currentSubscription, isCurrentActive }: PlanCatalogProps) {
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
                    const total = purchaseModal.type!.price * Number(months);
                    const hasBalance = walletBalance !== null;
                    const sufficient = hasBalance && walletBalance! >= total;
                    const missing = hasBalance ? Math.max(0, total - walletBalance!) : 0;

                    return (
                        <div className="flex flex-col">
                            <div className="bg-gradient-to-br from-violet-600 to-purple-800 rounded-t-2xl -m-6 mb-0 px-6 pt-6 pb-5 text-white">
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

                            <div className="px-6 pt-5 pb-6 space-y-4">
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
                                        <span className="font-medium text-gray-700 dark:text-gray-200">{formatCurrency(total)}</span>
                                    </div>
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

                                <div className="flex justify-end gap-2 pt-1">
                                    <Button variant="secondary" onClick={() => setPurchaseModal({ open: false })}>Cancelar</Button>
                                    {!loadingBalance && !sufficient ? (
                                        <Button variant="purple" onClick={() => handleRechargeBold(missing)} loading={rechargingBold}>
                                            Recargar con Bold
                                        </Button>
                                    ) : (
                                        <Button variant="success" onClick={handleBuy} loading={buying} disabled={loadingBalance}>Confirmar Compra</Button>
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
}
