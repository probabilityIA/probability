'use client';

import { useState, FormEvent, useEffect } from 'react';
import { Button, Input, Alert, Select } from '@/shared/ui';
import { ShipitConfig } from '../../domain/types';
import { createIntegrationAction } from '@/services/integrations/core/infra/actions';
import { useToast } from '@/shared/providers/toast-provider';
import { getBusinessesSimpleAction } from '@/services/auth/business/infra/actions';
import { TokenStorage } from '@/shared/utils/token-storage';
import { getActionError } from '@/shared/utils/action-result';
import {
    Cog6ToothIcon,
    TruckIcon,
    InformationCircleIcon,
    ArrowLeftIcon,
    CheckBadgeIcon,
} from '@heroicons/react/24/outline';

interface ShipitConfigFormProps {
    onSuccess?: () => void;
    onCancel?: () => void;
}

export function ShipitConfigForm({ onSuccess, onCancel }: ShipitConfigFormProps) {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [isSuperAdmin, setIsSuperAdmin] = useState(false);
    const [businesses, setBusinesses] = useState<Array<{ id: number; name: string }>>([]);
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
    const [loadingBusinesses, setLoadingBusinesses] = useState(false);

    const [formData, setFormData] = useState({
        name: 'Shipit',
    });

    useEffect(() => {
        const checkUserAndLoadBusinesses = async () => {
            const permissions = TokenStorage.getPermissions();
            const isSuperUser = permissions?.is_super || false;
            setIsSuperAdmin(isSuperUser);

            if (isSuperUser) {
                setLoadingBusinesses(true);
                try {
                    const response = await getBusinessesSimpleAction();
                    if (response.success && response.data) {
                        setBusinesses(response.data);
                    }
                } catch (err) {
                    console.error('Error loading businesses:', err);
                    showToast('Error al cargar la lista de negocios', 'error');
                } finally {
                    setLoadingBusinesses(false);
                }
            } else {
                if (permissions?.business_id) {
                    setSelectedBusinessId(permissions.business_id);
                }
            }
        };

        checkUserAndLoadBusinesses();
    }, []);

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);

        try {
            if (isSuperAdmin && !selectedBusinessId) {
                showToast('Debes seleccionar un negocio', 'warning');
                setLoading(false);
                return;
            }

            const config: ShipitConfig = {
                use_platform_token: true,
            };

            const response = await createIntegrationAction({
                name: formData.name,
                code: `shipit_${Date.now()}`,
                integration_type_id: 34,
                category: 'shipping',
                business_id: isSuperAdmin ? selectedBusinessId : null,
                config: config as any,
                credentials: {} as any,
                is_active: true,
                is_default: false,
            });

            if (response.success) {
                showToast('Transportadora Shipit activada exitosamente', 'success');
                onSuccess?.();
            } else {
                throw new Error(response.message || 'Error al activar la transportadora');
            }
        } catch (err: any) {
            setError(getActionError(err, 'Error al activar la transportadora'));
            showToast('Error al activar Shipit', 'error');
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
                        Shipit - Logistica
                    </h2>
                </div>
                <p className="text-sm text-gray-600 dark:text-gray-300 ml-14">
                    Activa Shipit (Chile) para generar envios, cotizar tarifas entre couriers y hacer seguimiento desde Probability.
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
                            No necesitas credenciales: los envios se generan con la cuenta Shipit de Probability.
                            Solo activa la transportadora y empieza a generar guias.
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

                {isSuperAdmin && (
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">
                            Negocio <span className="text-red-500">*</span>
                        </label>
                        {loadingBusinesses ? (
                            <div className="flex items-center gap-2 p-3 bg-gray-50 dark:bg-gray-700 rounded-lg">
                                <svg className="animate-spin h-5 w-5 text-teal-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                                </svg>
                                <span className="text-sm text-gray-600 dark:text-gray-300">Cargando negocios...</span>
                            </div>
                        ) : (
                            <Select
                                value={selectedBusinessId?.toString() || ''}
                                onChange={(e) => setSelectedBusinessId(Number(e.target.value))}
                                options={[
                                    { value: '', label: '-- Selecciona un negocio --' },
                                    ...businesses.map((business) => ({
                                        value: business.id.toString(),
                                        label: business.name,
                                    })),
                                ]}
                                required
                                className="bg-white dark:bg-gray-800"
                            />
                        )}
                        <p className="text-xs text-gray-500 dark:text-gray-400 mt-1.5 flex items-start gap-1">
                            <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <span>Selecciona el negocio al que pertenecerá esta transportadora</span>
                        </p>
                    </div>
                )}
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
                            Activando...
                        </>
                    ) : (
                        <>
                            <TruckIcon className="w-5 h-5 mr-2" />
                            Activar Transportadora
                        </>
                    )}
                </Button>
            </div>
        </form>
    );
}
