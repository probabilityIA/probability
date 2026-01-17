'use client';

import { useState } from 'react';
import { Input, Button, Alert } from '@/shared/ui';
import { TokenStorage } from '@/shared/config';

interface ShopifyOAuthFormProps {
    onCancel?: () => void;
    onSubmit?: (data: any) => void;
    onTestConnection?: (config: any, credentials: any) => Promise<boolean>;
    onGetWebhook?: () => Promise<any>;
    initialData?: {
        name?: string;
        code?: string;
        store_id?: string;
        config?: any;
        credentials?: any;
        business_id?: number | null;
    };
    isEdit?: boolean;
    integrationId?: number;
}

export default function ShopifyOAuthForm({
    onCancel,
    onSubmit,
    onTestConnection,
    onGetWebhook,
    initialData,
    isEdit,
    integrationId,
}: ShopifyOAuthFormProps) {
    const [formData, setFormData] = useState({
        name: initialData?.name || '',
        shop_domain: initialData?.store_id || '',
    });

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleConnectShopify = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!formData.name || !formData.shop_domain) {
            setError('Por favor completa todos los campos');
            return;
        }

        const token = TokenStorage.getSessionToken();

        if (!token) {
            setError('No se encontró token de autenticación. Por favor, inicia sesión de nuevo.');
            return;
        }

        setLoading(true);
        setError(null);

        try {
            // Llamar al backend para iniciar el flujo OAuth
            const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:3050/api/v1';
            const response = await fetch(`${apiBaseUrl}/integrations/shopify/connect`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    shop_domain: formData.shop_domain,
                    integration_name: formData.name
                })
            });

            const data = await response.json();

            if (!response.ok || !data.success) {
                throw new Error(data.error || data.message || 'Error al iniciar OAuth');
            }

            // Redirigir al usuario a Shopify para autorización
            if (data.authorization_url) {
                window.location.href = data.authorization_url;
            } else {
                throw new Error('No se recibió URL de autorización');
            }
        } catch (err: any) {
            console.error('Error al conectar con Shopify:', err);
            setError(err.message || 'Error al conectar con Shopify');
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleConnectShopify} className="space-y-6 w-full">
            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <div className="p-6 rounded-lg border border-gray-200 bg-white">
                <div className="space-y-4">
                    <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-4">
                        <div className="flex items-start gap-3">
                            <span className="text-2xl">ℹ️</span>
                            <div>
                                <p className="text-sm font-medium text-blue-900 mb-1">
                                    Conexión OAuth con Shopify
                                </p>
                                <p className="text-xs text-blue-700">
                                    Serás redirigido a Shopify para autorizar la conexión.
                                    No necesitas ingresar tokens manualmente.
                                </p>
                            </div>
                        </div>
                    </div>

                    {/* Nombre de la Integración */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Nombre de la Integración *
                        </label>
                        <Input
                            type="text"
                            required
                            placeholder="Ej: Tienda Principal"
                            value={formData.name}
                            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                            className="w-full"
                        />
                        <p className="text-xs text-gray-500 mt-1">
                            Un nombre descriptivo para identificar esta integración
                        </p>
                    </div>

                    {/* Store Domain */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Dominio de la Tienda *
                        </label>
                        <Input
                            type="text"
                            required
                            placeholder="mystore.myshopify.com"
                            value={formData.shop_domain}
                            onChange={(e) => setFormData({ ...formData, shop_domain: e.target.value })}
                            className="w-full"
                        />
                        <p className="text-xs text-gray-500 mt-1">
                            El dominio de tu tienda Shopify (ej: mitienda.myshopify.com)
                        </p>
                    </div>
                </div>
            </div>

            {/* Action Buttons */}
            <div className="flex flex-row justify-end gap-3 pt-4 border-t">
                {onCancel && (
                    <Button
                        type="button"
                        onClick={onCancel}
                        variant="outline"
                        disabled={loading}
                    >
                        Cancelar
                    </Button>
                )}
                <Button
                    type="submit"
                    disabled={loading || !formData.name || !formData.shop_domain}
                    loading={loading}
                    variant="primary"
                >
                    {loading ? 'Conectando...' : '🔗 Conectar con Shopify'}
                </Button>
            </div>
        </form>
    );
}
