'use client';

import { useState, useEffect, useCallback } from 'react';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { TokenStorage } from '@/shared/config';
import { Spinner } from '@/shared/ui';
import { BusinessApiRepository } from '@/services/auth/business/infra/repository/api-repository';
import {
    getWalletsAction,
    getPendingRequestsAction,
    processRequestAction,
    getWalletBalanceAction,
    rechargeWalletAction,
    Wallet
} from '@/services/modules/wallet/infra/actions';

const formatCurrency = (amount: number) =>
    new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP' }).format(amount);

export default function WalletPage() {
    const { isSuperAdmin } = usePermissions();

    return (
        <div className="p-6">
            <h1 className="text-2xl font-bold mb-6 text-gray-900">Billetera</h1>
            {isSuperAdmin ? <AdminWalletView /> : <BusinessWalletView />}
        </div>
    );
}

function AdminWalletView() {
    const [wallets, setWallets] = useState<Wallet[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [businesses, setBusinesses] = useState<Record<number, string>>({});

    useEffect(() => {
        const fetchWalletsAndBusinesses = async () => {
            try {
                // Fetch Wallets via Server Action
                const walletRes = await getWalletsAction();
                if (!walletRes.success) throw new Error(walletRes.error || 'Failed to fetch wallets');
                setWallets(walletRes.data || []);

                // Fetch Businesses (keeps using existing API repo, but could also be moved to action if needed, 
                // but let's stick to wallets for now as requested)
                const token = TokenStorage.getSessionToken();
                const businessRepo = new BusinessApiRepository(token);
                const businessesRes = await businessRepo.getBusinesses({ per_page: 1000 });
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
        };
        fetchWalletsAndBusinesses();
    }, []);

    if (loading) return <Spinner />;
    if (error) return <div className="text-red-500">{error}</div>;

    return (
        <div className="bg-gray-800 rounded-lg p-4">
            {/* Wallets Table */}
            <h2 className="text-lg font-semibold mb-4 text-white">Saldos de Negocios</h2>
            <div className="overflow-x-auto mb-8">
                <table className="min-w-full text-left text-sm whitespace-nowrap">
                    <thead className="uppercase tracking-wider border-b border-gray-700 bg-gray-800">
                        <tr>
                            <th scope="col" className="px-6 py-3 text-gray-400">Negocio</th>
                            <th scope="col" className="px-6 py-3 text-gray-400">Saldo</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-700">
                        {wallets.map((wallet) => (
                            <tr key={wallet.ID} className="hover:bg-gray-700">
                                <td className="px-6 py-4 font-medium text-white">
                                    {businesses[wallet.BusinessID] || `ID: ${wallet.BusinessID}`}
                                </td>
                                <td className="px-6 py-4 text-green-400 font-bold">{formatCurrency(wallet.Balance)}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>

            {/* Pending Requests Table */}
            <PendingRequestsView businesses={businesses} onRequestsChanged={() => {
                // Trigger refresh if needed
                // For now, we reload the page or we could extract fetchWalletsAndBusinesses to be reusable
                // Ideally, we should update the wallet list too if a request is approved.
                window.location.reload();
            }}
                allWallets={wallets}
            />
        </div>
    );
}

function PendingRequestsView({ businesses, onRequestsChanged, allWallets }: { businesses: Record<number, string>, onRequestsChanged: () => void, allWallets: Wallet[] }) {
    const [requests, setRequests] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);

    const fetchRequests = useCallback(async () => {
        try {
            const res = await getPendingRequestsAction();
            if (res.success) {
                setRequests(res.data as any[] || []);
            }
        } catch (error) {
            console.error(error);
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        fetchRequests();
    }, [fetchRequests]);

    const handleAction = async (id: string, action: 'approve' | 'reject') => {
        try {
            const res = await processRequestAction(id, action);
            if (res.success) {
                fetchRequests(); // Refresh list
                onRequestsChanged(); // Notify parent
                alert(`Solicitud ${action === 'approve' ? 'APROBADA' : 'RECHAZADA'} correctamente.`);
            } else {
                alert(`Error al procesar la solicitud: ${res.error}`);
            }
        } catch (e) {
            alert("Error de conexión.");
        }
    };

    if (loading) return <div className="mt-6 p-4 text-gray-400">Cargando solicitudes...</div>;

    if (requests.length === 0) {
        return (
            <div className="bg-gray-800 rounded-lg p-4 mt-6 border border-gray-700 bg-opacity-50">
                <h2 className="text-lg font-semibold mb-2 text-yellow-500">Solicitudes de Recarga Pendientes</h2>
                <p className="text-gray-400 text-sm">No hay solicitudes pendientes por revisar.</p>
            </div>
        );
    }

    return (
        <div className="bg-gray-800 rounded-lg p-4 mt-6 border border-gray-700 bg-opacity-50">
            <h2 className="text-lg font-semibold mb-4 text-yellow-500">Solicitudes de Recarga Pendientes</h2>
            <div className="overflow-x-auto">
                <table className="min-w-full text-left text-sm whitespace-nowrap">
                    <thead className="uppercase tracking-wider border-b border-gray-700 bg-gray-800">
                        <tr>
                            <th scope="col" className="px-6 py-3 text-gray-400">Fecha</th>
                            <th scope="col" className="px-6 py-3 text-gray-400">Negocio</th>
                            <th scope="col" className="px-6 py-3 text-gray-400">Monto</th>
                            <th scope="col" className="px-6 py-3 text-gray-400">Acciones</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-700">
                        {requests.map((req) => {
                            return (
                                <tr key={req.ID} className="hover:bg-gray-700">
                                    <td className="px-6 py-4 text-white">
                                        {new Date(req.CreatedAt).toLocaleDateString()} {new Date(req.CreatedAt).toLocaleTimeString()}
                                    </td>
                                    <td className="px-6 py-4 text-white font-medium">
                                        <WalletBusinessName walletId={req.WalletID} businesses={businesses} wallets={allWallets} />
                                    </td>
                                    <td className="px-6 py-4 text-green-400 font-bold">{formatCurrency(req.Amount)}</td>
                                    <td className="px-6 py-4">
                                        <div className="flex gap-2">
                                            <button
                                                onClick={() => handleAction(req.ID, 'approve')}
                                                className="bg-green-600 hover:bg-green-700 text-white px-3 py-1 rounded text-xs"
                                            >
                                                Aprobar
                                            </button>
                                            <button
                                                onClick={() => handleAction(req.ID, 'reject')}
                                                className="bg-red-600 hover:bg-red-700 text-white px-3 py-1 rounded text-xs"
                                            >
                                                Rechazar
                                            </button>
                                        </div>
                                    </td>
                                </tr>
                            );
                        })}
                    </tbody>
                </table>
            </div>
        </div>
    );
}

// Helper to find business name from wallet ID
function WalletBusinessName({ walletId, businesses, wallets }: { walletId: string, businesses: Record<number, string>, wallets: Wallet[] }) {
    const [name, setName] = useState("...");

    useEffect(() => {
        const wallet = wallets.find(w => w.ID === walletId);
        if (wallet) {
            const businessName = businesses[wallet.BusinessID];
            setName(businessName ? `${businessName} (ID: ${wallet.BusinessID})` : `Negocio ${wallet.BusinessID}`);
        } else {
            setName(walletId.substring(0, 8) + "...");
        }
    }, [walletId, businesses, wallets]);

    return <span>{name}</span>;
}

function BusinessWalletView() {
    const [wallet, setWallet] = useState<Wallet | null>(null);
    const [loading, setLoading] = useState(true);
    const [rechargeAmount, setRechargeAmount] = useState<string>('');
    const [showQr, setShowQr] = useState(false);
    const [message, setMessage] = useState<{ type: 'success' | 'alert' | 'error', text: string } | null>(null);

    const fetchBalance = async () => {
        try {
            const res = await getWalletBalanceAction();
            if (!res.success) throw new Error(res.error || 'Failed to fetch balance');
            setWallet(res.data || null);
        } catch (err: any) {
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchBalance();
    }, []);

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

        setLoading(true);
        setMessage(null);

        try {
            const res = await rechargeWalletAction(amount);

            if (!res.success) {
                throw new Error(res.error || 'Error al reportar pago');
            }

            // Success
            setMessage({
                type: 'success',
                text: 'Pago reportado exitosamente. Se reflejará en su cuenta en aproximadamente 2 horas una vez confirmado.'
            });
            setShowQr(true);
            setRechargeAmount('');

        } catch (err: any) {
            setMessage({ type: 'error', text: err.message });
        } finally {
            setLoading(false);
        }
    };

    if (loading && !wallet) return <Spinner />;

    return (
        <div className="grid gap-6 md:grid-cols-2">
            {/* Balance Card */}
            <div className="bg-gray-800 rounded-lg p-6 shadow-lg">
                <h2 className="text-lg font-semibold text-gray-400 mb-2">Saldo Disponible</h2>
                <div className="text-4xl font-bold text-green-400">
                    {wallet ? formatCurrency(wallet.Balance) : '$ 0.00'}
                </div>
            </div>

            {/* Recharge Card */}
            <div className="bg-gray-800 rounded-lg p-6 shadow-lg">
                <h2 className="text-lg font-semibold text-white mb-4">Recargar Saldo</h2>
                <div className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-300 mb-1">Monto a recargar (Mínimo $15.000)</label>
                        <input
                            type="number"
                            value={rechargeAmount}
                            onChange={(e) => setRechargeAmount(e.target.value)}
                            className="w-full bg-white border border-gray-300 rounded px-3 py-2 text-black font-medium focus:outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-500/50"
                            placeholder="Ej: 50000"
                        />
                    </div>

                    <div className="bg-yellow-900/30 border border-yellow-700/50 p-3 rounded text-sm text-yellow-200 mb-3">
                        <p className="font-bold mb-1">⚠ Importante</p>
                        <p>Debe consignar <strong>exactamente</strong> el valor ingresado. El saldo se actualizará tras validación administrativa.</p>
                    </div>

                    <button
                        onClick={handleRechargeRequest}
                        disabled={loading && !!wallet}
                        className={`w-full py-2 px-4 rounded font-bold text-white transition-colors ${loading && !!wallet ? 'bg-gray-600' : 'bg-blue-600 hover:bg-blue-700'}`}
                    >
                        {loading && !!wallet ? 'Procesando...' : 'Reportar Pago y Ver QR'}
                    </button>

                    {message && (
                        <div className={`p-3 rounded text-sm ${message.type === 'success' ? 'bg-green-900/50 text-green-200' :
                            message.type === 'alert' ? 'bg-yellow-900/50 text-yellow-200' :
                                'bg-red-900/50 text-red-200'
                            }`}>
                            {message.text}
                        </div>
                    )}

                    {showQr && (
                        <div className="mt-4 flex flex-col items-center p-4 bg-white rounded-lg">
                            {/* Static QR Image from public folder */}
                            <img src="/QR.png" alt="Nequi QR" className="max-w-[200px] h-auto" />
                            <p className="mt-2 text-gray-800 text-sm font-medium text-center">
                                Escanea para pagar
                            </p>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}

