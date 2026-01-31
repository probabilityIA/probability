"use client";

import { useState, useEffect } from "react";
import {
  NotificationConfig,
  CreateConfigDTO,
  UpdateConfigDTO,
  NotificationType,
  NotificationEventType,
} from "../../domain/types";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import { useToast } from "@/shared/providers/toast-provider";
import { Checkbox } from "@/shared/ui/checkbox";
import { TokenStorage } from "@/shared/utils/token-storage";
import { useBusinessesSimple } from "@/services/auth/business/ui/hooks/useBusinessesSimple";
import { useIntegrationsSimple } from "@/services/integrations/core/ui/hooks/useIntegrationsSimple";
import { useOrderStatuses } from "@/services/modules/orderstatus/ui";
import {
  getNotificationTypesAction,
  getNotificationEventTypesAction,
  createConfigAction,
  updateConfigAction,
} from "../../infra/actions";

interface NotificationConfigFormProps {
  config?: NotificationConfig;
  onSuccess: () => void;
  onCancel: () => void;
}

export function NotificationConfigForm({
  config,
  onSuccess,
  onCancel,
}: NotificationConfigFormProps) {
  const [loading, setLoading] = useState(false);
  const { showToast } = useToast();

  // Obtener permisos del usuario actual
  const permissions = TokenStorage.getPermissions();
  const isSuperAdmin = permissions?.is_super || false;
  const userBusinessId = permissions?.business_id || 0;

  // Hooks para obtener datos (usando endpoints ligeros)
  const { businesses, loading: loadingBusinesses } = useBusinessesSimple();
  const { integrations, loading: loadingIntegrations } = useIntegrationsSimple();
  const { orderStatuses, loading: loadingOrderStatuses } = useOrderStatuses(true); // Solo activos

  // Estados para notification types y event types
  const [notificationTypes, setNotificationTypes] = useState<NotificationType[]>([]);
  const [eventTypes, setEventTypes] = useState<NotificationEventType[]>([]);
  const [loadingTypes, setLoadingTypes] = useState(false);
  const [loadingEvents, setLoadingEvents] = useState(false);

  // Form data con nueva estructura
  const [formData, setFormData] = useState<CreateConfigDTO>({
    business_id: isSuperAdmin ? 0 : userBusinessId,
    integration_id: 0, // OBLIGATORIO - La integración origen
    notification_type_id: 0, // OBLIGATORIO - Canal de salida
    notification_event_type_id: 0, // OBLIGATORIO - Tipo de evento
    enabled: true,
    description: "",
    order_status_ids: [], // Estados de orden a notificar
  });

  // ============================================
  // CARGAR NOTIFICATION TYPES AL MONTAR
  // ============================================
  useEffect(() => {
    const loadNotificationTypes = async () => {
      setLoadingTypes(true);
      try {
        const result = await getNotificationTypesAction();
        if (result.success) {
          setNotificationTypes(result.data);
        } else {
          showToast("Error al cargar tipos de notificación", "error");
        }
      } catch (error) {
        showToast("Error al cargar tipos de notificación", "error");
      } finally {
        setLoadingTypes(false);
      }
    };

    loadNotificationTypes();
  }, []);

  // ============================================
  // CARGAR EVENT TYPES cuando cambia notification_type_id
  // ============================================
  useEffect(() => {
    if (formData.notification_type_id > 0) {
      const loadEventTypes = async () => {
        setLoadingEvents(true);
        try {
          const result = await getNotificationEventTypesAction(
            formData.notification_type_id
          );
          if (result.success) {
            setEventTypes(result.data);
          } else {
            showToast("Error al cargar eventos", "error");
            setEventTypes([]);
          }
        } catch (error) {
          showToast("Error al cargar eventos", "error");
          setEventTypes([]);
        } finally {
          setLoadingEvents(false);
        }
      };

      loadEventTypes();
    } else {
      setEventTypes([]);
    }
  }, [formData.notification_type_id]);

  // ============================================
  // CARGAR DATOS EXISTENTES (MODO EDICIÓN)
  // ============================================
  useEffect(() => {
    if (config) {
      setFormData({
        business_id: config.business_id,
        integration_id: config.integration_id,
        notification_type_id: config.notification_type_id,
        notification_event_type_id: config.notification_event_type_id,
        enabled: config.enabled,
        description: config.description || "",
        order_status_ids: config.order_status_ids || [],
      });
    } else if (!isSuperAdmin && userBusinessId) {
      // Auto-select business for business users when creating
      setFormData((prev) => ({ ...prev, business_id: userBusinessId }));
    }
  }, [config, isSuperAdmin, userBusinessId]);

  // ============================================
  // SUBMIT HANDLER
  // ============================================
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Validaciones
    if (!formData.business_id) {
      showToast("Selecciona un business", "error");
      return;
    }

    if (!formData.integration_id) {
      showToast("Selecciona una integración origen", "error");
      return;
    }

    if (!formData.notification_type_id) {
      showToast("Selecciona un tipo de notificación", "error");
      return;
    }

    if (!formData.notification_event_type_id) {
      showToast("Selecciona un evento", "error");
      return;
    }

    setLoading(true);

    try {
      let response;
      if (config) {
        const updateDto: UpdateConfigDTO = {
          integration_id: formData.integration_id,
          notification_type_id: formData.notification_type_id,
          notification_event_type_id: formData.notification_event_type_id,
          enabled: formData.enabled,
          description: formData.description,
          order_status_ids: formData.order_status_ids,
        };
        response = await updateConfigAction(config.id, updateDto);
      } else {
        response = await createConfigAction(formData);
      }

      if (response.success) {
        showToast(
          config ? "Configuración actualizada" : "Configuración creada",
          "success"
        );
        onSuccess();
      } else {
        showToast(
          response.error || "Error al guardar configuración",
          "error"
        );
      }
    } catch (error) {
      showToast("Error inesperado", "error");
    } finally {
      setLoading(false);
    }
  };

  // ============================================
  // HANDLER PARA ORDER STATUS CHECKBOX
  // ============================================
  const handleOrderStatusToggle = (statusId: number) => {
    setFormData((prev) => {
      const currentIds = prev.order_status_ids || [];
      if (currentIds.includes(statusId)) {
        return {
          ...prev,
          order_status_ids: currentIds.filter((id) => id !== statusId),
        };
      } else {
        return {
          ...prev,
          order_status_ids: [...currentIds, statusId],
        };
      }
    });
  };

  // Filtrar integraciones por business seleccionado
  const filteredIntegrations = integrations.filter(
    (integration) =>
      integration.business_id === formData.business_id && integration.is_active
  );

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="grid gap-6">
        {/* ============================================
            PASO 1: BUSINESS SELECTOR
            ============================================ */}
        {!config && (
          <div className="grid gap-2">
            <Label htmlFor="business_id">
              Business{" "}
              {!isSuperAdmin && (
                <span className="text-xs text-gray-500">
                  (asignado automáticamente)
                </span>
              )}
            </Label>
            {isSuperAdmin ? (
              <select
                id="business_id"
                value={formData.business_id}
                onChange={(e) => {
                  const businessId = parseInt(e.target.value) || 0;
                  setFormData({
                    ...formData,
                    business_id: businessId,
                    integration_id: 0, // Reset
                  });
                }}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                required
                disabled={loadingBusinesses}
              >
                <option value="0">Seleccionar Business</option>
                {businesses.map((business) => (
                  <option key={business.id} value={business.id}>
                    {business.name}
                  </option>
                ))}
              </select>
            ) : (
              <Input
                id="business_id"
                type="text"
                value={
                  permissions?.business_name ||
                  `Business ID: ${formData.business_id}`
                }
                disabled
                className="bg-gray-50"
              />
            )}
          </div>
        )}

        {/* ============================================
            PASO 2: INTEGRATION SELECTOR (OBLIGATORIO)
            ============================================ */}
        <div className="grid gap-2">
          <Label htmlFor="integration_id" className="flex items-center gap-2">
            Integración Origen
            <span className="text-red-500">*</span>
          </Label>
          <select
            id="integration_id"
            value={formData.integration_id}
            onChange={(e) => {
              const integrationId = parseInt(e.target.value) || 0;
              setFormData({ ...formData, integration_id: integrationId });
            }}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100"
            disabled={loadingIntegrations || !formData.business_id}
            required
          >
            <option value="0">Seleccionar Integración</option>
            {filteredIntegrations.map((integration) => (
              <option key={integration.id} value={integration.id}>
                {integration.name} ({integration.type})
              </option>
            ))}
          </select>
          <p className="text-xs text-gray-500">
            La integración que genera los eventos (ej: Shopify - Mi Tiendita)
          </p>
          {!formData.business_id && (
            <p className="text-xs text-amber-600">
              Selecciona un business primero
            </p>
          )}
        </div>

        {/* ============================================
            PASO 3: NOTIFICATION TYPE (RADIO BUTTONS)
            ============================================ */}
        <div className="grid gap-2">
          <Label className="flex items-center gap-2">
            Tipo de Notificación (Canal de Salida)
            <span className="text-red-500">*</span>
          </Label>
          <div className="space-y-2 border rounded-md p-4">
            {loadingTypes ? (
              <p className="text-sm text-gray-500">Cargando tipos...</p>
            ) : notificationTypes.length === 0 ? (
              <p className="text-sm text-gray-500">
                No hay tipos de notificación disponibles
              </p>
            ) : (
              notificationTypes.map((type) => (
                <label
                  key={type.id}
                  className="flex items-start space-x-3 cursor-pointer hover:bg-gray-50 p-2 rounded"
                >
                  <input
                    type="radio"
                    name="notification_type"
                    value={type.id}
                    checked={formData.notification_type_id === type.id}
                    onChange={() =>
                      setFormData({
                        ...formData,
                        notification_type_id: type.id,
                        notification_event_type_id: 0, // Reset event
                      })
                    }
                    className="mt-1"
                  />
                  <div className="flex-1">
                    <div className="font-medium">{type.name}</div>
                    {type.description && (
                      <div className="text-xs text-gray-500">
                        {type.description}
                      </div>
                    )}
                  </div>
                </label>
              ))
            )}
          </div>
          <p className="text-xs text-gray-500">
            Define por qué canal se enviará la notificación
          </p>
        </div>

        {/* ============================================
            PASO 4: EVENT TYPE (DROPDOWN DINÁMICO)
            ============================================ */}
        {formData.notification_type_id > 0 && (
          <div className="grid gap-2">
            <Label htmlFor="event_type" className="flex items-center gap-2">
              Tipo de Evento
              <span className="text-red-500">*</span>
            </Label>
            <select
              id="event_type"
              value={formData.notification_event_type_id}
              onChange={(e) => {
                const eventTypeId = parseInt(e.target.value) || 0;
                setFormData({
                  ...formData,
                  notification_event_type_id: eventTypeId,
                });
              }}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
              required
              disabled={loadingEvents}
            >
              <option value="0">Seleccionar Evento</option>
              {eventTypes.map((event) => (
                <option key={event.id} value={event.id}>
                  {event.event_name} ({event.event_code})
                </option>
              ))}
            </select>
            {loadingEvents && (
              <p className="text-xs text-gray-500">Cargando eventos...</p>
            )}
            {!loadingEvents && eventTypes.length === 0 && (
              <p className="text-xs text-amber-600">
                No hay eventos disponibles para este tipo de notificación
              </p>
            )}
            <p className="text-xs text-gray-500">
              Define qué evento disparará la notificación
            </p>
          </div>
        )}

        {/* ============================================
            PASO 5: ORDER STATUS CHECKLIST
            ============================================ */}
        <div className="grid gap-2">
          <Label>Estados de Orden a Notificar</Label>
          <div className="border rounded-md p-4">
            {loadingOrderStatuses ? (
              <p className="text-sm text-gray-500">Cargando estados...</p>
            ) : orderStatuses.length === 0 ? (
              <p className="text-sm text-gray-500">
                No hay estados de orden disponibles
              </p>
            ) : (
              <div className="grid grid-cols-2 gap-3">
                {orderStatuses.map((status) => (
                  <label
                    key={status.id}
                    className="flex items-center space-x-2 cursor-pointer hover:bg-gray-50 p-2 rounded"
                  >
                    <Checkbox
                      checked={
                        formData.order_status_ids?.includes(status.id) || false
                      }
                      onCheckedChange={() =>
                        handleOrderStatusToggle(status.id)
                      }
                    />
                    <span className="text-sm">
                      {status.name}
                      <span className="text-xs text-gray-500 ml-1">
                        ({status.code})
                      </span>
                    </span>
                  </label>
                ))}
              </div>
            )}
          </div>
          <p className="text-xs text-gray-500">
            Si no seleccionas ninguno, se notificarán TODOS los estados
          </p>
        </div>

        {/* ============================================
            PASO 6: DESCRIPCIÓN Y HABILITADO
            ============================================ */}
        <div className="grid gap-2">
          <Label htmlFor="description">Descripción</Label>
          <Input
            id="description"
            value={formData.description}
            onChange={(e) =>
              setFormData({ ...formData, description: e.target.value })
            }
            placeholder="Descripción opcional de esta configuración"
          />
        </div>

        <div className="flex items-center space-x-2">
          <Checkbox
            id="enabled"
            checked={formData.enabled}
            onCheckedChange={(checked: boolean) =>
              setFormData({ ...formData, enabled: checked })
            }
          />
          <label
            htmlFor="enabled"
            className="text-sm font-medium leading-none cursor-pointer"
          >
            Configuración habilitada
          </label>
        </div>
      </div>

      {/* ============================================
          BOTONES DE ACCIÓN
          ============================================ */}
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
          {loading ? "Guardando..." : config ? "Actualizar" : "Crear"}
        </Button>
      </div>
    </form>
  );
}
