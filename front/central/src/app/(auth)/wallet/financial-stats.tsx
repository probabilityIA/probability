'use client';

import { useState, useEffect, useCallback } from 'react';
import { Button, Input, Select, Table, TableColumn, Spinner, Alert } from '@/shared/ui';
import { getFinancialStatsAction, FinancialStatsResponse, BusinessFinancialStats } from '@/services/modules/wallet/infra/actions';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';
import { PieChart, Pie, Cell, Legend, Tooltip, ResponsiveContainer, ComposedChart, Bar, XAxis, YAxis, CartesianGrid } from 'recharts';

const formatCurrency = (amount: number) =>
    new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP' }).format(amount);

const COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899', '#06b6d4', '#f97316'];

type ViewMode = 'month' | 'custom';

export function FinancialStatsView() {
    const { businesses } = useBusinessesSimple();
    const [viewMode, setViewMode] = useState<ViewMode>('month');
    const [selectedMonth, setSelectedMonth] = useState<string>(getCurrentMonth());
    const [startDate, setStartDate] = useState<string>(getFirstDayOfMonth());
    const [endDate, setEndDate] = useState<string>(getTodayDate());
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
    const [stats, setStats] = useState<FinancialStatsResponse | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const fetchStats = useCallback(async () => {
        try {
            setLoading(true);
            setError(null);

            const params = {
                businessId: selectedBusinessId || undefined,
                month: viewMode === 'month' ? selectedMonth : undefined,
                startDate: viewMode === 'custom' ? startDate : undefined,
                endDate: viewMode === 'custom' ? endDate : undefined,
            };

            const result = await getFinancialStatsAction(
                params.businessId,
                params.startDate,
                params.endDate,
                params.month
            );

            if (!result.success) {
                setError(result.error || 'Error al cargar estadísticas financieras');
                return;
            }

            if (result.data) {
                setStats(result.data);
            }
        } catch (err: any) {
            setError(err.message || 'Error desconocido');
        } finally {
            setLoading(false);
        }
    }, [viewMode, selectedMonth, startDate, endDate, selectedBusinessId]);

    useEffect(() => {
        fetchStats();
    }, [fetchStats]);

    const handleViewModeChange = (mode: ViewMode) => {
        setViewMode(mode);
    };

    // Prepare pie chart data
    const pieChartData = stats?.businesses.map(b => ({
        name: b.business_name,
        value: b.total_income,
    })) || [];

    // Prepare breakdown table (only when a business is selected)
    const breakdownData = selectedBusinessId && stats?.businesses.length === 1
        ? [{
            type: 'Membresías',
            amount: stats.businesses[0].subscription_income,
            percentage: stats.total_income > 0 ? (stats.businesses[0].subscription_income / stats.total_income * 100).toFixed(1) : '0',
        },
        {
            type: 'Guías',
            amount: stats.businesses[0].guide_income,
            count: stats.businesses[0].guide_count,
            percentage: stats.total_income > 0 ? (stats.businesses[0].guide_income / stats.total_income * 100).toFixed(1) : '0',
        }]
        : null;

    return (
        <div className="space-y-6">
            {/* Filters Section */}
            <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-4 space-y-4">
                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                        Negocio
                    </label>
                    <select
                        value={selectedBusinessId?.toString() ?? ''}
                        onChange={(e) => setSelectedBusinessId(e.target.value ? Number(e.target.value) : null)}
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-sm bg-white dark:bg-gray-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                    >
                        <option value="">Todos los negocios</option>
                        {businesses.map((b) => (
                            <option key={b.id} value={b.id}>{b.name}</option>
                        ))}
                    </select>
                </div>

                <div className="flex gap-4">
                    <label className="flex items-center gap-2">
                        <input
                            type="radio"
                            name="viewMode"
                            value="month"
                            checked={viewMode === 'month'}
                            onChange={() => handleViewModeChange('month')}
                            className="w-4 h-4"
                        />
                        <span className="text-sm font-medium text-gray-700 dark:text-gray-300">Mes</span>
                    </label>
                    <label className="flex items-center gap-2">
                        <input
                            type="radio"
                            name="viewMode"
                            value="custom"
                            checked={viewMode === 'custom'}
                            onChange={() => handleViewModeChange('custom')}
                            className="w-4 h-4"
                        />
                        <span className="text-sm font-medium text-gray-700 dark:text-gray-300">Rango personalizado</span>
                    </label>
                </div>

                {viewMode === 'month' ? (
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                            Mes
                        </label>
                        <input
                            type="month"
                            value={selectedMonth}
                            onChange={(e) => setSelectedMonth(e.target.value)}
                            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-sm bg-white dark:bg-gray-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                        />
                    </div>
                ) : (
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                                Fecha inicio
                            </label>
                            <input
                                type="date"
                                value={startDate}
                                onChange={(e) => setStartDate(e.target.value)}
                                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-sm bg-white dark:bg-gray-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                                Fecha fin
                            </label>
                            <input
                                type="date"
                                value={endDate}
                                onChange={(e) => setEndDate(e.target.value)}
                                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-sm bg-white dark:bg-gray-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                            />
                        </div>
                    </div>
                )}
            </div>

            {error && <Alert type="error">{error}</Alert>}

            {loading ? (
                <div className="flex justify-center py-8">
                    <Spinner />
                </div>
            ) : stats ? (
                <div className="space-y-6">
                    {/* Summary Card */}
                    <div className="bg-white dark:bg-gray-800 rounded-lg p-6 border border-gray-200 dark:border-gray-700">
                        <div className="grid grid-cols-3 gap-4">
                            <div>
                                <p className="text-sm text-gray-600 dark:text-gray-400 mb-2">Período</p>
                                <p className="text-lg font-semibold text-gray-900 dark:text-white">
                                    {new Date(stats.period.start).toLocaleDateString('es-CO')} - {new Date(stats.period.end).toLocaleDateString('es-CO')}
                                </p>
                            </div>
                            <div>
                                <p className="text-sm text-gray-600 dark:text-gray-400 mb-2">Total de negocios</p>
                                <p className="text-lg font-semibold text-gray-900 dark:text-white">{stats.businesses.length}</p>
                            </div>
                            <div>
                                <p className="text-sm text-gray-600 dark:text-gray-400 mb-2">Ingresos totales</p>
                                <p className="text-2xl font-bold text-green-600 dark:text-green-400">{formatCurrency(stats.total_income)}</p>
                            </div>
                        </div>
                    </div>

                    {/* Charts and Breakdown */}
                    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                        {/* Pie Chart */}
                        <div className="lg:col-span-2 bg-white dark:bg-gray-800 rounded-lg p-6 border border-gray-200 dark:border-gray-700">
                            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
                                {selectedBusinessId ? `${stats.businesses[0]?.business_name || 'Negocio'} - Desglose` : 'Ingresos por negocio'}
                            </h3>
                            {pieChartData.length > 0 ? (
                                <ResponsiveContainer width="100%" height={300}>
                                    <PieChart>
                                        <Pie
                                            data={pieChartData}
                                            cx="50%"
                                            cy="50%"
                                            labelLine={false}
                                            label={({ name, value, percent }) => `${name}: ${percent ? (percent * 100).toFixed(0) : 0}%`}
                                            outerRadius={80}
                                            fill="#8884d8"
                                            dataKey="value"
                                        >
                                            {pieChartData.map((entry, index) => (
                                                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                                            ))}
                                        </Pie>
                                        <Tooltip formatter={(value: any) => formatCurrency(value as number)} />
                                        <Legend />
                                    </PieChart>
                                </ResponsiveContainer>
                            ) : (
                                <p className="text-center text-gray-500 dark:text-gray-400">Sin datos disponibles</p>
                            )}
                        </div>

                        {/* Breakdown (shown when business is selected) */}
                        {breakdownData ? (
                            <div className="bg-white dark:bg-gray-800 rounded-lg p-6 border border-gray-200 dark:border-gray-700">
                                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Desglose</h3>
                                <div className="space-y-4">
                                    {breakdownData.map((item, idx) => (
                                        <div key={idx} className="bg-gray-50 dark:bg-gray-700 rounded p-4">
                                            <div className="flex justify-between items-start mb-2">
                                                <span className="font-medium text-gray-900 dark:text-white">{item.type}</span>
                                                <span className="text-sm font-semibold text-gray-500 dark:text-gray-400">{item.percentage}%</span>
                                            </div>
                                            <p className="text-2xl font-bold text-gray-900 dark:text-white mb-1">
                                                {formatCurrency(item.amount)}
                                            </p>
                                            {item.count !== undefined && (
                                                <p className="text-sm text-gray-600 dark:text-gray-400">
                                                    {item.count.toLocaleString('es-CO')} guía{item.count !== 1 ? 's' : ''}
                                                </p>
                                            )}
                                        </div>
                                    ))}
                                </div>
                            </div>
                        ) : null}
                    </div>

                    {/* Detailed Table */}
                    <div className="bg-white dark:bg-gray-800 rounded-lg overflow-hidden border border-gray-200 dark:border-gray-700">
                        <div className="p-6 border-b border-gray-200 dark:border-gray-700">
                            <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                                {selectedBusinessId ? 'Detalles' : 'Por negocio'}
                            </h3>
                        </div>
                        <div className="overflow-x-auto">
                            <table className="w-full">
                                <thead className="bg-gray-50 dark:bg-gray-700">
                                    <tr>
                                        <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 dark:text-gray-300 uppercase tracking-wider">Negocio</th>
                                        <th className="px-6 py-3 text-right text-xs font-semibold text-gray-700 dark:text-gray-300 uppercase tracking-wider">Membresías</th>
                                        <th className="px-6 py-3 text-right text-xs font-semibold text-gray-700 dark:text-gray-300 uppercase tracking-wider">Guías</th>
                                        <th className="px-6 py-3 text-right text-xs font-semibold text-gray-700 dark:text-gray-300 uppercase tracking-wider">Cantidad</th>
                                        <th className="px-6 py-3 text-right text-xs font-semibold text-gray-700 dark:text-gray-300 uppercase tracking-wider">Total</th>
                                    </tr>
                                </thead>
                                <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                                    {stats.businesses.map((business, idx) => (
                                        <tr key={idx} className="hover:bg-gray-50 dark:hover:bg-gray-700">
                                            <td className="px-6 py-4 text-sm font-medium text-gray-900 dark:text-white">
                                                {business.business_name}
                                            </td>
                                            <td className="px-6 py-4 text-sm text-right text-gray-900 dark:text-white">
                                                {formatCurrency(business.subscription_income)}
                                            </td>
                                            <td className="px-6 py-4 text-sm text-right text-gray-900 dark:text-white">
                                                {formatCurrency(business.guide_income)}
                                            </td>
                                            <td className="px-6 py-4 text-sm text-right text-gray-900 dark:text-white">
                                                {business.guide_count.toLocaleString('es-CO')}
                                            </td>
                                            <td className="px-6 py-4 text-sm font-semibold text-right text-green-600 dark:text-green-400">
                                                {formatCurrency(business.total_income)}
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                                <tfoot className="bg-gray-50 dark:bg-gray-700 font-semibold">
                                    <tr>
                                        <td className="px-6 py-4 text-sm text-gray-900 dark:text-white">TOTAL</td>
                                        <td className="px-6 py-4 text-sm text-right text-gray-900 dark:text-white">
                                            {formatCurrency(stats.businesses.reduce((sum, b) => sum + b.subscription_income, 0))}
                                        </td>
                                        <td className="px-6 py-4 text-sm text-right text-gray-900 dark:text-white">
                                            {formatCurrency(stats.businesses.reduce((sum, b) => sum + b.guide_income, 0))}
                                        </td>
                                        <td className="px-6 py-4 text-sm text-right text-gray-900 dark:text-white">
                                            {stats.businesses.reduce((sum, b) => sum + b.guide_count, 0).toLocaleString('es-CO')}
                                        </td>
                                        <td className="px-6 py-4 text-sm text-right text-green-600 dark:text-green-400">
                                            {formatCurrency(stats.total_income)}
                                        </td>
                                    </tr>
                                </tfoot>
                            </table>
                        </div>
                    </div>
                </div>
            ) : null}
        </div>
    );
}

// Helper functions
function getCurrentMonth(): string {
    const now = new Date();
    return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`;
}

function getFirstDayOfMonth(): string {
    const now = new Date();
    return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-01`;
}

function getTodayDate(): string {
    const now = new Date();
    return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`;
}
