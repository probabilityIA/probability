'use client';

import { useState, FormEvent, useEffect } from 'react';
import { Button, Input, Alert, Select } from '@/shared/ui';
import { WompiCredentials } from '../../domain/types';
import { updateIntegrationAction } from '@/services/integrations/core/infra/actions';
import { useToast } from '@/shared/providers/toast-provider';
import { getBusinessesSimpleAction } from '@/services/auth/business/infra/actions';
import { TokenStorage } from '@/shared/utils/token-storage';
import {
    KeyIcon,
    Cog6ToothIcon,
    CreditCardIcon,
    InformationCircleIcon,
    ArrowLeftIcon,
    EyeIcon,
    EyeSlashIcon,
} from '@heroicons/react/24/outline';

interface WompiEditFormProps {
    integrationId: number;
    initialData: {
        name: string;
        config: Record<string, any>;
        credentials?: WompiCredentials;
        business_id?: number | null;
        is_testing?: boolean;
        base_url_test?: string;
    };
    onSuccess?: () => void;
    onCancel?: () => void;
}

export function WompiEditForm({ integrationId, initialData, onSuccess, onCancel }: WompiEditFormProps) {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [showPrivateKey, setShowPrivateKey] = useState(false);

    const [isSuperAdmin, setIsSuperAdmin] = useState(false);
    const [businesses, setBusinesses] = useState<Array<{ id: number; name: string }>>([]);
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(initialData.business_id || null);
    const [loadingBusinesses, setLoadingBusinesses] = useState(false);

    const [formData, setFormData] = useState({
        name: initialData.name,
        private_key: initialData.credentials?.private_key || '',
        environment: initialData.credentials?.environment || 'sandbox' as 'sandbox' | 'production',
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
            }
        };

        checkUserAndLoadBusinesses();
    }, []);

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);

        try {
            const updateData: any = {
                name: formData.name,
                config: {},
                is_testing: formData.environment === 'sandbox',
            };

            if (formData.private_key) {
                const credentials: WompiCredentials = {
                    private_key: formData.private_key,
                    environment: formData.environment,
                };
                updateData.credentials = credentials;
            }

            const response = await updateIntegrationAction(integrationId, updateData);

            if (response.success) {
                showToast('Integración Wompi actualizada exitosamente', 'success');
                onSuccess?.();
            } else {
                throw new Error(response.message || 'Error al actualizar integración');
            }
        } catch (err: any) {
            setError(err.message || 'Error al actualizar integración');
            showToast('Error al actualizar integración Wompi', 'error');
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-8" autoComplete="off">
            <div className="border-b border-gray-200 pb-6">
                <div className="flex items-center gap-3 mb-2">
                    <div className="p-2 bg-green-50 rounded-lg">
                        <CreditCardIcon className="w-6 h-6 text-green-600" />
                    </div>
                    <h2 className="text-2xl font-bold text-gray-900">Editar Wompi Pagos</h2>
                </div>
                <p className="text-sm text-gray-600 ml-14">Actualiza la configuración de tu integración con Wompi.</p>
            </div>

            {error && <Alert type="error">{error}</Alert>}

            <div className="bg-gray-50 rounded-xl p-6 space-y-4">
                <div className="flex items-center gap-2 mb-4">
                    <Cog6ToothIcon className="w-5 h-5 text-gray-700" />
                    <h3 className="text-lg font-semibold text-gray-900">Configuración General</h3>
                </div>
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Nombre de la Integración <span className="text-red-500">*</span>
                    </label>
                    <Input type="text" value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })} required className="bg-white" />
                </div>
                {isSuperAdmin && (
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">Negocio</label>
                        {loadingBusinesses ? (
                            <div className="flex items-center gap-2 p-3 bg-gray-100 rounded-lg">
                                <svg className="animate-spin h-5 w-5 text-green-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle><path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
                                <span className="text-sm text-gray-600">Cargando negocios...</span>
                            </div>
                        ) : (
                            <Select value={selectedBusinessId?.toString() || ''} onChange={(e) => setSelectedBusinessId(Number(e.target.value))} options={[{ value: '', label: '-- Selecciona un negocio --' }, ...businesses.map((b) => ({ value: b.id.toString(), label: b.name }))]} disabled={true} className="bg-white" />
                        )}
                        <p className="text-xs text-gray-500 mt-1.5 flex items-start gap-1">
                            <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <span>El negocio no puede ser modificado después de la creación</span>
                        </p>
                    </div>
                )}
            </div>

            <div className="bg-gradient-to-br from-green-50 to-emerald-50 rounded-xl p-6 space-y-4 border border-green-100">
                <div className="flex items-center gap-2 mb-4">
                    <KeyIcon className="w-5 h-5 text-green-700" />
                    <h3 className="text-lg font-semibold text-gray-900">Credenciales API</h3>
                </div>
                <p className="text-sm text-green-900 -mt-2 mb-4 flex items-start gap-2">
                    <InformationCircleIcon className="w-5 h-5 flex-shrink-0 mt-0.5" />
                    <span>Deja el campo vacío para mantener la Private Key actual.</span>
                </p>
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">Ambiente</label>
                    <Select value={formData.environment} onChange={(e) => setFormData({ ...formData, environment: e.target.value as 'sandbox' | 'production' })} options={[{ value: 'sandbox', label: 'Sandbox (Pruebas)' }, { value: 'production', label: 'Producción' }]} className="bg-white" />
                </div>
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">Private Key</label>
                    <div className="relative">
                        <Input type={showPrivateKey ? 'text' : 'password'} value={formData.private_key} onChange={(e) => setFormData({ ...formData, private_key: e.target.value })} placeholder="Dejar vacío para mantener la actual" autoComplete="new-password" data-1p-ignore className="bg-white font-mono text-sm pr-10" />
                        <button type="button" onClick={() => setShowPrivateKey(!showPrivateKey)} className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none" tabIndex={-1}>
                            {showPrivateKey ? <EyeSlashIcon className="w-5 h-5" /> : <EyeIcon className="w-5 h-5" />}
                        </button>
                    </div>
                </div>
            </div>

            <div className="flex justify-between items-center gap-3 pt-6 border-t border-gray-200">
                {onCancel && (
                    <Button type="button" onClick={onCancel} disabled={loading} className="min-w-[140px] bg-gray-100 hover:bg-gray-200 text-gray-700 border border-gray-300">
                        <ArrowLeftIcon className="w-4 h-4 mr-2" />
                        Cancelar
                    </Button>
                )}
                <Button type="submit" variant="primary" disabled={loading} className="min-w-[200px] bg-green-600 hover:bg-green-700 text-white font-semibold">
                    {loading ? (
                        <><svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle><path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>Actualizando...</>
                    ) : (
                        <><CreditCardIcon className="w-5 h-5 mr-2" />Actualizar Integración</>
                    )}
                </Button>
            </div>
        </form>
    );
}
