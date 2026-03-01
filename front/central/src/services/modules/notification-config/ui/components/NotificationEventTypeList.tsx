"use client";

import { useState, useEffect } from "react";
import { NotificationEventType, NotificationType } from "../../domain/types";
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
      if (result.success && result.data) {
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
            <button
              type="button"
              onClick={onCreate}
              className="p-2 rounded-lg bg-blue-50 text-blue-600 hover:bg-blue-100 transition-colors"
              title="Nuevo Evento"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
            </button>
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
                      className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${eventType.is_active
                          ? "bg-green-100 text-green-800"
                          : "bg-red-100 text-red-800"
                        }`}
                    >
                      {eventType.is_active ? "Activo" : "Inactivo"}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <div className="flex items-center justify-end gap-1">
                      <button
                        type="button"
                        onClick={() => handleToggleActive(eventType)}
                        className={`p-1.5 rounded-md transition-colors ${
                          eventType.is_active
                            ? "bg-green-50 text-green-600 hover:bg-green-100"
                            : "bg-gray-100 text-gray-400 hover:bg-gray-200"
                        }`}
                        title={eventType.is_active ? "Desactivar" : "Activar"}
                      >
                        <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          {eventType.is_active ? (
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z" />
                          ) : (
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                          )}
                        </svg>
                      </button>
                      <button
                        type="button"
                        onClick={() => onEdit(eventType)}
                        className="p-1.5 rounded-md bg-amber-50 text-amber-600 hover:bg-amber-100 transition-colors"
                        title="Editar"
                      >
                        <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                        </svg>
                      </button>
                      <button
                        type="button"
                        onClick={() => setDeleteModal({ isOpen: true, eventType })}
                        className="p-1.5 rounded-md bg-red-50 text-red-500 hover:bg-red-100 transition-colors"
                        title="Eliminar"
                      >
                        <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
                      </button>
                    </div>
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
