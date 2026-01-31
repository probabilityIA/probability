"use client";

import { useState, useEffect } from "react";
import {
  NotificationEventType,
  NotificationType,
  CreateNotificationEventTypeDTO,
  UpdateNotificationEventTypeDTO,
} from "../../domain/types";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import { Checkbox } from "@/shared/ui/checkbox";
import { useToast } from "@/shared/providers/toast-provider";
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

  const [formData, setFormData] = useState<
    CreateNotificationEventTypeDTO | UpdateNotificationEventTypeDTO
  >({
    notification_type_id: 0,
    event_code: "",
    event_name: "",
    description: "",
    is_active: true,
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

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="grid gap-4">
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

      <div className="flex justify-end gap-2 pt-4 border-t">
        <Button
          type="button"
          variant="outline"
          onClick={onCancel}
          disabled={loading}
        >
          Cancelar
        </Button>
        <Button type="submit" disabled={loading}>
          {loading ? "Guardando..." : eventType ? "Actualizar" : "Crear"}
        </Button>
      </div>
    </form>
  );
}
