"use client";

import { useState } from "react";
import { NotificationTypeList } from "@/services/modules/notification-config/ui/components/NotificationTypeList";
import { NotificationTypeForm } from "@/services/modules/notification-config/ui/components/NotificationTypeForm";
import { Modal } from "@/shared/ui/modal";
import { NotificationType } from "@/services/modules/notification-config/domain/types";

export default function NotificationTypesPage() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedType, setSelectedType] = useState<
    NotificationType | undefined
  >(undefined);
  const [refreshKey, setRefreshKey] = useState(0);

  const handleCreate = () => {
    setSelectedType(undefined);
    setIsModalOpen(true);
  };

  const handleEdit = (type: NotificationType) => {
    setSelectedType(type);
    setIsModalOpen(true);
  };

  const handleSuccess = () => {
    setIsModalOpen(false);
    setRefreshKey((prev) => prev + 1);
  };

  return (
    <div className="min-h-screen bg-gray-50 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">
            Tipos de Notificación
          </h1>
          <p className="mt-2 text-sm text-gray-600">
            Administra los tipos/canales de notificación disponibles en el
            sistema (WhatsApp, SSE, Email, SMS, etc.)
          </p>
        </div>
      </div>

      <NotificationTypeList
        onEdit={handleEdit}
        onCreate={handleCreate}
        refreshKey={refreshKey}
      />

      <Modal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title={selectedType ? "Editar Tipo de Notificación" : "Nuevo Tipo de Notificación"}
      >
        <NotificationTypeForm
          type={selectedType}
          onSuccess={handleSuccess}
          onCancel={() => setIsModalOpen(false)}
        />
      </Modal>
    </div>
  );
}
