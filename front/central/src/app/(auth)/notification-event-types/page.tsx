"use client";

import { useState } from "react";
import { NotificationEventTypeList } from "@/services/modules/notification-config/ui/components/NotificationEventTypeList";
import { NotificationEventTypeForm } from "@/services/modules/notification-config/ui/components/NotificationEventTypeForm";
import { Modal } from "@/shared/ui/modal";
import { NotificationEventType } from "@/services/modules/notification-config/domain/types";

export default function NotificationEventTypesPage() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedEventType, setSelectedEventType] = useState<
    NotificationEventType | undefined
  >(undefined);
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
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">
            Eventos de Notificación
          </h1>
          <p className="mt-2 text-sm text-gray-600">
            Administra los eventos específicos por tipo de notificación
            (order.created, order.shipped, invoice.created, etc.)
          </p>
        </div>
      </div>

      <NotificationEventTypeList
        onEdit={handleEdit}
        onCreate={handleCreate}
        refreshKey={refreshKey}
      />

      <Modal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title={
          selectedEventType
            ? "Editar Evento de Notificación"
            : "Nuevo Evento de Notificación"
        }
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
