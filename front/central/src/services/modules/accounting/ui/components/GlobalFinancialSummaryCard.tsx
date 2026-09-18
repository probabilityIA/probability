'use client';

import { useEffect, useMemo, useState } from 'react';
import { PieChart, Pie, Cell, BarChart, Bar, XAxis, YAxis, CartesianGrid, Legend, Tooltip, ResponsiveContainer } from 'recharts';
import { AccountingReport } from '../../domain/types';
import { getAccountingReportAction } from '../../infra/actions';
import { formatCOP, firstDayOfCurrentMonth, lastDayOfCurrentMonth, kindLabel } from '../format';
import { Alert, Spinner } from '@/shared/ui';

const PIE_COLORS = ['#8B5CF6', '#10B981', '#3B82F6', '#F59E0B', '#EF4444', '#6366F1', '#EC4899'];

function PieCard({ title, data, emptyLabel, keyPrefix }: { title: string; data: { name: string; value: number }[]; emptyLabel: string; keyPrefix: string }) {
    return (
        <div className="border border-gray-200 dark:border-gray-700 rounded-xl p-4 flex flex-col items-center">
            <h3 className="text-sm font-semibold text-gray-900 dark:text-white mb-2 self-start">{title}</h3>
            {data.length === 0 ? (
                <p className="text-sm text-gray-500 dark:text-gray-400 py-12">{emptyLabel}</p>
            ) : (
                <ResponsiveContainer width="100%" height={300}>
                    <PieChart>
                        <Pie
                            data={data}
                            cx="50%"
                            cy="50%"
                            labelLine={false}
                            label={({ name, percent }) => `${name}: ${((percent ?? 0) * 100).toFixed(0)}%`}
                            outerRadius={100}
                            dataKey="value"
                        >
                            {data.map((_, index) => (
                                <Cell key={`${keyPrefix}-${index}`} fill={PIE_COLORS[index % PIE_COLORS.length]} />
                            ))}
                        </Pie>
                        <Tooltip formatter={(value) => formatCOP(value as number)} />
                        <Legend />
                    </PieChart>
                </ResponsiveContainer>
            )}
        </div>
    );
}

export function GlobalFinancialSummaryCard() {
    const [from, setFrom] = useState(firstDayOfCurrentMonth());
    const [to, setTo] = useState(lastDayOfCurrentMonth());
    const [report, setReport] = useState<AccountingReport | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        let active = true;
        setLoading(true);
        getAccountingReportAction(from, to).then((result) => {
            if (!active) return;
            if (result.success) {
                setReport(result.data);
                setError(null);
            } else {
                setReport(null);
                setError(result.error);
            }
            setLoading(false);
        });
        return () => {
            active = false;
        };
    }, [from, to]);

    const totals = report?.totals;
    const cards = [
        { label: 'Ganancia real', value: totals?.real_income ?? 0, accent: 'text-emerald-600 dark:text-emerald-400' },
        { label: 'Gastos', value: totals?.real_expense ?? 0, accent: 'text-red-600 dark:text-red-400' },
        { label: 'Neto', value: totals?.net ?? 0, accent: (totals?.net ?? 0) >= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-red-600 dark:text-red-400' },
    ];
    const byConcept = report?.by_concept || [];
    const realIncomeRows = byConcept.filter((row) => row.kind === 'INCOME' && row.is_real_income);
    const cashInRows = byConcept.filter((row) => row.kind === 'INCOME');

    const realIncomePieData = useMemo(
        () => realIncomeRows.map((row) => ({ name: row.name, value: row.amount })),
        [realIncomeRows]
    );
    const subscriptionByBusiness = report?.by_subscription_business || [];
    const guideMarginByBusiness = report?.by_guide_margin_business || [];
    const subscriptionByBusinessPieData = useMemo(
        () => subscriptionByBusiness.map((row) => ({ name: row.business_name, value: row.amount })),
        [subscriptionByBusiness]
    );
    const guideMarginByBusinessPieData = useMemo(
        () => guideMarginByBusiness.map((row) => ({ name: row.business_name, value: row.amount })),
        [guideMarginByBusiness]
    );
    const cashInPieData = useMemo(
        () => cashInRows.map((row) => ({ name: row.name, value: row.amount })),
        [cashInRows]
    );
    const totalsBarData = useMemo(
        () => [
            { name: 'Ganancia real', value: totals?.real_income ?? 0 },
            { name: 'Entradas de caja', value: totals?.cash_in ?? 0 },
            { name: 'Gastos', value: totals?.real_expense ?? 0 },
            { name: 'Salidas de caja', value: totals?.cash_out ?? 0 },
            { name: 'Impuestos', value: totals?.tax_total ?? 0 },
            { name: 'Neto', value: totals?.net ?? 0 },
        ],
        [totals]
    );

    return (
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm dark:shadow-lg border border-gray-200 dark:border-gray-700 p-6 space-y-6">
            <div className="flex flex-col sm:flex-row sm:items-end sm:justify-between gap-4">
                <div>
                    <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Resumen financiero de Probability</h2>
                    <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
                        Todos los negocios juntos, excluyendo negocios de prueba
                    </p>
                </div>
                <div className="flex flex-wrap items-end gap-2">
                    <div>
                        <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">Desde</label>
                        <input
                            type="date"
                            value={from}
                            onChange={(e) => setFrom(e.target.value)}
                            className="px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                        />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">Hasta</label>
                        <input
                            type="date"
                            value={to}
                            onChange={(e) => setTo(e.target.value)}
                            className="px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                        />
                    </div>
                </div>
            </div>

            {loading && (
                <div className="flex justify-center py-8">
                    <Spinner size="md" color="primary" />
                </div>
            )}

            {!loading && error && <Alert type="error">{error}</Alert>}

            {!loading && !error && (
                <>
                    <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                        {cards.map((card) => (
                            <div key={card.label} className="bg-gray-50 dark:bg-gray-900/40 border border-gray-200 dark:border-gray-700 rounded-xl p-4">
                                <p className="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide">{card.label}</p>
                                <p className={`mt-2 text-xl font-bold ${card.accent}`}>{formatCOP(card.value)}</p>
                            </div>
                        ))}
                    </div>

                    <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
                        <PieCard
                            title="Ganancia real por fuente"
                            data={realIncomePieData}
                            emptyLabel="Sin ganancia real en el periodo"
                            keyPrefix="real"
                        />
                        <PieCard
                            title="Todo lo que entra a caja"
                            data={cashInPieData}
                            emptyLabel="Sin movimientos en el periodo"
                            keyPrefix="cash"
                        />
                    </div>

                    <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
                        <PieCard
                            title="Membresias: que negocio las genera"
                            data={subscriptionByBusinessPieData}
                            emptyLabel="Sin pagos de membresia en el periodo"
                            keyPrefix="sub-biz"
                        />
                        <PieCard
                            title="Ganancia por guias: que negocio la genera"
                            data={guideMarginByBusinessPieData}
                            emptyLabel="Sin ganancia por guias en el periodo"
                            keyPrefix="guide-biz"
                        />
                    </div>

                    <div className="border border-gray-200 dark:border-gray-700 rounded-xl p-4">
                        <h3 className="text-sm font-semibold text-gray-900 dark:text-white mb-4">Comparativo del periodo</h3>
                        <ResponsiveContainer width="100%" height={280}>
                            <BarChart data={totalsBarData}>
                                <CartesianGrid strokeDasharray="3 3" className="stroke-gray-200 dark:stroke-gray-700" />
                                <XAxis dataKey="name" tick={{ fontSize: 12 }} />
                                <YAxis tickFormatter={(v) => formatCOP(v)} width={90} tick={{ fontSize: 11 }} />
                                <Tooltip formatter={(value) => formatCOP(value as number)} />
                                <Bar dataKey="value" fill="#8B5CF6" radius={[4, 4, 0, 0]} />
                            </BarChart>
                        </ResponsiveContainer>
                    </div>

                    <div className="border border-gray-200 dark:border-gray-700 rounded-xl overflow-hidden">
                        <div className="px-4 py-3 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/40">
                            <h3 className="text-sm font-semibold text-gray-900 dark:text-white">Detalle por concepto</h3>
                        </div>
                        <div className="overflow-x-auto">
                            <table className="w-full text-sm">
                                <thead>
                                    <tr className="text-left text-xs text-gray-500 dark:text-gray-400 uppercase bg-gray-50 dark:bg-gray-900/40">
                                        <th className="px-4 py-2 font-medium">Fuente</th>
                                        <th className="px-4 py-2 font-medium">Tipo</th>
                                        <th className="px-4 py-2 font-medium text-right">Movs.</th>
                                        <th className="px-4 py-2 font-medium text-right">Monto</th>
                                    </tr>
                                </thead>
                                <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                                    {byConcept.length === 0 && (
                                        <tr>
                                            <td colSpan={4} className="px-4 py-6 text-center text-gray-500 dark:text-gray-400">
                                                Sin movimientos en el periodo
                                            </td>
                                        </tr>
                                    )}
                                    {byConcept.map((row) => (
                                        <tr key={row.concept_id} className="text-gray-800 dark:text-gray-200">
                                            <td className="px-4 py-2">
                                                <span className="font-medium">{row.name}</span>
                                                {!row.is_real_income && (
                                                    <span className="ml-2 text-[10px] font-semibold text-blue-600 dark:text-blue-400 uppercase">solo caja</span>
                                                )}
                                            </td>
                                            <td className="px-4 py-2">
                                                <span className={`inline-block px-2 py-0.5 rounded-full text-xs font-semibold ${row.kind === 'INCOME' ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300' : 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300'}`}>
                                                    {kindLabel(row.kind)}
                                                </span>
                                            </td>
                                            <td className="px-4 py-2 text-right">{row.entries_count}</td>
                                            <td className="px-4 py-2 text-right font-medium">{formatCOP(row.amount)}</td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </>
            )}
        </div>
    );
}
