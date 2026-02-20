'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Spinner } from '@/shared/ui/spinner';
import { Button } from '@/shared/ui/button';
import { useToast } from '@/shared/providers/toast-provider';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { InvoicingConfigForm } from '@/services/modules/invoicing/ui/components';
import { useIntegrations } from '@/services/integrations/core/ui/hooks/useIntegrations';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';

/**
 * Página de formulario para crear configuración de facturación
 */
export function InvoicingConfigFormPage() {
  const router = useRouter();
  const { showToast } = useToast();
  const { permissions } = usePermissions();
  const {
    integrations,
    loading: integrationsLoading,
    refresh: refreshIntegrations,
    setFilterCategory,
  } = useIntegrations();
  const { businesses, loading: businessesLoading } = useBusinessesSimple();

  const [selectedInvoicingIntegrationId, setSelectedInvoicingIntegrationId] = useState<number | null>(null);
  const [selectedEcommerceIntegrationIds, setSelectedEcommerceIntegrationIds] = useState<number[]>([]);
  const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
  const [currentCategory, setCurrentCategory] = useState<string>('invoicing');

  const businessId = permissions?.business_id || 0;
  const isSuperAdmin = !businessId || businessId === 0;

  // Para super admin, usar el negocio seleccionado; para usuario normal, usar su businessId
  const effectiveBusinessId = isSuperAdmin ? selectedBusinessId : businessId;

  useEffect(() => {
    if (effectiveBusinessId) {
      setFilterCategory(currentCategory); // Filtrar por categoría actual (ecommerce o invoicing)
      refreshIntegrations();
    }
  }, [effectiveBusinessId, currentCategory]);

  const handleSuccess = () => {
    const count = selectedEcommerceIntegrationIds.length;
    showToast(
      count > 1
        ? `${count} configuraciones creadas exitosamente`
        : 'Configuración creada exitosamente',
      'success'
    );
    router.push('/invoicing/configs');
  };

  const handleCancel = () => {
    router.push('/invoicing/configs');
  };

  if (integrationsLoading || (isSuperAdmin && businessesLoading)) {
    return (
      <div className="flex justify-center items-center h-screen">
        <Spinner />
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

      {/* Selector de negocio (solo para super admin) */}
      {isSuperAdmin && !selectedBusinessId && (
        <div className="bg-white rounded-lg shadow p-6 mb-6">
          <h2 className="text-lg font-semibold mb-4">
            Paso 1: Selecciona un negocio
          </h2>

          {businesses.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              <p>No hay negocios disponibles.</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {businesses.map((business) => (
                <div
                  key={business.id}
                  onClick={() => setSelectedBusinessId(business.id)}
                  className="border-2 border-gray-200 rounded-lg p-4 cursor-pointer hover:border-blue-500 hover:bg-blue-50 transition-all"
                >
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center">
                      <span className="text-blue-600 font-semibold text-lg">
                        {business.name.charAt(0).toUpperCase()}
                      </span>
                    </div>
                    <div>
                      <h3 className="font-semibold text-gray-900">
                        {business.name}
                      </h3>
                      <p className="text-xs text-gray-500">
                        ID: {business.id}
                      </p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* Selector de proveedor de facturación (PRIMERO) */}
      {effectiveBusinessId && !selectedInvoicingIntegrationId && (
        <div className="bg-white rounded-lg shadow p-6 mb-6">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h2 className="text-lg font-semibold">
                {isSuperAdmin ? 'Paso 2' : 'Paso 1'}: Selecciona un proveedor de facturación
              </h2>
              <p className="text-sm text-gray-600 mt-1">
                Selecciona el proveedor que emitirá las facturas electrónicas
              </p>
            </div>
            {isSuperAdmin && selectedBusinessId && (
              <Button
                variant="secondary"
                size="sm"
                onClick={() => {
                  setSelectedBusinessId(null);
                  setSelectedInvoicingIntegrationId(null);
                }}
              >
                Cambiar negocio
              </Button>
            )}
          </div>

          {integrations.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              <p>No hay proveedores de facturación disponibles.</p>
              <p className="text-sm mt-1">
                Primero debes conectar un proveedor de facturación electrónica (Softpymes, Siigo, etc.)
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
                  onClick={() => {
                    setSelectedInvoicingIntegrationId(integration.id);
                    setCurrentCategory('ecommerce,platform'); // Cambiar a e-commerce + plataforma para el siguiente paso
                  }}
                  className="border-2 border-gray-200 rounded-lg p-4 cursor-pointer hover:border-green-500 hover:bg-green-50 transition-all"
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

      {/* Selector de integración de E-commerce (SEGUNDO - MÚLTIPLE) */}
      {effectiveBusinessId && selectedInvoicingIntegrationId && selectedEcommerceIntegrationIds.length === 0 && (
        <div className="bg-white rounded-lg shadow p-6 mb-6">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h2 className="text-lg font-semibold">
                {isSuperAdmin ? 'Paso 3' : 'Paso 2'}: Selecciona origen de órdenes
              </h2>
              <p className="text-sm text-gray-600 mt-1">
                Selecciona una o más fuentes de órdenes cuyas órdenes se facturarán automáticamente
              </p>
            </div>
            <Button
              variant="secondary"
              size="sm"
              onClick={() => {
                setSelectedInvoicingIntegrationId(null);
                setSelectedEcommerceIntegrationIds([]);
                setCurrentCategory('invoicing');
              }}
            >
              Cambiar proveedor
            </Button>
          </div>

          {integrations.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              <p>No hay fuentes de órdenes disponibles.</p>
              <p className="text-sm mt-1">
                Primero debes conectar una tienda (Shopify, MercadoLibre, etc.)
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
            <>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                {integrations.map((integration) => {
                  const isSelected = selectedEcommerceIntegrationIds.includes(integration.id);

                  return (
                    <label
                      key={integration.id}
                      className={`border-2 rounded-lg p-4 cursor-pointer transition-all ${
                        isSelected
                          ? 'border-blue-500 bg-blue-50'
                          : 'border-gray-200 hover:border-blue-300 hover:bg-blue-25'
                      }`}
                    >
                      <div className="flex items-center gap-3">
                        <input
                          type="checkbox"
                          checked={isSelected}
                          onChange={(e) => {
                            if (e.target.checked) {
                              setSelectedEcommerceIntegrationIds([...selectedEcommerceIntegrationIds, integration.id]);
                            } else {
                              setSelectedEcommerceIntegrationIds(
                                selectedEcommerceIntegrationIds.filter((id) => id !== integration.id)
                              );
                            }
                          }}
                          className="w-5 h-5 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                        />
                        {integration.integration_type?.image_url ? (
                          <img
                            src={integration.integration_type.image_url}
                            alt={integration.name}
                            className="w-10 h-10 object-contain"
                          />
                        ) : (
                          <div className="w-10 h-10 bg-indigo-100 rounded-full flex items-center justify-center">
                            <span className="text-indigo-600 text-lg font-bold">P</span>
                          </div>
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
                    </label>
                  );
                })}
              </div>

              {/* Botón para continuar */}
              <div className="flex justify-end">
                <Button
                  variant="primary"
                  size="md"
                  onClick={() => {
                    if (selectedEcommerceIntegrationIds.length === 0) {
                      showToast('Debes seleccionar al menos una tienda', 'error');
                      return;
                    }
                    // Avanzar al siguiente paso
                    setCurrentCategory('ecommerce'); // Trigger para mostrar el formulario
                  }}
                  disabled={selectedEcommerceIntegrationIds.length === 0}
                >
                  Continuar ({selectedEcommerceIntegrationIds.length} tienda{selectedEcommerceIntegrationIds.length > 1 ? 's' : ''} seleccionada{selectedEcommerceIntegrationIds.length > 1 ? 's' : ''})
                </Button>
              </div>
            </>
          )}
        </div>
      )}

      {/* Formulario de configuración */}
      {effectiveBusinessId && selectedInvoicingIntegrationId && selectedEcommerceIntegrationIds.length > 0 && (
        <div className="bg-white rounded-lg shadow p-6">
          <div className="mb-6">
            <h2 className="text-lg font-semibold mb-4">
              {isSuperAdmin ? 'Paso 4' : 'Paso 3'}: Configura la facturación automática
            </h2>

            <div className="mb-4">
              <div className="bg-green-50 border border-green-200 rounded-lg p-4 mb-4">
                <p className="text-xs text-green-600 font-medium mb-1">PROVEEDOR DE FACTURACIÓN</p>
                <p className="font-semibold text-gray-900">
                  {integrations.find((i) => i.id === selectedInvoicingIntegrationId)?.name}
                </p>
                <p className="text-xs text-gray-600 mt-1">
                  Emitirá las facturas electrónicas
                </p>
              </div>

              <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                <p className="text-xs text-blue-600 font-medium mb-2">
                  FUENTES DE ÓRDENES ({selectedEcommerceIntegrationIds.length})
                </p>
                <div className="space-y-2">
                  {selectedEcommerceIntegrationIds.map((id) => {
                    const integration = integrations.find((i) => i.id === id);
                    return (
                      <div key={id} className="flex items-center gap-2 text-sm">
                        <span className="w-2 h-2 bg-blue-600 rounded-full"></span>
                        <span className="font-medium text-gray-900">{integration?.name}</span>
                        <span className="text-gray-500">({integration?.integration_type?.name})</span>
                      </div>
                    );
                  })}
                </div>
                <p className="text-xs text-gray-600 mt-2">
                  Las órdenes de estas tiendas se facturarán automáticamente
                </p>
              </div>
            </div>

            <Button
              variant="secondary"
              size="sm"
              onClick={() => {
                setSelectedEcommerceIntegrationIds([]);
                setCurrentCategory('ecommerce,platform');
              }}
              className="mr-2"
            >
              Cambiar fuentes
            </Button>
            <Button
              variant="secondary"
              size="sm"
              onClick={() => {
                setSelectedInvoicingIntegrationId(null);
                setSelectedEcommerceIntegrationIds([]);
                setCurrentCategory('invoicing');
              }}
            >
              Cambiar proveedor
            </Button>
          </div>

          <InvoicingConfigForm
            integrationIds={selectedEcommerceIntegrationIds}
            invoicingIntegrationId={selectedInvoicingIntegrationId}
            businessId={effectiveBusinessId}
            onSuccess={handleSuccess}
            onCancel={handleCancel}
          />
        </div>
      )}

    </div>
  );
}
