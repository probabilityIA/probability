'use client';

import { useState, useEffect, useCallback } from 'react';
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
import { CONCEPT_LABELS, CONCEPT_OPTIONS } from '@/services/modules/wallet/domain/concept';
import { getWalletKPISelectionAction, updateWalletKPISelectionAction } from '@/services/modules/wallet/infra/actions/wallet-kpi-selection';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useWalletBusiness } from '@/shared/contexts/wallet-business-context';
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

const getColorFromHash = (text: string): string => {
    const colors = ['#7c3aed', '#06b6d4', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899', '#14b8a6'];
    let hash = 0;
    for (let i = 0; i < text.length; i++) {
        hash = text.charCodeAt(i) + ((hash << 5) - hash);
    }
    return colors[Math.abs(hash) % colors.length];
};

export function AdminWalletView() {
    const { setSelectedBusinessId } = useWalletBusiness();
    const [wallets, setWallets] = useState<Wallet[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [businesses, setBusinesses] = useState<Record<number, string>>({});
    const [itemsPerPage, setItemsPerPage] = useState(10);
    const [searchBusiness, setSearchBusiness] = useState('');
    const [activeTab, setActiveTab] = useState<'review' | 'approved' | 'rejected'>('approved');
    const [showBusinessSelector, setShowBusinessSelector] = useState(false);
    const [selectedBusinessesForKPI, setSelectedBusinessesForKPI] = useState<Set<number>>(new Set());
    const [savingKPI, setSavingKPI] = useState(false);
    const [saveMessage, setSaveMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

    const fetchWalletsAndBusinesses = useCallback(async () => {
        try {
            setLoading(true);
            const walletRes = await getWalletsAction();
            if (!walletRes.success) throw new Error(walletRes.error || 'Failed to fetch wallets');
            const fetchedWallets = walletRes.data || [];
            setWallets(fetchedWallets);

            const { getBusinessesAction } = await import('@/services/auth/business/infra/actions');
            const businessesRes = await getBusinessesAction({ per_page: 10000 });
            if (businessesRes.data) {
                const businessMap: Record<number, string> = {};
                businessesRes.data.forEach((b: any) => {
                    businessMap[b.id] = b.name;
                });

                fetchedWallets.forEach((wallet: any) => {
                    if (!businessMap[wallet.BusinessID]) {
                        businessMap[wallet.BusinessID] = '';
                    }
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
    }, []);

    useEffect(() => {
        const loadKPISelection = async () => {
            const res = await getWalletKPISelectionAction();
            if (res.success && Array.isArray(res.data?.selected_business_ids)) {
                const ids = res.data.selected_business_ids;
                if (ids.length > 0) {
                    setSelectedBusinessesForKPI(new Set(ids));
                } else if (wallets.length > 0) {
                    const allIds = wallets.map(w => w.BusinessID);
                    setSelectedBusinessesForKPI(new Set(allIds));
                    await updateWalletKPISelectionAction(allIds);
                }
            } else if (wallets.length > 0) {
                const allIds = wallets.map(w => w.BusinessID);
                setSelectedBusinessesForKPI(new Set(allIds));
                await updateWalletKPISelectionAction(allIds);
            }
        };

        if (wallets.length > 0) {
            loadKPISelection();
        }
    }, [wallets]);

    const handleSaveKPISelection = async () => {
        setSavingKPI(true);
        setSaveMessage(null);
        try {
            const res = await updateWalletKPISelectionAction(Array.from(selectedBusinessesForKPI));
            if (res.success && res.data?.selected_business_ids) {
                setSaveMessage({ type: 'success', text: 'Selección guardada exitosamente ✓' });
                setTimeout(() => {
                    setShowBusinessSelector(false);
                    setSaveMessage(null);
                }, 1500);
            } else {
                setSaveMessage({ type: 'error', text: 'Error al guardar. Intenta de nuevo.' });
            }
        } catch (err) {
            setSaveMessage({ type: 'error', text: 'Error de conexión. Intenta de nuevo.' });
        } finally {
            setSavingKPI(false);
        }
    };

    const filteredWallets = wallets.filter(w => {
        const businessName = businesses[w.BusinessID] || '';
        return businessName.toLowerCase().includes(searchBusiness.toLowerCase());
    });

    const totalBalance = wallets
        .filter(w => selectedBusinessesForKPI.has(w.BusinessID))
        .reduce((sum, w) => sum + (typeof w.Balance === 'string' ? parseFloat(w.Balance) : w.Balance), 0);
    const activeBusinesses = wallets
        .filter(w => selectedBusinessesForKPI.has(w.BusinessID))
        .filter(w => (typeof w.Balance === 'string' ? parseFloat(w.Balance) : w.Balance) > 0).length;
    const selectedBusinessesCount = selectedBusinessesForKPI.size;

    const walletColumns: TableColumn<Wallet>[] = [
        {
            key: 'BusinessID',
            label: 'Negocio',
            render: (_val, row) => {
                const businessName = businesses[row.BusinessID] || `ID: ${row.BusinessID}`;
                const initial = businessName.charAt(0).toUpperCase();
                const bgColor = getColorFromHash(businessName);
                return (
                    <div className="flex items-center gap-3">
                        <div
                            className="w-8 h-8 rounded-lg flex items-center justify-center text-white text-sm font-bold"
                            style={{ backgroundColor: bgColor }}
                        >
                            {initial}
                        </div>
                        <span className="font-medium text-gray-900 dark:text-white">{businessName}</span>
                    </div>
                );
            }
        },
        {
            key: 'Balance',
            label: 'Saldo',
            render: (val) => {
                const balance = typeof val === 'string' ? parseFloat(val) : (val as number);
                const isNegative = balance < 0;
                return (
                    <span className="font-bold font-mono" style={{ color: isNegative ? '#dc2626' : (balance === 0 ? '#aab1c2' : '#16a34a') }}>
                        {isNegative && '-'}${formatCurrency(Math.abs(balance)).replace('$', '')}
                    </span>
                );
            }
        }
    ];

    if (error) return <Alert type="error">{error}</Alert>;

    return (
        <div className="space-y-8">
            <div>
                <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2" style={{ letterSpacing: '-0.03em' }}>
                    Saldos de Negocios
                </h1>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                <div className="bg-gradient-to-br from-blue-900 to-purple-900 rounded-2xl p-6 text-white relative cursor-pointer hover:shadow-lg transition-shadow" onClick={() => setShowBusinessSelector(!showBusinessSelector)}>
                    <div className="flex items-center justify-between mb-2">
                        <div className="text-sm font-medium opacity-80">Saldo Total en la Red</div>
                        <svg className="w-4 h-4 opacity-80" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                        </svg>
                    </div>
                    <div className="text-3xl font-bold font-mono">${formatCurrency(totalBalance).replace('$', '')}</div>
                    <div className="text-xs mt-3 opacity-70 flex items-center gap-1">
                        <span className="text-green-300">▲ 8.4%</span>
                        <span>vs. mes anterior</span>
                    </div>
                    {showBusinessSelector && (
                        <div className="absolute top-full left-0 mt-2 w-80 bg-white dark:bg-gray-800 rounded-lg shadow-xl z-50 p-4 border border-gray-200 dark:border-gray-700" onClick={(e) => e.stopPropagation()}>
                            <div className="mb-3">
                                <h3 className="font-semibold text-gray-900 dark:text-white mb-3">Selecciona negocios</h3>
                                <div className="max-h-64 overflow-y-auto space-y-2">
                                    {wallets.map(wallet => (
                                        <label key={wallet.BusinessID} className="flex items-center gap-2 cursor-pointer p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded">
                                            <input
                                                type="checkbox"
                                                checked={selectedBusinessesForKPI.has(wallet.BusinessID)}
                                                onChange={(e) => {
                                                    e.stopPropagation();
                                                    const newSet = new Set(selectedBusinessesForKPI);
                                                    if (e.target.checked) {
                                                        newSet.add(wallet.BusinessID);
                                                    } else {
                                                        newSet.delete(wallet.BusinessID);
                                                    }
                                                    setSelectedBusinessesForKPI(newSet);
                                                }}
                                                className="w-4 h-4 rounded"
                                            />
                                            <span className="text-sm text-gray-900 dark:text-white">
                                                {businesses[wallet.BusinessID] || `ID: ${wallet.BusinessID}`}
                                            </span>
                                        </label>
                                    ))}
                                </div>
                            </div>
                            {saveMessage && (
                                <div className={`mb-3 p-2 rounded text-xs font-medium text-center ${saveMessage.type === 'success' ? 'bg-green-100 dark:bg-green-900 text-green-700 dark:text-green-300' : 'bg-red-100 dark:bg-red-900 text-red-700 dark:text-red-300'}`}>
                                    {saveMessage.text}
                                </div>
                            )}
                            <div className="flex gap-2">
                                <button
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        setShowBusinessSelector(false);
                                    }}
                                    className="flex-1 px-3 py-2 text-xs font-medium text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white bg-gray-100 dark:bg-gray-700 rounded"
                                >
                                    Cancelar
                                </button>
                                <button
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        handleSaveKPISelection();
                                    }}
                                    disabled={savingKPI}
                                    className="flex-1 px-3 py-2 text-xs font-medium text-white bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 rounded"
                                >
                                    {savingKPI ? 'Guardando...' : 'Guardar'}
                                </button>
                            </div>
                        </div>
                    )}
                </div>

                <div className="bg-white dark:bg-gray-800 rounded-2xl p-6 border border-gray-200 dark:border-gray-700">
                    <div className="text-sm font-medium text-gray-600 dark:text-gray-400 mb-3">Negocios Activos</div>
                    <div className="text-2xl font-bold text-gray-900 dark:text-white mb-3">{activeBusinesses} / {selectedBusinessesCount}</div>
                    <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                        <div
                            className="bg-gradient-to-r from-green-400 to-green-600 h-2 rounded-full"
                            style={{ width: `${selectedBusinessesCount > 0 ? (activeBusinesses / selectedBusinessesCount) * 100 : 0}%` }}
                        />
                    </div>
                    <div className="text-xs text-gray-500 mt-2">{selectedBusinessesCount > 0 ? Math.round((activeBusinesses / selectedBusinessesCount) * 100) : 0}% con saldo</div>
                </div>

                <div className="bg-white dark:bg-gray-800 rounded-2xl p-6 border border-gray-200 dark:border-gray-700">
                    <div className="text-sm font-medium text-gray-600 dark:text-gray-400 mb-2">En Revisión</div>
                    <div className="text-2xl font-bold text-yellow-600 mb-3">-</div>
                    <button className="text-sm font-semibold text-blue-600 hover:text-blue-700">
                        Revisar ahora →
                    </button>
                </div>

                <div className="bg-white dark:bg-gray-800 rounded-2xl p-6 border border-gray-200 dark:border-gray-700">
                    <div className="text-sm font-medium text-gray-600 dark:text-gray-400 mb-2">Recargas Aprobadas (Jun)</div>
                    <div className="text-2xl font-bold text-green-600 font-mono">-</div>
                    <div className="text-xs text-gray-500 mt-3">Suma total</div>
                </div>
            </div>

            <div className="bg-white dark:bg-gray-800 rounded-2xl border border-gray-200 dark:border-gray-700 shadow-sm overflow-hidden">
                <div className="px-6 py-5 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-700/50 flex items-center justify-between">
                    <div>
                        <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Saldos por Negocio</h2>
                        <p className="text-xs text-gray-500 mt-1">{filteredWallets.length} de {wallets.length} negocios</p>
                    </div>
                    <div className="relative">
                        <svg className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                        </svg>
                        <input
                            type="text"
                            placeholder="Buscar negocio..."
                            value={searchBusiness}
                            onChange={(e) => setSearchBusiness(e.target.value)}
                            className="pl-9 pr-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm w-48 focus:outline-none focus:ring-2 focus:ring-purple-500"
                        />
                    </div>
                </div>
                <div className="overflow-x-auto">
                    <Table
                        columns={[
                            ...walletColumns,
                            {
                                key: 'actions',
                                label: 'Acciones',
                                render: (_, row) => (
                                    <div className="flex gap-2" onClick={(e) => e.stopPropagation()}>
                                        <RechargeWalletButton
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
                        data={filteredWallets}
                        loading={loading}
                        emptyMessage="No hay billeteras registradas"
                        onRowClick={(row) => setSelectedBusinessId(row.BusinessID)}
                    />
                </div>
            </div>

            <div className="bg-white dark:bg-gray-800 rounded-2xl border border-gray-200 dark:border-gray-700 shadow-sm overflow-hidden">
                <div className="border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-700/50">
                    <div className="flex">
                        {['approved', 'review', 'rejected'].map((tab) => (
                            <button
                                key={tab}
                                onClick={() => setActiveTab(tab as any)}
                                className={`flex-1 px-6 py-4 font-semibold text-sm transition-colors border-b-2 ${
                                    activeTab === tab
                                        ? 'text-blue-600 border-blue-600 dark:text-blue-400 dark:border-blue-400'
                                        : 'text-gray-600 dark:text-gray-400 border-transparent hover:text-gray-900 dark:hover:text-gray-200'
                                }`}
                            >
                                {tab === 'review' && 'En Revisión'}
                                {tab === 'approved' && 'Aprobados'}
                                {tab === 'rejected' && 'Rechazados'}
                            </button>
                        ))}
                    </div>
                </div>

                <div className="overflow-x-auto">
                    {activeTab === 'approved' && (
                        <RequestsTableView
                            title="Aprobados"
                            businesses={businesses}
                            onRequestsChanged={fetchWalletsAndBusinesses}
                            allWallets={wallets}
                            fetchAction={getProcessedRequestsAction}
                            filterStatus="COMPLETED"
                            showActions={false}
                            emptyMessage="Sin aprobados"
                            compact={false}
                            itemsPerPage={itemsPerPage}
                            onItemsPerPageChange={setItemsPerPage}
                            hideTitle={false}
                            showFiltersOnly={true}
                        />
                    )}
                    {activeTab === 'review' && (
                        <RequestsTableView
                            title="En revisión"
                            businesses={businesses}
                            onRequestsChanged={fetchWalletsAndBusinesses}
                            allWallets={wallets}
                            fetchAction={getPendingRequestsAction}
                            showActions={true}
                            emptyMessage="¡Todo al día! No hay solicitudes pendientes"
                            compact={false}
                            itemsPerPage={itemsPerPage}
                            onItemsPerPageChange={setItemsPerPage}
                            hideTitle={true}
                        />
                    )}
                    {activeTab === 'rejected' && (
                        <RequestsTableView
                            title="Rechazados"
                            businesses={businesses}
                            onRequestsChanged={fetchWalletsAndBusinesses}
                            allWallets={wallets}
                            fetchAction={getProcessedRequestsAction}
                            filterStatus="FAILED"
                            showActions={false}
                            emptyMessage="Sin rechazados"
                            compact={false}
                            itemsPerPage={itemsPerPage}
                            onItemsPerPageChange={setItemsPerPage}
                            hideTitle={false}
                            showFiltersOnly={true}
                        />
                    )}
                </div>
            </div>
        </div>
    );
}

interface BusinessWalletViewProps {
    businessId?: number;
    businessName?: string;
}

export function BusinessWalletView({ businessId, businessName }: BusinessWalletViewProps = {}) {
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
    const [activeTab, setActiveTab] = useState<'all' | 'completed' | 'pending'>('all');
    const [histView, setHistView] = useState<'timeline' | 'table'>('timeline');
    const [historyPage, setHistoryPage] = useState(1);

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
            const targetBusinessId = businessId;

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
            const targetBusinessId = businessId;
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

    const groupTransactionsByDay = (txns: any[]) => {
        const grouped: { [key: string]: any[] } = {};
        txns.forEach(tx => {
            const date = new Date(tx.CreatedAt);
            const today = new Date();
            const yesterday = new Date(today);
            yesterday.setDate(yesterday.getDate() - 1);

            let dayLabel = '';
            if (date.toDateString() === today.toDateString()) {
                dayLabel = 'Hoy';
            } else if (date.toDateString() === yesterday.toDateString()) {
                dayLabel = 'Ayer';
            } else {
                dayLabel = date.toLocaleDateString('es-CO', { month: 'short', day: 'numeric' });
            }

            if (!grouped[dayLabel]) grouped[dayLabel] = [];
            grouped[dayLabel].push(tx);
        });
        return grouped;
    };

    const filteredHistory = activeTab === 'all'
        ? history
        : activeTab === 'completed'
        ? history.filter(t => t.Status === 'COMPLETED' || t.Status === 'FAILED')
        : history.filter(t => t.Status === 'PENDING');

    const groupedByDay = groupTransactionsByDay(filteredHistory);

    return (
        <>
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
                <div className="flex">
                    <VirtualCard
                        balance={wallet?.Balance || 0}
                        cardLastFour={wallet?.ID?.slice(-4) || '7842'}
                        brand={displayName || 'ProbabilityIA'}
                        isActive={true}
                        brandTag="FINTECH"
                    />
                </div>

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

            <Modal
                isOpen={showQrModal}
                onClose={() => setShowQrModal(false)}
                showCloseButton={false}
                title="Escanea para Pagar"
                size="md"
            >
                <div className="flex flex-col items-center justify-center p-2 text-center">
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

            <PaymentMethodSelectorModal
                isOpen={showPaymentSelector}
                onClose={() => setShowPaymentSelector(false)}
                amount={rechargeAmount}
                onSelectNequi={handleSelectNequi}
                onSelectOther={handleSelectOtherGateway}
            />

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

            <div className="mt-12 space-y-8">
                <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
                    <h2 className="text-2xl font-bold text-gray-900 dark:text-white" style={{ letterSpacing: '-0.02em' }}>
                        Historial de Transacciones
                    </h2>
                </div>

                <div className="bg-white dark:bg-gray-800 rounded-2xl p-6 border border-gray-200 dark:border-gray-700">
                    <div className="flex gap-4 mb-6 border-b border-gray-200 dark:border-gray-700 pb-4 justify-between items-center flex-wrap">
                        <div className="flex gap-3">
                            {(['all', 'completed', 'pending'] as const).map((filter) => (
                                <button
                                    key={filter}
                                    onClick={() => setActiveTab(filter as any)}
                                    className={`px-4 py-2 text-sm font-medium transition border-b-2 -mb-4 ${
                                        activeTab === filter
                                            ? 'border-violet-600 text-violet-600'
                                            : 'border-transparent text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
                                    }`}
                                >
                                    {filter === 'all' ? 'Todas' : filter === 'completed' ? 'Completadas' : 'Pendientes'}
                                </button>
                            ))}
                        </div>
                        <div className="flex gap-2 bg-gray-100 dark:bg-gray-700 p-1 rounded-lg">
                            {(['timeline', 'table'] as const).map((view) => (
                                <button
                                    key={view}
                                    onClick={() => setHistView(view)}
                                    className={`px-3 py-1 text-xs font-medium rounded transition ${
                                        histView === view
                                            ? 'bg-white dark:bg-gray-600 text-gray-900 dark:text-white shadow-sm'
                                            : 'text-gray-600 dark:text-gray-400'
                                    }`}
                                >
                                    {view === 'timeline' ? 'Timeline' : 'Tabla'}
                                </button>
                            ))}
                        </div>
                    </div>

                    {histView === 'timeline' ? (
                        <div className="space-y-6">
                            {(() => {
                                const itemsPerPage = 10;
                                const flatTxns = Object.entries(groupedByDay).flatMap(([, dayTxns]) => dayTxns);
                                const totalPages = Math.ceil(flatTxns.length / itemsPerPage);
                                const startIdx = (historyPage - 1) * itemsPerPage;
                                const endIdx = startIdx + itemsPerPage;
                                const paginatedTxns = flatTxns.slice(startIdx, endIdx);
                                const paginatedGrouped: { [key: string]: any[] } = {};
                                paginatedTxns.forEach(tx => {
                                    const date = new Date(tx.CreatedAt);
                                    const today = new Date();
                                    const yesterday = new Date(today);
                                    yesterday.setDate(yesterday.getDate() - 1);
                                    let dayLabel = '';
                                    if (date.toDateString() === today.toDateString()) {
                                        dayLabel = 'Hoy';
                                    } else if (date.toDateString() === yesterday.toDateString()) {
                                        dayLabel = 'Ayer';
                                    } else {
                                        dayLabel = date.toLocaleDateString('es-CO', { month: 'short', day: 'numeric' });
                                    }
                                    if (!paginatedGrouped[dayLabel]) paginatedGrouped[dayLabel] = [];
                                    paginatedGrouped[dayLabel].push(tx);
                                });

                                if (flatTxns.length === 0) {
                                    return <p className="text-center text-gray-500 dark:text-gray-400 py-8">No hay transacciones</p>;
                                }

                                return (
                                    <>
                                        <div className="space-y-6">
                                            {Object.entries(paginatedGrouped).map(([dayLabel, dayTxns]) => (
                                                <div key={dayLabel}>
                                                    <h4 className="text-xs font-semibold text-gray-500 dark:text-gray-400 mb-3 uppercase tracking-wider">{dayLabel}</h4>
                                                    <div className="space-y-3">
                                                        {dayTxns.map((tx, idx) => {
                                                const isIncome = tx.Type === 'RECHARGE';
                                                const statusColor = tx.Status === 'COMPLETED' ? '#16a34a' : tx.Status === 'FAILED' ? '#dc2626' : '#f59e0b';
                                                const icon = isIncome ? '✓' : '📦';
                                                const iconBg = isIncome ? '#dcfce7' : '#f5f5f5';
                                                const time = new Date(tx.CreatedAt).toLocaleTimeString('es-CO', { hour: '2-digit', minute: '2-digit' });
                                                const concept = tx.Concept ? (CONCEPT_LABELS[tx.Concept] || tx.Concept) : '';

                                                return (
                                                    <div key={idx} className="flex gap-4 items-center pb-3 border-b border-gray-100 dark:border-gray-700 last:border-0">
                                                        <div className="w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0 text-base" style={{ backgroundColor: iconBg }}>
                                                            <span>{icon}</span>
                                                        </div>
                                                        <div className="flex-1 min-w-0">
                                                            <p className="text-sm font-medium text-gray-900 dark:text-white">
                                                                {isIncome ? 'Recarga' : `Consumo de saldo${concept ? ` - ${concept}` : ''}`}
                                                            </p>
                                                            <div className="flex gap-2 items-center mt-1 flex-wrap">
                                                                <span className="text-xs text-gray-500 dark:text-gray-400 font-mono">{time}</span>
                                                                {tx.integration_name && (
                                                                    <span className="text-xs px-2 py-0.5 rounded bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300">
                                                                        {tx.integration_name}
                                                                    </span>
                                                                )}
                                                            </div>
                                                        </div>
                                                        <div className="text-right flex-shrink-0">
                                                            <p className="text-sm font-bold font-mono" style={{ color: isIncome ? '#16a34a' : '#dc2626' }}>
                                                                {isIncome ? '+' : '−'} {formatCurrency(Math.abs(tx.Amount))}
                                                            </p>
                                                            <p className="text-xs text-gray-500 dark:text-gray-400 mt-0.5" style={{ color: statusColor }}>
                                                                {tx.Status === 'COMPLETED' ? 'Completado' : tx.Status === 'PENDING' ? 'Pendiente' : 'Fallido'}
                                                            </p>
                                                        </div>
                                                    </div>
                                                );
                                            })}
                                                    </div>
                                                </div>
                                            ))}
                                        </div>
                                        <div className="flex items-center justify-between pt-4 border-t border-gray-200 dark:border-gray-700">
                                            <span className="text-xs text-gray-500 dark:text-gray-400">
                                                Página {historyPage} de {totalPages} ({flatTxns.length} transacciones)
                                            </span>
                                            <div className="flex gap-2">
                                                <button
                                                    onClick={() => setHistoryPage(p => Math.max(1, p - 1))}
                                                    disabled={historyPage === 1}
                                                    className="px-3 py-1 text-xs font-medium rounded border border-gray-300 dark:border-gray-600 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                                                >
                                                    ← Anterior
                                                </button>
                                                <button
                                                    onClick={() => setHistoryPage(p => Math.min(totalPages, p + 1))}
                                                    disabled={historyPage === totalPages}
                                                    className="px-3 py-1 text-xs font-medium rounded border border-gray-300 dark:border-gray-600 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                                                >
                                                    Siguiente →
                                                </button>
                                            </div>
                                        </div>
                                    </>
                                );
                            })()}
                        </div>
                    ) : (
                        <HistoryTable
                            title=""
                            data={filteredHistory}
                            emptyMessage={
                                activeTab === 'all'
                                    ? 'No hay transacciones'
                                    : activeTab === 'completed'
                                    ? 'No hay transacciones completadas'
                                    : 'No hay transacciones pendientes'
                            }
                        />
                    )}
                </div>
            </div>
        </>
    );
}

function RequestsTableAccordion({
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
    allWallets: any[],
    fetchAction: () => Promise<any>,
    showActions: boolean,
    emptyMessage: string,
    filterStatus?: string,
    compact?: boolean,
    itemsPerPage: number,
    onItemsPerPageChange: (total: number) => void
}) {
    const [isExpanded, setIsExpanded] = useState(false);

    return (
        <div className="border border-gray-200 dark:border-gray-700 rounded-2xl overflow-hidden bg-white dark:bg-gray-800 shadow-sm">
            <button
                onClick={() => setIsExpanded(!isExpanded)}
                className="w-full flex items-center justify-between px-6 py-4 bg-gray-50 dark:bg-gray-700/50 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors border-b border-gray-200 dark:border-gray-700"
            >
                <span className="font-semibold text-base text-gray-900 dark:text-white">{title}</span>
                <svg
                    className={`w-5 h-5 text-gray-600 dark:text-gray-400 transition-transform duration-300 ${isExpanded ? 'rotate-180' : ''}`}
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                >
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 14l-7 7m0 0l-7-7m7 7V3" />
                </svg>
            </button>
            {isExpanded && (
                <div className="overflow-x-auto">
                    <RequestsTableView
                        title={title}
                        businesses={businesses}
                        onRequestsChanged={onRequestsChanged}
                        allWallets={allWallets}
                        fetchAction={fetchAction}
                        showActions={showActions}
                        emptyMessage={emptyMessage}
                        filterStatus={filterStatus}
                        compact={compact}
                        itemsPerPage={itemsPerPage}
                        onItemsPerPageChange={onItemsPerPageChange}
                        hideTitle={true}
                    />
                </div>
            )}
        </div>
    );
}

function ManualDebitAccordion({ businessId, businessName, onSuccess, isSuperAdmin }: { businessId: number, businessName: string, onSuccess: () => void, isSuperAdmin: boolean }) {
    const [isExpanded, setIsExpanded] = useState(false);
    const [amount, setAmount] = useState('');
    const [reference, setReference] = useState('');
    const [concept, setConcept] = useState('');
    const [loading, setLoading] = useState(false);

    if (!isSuperAdmin) return null;

    const handleDebit = async () => {
        if (!amount || isNaN(Number(amount))) return;
        if (!concept) {
            alert('Selecciona una categoria');
            return;
        }
        setLoading(true);
        try {
            const res = await manualDebitAction(businessId, Number(amount), reference, concept);
            if (res.success) {
                setIsExpanded(false);
                setAmount('');
                setReference('');
                setConcept('');
                onSuccess();
                alert('Saldo restado exitosamente');
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
        <div className="border border-red-200 dark:border-red-900 rounded-lg overflow-hidden">
            <button
                onClick={() => setIsExpanded(!isExpanded)}
                className="w-full flex items-center justify-between p-4 bg-red-50 dark:bg-red-950/20 hover:bg-red-100 dark:hover:bg-red-950/30 transition-colors"
            >
                <span className="font-semibold text-red-700 dark:text-red-400">⚠️ Restar Saldo (Manual)</span>
                <svg
                    className={`w-5 h-5 text-red-700 dark:text-red-400 transition-transform ${isExpanded ? 'rotate-180' : ''}`}
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                >
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 14l-7 7m0 0l-7-7m7 7V3" />
                </svg>
            </button>
            {isExpanded && (
                <div className="p-4 bg-white dark:bg-gray-800 space-y-4 border-t border-red-200 dark:border-red-900">
                    <Alert type="warning">
                        Esta es una operación manual. Úsala solo si Bold está fuera de servicio.
                    </Alert>
                    <Input
                        label="Monto a restar"
                        type="number"
                        value={amount}
                        onChange={e => setAmount(e.target.value)}
                        placeholder="Ej: 5000"
                    />
                    <div>
                        <label className="text-xs font-semibold text-gray-600 dark:text-gray-400 block mb-2">Categoria</label>
                        <select
                            value={concept}
                            onChange={e => setConcept(e.target.value)}
                            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm focus:outline-none focus:ring-2 focus:ring-purple-500"
                        >
                            <option value="">Selecciona una categoria</option>
                            {CONCEPT_OPTIONS.map(opt => (
                                <option key={opt.value} value={opt.value}>{opt.label}</option>
                            ))}
                        </select>
                    </div>
                    <Input
                        label="Referencia / Motivo"
                        value={reference}
                        onChange={e => setReference(e.target.value)}
                        placeholder="Ej: Ajuste de saldo"
                    />
                    <div className="flex justify-end gap-2">
                        <Button variant="secondary" onClick={() => {
                            setIsExpanded(false);
                            setAmount('');
                            setReference('');
                            setConcept('');
                        }}>Cancelar</Button>
                        <Button variant="danger" onClick={handleDebit} loading={loading}>Restar Saldo</Button>
                    </div>
                </div>
            )}
        </div>
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
            <button
                onClick={() => setIsOpen(true)}
                className="px-4 py-2 text-sm font-semibold rounded-lg bg-white dark:bg-gray-800 border-2 transition-colors"
                style={{ color: '#dc2626', borderColor: '#dc2626' }}
            >
                Borrar Historial
            </button>
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
    const [mode, setMode] = useState<'add' | 'subtract'>('add');
    const [amount, setAmount] = useState('');
    const [reason, setReason] = useState('');
    const [concept, setConcept] = useState('RECHARGE');
    const [loading, setLoading] = useState(false);

    const isSubtract = mode === 'subtract';

    const reset = () => {
        setMode('add');
        setAmount('');
        setReason('');
        setConcept('RECHARGE');
    };

    const changeMode = (next: 'add' | 'subtract') => {
        setMode(next);
        setConcept(next === 'add' ? 'RECHARGE' : '');
    };

    const handleSubmit = async () => {
        if (!amount || isNaN(Number(amount)) || Number(amount) <= 0) {
            alert('Ingresa un monto válido (mayor a 0)');
            return;
        }
        if (!reason.trim()) {
            alert('Ingresa un motivo');
            return;
        }
        if (!concept) {
            alert('Selecciona una categoria');
            return;
        }
        const signedAmount = isSubtract ? -Number(amount) : Number(amount);
        setLoading(true);
        try {
            const res = await adminAdjustBalanceAction(businessId, signedAmount, reason, concept);
            if (res.success) {
                setIsOpen(false);
                reset();
                onSuccess();
                alert(isSubtract ? 'Saldo descontado exitosamente' : 'Saldo agregado exitosamente');
            } else {
                alert(res.error || 'Error al ajustar el saldo');
            }
        } catch (e: any) {
            alert(`Error al procesar: ${e.message || 'Error desconocido'}`);
        } finally {
            setLoading(false);
        }
    };

    return (
        <>
            <button
                onClick={() => setIsOpen(true)}
                className="px-4 py-2 text-sm font-semibold rounded-lg text-white hover:opacity-90 transition-opacity"
                style={{ backgroundColor: '#0f1729' }}
            >
                Ajustar saldo
            </button>
            <Modal isOpen={isOpen} onClose={() => { setIsOpen(false); reset(); }} title={`Ajustar saldo de ${businessName}`}>
                <div className="space-y-4 p-4">
                    <div className="grid grid-cols-2 gap-2 p-1 rounded-lg bg-gray-100 dark:bg-gray-800">
                        <button
                            type="button"
                            onClick={() => changeMode('add')}
                            className={`py-2 text-sm font-semibold rounded-md transition-colors ${!isSubtract ? 'bg-emerald-600 text-white' : 'text-gray-600 dark:text-gray-300'}`}
                        >
                            + Agregar
                        </button>
                        <button
                            type="button"
                            onClick={() => changeMode('subtract')}
                            className={`py-2 text-sm font-semibold rounded-md transition-colors ${isSubtract ? 'bg-red-600 text-white' : 'text-gray-600 dark:text-gray-300'}`}
                        >
                            - Descontar
                        </button>
                    </div>
                    <Input
                        label={isSubtract ? 'Monto a descontar' : 'Monto a agregar'}
                        type="number"
                        value={amount}
                        onChange={e => setAmount(e.target.value)}
                        placeholder="Ej: 10000"
                    />
                    <div>
                        <label className="text-xs font-semibold text-gray-600 dark:text-gray-400 block mb-2">Categoria</label>
                        <select
                            value={concept}
                            onChange={e => setConcept(e.target.value)}
                            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm focus:outline-none focus:ring-2 focus:ring-purple-500"
                        >
                            <option value="">Selecciona una categoria</option>
                            {CONCEPT_OPTIONS.map(opt => (
                                <option key={opt.value} value={opt.value}>{opt.label}</option>
                            ))}
                        </select>
                    </div>
                    <Input
                        label="Motivo"
                        type="text"
                        value={reason}
                        onChange={e => setReason(e.target.value)}
                        placeholder={isSubtract ? 'Ej: Cobro de mensualidad, uso extra, etc.' : 'Ej: Ajuste por error, promoción, etc.'}
                    />
                    <div className={`text-sm font-semibold ${isSubtract ? 'text-red-600' : 'text-emerald-600'}`}>
                        {amount && !isNaN(Number(amount)) && Number(amount) > 0
                            ? `${isSubtract ? '-' : '+'} $${Number(amount).toLocaleString('es-CO')}`
                            : ''}
                    </div>
                    <div className="flex justify-end gap-2">
                        <Button variant="secondary" onClick={() => { setIsOpen(false); reset(); }}>Cancelar</Button>
                        <Button variant={isSubtract ? 'danger' : 'primary'} onClick={handleSubmit} loading={loading}>
                            {isSubtract ? 'Descontar saldo' : 'Agregar saldo'}
                        </Button>
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
    onItemsPerPageChange,
    hideTitle = false,
    showFiltersOnly = false
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
    onItemsPerPageChange: (total: number) => void,
    hideTitle?: boolean,
    showFiltersOnly?: boolean
}) {
    const [requests, setRequests] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);
    const [processingId, setProcessingId] = useState<string | null>(null);
    const [currentPage, setCurrentPage] = useState(1);
    const [dateFrom, setDateFrom] = useState<string>('');
    const [dateTo, setDateTo] = useState<string>('');
    const [typeFilter, setTypeFilter] = useState<'all' | 'RECHARGE' | 'USAGE'>('all');
    const [businessFilter, setBusinessFilter] = useState<number | 'all'>('all');

    const fetchRequests = useCallback(async () => {
        try {
            setLoading(true);
            const res = await fetchAction();
            if (res.success) {
                let data = res.data as any[] || [];
                if (filterStatus) {
                    data = data.filter(r => r.Status === filterStatus);
                }
                if (dateFrom || dateTo) {
                    data = data.filter(r => {
                        const createdDate = r.CreatedAt ? new Date(r.CreatedAt).toLocaleDateString('en-CA') : '';

                        if (dateFrom && createdDate < dateFrom) {
                            return false;
                        }
                        if (dateTo && createdDate > dateTo) {
                            return false;
                        }
                        return true;
                    });
                }
                if (typeFilter !== 'all') {
                    data = data.filter(r => r.Type === typeFilter);
                }
                if (businessFilter !== 'all') {
                    data = data.filter(r => r.BusinessID === businessFilter);
                }
                setRequests(data);
            }
        } catch (error) {
            console.error(error);
        } finally {
            setLoading(false);
        }
    }, [fetchAction, filterStatus, dateFrom, dateTo, typeFilter, businessFilter]);

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
            render: (val, row) => {
                const request = row as any;
                let businessId: number | null = null;
                let businessName: string = '';

                if (request.BusinessID) {
                    businessId = request.BusinessID;
                    businessName = businesses[request.BusinessID] || '';
                }

                if (!businessName && val) {
                    const wallet = allWallets.find(w => w.ID === val);
                    if (wallet) {
                        businessId = wallet.BusinessID;
                        businessName = businesses[wallet.BusinessID] || '';
                    }
                }

                if (businessName) {
                    return <span className="font-medium text-gray-900 dark:text-white">{businessName}</span>;
                }

                if (businessId) {
                    return <span className="font-medium text-gray-900 dark:text-white">{businessId}</span>;
                }

                return <span className="text-gray-500 dark:text-gray-400">...</span>;
            }
        },
        {
            key: 'Reference',
            label: 'Referencia',
            width: '350px',
            render: (val) => {
                const reference = val as string;
                if (!reference) {
                    return <span className="text-gray-400 dark:text-gray-500">-</span>;
                }

                let methodType = 'OTRO';
                let badgeBg = 'bg-gray-200 dark:bg-gray-700';
                let badgeText = 'text-gray-800 dark:text-gray-200';

                if (reference.startsWith('WLT') || reference.startsWith('BOLD_SANDBOX_')) {
                    methodType = 'BOLD';
                    badgeBg = 'bg-orange-100 dark:bg-orange-900';
                    badgeText = 'text-orange-800 dark:text-orange-200';
                } else if (reference.startsWith('MANUAL_')) {
                    methodType = 'NEQUI';
                    badgeBg = 'bg-pink-100 dark:bg-pink-900';
                    badgeText = 'text-pink-800 dark:text-pink-200';
                } else if (reference.startsWith('MAN_DEB_')) {
                    methodType = 'DÉBITO';
                    badgeBg = 'bg-purple-100 dark:bg-purple-900';
                    badgeText = 'text-purple-800 dark:text-purple-200';
                }

                return (
                    <div className="flex items-center gap-2">
                        <span className={`px-2 py-1 rounded text-xs font-semibold whitespace-nowrap ${badgeBg} ${badgeText}`}>
                            {methodType}
                        </span>
                        <span className="text-gray-700 dark:text-gray-300 text-sm break-words">
                            {reference}
                        </span>
                    </div>
                );
            }
        },
        {
            key: 'Concept',
            label: 'Categoria',
            render: (val) => {
                const concept = val as string;
                const label = concept ? (CONCEPT_LABELS[concept] || concept) : '';
                return <span className="text-gray-700 dark:text-gray-300 text-sm">{label || '-'}</span>;
            }
        },
        {
            key: 'Amount',
            label: 'Monto',
            render: (val, row) => {
                const amount = val as number;
                const isIncome = row.Type === 'RECHARGE' || row.Type === 'income';
                const color = isIncome ? '#16a34a' : '#dc2626';
                const sign = isIncome ? '+' : '-';
                return (
                    <span className="font-bold font-mono" style={{ color }}>
                        {sign}${formatCurrency(Math.abs(amount)).replace('$', '')}
                    </span>
                );
            }
        },
    ];

    if (showActions) {
        requestColumns.push({
            key: 'actions',
            label: 'Acciones',
            render: (_, row) => (
                <div className="flex gap-2">
                    <button
                        onClick={() => handleAction(row.ID, 'approve')}
                        disabled={!!processingId}
                        className="px-3 py-1.5 text-xs font-semibold rounded-lg transition-all duration-200 text-white hover:opacity-90 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed"
                        style={{ backgroundColor: '#16a34a' }}
                    >
                        {processingId === row.ID ? '...' : 'Aprobar'}
                    </button>
                    <button
                        onClick={() => handleAction(row.ID, 'reject')}
                        disabled={!!processingId}
                        className="px-3 py-1.5 text-xs font-semibold rounded-lg transition-all duration-200 text-white hover:opacity-90 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed"
                        style={{ backgroundColor: '#dc2626' }}
                    >
                        {processingId === row.ID ? '...' : 'Rechazar'}
                    </button>
                </div>
            )
        });
    }

    const totalItems = requests.length;
    const totalPages = Math.ceil(totalItems / itemsPerPage);
    const paginatedData = requests.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage);

    return (
        <div className={`bg-white dark:bg-gray-800 rounded-2xl border border-gray-200 dark:border-gray-700 shadow-sm overflow-hidden flex flex-col ${compact ? 'p-2' : ''}`}>
            {!compact && !hideTitle && (
                <div className="px-6 py-5 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-700/50 space-y-4">
                    <div className="flex items-center justify-between gap-4">
                        {!showFiltersOnly && <h2 className="text-lg font-semibold text-gray-900 dark:text-white">{title}</h2>}
                    </div>
                    <div className="flex flex-wrap gap-3 items-end">
                        <div>
                            <label className="text-xs font-semibold text-gray-600 dark:text-gray-400 block mb-2">Fecha</label>
                            <div className="flex gap-2">
                                <input
                                    type="date"
                                    value={dateFrom}
                                    onChange={(e) => {
                                        setDateFrom(e.target.value);
                                        setCurrentPage(1);
                                    }}
                                    className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm w-36 focus:outline-none focus:ring-2 focus:ring-purple-500"
                                    placeholder="Desde"
                                />
                                <input
                                    type="date"
                                    value={dateTo}
                                    onChange={(e) => {
                                        setDateTo(e.target.value);
                                        setCurrentPage(1);
                                    }}
                                    className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm w-36 focus:outline-none focus:ring-2 focus:ring-purple-500"
                                    placeholder="Hasta"
                                />
                            </div>
                        </div>
                        <div>
                            <label className="text-xs font-semibold text-gray-600 dark:text-gray-400 block mb-2">Tipo</label>
                            <div className="flex gap-2">
                                <button
                                    onClick={() => {
                                        setTypeFilter('all');
                                        setCurrentPage(1);
                                    }}
                                    className={`px-4 py-2 text-xs font-semibold rounded-lg transition-all ${
                                        typeFilter === 'all'
                                            ? 'bg-blue-600 text-white'
                                            : 'bg-white dark:bg-gray-700 text-gray-900 dark:text-white border border-gray-300 dark:border-gray-600 hover:bg-gray-100 dark:hover:bg-gray-600'
                                    }`}
                                >
                                    Todas
                                </button>
                                <button
                                    onClick={() => {
                                        setTypeFilter('RECHARGE');
                                        setCurrentPage(1);
                                    }}
                                    className={`px-4 py-2 text-xs font-semibold rounded-lg transition-all ${
                                        typeFilter === 'RECHARGE'
                                            ? 'bg-green-600 text-white'
                                            : 'bg-white dark:bg-gray-700 text-gray-900 dark:text-white border border-gray-300 dark:border-gray-600 hover:bg-gray-100 dark:hover:bg-gray-600'
                                    }`}
                                >
                                    Ingresos
                                </button>
                                <button
                                    onClick={() => {
                                        setTypeFilter('USAGE');
                                        setCurrentPage(1);
                                    }}
                                    className={`px-4 py-2 text-xs font-semibold rounded-lg transition-all ${
                                        typeFilter === 'USAGE'
                                            ? 'bg-red-600 text-white'
                                            : 'bg-white dark:bg-gray-700 text-gray-900 dark:text-white border border-gray-300 dark:border-gray-600 hover:bg-gray-100 dark:hover:bg-gray-600'
                                    }`}
                                >
                                    Egresos
                                </button>
                            </div>
                        </div>
                        <div>
                            <label className="text-xs font-semibold text-gray-600 dark:text-gray-400 block mb-2">Negocio</label>
                            <select
                                value={businessFilter}
                                onChange={(e) => {
                                    setBusinessFilter(e.target.value === 'all' ? 'all' : parseInt(e.target.value));
                                    setCurrentPage(1);
                                }}
                                className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm w-48 focus:outline-none focus:ring-2 focus:ring-purple-500"
                            >
                                <option value="all">Todos los negocios</option>
                                {Object.entries(businesses).map(([id, name]) => (
                                    <option key={id} value={id}>{name}</option>
                                ))}
                            </select>
                        </div>
                        {(dateFrom || dateTo || typeFilter !== 'all' || businessFilter !== 'all') && (
                            <button
                                onClick={() => {
                                    setDateFrom('');
                                    setDateTo('');
                                    setTypeFilter('all');
                                    setBusinessFilter('all');
                                    setCurrentPage(1);
                                }}
                                className="px-3 py-2 text-xs font-medium text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white transition-colors"
                            >
                                Limpiar
                            </button>
                        )}
                    </div>
                </div>
            )}
            {compact && <div className="px-6 py-3 bg-gray-50 dark:bg-gray-700/50 border-b border-gray-200 dark:border-gray-700 font-bold text-gray-700 dark:text-gray-100 text-sm">{title}</div>}

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
            key: 'Concept',
            label: 'Categoria',
            render: (val) => {
                const concept = val as string;
                const label = concept ? (CONCEPT_LABELS[concept] || concept) : '';
                return <span className="text-gray-700 dark:text-gray-300 text-sm">{label || '-'}</span>;
            }
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
        <div className="bg-white dark:bg-gray-800 rounded-2xl border border-gray-200 dark:border-gray-700 overflow-hidden shadow-sm">
            <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-700/50">
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
