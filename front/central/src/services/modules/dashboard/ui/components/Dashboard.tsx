'use client';

import { useEffect, useState, useCallback, useRef } from 'react';
import { getDashboardStatsAction } from '../../infra/actions';
import { getBusinessesAction } from '@/services/auth/business/infra/actions';
import { DashboardStats } from '../../domain/types';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { TokenStorage } from '@/shared/utils/token-storage';
import { TopCustomersTable } from './TopCustomersTable';
import { TopProductsTable } from './TopProductsTable';
import { Spinner, Alert, Select } from '@/shared/ui';
import {
    ChartContainer,
    ChartTooltip,
    ChartTooltipContent,
    ChartLegend,
    ChartGradientDefs,

    CHART_GRADIENTS,
    ChartCustomGradientBar,
} from '@/shared/ui/shadcn/Chart';
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
    LabelList,
} from 'recharts';

// Colores para los gráficos
const COLORS = [
    '#3B82F6', // blue-500
    '#8B5CF6', // purple-500
    '#6366F1', // indigo-500

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
    const [userName, setUserName] = useState<string>('');

    useEffect(() => {
        const userData = TokenStorage.getUser();
        if (userData) {
            setUserName(userData.name);
        }
    }, []);

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

    // Simple card header menu
    const CardMenu = ({ items = [] }: { items?: string[] }) => {
        const [open, setOpen] = useState(false);
        const ref = useRef<HTMLDivElement | null>(null);
        useEffect(() => {
            const onDoc = (e: MouseEvent) => {
                if (!ref.current) return;
                if (!ref.current.contains(e.target as Node)) setOpen(false);
            };
            document.addEventListener('click', onDoc);
            return () => document.removeEventListener('click', onDoc);
        }, []);
        return (
            <div className="relative" ref={ref}>
                <button onClick={() => setOpen(v => !v)} className="h-6 w-6 rounded-full bg-gray-100 flex items-center justify-center text-gray-500">⋯</button>
                {open && (
                    <div className="absolute right-0 mt-2 w-40 bg-white border rounded shadow-md z-50">
                        {items.map((it, i) => (
                            <button key={i} className="w-full text-left px-3 py-2 text-sm hover:bg-gray-50">{it}</button>
                        ))}
                    </div>
                )}
            </div>
        );
    };

    // Simple SVG sparkline from numeric array
    const Sparkline = ({ data = [], color = '#8B5CF6' }: { data?: number[]; color?: string }) => {
        const w = 80; const h = 28;
        if (!Array.isArray(data) || data.length === 0) {
            return <div className="w-20 h-7" />;
        }
        const max = Math.max(...data);
        const min = Math.min(...data);
        const range = max - min || 1;
        const step = w / (data.length - 1 || 1);
        const points = data.map((v, i) => `${i * step},${h - ((v - min) / range) * h}`).join(' ');
        return (
            <svg width={w} height={h} className="w-20 h-7">
                <polyline fill="none" stroke={color} strokeWidth={2} points={points} strokeLinecap="round" strokeLinejoin="round" />
            </svg>
        );
    };

    // ModernBarChart - gráfico de barras con estilos shadcn y gradientes premium (estilo desvanecido)
    const ModernBarChart = ({ data, xKey = 'name', dataKey = 'value', height = 300, gradientType = 'purple' }: any) => {
        const gradient = CHART_GRADIENTS[gradientType as keyof typeof CHART_GRADIENTS] || CHART_GRADIENTS.purple;
        const mainColor = gradient.colors[0];

        return (
            <ChartContainer config={{}} className="h-full w-full">
                <ResponsiveContainer width="100%" height={height}>
                    <BarChart data={data} margin={{ top: 5, right: 30, left: 20, bottom: 60 }}>
                        <CartesianGrid
                            strokeDasharray="3 3"
                            stroke="#e5e7eb"
                            strokeOpacity={0.5}
                            vertical={false}
                        />
                        <XAxis
                            dataKey={xKey}
                            angle={-45}
                            textAnchor="end"
                            height={80}
                            interval={0}
                            tick={{ fontSize: 12, fill: '#6b7280' }}
                            stroke="#d1d5db"
                            axisLine={false}
                            tickLine={false}
                        />
                        <YAxis
                            tick={{ fontSize: 12, fill: '#6b7280' }}
                            stroke="#d1d5db"
                            axisLine={false}
                            tickLine={false}
                        />
                        <Tooltip cursor={false} content={<CustomTooltip />} />
                        <Bar
                            dataKey={dataKey}
                            shape={(props: any) => <ChartCustomGradientBar {...props} fill={mainColor} />}
                            className="transition-all duration-300 hover:opacity-80"
                        />
                    </BarChart>
                </ResponsiveContainer>
            </ChartContainer>
        );
    };

    const ModernPieChart = ({ data, height = 300 }: any) => {
        const chartData = (data || []).map((item: any, index: number) => ({
            ...item,
            fill: COLORS[index % COLORS.length]
        }));

        return (
            <ChartContainer config={{}} className="h-full w-full [&_.recharts-pie-label-text]:fill-gray-600">
                <ResponsiveContainer width="100%" height={height}>
                    <PieChart>
                        {/* No necesitamos GradientDefs si usamos colores sólidos de COLORS */}
                        <Tooltip cursor={false} content={<CustomTooltip />} />
                        <Pie
                            data={chartData}
                            dataKey="value"
                            nameKey="name"
                            label={({ percent }) => `${((percent ?? 0) * 100).toFixed(0)}%`}
                            labelLine={true}
                        >
                            <LabelList
                                dataKey="name"
                                position="inside"
                                fill="white"
                                className="fill-white text-xs font-bold"
                                stroke="none"
                                formatter={(val: any) => val?.toString().length > 10 ? val.toString().substring(0, 10) + '...' : val}
                            />
                        </Pie>
                    </PieChart>
                </ResponsiveContainer>
            </ChartContainer>
        );
    };

    // ModernHorizontalBarChart - gráfico de barras horizontal con gradientes
    const ModernHorizontalBarChart = ({ data, height = 300 }: any) => {
        return (
            <ChartContainer config={{}} className="h-full w-full">
                <ResponsiveContainer width="100%" height={height}>
                    <BarChart data={data} layout="vertical" margin={{ top: 5, right: 30, left: 100, bottom: 5 }}>
                        <ChartGradientDefs />
                        <CartesianGrid
                            strokeDasharray="3 3"
                            stroke="#e5e7eb"
                            strokeOpacity={0.5}
                        />
                        <XAxis
                            type="number"
                            tick={{ fontSize: 12, fill: '#6b7280' }}
                            stroke="#d1d5db"
                        />
                        <YAxis
                            type="category"
                            dataKey="name"
                            width={90}
                            tick={{ fontSize: 12, fill: '#6b7280' }}
                            stroke="#d1d5db"
                        />
                        <Tooltip cursor={false} content={<CustomTooltip />} />
                        <Bar dataKey="value" radius={[0, 8, 8, 0]}>
                            {data.map((entry: any, index: number) => (
                                <Cell
                                    key={`cell-${index}`}
                                    fill={COLORS[index % COLORS.length]}
                                    className="transition-all duration-300 hover:opacity-80"
                                />
                            ))}
                        </Bar>
                    </BarChart>
                </ResponsiveContainer>
            </ChartContainer>
        );
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

    // Agrupar y normalizar ciudades: eliminar diacríticos y capitalizar palabras (Bogotá -> Bogota)
    const locationData = (() => {
        const arr = stats.orders_by_location || [];
        const map = new Map<string, { name: string; fullName: string; value: number }>();
        const titleCase = (s: string) =>
            s
                .toLowerCase()
                .split(/\s+/)
                .filter(Boolean)
                .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
                .join(' ');
        for (const item of arr) {
            const rawCity = (item.city || '').toString();
            const rawState = (item.state || '').toString();
            const noDiacritic = rawCity.normalize('NFD').replace(/[\u0300-\u036f]/g, '');
            const displayRaw = noDiacritic || rawCity || '';
            const key = displayRaw.toLowerCase();
            const count = Number(item.order_count ?? (item as any).count ?? (item as any).value ?? 0) || 0;
            if (map.has(key)) {
                map.get(key)!.value += count;
            } else {
                const display = titleCase(displayRaw);
                map.set(key, {
                    name: display.length > 12 ? `${display.substring(0, 12)}...` : display,
                    fullName: `${display}${rawState ? `, ${rawState}` : ''}`,
                    value: count,
                });
            }
        }
        return Array.from(map.values());
    })();

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

    // Tabla para Top Products: intentar extraer precio y calcular total ganado
    const currency = (stats as any).currency || 'USD';
    const productsTableData = (stats.top_products || []).map((item: any) => {
        const units = Number(item.total_sold ?? item.units_sold ?? item.order_count ?? item.quantity_sold ?? 0) || 0;
        let price = null;
        if (item.price != null) price = Number(item.price);
        else if (item.unit_price != null) price = Number(item.unit_price);
        else if (item.average_price != null) price = Number(item.average_price);
        else if (item.total_revenue != null && units) price = Number(item.total_revenue) / units;
        const totalEarned = price != null ? price * units : Number(item.total_revenue ?? 0);
        return {
            name: item.product_name,
            sku: item.sku,
            units,
            price,
            totalEarned,
        };
    });

    // Traducciones de estados
    const STATUS_TRANSLATIONS: Record<string, string> = {
        'pending': 'Pendiente',
        'processing': 'Procesando',
        'shipped': 'Enviado',
        'in_transit': 'En Tránsito',
        'delivered': 'Entregado',
        'cancelled': 'Cancelado',
        'returned': 'Devuelto',
        'failed': 'Fallido',
        'out_for_delivery': 'En Reparto',
        'ready_to_ship': 'Listo para Enviar',
        'payment_pending': 'Pago Pendiente',
        'completed': 'Completado',
        'new': 'Nuevo',
    };

    const getTranslatedStatus = (status: string) => {
        const lowerStatus = String(status).toLowerCase();
        return STATUS_TRANSLATIONS[lowerStatus] ||
            (String(status).charAt(0).toUpperCase() + String(status).slice(1).replace(/_/g, ' '));
    };

    // Preparar datos para nuevos gráficos: Envíos
    const shipmentsByStatusData = (stats.shipments_by_status || []).map((item) => ({
        name: getTranslatedStatus(item.status),
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

    // series para sparklines (si existen)
    const ordersByDateSeries = Array.isArray((stats as any).orders_by_date)
        ? (stats as any).orders_by_date.map((d: any) => d.count ?? d.order_count ?? d.value ?? 0)
        : [];

    // Formatear revenue para el encabezado superior (fallback si no existe)
    const revenueNumber = (stats as any).total_revenue ?? (stats.total_orders ?? 90239);
    const formattedRevenue = typeof revenueNumber === 'number'
        ? revenueNumber.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })
        : String(revenueNumber);

    // Valores para la tarjeta de Total Orders (antes 'New subscriptions')
    const totalOrders = stats.total_orders ?? 0;
    const totalOrdersChangePct = (stats as any).total_orders_change_percentage ?? null;
    const totalOrdersLastWeek = (stats as any).total_orders_last_week ?? null;
    const computedTotalOrdersChange = totalOrdersChangePct !== null
        ? Math.round(totalOrdersChangePct)
        : totalOrdersLastWeek
            ? Math.round(((totalOrders - totalOrdersLastWeek) / (totalOrdersLastWeek || 1)) * 100)
            : null;

    // New orders today: try multiple possible fields or derive from a daily series
    const newOrdersToday = (() => {
        const s: any = stats as any;
        if (typeof s.orders_today === 'number') return s.orders_today;
        if (typeof s.today_orders === 'number') return s.today_orders;
        if (Array.isArray(s.orders_by_date) && s.orders_by_date.length > 0) {
            const last = s.orders_by_date[s.orders_by_date.length - 1];
            return last.count ?? last.order_count ?? last.value ?? 0;
        }
        return 0;
    })();

    // Pending orders: look for orders_by_status or fallback to shipments_by_status
    const pendingOrders = (() => {
        const s: any = stats as any;
        const candidates = s.orders_by_status || s.orders_status || s.orders_by_state || s.orders_by_statuses || null;
        const findPending = (arr: any[]) => {
            if (!Array.isArray(arr)) return null;
            const item = arr.find((it: any) => {
                const key = (it.status || it.name || it.state || '').toString().toLowerCase();
                return key === 'pending' || key === 'pendiente';
            });
            if (!item) return null;
            return item.count ?? item.order_count ?? item.value ?? null;
        };
        const fromOrders = findPending(candidates || []);
        if (fromOrders !== null) return fromOrders;
        // fallback to shipments_by_status (in case stats uses shipments)
        const ship = findPending(s.shipments_by_status || []);
        return ship ?? 0;
    })();

    return (
        <div className="space-y-6">
            {/* Top revenue header (estilo similar a la imagen) */}
            <div className="flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
                <div>
                    <h1 className="text-2xl font-semibold text-gray-700">¡Hola, {userName}! 👋  <br /> ¿Cómo va tu día? </h1>
                    <p className="mt-2 text-4xl font-extrabold bg-gradient-to-r from-purple-500 via-pink-500 to-yellow-400 bg-clip-text text-transparent capitalize">
                        Dashboard
                    </p>
                </div>
                <div className="flex items-center space-x-2">
                    {isSuperAdmin && (
                        <div className="w-64">
                            <Select
                                value={selectedBusinessId ? selectedBusinessId.toString() : 'all'}
                                onChange={(e) => {
                                    const value = e.target.value;
                                    setSelectedBusinessId(value === 'all' ? undefined : Number(value));
                                }}
                                options={[
                                    { value: 'all', label: 'Todos los negocios' },
                                    ...businesses.map((b) => ({ value: b.id.toString(), label: b.name })),
                                ]}
                                label="Seleccionar Negocio"
                                id="business-select"
                                name="business-select"
                            />
                        </div>
                    )}

                </div>
            </div>

            {/* Small summary cards under the revenue header (like the screenshot) */}
            <div className="mt-4">
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                    <div className="p-4 bg-white rounded-lg shadow-md flex items-center justify-between">
                        <div>
                            <p className="text-sm text-gray-500">Órdenes Totales</p>
                            <div className="mt-2 flex items-center space-x-4">
                                <div>
                                    <p className="text-2xl font-bold text-gray-900">{totalOrders.toLocaleString()}</p>
                                    {computedTotalOrdersChange !== null ? (
                                        <p className={`text-xs ${computedTotalOrdersChange >= 0 ? 'text-green-600' : 'text-amber-500'}`}>
                                            {computedTotalOrdersChange >= 0 ? '↑' : '↓'} {Math.abs(computedTotalOrdersChange)}% compared to last week
                                        </p>
                                    ) : (
                                        <p className="text-xs text-gray-400">Aún no hay comparación semanal</p>
                                    )}
                                </div>
                                <div className="w-20 h-10">
                                    <svg viewBox="0 0 100 40" preserveAspectRatio="none" className="w-full h-full">
                                        <path d="M0,30 C25,10 50,12 75,2 100,0" fill="none" stroke="#8B5CF6" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round" />
                                    </svg>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div className="p-4 bg-white rounded-lg shadow-md flex items-center justify-between">
                        <div>
                            <p className="text-sm text-gray-500">Órdenes del día</p>
                            <div className="mt-2 flex items-center space-x-4">
                                <div>
                                    <p className="text-2xl font-bold text-gray-900">{newOrdersToday.toLocaleString()}</p>
                                    <p className="text-xs text-gray-400">Actuales</p>
                                </div>
                                <div className="w-20 h-10">
                                    <svg viewBox="0 0 100 40" preserveAspectRatio="none" className="w-full h-full">
                                        <path d="M0,30 C20,20 40,15 60,10 80,6 100,8" fill="none" stroke="#F97316" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round" />
                                    </svg>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div className="p-4 bg-white rounded-lg shadow-md flex items-center justify-between">
                        <div>
                            <p className="text-sm text-gray-500">Órdenes Pendientes</p>
                            <div className="mt-2 flex items-center space-x-4">
                                <div>
                                    <p className="text-2xl font-bold text-gray-900">{(pendingOrders || 0).toLocaleString()}</p>
                                    <p className="text-xs text-gray-400">Ordenes actualmente en estado pendiente</p>
                                </div>
                                <div className="w-20 h-10">
                                    <svg viewBox="0 0 100 40" preserveAspectRatio="none" className="w-full h-full">
                                        <path d="M0,30 C15,25 30,20 50,18 70,16 85,14 100,10" fill="none" stroke="#06B6D4" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round" />
                                    </svg>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {loading && stats && (
                <div className="flex items-center justify-center py-4">
                    <Spinner size="md" color="primary" text="Actualizando estadísticas..." />
                </div>
            )}


            {/* Gráficas en 3 columnas */}
            <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
                {/* Orders by Location - Gráfico de Barras (primary, wide) */}
                <div className="bg-white rounded-2xl shadow-md p-6 xl:col-span-2">
                    <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center">
                            <MapPinIcon className="w-5 h-5 text-gray-400 mr-2" />
                            <h2 className="text-lg font-semibold text-gray-900">Órdenes por Ubicación</h2>
                        </div>
                        <CardMenu items={["Ver detalles", "Exportar", "Refrescar"]} />
                    </div>
                    {(locationData || []).length > 0 ? (
                        <ModernBarChart data={locationData} xKey="name" dataKey="value" height={340} gradientType="purple" />
                    ) : (
                        <p className="text-sm text-gray-500">No hay datos disponibles</p>
                    )}
                </div>


                {/* Shipments by Carrier */}
                <div className="bg-white rounded-2xl shadow-md p-6">
                    <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center">
                            <ArchiveBoxIcon className="w-5 h-5 text-gray-400 mr-2" />
                            <h2 className="text-lg font-semibold text-gray-900">Envíos por Transportista</h2>
                        </div>
                        <CardMenu items={["Ver detalles", "Exportar", "Refrescar"]} />
                    </div>
                    {(stats.shipments_by_carrier || []).length > 0 ? (
                        <ModernBarChart data={shipmentsByCarrierData} xKey="name" dataKey="value" height={340} gradientType="indigo" />
                    ) : (
                        <p className="text-sm text-gray-500">No hay datos disponibles</p>
                    )}
                </div>


                {/* Orders by Integration Type - Gráfico de Pastel 
                <div className="bg-white rounded-2xl shadow-md p-6">
                    <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center">
                            <ChartBarIcon className="w-5 h-5 text-gray-400 mr-2" />
                            <h2 className="text-lg font-semibold text-gray-900">Órdenes por Tipo de Integración</h2>
                        </div>
                        <CardMenu items={["Ver detalles","Exportar","Refrescar"]} />
                    </div>
                    {(integrationData || []).length > 0 ? (
                        <ModernPieChart data={integrationData} height={300} />
                    ) : (
                        <p className="text-sm text-gray-500">No hay datos disponibles</p>
                    )}
                </div>  
                */}

                {/* Top Customers - Gráfico de Barras */}
                {/* Top Customers - Tabla Interactiva */}
                <div className="bg-white rounded-2xl shadow-md p-6 lg:col-span-2">
                    <div className="flex items-center justify-between mb-2">
                        <div className="flex items-center">
                            <UserGroupIcon className="w-5 h-5 text-gray-400 mr-2" />
                            <h2 className="text-lg font-semibold text-gray-900">Mejores Clientes</h2>
                        </div>
                        <CardMenu items={["Exportar", "Refrescar"]} />
                    </div>
                    {(stats.top_customers || []).length > 0 ? (
                        <TopCustomersTable data={stats.top_customers} />
                    ) : (
                        <p className="text-sm text-gray-500">No hay datos disponibles</p>
                    )}
                </div>

                {/* Shipments by Status */}
                <div className="bg-white rounded-2xl shadow-md p-6">
                    <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center">
                            <ArchiveBoxIcon className="w-5 h-5 text-gray-400 mr-2" />
                            <h2 className="text-lg font-semibold text-gray-900">Estado de los Envios</h2>
                        </div>
                    </div>
                    {(stats.shipments_by_status || []).length > 0 ? (
                        <ModernPieChart data={shipmentsByStatusData} height={300} />
                    ) : (
                        <p className="text-sm text-gray-500">No hay datos disponibles</p>
                    )}
                </div>

                {/* Top Drivers
                <div className="bg-white rounded-2xl shadow-md p-6">
                    <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center">
                            <TruckIcon className="w-5 h-5 text-gray-400 mr-2" />
                            <h2 className="text-lg font-semibold text-gray-900">Top Transportadores</h2>
                        </div>
                        <CardMenu items={["Ver detalles","Exportar","Refrescar"]} />
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
                </div>  */}

                {/* Drivers by Location 
                <div className="bg-white rounded-2xl shadow-md p-6">
                    <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center">
                            <TruckIcon className="w-5 h-5 text-gray-400 mr-2" />
                            <h2 className="text-lg font-semibold text-gray-900">Transportadores por Ubicación</h2>
                        </div>
                        <CardMenu items={["Ver detalles","Exportar","Refrescar"]} />
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
                </div> */}

                {/* Top Products as Table */}
                {/* Top Products - Tabla Interactiva */}
                <div className="bg-white rounded-2xl shadow-md p-6 lg:col-span-2">
                    <div className="flex items-center justify-between mb-2">
                        <div className="flex items-center">
                            <CubeIcon className="w-5 h-5 text-gray-400 mr-2" />
                            <h2 className="text-lg font-semibold text-gray-900">Productos Más Vendidos</h2>
                        </div>
                        <CardMenu items={["Exportar", "Refrescar"]} />
                    </div>
                    {productsTableData.length > 0 ? (
                        <TopProductsTable data={productsTableData} />
                    ) : (
                        <p className="text-sm text-gray-500">No hay datos disponibles</p>
                    )}
                </div>

                {/* Products by Category 
                <div className="bg-white rounded-2xl shadow-md p-6">
                    <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center">
                            <CubeIcon className="w-5 h-5 text-gray-400 mr-2" />
                            <h2 className="text-lg font-semibold text-gray-900">Productos por Categoría</h2>
                        </div>
                        <CardMenu items={["Ver detalles","Exportar","Refrescar"]} />
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
                </div> */}

                {/* Products by Brand 
                <div className="bg-white rounded-2xl shadow-md p-6">
                    <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center">
                            <CubeIcon className="w-5 h-5 text-gray-400 mr-2" />
                            <h2 className="text-lg font-semibold text-gray-900">Productos por Marca</h2>
                        </div>
                        <CardMenu items={["Ver detalles","Exportar","Refrescar"]} />
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
                </div> */}




                {/* Shipments by Warehouse 
                <div className="bg-white rounded-2xl shadow-md p-6">
                    <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center">
                            <ArchiveBoxIcon className="w-5 h-5 text-gray-400 mr-2" />
                            <h2 className="text-lg font-semibold text-gray-900">Envíos por Almacén</h2>
                        </div>
                        <CardMenu items={["Ver detalles","Exportar","Refrescar"]} />
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
                </div> */}
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
