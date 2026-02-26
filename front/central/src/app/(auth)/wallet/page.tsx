'use client';

import { useState, useEffect, useCallback } from 'react';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { Spinner, Button, Input, Table, TableColumn, Alert, Modal } from '@/shared/ui';
import {
    getWalletsAction,
    getPendingRequestsAction,
    getProcessedRequestsAction,
    processRequestAction,
    getWalletBalanceAction,
    rechargeWalletAction,
    reportPaymentAction,
    manualDebitAction,
    getWalletHistoryAction,
    clearRechargeHistoryAction,
    Wallet
} from '@/services/modules/wallet/infra/actions';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';
import { VirtualCard } from './virtual-card';

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
        <div className="min-h-screen bg-gray-50 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
                <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">Billetera</h1>

                {/* Business selector - solo para super admin */}
                {isSuperAdmin && businesses.length > 0 && (
                    <div className="flex items-center gap-3 bg-blue-50 border border-blue-200 rounded-lg px-4 py-2">
                        <label className="text-sm font-medium text-blue-800 whitespace-nowrap">
                            Negocio:
                        </label>
                        <select
                            value={selectedBusinessId?.toString() ?? ''}
                            onChange={(e) => {
                                const val = e.target.value;
                                setSelectedBusinessId(val ? Number(val) : null);
                            }}
                            className="px-3 py-1.5 border border-blue-300 rounded-md text-sm bg-white focus:outline-none focus:ring-2 focus:ring-blue-500 min-w-[200px]"
                        >
                            <option value="">Vista Administrativa</option>
                            {businesses.map((b) => (
                                <option key={b.id} value={b.id}>{b.name}</option>
                            ))}
                        </select>
                    </div>
                )}
            </div>

            <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
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
            setError(err.message);
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
                <span className="font-medium text-gray-900">
                    {businesses[row.BusinessID] || `ID: ${row.BusinessID}`}
                </span>
            )
        },
        {
            key: 'Balance',
            label: 'Saldo',
            render: (val) => (
                <span className="font-bold text-green-600">
                    {formatCurrency(val as number)}
                </span>
            )
        }
    ];

    if (error) return <Alert type="error">{error}</Alert>;

    return (
        <div className="space-y-8">
            {/* Wallets Section */}
            <div>
                <h2 className="text-lg font-semibold text-gray-900 mb-4">Saldos de Negocios</h2>
                <Table
                    columns={[
                        ...walletColumns,
                        {
                            key: 'actions',
                            label: 'Acciones',
                            render: (_, row) => (
                                <div className="flex gap-2">
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
            <div className="bg-white rounded-xl border border-gray-100 shadow-sm overflow-hidden flex flex-col">
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
            </div>

            {/* Bottom Row - Approved and Rejected (Side by Side) */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <div className="bg-white rounded-xl border border-gray-100 shadow-sm overflow-hidden flex flex-col">
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
                </div>

                <div className="bg-white rounded-xl border border-gray-100 shadow-sm overflow-hidden flex flex-col">
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
            const res = await clearRechargeHistoryAction(businessId);
            if (res.success) {
                setIsOpen(false);
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
            <Button size="sm" variant="outline" className="text-red-600 border-red-200 hover:bg-red-50" onClick={() => setIsOpen(true)}>
                Borrar Historial
            </Button>
            <Modal isOpen={isOpen} onClose={() => setIsOpen(false)} title="Confirmar Eliminación">
                <div className="space-y-4 p-4">
                    <p className="text-gray-600">
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
            render: (val) => <span className="text-gray-600 font-mono text-sm">{formatDate(val as string)}</span>
        },
        {
            key: 'WalletID',
            label: 'Negocio',
            render: (val) => {
                const wallet = allWallets.find(w => w.ID === val);
                if (wallet) {
                    const name = businesses[wallet.BusinessID];
                    return <span className="font-medium text-gray-900">{name || `ID: ${wallet.BusinessID}`}</span>;
                }
                return <span className="text-gray-500">...</span>;
            }
        },
        {
            key: 'Amount',
            label: 'Monto',
            render: (val) => <span className="font-bold text-gray-900">{formatCurrency(val as number)}</span>
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
        <div className={`${compact ? 'p-2' : 'pt-4 border-t border-gray-100 mt-8'}`}>
            {!compact && <h2 className="text-lg font-semibold text-gray-900 mb-4">{title}</h2>}
            {compact && <div className="px-4 py-2 bg-gray-50 border-b border-gray-100 font-bold text-gray-700 text-sm uppercase tracking-wider">{title}</div>}
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
    const { permissions } = usePermissions();
    const isSuperAdminView = !!businessId;

    const [wallet, setWallet] = useState<Wallet | null>(null);
    const [loading, setLoading] = useState(true);
    const [rechargeAmount, setRechargeAmount] = useState<string>('');
    const [showQrModal, setShowQrModal] = useState(false);
    const [showConfirmationModal, setShowConfirmationModal] = useState(false);
    const [message, setMessage] = useState<{ type: 'success' | 'warning' | 'error', text: string } | null>(null);
    const [processing, setProcessing] = useState(false);
    const [currentRequestId, setCurrentRequestId] = useState<string | null>(null);

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

    const handleRechargeRequest = async () => {
        if (!rechargeAmount || isNaN(Number(rechargeAmount))) {
            setMessage({ type: 'error', text: 'Ingrese un monto válido' });
            return;
        }

        const amount = Number(rechargeAmount);
        if (amount < 15000) {
            setMessage({ type: 'error', text: 'El monto mínimo de recarga es de $15.000' });
            return;
        }

        setProcessing(true);
        setMessage(null);

        try {
            const res = await rechargeWalletAction(amount, businessId);

            if (!res.success) {
                throw new Error(res.error || 'Error al reportar pago');
            }

            if (res.data?.ID) {
                setCurrentRequestId(res.data.ID);
            }

            setShowQrModal(true);

        } catch (err: any) {
            setMessage({ type: 'error', text: err.message });
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
                <div className="mb-6 flex items-center gap-3 bg-blue-50 border border-blue-200 rounded-lg px-4 py-3">
                    <svg className="w-5 h-5 text-blue-600 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <p className="text-sm text-blue-800">
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
                <div className="bg-white rounded-xl border border-gray-200 shadow-sm p-6 lg:p-8">
                    <h2 className="text-xl font-bold text-gray-900 mb-6">
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
                            leftIcon={<span className="text-gray-500 font-bold"> </span>}
                        />

                        {/* Quick Amounts Chips */}
                        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-2">
                            {QUICK_AMOUNTS.map((amt) => (
                                <button
                                    key={amt}
                                    onClick={() => setRechargeAmount(String(amt))}
                                    className={`px-2 py-2 text-xs font-medium rounded-lg border transition-all ${Number(rechargeAmount) === amt
                                        ? 'bg-[#7c3aed] text-white border-[#7c3aed] shadow-md transform scale-105'
                                        : 'bg-white text-gray-700 border-gray-200 hover:border-purple-300 hover:bg-purple-50'
                                        }`}
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
                                backgroundColor: rechargeAmount ? '#7c3aed' : '#505050ff',
                                borderColor: rechargeAmount ? '#7c3aed' : '#505050ff',
                                boxShadow: rechargeAmount ? '0 10px 15px -3px rgba(124, 58, 237, 0.2)' : '0 10px 15px -3px rgba(59, 130, 246, 0.2)'
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
                    <div className="bg-white p-2 rounded-xl border border-gray-100 mb-4 w-full max-w-[200px] flex justify-center">
                        <img src="/QR.png" alt="Nequi QR" className="w-full h-auto object-contain mix-blend-multiply" />
                    </div>

                    <h3 className="text-xl font-bold text-gray-900 mb-1">
                        {rechargeAmount ? formatCurrency(Number(rechargeAmount)) : '$ --'}
                    </h3>
                    <p className="text-gray-500 mb-4 text-xs">
                        Total a pagar vía Nequi/Bancolombia
                    </p>

                    <div className="w-full bg-blue-50 border border-blue-100 rounded-lg p-3 mb-4 text-left">
                        <h4 className="font-semibold text-blue-900 text-xs mb-1">Siguientes pasos:</h4>
                        <ul className="list-disc list-inside text-[11px] text-blue-800 space-y-0.5">
                            <li>Escanea el código QR desde tu App bancaria.</li>
                            <li>Verifica que el monto sea exacto.</li>
                            <li>Realiza el pago.</li>
                            <li>Tu saldo se verá reflejado en aprox. 2 horas.</li>
                        </ul>
                    </div>

                    <Button
                        variant="secondary"
                        onClick={async () => {
                            if (currentRequestId) {
                                await reportPaymentAction(currentRequestId);
                            }
                            setShowQrModal(false);
                            setShowConfirmationModal(true);
                            setCurrentRequestId(null);
                        }}
                        className="w-full mb-2 bg-[#7c3aed] hover:bg-[#6d28d9] border-[#7c3aed] py-2 text-sm text-white"
                    >
                        Ya generé el pago
                    </Button>

                    <Button
                        variant="primary"
                        onClick={() => {
                            setShowQrModal(false);
                            setRechargeAmount('');
                        }}
                        className="w-full py-2 text-sm"
                    >
                        Regresar
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
                    <div className="w-20 h-20 bg-green-100 rounded-full flex items-center justify-center mb-6">
                        <svg className="w-12 h-12 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                        </svg>
                    </div>

                    <h3 className="text-2xl font-bold text-gray-900 mb-3">
                        ¡Pago Reportado!
                    </h3>

                    <p className="text-gray-600 mb-6 text-base leading-relaxed">
                        Revisaremos su pago y será acreditado en unos minutos.
                    </p>

                    <Button
                        variant="primary"
                        onClick={() => {
                            setShowConfirmationModal(false);
                            setRechargeAmount('');
                        }}
                        className="w-full bg-[#7c3aed] hover:bg-[#6d28d9] border-[#7c3aed]"
                    >
                        Cerrar
                    </Button>
                </div>
            </Modal>

            {/* Business Transaction History */}
            <div className="mt-12 space-y-8">
                <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                    <h2 className="text-2xl font-bold text-gray-900">Historial de Transacciones</h2>
                    <div className="flex gap-4">
                        <div className="bg-green-50 px-4 py-2 rounded-lg border border-green-100">
                            <p className="text-xs text-green-700 font-medium">Total Aprobado</p>
                            <p className="text-lg font-bold text-green-800">
                                {formatCurrency(history.filter(t => t.Status === 'COMPLETED').reduce((acc, t) => acc + t.Amount, 0))}
                            </p>
                        </div>
                        <div className="bg-yellow-50 px-4 py-2 rounded-lg border border-yellow-100">
                            <p className="text-xs text-yellow-700 font-medium">Total Pendiente</p>
                            <p className="text-lg font-bold text-yellow-800">
                                {formatCurrency(history.filter(t => t.Status === 'PENDING').reduce((acc, t) => acc + t.Amount, 0))}
                            </p>
                        </div>
                    </div>
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
            render: (val) => <span className="text-gray-600 font-mono text-sm">{formatDate(val as string)}</span>
        },
        {
            key: 'Reference',
            label: 'Referencia',
            render: (val) => <span className="text-gray-500 text-sm">{(val as string) || '---'}</span>
        },
        {
            key: 'Amount',
            label: 'Monto',
            render: (val) => <span className="font-bold text-gray-900">{formatCurrency(val as number)}</span>
        },
        {
            key: 'Status',
            label: 'Estado',
            render: (val) => {
                const colors = {
                    'PENDING': 'bg-yellow-100 text-yellow-800',
                    'COMPLETED': 'bg-green-100 text-green-800',
                    'FAILED': 'bg-red-100 text-red-800'
                };
                return (
                    <span className={`px-2 py-1 rounded-full text-xs font-medium ${(colors as any)[val as string] || 'bg-gray-100'}`}>
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
        <div className="bg-white rounded-xl border border-gray-100 overflow-hidden">
            <div className="p-4 border-b border-gray-50 bg-gray-50/50">
                <h3 className="font-semibold text-gray-900">{title}</h3>
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
