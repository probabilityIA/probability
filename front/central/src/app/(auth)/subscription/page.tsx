'use client';

import { useState, useEffect, useCallback } from 'react';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { Spinner, Button, Modal, Alert, Input } from '@/shared/ui';
import {
    getMySubscriptionAction,
    registerSubscriptionPaymentAction,
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
    BusinessSubscription,
    SubscriptionType,
    BusinessModuleOverride,
    ModuleInfo,
} from '@/services/modules/wallet/infra/subscription-actions';
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

export default function SubscriptionPage() {
    const { isSuperAdmin } = usePermissions();
    const { businesses } = useBusinessesSimple();
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
    const [adminTab, setAdminTab] = useState<'businesses' | 'types'>('businesses');

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
                                </div>
                                {adminTab === 'businesses'
                                    ? <AdminSubscriptionsView businesses={businesses} />
                                    : <SubscriptionTypesAdminPanel />}
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

function AdminSubscriptionsView({ businesses }: { businesses: Array<{ id: number; name: string }> }) {
    const [filter, setFilter] = useState('all');
    const [search, setSearch] = useState('');
    const [registerModal, setRegisterModal] = useState<{ open: boolean; business?: { id: number; name: string } }>({ open: false });
    const [subscriptionTypes, setSubscriptionTypes] = useState<SubscriptionType[]>([]);
    const [selectedTypeId, setSelectedTypeId] = useState('');
    const [months, setMonths] = useState('1');
    const [payRef, setPayRef] = useState('');
    const [notes, setNotes] = useState('');
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
    const [disablingId, setDisablingId] = useState<number | null>(null);

    const [subStatuses, setSubStatuses] = useState<Record<number, { status: string; endDate?: string; typeName?: string }>>({});

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
                        endDate: res.data?.end_date,
                        typeName: res.data?.subscription_type_name,
                    },
                }));
            }
        });
    }, [businesses]);

    const handleRegisterPayment = async () => {
        if (!registerModal.business || !selectedTypeId) return;
        setLoading(true);
        const res = await registerSubscriptionPaymentAction({
            businessId: registerModal.business.id,
            subscriptionTypeId: Number(selectedTypeId),
            monthsToAdd: Number(months),
            paymentReference: payRef || undefined,
            notes: notes || undefined,
        });
        setLoading(false);
        if (res.success) {
            setMessage({ type: 'success', text: `Pago registrado para ${registerModal.business.name}. Ahora puede usar la plataforma.` });
            setRegisterModal({ open: false });
            setSelectedTypeId(''); setMonths('1'); setPayRef(''); setNotes('');
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
                                    {subInfo?.endDate && (
                                        <p className="text-xs text-gray-400 mt-0.5">Vence: {formatDate(subInfo.endDate)}</p>
                                    )}
                                </div>
                                {subInfo
                                    ? <StatusBadge status={subInfo.status} />
                                    : <span className="text-xs text-gray-400 animate-pulse">Cargando...</span>
                                }
                            </div>
                            <div className="flex gap-2 pt-2 border-t border-gray-100 dark:border-gray-600">
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
                    <Input label="Referencia de pago (opcional)" value={payRef} onChange={(e) => setPayRef(e.target.value)} placeholder="Nro. de transferencia, comprobante..." />
                    <Input label="Notas (opcional)" value={notes} onChange={(e) => setNotes(e.target.value)} placeholder="Observaciones internas..." />
                    <div className="flex justify-end gap-2 pt-2">
                        <Button variant="secondary" onClick={() => setRegisterModal({ open: false })}>Cancelar</Button>
                        <Button variant="success" onClick={handleRegisterPayment} loading={loading} disabled={!selectedTypeId}>Confirmar Pago</Button>
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

            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                {types.map((t) => (
                    <div key={t.id} className="bg-white dark:bg-gray-700 rounded-xl border border-gray-200 dark:border-gray-600 shadow-sm p-5 space-y-3">
                        <div className="flex items-start justify-between">
                            <div>
                                <h3 className="font-semibold text-gray-900 dark:text-white">{t.name}</h3>
                                <span className="text-xs text-gray-400">{t.code}</span>
                            </div>
                            <span className={`text-xs px-2 py-1 rounded-full font-medium ${t.active ? 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-400' : 'bg-gray-100 text-gray-500 dark:bg-gray-800 dark:text-gray-400'}`}>
                                {t.active ? 'Activo' : 'Inactivo'}
                            </span>
                        </div>
                        <p className="text-lg font-bold text-gray-900 dark:text-white">{formatCurrency(t.price)}<span className="text-xs font-normal text-gray-400">/{t.billing_period === 'monthly' ? 'mes' : 'año'}</span></p>
                        {t.description && <p className="text-sm text-gray-500 dark:text-gray-400">{t.description}</p>}
                        <div className="flex flex-wrap gap-1">
                            {(t.module_codes ?? []).map((m) => (
                                <span key={m} className="text-xs bg-violet-50 dark:bg-violet-900/30 text-violet-700 dark:text-violet-300 px-2 py-0.5 rounded-full">{moduleName(m)}</span>
                            ))}
                        </div>
                        <div className="flex gap-2 pt-2 border-t border-gray-100 dark:border-gray-600">
                            <Button size="sm" variant="secondary" onClick={() => openEdit(t)} className="flex-1 text-xs">Editar</Button>
                            <Button size="sm" variant="danger" onClick={() => handleDelete(t)} className="flex-1 text-xs">Eliminar</Button>
                        </div>
                    </div>
                ))}
            </div>

            <Modal isOpen={editModal.open} onClose={() => setEditModal({ open: false })} title={editModal.type ? `Editar ${editModal.type.name}` : 'Nuevo tipo de suscripción'} size="md">
                <div className="space-y-4 p-4">
                    <Input label="Nombre" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
                    {!editModal.type && (
                        <Input label="Código (unico)" value={form.code} onChange={(e) => setForm({ ...form, code: e.target.value })} placeholder="ej: basico" />
                    )}
                    <Input label="Descripción" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} />
                    <Input label="Precio" type="number" value={form.price} onChange={(e) => setForm({ ...form, price: e.target.value })} />
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Periodo de facturación</label>
                        <select value={form.billing_period} onChange={(e) => setForm({ ...form, billing_period: e.target.value })} className="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white">
                            <option value="monthly">Mensual</option>
                            <option value="annual">Anual</option>
                        </select>
                    </div>
                    {editModal.type && (
                        <label className="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
                            <input type="checkbox" checked={form.active} onChange={(e) => setForm({ ...form, active: e.target.checked })} />
                            Activo
                        </label>
                    )}
                    <Input label="Límite de canales E-commerce (0 = sin límite)" type="number" value={form.max_ecommerce_channels} onChange={(e) => setForm({ ...form, max_ecommerce_channels: e.target.value })} />
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Módulos incluidos</label>
                        <div className="grid grid-cols-2 gap-2">
                            {moduleCatalog.map(({ code, name }) => (
                                <label key={code} className="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
                                    <input type="checkbox" checked={form.module_codes.includes(code)} onChange={() => toggleModule(code)} />
                                    {name}
                                </label>
                            ))}
                        </div>
                    </div>
                    <div className="flex justify-end gap-2 pt-2">
                        <Button variant="secondary" onClick={() => setEditModal({ open: false })}>Cancelar</Button>
                        <Button variant="primary" onClick={handleSave}>Guardar</Button>
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

    const moduleName = (code: string) => moduleCatalog.find((m) => m.code === code)?.name ?? code;

    useEffect(() => {
        Promise.all([listSubscriptionTypesAction(true), getModuleCatalogAction()]).then(([typesRes, catalogRes]) => {
            if (typesRes.success && typesRes.data) setTypes(typesRes.data);
            if (catalogRes.success && catalogRes.data) setModuleCatalog(catalogRes.data);
            setLoading(false);
        });
    }, []);

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
                                        onClick={() => { setMonths('1'); setPurchaseModal({ open: true, type: t }); }}
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

            <Modal isOpen={purchaseModal.open} onClose={() => setPurchaseModal({ open: false })} title={`Comprar ${purchaseModal.type?.name}`} size="sm">
                <div className="space-y-4 p-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Meses</label>
                        <select value={months} onChange={(e) => setMonths(e.target.value)} className="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white">
                            <option value="1">1 mes</option>
                            <option value="3">3 meses</option>
                            <option value="6">6 meses</option>
                            <option value="12">12 meses (anual)</option>
                        </select>
                    </div>
                    {purchaseModal.type && (
                        <p className="text-sm text-gray-600 dark:text-gray-300">
                            Total a debitar de tu billetera: <strong>{formatCurrency(purchaseModal.type.price * Number(months))}</strong>
                        </p>
                    )}
                    <div className="flex justify-end gap-2 pt-2">
                        <Button variant="secondary" onClick={() => setPurchaseModal({ open: false })}>Cancelar</Button>
                        <Button variant="success" onClick={handleBuy} loading={buying}>Confirmar Compra</Button>
                    </div>
                </div>
            </Modal>
        </div>
    );
}
