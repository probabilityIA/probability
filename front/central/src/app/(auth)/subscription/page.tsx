'use client';

import { useState, useEffect, useCallback } from 'react';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { Spinner, Button, Modal, Alert, Input } from '@/shared/ui';
import {
    getMySubscriptionAction,
    registerSubscriptionPaymentAction,
    disableSubscriptionAction,
    BusinessSubscription,
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

// Badge de estado de suscripción con color dinámico
function StatusBadge({ status }: { status?: string }) {
    const map: Record<string, { label: string; cls: string }> = {
        active: { label: 'Activo', cls: 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-400' },
        paid: { label: 'Activo', cls: 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-400' },
        expired: { label: 'Vencido', cls: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400' },
        cancelled: { label: 'Suspendido', cls: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400' },
        pending: { label: 'Pendiente', cls: 'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-400' },
    };
    const entry = map[status ?? ''] ?? { label: 'Sin suscripción', cls: 'bg-gray-100 text-gray-500 dark:text-gray-400 dark:bg-gray-700 dark:text-gray-400' };
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
                            <BusinessSubscriptionView
                                businessId={selectedBusinessId}
                                businessName={selectedBusiness?.name}
                                isSuperAdminView
                            />
                        ) : (
                            <AdminSubscriptionsView businesses={businesses} />
                        )
                    ) : (
                        <BusinessSubscriptionView />
                    )}
                </div>
            </div>
        </div>
    );
}

// ─────────────────────────────────────────────────────────────────────────────
// Admin View: Lista todos los negocios con estado REAL de suscripción
// ─────────────────────────────────────────────────────────────────────────────
function AdminSubscriptionsView({ businesses }: { businesses: Array<{ id: number; name: string }> }) {
    const [filter, setFilter] = useState('all');
    const [registerModal, setRegisterModal] = useState<{ open: boolean; business?: { id: number; name: string } }>({ open: false });
    const [amount, setAmount] = useState('');
    const [months, setMonths] = useState('1');
    const [payRef, setPayRef] = useState('');
    const [notes, setNotes] = useState('');
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
    const [disablingId, setDisablingId] = useState<number | null>(null);

    // Estado real de suscripción por negocio
    const [subStatuses, setSubStatuses] = useState<Record<number, { status: string; endDate?: string }>>({});

    useEffect(() => {
        if (!businesses.length) return;
        businesses.forEach(async (biz) => {
            const res = await getMySubscriptionAction(biz.id);
            if (res.success) {
                setSubStatuses((prev) => ({
                    ...prev,
                    [biz.id]: {
                        status: res.data?.status ?? 'pending',
                        endDate: res.data?.endDate,
                    },
                }));
            }
        });
    }, [businesses]);

    const handleRegisterPayment = async () => {
        if (!registerModal.business) return;
        if (!amount || isNaN(Number(amount))) return;
        setLoading(true);
        const res = await registerSubscriptionPaymentAction({
            businessId: registerModal.business.id,
            amount: Number(amount),
            monthsToAdd: Number(months),
            paymentReference: payRef || undefined,
            notes: notes || undefined,
        });
        setLoading(false);
        if (res.success) {
            setMessage({ type: 'success', text: `Pago registrado para ${registerModal.business.name}. Ahora puede usar la plataforma.` });
            setRegisterModal({ open: false });
            setAmount(''); setMonths('1'); setPayRef(''); setNotes('');
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
            setSubStatuses((prev) => ({ ...prev, [biz.id]: { status: 'expired' } }));
        } else {
            setMessage({ type: 'error', text: res.error || 'Error al suspender' });
        }
    };

    const filteredBusinesses = businesses.filter((biz) => {
        const s = subStatuses[biz.id]?.status;
        if (filter === 'active') return s === 'active' || s === 'paid';
        if (filter === 'expired') return s === 'expired' || s === 'cancelled' || s === 'pending' || !s;
        return true;
    });

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Clientes y Suscripciones</h2>
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

            {message && (
                <Alert type={message.type} onClose={() => setMessage(null)}>
                    {message.text}
                </Alert>
            )}

            <div className="bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-700 rounded-lg p-4">
                <div className="flex gap-3">
                    <span className="text-2xl">💡</span>
                    <div className="text-sm text-amber-800 dark:text-amber-300">
                        <p className="font-semibold mb-1">Panel de control de suscripciones</p>
                        <p>Registra el pago mensual de cada cliente para habilitarles el acceso. Si un cliente no paga, suspéndelo.</p>
                    </div>
                </div>
            </div>

            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                {filteredBusinesses.map((biz) => {
                    const subInfo = subStatuses[biz.id];
                    return (
                        <div key={biz.id} className="bg-white dark:bg-gray-700 rounded-xl border border-gray-200 dark:border-gray-600 shadow-sm p-5 space-y-3">
                            <div className="flex items-start justify-between">
                                <div>
                                    <h3 className="font-semibold text-gray-900 dark:text-white">{biz.name}</h3>
                                    <span className="text-xs text-gray-400">ID: {biz.id}</span>
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
                                    💳 Registrar Pago
                                </Button>
                                <Button size="sm" variant="danger" onClick={() => handleDisable(biz)} loading={disablingId === biz.id} className="flex-1 text-xs">
                                    🔒 Suspender
                                </Button>
                            </div>
                        </div>
                    );
                })}
            </div>

            <Modal isOpen={registerModal.open} onClose={() => setRegisterModal({ open: false })} title={`Registrar Pago — ${registerModal.business?.name}`} size="md">
                <div className="space-y-4 p-4">
                    <div className="bg-blue-50 dark:bg-blue-900/30 rounded-lg p-3 text-sm text-blue-800 dark:text-blue-300">
                        Al registrar el pago, el cliente quedará activo y podrá usar todas las funciones.
                    </div>
                    <Input label="Monto pagado (COP)" type="number" value={amount} onChange={(e) => setAmount(e.target.value)} placeholder="Ej: 150000" />
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 dark:text-gray-300 mb-1">Meses a habilitar</label>
                        <select value={months} onChange={(e) => setMonths(e.target.value)} className="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white dark:bg-gray-700 dark:text-white">
                            <option value="1">1 mes</option>
                            <option value="2">2 meses</option>
                            <option value="3">3 meses</option>
                            <option value="6">6 meses</option>
                            <option value="12">12 meses (anual)</option>
                        </select>
                    </div>
                    <Input label="Referencia de pago (opcional)" value={payRef} onChange={(e) => setPayRef(e.target.value)} placeholder="Nro. de transferencia, comprobante..." />
                    <Input label="Notas (opcional)" value={notes} onChange={(e) => setNotes(e.target.value)} placeholder="Observaciones internas..." />
                    <div className="flex justify-end gap-2 pt-2">
                        <Button variant="secondary" onClick={() => setRegisterModal({ open: false })}>Cancelar</Button>
                        <Button variant="success" onClick={handleRegisterPayment} loading={loading}>✅ Confirmar Pago</Button>
                    </div>
                </div>
            </Modal>
        </div>
    );
}

// ─────────────────────────────────────────────────────────────────────────────
// Business View: El negocio ve su estado y el QR/info para pagar
// ─────────────────────────────────────────────────────────────────────────────
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

    const isExpired = subscription?.status === 'pending' || !subscription;

    return (
        <div className="space-y-6">
            {isSuperAdminView && (
                <div className="flex items-center gap-3 bg-blue-50 dark:bg-blue-900/30 border border-blue-200 dark:border-blue-700 rounded-lg px-4 py-3">
                    <p className="text-sm text-blue-800 dark:text-blue-300">
                        Vista de suscripción de <strong>{businessName}</strong> — modo super admin
                    </p>
                </div>
            )}

            {/* Status Card */}
            <div className={`rounded-2xl p-6 text-white ${isExpired ? 'bg-gradient-to-br from-red-500 to-red-700' : 'bg-gradient-to-br from-violet-600 to-purple-800'} shadow-lg`}>
                <div className="flex items-center justify-between mb-4">
                    <div>
                        <p className="text-white/70 text-sm font-medium uppercase tracking-wider">Estado de Suscripción</p>
                        <h2 className="text-2xl font-bold mt-1">{isExpired ? '⚠️ Suspendida' : '✅ Activa'}</h2>
                    </div>
                    <div className="text-5xl opacity-30">{isExpired ? '🔒' : '🚀'}</div>
                </div>
                {subscription && (
                    <div className="grid grid-cols-2 gap-4 mt-4 pt-4 border-t border-white/20">
                        <div>
                            <p className="text-white/60 text-xs">Último pago</p>
                            <p className="font-semibold">{formatCurrency(subscription.amount)}</p>
                        </div>
                        <div>
                            <p className="text-white/60 text-xs">Válida hasta</p>
                            <p className="font-semibold">{formatDate(subscription.endDate)}</p>
                        </div>
                    </div>
                )}
            </div>

            {/* Payment instructions + QR */}
            <div className="bg-white dark:bg-gray-700 rounded-xl border border-gray-200 dark:border-gray-600 shadow-sm p-6">
                <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-4">📋 Información de Pago</h3>
                <p className="text-gray-600 dark:text-gray-300 dark:text-gray-300 text-sm mb-4">
                    Para renovar tu suscripción, escanea el QR con tu app Nequi o realiza una transferencia
                    con los datos indicados y notifica a tu asesor con el comprobante.
                </p>

                <div className="flex flex-col sm:flex-row gap-6 items-start">
                    {/* QR de pago */}
                    <div className="flex flex-col items-center gap-2 flex-shrink-0">
                        <div className="w-48 h-48 rounded-xl border-2 border-violet-200 dark:border-violet-700 overflow-hidden bg-white p-1 shadow-sm">
                            {/* eslint-disable-next-line @next/next/no-img-element */}
                            <img src="/QR.png" alt="QR de pago Nequi" className="w-full h-full object-contain" />
                        </div>
                        <p className="text-xs text-gray-500 dark:text-gray-400 text-center">Escanear con Nequi</p>
                    </div>

                    {/* Datos bancarios */}
                    <div className="flex-1 bg-gray-50 dark:bg-gray-800 rounded-xl p-4 space-y-3 border border-gray-100 dark:border-gray-600">
                        <div className="flex items-center gap-3">
                            <div className="w-10 h-10 bg-green-100 dark:bg-green-900/40 rounded-full flex items-center justify-center text-xl">🏦</div>
                            <div>
                                <p className="font-semibold text-gray-900 dark:text-white text-sm">Bancolombia / Nequi</p>
                                <p className="text-xs text-gray-500 dark:text-gray-400">Cuenta de Ahorros</p>
                            </div>
                        </div>
                        <div className="grid grid-cols-2 gap-2 text-sm">
                            <div className="bg-white dark:bg-gray-700 rounded-lg p-3">
                                <p className="text-gray-400 text-xs uppercase tracking-wider">Titular</p>
                                <p className="font-semibold text-gray-900 dark:text-white mt-0.5">ProbabilityIA SAS</p>
                            </div>
                            <div className="bg-white dark:bg-gray-700 rounded-lg p-3">
                                <p className="text-gray-400 text-xs uppercase tracking-wider">Número de cuenta</p>
                                <p className="font-semibold text-gray-900 dark:text-white font-mono mt-0.5">*** *** **** ****</p>
                            </div>
                            <div className="bg-white dark:bg-gray-700 rounded-lg p-3">
                                <p className="text-gray-400 text-xs uppercase tracking-wider">NIT</p>
                                <p className="font-semibold text-gray-900 dark:text-white mt-0.5">***.***.***-*</p>
                            </div>
                            <div className="bg-white dark:bg-gray-700 rounded-lg p-3">
                                <p className="text-gray-400 text-xs uppercase tracking-wider">Celular (Nequi)</p>
                                <p className="font-semibold text-gray-900 dark:text-white mt-0.5">+57 *** *** ****</p>
                            </div>
                        </div>
                    </div>
                </div>

                <div className="mt-4 flex items-start gap-3 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-700 rounded-lg p-4">
                    <span className="text-lg">⚠️</span>
                    <p className="text-sm text-amber-800 dark:text-amber-300">
                        Una vez realices el pago, envía el comprobante a tu asesor. Activaremos tu cuenta en un máximo de <strong>2 horas hábiles</strong>.
                    </p>
                </div>
            </div>
        </div>
    );
}
