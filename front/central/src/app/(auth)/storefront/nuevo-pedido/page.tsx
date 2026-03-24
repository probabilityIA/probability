import { getStorefrontBusinessId } from '@/shared/utils/storefront-business';
import { NewOrderForm } from '@/services/modules/storefront/ui/components/NewOrderForm';

export default async function NuevoPedidoPage() {
    const businessId = await getStorefrontBusinessId();

    return (
        <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white dark:text-white mb-6">Nuevo Pedido</h1>
            <NewOrderForm businessId={businessId} />
        </div>
    );
}
