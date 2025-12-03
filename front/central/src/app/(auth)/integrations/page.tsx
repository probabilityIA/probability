import { IntegrationList } from '@/services/integrations/core/ui';

export default function IntegrationsPage() {
    return (
        <div className="container mx-auto px-4 py-8">
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-3xl font-bold text-gray-900">Integraciones</h1>
                <button className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">
                    Crear Integración
                </button>
            </div>
            <IntegrationList />
        </div>
    );
}
