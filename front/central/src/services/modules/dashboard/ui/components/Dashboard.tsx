'use client';

import { useEffect, useState, useCallback } from 'react';
import { getDashboardStatsAction } from '../../infra/actions';
import { getBusinessesAction } from '@/services/auth/business/infra/actions';
import { DashboardStats } from '../../domain/types';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { Spinner, Alert, Select } from '@/shared/ui';
import {
    ShoppingBagIcon,
    UserGroupIcon,
    MapPinIcon,
    ChartBarIcon,
    TruckIcon,
    CubeIcon,
    ArchiveBoxIcon,
    BuildingOfficeIcon,
} from '@heroicons/react/24/outline';
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

// Colores para los gráficos
const COLORS = [
    '#3B82F6', // blue-500
    '#10B981', // green-500
    '#F59E0B', // amber-500
    '#EF4444', // red-500
    '#8B5CF6', // purple-500
    '#F97316', // orange-500
    '#06B6D4', // cyan-500
    '#EC4899', // pink-500
    '#6366F1', // indigo-500
    '#14B8A6', // teal-500
];

interface Business {
    id: number;
    name: string;
}

export default function Dashboard() {
    const { isSuperAdmin } = usePermissions();
    const [stats, setStats] = useState<DashboardStats | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | undefined>(undefined);
    const [businesses, setBusinesses] = useState<Business[]>([]);
    const [loadingBusinesses, setLoadingBusinesses] = useState(false);

    // Cargar lista de businesses si es super admin
    useEffect(() => {
        if (isSuperAdmin) {
            const fetchBusinesses = async () => {
                try {
                    setLoadingBusinesses(true);
                    const response = await getBusinessesAction({ page: 1, per_page: 100 });
                    setBusinesses(response.data || []);
                } catch (err: any) {
                    console.error('Error fetching businesses:', err);
                } finally {
                    setLoadingBusinesses(false);
                }
            };
            fetchBusinesses();
        }
    }, [isSuperAdmin]);

    // Cargar estadísticas
    const fetchStats = useCallback(async () => {
        try {
            setLoading(true);
            setError(null);
            const response = await getDashboardStatsAction(selectedBusinessId);
            setStats(response.data);
        } catch (err: any) {
            console.error('Error fetching dashboard stats:', err);
            setError(err.message || 'Error al cargar las estadísticas');
        } finally {
            setLoading(false);
        }
    }, [selectedBusinessId]);

    useEffect(() => {
        fetchStats();
    }, [fetchStats]);

    // Custom tooltip para mostrar el nombre completo
    const CustomTooltip = ({ active, payload }: any) => {
        if (active && payload && payload.length) {
            return (
                <div className="bg-white p-3 border border-gray-200 rounded-lg shadow-lg">
                    <p className="font-semibold text-gray-900">
                        {payload[0].payload.fullName || payload[0].payload.name}
                    </p>
                    {payload[0].payload.email && (
                        <p className="text-xs text-gray-500 mt-1">{payload[0].payload.email}</p>
                    )}
                    <p className="text-sm text-gray-700 mt-1">
                        <span className="font-bold">{payload[0].value.toLocaleString()}</span> {payload[0].payload.unit || 'órdenes'}
                    </p>
                </div>
            );
        }
        return null;
    };

    if (loading && !stats) {
        return (
            <div className="flex items-center justify-center py-12">
                <Spinner size="xl" color="primary" text="Cargando estadísticas..." />
            </div>
        );
    }

    if (error && !stats) {
        return (
            <div className="p-4">
                <Alert type="error">{error}</Alert>
            </div>
        );
    }

    if (!stats) {
        return (
            <div className="p-4">
                <Alert type="warning">No hay estadísticas disponibles</Alert>
            </div>
        );
    }

    // Preparar datos para gráficos existentes
    const integrationData = (stats.orders_by_integration_type || []).map((item) => ({
        name: item.integration_type.charAt(0).toUpperCase() + item.integration_type.slice(1),
        value: item.count,
    }));

    const locationData = (stats.orders_by_location || []).map((item) => ({
        name: item.city.length > 12 ? `${item.city.substring(0, 12)}...` : item.city,
        fullName: `${item.city}${item.state ? `, ${item.state}` : ''}`,
        value: item.order_count,
    }));

    const customersData = (stats.top_customers || []).map((item) => ({
        name: item.customer_name.length > 15 ? `${item.customer_name.substring(0, 15)}...` : item.customer_name,
        fullName: item.customer_name,
        email: item.customer_email,
        value: item.order_count,
        unit: 'órdenes',
    }));

    // Preparar datos para nuevos gráficos: Transportadores
    const driversData = (stats.top_drivers || []).map((item) => ({
        name: item.driver_name.length > 15 ? `${item.driver_name.substring(0, 15)}...` : item.driver_name,
        fullName: item.driver_name,
        value: item.order_count,
        unit: 'órdenes',
    }));

    const driversByLocationData = (stats.drivers_by_location || []).map((item) => ({
        name: item.city.length > 12 ? `${item.city.substring(0, 12)}...` : item.city,
        fullName: `${item.driver_name} - ${item.city}${item.state ? `, ${item.state}` : ''}`,
        driverName: item.driver_name,
        city: item.city,
        state: item.state,
        value: item.order_count,
        unit: 'órdenes',
    }));

    // Preparar datos para nuevos gráficos: Productos
    const productsData = (stats.top_products || []).map((item) => ({
        name: item.product_name.length > 20 ? `${item.product_name.substring(0, 20)}...` : item.product_name,
        fullName: item.product_name,
        sku: item.sku,
        value: item.order_count,
        totalSold: item.total_sold,
        unit: 'órdenes',
    }));

    const productsByCategoryData = (stats.products_by_category || []).map((item) => ({
        name: item.category || 'Sin categoría',
        value: item.count,
    }));

    const productsByBrandData = (stats.products_by_brand || []).map((item) => ({
        name: item.brand || 'Sin marca',
        value: item.count,
    }));

    // Preparar datos para nuevos gráficos: Envíos
    const shipmentsByStatusData = (stats.shipments_by_status || []).map((item) => ({
        name: item.status.charAt(0).toUpperCase() + item.status.slice(1).replace(/_/g, ' '),
        value: item.count,
    }));

    const shipmentsByCarrierData = (stats.shipments_by_carrier || []).map((item) => ({
        name: item.carrier.length > 15 ? `${item.carrier.substring(0, 15)}...` : item.carrier,
        fullName: item.carrier,
        value: item.count,
        unit: 'envíos',
    }));

    const shipmentsByWarehouseData = (stats.shipments_by_warehouse || []).map((item) => ({
        name: item.warehouse_name.length > 15 ? `${item.warehouse_name.substring(0, 15)}...` : item.warehouse_name,
        fullName: item.warehouse_name,
        value: item.count,
        unit: 'envíos',
    }));

    // Preparar datos para businesses (solo super admin)
    const businessesData = (stats.orders_by_business || []).map((item) => ({
        name: item.business_name.length > 15 ? `${item.business_name.substring(0, 15)}...` : item.business_name,
        fullName: item.business_name,
        value: item.order_count,
        unit: 'órdenes',
    }));

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
                    <p className="mt-1 text-sm text-gray-500">
                        Resumen general de tus órdenes y estadísticas
                    </p>
                </div>
                {/* Filtro de business (solo super admin) */}
                {isSuperAdmin && (
                    <div className="w-64">
                        <Select
                            label="Filtrar por Business"
                            options={[
                                { value: '', label: 'Todos los businesses' },
                                ...businesses.map(b => ({ value: String(b.id), label: b.name })),
                            ]}
                            value={selectedBusinessId ? String(selectedBusinessId) : ''}
                            onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {
                                const value = e.target.value;
                                setSelectedBusinessId(value ? Number(value) : undefined);
                            }}
                        />
                    </div>
                )}
            </div>

            {loading && stats && (
                <div className="flex items-center justify-center py-4">
                    <Spinner size="md" color="primary" text="Actualizando estadísticas..." />
                </div>
            )}

            {/* Total Orders Card */}
            <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
                <div className="flex items-center justify-between">
                    <div>
                        <p className="text-sm font-medium text-gray-600">Total de Órdenes</p>
                        <p className="mt-2 text-4xl font-bold text-gray-900">
                            {stats.total_orders.toLocaleString()}
                        </p>
                    </div>
                    <div className="p-3 bg-blue-100 rounded-full">
                        <ShoppingBagIcon className="w-8 h-8 text-blue-600" />
                    </div>
                </div>
            </div>

            {/* Gráficas en 3 columnas */}
            <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
                {/* Orders by Integration Type - Gráfico de Torta */}
                <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
                    <div className="flex items-center mb-4">
                        <ChartBarIcon className="w-5 h-5 text-gray-400 mr-2" />
                        <h2 className="text-lg font-semibold text-gray-900">
                            Órdenes por Tipo de Integración
                        </h2>
                    </div>
                    {(stats.orders_by_integration_type || []).length > 0 ? (
                        <ResponsiveContainer width="100%" height={300}>
                            <PieChart>
                                <Pie
                                    data={integrationData}
                                    cx="50%"
                                    cy="50%"
                                    labelLine={false}
                                    label={({ name, percent }) =>
                                        `${name}: ${((percent || 0) * 100).toFixed(0)}%`
                                    }
                                    outerRadius={100}
                                    fill="#8884d8"
                                    dataKey="value"
                                >
                                    {integrationData.map((entry, index) => (
                                        <Cell
                                            key={`cell-${index}`}
                                            fill={COLORS[index % COLORS.length]}
                                        />
                                    ))}
                                </Pie>
                                <Tooltip formatter={(value: number | undefined) => (value || 0).toLocaleString()} />
                                <Legend
                                    verticalAlign="bottom"
                                    height={36}
                                    formatter={(value) => (
                                        <span className="text-sm text-gray-700">{value}</span>
                                    )}
                                />
                            </PieChart>
                        </ResponsiveContainer>
                    ) : (
                        <p className="text-sm text-gray-500">No hay datos disponibles</p>
                    )}
                </div>

                {/* Orders by Location - Gráfico de Barras */}
                <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
                    <div className="flex items-center mb-4">
                        <MapPinIcon className="w-5 h-5 text-gray-400 mr-2" />
                        <h2 className="text-lg font-semibold text-gray-900">
                            Órdenes por Ubicación
                        </h2>
                    </div>
                    {(stats.orders_by_location || []).length > 0 ? (
                        <ResponsiveContainer width="100%" height={300}>
                            <BarChart data={locationData} margin={{ top: 5, right: 30, left: 20, bottom: 60 }}>
                                <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                                <XAxis
                                    dataKey="name"
                                    angle={-45}
                                    textAnchor="end"
                                    height={80}
                                    interval={0}
                                    tick={{ fontSize: 12 }}
                                />
                                <YAxis tick={{ fontSize: 12 }} />
                                <Tooltip content={<CustomTooltip />} />
                                <Bar dataKey="value" fill="#3B82F6" radius={[4, 4, 0, 0]}>
                                    {locationData.map((entry, index) => (
                                        <Cell
                                            key={`cell-${index}`}
                                            fill={COLORS[index % COLORS.length]}
                                        />
                                    ))}
                                </Bar>
                            </BarChart>
                        </ResponsiveContainer>
                    ) : (
                        <p className="text-sm text-gray-500">No hay datos disponibles</p>
                    )}
                </div>

                {/* Top Customers - Gráfico de Barras */}
                <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
                    <div className="flex items-center mb-4">
                        <UserGroupIcon className="w-5 h-5 text-gray-400 mr-2" />
                        <h2 className="text-lg font-semibold text-gray-900">Top Clientes</h2>
                    </div>
                    {(stats.top_customers || []).length > 0 ? (
                        <ResponsiveContainer width="100%" height={300}>
                            <BarChart data={customersData} layout="vertical" margin={{ top: 5, right: 30, left: 100, bottom: 5 }}>
                                <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                                <XAxis type="number" tick={{ fontSize: 12 }} />
                                <YAxis
                                    type="category"
                                    dataKey="name"
                                    width={90}
                                    tick={{ fontSize: 12 }}
                                />
                                <Tooltip content={<CustomTooltip />} />
                                <Bar dataKey="value" fill="#10B981" radius={[0, 4, 4, 0]}>
                                    {customersData.map((entry, index) => (
                                        <Cell
                                            key={`cell-${index}`}
                                            fill={COLORS[index % COLORS.length]}
                                        />
                                    ))}
                                </Bar>
                            </BarChart>
                        </ResponsiveContainer>
                    ) : (
                        <p className="text-sm text-gray-500">No hay datos disponibles</p>
                    )}
                </div>

                {/* Top Drivers */}
                <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
                    <div className="flex items-center mb-4">
                        <TruckIcon className="w-5 h-5 text-gray-400 mr-2" />
                        <h2 className="text-lg font-semibold text-gray-900">Top Transportadores</h2>
                    </div>
                    {(stats.top_drivers || []).length > 0 ? (
                        <ResponsiveContainer width="100%" height={300}>
                            <BarChart data={driversData} layout="vertical" margin={{ top: 5, right: 30, left: 80, bottom: 5 }}>
                                <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                                <XAxis type="number" tick={{ fontSize: 12 }} />
                                <YAxis
                                    type="category"
                                    dataKey="name"
                                    width={75}
                                    tick={{ fontSize: 12 }}
                                />
                                <Tooltip content={<CustomTooltip />} />
                                <Bar dataKey="value" fill="#8B5CF6" radius={[0, 4, 4, 0]}>
                                    {driversData.map((entry, index) => (
                                        <Cell
                                            key={`cell-${index}`}
                                            fill={COLORS[index % COLORS.length]}
                                        />
                                    ))}
                                </Bar>
                            </BarChart>
                        </ResponsiveContainer>
                    ) : (
                        <p className="text-sm text-gray-500">No hay datos disponibles</p>
                    )}
                </div>

                {/* Drivers by Location */}
                <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
                    <div className="flex items-center mb-4">
                        <TruckIcon className="w-5 h-5 text-gray-400 mr-2" />
                        <h2 className="text-lg font-semibold text-gray-900">Transportadores por Ubicación</h2>
                    </div>
                    {(stats.drivers_by_location || []).length > 0 ? (
                        <ResponsiveContainer width="100%" height={300}>
                            <BarChart data={driversByLocationData} margin={{ top: 5, right: 30, left: 20, bottom: 60 }}>
                                <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                                <XAxis
                                    dataKey="name"
                                    angle={-45}
                                    textAnchor="end"
                                    height={80}
                                    interval={0}
                                    tick={{ fontSize: 12 }}
                                />
                                <YAxis tick={{ fontSize: 12 }} />
                                <Tooltip content={<CustomTooltip />} />
                                <Bar dataKey="value" fill="#F59E0B" radius={[4, 4, 0, 0]}>
                                    {driversByLocationData.map((entry, index) => (
                                        <Cell
                                            key={`cell-${index}`}
                                            fill={COLORS[index % COLORS.length]}
                                        />
                                    ))}
                                </Bar>
                            </BarChart>
                        </ResponsiveContainer>
                    ) : (
                        <p className="text-sm text-gray-500">No hay datos disponibles</p>
                    )}
                </div>

                {/* Top Products */}
                <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
                    <div className="flex items-center mb-4">
                        <CubeIcon className="w-5 h-5 text-gray-400 mr-2" />
                        <h2 className="text-lg font-semibold text-gray-900">Top Productos</h2>
                    </div>
                    {(stats.top_products || []).length > 0 ? (
                        <ResponsiveContainer width="100%" height={300}>
                            <BarChart data={productsData} layout="vertical" margin={{ top: 5, right: 30, left: 100, bottom: 5 }}>
                                <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                                <XAxis type="number" tick={{ fontSize: 12 }} />
                                <YAxis
                                    type="category"
                                    dataKey="name"
                                    width={95}
                                    tick={{ fontSize: 11 }}
                                />
                                <Tooltip content={<CustomTooltip />} />
                                <Bar dataKey="value" fill="#06B6D4" radius={[0, 4, 4, 0]}>
                                    {productsData.map((entry, index) => (
                                        <Cell
                                            key={`cell-${index}`}
                                            fill={COLORS[index % COLORS.length]}
                                        />
                                    ))}
                                </Bar>
                            </BarChart>
                        </ResponsiveContainer>
                    ) : (
                        <p className="text-sm text-gray-500">No hay datos disponibles</p>
                    )}
                </div>

                {/* Products by Category */}
                <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
                    <div className="flex items-center mb-4">
                        <CubeIcon className="w-5 h-5 text-gray-400 mr-2" />
                        <h2 className="text-lg font-semibold text-gray-900">Productos por Categoría</h2>
                    </div>
                    {(stats.products_by_category || []).length > 0 ? (
                        <ResponsiveContainer width="100%" height={300}>
                            <PieChart>
                                <Pie
                                    data={productsByCategoryData}
                                    cx="50%"
                                    cy="50%"
                                    labelLine={false}
                                    label={({ name, percent }) =>
                                        `${name}: ${((percent || 0) * 100).toFixed(0)}%`
                                    }
                                    outerRadius={80}
                                    fill="#8884d8"
                                    dataKey="value"
                                >
                                    {productsByCategoryData.map((entry, index) => (
                                        <Cell
                                            key={`cell-${index}`}
                                            fill={COLORS[index % COLORS.length]}
                                        />
                                    ))}
                                </Pie>
                                <Tooltip formatter={(value: number | undefined) => (value || 0).toLocaleString()} />
                            </PieChart>
                        </ResponsiveContainer>
                    ) : (
                        <p className="text-sm text-gray-500">No hay datos disponibles</p>
                    )}
                </div>

                {/* Products by Brand */}
                <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
                    <div className="flex items-center mb-4">
                        <CubeIcon className="w-5 h-5 text-gray-400 mr-2" />
                        <h2 className="text-lg font-semibold text-gray-900">Productos por Marca</h2>
                    </div>
                    {(stats.products_by_brand || []).length > 0 ? (
                        <ResponsiveContainer width="100%" height={300}>
                            <BarChart data={productsByBrandData} margin={{ top: 5, right: 30, left: 20, bottom: 60 }}>
                                <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                                <XAxis
                                    dataKey="name"
                                    angle={-45}
                                    textAnchor="end"
                                    height={80}
                                    interval={0}
                                    tick={{ fontSize: 12 }}
                                />
                                <YAxis tick={{ fontSize: 12 }} />
                                <Tooltip formatter={(value: number | undefined) => (value || 0).toLocaleString()} />
                                <Bar dataKey="value" fill="#EC4899" radius={[4, 4, 0, 0]}>
                                    {productsByBrandData.map((entry, index) => (
                                        <Cell
                                            key={`cell-${index}`}
                                            fill={COLORS[index % COLORS.length]}
                                        />
                                    ))}
                                </Bar>
                            </BarChart>
                        </ResponsiveContainer>
                    ) : (
                        <p className="text-sm text-gray-500">No hay datos disponibles</p>
                    )}
                </div>

                {/* Shipments by Status */}
                <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
                    <div className="flex items-center mb-4">
                        <ArchiveBoxIcon className="w-5 h-5 text-gray-400 mr-2" />
                        <h2 className="text-lg font-semibold text-gray-900">Envíos por Estado</h2>
                    </div>
                    {(stats.shipments_by_status || []).length > 0 ? (
                        <ResponsiveContainer width="100%" height={300}>
                            <PieChart>
                                <Pie
                                    data={shipmentsByStatusData}
                                    cx="50%"
                                    cy="50%"
                                    labelLine={false}
                                    label={({ name, percent }) =>
                                        `${name}: ${((percent || 0) * 100).toFixed(0)}%`
                                    }
                                    outerRadius={80}
                                    fill="#8884d8"
                                    dataKey="value"
                                >
                                    {shipmentsByStatusData.map((entry, index) => (
                                        <Cell
                                            key={`cell-${index}`}
                                            fill={COLORS[index % COLORS.length]}
                                        />
                                    ))}
                                </Pie>
                                <Tooltip formatter={(value: number | undefined) => (value || 0).toLocaleString()} />
                                <Legend
                                    verticalAlign="bottom"
                                    height={36}
                                    formatter={(value) => (
                                        <span className="text-sm text-gray-700">{value}</span>
                                    )}
                                />
                            </PieChart>
                        </ResponsiveContainer>
                    ) : (
                        <p className="text-sm text-gray-500">No hay datos disponibles</p>
                    )}
                </div>

                {/* Shipments by Carrier */}
                <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
                    <div className="flex items-center mb-4">
                        <ArchiveBoxIcon className="w-5 h-5 text-gray-400 mr-2" />
                        <h2 className="text-lg font-semibold text-gray-900">Envíos por Transportista</h2>
                    </div>
                    {(stats.shipments_by_carrier || []).length > 0 ? (
                        <ResponsiveContainer width="100%" height={300}>
                            <BarChart data={shipmentsByCarrierData} margin={{ top: 5, right: 30, left: 20, bottom: 60 }}>
                                <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                                <XAxis
                                    dataKey="name"
                                    angle={-45}
                                    textAnchor="end"
                                    height={80}
                                    interval={0}
                                    tick={{ fontSize: 12 }}
                                />
                                <YAxis tick={{ fontSize: 12 }} />
                                <Tooltip content={<CustomTooltip />} />
                                <Bar dataKey="value" fill="#6366F1" radius={[4, 4, 0, 0]}>
                                    {shipmentsByCarrierData.map((entry, index) => (
                                        <Cell
                                            key={`cell-${index}`}
                                            fill={COLORS[index % COLORS.length]}
                                        />
                                    ))}
                                </Bar>
                            </BarChart>
                        </ResponsiveContainer>
                    ) : (
                        <p className="text-sm text-gray-500">No hay datos disponibles</p>
                    )}
                </div>

                {/* Shipments by Warehouse */}
                <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
                    <div className="flex items-center mb-4">
                        <ArchiveBoxIcon className="w-5 h-5 text-gray-400 mr-2" />
                        <h2 className="text-lg font-semibold text-gray-900">Envíos por Almacén</h2>
                    </div>
                    {(stats.shipments_by_warehouse || []).length > 0 ? (
                        <ResponsiveContainer width="100%" height={300}>
                            <BarChart data={shipmentsByWarehouseData} margin={{ top: 5, right: 30, left: 20, bottom: 60 }}>
                                <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                                <XAxis
                                    dataKey="name"
                                    angle={-45}
                                    textAnchor="end"
                                    height={80}
                                    interval={0}
                                    tick={{ fontSize: 12 }}
                                />
                                <YAxis tick={{ fontSize: 12 }} />
                                <Tooltip content={<CustomTooltip />} />
                                <Bar dataKey="value" fill="#14B8A6" radius={[4, 4, 0, 0]}>
                                    {shipmentsByWarehouseData.map((entry, index) => (
                                        <Cell
                                            key={`cell-${index}`}
                                            fill={COLORS[index % COLORS.length]}
                                        />
                                    ))}
                                </Bar>
                            </BarChart>
                        </ResponsiveContainer>
                    ) : (
                        <p className="text-sm text-gray-500">No hay datos disponibles</p>
                    )}
                </div>
            </div>

            {/* Businesses (solo super admin, solo cuando NO hay filtro aplicado) */}
            {isSuperAdmin && !selectedBusinessId && stats.orders_by_business && stats.orders_by_business.length > 0 && (
                <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
                    <div className="flex items-center mb-4">
                        <BuildingOfficeIcon className="w-5 h-5 text-gray-400 mr-2" />
                        <h2 className="text-lg font-semibold text-gray-900">Órdenes por Business</h2>
                    </div>
                    {businessesData.length > 0 ? (
                        <ResponsiveContainer width="100%" height={400}>
                            <BarChart data={businessesData} layout="vertical" margin={{ top: 5, right: 30, left: 150, bottom: 5 }}>
                                <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                                <XAxis type="number" tick={{ fontSize: 12 }} />
                                <YAxis
                                    type="category"
                                    dataKey="name"
                                    width={145}
                                    tick={{ fontSize: 12 }}
                                />
                                <Tooltip content={<CustomTooltip />} />
                                <Bar dataKey="value" fill="#F97316" radius={[0, 4, 4, 0]}>
                                    {businessesData.map((entry, index) => (
                                        <Cell
                                            key={`cell-${index}`}
                                            fill={COLORS[index % COLORS.length]}
                                        />
                                    ))}
                                </Bar>
                            </BarChart>
                        </ResponsiveContainer>
                    ) : (
                        <p className="text-sm text-gray-500">No hay datos disponibles</p>
                    )}
                </div>
            )}
        </div>
    );
}
