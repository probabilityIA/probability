import { WebsiteConfigManager } from '@/services/modules/website-config/ui/components/WebsiteConfigManager';

export default function WebsiteConfigPage() {
    return (
        <div className="p-6 max-w-4xl mx-auto">
            <h1 className="text-2xl font-bold text-gray-900 mb-6">Configurar Sitio Web</h1>
            <WebsiteConfigManager />
        </div>
    );
}
