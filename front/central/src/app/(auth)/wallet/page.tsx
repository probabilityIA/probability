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
    Wallet
} from '@/services/modules/wallet/infra/actions';
import { getBusinessesAction } from '@/services/auth/business/infra/actions';

const formatCurrency = (amount: number) =>
    new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP' }).format(amount);

const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return `${date.toLocaleDateString()} ${date.toLocaleTimeString()}`;
};

export default function WalletPage() {
    const { isSuperAdmin } = usePermissions();

    return (
        <div className="min-h-screen bg-gray-50 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
                <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">Billetera</h1>
            </div>

            <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
                <div className="p-6">
                    {isSuperAdmin ? <AdminWalletView /> : <BusinessWalletView />}
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

    const fetchWalletsAndBusinesses = useCallback(async () => {
        try {
            setLoading(true);
            // Fetch Wallets
            const walletRes = await getWalletsAction();
            if (!walletRes.success) throw new Error(walletRes.error || 'Failed to fetch wallets');
            setWallets(walletRes.data || []);

            // Fetch Businesses
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
                    columns={walletColumns}
                    data={wallets}
                    loading={loading}
                    emptyMessage="No hay billeteras registradas"
                />
            </div>

            {/* In Review Section (Renamed from Pending) */}
            <RequestsTableView
                title="Pagos para revisión"
                businesses={businesses}
                onRequestsChanged={fetchWalletsAndBusinesses}
                allWallets={wallets}
                fetchAction={getPendingRequestsAction}
                showActions={true}
                emptyMessage="No hay pagos para revisión"
            />

            {/* Approved Section */}
            <RequestsTableView
                title="Pagos aprobados"
                businesses={businesses}
                onRequestsChanged={fetchWalletsAndBusinesses}
                allWallets={wallets}
                fetchAction={getProcessedRequestsAction}
                showActions={false}
                emptyMessage="No hay pagos aprobados"
            />
        </div>
    );
}

function RequestsTableView({
    title,
    businesses,
    onRequestsChanged,
    allWallets,
    fetchAction,
    showActions,
    emptyMessage
}: {
    title: string,
    businesses: Record<number, string>,
    onRequestsChanged: () => void,
    allWallets: Wallet[],
    fetchAction: () => Promise<any>,
    showActions: boolean,
    emptyMessage: string
}) {
    const [requests, setRequests] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);
    const [processingId, setProcessingId] = useState<string | null>(null);

    const fetchRequests = useCallback(async () => {
        try {
            setLoading(true);
            const res = await fetchAction();
            if (res.success) {
                setRequests(res.data as any[] || []);
            }
        } catch (error) {
            console.error(error);
        } finally {
            setLoading(false);
        }
    }, [fetchAction]);

    useEffect(() => {
        fetchRequests();
    }, [fetchRequests]);

    const handleAction = async (id: string, action: 'approve' | 'reject') => {
        setProcessingId(id);
        try {
            const res = await processRequestAction(id, action);
            if (res.success) {
                await fetchRequests(); // Refresh list
                onRequestsChanged(); // Notify parent
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
                <div className="flex gap-2">
                    <Button
                        size="sm"
                        variant="success"
                        onClick={() => handleAction(row.ID, 'approve')}
                        loading={processingId === row.ID}
                        disabled={!!processingId}
                    >
                        Aprobar
                    </Button>
                    <Button
                        size="sm"
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

    return (
        <div className="pt-4 border-t border-gray-100 mt-8">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">{title}</h2>
            <Table
                columns={requestColumns}
                data={requests}
                loading={loading}
                emptyMessage={emptyMessage}
            />
        </div>
    );
}

function BusinessWalletView() {
    const { permissions } = usePermissions();
    const [wallet, setWallet] = useState<Wallet | null>(null);
    const [loading, setLoading] = useState(true);
    const [rechargeAmount, setRechargeAmount] = useState<string>('');
    const [showQrModal, setShowQrModal] = useState(false);
    const [showConfirmationModal, setShowConfirmationModal] = useState(false);
    const [message, setMessage] = useState<{ type: 'success' | 'warning' | 'error', text: string } | null>(null);
    const [processing, setProcessing] = useState(false);
    const [currentRequestId, setCurrentRequestId] = useState<string | null>(null);

    const QUICK_AMOUNTS = [15000, 50000, 100000, 200000, 500000];

    const fetchBalance = useCallback(async () => {
        try {
            const res = await getWalletBalanceAction();
            if (!res.success) throw new Error(res.error || 'Failed to fetch balance');
            setWallet(res.data || null);
        } catch (err: any) {
            console.error(err);
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        fetchBalance();
    }, [fetchBalance]);

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
            const res = await rechargeWalletAction(amount);

            if (!res.success) {
                throw new Error(res.error || 'Error al reportar pago');
            }

            // Save request ID for later reporting
            if (res.data?.ID) {
                setCurrentRequestId(res.data.ID);
            }

            // Show standardized success message in modal or alert, here we use the QR modal as "Next Step"
            setShowQrModal(true);
            // Keep rechargeAmount to show in modal

        } catch (err: any) {
            setMessage({ type: 'error', text: err.message });
        } finally {
            setProcessing(false);
        }
    };

    if (loading && !wallet) return <div className="p-8 text-center"><Spinner /></div>;

    return (
        <>
            <div className="grid gap-8 lg:grid-cols-2">
                {/* Balance Card - Premium Design */}
                <div className="relative overflow-hidden bg-gradient-to-br from-gray-900 to-gray-800 rounded-2xl p-8 text-white shadow-xl flex flex-col justify-between min-h-[240px]">
                    <div className="absolute top-0 right-0 p-4 opacity-10">
                        <svg width="200" height="200" viewBox="0 0 24 24" fill="currentColor">
                            <path d="M21 18v1c0 1.1-.9 2-2 2H5c-1.11 0-2-.9-2-2V5c0-1.1.89-2 2-2h14c1.1 0 2 .9 2 2v1h-9c-1.11 0-2 .9-2 2v8c0 1.1.89 2 2 2h9zm-9-2h10V8H12v8zm4-2.5c-.83 0-1.5-.67-1.5-1.5s.67-1.5 1.5-1.5 1.5.67 1.5 1.5-.67 1.5-1.5 1.5z" />
                        </svg>
                    </div>

                    <div className="relative z-10 flex justify-between items-start">
                        <div>
                            <p className="text-gray-400 text-sm font-medium tracking-wider uppercase mb-1">Saldo Disponible</p>
                            <h3 className="text-4xl font-bold tracking-tight text-white drop-shadow-md">
                                {wallet ? formatCurrency(wallet.Balance) : '$ 0.00'}
                            </h3>
                        </div>
                        <div className="bg-white/10 backdrop-blur-md p-2 rounded-lg border border-white/20">
                            <svg className="w-8 h-8 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />
                            </svg>
                        </div>
                    </div>

                    <div className="relative z-10 mt-auto pt-8">
                        <div className="flex items-center gap-2 mb-2">
                            <div className="h-6 w-10 bg-yellow-500/80 rounded flex items-center justify-center overflow-hidden relative">
                                <div className="absolute w-8 h-8 rounded-full border border-yellow-300 top-[-4px] left-[-8px]"></div>
                            </div>
                            <span className="text-gray-400 text-xs tracking-widest">BILLETERA EMPRESARIAL</span>
                        </div>
                        <p className="font-mono text-gray-300 tracking-wider">
                            {permissions?.business_name || '....'}
                        </p>
                    </div>
                </div>

                {/* Recharge Section */}
                <div className="bg-white rounded-xl border border-gray-200 shadow-sm p-6 lg:p-8">
                    <h2 className="text-xl font-bold text-gray-900 mb-6">Recargar Saldo</h2>

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
                    {/* Success Icon */}
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
        </>
    );
}
