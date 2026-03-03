"use client";

import { useState, useEffect } from "react";
import {
  NotificationEventType,
  NotificationType,
  CreateNotificationEventTypeDTO,
  UpdateNotificationEventTypeDTO,
} from "../../domain/types";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import { Checkbox } from "@/shared/ui/checkbox";
import { useToast } from "@/shared/providers/toast-provider";
import { useOrderStatuses } from "@/services/modules/orderstatus/ui";
import {
  getNotificationTypesAction,
  createNotificationEventTypeAction,
  updateNotificationEventTypeAction,
} from "../../infra/actions";

interface NotificationEventTypeFormProps {
  eventType?: NotificationEventType;
  onSuccess: () => void;
  onCancel: () => void;
}

export function NotificationEventTypeForm({
  eventType,
  onSuccess,
  onCancel,
}: NotificationEventTypeFormProps) {
  const [loading, setLoading] = useState(false);
  const [notificationTypes, setNotificationTypes] = useState<
    NotificationType[]
  >([]);
  const [loadingTypes, setLoadingTypes] = useState(false);
  const { showToast } = useToast();
  const { orderStatuses, loading: loadingOrderStatuses } = useOrderStatuses(true);

  const [formData, setFormData] = useState<
    CreateNotificationEventTypeDTO | UpdateNotificationEventTypeDTO
  >({
    notification_type_id: 0,
    event_code: "",
    event_name: "",
    description: "",
    is_active: true,
    allowed_order_status_ids: [],
  } as CreateNotificationEventTypeDTO);

  // Cargar notification types
  useEffect(() => {
    const loadNotificationTypes = async () => {
      setLoadingTypes(true);
      try {
        const result = await getNotificationTypesAction();
        if (result.success) {
          setNotificationTypes(result.data);
        }
      } catch (error) {
        showToast("Error al cargar tipos de notificación", "error");
      } finally {
        setLoadingTypes(false);
      }
    };

    loadNotificationTypes();
  }, []);

  // Cargar datos existentes (modo edición)
  useEffect(() => {
    if (eventType) {
      setFormData({
        event_name: eventType.event_name,
        description: eventType.description || "",
        is_active: eventType.is_active,
        allowed_order_status_ids: eventType.allowed_order_status_ids || [],
      } as UpdateNotificationEventTypeDTO);
    }
  }, [eventType]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (
      !eventType &&
      !(formData as CreateNotificationEventTypeDTO).notification_type_id
    ) {
      showToast("Selecciona un tipo de notificación", "error");
      return;
    }

    if (
      !eventType &&
      !(formData as CreateNotificationEventTypeDTO).event_code
    ) {
      showToast("El código de evento es requerido", "error");
      return;
    }

    if (!formData.event_name) {
      showToast("El nombre del evento es requerido", "error");
      return;
    }

    setLoading(true);

    try {
      let response;
      if (eventType) {
        response = await updateNotificationEventTypeAction(
          eventType.id,
          formData as UpdateNotificationEventTypeDTO
        );
      } else {
        response = await createNotificationEventTypeAction(
          formData as CreateNotificationEventTypeDTO
        );
      }

      if (response.success) {
        showToast(
          eventType
            ? "Evento actualizado exitosamente"
            : "Evento creado exitosamente",
          "success"
        );
        onSuccess();
      } else {
        showToast(response.error || "Error al guardar", "error");
      }
    } catch (error: any) {
      showToast(error.message || "Error inesperado", "error");
    } finally {
      setLoading(false);
    }
  };

  const selectedStatusIds = (formData as any).allowed_order_status_ids || [];

  const toggleStatusId = (statusId: number) => {
    const current: number[] = (formData as any).allowed_order_status_ids || [];
    const newIds = current.includes(statusId)
      ? current.filter((id: number) => id !== statusId)
      : [...current, statusId];
    setFormData({ ...formData, allowed_order_status_ids: newIds } as any);
  };

  const handleSelectAll = () => {
    const allIds = orderStatuses.map((s) => s.id);
    setFormData({ ...formData, allowed_order_status_ids: allIds } as any);
  };

  const handleClearAll = () => {
    setFormData({ ...formData, allowed_order_status_ids: [] } as any);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Left column: basic fields */}
        <div className="space-y-4">
          {/* Notification Type Selector (solo en creación) */}
          {!eventType && (
            <div className="grid gap-2">
              <Label htmlFor="notification_type_id" className="flex items-center gap-1">
                Tipo de Notificación
                <span className="text-red-500">*</span>
              </Label>
              <select
                id="notification_type_id"
                value={
                  (formData as CreateNotificationEventTypeDTO)
                    .notification_type_id || 0
                }
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    notification_type_id: parseInt(e.target.value),
                  } as CreateNotificationEventTypeDTO)
                }
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                required
                disabled={loadingTypes}
              >
                <option value="0">Seleccionar Tipo</option>
                {notificationTypes.map((type) => (
                  <option key={type.id} value={type.id}>
                    {type.name} ({type.code})
                  </option>
                ))}
              </select>
            </div>
          )}

          {/* Event Code (solo en creación) */}
          {!eventType && (
            <div className="grid gap-2">
              <Label htmlFor="event_code" className="flex items-center gap-1">
                Código de Evento
                <span className="text-red-500">*</span>
              </Label>
              <Input
                id="event_code"
                value={
                  (formData as CreateNotificationEventTypeDTO).event_code || ""
                }
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    event_code: e.target.value,
                  } as CreateNotificationEventTypeDTO)
                }
                placeholder="ej: order.created, order.shipped"
                required
              />
              <p className="text-xs text-gray-500">
                Código único del evento (ej: order.created, invoice.generated)
              </p>
            </div>
          )}

          {/* Event Name */}
          <div className="grid gap-2">
            <Label htmlFor="event_name" className="flex items-center gap-1">
              Nombre del Evento
              <span className="text-red-500">*</span>
            </Label>
            <Input
              id="event_name"
              value={formData.event_name}
              onChange={(e) =>
                setFormData({ ...formData, event_name: e.target.value })
              }
              placeholder="ej: Confirmación de Pedido, Pedido Enviado"
              required
            />
          </div>

          {/* Description */}
          <div className="grid gap-2">
            <Label htmlFor="description">Descripción</Label>
            <Input
              id="description"
              value={formData.description}
              onChange={(e) =>
                setFormData({ ...formData, description: e.target.value })
              }
              placeholder="Descripción opcional del evento"
            />
          </div>

          {/* Is Active */}
          <div className="flex items-center space-x-2">
            <Checkbox
              id="is_active"
              checked={formData.is_active}
              onCheckedChange={(checked: boolean) =>
                setFormData({ ...formData, is_active: checked })
              }
            />
            <label
              htmlFor="is_active"
              className="text-sm font-medium leading-none cursor-pointer"
            >
              Evento activo
            </label>
          </div>
        </div>

        {/* Right column: order statuses */}
        <div className="space-y-2">
          <Label className="flex items-center gap-1">
            Estados de Orden Permitidos
          </Label>
          <p className="text-xs text-gray-500">
            Selecciona qué estados de orden puede usar este evento. Si no seleccionas ninguno, se permiten todos.
          </p>

          {/* Select all / Clear all buttons */}
          {orderStatuses.length > 0 && (
            <div className="flex gap-2">
              <button
                type="button"
                onClick={handleSelectAll}
                className="text-xs text-blue-600 hover:text-blue-800 font-medium"
              >
                Seleccionar todos
              </button>
              <span className="text-xs text-gray-300">|</span>
              <button
                type="button"
                onClick={handleClearAll}
                className="text-xs text-gray-500 hover:text-gray-700 font-medium"
              >
                Limpiar
              </button>
            </div>
          )}

          <div className="border rounded-lg max-h-[400px] overflow-y-auto p-3">
            {loadingOrderStatuses ? (
              <p className="text-sm text-gray-500">Cargando estados...</p>
            ) : orderStatuses.length === 0 ? (
              <p className="text-sm text-gray-500">No hay estados disponibles</p>
            ) : (
              <div className="grid grid-cols-2 xl:grid-cols-3 gap-2">
                {orderStatuses.map((status) => {
                  const isChecked = selectedStatusIds.includes(status.id);
                  const statusColor = status.color || "#9CA3AF";
                  return (
                    <button
                      key={status.id}
                      type="button"
                      onClick={() => toggleStatusId(status.id)}
                      className={`flex items-center justify-between gap-2 px-3 py-2 rounded-lg border transition-colors ${
                        isChecked
                          ? "bg-blue-50 border-blue-200"
                          : "bg-white border-gray-200 hover:bg-gray-50"
                      }`}
                    >
                      <span
                        className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium truncate"
                        style={{
                          backgroundColor: statusColor + "20",
                          color: statusColor,
                        }}
                      >
                        {status.name}
                      </span>
                      <div
                        className={`relative w-9 h-5 rounded-full shrink-0 transition-colors ${
                          isChecked ? "bg-blue-500" : "bg-gray-300"
                        }`}
                      >
                        <div
                          className={`absolute top-0.5 w-4 h-4 rounded-full bg-white shadow transition-transform ${
                            isChecked ? "translate-x-4" : "translate-x-0.5"
                          }`}
                        />
                      </div>
                    </button>
                  );
                })}
              </div>
            )}
          </div>
          {selectedStatusIds.length === 0 && (
            <p className="text-xs text-amber-600">
              Sin selección = todos los estados permitidos
            </p>
          )}
          {selectedStatusIds.length > 0 && (
            <p className="text-xs text-blue-600">
              {selectedStatusIds.length} estado(s) seleccionado(s)
            </p>
          )}
        </div>
      </div>

      <div className="flex justify-end gap-2 pt-4 border-t">
        <button
          type="button"
          onClick={onCancel}
          disabled={loading}
          className="p-2 rounded-lg bg-gray-100 text-gray-500 hover:bg-gray-200 transition-colors disabled:opacity-40"
          title="Cancelar"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
        <button
          type="submit"
          disabled={loading}
          className="p-2 rounded-lg bg-green-50 text-green-600 hover:bg-green-100 transition-colors disabled:opacity-40"
          title={loading ? "Guardando..." : eventType ? "Actualizar" : "Crear"}
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
          </svg>
        </button>
      </div>
    </form>
  );
}
