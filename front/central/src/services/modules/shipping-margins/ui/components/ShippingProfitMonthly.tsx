'use client';

import { useEffect, useState, useCallback } from 'react';
import { Spinner, Alert } from '@/shared/ui';
import { ProfitReportDetailResponse } from '../../domain/types';
import { shippingProfitReportDetailAction } from '../../infra/actions';
import { PieChart, Pie, Cell, Legend, Tooltip, ResponsiveContainer } from 'recharts';
import { getActionError } from '@/shared/utils/action-result';

interface Props {
    selectedBusinessId?: number;
}

interface MonthlyData {
    month: string;
    monthKey: string;
    shipments: number;
    customer_charge: number;
    carrier_cost: number;
    profit: number;
}

interface CarrierMonthlyData {
    carrier: string;
    shipments: number;
    customer_charge: number;
    carrier_cost: number;
    profit: number;
}

interface DailyShipmentData {
    shipment_id: number;
    order_number: string;
    tracking_number: string;
    carrier: string;
    customer_charge: number;
    carrier_cost: number;
    profit: number;
    created_at: string;
}

const fmt = (n: number) => '$ ' + Math.round(n).toLocaleString('es-CO');

function getMonthsRange(monthsBack: number = 12) {
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    const to = today;
    const from = new Date(today.getFullYear(), today.getMonth() - (monthsBack - 1), 1);
    from.setHours(0, 0, 0, 0);

    return {
        from: `${from.getFullYear()}-${String(from.getMonth() + 1).padStart(2, '0')}-${String(from.getDate()).padStart(2, '0')}`,
        to: `${to.getFullYear()}-${String(to.getMonth() + 1).padStart(2, '0')}-${String(to.getDate()).padStart(2, '0')}`
    };
}

function getMonthKey(dateStr: string): string {
    const date = new Date(dateStr);
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`;
}

function getMonthLabel(monthKey: string): string {
    const [year, month] = monthKey.split('-');
    const date = new Date(parseInt(year), parseInt(month) - 1);
    return date.toLocaleString('es-CO', { month: 'long', year: 'numeric' });
}

export default function ShippingProfitMonthly({ selectedBusinessId }: Props) {
    const [monthlyData, setMonthlyData] = useState<MonthlyData[]>([]);
    const [carriersByMonth, setCarriersByMonth] = useState<Map<string, CarrierMonthlyData[]>>(new Map());
    const [dailyShipments, setDailyShipments] = useState<DailyShipmentData[]>([]);
    const [expandedMonth, setExpandedMonth] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [totals, setTotals] = useState({ shipments: 0, customer_charge: 0, carrier_cost: 0, profit: 0 });
    const [filterFromDate, setFilterFromDate] = useState<string>('');
    const [filterToDate, setFilterToDate] = useState<string>('');
    const [detailPage, setDetailPage] = useState(1);
    const pageSize = 20;

    const loadData = useCallback(async () => {
        try {
            setLoading(true);
            setError(null);

            const range = getMonthsRange(60);
            const result = await shippingProfitReportDetailAction({
                business_id: selectedBusinessId,
                from: range.from,
                to: range.to,
                page: 1,
                page_size: 1000
            });

            if (result instanceof Error) {
                throw result;
            }

            const data = result as ProfitReportDetailResponse;
            const monthMap = new Map<string, MonthlyData>();
            const carrierMonthMap = new Map<string, Map<string, CarrierMonthlyData>>();

            data.data.forEach(row => {
                const monthKey = getMonthKey(row.created_at);

                if (!monthMap.has(monthKey)) {
                    monthMap.set(monthKey, {
                        month: getMonthLabel(monthKey),
                        monthKey,
                        shipments: 0,
                        customer_charge: 0,
                        carrier_cost: 0,
                        profit: 0
                    });
                }
                const month = monthMap.get(monthKey)!;
                month.shipments += 1;
                month.customer_charge += row.customer_charge;
                month.carrier_cost += row.carrier_cost;
                month.profit += row.profit;

                if (!carrierMonthMap.has(monthKey)) {
                    carrierMonthMap.set(monthKey, new Map());
                }
                const carrierMap = carrierMonthMap.get(monthKey)!;
                const carrierKey = row.carrier;

                if (!carrierMap.has(carrierKey)) {
                    carrierMap.set(carrierKey, {
                        carrier: row.carrier,
                        shipments: 0,
                        customer_charge: 0,
                        carrier_cost: 0,
                        profit: 0
                    });
                }

                const carrier = carrierMap.get(carrierKey)!;
                carrier.shipments += 1;
                carrier.customer_charge += row.customer_charge;
                carrier.carrier_cost += row.carrier_cost;
                carrier.profit += row.profit;
            });

            const sorted = Array.from(monthMap.values())
                .sort((a, b) => a.monthKey.localeCompare(b.monthKey));

            setMonthlyData(sorted);

            const carriersByMonthFinal = new Map<string, CarrierMonthlyData[]>();
            carrierMonthMap.forEach((carrierMap, monthKey) => {
                const carriersList = Array.from(carrierMap.values())
                    .sort((a, b) => b.shipments - a.shipments);
                carriersByMonthFinal.set(monthKey, carriersList);
            });
            setCarriersByMonth(carriersByMonthFinal);

            const dailyShipmentsData: DailyShipmentData[] = data.data.map(row => ({
                shipment_id: row.shipment_id,
                order_number: row.order_number,
                tracking_number: row.tracking_number,
                carrier: row.carrier,
                customer_charge: row.customer_charge,
                carrier_cost: row.carrier_cost,
                profit: row.profit,
                created_at: row.created_at
            })).sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());

            setDailyShipments(dailyShipmentsData);

            const totalShipments = data.data.length;
            const totalCustomerCharge = data.data.reduce((sum, row) => sum + row.customer_charge, 0);
            const totalCarrierCost = data.data.reduce((sum, row) => sum + row.carrier_cost, 0);
            const totalProfit = totalCustomerCharge - totalCarrierCost;

            setTotals({
                shipments: totalShipments,
                customer_charge: totalCustomerCharge,
                carrier_cost: totalCarrierCost,
                profit: totalProfit
            });

            setFilterFromDate(range.from);
            setFilterToDate(range.to);
        } catch (err: any) {
            setError(getActionError(err));
        } finally {
            setLoading(false);
        }
    }, [selectedBusinessId]);

    useEffect(() => {
        loadData();
    }, [loadData]);

    if (loading) return <Spinner />;
    if (error) return <Alert type="error">{error}</Alert>;

    const pieData = [
        { name: 'Cobrado Cliente', value: totals.customer_charge, fill: '#3b82f6' },
        { name: 'Costo Carrier', value: totals.carrier_cost, fill: '#ef4444' }
    ];

    return (
        <div className="space-y-8">
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                <div className="bg-white dark:bg-gray-800 p-4 rounded-lg border border-gray-200 dark:border-gray-700">
                    <p className="text-xs text-gray-500 dark:text-gray-400 mb-1">Guías Generadas</p>
                    <p className="text-2xl font-bold text-gray-900 dark:text-white">{totals.shipments}</p>
                </div>
                <div className="bg-white dark:bg-gray-800 p-4 rounded-lg border border-gray-200 dark:border-gray-700">
                    <p className="text-xs text-gray-500 dark:text-gray-400 mb-1">Cobrado Cliente</p>
                    <p className="text-2xl font-bold text-blue-600">{fmt(totals.customer_charge)}</p>
                </div>
                <div className="bg-white dark:bg-gray-800 p-4 rounded-lg border border-gray-200 dark:border-gray-700">
                    <p className="text-xs text-gray-500 dark:text-gray-400 mb-1">Costo Carrier</p>
                    <p className="text-2xl font-bold text-red-600">{fmt(totals.carrier_cost)}</p>
                </div>
                <div className="bg-white dark:bg-gray-800 p-4 rounded-lg border border-gray-200 dark:border-gray-700">
                    <p className="text-xs text-gray-500 dark:text-gray-400 mb-1">Ganancia</p>
                    <p className="text-2xl font-bold text-green-600">{fmt(totals.profit)}</p>
                </div>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <div className="bg-white dark:bg-gray-800 p-6 rounded-lg border border-gray-200 dark:border-gray-700 overflow-x-auto">
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Resumen Mensual</h3>
                    <table className="w-full text-sm">
                        <thead>
                            <tr className="border-b border-gray-200 dark:border-gray-700">
                                <th className="text-left py-2 px-2 font-semibold text-gray-700 dark:text-gray-300">Mes</th>
                                <th className="text-right py-2 px-2 font-semibold text-gray-700 dark:text-gray-300">Guías</th>
                                <th className="text-right py-2 px-2 font-semibold text-gray-700 dark:text-gray-300">Cobrado</th>
                                <th className="text-right py-2 px-2 font-semibold text-gray-700 dark:text-gray-300">Costo</th>
                                <th className="text-right py-2 px-2 font-semibold text-gray-700 dark:text-gray-300">Ganancia</th>
                            </tr>
                        </thead>
                        <tbody>
                            {monthlyData.map(row => (
                                <tr key={row.monthKey} className="border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50">
                                    <td className="py-3 px-2 text-gray-900 dark:text-gray-100">{row.month}</td>
                                    <td className="py-3 px-2 text-right text-gray-600 dark:text-gray-300">{row.shipments}</td>
                                    <td className="py-3 px-2 text-right text-blue-600 font-medium">{fmt(row.customer_charge)}</td>
                                    <td className="py-3 px-2 text-right text-red-600 font-medium">{fmt(row.carrier_cost)}</td>
                                    <td className="py-3 px-2 text-right text-green-600 font-medium">{fmt(row.profit)}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>

                <div className="bg-white dark:bg-gray-800 p-6 rounded-lg border border-gray-200 dark:border-gray-700 flex flex-col items-center justify-center">
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4 w-full">Distribución</h3>
                    <ResponsiveContainer width="100%" height={400}>
                        <PieChart>
                            <Pie
                                data={pieData}
                                cx="50%"
                                cy="50%"
                                labelLine={false}
                                label={({ name, value }) => `${name}: ${fmt(value)}`}
                                outerRadius={120}
                                fill="#8884d8"
                                dataKey="value"
                            >
                                {pieData.map((entry, index) => (
                                    <Cell key={`cell-${index}`} fill={entry.fill} />
                                ))}
                            </Pie>
                            <Tooltip formatter={(value) => fmt(value as number)} />
                            <Legend />
                        </PieChart>
                    </ResponsiveContainer>
                </div>
            </div>

            <div className="bg-white dark:bg-gray-800 p-6 rounded-lg border border-gray-200 dark:border-gray-700">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-6">Detalle por Día</h3>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                            Desde
                        </label>
                        <input
                            type="date"
                            value={filterFromDate}
                            onChange={(e) => setFilterFromDate(e.target.value)}
                            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                            Hasta
                        </label>
                        <input
                            type="date"
                            value={filterToDate}
                            onChange={(e) => setFilterToDate(e.target.value)}
                            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        />
                    </div>
                </div>

                <div className="overflow-x-auto">
                    <table className="w-full text-sm">
                        <thead>
                            <tr className="border-b border-gray-200 dark:border-gray-700">
                                <th className="text-left py-3 px-2 font-semibold text-gray-700 dark:text-gray-300">Fecha</th>
                                <th className="text-left py-3 px-2 font-semibold text-gray-700 dark:text-gray-300">Orden</th>
                                <th className="text-left py-3 px-2 font-semibold text-gray-700 dark:text-gray-300">Tracking</th>
                                <th className="text-left py-3 px-2 font-semibold text-gray-700 dark:text-gray-300">Transportadora</th>
                                <th className="text-right py-3 px-2 font-semibold text-gray-700 dark:text-gray-300">Cobrado</th>
                                <th className="text-right py-3 px-2 font-semibold text-gray-700 dark:text-gray-300">Costo Real</th>
                                <th className="text-right py-3 px-2 font-semibold text-gray-700 dark:text-gray-300">Ganancia</th>
                            </tr>
                        </thead>
                        <tbody>
                            {(() => {
                                const filtered = dailyShipments.filter(shipment => {
                                    const shipmentDate = shipment.created_at.split('T')[0];
                                    return shipmentDate >= filterFromDate && shipmentDate <= filterToDate;
                                });
                                const totalPages = Math.ceil(filtered.length / pageSize);
                                const start = (detailPage - 1) * pageSize;
                                const paginated = filtered.slice(start, start + pageSize);
                                return paginated.map(shipment => (
                                    <tr key={shipment.shipment_id} className="border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50">
                                        <td className="py-3 px-2 text-gray-900 dark:text-gray-100 text-xs">
                                            {new Date(shipment.created_at).toLocaleDateString('es-CO')}
                                        </td>
                                        <td className="py-3 px-2 text-gray-900 dark:text-gray-100 font-medium">{shipment.order_number}</td>
                                        <td className="py-3 px-2 text-gray-600 dark:text-gray-400 text-xs">{shipment.tracking_number}</td>
                                        <td className="py-3 px-2 text-gray-600 dark:text-gray-400">{shipment.carrier}</td>
                                        <td className="py-3 px-2 text-right text-blue-600 font-medium">{fmt(shipment.customer_charge)}</td>
                                        <td className="py-3 px-2 text-right text-red-600 font-medium">{fmt(shipment.carrier_cost)}</td>
                                        <td className="py-3 px-2 text-right text-green-600 font-medium">{fmt(shipment.profit)}</td>
                                    </tr>
                                ));
                            })()}
                        </tbody>
                    </table>
                </div>

                {(() => {
                    const filtered = dailyShipments.filter(shipment => {
                        const shipmentDate = shipment.created_at.split('T')[0];
                        return shipmentDate >= filterFromDate && shipmentDate <= filterToDate;
                    });
                    const totalPages = Math.ceil(filtered.length / pageSize);
                    return (
                        <div className="mt-4 flex items-center justify-between">
                            <div className="text-gray-600 dark:text-gray-300">
                                Página <span className="font-semibold">{detailPage}</span> de <span className="font-semibold">{totalPages}</span> &middot; {filtered.length} guías
                            </div>
                            <div className="flex items-center gap-2">
                                <button
                                    onClick={() => setDetailPage((p) => Math.max(1, p - 1))}
                                    disabled={detailPage <= 1}
                                    className="px-3 py-1.5 border border-gray-300 dark:border-gray-600 rounded-md disabled:opacity-50 text-gray-700 dark:text-gray-200"
                                >
                                    Anterior
                                </button>
                                <button
                                    onClick={() => setDetailPage((p) => Math.min(totalPages, p + 1))}
                                    disabled={detailPage >= totalPages}
                                    className="px-3 py-1.5 border border-gray-300 dark:border-gray-600 rounded-md disabled:opacity-50 text-gray-700 dark:text-gray-200"
                                >
                                    Siguiente
                                </button>
                            </div>
                        </div>
                    );
                })()}
            </div>

            <div className="bg-white dark:bg-gray-800 p-6 rounded-lg border border-gray-200 dark:border-gray-700">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Resumen por Transportadora</h3>
                <div className="space-y-3">
                    {monthlyData.map(month => (
                        <div key={month.monthKey} className="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
                            <button
                                onClick={() => setExpandedMonth(expandedMonth === month.monthKey ? null : month.monthKey)}
                                className="w-full p-4 bg-gray-50 dark:bg-gray-700/50 hover:bg-gray-100 dark:hover:bg-gray-700 flex items-center justify-between transition-colors"
                            >
                                <span className="font-medium text-gray-900 dark:text-white">{month.month}</span>
                                <span className="text-sm text-gray-600 dark:text-gray-400">{month.shipments} guías</span>
                                <svg
                                    className={`w-5 h-5 transition-transform ${expandedMonth === month.monthKey ? 'rotate-180' : ''}`}
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                >
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 14l-7 7m0 0l-7-7m7 7V3" />
                                </svg>
                            </button>

                            {expandedMonth === month.monthKey && (
                                <div className="p-4 overflow-x-auto">
                                    <table className="w-full text-sm">
                                        <thead>
                                            <tr className="border-b border-gray-200 dark:border-gray-700">
                                                <th className="text-left py-2 px-2 font-semibold text-gray-700 dark:text-gray-300">Transportadora</th>
                                                <th className="text-right py-2 px-2 font-semibold text-gray-700 dark:text-gray-300">Guías</th>
                                                <th className="text-right py-2 px-2 font-semibold text-gray-700 dark:text-gray-300">Cobrado</th>
                                                <th className="text-right py-2 px-2 font-semibold text-gray-700 dark:text-gray-300">Costo Real</th>
                                                <th className="text-right py-2 px-2 font-semibold text-gray-700 dark:text-gray-300">Ganancia</th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {(carriersByMonth.get(month.monthKey) || []).map(carrier => (
                                                <tr key={carrier.carrier} className="border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50">
                                                    <td className="py-3 px-2 text-gray-900 dark:text-gray-100">{carrier.carrier}</td>
                                                    <td className="py-3 px-2 text-right text-gray-600 dark:text-gray-300">{carrier.shipments}</td>
                                                    <td className="py-3 px-2 text-right text-blue-600 font-medium">{fmt(carrier.customer_charge)}</td>
                                                    <td className="py-3 px-2 text-right text-red-600 font-medium">{fmt(carrier.carrier_cost)}</td>
                                                    <td className="py-3 px-2 text-right text-green-600 font-medium">{fmt(carrier.profit)}</td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                </div>
                            )}
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}
