'use client';

import { useState, FormEvent } from 'react';
import { Button, Input, Alert } from '@/shared/ui';
import { ShipitConfig, ShipitCredentials } from '../../domain/types';
import { updateIntegrationAction } from '@/services/integrations/core/infra/actions';
import { useToast } from '@/shared/providers/toast-provider';
import { getActionError } from '@/shared/utils/action-result';
import {
    Cog6ToothIcon,
    TruckIcon,
    InformationCircleIcon,
    ArrowLeftIcon,
    CheckBadgeIcon,
} from '@heroicons/react/24/outline';

interface ShipitEditFormProps {
    integrationId: number;
    initialData: {
        name: string;
        config: ShipitConfig;
        credentials?: ShipitCredentials;
        business_id?: number | null;
    };
    onSuccess?: () => void;
    onCancel?: () => void;
}

export function ShipitEditForm({ integrationId, initialData, onSuccess, onCancel }: ShipitEditFormProps) {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const [formData, setFormData] = useState({
        name: initialData.name,
    });

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);

        try {
            const config: ShipitConfig = {
                ...initialData.config,
                use_platform_token: true,
            };

            const response = await updateIntegrationAction(integrationId, {
                name: formData.name,
                config: config,
            } as any);

            if (response.success) {
                showToast('Integracion Shipit actualizada exitosamente', 'success');
                onSuccess?.();
            } else {
                throw new Error(response.message || 'Error al actualizar integración');
            }
        } catch (err: any) {
            setError(getActionError(err, 'Error al actualizar integración'));
            showToast('Error al actualizar integracion Shipit', 'error');
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-8" autoComplete="off">
            <div className="border-b border-gray-200 dark:border-gray-700 pb-6">
                <div className="flex items-center gap-3 mb-2">
                    <div className="p-2 bg-teal-50 rounded-lg">
                        <TruckIcon className="w-6 h-6 text-teal-600" />
                    </div>
                    <h2 className="text-2xl font-bold text-gray-900 dark:text-white">
                        Editar Shipit - Logistica
                    </h2>
                </div>
                <p className="text-sm text-gray-600 dark:text-gray-300 ml-14">
                    Actualiza la configuracion de tu transportadora Shipit.
                </p>
            </div>

            {error && (
                <Alert type="error">
                    {error}
                </Alert>
            )}

            <div className="bg-gradient-to-br from-teal-50 to-cyan-50 rounded-xl p-6 border border-teal-100">
                <div className="flex items-start gap-3">
                    <CheckBadgeIcon className="w-6 h-6 text-teal-700 flex-shrink-0 mt-0.5" />
                    <div>
                        <h3 className="text-sm font-semibold text-gray-900 mb-1">
                            Conexion gestionada por la plataforma
                        </h3>
                        <p className="text-sm text-teal-900">
                            Los envios se generan con la cuenta Shipit de Probability, no necesitas credenciales propias.
                            Puedes activar o desactivar la transportadora desde el listado de integraciones.
                        </p>
                    </div>
                </div>
            </div>

            <div className="bg-gray-50 dark:bg-gray-700 rounded-xl p-6 space-y-4">
                <div className="flex items-center gap-2 mb-4">
                    <Cog6ToothIcon className="w-5 h-5 text-gray-700 dark:text-gray-200" />
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                        Configuracion General
                    </h3>
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">
                        Nombre de la Integracion <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="text"
                        value={formData.name}
                        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                        placeholder="Ej: Shipit"
                        required
                        className="bg-white dark:bg-gray-800"
                    />
                    <p className="text-xs text-gray-500 dark:text-gray-400 mt-1.5 flex items-start gap-1">
                        <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                        <span>Nombre descriptivo para identificar esta transportadora en el sistema</span>
                    </p>
                </div>
            </div>

            <div className="flex justify-between items-center gap-3 pt-6 border-t border-gray-200 dark:border-gray-700">
                {onCancel && (
                    <Button
                        type="button"
                        onClick={onCancel}
                        disabled={loading}
                        className="min-w-[140px] bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 text-gray-700 dark:text-gray-200 border border-gray-300 dark:border-gray-600"
                    >
                        <ArrowLeftIcon className="w-4 h-4 mr-2" />
                        Cancelar
                    </Button>
                )}
                <Button
                    type="submit"
                    variant="primary"
                    disabled={loading}
                    className="min-w-[200px] bg-teal-600 hover:bg-teal-700 text-white font-semibold"
                >
                    {loading ? (
                        <>
                            <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                            </svg>
                            Actualizando...
                        </>
                    ) : (
                        <>
                            <TruckIcon className="w-5 h-5 mr-2" />
                            Actualizar Integracion
                        </>
                    )}
                </Button>
            </div>
        </form>
    );
}
