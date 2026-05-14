'use client';

import { useState, useEffect, useCallback } from 'react';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { Spinner, Button, Input, Table, TableColumn, Alert, Modal } from '@/shared/ui';
import { PaymentMethodSelectorModal } from '@/services/modules/pay/ui';
import {
    getWalletsAction,
    getPendingRequestsAction,
    getProcessedRequestsAction,
    processRequestAction,
    getWalletBalanceAction,
    rechargeWalletAction,
    manualDebitAction,
    getWalletHistoryAction,
    clearRechargeHistoryAction,
    adminAdjustBalanceAction,
    Wallet
} from '@/services/modules/wallet/infra/actions';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';
import { VirtualCard } from './virtual-card';
import { FinancialStatsView } from './financial-stats';
import { BoldPaymentProcessingModal } from './bold-payment-processing-modal';
import { getActionError } from '@/shared/utils/action-result';
import { getBoldSignatureAction } from '@/services/modules/pay/infra/actions';

const formatCurrency = (amount: number) =>
    new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP' }).format(amount);

const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return `${date.toLocaleDateString()} ${date.toLocaleTimeString()}`;
};

export default function WalletPage() {
    const { isSuperAdmin } = usePermissions();
    const { businesses } = useBusinessesSimple();
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);

    const selectedBusiness = businesses.find(b => b.id === selectedBusinessId);

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
                <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white">Billetera</h1>

                {/* Business selector - solo para super admin */}
                {isSuperAdmin && businesses.length > 0 && (
                    <div className="flex items-center gap-3 rounded-lg px-4 py-2" style={{ backgroundColor: 'var(--color-primary-50)', border: `1px solid var(--color-primary-200)` }}>
                        <label className="text-sm font-medium whitespace-nowrap" style={{ color: 'var(--color-primary-900)' }}>
                            Negocio:
                        </label>
                        <select
                            value={selectedBusinessId?.toString() ?? ''}
                            onChange={(e) => {
                                const val = e.target.value;
                                setSelectedBusinessId(val ? Number(val) : null);
                            }}
                            className="px-3 py-1.5 rounded-md text-sm bg-white dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-2 min-w-[200px]"
                            style={{ border: `1px solid var(--color-primary-300)` }}
                        >
                            <option value="">Vista Administrativa</option>
                            {businesses.map((b) => (
                                <option key={b.id} value={b.id}>{b.name}</option>
                            ))}
                        </select>
                    </div>
                )}
            </div>

            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm dark:shadow-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
                <div className="p-6">
                    {isSuperAdmin ? (
                        selectedBusinessId ? (
                            <BusinessWalletView
                                businessId={selectedBusinessId}
                                businessName={selectedBusiness?.name}
                            />
                        ) : (
                            <AdminWalletView />
                        )
                    ) : (
                        <BusinessWalletView />
                    )}
                </div>
            </div>
        </div>
    );
}

function AdminWalletView() {
    const [wallets, setWallets] = useState<Wallet[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [businesses, setBusinesses] = useState<Record<number, string>>({});
    const [itemsPerPage, setItemsPerPage] = useState(10);
    const [activeTab, setActiveTab] = useState<'saldos' | 'finanzas'>('saldos');

    const fetchWalletsAndBusinesses = useCallback(async () => {
        try {
            setLoading(true);
            const walletRes = await getWalletsAction();
            if (!walletRes.success) throw new Error(walletRes.error || 'Failed to fetch wallets');
            setWallets(walletRes.data || []);

            const { getBusinessesAction } = await import('@/services/auth/business/infra/actions');
            const businessesRes = await getBusinessesAction({ per_page: 1000 });
            if (businessesRes.data) {
                const businessMap: Record<number, string> = {};
                businessesRes.data.forEach((b: any) => {
                    businessMap[b.id] = b.name;
                });
                setBusinesses(businessMap);
            }
        } catch (err: any) {
            setError(getActionError(err));
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        fetchWalletsAndBusinesses();
    }, [fetchWalletsAndBusinesses]);

    const walletColumns: TableColumn<Wallet>[] = [
        {
            key: 'BusinessID',
            label: 'Negocio',
            render: (_val, row) => (
                <span className="font-medium text-gray-900 dark:text-white">
                    {businesses[row.BusinessID] || `ID: ${row.BusinessID}`}
                </span>
            )
        },
        {
            key: 'Balance',
            label: 'Saldo',
            render: (val) => {
                const balance = typeof val === 'string' ? parseFloat(val) : (val as number);
                const isNegative = balance < 0;
                return (
                    <span className="font-bold" style={{ color: isNegative ? '#dc2626' : '#16a34a' }}>
                        {isNegative && '-'}${formatCurrency(Math.abs(balance)).replace('$', '')}
                    </span>
                );
            }
        }
    ];

    if (error) return <Alert type="error">{error}</Alert>;

    return (
        <div className="space-y-6">
            {/* Tabs */}
            <div className="flex gap-2 border-b border-gray-200 dark:border-gray-700">
                <button
                    onClick={() => setActiveTab('saldos')}
                    className={`px-4 py-3 text-sm font-medium border-b-2 transition-colors ${
                        activeTab === 'saldos'
                            ? 'border-b-2'
                            : 'text-gray-600 dark:text-gray-400 border-transparent hover:text-gray-900 dark:hover:text-gray-200'
                    }`}
                    style={activeTab === 'saldos' ? { color: 'var(--color-primary-600)', borderColor: 'var(--color-primary-600)' } : {}}
                >
                    Saldos
                </button>
                <button
                    onClick={() => setActiveTab('finanzas')}
                    className={`px-4 py-3 text-sm font-medium border-b-2 transition-colors ${
                        activeTab === 'finanzas'
                            ? 'border-b-2'
                            : 'text-gray-600 dark:text-gray-400 border-transparent hover:text-gray-900 dark:hover:text-gray-200'
                    }`}
                    style={activeTab === 'finanzas' ? { color: 'var(--color-primary-600)', borderColor: 'var(--color-primary-600)' } : {}}
                >
                    📊 Finanzas
                </button>
            </div>

            {/* Tab Content */}
            {activeTab === 'saldos' && (
                <div className="space-y-8">
                    {/* Wallets Section */}
                    <div>
                        <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Saldos de Negocios</h2>
                        <Table
                            columns={[
                                ...walletColumns,
                                {
                                    key: 'actions',
                                    label: 'Acciones',
                                    render: (_, row) => (
                                        <div className="flex gap-2">
                                            <RechargeWalletButton
                                                businessId={row.BusinessID}
                                                businessName={businesses[row.BusinessID] || `ID: ${row.BusinessID}`}
                                                onSuccess={fetchWalletsAndBusinesses}
                                            />
                                            <ManualDebitButton
                                                businessId={row.BusinessID}
                                                businessName={businesses[row.BusinessID] || `ID: ${row.BusinessID}`}
                                                onSuccess={fetchWalletsAndBusinesses}
                                            />
                                            <ClearHistoryButton
                                                businessId={row.BusinessID}
                                                businessName={businesses[row.BusinessID] || `ID: ${row.BusinessID}`}
                                                onSuccess={fetchWalletsAndBusinesses}
                                            />
                                        </div>
                                    )
                                }
                            ]}
                            data={wallets}
                            loading={loading}
                            emptyMessage="No hay billeteras registradas"
                        />
                    </div>

                    {/* Top Section - En Revisión (Full Width) */}
                    <RequestsTableView
                        title="En revisión"
                        businesses={businesses}
                        onRequestsChanged={fetchWalletsAndBusinesses}
                        allWallets={wallets}
                        fetchAction={getPendingRequestsAction}
                        showActions={true}
                        emptyMessage="Sin pendientes"
                        compact={false}
                        itemsPerPage={itemsPerPage}
                        onItemsPerPageChange={setItemsPerPage}
                    />

                    {/* Bottom Row - Approved and Rejected (Side by Side) */}
                    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                        <RequestsTableView
                            title="Aprobados"
                            businesses={businesses}
                            onRequestsChanged={fetchWalletsAndBusinesses}
                            allWallets={wallets}
                            fetchAction={getProcessedRequestsAction}
                            filterStatus="COMPLETED"
                            showActions={false}
                            emptyMessage="Sin aprobados"
                            compact={true}
                            itemsPerPage={itemsPerPage}
                            onItemsPerPageChange={setItemsPerPage}
                        />

                        <RequestsTableView
                            title="Rechazados"
                            businesses={businesses}
                            onRequestsChanged={fetchWalletsAndBusinesses}
                            allWallets={wallets}
                            fetchAction={getProcessedRequestsAction}
                            filterStatus="FAILED"
                            showActions={false}
                            emptyMessage="Sin rechazados"
                            compact={true}
                            itemsPerPage={itemsPerPage}
                            onItemsPerPageChange={setItemsPerPage}
                        />
                    </div>
                </div>
            )}

            {activeTab === 'finanzas' && <FinancialStatsView />}
        </div>
    );
}

function ManualDebitButton({ businessId, businessName, onSuccess }: { businessId: number, businessName: string, onSuccess: () => void }) {
    const [isOpen, setIsOpen] = useState(false);
    const [amount, setAmount] = useState('');
    const [reference, setReference] = useState('');
    const [loading, setLoading] = useState(false);

    const handleDebit = async () => {
        if (!amount || isNaN(Number(amount))) return;
        setLoading(true);
        try {
            const res = await manualDebitAction(businessId, Number(amount), reference);
            if (res.success) {
                setIsOpen(false);
                setAmount('');
                setReference('');
                onSuccess();
            } else {
                alert(res.error);
            }
        } catch (e) {
            alert("Error al procesar");
        } finally {
            setLoading(false);
        }
    };

    return (
        <>
            <Button size="sm" variant="danger" onClick={() => setIsOpen(true)}>Restar Saldo</Button>
            <Modal isOpen={isOpen} onClose={() => setIsOpen(false)} title={`Restar saldo a ${businessName}`}>
                <div className="space-y-4 p-4">
                    <Input
                        label="Monto a restar"
                        type="number"
                        value={amount}
                        onChange={e => setAmount(e.target.value)}
                        placeholder="Ej: 5000"
                    />
                    <Input
                        label="Referencia / Motivo"
                        value={reference}
                        onChange={e => setReference(e.target.value)}
                        placeholder="Ej: Ajuste de saldo"
                    />
                    <div className="flex justify-end gap-2">
                        <Button variant="secondary" onClick={() => setIsOpen(false)}>Cancelar</Button>
                        <Button variant="danger" onClick={handleDebit} loading={loading}>Restar Saldo</Button>
                    </div>
                </div>
            </Modal>
        </>
    );
}

function ClearHistoryButton({ businessId, businessName, onSuccess }: { businessId: number, businessName: string, onSuccess: () => void }) {
    const [isOpen, setIsOpen] = useState(false);
    const [loading, setLoading] = useState(false);

    const handleClear = async () => {
        setLoading(true);
        try {
            console.log('Borrando historial para business:', businessId);
            const res = await clearRechargeHistoryAction(businessId);
            console.log('Respuesta del servidor:', res);
            if (res.success) {
                setIsOpen(false);
                alert('Historial borrado exitosamente');
                onSuccess();
            } else {
                console.error('Error del servidor:', res.error);
                alert(`Error: ${res.error || 'Error desconocido'}`);
            }
        } catch (e: any) {
            console.error('Error en handleClear:', e);
            alert(`Error al procesar: ${e.message || e}`);
        } finally {
            setLoading(false);
        }
    };

    return (
        <>
            <Button size="sm" variant="outline" onClick={() => setIsOpen(true)} style={{ color: '#dc2626', borderColor: '#fecaca' }}>
                Borrar Historial
            </Button>
            <Modal isOpen={isOpen} onClose={() => setIsOpen(false)} title="Confirmar Eliminación">
                <div className="space-y-4 p-4">
                    <p className="text-gray-600 dark:text-gray-300">
                        ¿Estás seguro de que deseas borrar <strong>todo el historial de recargas</strong> de <strong>{businessName}</strong>?
                    </p>
                    <Alert type="warning">
                        Esta acción es irreversible y eliminará todos los registros de recargas (aprobadas, rechazadas y pendientes). El saldo actual no se verá afectado.
                    </Alert>
                    <div className="flex justify-end gap-2">
                        <Button variant="secondary" onClick={() => setIsOpen(false)}>Cancelar</Button>
                        <Button variant="danger" onClick={handleClear} loading={loading}>Borrar Historial</Button>
                    </div>
                </div>
            </Modal>
        </>
    );
}

function RechargeWalletButton({ businessId, businessName, onSuccess }: { businessId: number, businessName: string, onSuccess: () => void }) {
    const [isOpen, setIsOpen] = useState(false);
    const [amount, setAmount] = useState('');
    const [reason, setReason] = useState('');
    const [loading, setLoading] = useState(false);

    const handleRecharge = async () => {
        if (!amount || isNaN(Number(amount)) || Number(amount) <= 0) {
            alert('Ingresa un monto válido');
            return;
        }
        if (!reason.trim()) {
            alert('Ingresa un motivo para la recarga');
            return;
        }
        setLoading(true);
        try {
            const res = await adminAdjustBalanceAction(businessId, Number(amount), reason);
            if (res.success) {
                setIsOpen(false);
                setAmount('');
                setReason('');
                onSuccess();
                alert('Saldo agregado exitosamente');
            } else {
                alert(res.error || 'Error al recargar');
            }
        } catch (e: any) {
            alert(`Error al procesar: ${e.message || 'Error desconocido'}`);
        } finally {
            setLoading(false);
        }
    };

    return (
        <>
            <Button size="sm" variant="primary" onClick={() => setIsOpen(true)}>Agregar Saldo</Button>
            <Modal isOpen={isOpen} onClose={() => setIsOpen(false)} title={`Agregar saldo a ${businessName}`}>
                <div className="space-y-4 p-4">
                    <Input
                        label="Monto a agregar"
                        type="number"
                        value={amount}
                        onChange={e => setAmount(e.target.value)}
                        placeholder="Ej: 10000"
                    />
                    <Input
                        label="Motivo de la recarga"
                        type="text"
                        value={reason}
                        onChange={e => setReason(e.target.value)}
                        placeholder="Ej: Ajuste por error, promoción, etc."
                    />
                    <div className="flex justify-end gap-2">
                        <Button variant="secondary" onClick={() => setIsOpen(false)}>Cancelar</Button>
                        <Button variant="primary" onClick={handleRecharge} loading={loading}>Agregar Saldo</Button>
                    </div>
                </div>
            </Modal>
        </>
    );
}

function RequestsTableView({
    title,
    businesses,
    onRequestsChanged,
    allWallets,
    fetchAction,
    showActions,
    emptyMessage,
    filterStatus,
    compact,
    itemsPerPage,
    onItemsPerPageChange
}: {
    title: string,
    businesses: Record<number, string>,
    onRequestsChanged: () => void,
    allWallets: Wallet[],
    fetchAction: () => Promise<any>,
    showActions: boolean,
    emptyMessage: string,
    filterStatus?: string,
    compact?: boolean,
    itemsPerPage: number,
    onItemsPerPageChange: (total: number) => void
}) {
    const [requests, setRequests] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);
    const [processingId, setProcessingId] = useState<string | null>(null);
    const [currentPage, setCurrentPage] = useState(1);

    const fetchRequests = useCallback(async () => {
        try {
            setLoading(true);
            const res = await fetchAction();
            if (res.success) {
                let data = res.data as any[] || [];
                if (filterStatus) {
                    data = data.filter(r => r.Status === filterStatus);
                }
                setRequests(data);
            }
        } catch (error) {
            console.error(error);
        } finally {
            setLoading(false);
        }
    }, [fetchAction, filterStatus]);

    useEffect(() => {
        fetchRequests();
    }, [fetchRequests]);

    const handleAction = async (id: string, action: 'approve' | 'reject') => {
        setProcessingId(id);
        try {
            const res = await processRequestAction(id, action);
            if (res.success) {
                await fetchRequests();
                onRequestsChanged();
            } else {
                alert(`Error al procesar la solicitud: ${res.error}`);
            }
        } catch (e) {
            alert("Error de conexión.");
        } finally {
            setProcessingId(null);
        }
    };

    const requestColumns: TableColumn<any>[] = [
        {
            key: 'CreatedAt',
            label: 'Fecha',
            render: (val) => <span className="text-gray-600 dark:text-gray-300 font-mono text-sm">{formatDate(val as string)}</span>
        },
        {
            key: 'WalletID',
            label: 'Negocio',
            render: (val) => {
                const wallet = allWallets.find(w => w.ID === val);
                if (wallet) {
                    const name = businesses[wallet.BusinessID];
                    return <span className="font-medium text-gray-900 dark:text-white">{name || `ID: ${wallet.BusinessID}`}</span>;
                }
                return <span className="text-gray-500 dark:text-gray-400">...</span>;
            }
        },
        {
            key: 'Amount',
            label: 'Monto',
            render: (val) => <span className="font-bold text-gray-900 dark:text-white">{formatCurrency(val as number)}</span>
        },
    ];

    if (showActions) {
        requestColumns.push({
            key: 'actions',
            label: 'Acciones',
            render: (_, row) => (
                <div className="flex gap-1.5">
                    <Button
                        size="sm"
                        className="px-2 py-0.5 text-[10px] h-auto min-h-0"
                        variant="success"
                        onClick={() => handleAction(row.ID, 'approve')}
                        loading={processingId === row.ID}
                        disabled={!!processingId}
                    >
                        Aprobar
                    </Button>
                    <Button
                        size="sm"
                        className="px-2 py-0.5 text-[10px] h-auto min-h-0"
                        variant="danger"
                        onClick={() => handleAction(row.ID, 'reject')}
                        loading={processingId === row.ID}
                        disabled={!!processingId}
                    >
                        Rechazar
                    </Button>
                </div>
            )
        });
    }

    const totalItems = requests.length;
    const totalPages = Math.ceil(totalItems / itemsPerPage);
    const paginatedData = requests.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage);

    return (
        <div className={`bg-white dark:bg-gray-800 rounded-xl border border-gray-100 dark:border-gray-700 shadow-sm dark:shadow-lg overflow-hidden flex flex-col ${compact ? 'p-2' : 'pt-4 border-t border-gray-100 mt-8'}`}>
            {!compact && <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">{title}</h2>}
            {compact && <div className="px-4 py-2 bg-gray-50 dark:bg-gray-700 border-b border-gray-100 dark:border-gray-600 font-bold text-gray-700 dark:text-gray-100 text-sm uppercase tracking-wider">{title}</div>}
            <Table
                columns={requestColumns}
                data={paginatedData}
                loading={loading}
                emptyMessage={emptyMessage}
                pagination={{
                    currentPage,
                    totalPages,
                    totalItems,
                    itemsPerPage,
                    onPageChange: setCurrentPage,
                    onItemsPerPageChange,
                    itemsPerPageOptions: [5, 10, 20, 50],
                    showItemsPerPageSelector: true
                }}
            />
        </div>
    );
}

interface BusinessWalletViewProps {
    /** Super admin: ID del negocio a gestionar. Undefined = negocio propio del usuario */
    businessId?: number;
    /** Super admin: nombre del negocio seleccionado para mostrar en la UI */
    businessName?: string;
}

function BusinessWalletView({ businessId, businessName }: BusinessWalletViewProps = {}) {
    const { permissions, isSuperAdmin } = usePermissions();
    const isSuperAdminView = !!businessId;

    const [wallet, setWallet] = useState<Wallet | null>(null);
    const [loading, setLoading] = useState(true);
    const [rechargeAmount, setRechargeAmount] = useState<string>('');
    const [showQrModal, setShowQrModal] = useState(false);
    const [showConfirmationModal, setShowConfirmationModal] = useState(false);
    const [showPaymentSelector, setShowPaymentSelector] = useState(false);
    const [showComingSoonModal, setShowComingSoonModal] = useState(false);
    const [comingSoonGatewayName, setComingSoonGatewayName] = useState('');
    const [message, setMessage] = useState<{ type: 'success' | 'warning' | 'error', text: string } | null>(null);
    const [processing, setProcessing] = useState(false);
    const [currentRequestId, setCurrentRequestId] = useState<string | null>(null);
    const [qrCode, setQrCode] = useState<string | null>(null);
    const [boldProcessing, setBoldProcessing] = useState<{ orderId: string; amount: number; pollingEnabled: boolean } | null>(null);

    const QUICK_AMOUNTS = [15000, 50000, 100000, 200000, 500000];

    const [history, setHistory] = useState<any[]>([]);

    const fetchHistory = useCallback(async () => {
        try {
            const res = await getWalletHistoryAction(businessId);
            if (res.success) {
                setHistory(res.data || []);
            }
        } catch (e) {
            console.error(e);
        }
    }, [businessId]);

    const fetchBalance = useCallback(async () => {
        try {
            const res = await getWalletBalanceAction(businessId);
            if (!res.success) throw new Error(res.error || 'Failed to fetch balance');
            setWallet(res.data || null);
        } catch (err: any) {
            console.error(err);
        } finally {
            setLoading(false);
        }
    }, [businessId]);

    useEffect(() => {
        setLoading(true);
        setWallet(null);
        setHistory([]);
        fetchBalance();
        fetchHistory();
    }, [fetchBalance, fetchHistory]);

    const handleRechargeRequest = () => {
        if (!rechargeAmount || isNaN(Number(rechargeAmount)) || Number(rechargeAmount) <= 0) {
            setMessage({ type: 'error', text: 'Ingrese un monto válido' });
            return;
        }

        setMessage(null);
        setShowPaymentSelector(true);
    };

    const handleSelectNequi = async () => {
        setShowPaymentSelector(false);
        setProcessing(true);
        setMessage(null);

        try {
            // Si no hay businessId (vista de usuario normal), usar el del permissions
            const targetBusinessId = businessId || permissions?.business_id;

            console.log('Iniciando recarga:', { amount: rechargeAmount, businessId: targetBusinessId });

            const res = await rechargeWalletAction(Number(rechargeAmount), targetBusinessId);

            console.log('Respuesta recarga:', res);

            if (!res.success) {
                throw new Error(res.error || 'Error al reportar pago');
            }

            if (res.data?.ID) {
                setCurrentRequestId(res.data.ID);
            }

            if (res.data?.qr_code) {
                setQrCode(res.data.qr_code);
            }

            setShowQrModal(true);
        } catch (err: any) {
            console.error('Error en handleSelectNequi:', err);
            setMessage({ type: 'error', text: err.message || 'Error al procesar la recarga' });
            setShowPaymentSelector(true);
        } finally {
            setProcessing(false);
        }
    };

    const handleSelectOtherGateway = (gatewayName: string) => {
        if (gatewayName.toLowerCase().includes('bold')) {
            handleSelectBold();
            return;
        }
        setShowPaymentSelector(false);
        setComingSoonGatewayName(gatewayName);
        setShowComingSoonModal(true);
    };

    const handleSelectBold = async () => {
        setShowPaymentSelector(false);
        setProcessing(true);
        setMessage({
            type: 'warning',
            text: 'Te llevaremos al checkout de Bold. La confirmación del pago puede tardar hasta 5 minutos en aparecer en tu billetera.',
        });

        try {
            const targetBusinessId = businessId || permissions?.business_id;
            const res = await getBoldSignatureAction(Number(rechargeAmount), targetBusinessId);

            if (!res?.success) {
                throw new Error(res?.message || 'Error al obtener firma de Bold');
            }

            const { order_id, currency, amount, hash, public_key, redirection_url, polling_enabled } = res.data;

            if (!window.hasOwnProperty('BoldCheckout')) {
                await new Promise<void>((resolve, reject) => {
                    const script = document.createElement('script');
                    script.src = 'https://checkout.bold.co/library/boldPaymentButton.js';
                    script.async = true;
                    script.onload = () => resolve();
                    script.onerror = () => reject(new Error('No se pudo cargar el script de Bold (revisa conexion)'));
                    document.body.appendChild(script);
                    setTimeout(() => reject(new Error('Timeout cargando script de Bold')), 10000);
                });
            }

            const checkoutConfig: Record<string, unknown> = {
                orderId: order_id,
                currency,
                amount,
                apiKey: public_key,
                integritySignature: hash,
                description: `Recarga de Billetera Probability - Orden ${order_id}`,
            };
            if (redirection_url) {
                checkoutConfig.redirectionUrl = redirection_url;
            }
            // @ts-expect-error BoldCheckout is loaded from external script
            const checkout = new BoldCheckout(checkoutConfig);

            checkout.open();
            setBoldProcessing({ orderId: order_id, amount, pollingEnabled: !!polling_enabled });
        } catch (err: any) {
            console.error('Bold error:', err);
            setMessage({ type: 'error', text: err.message || 'Error al iniciar pago con Bold' });
        } finally {
            setProcessing(false);
        }
    };

    if (loading && !wallet) return <div className="p-8 text-center"><Spinner /></div>;

    const displayName = businessName || permissions?.business_name || '....';

    return (
        <>
            {/* Banner de contexto para super admin */}
            {isSuperAdminView && (
                <div className="mb-6 flex items-center gap-3 rounded-lg px-4 py-3" style={{ backgroundColor: 'var(--color-primary-50)', border: `1px solid var(--color-primary-200)` }}>
                    <svg className="w-5 h-5 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" style={{ color: 'var(--color-primary-600)' }}>
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <p className="text-sm" style={{ color: 'var(--color-primary-900)' }}>
                        Vista de billetera de <strong>{displayName}</strong> (ID: {businessId}) — modo super admin
                    </p>
                </div>
            )}

            <div className="grid gap-8 lg:grid-cols-2">
                {/* Virtual Card - Pixel Perfect from Pencil Design */}
                <div className="flex">
                    <VirtualCard
                        balance={wallet?.Balance || 0}
                        cardLastFour={wallet?.ID?.slice(-4) || '7842'}
                        brand={displayName || 'ProbabilityIA'}
                        isActive={true}
                        brandTag="FINTECH"
                    />
                </div>

                {/* Recharge Section */}
                <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 shadow-sm dark:shadow-lg p-6 lg:p-8">
                    <h2 className="text-xl font-bold text-gray-900 dark:text-white mb-6">
                        {isSuperAdminView ? `Recargar Saldo — ${displayName}` : 'Recargar Saldo'}
                    </h2>

                    <div className="space-y-6">
                        <Input
                            label="Monto a recargar"
                            type="number"
                            placeholder=" Ej: 50000"
                            value={rechargeAmount}
                            onChange={(e) => setRechargeAmount(e.target.value)}
                            helperText="El monto mínimo es de $15.000"
                            leftIcon={<span className="text-gray-500 dark:text-gray-400 font-bold"> </span>}
                        />

                        {/* Quick Amounts Chips */}
                        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-2">
                            {QUICK_AMOUNTS.map((amt) => (
                                <button
                                    key={amt}
                                    onClick={() => setRechargeAmount(String(amt))}
                                    className={`px-2 py-2 text-xs font-medium rounded-lg border transition-all ${Number(rechargeAmount) === amt
                                        ? 'text-white shadow-md transform scale-105'
                                        : 'bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-200 border-gray-200 dark:border-gray-600 dark:hover:bg-gray-600'
                                        }`}
                                    style={Number(rechargeAmount) === amt ? {
                                        backgroundColor: 'var(--color-tertiary-500)',
                                        borderColor: 'var(--color-tertiary-500)'
                                    } : {
                                        borderColor: Number(rechargeAmount) === amt ? 'var(--color-tertiary-300)' : 'inherit'
                                    }}
                                    onMouseEnter={(e) => {
                                        if (Number(rechargeAmount) !== amt) {
                                            (e.target as HTMLButtonElement).style.borderColor = 'var(--color-tertiary-300)';
                                            (e.target as HTMLButtonElement).style.backgroundColor = 'var(--color-tertiary-50)';
                                        }
                                    }}
                                    onMouseLeave={(e) => {
                                        if (Number(rechargeAmount) !== amt) {
                                            (e.target as HTMLButtonElement).style.borderColor = '#e5e7eb';
                                            (e.target as HTMLButtonElement).style.backgroundColor = 'white';
                                        }
                                    }}
                                >
                                    $ {amt / 1000}k
                                </button>
                            ))}
                        </div>

                        {rechargeAmount && (
                            <Alert type="warning">
                                <span className="text-xs">
                                    Debe consignar <strong>exactamente</strong> el valor ingresado. El saldo se actualizará tras validación.
                                </span>
                            </Alert>
                        )}

                        <button
                            onClick={handleRechargeRequest}
                            disabled={processing}
                            className="w-full py-3 text-lg font-semibold rounded-lg transition-all shadow-lg disabled:opacity-50 disabled:cursor-not-allowed text-white"
                            style={{
                                backgroundColor: rechargeAmount ? 'var(--color-tertiary-500)' : '#505050ff',
                                borderColor: rechargeAmount ? 'var(--color-tertiary-500)' : '#505050ff',
                                boxShadow: rechargeAmount ? '0 10px 15px -3px var(--color-tertiary-200)' : '0 10px 15px -3px rgba(59, 130, 246, 0.2)'
                            }}
                        >
                            {processing ? (
                                <div className="flex items-center justify-center gap-2">
                                    <div className="spinner w-4 h-4" />
                                    <span>Procesando...</span>
                                </div>
                            ) : (
                                'Proceder al Pago'
                            )}
                        </button>

                        {message && (
                            <Alert type={message.type} onClose={() => setMessage(null)}>
                                {message.text}
                            </Alert>
                        )}
                    </div>
                </div>
            </div>

            {/* QR Payment Modal */}
            <Modal
                isOpen={showQrModal}
                onClose={() => setShowQrModal(false)}
                showCloseButton={false}
                title="Escanea para Pagar"
                size="md"
            >
                <div className="flex flex-col items-center justify-center p-2 text-center">
                    {/* QR Nequi - Imagen estática */}
                    <div className="bg-white dark:bg-gray-800 p-4 rounded-xl border border-gray-100 dark:border-gray-700 mb-4 w-full max-w-[380px] flex justify-center">
                        <img
                            src="https://images-cam93.s3.us-east-1.amazonaws.com/QR_Cuenta_de_probability.jpeg"
                            alt="Nequi QR - Probability"
                            className="w-full h-auto object-contain"
                        />
                    </div>

                    <h3 className="text-xl font-bold text-gray-900 dark:text-white mb-1">
                        {rechargeAmount ? formatCurrency(Number(rechargeAmount)) : '$ --'}
                    </h3>
                    <p className="text-gray-500 dark:text-gray-400 mb-4 text-xs">
                        Total a pagar vía Nequi/Bancolombia
                    </p>

                    <div className="w-full border rounded-lg p-3 mb-4 text-left" style={{ backgroundColor: 'var(--color-primary-50)', borderColor: 'var(--color-primary-200)' }}>
                        <h4 className="font-semibold text-xs mb-1" style={{ color: 'var(--color-primary-900)' }}>Siguientes pasos:</h4>
                        <ul className="list-disc list-inside text-[11px] space-y-0.5" style={{ color: 'var(--color-primary-900)' }}>
                            <li>Escanea el código QR desde tu App bancaria.</li>
                            <li>Verifica que el monto y la llave sean correctos.</li>
                            <li>Realiza el pago.</li>
                            <li>Tu saldo se verá reflejado cuando Nequi confirme la transacción.</li>
                        </ul>
                    </div>

                    <Button
                        variant="primary"
                        onClick={async () => {
                            setShowQrModal(false);
                            setShowConfirmationModal(true);
                            setCurrentRequestId(null);
                        }}
                        className="w-full mb-2"
                    >
                        Ya generé el pago
                    </Button>

                    <Button
                        variant="primary"
                        onClick={() => {
                            setShowQrModal(false);
                            setRechargeAmount('');
                            setQrCode(null);
                        }}
                        className="w-full py-2 text-sm"
                    >
                        Regresar
                    </Button>
                </div>
            </Modal>

            {/* Payment Method Selector Modal */}
            <PaymentMethodSelectorModal
                isOpen={showPaymentSelector}
                onClose={() => setShowPaymentSelector(false)}
                amount={rechargeAmount}
                onSelectNequi={handleSelectNequi}
                onSelectOther={handleSelectOtherGateway}
            />

            {/* Bold Payment Processing Modal (SSE) */}
            <BoldPaymentProcessingModal
                open={boldProcessing !== null}
                orderId={boldProcessing?.orderId || ''}
                amount={boldProcessing?.amount || 0}
                businessId={businessId}
                pollingEnabled={boldProcessing?.pollingEnabled ?? false}
                onClose={() => {
                    setBoldProcessing(null);
                    fetchBalance();
                    fetchHistory();
                }}
                onResolved={(status) => {
                    if (status === 'success') {
                        fetchBalance();
                        fetchHistory();
                        setMessage({ type: 'success', text: 'Pago confirmado por Bold. Saldo actualizado.' });
                    } else if (status === 'failed') {
                        setMessage({ type: 'error', text: 'Bold rechazó el pago.' });
                    } else if (status === 'timeout') {
                        setMessage({ type: 'warning', text: 'Pago en proceso. Aparecerá en el historial cuando Bold confirme.' });
                    }
                }}
            />

            {/* Coming Soon Modal */}
            <Modal
                isOpen={showComingSoonModal}
                onClose={() => setShowComingSoonModal(false)}
                title="Próximamente"
                size="sm"
            >
                <div className="flex flex-col items-center text-center p-4 space-y-4">
                    <div className="w-16 h-16 rounded-full flex items-center justify-center" style={{ backgroundColor: 'var(--color-tertiary-100)' }}>
                        <span className="text-3xl">🚀</span>
                    </div>
                    <div>
                        <p className="font-semibold text-gray-900 dark:text-white mb-1">
                            {comingSoonGatewayName} estará disponible próximamente
                        </p>
                        <p className="text-sm text-gray-500 dark:text-gray-400">
                            Por ahora puedes usar Nequi para recargar tu billetera.
                        </p>
                    </div>
                    <Button
                        variant="secondary"
                        className="w-full"
                        onClick={() => {
                            setShowComingSoonModal(false);
                            setShowPaymentSelector(true);
                        }}
                    >
                        Volver a métodos de pago
                    </Button>
                </div>
            </Modal>

            {/* Confirmation Modal */}
            <Modal
                isOpen={showConfirmationModal}
                onClose={() => {
                    setShowConfirmationModal(false);
                    setRechargeAmount('');
                }}
                size="md"
            >
                <div className="flex flex-col items-center justify-center p-6 text-center">
                    <div className="w-20 h-20 rounded-full flex items-center justify-center mb-6" style={{ backgroundColor: '#dcfce7' }}>
                        <svg className="w-12 h-12" fill="none" viewBox="0 0 24 24" stroke="currentColor" style={{ color: '#16a34a' }}>
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                        </svg>
                    </div>

                    <h3 className="text-2xl font-bold text-gray-900 dark:text-white mb-3">
                        ¡Pago Reportado!
                    </h3>

                    <p className="text-gray-600 dark:text-gray-300 mb-6 text-base leading-relaxed">
                        Tu pago está en estado <strong>PENDIENTE</strong> de revisión. Será acreditado a tu billetera cuando Nequi confirme la transacción o un administrador lo apruebe manualmente.
                    </p>

                    <Button
                        variant="primary"
                        onClick={() => {
                            setShowConfirmationModal(false);
                            setRechargeAmount('');
                        }}
                        className="w-full"
                    >
                        Cerrar
                    </Button>
                </div>
            </Modal>

            {/* Business Transaction History */}
            <div className="mt-12 space-y-8">
                <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                    <h2 className="text-2xl font-bold text-gray-900 dark:text-white">Historial de Transacciones</h2>
                </div>

                <div className="space-y-8">
                    <HistoryTable
                        title="Transacciones Recientes / Procesadas"
                        data={history.filter(t => t.Status === 'COMPLETED' || t.Status === 'FAILED')}
                        emptyMessage="No hay transacciones procesadas"
                    />
                    <HistoryTable
                        title="Transacciones Pendientes"
                        data={history.filter(t => t.Status === 'PENDING')}
                        emptyMessage="No hay transacciones pendientes"
                    />
                </div>
            </div>
        </>
    );
}

function HistoryTable({ title, data, emptyMessage }: { title: string, data: any[], emptyMessage: string }) {
    const [currentPage, setCurrentPage] = useState(1);
    const [itemsPerPage, setItemsPerPage] = useState(10);

    const columns: TableColumn<any>[] = [
        {
            key: 'CreatedAt',
            label: 'Fecha',
            render: (val) => <span className="text-gray-600 dark:text-gray-300 font-mono text-sm">{formatDate(val as string)}</span>
        },
        {
            key: 'Reference',
            label: 'Referencia',
            render: (val) => <span className="text-gray-500 dark:text-gray-400 text-sm">{(val as string) || '---'}</span>
        },
        {
            key: 'integration_name',
            label: 'Metodo',
            render: (val, row) => {
                const r = row as any;
                const name = (val as string) || r.integration_name || '---';
                const imgUrl = r.integration_image_url as string | undefined;
                if (name === '---') return <span className="text-gray-400 text-sm">---</span>;
                const colorMap: Record<string, { bg: string; text: string; border: string }> = {
                    'Bold': { bg: 'var(--color-tertiary-50)', text: 'var(--color-tertiary-900)', border: 'var(--color-tertiary-200)' },
                    'Nequi': { bg: 'var(--color-quaternary-50)', text: 'var(--color-quaternary-900)', border: 'var(--color-quaternary-200)' },
                    'Debito manual': { bg: '#f3f4f6', text: '#1f2937', border: '#e5e7eb' },
                };
                const colors = colorMap[name] || { bg: '#f3f4f6', text: '#1f2937', border: '#e5e7eb' };
                return (
                    <span className="inline-flex items-center gap-1.5 px-2 py-1 rounded-full text-xs font-medium border" style={{ backgroundColor: colors.bg, color: colors.text, borderColor: colors.border }}>
                        {imgUrl && (
                            // eslint-disable-next-line @next/next/no-img-element
                            <img src={imgUrl} alt={name} className="h-4 w-4 object-contain" />
                        )}
                        {name}
                    </span>
                );
            }
        },
        {
            key: 'Amount',
            label: 'Monto',
            render: (val, row) => {
                const amount = val as number;
                const type = (row as any).Type as string;
                const isDebit = type === 'USAGE';
                return (
                    <span className="font-bold" style={{ color: isDebit ? '#dc2626' : '#16a34a' }}>
                        {isDebit ? '-' : '+'}{formatCurrency(amount)}
                    </span>
                );
            }
        },
        {
            key: 'Status',
            label: 'Estado',
            render: (val) => {
                const colors: Record<string, { bg: string; text: string }> = {
                    'PENDING': { bg: '#fef3c7', text: '#92400e' },
                    'COMPLETED': { bg: '#dcfce7', text: '#166534' },
                    'FAILED': { bg: '#fee2e2', text: '#991b1b' }
                };
                const color = (colors as any)[val as string] || { bg: '#f3f4f6', text: '#1f2937' };
                return (
                    <span className="px-2 py-1 rounded-full text-xs font-medium" style={{ backgroundColor: color.bg, color: color.text }}>
                        {val === 'PENDING' ? 'Pendiente' : val === 'COMPLETED' ? 'Completado' : 'Rechazado'}
                    </span>
                );
            }
        }
    ];

    const totalItems = data.length;
    const totalPages = Math.ceil(totalItems / itemsPerPage);
    const paginatedData = data.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage);

    return (
        <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-100 dark:border-gray-700 overflow-hidden">
            <div className="p-4 border-b border-gray-50 dark:border-gray-700 bg-gray-50/50 dark:bg-gray-700/50">
                <h3 className="font-semibold text-gray-900 dark:text-white">{title}</h3>
            </div>
            <Table
                columns={columns}
                data={paginatedData}
                emptyMessage={emptyMessage}
                pagination={{
                    currentPage,
                    totalPages,
                    totalItems,
                    itemsPerPage,
                    onPageChange: setCurrentPage,
                    onItemsPerPageChange: setItemsPerPage,
                    itemsPerPageOptions: [5, 10, 20, 50]
                }}
            />
        </div>
    );
}
