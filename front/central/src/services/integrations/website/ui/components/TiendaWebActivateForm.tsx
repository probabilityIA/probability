'use client';

import { useState } from 'react';
import { Alert, Button } from '@/shared/ui';
import { GlobeAltIcon } from '@heroicons/react/24/outline';
import { IntegrationType } from '@/services/integrations/core/domain/types';
import { createIntegrationAction } from '@/services/integrations/core/infra/actions';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';
import { getActionError } from '@/shared/utils/action-result';

interface TiendaWebActivateFormProps {
    integrationType: IntegrationType;
    onSuccess: () => void;
    onBack: () => void;
}

export function TiendaWebActivateForm({ integrationType, onSuccess, onBack }: TiendaWebActivateFormProps) {
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
            const result = await createIntegrationAction({
                name: 'Sitio Web',
                code: 'tienda_web',
                integration_type_id: integrationType.id,
                category: integrationType.category?.code || integrationType.integration_category?.code || 'storefront',
                business_id: isSuperAdmin ? selectedBusinessId : null,
                is_active: true,
                is_default: true,
                config: {},
                credentials: {},
            });

            if (result.success) {
                onSuccess();
            } else {
                setError(result.message || 'Error al activar Sitio Web');
            }
        } catch (err: any) {
            setError(getActionError(err, 'Error al activar Sitio Web'));
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="max-w-md mx-auto py-8 space-y-6">
            <div className="flex flex-col items-center text-center">
                <div className="w-16 h-16 bg-emerald-100 rounded-full flex items-center justify-center mb-4">
                    <GlobeAltIcon className="w-8 h-8 text-emerald-600" />
                </div>
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white">Activar Sitio Web</h3>
                <p className="text-sm text-gray-500 dark:text-gray-400 mt-2">
                    Crea un sitio web público para tu negocio.
                    La configuración de secciones se gestiona desde Sitio Web.
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
                            className="w-full px-3 py-2 border border-blue-300 rounded-lg text-sm bg-white focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                        >
                            <option value="">— Selecciona un negocio —</option>
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
                    Activar Sitio Web
                </Button>
            </div>
        </div>
    );
}
