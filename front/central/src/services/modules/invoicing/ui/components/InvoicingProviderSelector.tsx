'use client';

import type { InvoicingProvider } from '@/services/modules/invoicing/domain/types';

interface InvoicingProviderSelectorProps {
  providers: InvoicingProvider[];
  value?: number;
  onChange: (providerId: number) => void;
  disabled?: boolean;
}

/**
 * Selector de proveedor de facturación electrónica
 */
export function InvoicingProviderSelector({
  providers,
  value,
  onChange,
  disabled = false,
}: InvoicingProviderSelectorProps) {
  const selectedProvider = providers.find((p) => p.id === value);

  return (
    <div>
      <label className="block text-sm font-medium text-gray-700 mb-2">
        Proveedor de Facturación <span className="text-red-500">*</span>
      </label>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {providers.map((provider) => {
          const isSelected = value === provider.id;
          const isDisabled = disabled || !provider.is_active;

          return (
            <div
              key={provider.id}
              onClick={() => !isDisabled && onChange(provider.id)}
              className={`
                relative border-2 rounded-lg p-4 cursor-pointer transition-all
                ${
                  isSelected
                    ? 'border-blue-500 bg-blue-50'
                    : 'border-gray-200 bg-white hover:border-gray-300'
                }
                ${isDisabled ? 'opacity-50 cursor-not-allowed' : ''}
              `}
            >
              {/* Radio button visual */}
              <div className="flex items-start gap-3">
                <div
                  className={`
                  w-5 h-5 rounded-full border-2 flex items-center justify-center mt-0.5
                  ${
                    isSelected
                      ? 'border-blue-500 bg-blue-500'
                      : 'border-gray-300 bg-white'
                  }
                `}
                >
                  {isSelected && (
                    <div className="w-2 h-2 bg-white rounded-full" />
                  )}
                </div>

                <div className="flex-1">
                  {/* Nombre del proveedor */}
                  <h3 className="font-semibold text-gray-900">
                    {provider.name}
                  </h3>

                  {/* Tipo de proveedor */}
                  <p className="text-xs text-gray-500 mt-1">
                    {provider.provider_type_code}
                  </p>

                  {/* Descripción */}
                  {provider.description && (
                    <p className="text-sm text-gray-600 mt-2">
                      {provider.description}
                    </p>
                  )}

                  {/* Estado */}
                  {!provider.is_active && (
                    <span className="inline-block mt-2 px-2 py-1 text-xs bg-gray-200 text-gray-600 rounded">
                      No disponible
                    </span>
                  )}
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {/* Información del proveedor seleccionado */}
      {selectedProvider && (
        <div className="mt-4 p-4 bg-gray-50 border border-gray-200 rounded-lg">
          <h4 className="text-sm font-medium text-gray-900 mb-2">
            Información del proveedor
          </h4>
          <dl className="space-y-1 text-sm">
            <div className="flex justify-between">
              <dt className="text-gray-600">Tipo:</dt>
              <dd className="text-gray-900 font-mono">
                {selectedProvider.provider_type_code}
              </dd>
            </div>
            {selectedProvider.description && (
              <div className="mt-2">
                <dt className="text-gray-600">Descripción:</dt>
                <dd className="text-gray-900">{selectedProvider.description}</dd>
              </div>
            )}
          </dl>
        </div>
      )}

      {providers.length === 0 && (
        <div className="text-center py-8 text-gray-500">
          <p>No hay proveedores de facturación disponibles.</p>
          <p className="text-sm mt-1">
            Contacte con el administrador para configurar proveedores.
          </p>
        </div>
      )}
    </div>
  );
}
