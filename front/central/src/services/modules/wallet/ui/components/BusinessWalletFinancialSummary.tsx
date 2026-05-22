'use client';

import { useEffect, useState } from 'react';
import { Spinner, Alert } from '@/shared/ui';
import { getWalletHistoryAction, getFinancialStatsAction } from '../../infra/actions';

interface Props {
    selectedBusinessId: number;
}

interface Transaction {
    ID: string;
    Amount: number;
    Description?: string;
    Type: string;
    CreatedAt: string;
}

const fmt = (n: number) => '$ ' + Math.round(n).toLocaleString('es-CO');

export default function BusinessWalletFinancialSummary({ selectedBusinessId }: Props) {
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [transactions, setTransactions] = useState<Transaction[]>([]);
    const [stats, setStats] = useState<any>(null);
    const [totals, setTotals] = useState({
        recharged: 0,
        membershipExpense: 0,
        guideExpense: 0
    });

    useEffect(() => {
        const loadData = async () => {
            try {
                setLoading(true);
                setError(null);

                const historyRes = await getWalletHistoryAction(selectedBusinessId);

                if (!historyRes.success) {
                    console.error('History error:', historyRes.error);
                    throw new Error(historyRes.error || 'Error al cargar historial');
                }

                console.log('History data:', historyRes.data);

                const txData = (historyRes.data || []) as Transaction[];
                setTransactions(txData.sort((a, b) => new Date(b.CreatedAt).getTime() - new Date(a.CreatedAt).getTime()));

                const recharged = txData
                    .filter(tx => tx.Type?.toLowerCase().includes('recharge') || tx.Type?.toLowerCase().includes('recarga'))
                    .reduce((sum, tx) => sum + Math.abs(tx.Amount), 0);

                const expenses = txData
                    .filter(tx => tx.Amount < 0)
                    .reduce((sum, tx) => sum + Math.abs(tx.Amount), 0);

                setTotals({
                    recharged,
                    membershipExpense: expenses,
                    guideExpense: 0
                });
            } catch (err: any) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        loadData();
    }, [selectedBusinessId]);

    if (loading) return <Spinner />;
    if (error) return <Alert type="error">{error}</Alert>;

    return (
        <div className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="bg-white dark:bg-gray-800 p-4 rounded-lg border border-gray-200 dark:border-gray-700">
                    <p className="text-xs text-gray-500 dark:text-gray-400 mb-1">Total Recargado</p>
                    <p className="text-2xl font-bold text-blue-600">{fmt(totals.recharged)}</p>
                </div>
                <div className="bg-white dark:bg-gray-800 p-4 rounded-lg border border-gray-200 dark:border-gray-700">
                    <p className="text-xs text-gray-500 dark:text-gray-400 mb-1">Gastos de Membresía</p>
                    <p className="text-2xl font-bold text-red-600">{fmt(totals.membershipExpense)}</p>
                </div>
            </div>

            <div className="bg-white dark:bg-gray-800 p-6 rounded-lg border border-gray-200 dark:border-gray-700">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Historial de Transacciones</h3>
                <div className="overflow-x-auto">
                    <table className="w-full text-sm">
                        <thead>
                            <tr className="border-b border-gray-200 dark:border-gray-700">
                                <th className="text-left py-3 px-2 font-semibold text-gray-700 dark:text-gray-300">Fecha</th>
                                <th className="text-left py-3 px-2 font-semibold text-gray-700 dark:text-gray-300">Tipo</th>
                                <th className="text-left py-3 px-2 font-semibold text-gray-700 dark:text-gray-300">Descripción</th>
                                <th className="text-right py-3 px-2 font-semibold text-gray-700 dark:text-gray-300">Monto</th>
                            </tr>
                        </thead>
                        <tbody>
                            {transactions.length > 0 ? (
                                transactions.slice(0, 20).map(tx => (
                                    <tr key={tx.ID} className="border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50">
                                        <td className="py-3 px-2 text-gray-900 dark:text-gray-100 text-xs">
                                            {new Date(tx.CreatedAt).toLocaleDateString('es-CO')}
                                        </td>
                                        <td className="py-3 px-2 text-gray-600 dark:text-gray-400">{tx.Type}</td>
                                        <td className="py-3 px-2 text-gray-600 dark:text-gray-400 text-xs">{tx.Description || '-'}</td>
                                        <td className={`py-3 px-2 text-right font-medium ${tx.Amount > 0 ? 'text-blue-600' : 'text-red-600'}`}>
                                            {tx.Amount > 0 ? '+' : ''}{fmt(tx.Amount)}
                                        </td>
                                    </tr>
                                ))
                            ) : (
                                <tr>
                                    <td colSpan={4} className="py-4 px-2 text-center text-gray-600 dark:text-gray-400">
                                        No hay transacciones registradas
                                    </td>
                                </tr>
                            )}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
}
