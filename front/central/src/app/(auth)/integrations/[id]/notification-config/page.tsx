/**
 * Página de configuración de notificaciones para una integración
 */

'use client';

import { useState } from 'react';
import { useParams } from 'next/navigation';
import {
  ConfigList,
  ConfigForm,
  useNotificationConfigs,
  type IntegrationNotificationConfig,
  type CreateNotificationConfigDTO,
  type UpdateNotificationConfigDTO,
} from '@/services/integrations/notification-config';

export default function NotificationConfigPage() {
  const params = useParams();
  const integrationId = parseInt(params.id as string);

  const [showForm, setShowForm] = useState(false);
  const [editingConfig, setEditingConfig] = useState<IntegrationNotificationConfig | null>(null);
  const [deletingConfig, setDeletingConfig] = useState<IntegrationNotificationConfig | null>(null);

  const { createConfig, updateConfig, deleteConfig, refetch } = useNotificationConfigs({
    integration_id: integrationId,
  });

  const handleCreate = async (data: CreateNotificationConfigDTO) => {
    await createConfig(data);
    setShowForm(false);
    refetch();
  };

  const handleUpdate = async (data: UpdateNotificationConfigDTO) => {
    if (editingConfig) {
      await updateConfig(editingConfig.id, data);
      setEditingConfig(null);
      refetch();
    }
  };

  // Wrapper que maneja ambos tipos de submit
  const handleSubmit = async (data: CreateNotificationConfigDTO | UpdateNotificationConfigDTO) => {
    if (editingConfig) {
      await handleUpdate(data as UpdateNotificationConfigDTO);
    } else {
      await handleCreate(data as CreateNotificationConfigDTO);
    }
    setShowForm(false);
  };

  const handleEdit = (config: IntegrationNotificationConfig) => {
    setEditingConfig(config);
    setShowForm(true);
  };

  const handleDelete = (config: IntegrationNotificationConfig) => {
    setDeletingConfig(config);
  };

  const confirmDelete = async () => {
    if (deletingConfig) {
      await deleteConfig(deletingConfig.id);
      setDeletingConfig(null);
      refetch();
    }
  };

  return (
    <div className="container mx-auto px-4 py-8 max-w-6xl">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">
          Configuración de Notificaciones
        </h1>
        <p className="text-gray-600">
          Configura cuándo y cómo enviar notificaciones automáticas
        </p>
      </div>

      {/* Action Button */}
      {!showForm && (
        <div className="mb-6">
          <button
            onClick={() => {
              setEditingConfig(null);
              setShowForm(true);
            }}
            className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
          >
            <svg
              className="w-5 h-5 mr-2"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 4v16m8-8H4"
              />
            </svg>
            Nueva Configuración
          </button>
        </div>
      )}

      {/* Form */}
      {showForm && (
        <div className="mb-8">
          <ConfigForm
            integrationId={integrationId}
            initialData={editingConfig || undefined}
            onSubmit={handleSubmit}
            onCancel={() => {
              setShowForm(false);
              setEditingConfig(null);
            }}
          />
        </div>
      )}

      {/* List */}
      {!showForm && (
        <ConfigList
          integrationId={integrationId}
          onEdit={handleEdit}
          onDelete={handleDelete}
        />
      )}

      {/* Delete Confirmation Modal */}
      {deletingConfig && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4">
            <h3 className="text-lg font-semibold text-gray-900 mb-2">
              Confirmar Eliminación
            </h3>
            <p className="text-gray-600 mb-6">
              ¿Estás seguro de que deseas eliminar esta configuración de notificación?
            </p>
            <div className="flex justify-end gap-3">
              <button
                onClick={() => setDeletingConfig(null)}
                className="px-4 py-2 text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
              >
                Cancelar
              </button>
              <button
                onClick={confirmDelete}
                className="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700"
              >
                Eliminar
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
