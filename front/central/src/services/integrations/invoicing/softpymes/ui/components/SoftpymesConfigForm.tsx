'use client';

import { useState, FormEvent } from 'react';
import { Button, Input, Alert } from '@/shared/ui';
import { SoftpymesConfig, SoftpymesCredentials } from '../../domain/types';
import { createIntegrationAction } from '@/services/integrations/core/infra/actions';
import { useToast } from '@/shared/providers/toast-provider';

interface SoftpymesConfigFormProps {
    onSuccess?: () => void;
    onCancel?: () => void;
}

export function SoftpymesConfigForm({ onSuccess, onCancel }: SoftpymesConfigFormProps) {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const [formData, setFormData] = useState({
        name: '',
        company_nit: '',
        company_name: '',
        api_url: 'https://api.softpymes.com',
        test_mode: false,
        api_key: '',
        api_secret: '',
    });

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);

        try {
            const config: SoftpymesConfig = {
                company_nit: formData.company_nit,
                company_name: formData.company_name,
                api_url: formData.api_url,
                test_mode: formData.test_mode,
            };

            const credentials: SoftpymesCredentials = {
                api_key: formData.api_key,
                api_secret: formData.api_secret,
            };

            // Create integration
            // Note: integration_type_id should be the ID for Softpymes from backend
            // You may need to fetch this or hardcode it based on your backend setup
            const response = await createIntegrationAction({
                name: formData.name,
                code: `softpymes_${Date.now()}`,
                integration_type_id: 7, // TODO: Get this from backend or config
                category: 'invoicing',
                business_id: null, // Will be set by backend from JWT
                config: config as any,
                credentials: credentials as any,
                is_active: true,
                is_default: false,
            });

            if (response.success) {
                showToast('Integración Softpymes creada exitosamente', 'success');
                onSuccess?.();
            } else {
                throw new Error(response.message || 'Error al crear integración');
            }
        } catch (err: any) {
            setError(err.message || 'Error al crear integración');
            showToast('Error al crear integración Softpymes', 'error');
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-6 p-6">
            <div>
                <h3 className="text-lg font-semibold text-gray-900 mb-4">
                    Configurar Softpymes
                </h3>
                <p className="text-sm text-gray-600 mb-6">
                    Ingresa tus credenciales de Softpymes para conectar la integración de facturación electrónica.
                </p>
            </div>

            {error && (
                <Alert type="error">
                    {error}
                </Alert>
            )}

            {/* Nombre de la integración */}
            <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                    Nombre de la integración *
                </label>
                <Input
                    type="text"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    placeholder="Softpymes - Producción"
                    required
                />
                <p className="text-xs text-gray-500 mt-1">
                    Nombre descriptivo para identificar esta integración
                </p>
            </div>

            {/* Información de la empresa */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        NIT de la empresa *
                    </label>
                    <Input
                        type="text"
                        value={formData.company_nit}
                        onChange={(e) => setFormData({ ...formData, company_nit: e.target.value })}
                        placeholder="900123456-7"
                        required
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Nombre de la empresa *
                    </label>
                    <Input
                        type="text"
                        value={formData.company_name}
                        onChange={(e) => setFormData({ ...formData, company_name: e.target.value })}
                        placeholder="Mi Empresa SAS"
                        required
                    />
                </div>
            </div>

            {/* Credenciales */}
            <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                    API Key *
                </label>
                <Input
                    type="password"
                    value={formData.api_key}
                    onChange={(e) => setFormData({ ...formData, api_key: e.target.value })}
                    placeholder="sk_live_..."
                    required
                />
                <p className="text-xs text-gray-500 mt-1">
                    Encuentra tu API Key en el panel de Softpymes
                </p>
            </div>

            <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                    API Secret *
                </label>
                <Input
                    type="password"
                    value={formData.api_secret}
                    onChange={(e) => setFormData({ ...formData, api_secret: e.target.value })}
                    placeholder="********"
                    required
                />
            </div>

            {/* API URL */}
            <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                    URL de la API
                </label>
                <Input
                    type="url"
                    value={formData.api_url}
                    onChange={(e) => setFormData({ ...formData, api_url: e.target.value })}
                    placeholder="https://api.softpymes.com"
                />
                <p className="text-xs text-gray-500 mt-1">
                    URL base de la API de Softpymes
                </p>
            </div>

            {/* Test Mode */}
            <div className="flex items-center">
                <input
                    type="checkbox"
                    id="test_mode"
                    checked={formData.test_mode}
                    onChange={(e) => setFormData({ ...formData, test_mode: e.target.checked })}
                    className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                />
                <label htmlFor="test_mode" className="ml-2 block text-sm text-gray-700">
                    Modo de pruebas
                </label>
            </div>

            {/* Buttons */}
            <div className="flex justify-end gap-3 pt-4 border-t">
                {onCancel && (
                    <Button
                        type="button"
                        variant="outline"
                        onClick={onCancel}
                        disabled={loading}
                    >
                        Cancelar
                    </Button>
                )}
                <Button
                    type="submit"
                    variant="primary"
                    disabled={loading}
                >
                    {loading ? 'Conectando...' : 'Conectar Softpymes'}
                </Button>
            </div>
        </form>
    );
}
