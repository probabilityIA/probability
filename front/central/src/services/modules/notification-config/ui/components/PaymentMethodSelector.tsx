/**
 * PaymentMethodSelector - Selector de métodos de pago
 * (Migrado desde services/integrations/notification-config)
 */

'use client';

import { usePaymentMethods } from '../hooks/useIntegrationNotificationConfigs';

interface PaymentMethodSelectorProps {
  selectedMethods: number[];
  onChange: (methods: number[]) => void;
  disabled?: boolean;
}

export function PaymentMethodSelector({
  selectedMethods,
  onChange,
  disabled = false,
}: PaymentMethodSelectorProps) {
  const { paymentMethods, loading, error } = usePaymentMethods();

  const handleToggle = (methodId: number) => {
    if (selectedMethods.includes(methodId)) {
      onChange(selectedMethods.filter((id) => id !== methodId));
    } else {
      onChange([...selectedMethods, methodId]);
    }
  };

  const handleSelectAll = () => {
    onChange(paymentMethods.map((m) => m.id));
  };

  const handleClearAll = () => {
    onChange([]);
  };

  if (loading) {
    return (
      <div className="animate-pulse space-y-2">
        {[1, 2, 3].map((i) => (
          <div key={i} className="h-10 bg-gray-200 rounded" />
        ))}
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-red-600 text-sm p-3 bg-red-50 rounded border border-red-200">
        {error}
      </div>
    );
  }

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <label className="block text-sm font-medium text-gray-700">
          Métodos de Pago
        </label>
        <div className="flex gap-2">
          <button
            type="button"
            onClick={handleSelectAll}
            disabled={disabled}
            className="text-xs text-blue-600 hover:text-blue-800 disabled:opacity-50"
          >
            Seleccionar todos
          </button>
          <span className="text-xs text-gray-400">|</span>
          <button
            type="button"
            onClick={handleClearAll}
            disabled={disabled}
            className="text-xs text-gray-600 hover:text-gray-800 disabled:opacity-50"
          >
            Limpiar
          </button>
        </div>
      </div>

      <p className="text-xs text-gray-500">
        Deja vacío para aplicar a todos los métodos de pago
      </p>

      <div className="max-h-60 overflow-y-auto border border-gray-300 rounded-lg p-3 space-y-2">
        {paymentMethods.length === 0 ? (
          <p className="text-sm text-gray-500 text-center py-4">
            No hay métodos de pago disponibles
          </p>
        ) : (
          paymentMethods.map((method) => (
            <label
              key={method.id}
              className={`flex items-center p-2 rounded hover:bg-gray-50 cursor-pointer ${
                !method.is_active ? 'opacity-50' : ''
              }`}
            >
              <input
                type="checkbox"
                checked={selectedMethods.includes(method.id)}
                onChange={() => handleToggle(method.id)}
                disabled={disabled || !method.is_active}
                className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
              />
              <div className="ml-3 flex-1">
                <div className="text-sm font-medium text-gray-900">
                  {method.name}
                  {!method.is_active && (
                    <span className="ml-2 text-xs text-gray-500">(Inactivo)</span>
                  )}
                </div>
                <div className="text-xs text-gray-500">{method.code}</div>
              </div>
            </label>
          ))
        )}
      </div>

      {selectedMethods.length > 0 && (
        <div className="text-xs text-gray-600">
          {selectedMethods.length} método{selectedMethods.length !== 1 ? 's' : ''} seleccionado
          {selectedMethods.length !== 1 ? 's' : ''}
        </div>
      )}
    </div>
  );
}
