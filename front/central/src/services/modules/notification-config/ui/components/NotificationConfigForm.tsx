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
  selectedBusinessId?: number;
}

export function NotificationConfigForm({
  config,
  onSuccess,
  onCancel,
  selectedBusinessId,
}: NotificationConfigFormProps) {
  const [loading, setLoading] = useState(false);
  const { showToast } = useToast();

  // Obtener permisos del usuario actual
  const permissions = TokenStorage.getPermissions();
  const isSuperAdmin = permissions?.is_super || false;
  const userBusinessId = permissions?.business_id || 0;

  // Si viene selectedBusinessId de la página (super admin ya seleccionó), usar ese
  const effectiveBusinessId = selectedBusinessId || userBusinessId;

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
    business_id: effectiveBusinessId,
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
    } else {
      // Auto-select business: from page-level selector or JWT
      setFormData((prev) => ({ ...prev, business_id: effectiveBusinessId }));
    }
  }, [config, effectiveBusinessId]);

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

  // Filtrar integraciones: solo ecommerce activas del business seleccionado
  const filteredIntegrations = integrations.filter(
    (integration) =>
      integration.business_id === formData.business_id &&
      integration.is_active &&
      integration.category === 'ecommerce'
  );

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* Business selector (solo si no viene de la página) */}
      {!config && !selectedBusinessId && (
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
                  integration_id: 0,
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

      {/* Dos columnas */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* ============================================
            COLUMNA IZQUIERDA: Configuración del trigger
            ============================================ */}
        <div className="space-y-5">
          {/* Integración Origen */}
          <div className="grid gap-2">
            <Label htmlFor="integration_id" className="flex items-center gap-2">
              Integración Origen
              <span className="text-red-500">*</span>
            </Label>
            {(() => {
              if (loadingIntegrations) return <p className="text-sm text-gray-500">Cargando integraciones...</p>;
              if (!formData.business_id) return <p className="text-xs text-amber-600">Selecciona un business primero</p>;

              // En modo edición: solo mostrar la integración seleccionada (no editable)
              if (config) {
                const selected = integrations.find(i => i.id === formData.integration_id);
                if (!selected) return <p className="text-sm text-gray-500">Integración no encontrada</p>;
                return (
                  <div className="flex items-center gap-3 px-3 py-2.5 rounded-lg border border-gray-200 bg-gray-50">
                    {selected.image_url ? (
                      <img src={selected.image_url} alt={selected.name} className="w-7 h-7 object-contain rounded shrink-0" />
                    ) : (
                      <div className="w-7 h-7 rounded bg-gray-200 flex items-center justify-center shrink-0">
                        <span className="text-xs font-bold text-gray-500">{selected.type?.charAt(0).toUpperCase() || "?"}</span>
                      </div>
                    )}
                    <div className="min-w-0">
                      <span className="text-sm block truncate font-medium text-gray-900">{selected.name}</span>
                      <span className="text-xs text-gray-400">{selected.type}</span>
                    </div>
                  </div>
                );
              }

              // En modo creación: lista seleccionable
              if (filteredIntegrations.length === 0) return <p className="text-sm text-gray-500">No hay integraciones disponibles</p>;
              return (
                <div className="flex flex-col gap-2">
                  {filteredIntegrations.map((integration) => {
                    const isSelected = formData.integration_id === integration.id;
                    return (
                      <button
                        key={integration.id}
                        type="button"
                        onClick={() => setFormData({ ...formData, integration_id: integration.id })}
                        className={`flex items-center gap-3 px-3 py-2.5 rounded-lg border transition-all text-left ${
                          isSelected
                            ? "bg-blue-50 border-blue-300 ring-1 ring-blue-200"
                            : "bg-white border-gray-200 hover:bg-gray-50"
                        }`}
                      >
                        {integration.image_url ? (
                          <img src={integration.image_url} alt={integration.name} className="w-7 h-7 object-contain rounded shrink-0" />
                        ) : (
                          <div className="w-7 h-7 rounded bg-gray-200 flex items-center justify-center shrink-0">
                            <span className="text-xs font-bold text-gray-500">{integration.type?.charAt(0).toUpperCase() || "?"}</span>
                          </div>
                        )}
                        <div className="min-w-0">
                          <span className={`text-sm block truncate ${isSelected ? "font-medium text-gray-900" : "text-gray-700"}`}>{integration.name}</span>
                          <span className="text-xs text-gray-400">{integration.type}</span>
                        </div>
                      </button>
                    );
                  })}
                </div>
              );
            })()}
            <p className="text-xs text-gray-500">
              La integración que genera los eventos
            </p>
          </div>

        {/* ============================================
            PASO 3: NOTIFICATION TYPE (BOTONES DE COLOR)
            ============================================ */}
        <div className="grid gap-2">
          <Label className="flex items-center gap-2">
            Canal de Notificación
            <span className="text-red-500">*</span>
          </Label>
          <div className="flex flex-wrap gap-2">
            {loadingTypes ? (
              <p className="text-sm text-gray-500">Cargando canales...</p>
            ) : notificationTypes.length === 0 ? (
              <p className="text-sm text-gray-500">
                No hay canales disponibles
              </p>
            ) : (
              notificationTypes.map((type) => {
                const code = type.code?.toLowerCase() || '';
                const isSelected = formData.notification_type_id === type.id;
                const colorMap: Record<string, { bg: string; selectedBg: string; text: string; ring: string }> = {
                  whatsapp: { bg: 'bg-green-50 border-green-200 text-green-700', selectedBg: 'bg-green-500 border-green-600 text-white', text: 'text-green-600', ring: 'ring-green-300' },
                  email: { bg: 'bg-orange-50 border-orange-200 text-orange-700', selectedBg: 'bg-orange-500 border-orange-600 text-white', text: 'text-orange-600', ring: 'ring-orange-300' },
                  sms: { bg: 'bg-purple-50 border-purple-200 text-purple-700', selectedBg: 'bg-purple-500 border-purple-600 text-white', text: 'text-purple-600', ring: 'ring-purple-300' },
                  sse: { bg: 'bg-blue-50 border-blue-200 text-blue-700', selectedBg: 'bg-blue-500 border-blue-600 text-white', text: 'text-blue-600', ring: 'ring-blue-300' },
                };
                const colors = colorMap[code] || colorMap.sse;

                return (
                  <button
                    key={type.id}
                    type="button"
                    onClick={() =>
                      setFormData({
                        ...formData,
                        notification_type_id: type.id,
                        notification_event_type_id: 0,
                      })
                    }
                    className={`
                      px-4 py-2.5 rounded-lg border-2 font-medium text-sm transition-all duration-150
                      focus:outline-none focus:ring-2 ${colors.ring}
                      ${isSelected ? colors.selectedBg + ' shadow-md scale-105' : colors.bg + ' hover:scale-102'}
                    `}
                  >
                    {type.name}
                  </button>
                );
              })
            )}
          </div>
          <p className="text-xs text-gray-500">
            Selecciona el canal por donde se enviará la notificación
          </p>
        </div>

          {/* Tipo de Evento */}
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
                  No hay eventos disponibles para este tipo
                </p>
              )}
              <p className="text-xs text-gray-500">
                Define qué evento disparará la notificación
              </p>
            </div>
          )}

        </div>

        {/* ============================================
            COLUMNA DERECHA: Estados de orden
            ============================================ */}
        <div className="space-y-2">
          <Label>Estados de Orden a Notificar</Label>
          <div className="border rounded-lg max-h-80 overflow-y-auto p-3">
            {loadingOrderStatuses ? (
              <p className="text-sm text-gray-500 p-4">Cargando estados...</p>
            ) : orderStatuses.length === 0 ? (
              <p className="text-sm text-gray-500 p-4">
                No hay estados disponibles
              </p>
            ) : (
              <div className="grid grid-cols-2 gap-2">
                {orderStatuses.map((status) => {
                  const isChecked = (formData.order_status_ids || []).includes(status.id);
                  const statusColor = status.color || "#9CA3AF";
                  const iconByCode: Record<string, React.ReactNode> = {
                    pending: (
                      <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={2}><circle cx="12" cy="12" r="10" /><path d="M12 6v6l4 2" strokeLinecap="round" /></svg>
                    ),
                    processing: (
                      <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={2}><path d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" strokeLinecap="round" strokeLinejoin="round" /></svg>
                    ),
                    on_hold: (
                      <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={2}><path d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z" strokeLinecap="round" strokeLinejoin="round" /></svg>
                    ),
                    shipped: (
                      <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={2}><path d="M13 16V6a1 1 0 00-1-1H4a1 1 0 00-1 1v10m10 0H3m10 0a2 2 0 104 0m-4 0a2 2 0 114 0m6-6v6a1 1 0 01-1 1h-1m-6-1a2 2 0 104 0M9 16H5" strokeLinecap="round" strokeLinejoin="round" /></svg>
                    ),
                    delivered: (
                      <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={2}><path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" strokeLinecap="round" strokeLinejoin="round" /></svg>
                    ),
                    completed: (
                      <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={2}><path d="M5 13l4 4L19 7" strokeLinecap="round" strokeLinejoin="round" /></svg>
                    ),
                    refunded: (
                      <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={2}><path d="M3 10h10a5 5 0 010 10H9m-6-4l4 4m0 0l4-4" strokeLinecap="round" strokeLinejoin="round" /></svg>
                    ),
                    cancelled: (
                      <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={2}><path d="M6 18L18 6M6 6l12 12" strokeLinecap="round" strokeLinejoin="round" /></svg>
                    ),
                    failed: (
                      <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={2}><path d="M12 9v2m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" strokeLinecap="round" strokeLinejoin="round" /></svg>
                    ),
                  };
                  const icon = iconByCode[status.code] || null;
                  return (
                    <button
                      key={status.id}
                      type="button"
                      onClick={() => {
                        setFormData((prev) => {
                          const currentIds = prev.order_status_ids || [];
                          const newIds = isChecked
                            ? currentIds.filter((id) => id !== status.id)
                            : [...currentIds, status.id];
                          return { ...prev, order_status_ids: newIds };
                        });
                      }}
                      className={`flex items-center justify-between gap-2 px-3 py-2 rounded-lg border transition-colors ${isChecked ? "bg-blue-50 border-blue-200" : "bg-white border-gray-200 hover:bg-gray-50"}`}
                    >
                      <div className="flex items-center gap-2 min-w-0">
                        <span
                          className="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium truncate"
                          style={{ backgroundColor: statusColor + "20", color: statusColor }}
                        >
                          {icon}
                          {status.name}
                        </span>
                      </div>
                      {/* Toggle switch */}
                      <div className={`relative w-9 h-5 rounded-full shrink-0 transition-colors ${isChecked ? "bg-blue-500" : "bg-gray-300"}`}>
                        <div className={`absolute top-0.5 w-4 h-4 rounded-full bg-white shadow transition-transform ${isChecked ? "translate-x-4" : "translate-x-0.5"}`} />
                      </div>
                    </button>
                  );
                })}
              </div>
            )}
          </div>
          <p className="text-xs text-gray-500">
            Marca los estados que activarán esta notificación
          </p>
        </div>
      </div>

      {/* Descripción - ancho completo */}
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

      {/* ============================================
          BOTONES DE ACCIÓN
          ============================================ */}
      <div className="flex items-center justify-between pt-4 border-t">
        <button
          type="button"
          onClick={() => setFormData({ ...formData, enabled: !formData.enabled })}
          className={`flex items-center gap-2 px-4 py-2 rounded-lg border text-sm font-medium transition-colors ${
            formData.enabled
              ? "bg-green-50 border-green-300 text-green-700 hover:bg-green-100"
              : "bg-red-50 border-red-300 text-red-700 hover:bg-red-100"
          }`}
        >
          <div className={`relative w-9 h-5 rounded-full shrink-0 transition-colors ${formData.enabled ? "bg-green-500" : "bg-red-400"}`}>
            <div className={`absolute top-0.5 w-4 h-4 rounded-full bg-white shadow transition-transform ${formData.enabled ? "translate-x-4" : "translate-x-0.5"}`} />
          </div>
          {formData.enabled ? "Activa" : "Inactiva"}
        </button>
        <div className="flex gap-2">
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
      </div>
    </form>
  );
}
