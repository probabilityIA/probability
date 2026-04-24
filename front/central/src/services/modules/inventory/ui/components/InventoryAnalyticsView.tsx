'use client';

import { useState, useEffect } from 'react';
import {
    PieChart,
    Pie,
    Cell,
    BarChart,
    Bar,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    Legend,
    ResponsiveContainer,
} from 'recharts';
import {
    BuildingStorefrontIcon,
    CubeIcon,
    ArchiveBoxIcon,
    CheckCircleIcon,
} from '@heroicons/react/24/outline';
import { Spinner } from '@/shared/ui';
import { getWarehousesAction } from '@/services/modules/warehouses/infra/actions';
import { getWarehouseInventoryAction } from '../../infra/actions';
import { getProductsAction } from '@/services/modules/products/infra/actions';
import { InventoryLevel } from '../../domain/types';
import { Warehouse } from '@/services/modules/warehouses/domain/types';

interface Props {
    businessId?: number;
}

interface WarehouseStat {
    name: string;
    code: string;
    units: number;
    products: number;
}

interface AnalyticsData {
    activeWarehouses: number;
    totalProducts: number;
    totalUnits: number;
    productsWithStock: number;
    warehouseStats: WarehouseStat[];
}

const DONUT_COLORS = ['#10b981', '#e5e7eb'];
const BAR_COLOR = '#6366f1';

export default function InventoryAnalyticsView({ businessId }: Props) {
    const [data, setData] = useState<AnalyticsData | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const load = async () => {
            setLoading(true);
            setError(null);
            try {
                const [warehousesRes, productsRes] = await Promise.all([
                    getWarehousesAction({ page: 1, page_size: 100, is_active: true, business_id: businessId }),
                    getProductsAction({ page: 1, page_size: 1, business_id: businessId } as any),
                ]);

                const warehouses: Warehouse[] = warehousesRes?.data ?? [];
                const totalProducts: number = (productsRes as any)?.total ?? 0;

                const inventoryResults = await Promise.all(
                    warehouses.map((w) =>
                        getWarehouseInventoryAction(w.id, { page: 1, page_size: 200, business_id: businessId })
                    )
                );

                let totalUnits = 0;
                const productIdsWithStock = new Set<string>();
                const warehouseStats: WarehouseStat[] = [];

                warehouses.forEach((warehouse, idx) => {
                    const levels: InventoryLevel[] = inventoryResults[idx]?.data ?? [];
                    let warehouseUnits = 0;
                    let warehouseProducts = 0;

                    levels.forEach((level) => {
                        totalUnits += level.quantity;
                        warehouseUnits += level.quantity;
                        if (level.quantity > 0) {
                            productIdsWithStock.add(level.product_id);
                            warehouseProducts++;
                        }
                    });

                    warehouseStats.push({
                        name: warehouse.name,
                        code: warehouse.code,
                        units: warehouseUnits,
                        products: warehouseProducts,
                    });
                });

                warehouseStats.sort((a, b) => b.units - a.units);

                setData({
                    activeWarehouses: warehouses.length,
                    totalProducts,
                    totalUnits,
                    productsWithStock: productIdsWithStock.size,
                    warehouseStats,
                });
            } catch (err: any) {
                setError(err?.message ?? 'Error al cargar datos');
            } finally {
                setLoading(false);
            }
        };

        load();
    }, [businessId]);

    if (loading) {
        return (
            <div className="flex items-center justify-center py-20 w-full">
                <Spinner size="lg" />
            </div>
        );
    }

    if (error || !data) {
        return (
            <div className="flex items-center justify-center py-20 w-full">
                <p className="text-gray-500 text-sm">{error ?? 'Sin datos disponibles'}</p>
            </div>
        );
    }

    const stockPercentage =
        data.totalProducts > 0 ? Math.round((data.productsWithStock / data.totalProducts) * 100) : 0;

    const donutData = [
        { name: 'Con stock', value: data.productsWithStock },
        { name: 'Sin stock', value: Math.max(0, data.totalProducts - data.productsWithStock) },
    ];

    const kpiCards = [
        {
            label: 'Bodegas activas',
            value: data.activeWarehouses,
            icon: BuildingStorefrontIcon,
            iconBg: 'bg-teal-100',
            iconColor: 'text-teal-600',
        },
        {
            label: 'Total productos',
            value: data.totalProducts,
            icon: CubeIcon,
            iconBg: 'bg-indigo-100',
            iconColor: 'text-indigo-600',
        },
        {
            label: 'Unidades en stock',
            value: data.totalUnits.toLocaleString(),
            icon: ArchiveBoxIcon,
            iconBg: 'bg-amber-100',
            iconColor: 'text-amber-600',
        },
        {
            label: 'Productos con stock',
            value: data.productsWithStock,
            icon: CheckCircleIcon,
            iconBg: 'bg-emerald-100',
            iconColor: 'text-emerald-600',
        },
    ];

    return (
        <div className="space-y-6">
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
                {kpiCards.map((card) => {
                    const Icon = card.icon;
                    return (
                        <div
                            key={card.label}
                            className="bg-white rounded-xl border border-gray-200 p-4 flex flex-col items-center gap-3"
                        >
                            <div className={`w-11 h-11 rounded-full flex items-center justify-center ${card.iconBg}`}>
                                <Icon className={`w-6 h-6 ${card.iconColor}`} />
                            </div>
                            <p className="text-2xl font-bold text-gray-900">{card.value}</p>
                            <p className="text-xs text-gray-500 text-center">{card.label}</p>
                        </div>
                    );
                })}
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <div className="bg-white rounded-xl border border-gray-200 p-5">
                    <h3 className="text-sm font-semibold text-gray-700 mb-1">Productos con / sin stock</h3>
                    <p className="text-xs text-gray-400 mb-4">{stockPercentage}% con stock disponible</p>
                    <ResponsiveContainer width="100%" height={220}>
                        <PieChart>
                            <Pie
                                data={donutData}
                                cx="50%"
                                cy="50%"
                                innerRadius={60}
                                outerRadius={90}
                                paddingAngle={2}
                                dataKey="value"
                            >
                                {donutData.map((_, index) => (
                                    <Cell key={index} fill={DONUT_COLORS[index]} />
                                ))}
                            </Pie>
                            <Tooltip formatter={(value) => [value, 'Productos']} />
                            <Legend />
                        </PieChart>
                    </ResponsiveContainer>
                </div>

                <div className="bg-white rounded-xl border border-gray-200 p-5">
                    <h3 className="text-sm font-semibold text-gray-700 mb-4">Unidades por bodega</h3>
                    <ResponsiveContainer width="100%" height={220}>
                        <BarChart data={data.warehouseStats.slice(0, 8)} margin={{ top: 4, right: 8, left: 0, bottom: 0 }}>
                            <CartesianGrid strokeDasharray="3 3" vertical={false} />
                            <XAxis dataKey="code" tick={{ fontSize: 11 }} />
                            <YAxis tick={{ fontSize: 11 }} />
                            <Tooltip formatter={(value) => [Number(value).toLocaleString(), 'Unidades']} />
                            <Bar dataKey="units" fill={BAR_COLOR} radius={[4, 4, 0, 0]} name="Unidades" />
                        </BarChart>
                    </ResponsiveContainer>
                </div>
            </div>

            <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden">
                <table className="table w-full">
                    <thead>
                        <tr>
                            <th>Bodega</th>
                            <th>Codigo</th>
                            <th>Productos con stock</th>
                            <th>Unidades totales</th>
                        </tr>
                    </thead>
                    <tbody>
                        {data.warehouseStats.map((ws) => (
                            <tr key={ws.code}>
                                <td className="font-medium text-gray-900">{ws.name}</td>
                                <td className="text-gray-500 font-mono text-xs">{ws.code}</td>
                                <td>{ws.products}</td>
                                <td>
                                    {ws.units === 0 ? (
                                        <span className="text-gray-400">&mdash;</span>
                                    ) : (
                                        ws.units.toLocaleString()
                                    )}
                                </td>
                            </tr>
                        ))}
                        {data.warehouseStats.length === 0 && (
                            <tr>
                                <td colSpan={4} className="text-center text-gray-400 py-8">
                                    Sin bodegas activas
                                </td>
                            </tr>
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
}
