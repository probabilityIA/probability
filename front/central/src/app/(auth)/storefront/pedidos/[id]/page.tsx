import Link from 'next/link';
import { getOrderAction } from '@/services/modules/storefront/infra/actions';
import { getStorefrontBusinessId } from '@/shared/utils/storefront-business';

interface PageProps {
    params: Promise<{ id: string }>;
}

const statusColors: Record<string, string> = {
    pending: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400',
    processing: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400',
    completed: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400',
    cancelled: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400',
};

const formatPrice = (price: number) =>
    new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP', minimumFractionDigits: 0 }).format(price);

const formatDate = (date: string) =>
    new Date(date).toLocaleDateString('es-CO', { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit' });

export default async function OrderDetailPage({ params }: PageProps) {
    const { id } = await params;
    const businessId = await getStorefrontBusinessId();
    const order = await getOrderAction(id, businessId);

    if (!order) {
        return (
            <div className="text-center py-12">
                <p className="text-gray-500 dark:text-gray-400 text-lg mb-4">Orden no encontrada</p>
                <Link href="/storefront/pedidos" className="text-indigo-600 hover:text-indigo-700 font-medium">
                    Volver a pedidos
                </Link>
            </div>
        );
    }

    return (
        <div>
            <div className="flex items-center gap-4 mb-6">
                <Link href="/storefront/pedidos" className="text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:text-gray-200 dark:text-gray-400 dark:hover:text-gray-200">
                    &larr; Volver
                </Link>
                <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
                    Orden {order.order_number}
                </h1>
                <span className={`inline-flex px-3 py-1 text-sm font-medium rounded-full ${statusColors[order.status] || 'bg-gray-100 text-gray-800 dark:text-gray-100 dark:bg-gray-700 dark:text-gray-300'}`}>
                    {order.status}
                </span>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <div className="lg:col-span-2">
                    <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
                        <div className="px-4 py-3 bg-gray-50 dark:bg-gray-700">
                            <h2 className="text-sm font-medium text-gray-700 dark:text-gray-200 dark:text-gray-300 uppercase">Productos</h2>
                        </div>
                        <div className="divide-y divide-gray-200 dark:divide-gray-700">
                            {order.items && order.items.length > 0 ? (
                                order.items.map((item: any, index: number) => (
                                    <div key={index} className="px-4 py-3 flex items-center gap-4">
                                        {item.image_url && (
                                            <img src={item.image_url} alt={item.product_name} className="w-12 h-12 object-cover rounded" />
                                        )}
                                        <div className="flex-1 min-w-0">
                                            <p className="font-medium text-gray-900 dark:text-white">{item.product_name}</p>
                                            <p className="text-sm text-gray-500 dark:text-gray-400">
                                                {formatPrice(item.unit_price)} x {item.quantity}
                                            </p>
                                        </div>
                                        <p className="font-medium text-gray-900 dark:text-white">
                                            {formatPrice(item.total_price)}
                                        </p>
                                    </div>
                                ))
                            ) : (
                                <div className="px-4 py-6 text-center text-gray-500 dark:text-gray-400">Sin items</div>
                            )}
                        </div>
                    </div>
                </div>

                <div className="lg:col-span-1">
                    <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
                        <h2 className="text-sm font-medium text-gray-700 dark:text-gray-200 dark:text-gray-300 uppercase">Resumen</h2>
                        <div className="space-y-2 text-sm">
                            <div className="flex justify-between text-gray-600 dark:text-gray-300 dark:text-gray-400">
                                <span>Fecha</span>
                                <span>{formatDate(order.created_at)}</span>
                            </div>
                            <div className="flex justify-between text-gray-600 dark:text-gray-300 dark:text-gray-400">
                                <span>Estado</span>
                                <span className="capitalize">{order.status}</span>
                            </div>
                            {order.currency && (
                                <div className="flex justify-between text-gray-600 dark:text-gray-300 dark:text-gray-400">
                                    <span>Moneda</span>
                                    <span>{order.currency}</span>
                                </div>
                            )}
                        </div>
                        <div className="border-t border-gray-200 dark:border-gray-600 pt-3">
                            <div className="flex justify-between text-lg font-bold text-gray-900 dark:text-white">
                                <span>Total</span>
                                <span>{formatPrice(order.total_amount)}</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
