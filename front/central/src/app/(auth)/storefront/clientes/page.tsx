import { getStorefrontBusinessId } from '@/shared/utils/storefront-business';
import { ClientsManager } from '@/services/modules/storefront/ui/components/ClientsManager';

export default async function ClientesPage() {
    const businessId = await getStorefrontBusinessId();

    return (
        <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Clientes</h1>
            <ClientsManager businessId={businessId} />
        </div>
    );
}
