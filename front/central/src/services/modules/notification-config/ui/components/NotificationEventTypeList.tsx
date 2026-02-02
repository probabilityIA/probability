"use client";

import { useState, useEffect } from "react";
import { NotificationEventType, NotificationType } from "../../domain/types";
import { Button } from "@/shared/ui/button";
import { useToast } from "@/shared/providers/toast-provider";
import {
  getNotificationTypesAction,
  getNotificationEventTypesAction,
  deleteNotificationEventTypeAction,
  toggleNotificationEventTypeActiveAction,
} from "../../infra/actions";
import { ConfirmModal } from "@/shared/ui/confirm-modal";

interface NotificationEventTypeListProps {
  onEdit: (eventType: NotificationEventType) => void;
  onCreate: () => void;
  refreshKey: number;
}

export function NotificationEventTypeList({
  onEdit,
  onCreate,
  refreshKey,
}: NotificationEventTypeListProps) {
  const [eventTypes, setEventTypes] = useState<NotificationEventType[]>([]);
  const [notificationTypes, setNotificationTypes] = useState<
    NotificationType[]
  >([]);
  const [selectedTypeFilter, setSelectedTypeFilter] = useState<number>(0);
  const [loading, setLoading] = useState(false);
  const [deleteModal, setDeleteModal] = useState<{
    isOpen: boolean;
    eventType?: NotificationEventType;
  }>({ isOpen: false });
  const { showToast } = useToast();

  // Cargar tipos de notificación para el filtro
  useEffect(() => {
    const loadNotificationTypes = async () => {
      try {
        const result = await getNotificationTypesAction();
        if (result.success) {
          setNotificationTypes(result.data);
        }
      } catch (error) {
        console.error("Error loading notification types:", error);
      }
    };

    loadNotificationTypes();
  }, []);

  const fetchEventTypes = async () => {
    setLoading(true);
    try {
      const result = await getNotificationEventTypesAction(
        selectedTypeFilter || undefined
      );
      if (result.success) {
        setEventTypes(result.data);
      } else {
        showToast("Error al cargar eventos", "error");
      }
    } catch (error) {
      showToast("Error al cargar eventos", "error");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchEventTypes();
  }, [refreshKey, selectedTypeFilter]);

  const handleDelete = async () => {
    if (!deleteModal.eventType) return;

    try {
      const result = await deleteNotificationEventTypeAction(
        deleteModal.eventType.id
      );
      if (result.success) {
        showToast("Evento eliminado exitosamente", "success");
        setDeleteModal({ isOpen: false });
        fetchEventTypes();
      } else {
        showToast(result.error || "Error al eliminar", "error");
      }
    } catch (error: any) {
      showToast(error.message || "Error al eliminar", "error");
    }
  };

  const handleToggleActive = async (eventType: NotificationEventType) => {
    try {
      const result = await toggleNotificationEventTypeActiveAction(
        eventType.id
      );
      if (result.success) {
        const newStatus = result.data.is_active ? "activado" : "desactivado";
        showToast(`Evento ${newStatus} exitosamente`, "success");
        fetchEventTypes();
      } else {
        showToast(result.error || "Error al cambiar estado", "error");
      }
    } catch (error: any) {
      showToast(error.message || "Error al cambiar estado", "error");
    }
  };

  if (loading) {
    return <div className="text-center py-8">Cargando...</div>;
  }

  return (
    <div className="bg-white shadow-md rounded-lg overflow-hidden">
      <div className="p-4 border-b">
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          <h2 className="text-lg font-semibold">Eventos de Notificación</h2>
          <div className="flex gap-2 items-center w-full sm:w-auto">
            <select
              value={selectedTypeFilter}
              onChange={(e) => setSelectedTypeFilter(parseInt(e.target.value))}
              className="px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500 text-sm"
            >
              <option value="0">Todos los tipos</option>
              {notificationTypes.map((type) => (
                <option key={type.id} value={type.id}>
                  {type.name}
                </option>
              ))}
            </select>
            <Button onClick={onCreate}>+ Nuevo Evento</Button>
          </div>
        </div>
      </div>

      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Tipo
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Evento
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Descripción
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Estado
              </th>
              <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Acciones
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {eventTypes.length === 0 ? (
              <tr>
                <td colSpan={5} className="px-6 py-4 text-center text-gray-500">
                  No hay eventos de notificación
                </td>
              </tr>
            ) : (
              eventTypes.map((eventType) => (
                <tr key={eventType.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className="text-sm font-medium text-gray-900">
                      {eventType.notification_type?.name || "-"}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm font-medium text-gray-900">
                      {eventType.event_name}
                    </div>
                  </td>
                  <td className="px-6 py-4">
                    <div className="text-sm text-gray-500">
                      {eventType.description || "-"}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span
                      className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                        eventType.is_active
                          ? "bg-green-100 text-green-800"
                          : "bg-red-100 text-red-800"
                      }`}
                    >
                      {eventType.is_active ? "Activo" : "Inactivo"}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <Button
                      variant={eventType.is_active ? "danger" : "outline"}
                      size="sm"
                      onClick={() => handleToggleActive(eventType)}
                      className="mr-2"
                    >
                      {eventType.is_active ? "Desactivar" : "Activar"}
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => onEdit(eventType)}
                      className="mr-2"
                    >
                      Editar
                    </Button>
                    <Button
                      variant="danger"
                      size="sm"
                      onClick={() =>
                        setDeleteModal({ isOpen: true, eventType })
                      }
                    >
                      Eliminar
                    </Button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      <ConfirmModal
        isOpen={deleteModal.isOpen}
        onClose={() => setDeleteModal({ isOpen: false })}
        onConfirm={handleDelete}
        title="Eliminar Evento de Notificación"
        message={`¿Estás seguro de eliminar el evento "${deleteModal.eventType?.event_name}"? Esta acción no se puede deshacer.`}
        confirmText="Eliminar"
        cancelText="Cancelar"
      />
    </div>
  );
}
