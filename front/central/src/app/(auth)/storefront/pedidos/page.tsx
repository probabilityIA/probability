import Link from 'next/link';
import { getOrdersAction } from '@/services/modules/storefront/infra/actions';
import { getStorefrontBusinessId } from '@/shared/utils/storefront-business';
import { StorefrontPagination } from '../catalogo/pagination';

interface PageProps {
    searchParams: Promise<{ page?: string }>;
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
    new Date(date).toLocaleDateString('es-CO', { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' });

export default async function PedidosPage({ searchParams }: PageProps) {
    const params = await searchParams;
    const businessId = await getStorefrontBusinessId();
    const page = params.page ? parseInt(params.page) : 1;

    const data = await getOrdersAction({ page, page_size: 10, business_id: businessId });

    return (
        <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white dark:text-white mb-6">Mis Pedidos</h1>

            {data.data.length === 0 ? (
                <div className="text-center py-12">
                    <p className="text-gray-500 dark:text-gray-400 dark:text-gray-400 text-lg mb-4">No tienes pedidos aun</p>
                    <Link
                        href="/storefront/nuevo-pedido"
                        className="px-6 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
                    >
                        Crear primer pedido
                    </Link>
                </div>
            ) : (
                <>
                    <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
                        <table className="w-full">
                            <thead className="bg-gray-50 dark:bg-gray-700">
                                <tr>
                                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 dark:text-gray-400 uppercase">Orden</th>
                                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 dark:text-gray-400 uppercase">Fecha</th>
                                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 dark:text-gray-400 uppercase">Estado</th>
                                    <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 dark:text-gray-400 uppercase">Total</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                                {data.data.map((order: any) => (
                                    <tr key={order.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/50">
                                        <td className="px-4 py-3 text-sm font-medium">
                                            <Link href={`/storefront/pedidos/${order.id}`} className="text-indigo-600 dark:text-indigo-400 hover:underline">
                                                {order.order_number}
                                            </Link>
                                        </td>
                                        <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-300 dark:text-gray-400">
                                            {formatDate(order.created_at)}
                                        </td>
                                        <td className="px-4 py-3">
                                            <span className={`inline-flex px-2 py-1 text-xs font-medium rounded-full ${statusColors[order.status] || 'bg-gray-100 text-gray-800 dark:text-gray-100 dark:bg-gray-700 dark:text-gray-300'}`}>
                                                {order.status}
                                            </span>
                                        </td>
                                        <td className="px-4 py-3 text-sm text-right font-medium text-gray-900 dark:text-white dark:text-white">
                                            {formatPrice(order.total_amount)}
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>

                    {data.total_pages > 1 && (
                        <StorefrontPagination
                            currentPage={data.page}
                            totalPages={data.total_pages}
                            total={data.total}
                            basePath="/storefront/pedidos"
                            label="pedidos"
                        />
                    )}
                </>
            )}
        </div>
    );
}
