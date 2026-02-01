'use client';

import { Integration } from '@/services/integrations/core/domain/types';
import { Badge, Button } from '@/shared/ui';
import { CheckCircleIcon, XCircleIcon } from '@heroicons/react/24/solid';

interface SoftpymesIntegrationViewProps {
    integration: Integration;
    onEdit?: () => void;
    onTest?: () => void;
    onToggleActive?: () => void;
}

export function SoftpymesIntegrationView({
    integration,
    onEdit,
    onTest,
    onToggleActive
}: SoftpymesIntegrationViewProps) {
    const config = integration.config as any;

    return (
        <div className="bg-white border rounded-lg p-6 hover:shadow-md transition-shadow">
            {/* Header */}
            <div className="flex items-start justify-between mb-4">
                <div className="flex items-center gap-3">
                    <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
                        <span className="text-2xl">📄</span>
                    </div>
                    <div>
                        <h3 className="text-lg font-semibold text-gray-900">
                            {integration.name}
                        </h3>
                        <p className="text-sm text-gray-500">
                            Softpymes - Facturación Electrónica
                        </p>
                    </div>
                </div>

                <div className="flex items-center gap-2">
                    {integration.is_active ? (
                        <Badge type="success">
                            <CheckCircleIcon className="w-4 h-4 mr-1" />
                            Activo
                        </Badge>
                    ) : (
                        <Badge type="error">
                            <XCircleIcon className="w-4 h-4 mr-1" />
                            Inactivo
                        </Badge>
                    )}
                </div>
            </div>

            {/* Configuration Info */}
            <div className="space-y-2 mb-4">
                <div className="flex justify-between text-sm">
                    <span className="text-gray-600">Empresa:</span>
                    <span className="font-medium text-gray-900">
                        {config?.company_name || 'No configurado'}
                    </span>
                </div>
                <div className="flex justify-between text-sm">
                    <span className="text-gray-600">NIT:</span>
                    <span className="font-medium text-gray-900">
                        {config?.company_nit || 'No configurado'}
                    </span>
                </div>
                <div className="flex justify-between text-sm">
                    <span className="text-gray-600">Modo:</span>
                    <span className="font-medium text-gray-900">
                        {config?.test_mode ? (
                            <Badge type="warning">Pruebas</Badge>
                        ) : (
                            <Badge type="primary">Producción</Badge>
                        )}
                    </span>
                </div>
            </div>

            {/* Actions */}
            <div className="flex gap-2 pt-4 border-t">
                {onEdit && (
                    <Button
                        size="sm"
                        variant="outline"
                        onClick={onEdit}
                    >
                        Editar
                    </Button>
                )}
                {onTest && (
                    <Button
                        size="sm"
                        variant="outline"
                        onClick={onTest}
                    >
                        Probar Conexión
                    </Button>
                )}
                {onToggleActive && (
                    <Button
                        size="sm"
                        variant={integration.is_active ? 'outline' : 'primary'}
                        onClick={onToggleActive}
                    >
                        {integration.is_active ? 'Desactivar' : 'Activar'}
                    </Button>
                )}
            </div>
        </div>
    );
}
