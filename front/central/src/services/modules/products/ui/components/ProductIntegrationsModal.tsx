'use client';

import { useState, useEffect } from 'react';
import { Modal } from '@/shared/ui/modal';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Select } from '@/shared/ui/select';
import { Alert } from '@/shared/ui/alert';
import { Spinner } from '@/shared/ui/spinner';
import { Badge } from '@/shared/ui/badge';
import { Product, ProductIntegration } from '../../domain/types';
import {
    getProductIntegrationsAction,
    addProductIntegrationAction,
    removeProductIntegrationAction
} from '../../infra/actions';
import { getIntegrationsAction } from '@/services/integrations/core/infra/actions';
import { Integration } from '@/services/integrations/core/domain/types';

interface ProductIntegrationsModalProps {
    product: Product;
    isOpen: boolean;
    onClose: () => void;
    onSuccess?: () => void;
}

export default function ProductIntegrationsModal({
    product,
    isOpen,
    onClose,
    onSuccess
}: ProductIntegrationsModalProps) {
    const [integrations, setIntegrations] = useState<ProductIntegration[]>([]);
    const [allIntegrations, setAllIntegrations] = useState<Integration[]>([]);
    const [availableIntegrations, setAvailableIntegrations] = useState<Integration[]>([]);
    const [loading, setLoading] = useState(false);
    const [loadingIntegrations, setLoadingIntegrations] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    // Form state for adding new integration
    const [selectedIntegrationId, setSelectedIntegrationId] = useState<string>('');
    const [externalProductId, setExternalProductId] = useState('');

    useEffect(() => {
        if (isOpen && product.id) {
            loadProductIntegrations();
            loadAvailableIntegrations();
        }
    }, [isOpen, product.id]);

    const loadProductIntegrations = async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getProductIntegrationsAction(product.id);
            if (response.success && response.data) {
                setIntegrations(response.data);
            } else {
                setError(response.message || 'Error al cargar las integraciones');
            }
        } catch (err: any) {
            setError(err.message || 'Error al cargar las integraciones');
        } finally {
            setLoading(false);
        }
    };

    const loadAvailableIntegrations = async (currentIntegrations?: ProductIntegration[]) => {
        setLoadingIntegrations(true);
        try {
            const response = await getIntegrationsAction({
                business_id: product.business_id,
                is_active: true,
                page_size: 100
            });
            if (response.success && response.data) {
                // Guardar todas las integraciones para mostrar nombres
                setAllIntegrations(response.data);
                
                // Filtrar integraciones que ya están asignadas para el selector
                const integrationsToCheck = currentIntegrations || integrations;
                const assignedIds = integrationsToCheck.map(i => i.integration_id);
                const available = response.data.filter(
                    (integration: Integration) => !assignedIds.includes(integration.id)
                );
                setAvailableIntegrations(available);
            }
        } catch (err: any) {
            console.error('Error al cargar integraciones disponibles:', err);
            // Si falla, al menos intentar cargar todas las integraciones
            setAllIntegrations([]);
            setAvailableIntegrations([]);
        } finally {
            setLoadingIntegrations(false);
        }
    };

    const handleAddIntegration = async () => {
        if (!selectedIntegrationId || !externalProductId.trim()) {
            setError('Por favor completa todos los campos');
            return;
        }

        setLoading(true);
        setError(null);
        setSuccess(null);

        try {
            const response = await addProductIntegrationAction(product.id, {
                integration_id: parseInt(selectedIntegrationId),
                external_product_id: externalProductId.trim()
            });

            if (response.success) {
                setSuccess('Integración agregada exitosamente');
                setSelectedIntegrationId('');
                setExternalProductId('');
                // Recargar integraciones del producto
                const updatedResponse = await getProductIntegrationsAction(product.id);
                if (updatedResponse.success && updatedResponse.data) {
                    setIntegrations(updatedResponse.data);
                    // Actualizar integraciones disponibles con la nueva lista
                    await loadAvailableIntegrations(updatedResponse.data);
                }
                if (onSuccess) onSuccess();
            } else {
                setError(response.message || 'Error al agregar la integración');
            }
        } catch (err: any) {
            setError(err.message || 'Error al agregar la integración');
        } finally {
            setLoading(false);
        }
    };

    const handleRemoveIntegration = async (integrationId: number) => {
        if (!confirm('¿Estás seguro de que deseas remover esta integración?')) return;

        setLoading(true);
        setError(null);
        setSuccess(null);

        try {
            const response = await removeProductIntegrationAction(product.id, integrationId);
            if (response.success) {
                setSuccess('Integración removida exitosamente');
                // Recargar integraciones del producto
                const updatedResponse = await getProductIntegrationsAction(product.id);
                if (updatedResponse.success && updatedResponse.data) {
                    setIntegrations(updatedResponse.data);
                    // Actualizar integraciones disponibles con la nueva lista
                    await loadAvailableIntegrations(updatedResponse.data);
                }
                if (onSuccess) onSuccess();
            } else {
                setError(response.message || 'Error al remover la integración');
            }
        } catch (err: any) {
            setError(err.message || 'Error al remover la integración');
        } finally {
            setLoading(false);
        }
    };

    const getIntegrationName = (integration: ProductIntegration): string => {
        // Primero intentar usar el nombre del backend si está disponible
        if (integration.integration_name) {
            return integration.integration_name;
        }
        // Si no, buscar en todas las integraciones cargadas
        const found = allIntegrations.find(i => i.id === integration.integration_id);
        return found?.name || `Integración #${integration.integration_id}`;
    };

    const getIntegrationType = (integration: ProductIntegration): string => {
        // Primero intentar usar el tipo del backend si está disponible
        if (integration.integration_type) {
            return integration.integration_type;
        }
        // Si no, buscar en todas las integraciones cargadas
        const found = allIntegrations.find(i => i.id === integration.integration_id);
        return found?.type || 'unknown';
    };

    return (
        <Modal
            isOpen={isOpen}
            onClose={onClose}
            title={`Gestionar Integraciones - ${product.name}`}
            size="lg"
        >
            <div className="space-y-6">
                {error && (
                    <Alert type="error" onClose={() => setError(null)}>
                        {error}
                    </Alert>
                )}

                {success && (
                    <Alert type="success" onClose={() => setSuccess(null)}>
                        {success}
                    </Alert>
                )}

                {/* Agregar nueva integración */}
                <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                    <h3 className="text-sm font-semibold text-gray-700 mb-3">Agregar Integración</h3>
                    <div className="space-y-3">
                        <Select
                            label="Integración"
                            value={selectedIntegrationId}
                            onChange={(e) => setSelectedIntegrationId(e.target.value)}
                            options={[
                                { label: 'Seleccionar integración...', value: '' },
                                ...(availableIntegrations.length > 0
                                    ? availableIntegrations.map(integration => ({
                                        label: `${integration.name} (${integration.type})`,
                                        value: String(integration.id)
                                    }))
                                    : [{ label: 'No hay integraciones disponibles', value: '', disabled: true }]
                                )
                            ]}
                            disabled={loading || loadingIntegrations || availableIntegrations.length === 0}
                        />
                        {availableIntegrations.length === 0 && !loadingIntegrations && (
                            <p className="text-xs text-gray-500">
                                Todas las integraciones activas ya están asignadas a este producto
                            </p>
                        )}
                        <Input
                            label="ID de Producto Externo"
                            value={externalProductId}
                            onChange={(e) => setExternalProductId(e.target.value)}
                            placeholder="ID del producto en la plataforma externa"
                            disabled={loading}
                        />
                        <Button
                            onClick={handleAddIntegration}
                            disabled={loading || !selectedIntegrationId || !externalProductId.trim()}
                            className="w-full"
                        >
                            {loading ? <Spinner size="sm" /> : 'Agregar Integración'}
                        </Button>
                    </div>
                </div>

                {/* Lista de integraciones actuales */}
                <div>
                    <h3 className="text-sm font-semibold text-gray-700 mb-3">
                        Integraciones Asignadas ({integrations.length})
                    </h3>
                    {loading && integrations.length === 0 ? (
                        <div className="text-center py-8">
                            <Spinner size="sm" />
                            <p className="text-sm text-gray-500 mt-2">Cargando integraciones...</p>
                        </div>
                    ) : integrations.length === 0 ? (
                        <div className="text-center py-8 text-gray-500">
                            <p className="text-sm">No hay integraciones asignadas a este producto</p>
                        </div>
                    ) : (
                        <div className="space-y-2">
                            {integrations.map((integration) => (
                                <div
                                    key={integration.id}
                                    className="flex items-center justify-between p-3 bg-white border border-gray-200 rounded-lg hover:bg-gray-50"
                                >
                                    <div className="flex-1">
                                        <div className="flex items-center gap-2 mb-1">
                                            <span className="text-sm font-medium text-gray-900">
                                                {getIntegrationName(integration)}
                                            </span>
                                            <Badge type="secondary" className="text-xs">
                                                {getIntegrationType(integration)}
                                            </Badge>
                                        </div>
                                        <div className="text-xs text-gray-500">
                                            ID Externo: <span className="font-mono">{integration.external_product_id}</span>
                                        </div>
                                        <div className="text-xs text-gray-400 mt-1">
                                            Asignado: {new Date(integration.created_at).toLocaleDateString('es-CO')}
                                        </div>
                                    </div>
                                    <Button
                                        variant="danger"
                                        size="sm"
                                        onClick={() => handleRemoveIntegration(integration.integration_id)}
                                        disabled={loading}
                                    >
                                        Remover
                                    </Button>
                                </div>
                            ))}
                        </div>
                    )}
                </div>

                {/* Botones de acción */}
                <div className="flex justify-end gap-2 pt-4 border-t border-gray-200">
                    <Button variant="secondary" onClick={onClose} disabled={loading}>
                        Cerrar
                    </Button>
                </div>
            </div>
        </Modal>
    );
}
