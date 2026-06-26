'use client';

import { useState } from 'react';
import { Alert, Button } from '@/shared/ui';
import { CubeIcon } from '@heroicons/react/24/outline';
import { IntegrationType } from '@/services/integrations/core/domain/types';
import { createIntegrationAction } from '@/services/integrations/core/infra/actions';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';
import { getActionError } from '@/shared/utils/action-result';

interface ModuleActivateFormProps {
    integrationType: IntegrationType;
    onSuccess: () => void;
    onBack: () => void;
}

export function ModuleActivateForm({ integrationType, onSuccess, onBack }: ModuleActivateFormProps) {
    const { isSuperAdmin } = usePermissions();
    const { businesses, loading: loadingBusinesses } = useBusinessesSimple();
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const requiresBusinessSelection = isSuperAdmin && selectedBusinessId === null;

    const handleActivate = async () => {
        if (isSuperAdmin && !selectedBusinessId) return;

        setLoading(true);
        setError(null);

        try {
            const businessId = isSuperAdmin ? selectedBusinessId : null;
            const code = businessId
                ? `${integrationType.code}_business_${businessId}`
                : `${integrationType.code}_${Date.now()}`;

            const result = await createIntegrationAction({
                name: integrationType.name,
                code,
                integration_type_id: integrationType.id,
                category: integrationType.category?.code || integrationType.integration_category?.code || 'internal',
                business_id: businessId,
                is_active: true,
                is_default: true,
                config: { use_platform_token: true },
                credentials: {},
            });

            if (result.success) {
                onSuccess();
            } else {
                setError(result.message || `Error al activar ${integrationType.name}`);
            }
        } catch (err: any) {
            setError(getActionError(err, `Error al activar ${integrationType.name}`));
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="max-w-md mx-auto py-8 space-y-6">
            <div className="flex flex-col items-center text-center">
                {integrationType.image_url ? (
                    <img
                        src={integrationType.image_url}
                        alt={integrationType.name}
                        className="w-16 h-16 object-contain rounded-lg shadow-md mb-4"
                    />
                ) : (
                    <div className="w-16 h-16 bg-indigo-100 rounded-full flex items-center justify-center mb-4">
                        <CubeIcon className="w-8 h-8 text-indigo-600" />
                    </div>
                )}
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white">Activar {integrationType.name}</h3>
                <p className="text-sm text-gray-500 dark:text-gray-400 dark:text-gray-400 mt-2">
                    Habilita el modulo de {integrationType.name} para este negocio.
                    Este es un modulo interno de la plataforma, no requiere credenciales externas.
                </p>
            </div>

            {isSuperAdmin && (
                <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                    <label className="block text-sm font-medium text-blue-800 mb-2">
                        Selecciona un negocio *
                    </label>
                    {loadingBusinesses ? (
                        <p className="text-sm text-blue-600">Cargando negocios...</p>
                    ) : (
                        <select
                            value={selectedBusinessId?.toString() ?? ''}
                            onChange={(e) => setSelectedBusinessId(e.target.value ? Number(e.target.value) : null)}
                            className="w-full px-3 py-2 border border-blue-300 rounded-lg text-sm bg-white dark:bg-gray-800 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                        >
                            <option value="">Selecciona un negocio</option>
                            {businesses.map(b => (
                                <option key={b.id} value={b.id}>{b.name} (ID: {b.id})</option>
                            ))}
                        </select>
                    )}
                </div>
            )}

            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <div className="flex gap-3">
                <Button
                    type="button"
                    variant="outline"
                    onClick={onBack}
                    className="flex-1"
                >
                    Cancelar
                </Button>
                <Button
                    type="button"
                    variant="primary"
                    onClick={handleActivate}
                    disabled={loading || requiresBusinessSelection}
                    loading={loading}
                    className="flex-1"
                >
                    Activar {integrationType.name}
                </Button>
            </div>
        </div>
    );
}
