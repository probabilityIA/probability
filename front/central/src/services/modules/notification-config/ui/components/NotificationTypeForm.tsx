"use client";

import { useState, useEffect } from "react";
import {
  NotificationType,
  CreateNotificationTypeDTO,
  UpdateNotificationTypeDTO,
} from "../../domain/types";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import { Checkbox } from "@/shared/ui/checkbox";
import { useToast } from "@/shared/providers/toast-provider";
import {
  createNotificationTypeAction,
  updateNotificationTypeAction,
} from "../../infra/actions";

interface NotificationTypeFormProps {
  type?: NotificationType;
  onSuccess: () => void;
  onCancel: () => void;
}

export function NotificationTypeForm({
  type,
  onSuccess,
  onCancel,
}: NotificationTypeFormProps) {
  const [loading, setLoading] = useState(false);
  const { showToast } = useToast();

  const [formData, setFormData] = useState<
    CreateNotificationTypeDTO | UpdateNotificationTypeDTO
  >({
    name: "",
    code: "",
    description: "",
    icon: "",
    is_active: true,
  });

  useEffect(() => {
    if (type) {
      setFormData({
        name: type.name,
        description: type.description || "",
        icon: type.icon || "",
        is_active: type.is_active,
      });
    }
  }, [type]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!type && !(formData as CreateNotificationTypeDTO).code) {
      showToast("El código es requerido", "error");
      return;
    }

    if (!formData.name) {
      showToast("El nombre es requerido", "error");
      return;
    }

    setLoading(true);

    try {
      let response;
      if (type) {
        response = await updateNotificationTypeAction(
          type.id,
          formData as UpdateNotificationTypeDTO
        );
      } else {
        response = await createNotificationTypeAction(
          formData as CreateNotificationTypeDTO
        );
      }

      if (response.success) {
        showToast(
          type ? "Tipo actualizado exitosamente" : "Tipo creado exitosamente",
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
        {/* Name */}
        <div className="grid gap-2">
          <Label htmlFor="name" className="flex items-center gap-1">
            Nombre
            <span className="text-red-500">*</span>
          </Label>
          <Input
            id="name"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            placeholder="ej: WhatsApp, SSE, Email"
            required
          />
        </div>

        {/* Code (solo en creación) */}
        {!type && (
          <div className="grid gap-2">
            <Label htmlFor="code" className="flex items-center gap-1">
              Código
              <span className="text-red-500">*</span>
            </Label>
            <Input
              id="code"
              value={(formData as CreateNotificationTypeDTO).code || ""}
              onChange={(e) =>
                setFormData({
                  ...formData,
                  code: e.target.value.toLowerCase(),
                } as CreateNotificationTypeDTO)
              }
              placeholder="ej: whatsapp, sse, email"
              required
            />
            <p className="text-xs text-gray-500">
              Identificador único (solo letras minúsculas, números y guiones)
            </p>
          </div>
        )}

        {/* Description */}
        <div className="grid gap-2">
          <Label htmlFor="description">Descripción</Label>
          <Input
            id="description"
            value={formData.description}
            onChange={(e) =>
              setFormData({ ...formData, description: e.target.value })
            }
            placeholder="Descripción opcional"
          />
        </div>

        {/* Icon */}
        <div className="grid gap-2">
          <Label htmlFor="icon">Ícono</Label>
          <Input
            id="icon"
            value={formData.icon}
            onChange={(e) => setFormData({ ...formData, icon: e.target.value })}
            placeholder="ej: message-circle, bell, mail"
          />
          <p className="text-xs text-gray-500">
            Nombre del ícono (Heroicons o similar)
          </p>
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
            Tipo activo
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
          {loading ? "Guardando..." : type ? "Actualizar" : "Crear"}
        </Button>
      </div>
    </form>
  );
}
