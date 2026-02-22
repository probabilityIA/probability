/**
 * StatusSelector - Selector de estados de orden
 * (Migrado desde services/integrations/notification-config)
 */

'use client';

import { useIntegrationOrderStatuses } from '../hooks/useIntegrationNotificationConfigs';

interface StatusSelectorProps {
  selectedStatuses: string[];
  onChange: (statuses: string[]) => void;
  disabled?: boolean;
}

export function StatusSelector({
  selectedStatuses,
  onChange,
  disabled = false,
}: StatusSelectorProps) {
  const { orderStatuses, loading, error } = useIntegrationOrderStatuses();

  const handleToggle = (statusCode: string) => {
    if (selectedStatuses.includes(statusCode)) {
      onChange(selectedStatuses.filter((code) => code !== statusCode));
    } else {
      onChange([...selectedStatuses, statusCode]);
    }
  };

  const handleSelectAll = () => {
    onChange(orderStatuses.map((s) => s.code));
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
          Estados de Orden
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
        Deja vacío para aplicar a todos los estados
      </p>

      <div className="max-h-60 overflow-y-auto border border-gray-300 rounded-lg p-3 space-y-2">
        {orderStatuses.length === 0 ? (
          <p className="text-sm text-gray-500 text-center py-4">
            No hay estados de orden disponibles
          </p>
        ) : (
          orderStatuses.map((status) => (
            <label
              key={status.id}
              className="flex items-center p-2 rounded hover:bg-gray-50 cursor-pointer"
            >
              <input
                type="checkbox"
                checked={selectedStatuses.includes(status.code)}
                onChange={() => handleToggle(status.code)}
                disabled={disabled}
                className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
              />
              <div className="ml-3 flex-1 flex items-center gap-2">
                {status.color && (
                  <div
                    className="w-3 h-3 rounded-full"
                    style={{ backgroundColor: status.color }}
                  />
                )}
                <div>
                  <div className="text-sm font-medium text-gray-900">
                    {status.name}
                  </div>
                  <div className="text-xs text-gray-500">{status.code}</div>
                </div>
              </div>
            </label>
          ))
        )}
      </div>

      {selectedStatuses.length > 0 && (
        <div className="text-xs text-gray-600">
          {selectedStatuses.length} estado{selectedStatuses.length !== 1 ? 's' : ''} seleccionado
          {selectedStatuses.length !== 1 ? 's' : ''}
        </div>
      )}
    </div>
  );
}
