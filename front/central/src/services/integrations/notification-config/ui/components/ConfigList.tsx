/**
 * ConfigList - Lista de configuraciones de notificación
 */

'use client';

import { useState } from 'react';
import type { IntegrationNotificationConfig } from '../../domain/types';
import { useNotificationConfigs } from '../hooks/useNotificationConfigs';

interface ConfigListProps {
  integrationId: number;
  onEdit: (config: IntegrationNotificationConfig) => void;
  onDelete: (config: IntegrationNotificationConfig) => void;
}

export function ConfigList({ integrationId, onEdit, onDelete }: ConfigListProps) {
  const { configs, loading, error, toggleConfig } = useNotificationConfigs({
    integration_id: integrationId,
  });

  const [togglingId, setTogglingId] = useState<number | null>(null);

  const handleToggle = async (config: IntegrationNotificationConfig) => {
    setTogglingId(config.id);
    try {
      await toggleConfig(config.id, !config.is_active);
    } catch (err) {
      console.error('Error toggling config:', err);
    } finally {
      setTogglingId(null);
    }
  };

  if (loading) {
    return (
      <div className="space-y-4">
        {[1, 2, 3].map((i) => (
          <div key={i} className="animate-pulse bg-gray-200 h-32 rounded-lg" />
        ))}
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
        {error}
      </div>
    );
  }

  if (configs.length === 0) {
    return (
      <div className="text-center py-12 bg-gray-50 rounded-lg border-2 border-dashed border-gray-300">
        <svg
          className="mx-auto h-12 w-12 text-gray-400"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
          />
        </svg>
        <h3 className="mt-2 text-sm font-medium text-gray-900">
          No hay configuraciones de notificación
        </h3>
        <p className="mt-1 text-sm text-gray-500">
          Comienza creando una nueva configuración
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {configs.map((config) => (
        <div
          key={config.id}
          className={`bg-white border rounded-lg p-4 transition-all ${
            config.is_active ? 'border-gray-200' : 'border-gray-300 bg-gray-50'
          }`}
        >
          <div className="flex items-start justify-between">
            <div className="flex-1">
              {/* Header */}
              <div className="flex items-center gap-3 mb-2">
                <span
                  className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                    config.notification_type === 'whatsapp'
                      ? 'bg-green-100 text-green-800'
                      : config.notification_type === 'email'
                      ? 'bg-blue-100 text-blue-800'
                      : 'bg-purple-100 text-purple-800'
                  }`}
                >
                  {config.notification_type.toUpperCase()}
                </span>
                <span className="text-sm text-gray-500">
                  Trigger: <span className="font-mono">{config.conditions.trigger}</span>
                </span>
                <span className="text-sm text-gray-500">
                  Prioridad: <span className="font-semibold">{config.priority}</span>
                </span>
              </div>

              {/* Description */}
              {config.description && (
                <p className="text-sm text-gray-700 mb-3">{config.description}</p>
              )}

              {/* Conditions */}
              <div className="flex flex-wrap gap-2 text-xs">
                {config.conditions.statuses.length > 0 && (
                  <div className="bg-blue-50 text-blue-700 px-2 py-1 rounded">
                    Estados: {config.conditions.statuses.join(', ')}
                  </div>
                )}
                {config.conditions.payment_methods.length > 0 && (
                  <div className="bg-orange-50 text-orange-700 px-2 py-1 rounded">
                    Métodos de pago: {config.conditions.payment_methods.join(', ')}
                  </div>
                )}
                {config.conditions.statuses.length === 0 &&
                  config.conditions.payment_methods.length === 0 && (
                    <div className="bg-gray-100 text-gray-600 px-2 py-1 rounded">
                      Sin filtros (aplica a todas las órdenes)
                    </div>
                  )}
              </div>

              {/* Template info */}
              <div className="mt-2 text-xs text-gray-500">
                Plantilla: <span className="font-mono">{config.config.template_name}</span> |
                Idioma: {config.config.language} | Destinatario: {config.config.recipient_type}
              </div>
            </div>

            {/* Actions */}
            <div className="flex items-center gap-2 ml-4">
              {/* Toggle Active */}
              <button
                onClick={() => handleToggle(config)}
                disabled={togglingId === config.id}
                className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                  config.is_active ? 'bg-green-600' : 'bg-gray-200'
                } ${togglingId === config.id ? 'opacity-50' : ''}`}
                title={config.is_active ? 'Desactivar' : 'Activar'}
              >
                <span
                  className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                    config.is_active ? 'translate-x-6' : 'translate-x-1'
                  }`}
                />
              </button>

              {/* Edit */}
              <button
                onClick={() => onEdit(config)}
                className="p-2 text-blue-600 hover:bg-blue-50 rounded transition-colors"
                title="Editar"
              >
                <svg
                  className="w-5 h-5"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                  />
                </svg>
              </button>

              {/* Delete */}
              <button
                onClick={() => onDelete(config)}
                className="p-2 text-red-600 hover:bg-red-50 rounded transition-colors"
                title="Eliminar"
              >
                <svg
                  className="w-5 h-5"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                  />
                </svg>
              </button>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
