'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Spinner } from '@/shared/ui/spinner';
import { Button } from '@/shared/ui/button';
import { useToast } from '@/shared/providers/toast-provider';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { InvoicingConfigForm } from '@/services/modules/invoicing/ui/components';
import { useInvoicingConfig } from '@/services/modules/invoicing/ui/hooks/useInvoicingConfig';
import { useIntegrations } from '@/services/integrations/core/ui/hooks/useIntegrations';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';

/**
 * Página de formulario para crear configuración de facturación
 */
export function InvoicingConfigFormPage() {
  const router = useRouter();
  const { showToast } = useToast();
  const { permissions } = usePermissions();
  const { providers, loading: providersLoading, fetchProviders } = useInvoicingConfig();
  const {
    integrations,
    loading: integrationsLoading,
    refresh: refreshIntegrations,
    setFilterCategory,
  } = useIntegrations();

  const [selectedIntegrationId, setSelectedIntegrationId] = useState<number | null>(null);

  const businessId = permissions?.business_id || 0;

  useEffect(() => {
    if (businessId) {
      fetchProviders('CO'); // Cargar proveedores de Colombia
      setFilterCategory('external'); // Filtrar solo integraciones externas
      refreshIntegrations();
    }
  }, [businessId]);

  const handleSuccess = () => {
    showToast('Configuración creada exitosamente', 'success');
    router.push('/invoicing/configs');
  };

  const handleCancel = () => {
    router.push('/invoicing/configs');
  };

  if (providersLoading || integrationsLoading) {
    return (
      <div className="flex justify-center items-center h-screen">
        <Spinner />
      </div>
    );
  }

  if (!businessId || businessId === 0) {
    return (
      <div className="p-8">
        <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-6">
          <p className="text-yellow-800">
            No se ha seleccionado un negocio. Por favor, seleccione un negocio para
            continuar.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="p-8 max-w-6xl mx-auto">
      {/* Header */}
      <div className="mb-6">
        <Button
          variant="secondary"
          size="sm"
          onClick={handleCancel}
          className="mb-4"
        >
          <ArrowLeftIcon className="w-4 h-4 mr-2" />
          Volver a Configuraciones
        </Button>

        <h1 className="text-3xl font-bold text-gray-900">
          Nueva Configuración de Facturación
        </h1>
        <p className="text-gray-600 mt-2">
          Configura la facturación automática para una integración específica
        </p>
      </div>

      {/* Selector de integración */}
      {!selectedIntegrationId && (
        <div className="bg-white rounded-lg shadow p-6 mb-6">
          <h2 className="text-lg font-semibold mb-4">
            Paso 1: Selecciona una integración
          </h2>

          {integrations.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              <p>No hay integraciones externas disponibles.</p>
              <p className="text-sm mt-1">
                Primero debes conectar una integración (Shopify, MercadoLibre, etc.)
              </p>
              <Button
                variant="primary"
                size="sm"
                onClick={() => router.push('/integrations')}
                className="mt-4"
              >
                Ir a Integraciones
              </Button>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {integrations.map((integration) => (
                <div
                  key={integration.id}
                  onClick={() => setSelectedIntegrationId(integration.id)}
                  className="border-2 border-gray-200 rounded-lg p-4 cursor-pointer hover:border-blue-500 hover:bg-blue-50 transition-all"
                >
                  <div className="flex items-center gap-3">
                    {integration.integration_type?.image_url && (
                      <img
                        src={integration.integration_type.image_url}
                        alt={integration.name}
                        className="w-10 h-10 object-contain"
                      />
                    )}
                    <div>
                      <h3 className="font-semibold text-gray-900">
                        {integration.name}
                      </h3>
                      <p className="text-xs text-gray-500">
                        {integration.integration_type?.name || integration.type}
                      </p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* Formulario de configuración */}
      {selectedIntegrationId && (
        <div className="bg-white rounded-lg shadow p-6">
          <div className="mb-4">
            <h2 className="text-lg font-semibold">
              Paso 2: Configura la facturación
            </h2>
            <p className="text-sm text-gray-600">
              Integración seleccionada:{' '}
              <span className="font-medium">
                {integrations.find((i) => i.id === selectedIntegrationId)?.name}
              </span>
            </p>
            <Button
              variant="secondary"
              size="sm"
              onClick={() => setSelectedIntegrationId(null)}
              className="mt-2"
            >
              Cambiar integración
            </Button>
          </div>

          <InvoicingConfigForm
            integrationId={selectedIntegrationId}
            businessId={businessId}
            providers={providers}
            onSuccess={handleSuccess}
            onCancel={handleCancel}
          />
        </div>
      )}

      {/* Proveedores disponibles */}
      {providers.length === 0 && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-6 mt-6">
          <h3 className="text-red-900 font-semibold mb-2">
            No hay proveedores de facturación disponibles
          </h3>
          <p className="text-red-800 text-sm">
            Contacte con el administrador del sistema para configurar proveedores
            de facturación electrónica.
          </p>
        </div>
      )}
    </div>
  );
}
