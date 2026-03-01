'use client';

import { useState } from 'react';
import { NotificationTypeList } from '@/services/modules/notification-config/ui/components/NotificationTypeList';
import { NotificationTypeForm } from '@/services/modules/notification-config/ui/components/NotificationTypeForm';
import { Modal } from '@/shared/ui/modal';
import type { NotificationType } from '@/services/modules/notification-config/domain/types';

export default function NotificationChannelsPage() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedChannel, setSelectedChannel] = useState<NotificationType | undefined>(undefined);
  const [refreshKey, setRefreshKey] = useState(0);

  const handleCreate = () => {
    setSelectedChannel(undefined);
    setIsModalOpen(true);
  };

  const handleEdit = (channel: NotificationType) => {
    setSelectedChannel(channel);
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
          <h1 className="text-xl font-semibold text-gray-900">Canales de Notificación</h1>
          <p className="text-sm text-gray-500 mt-0.5">
            Gestiona los canales disponibles (WhatsApp, Email, SMS, SSE)
          </p>
        </div>

        <NotificationTypeList
          onEdit={handleEdit}
          onCreate={handleCreate}
          refreshKey={refreshKey}
        />
      </div>

      <Modal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title={selectedChannel ? 'Editar Canal' : 'Nuevo Canal'}
        size="xl"
      >
        <NotificationTypeForm
          type={selectedChannel}
          onSuccess={handleSuccess}
          onCancel={() => setIsModalOpen(false)}
        />
      </Modal>
    </div>
  );
}
