'use client';

import { useState, FormEvent, useEffect } from 'react';
import { Button, Input, Alert, Select } from '@/shared/ui';
import { PayUConfig, PayUCredentials } from '../../domain/types';
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

interface PayUEditFormProps {
    integrationId: number;
    initialData: {
        name: string;
        config: PayUConfig;
        credentials?: PayUCredentials;
        business_id?: number | null;
        is_testing?: boolean;
        base_url_test?: string;
    };
    onSuccess?: () => void;
    onCancel?: () => void;
}

export function PayUEditForm({ integrationId, initialData, onSuccess, onCancel }: PayUEditFormProps) {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [showApiKey, setShowApiKey] = useState(false);

    const [isSuperAdmin, setIsSuperAdmin] = useState(false);
    const [businesses, setBusinesses] = useState<Array<{ id: number; name: string }>>([]);
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(initialData.business_id || null);
    const [loadingBusinesses, setLoadingBusinesses] = useState(false);

    const [formData, setFormData] = useState({
        name: initialData.name,
        api_key: initialData.credentials?.api_key || '',
        api_login: initialData.credentials?.api_login || '',
        account_id: initialData.config.account_id || '',
        merchant_id: initialData.config.merchant_id || '',
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
            const config: PayUConfig = {
                account_id: formData.account_id || undefined,
                merchant_id: formData.merchant_id || undefined,
            };

            const updateData: any = {
                name: formData.name,
                config: config,
                is_testing: formData.environment === 'sandbox',
            };

            if (formData.api_key && formData.api_login) {
                const credentials: PayUCredentials = {
                    api_key: formData.api_key,
                    api_login: formData.api_login,
                    environment: formData.environment,
                };
                updateData.credentials = credentials;
            }

            const response = await updateIntegrationAction(integrationId, updateData);

            if (response.success) {
                showToast('Integración PayU actualizada exitosamente', 'success');
                onSuccess?.();
            } else {
                throw new Error(response.message || 'Error al actualizar integración');
            }
        } catch (err: any) {
            setError(err.message || 'Error al actualizar integración');
            showToast('Error al actualizar integración PayU', 'error');
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-8" autoComplete="off">
            <div className="border-b border-gray-200 pb-6">
                <div className="flex items-center gap-3 mb-2">
                    <div className="p-2 bg-orange-50 rounded-lg">
                        <CreditCardIcon className="w-6 h-6 text-orange-600" />
                    </div>
                    <h2 className="text-2xl font-bold text-gray-900">Editar PayU Pagos</h2>
                </div>
                <p className="text-sm text-gray-600 ml-14">Actualiza la configuración de tu integración con PayU.</p>
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
                                <svg className="animate-spin h-5 w-5 text-orange-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle><path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
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
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">Account ID</label>
                        <Input type="text" value={formData.account_id} onChange={(e) => setFormData({ ...formData, account_id: e.target.value })} placeholder="ID de tu cuenta PayU" className="bg-white" />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">Merchant ID</label>
                        <Input type="text" value={formData.merchant_id} onChange={(e) => setFormData({ ...formData, merchant_id: e.target.value })} placeholder="ID de tu comercio PayU" className="bg-white" />
                    </div>
                </div>
            </div>

            <div className="bg-gradient-to-br from-orange-50 to-amber-50 rounded-xl p-6 space-y-4 border border-orange-100">
                <div className="flex items-center gap-2 mb-4">
                    <KeyIcon className="w-5 h-5 text-orange-700" />
                    <h3 className="text-lg font-semibold text-gray-900">Credenciales API</h3>
                </div>
                <p className="text-sm text-orange-900 -mt-2 mb-4 flex items-start gap-2">
                    <InformationCircleIcon className="w-5 h-5 flex-shrink-0 mt-0.5" />
                    <span>Deja los campos vacíos para mantener las credenciales actuales.</span>
                </p>
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">Ambiente</label>
                    <Select value={formData.environment} onChange={(e) => setFormData({ ...formData, environment: e.target.value as 'sandbox' | 'production' })} options={[{ value: 'sandbox', label: 'Sandbox (Pruebas)' }, { value: 'production', label: 'Producción' }]} className="bg-white" />
                </div>
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">API Login</label>
                    <Input type="text" value={formData.api_login} onChange={(e) => setFormData({ ...formData, api_login: e.target.value })} placeholder="Dejar vacío para mantener el actual" autoComplete="off" data-1p-ignore className="bg-white font-mono text-sm" />
                </div>
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">API Key</label>
                    <div className="relative">
                        <Input type={showApiKey ? 'text' : 'password'} value={formData.api_key} onChange={(e) => setFormData({ ...formData, api_key: e.target.value })} placeholder="Dejar vacío para mantener la actual" autoComplete="new-password" data-1p-ignore className="bg-white font-mono text-sm pr-10" />
                        <button type="button" onClick={() => setShowApiKey(!showApiKey)} className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none" tabIndex={-1}>
                            {showApiKey ? <EyeSlashIcon className="w-5 h-5" /> : <EyeIcon className="w-5 h-5" />}
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
                <Button type="submit" variant="primary" disabled={loading} className="min-w-[200px] bg-orange-600 hover:bg-orange-700 text-white font-semibold">
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
