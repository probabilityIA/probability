'use client';

import { useState, useEffect } from 'react';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { TokenStorage } from '@/shared/config';
import { Spinner } from '@/shared/ui';

interface Wallet {
    ID: string;
    BusinessID: number;
    Balance: number;
}

const formatCurrency = (amount: number) =>
    new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP' }).format(amount);

export default function WalletPage() {
    const { isSuperAdmin } = usePermissions();

    return (
        <div className="p-6">
            <h1 className="text-2xl font-bold mb-6 text-white">Billetera</h1>
            {isSuperAdmin ? <AdminWalletView /> : <BusinessWalletView />}
        </div>
    );
}

function AdminWalletView() {
    const [wallets, setWallets] = useState<Wallet[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchWallets = async () => {
            try {
                const token = TokenStorage.getSessionToken();
                const res = await fetch('http://localhost:3050/api/v1/wallet/all', {
                    headers: {
                        'Authorization': `Bearer ${token}`
                    }
                });
                if (!res.ok) throw new Error('Failed to fetch wallets');
                const data = await res.json();
                setWallets(data || []);
            } catch (err: any) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };
        fetchWallets();
    }, []);

    if (loading) return <Spinner />;
    if (error) return <div className="text-red-500">{error}</div>;

    return (
        <div className="bg-gray-800 rounded-lg p-4">
            <h2 className="text-lg font-semibold mb-4 text-white">Saldos de Negocios</h2>
            <div className="overflow-x-auto">
                <table className="min-w-full text-left text-sm whitespace-nowrap">
                    <thead className="uppercase tracking-wider border-b border-gray-700 bg-gray-800">
                        <tr>
                            <th scope="col" className="px-6 py-3 text-gray-400">ID Negocio</th>
                            <th scope="col" className="px-6 py-3 text-gray-400">Saldo</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-700">
                        {wallets.map((wallet) => (
                            <tr key={wallet.ID} className="hover:bg-gray-700">
                                <td className="px-6 py-4 font-medium text-white">{wallet.BusinessID}</td>
                                <td className="px-6 py-4 text-green-400 font-bold">{formatCurrency(wallet.Balance)}</td>
                            </tr>
                        ))}
                        {wallets.length === 0 && (
                            <tr>
                                <td colSpan={2} className="px-6 py-4 text-center text-gray-500">No hay billeteras activas.</td>
                            </tr>
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
}

function BusinessWalletView() {
    const [wallet, setWallet] = useState<Wallet | null>(null);
    const [loading, setLoading] = useState(true);
    const [rechargeAmount, setRechargeAmount] = useState<string>('');
    const [qrCode, setQrCode] = useState<string | null>(null);
    const [processing, setProcessing] = useState(false);
    const [message, setMessage] = useState<{ type: 'success' | 'error', text: string } | null>(null);

    const fetchBalance = async () => {
        try {
            const token = TokenStorage.getSessionToken();
            // Use BusinessToken if available for business-specific requests?
            // Based on sidebar/layout logic, session token might be enough or we need business token.
            // Layout says: "Usuario business... el token de sesion ya tiene permisios" but sidebar uses specialized tokens?
            // Lets try session token first.
            const res = await fetch('http://localhost:3050/api/v1/wallet/balance', {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
            if (!res.ok) throw new Error('Failed to fetch balance');
            const data = await res.json();
            setWallet(data);
        } catch (err: any) {
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchBalance();
    }, []);

    const handleRecharge = async () => {
        if (!rechargeAmount || isNaN(Number(rechargeAmount)) || Number(rechargeAmount) <= 0) {
            setMessage({ type: 'error', text: 'Ingrese un monto válido' });
            return;
        }
        setProcessing(true);
        setQrCode(null);
        setMessage(null);

        try {
            const token = TokenStorage.getSessionToken();
            const res = await fetch('http://localhost:3050/api/v1/wallet/recharge', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({ amount: Number(rechargeAmount) })
            });

            if (!res.ok) {
                const errData = await res.json();
                throw new Error(errData.error || 'Error en recarga');
            }

            const data = await res.json();
            setQrCode(data.qr_code);
            setMessage({ type: 'success', text: 'Código QR generado. Escanéelo para pagar.' });

            // Refresh balance immediately as per requirement "numeric value stored as active balance"
            fetchBalance();
        } catch (err: any) {
            setMessage({ type: 'error', text: err.message });
        } finally {
            setProcessing(false);
        }
    };

    if (loading) return <Spinner />;

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
                        <label className="block text-sm font-medium text-gray-300 mb-1">Monto a recargar</label>
                        <input
                            type="number"
                            value={rechargeAmount}
                            onChange={(e) => setRechargeAmount(e.target.value)}
                            className="w-full bg-white border border-gray-300 rounded px-3 py-2 text-black font-medium focus:outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-500/50"
                            placeholder="Ej: 50000"
                        />
                    </div>
                    <button
                        onClick={handleRecharge}
                        disabled={processing}
                        className={`w-full py-2 px-4 rounded font-bold text-white transition-colors ${processing ? 'bg-gray-600 cursor-not-allowed' : 'bg-blue-600 hover:bg-blue-700'
                            }`}
                    >
                        {processing ? 'Generando QR...' : 'Generar QR Nequi'}
                    </button>

                    {message && (
                        <div className={`p-3 rounded text-sm ${message.type === 'success' ? 'bg-green-900/50 text-green-200' : 'bg-red-900/50 text-red-200'}`}>
                            {message.text}
                        </div>
                    )}

                    {qrCode && (
                        <div className="mt-4 flex flex-col items-center p-4 bg-white rounded-lg">
                            <img src={`https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(qrCode)}`} alt="Nequi QR" />
                            <p className="mt-2 text-gray-800 text-sm font-medium break-all text-center">{qrCode}</p>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
