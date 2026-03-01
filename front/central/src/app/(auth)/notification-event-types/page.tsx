'use client';

import { useState } from 'react';
import { NotificationEventTypeList } from '@/services/modules/notification-config/ui/components/NotificationEventTypeList';
import { NotificationEventTypeForm } from '@/services/modules/notification-config/ui/components/NotificationEventTypeForm';
import { Modal } from '@/shared/ui/modal';
import type { NotificationEventType } from '@/services/modules/notification-config/domain/types';

export default function NotificationEventTypesPage() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedEventType, setSelectedEventType] = useState<NotificationEventType | undefined>(undefined);
  const [refreshKey, setRefreshKey] = useState(0);

  const handleCreate = () => {
    setSelectedEventType(undefined);
    setIsModalOpen(true);
  };

  const handleEdit = (eventType: NotificationEventType) => {
    setSelectedEventType(eventType);
    setIsModalOpen(true);
  };

  const handleSuccess = () => {
    setIsModalOpen(false);
    setRefreshKey((prev) => prev + 1);
  };

  return (
    <div className="min-h-screen bg-gray-50 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
      <div className="space-y-6">
        <div>
          <h1 className="text-xl font-semibold text-gray-900">Tipos de Eventos</h1>
          <p className="text-sm text-gray-500 mt-0.5">
            Gestiona los tipos de eventos de notificación
          </p>
        </div>

        <NotificationEventTypeList
          onEdit={handleEdit}
          onCreate={handleCreate}
          refreshKey={refreshKey}
        />
      </div>

      <Modal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title={selectedEventType ? 'Editar Tipo de Evento' : 'Nuevo Tipo de Evento'}
        size="4xl"
      >
        <NotificationEventTypeForm
          eventType={selectedEventType}
          onSuccess={handleSuccess}
          onCancel={() => setIsModalOpen(false)}
        />
      </Modal>
    </div>
  );
}
