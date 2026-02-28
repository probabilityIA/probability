"use client";

import { useState, useEffect } from "react";
import {
  NotificationType,
  NotificationEventType,
} from "../../domain/types";
import { Label } from "@/shared/ui/label";
import { Input } from "@/shared/ui/input";
import {
  getNotificationTypesAction,
  getNotificationEventTypesAction,
} from "../../infra/actions";

export interface LocalRule {
  _tempId: string;
  id?: number;
  notification_type_id: number;
  notification_event_type_id: number;
  enabled: boolean;
  description: string;
  order_status_ids: number[];
  _deleted: boolean;
}

interface OrderStatusOption {
  id: number;
  name: string;
  code: string;
  color?: string;
  is_active: boolean;
}

interface RuleCardProps {
  rule: LocalRule;
  index: number;
  orderStatuses: OrderStatusOption[];
  onChange: (updated: LocalRule) => void;
  onDelete: () => void;
}

const CHANNEL_COLORS: Record<string, { bg: string; selectedBg: string; ring: string }> = {
  whatsapp: { bg: "bg-green-50 border-green-200 text-green-700", selectedBg: "bg-green-500 border-green-600 text-white", ring: "ring-green-300" },
  email: { bg: "bg-orange-50 border-orange-200 text-orange-700", selectedBg: "bg-orange-500 border-orange-600 text-white", ring: "ring-orange-300" },
  sms: { bg: "bg-purple-50 border-purple-200 text-purple-700", selectedBg: "bg-purple-500 border-purple-600 text-white", ring: "ring-purple-300" },
  sse: { bg: "bg-blue-50 border-blue-200 text-blue-700", selectedBg: "bg-blue-500 border-blue-600 text-white", ring: "ring-blue-300" },
};

export function RuleCard({ rule, index, orderStatuses, onChange, onDelete }: RuleCardProps) {
  const [notificationTypes, setNotificationTypes] = useState<NotificationType[]>([]);
  const [eventTypes, setEventTypes] = useState<NotificationEventType[]>([]);
  const [loadingTypes, setLoadingTypes] = useState(false);
  const [loadingEvents, setLoadingEvents] = useState(false);

  // Load notification types on mount
  useEffect(() => {
    const load = async () => {
      setLoadingTypes(true);
      try {
        const result = await getNotificationTypesAction();
        if (result.success) setNotificationTypes(result.data);
      } finally {
        setLoadingTypes(false);
      }
    };
    load();
  }, []);

  // Load event types when notification_type_id changes
  useEffect(() => {
    if (rule.notification_type_id > 0) {
      const load = async () => {
        setLoadingEvents(true);
        try {
          const result = await getNotificationEventTypesAction(rule.notification_type_id);
          if (result.success) setEventTypes(result.data);
          else setEventTypes([]);
        } finally {
          setLoadingEvents(false);
        }
      };
      load();
    } else {
      setEventTypes([]);
    }
  }, [rule.notification_type_id]);

  // Get allowed statuses for selected event type
  const selectedEvent = eventTypes.find((e) => e.id === rule.notification_event_type_id);
  const allowedStatusIds = selectedEvent?.allowed_order_status_ids;
  const filteredStatuses =
    allowedStatusIds && allowedStatusIds.length > 0
      ? orderStatuses.filter((s) => allowedStatusIds.includes(s.id))
      : orderStatuses;

  if (rule._deleted) return null;

  return (
    <div className="border border-gray-200 rounded-lg p-4 bg-white relative group">
      {/* Header row: index + enabled toggle + delete */}
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-3">
          <span className="text-xs font-bold text-gray-400 bg-gray-100 rounded-full w-6 h-6 flex items-center justify-center">
            {index + 1}
          </span>
          <button
            type="button"
            onClick={() => onChange({ ...rule, enabled: !rule.enabled })}
            className={`flex items-center gap-2 px-3 py-1 rounded-full text-xs font-medium transition-colors ${
              rule.enabled
                ? "bg-green-50 text-green-700 border border-green-200"
                : "bg-gray-50 text-gray-500 border border-gray-200"
            }`}
          >
            <div
              className={`relative w-7 h-4 rounded-full shrink-0 transition-colors ${
                rule.enabled ? "bg-green-500" : "bg-gray-300"
              }`}
            >
              <div
                className={`absolute top-0.5 w-3 h-3 rounded-full bg-white shadow transition-transform ${
                  rule.enabled ? "translate-x-3" : "translate-x-0.5"
                }`}
              />
            </div>
            {rule.enabled ? "Activa" : "Inactiva"}
          </button>
        </div>
        <button
          type="button"
          onClick={onDelete}
          className="text-gray-400 hover:text-red-500 transition-colors p-1"
          title="Eliminar regla"
        >
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </button>
      </div>

      {/* Canal (notification type) - color buttons */}
      <div className="mb-3">
        <Label className="text-xs text-gray-500 mb-1.5 block">Canal</Label>
        <div className="flex flex-wrap gap-1.5">
          {loadingTypes ? (
            <span className="text-xs text-gray-400">Cargando...</span>
          ) : (
            notificationTypes.map((type) => {
              const code = type.code?.toLowerCase() || "";
              const isSelected = rule.notification_type_id === type.id;
              const colors = CHANNEL_COLORS[code] || CHANNEL_COLORS.sse;
              return (
                <button
                  key={type.id}
                  type="button"
                  onClick={() =>
                    onChange({
                      ...rule,
                      notification_type_id: type.id,
                      notification_event_type_id: 0,
                      order_status_ids: [],
                    })
                  }
                  className={`px-3 py-1.5 rounded-md border text-xs font-medium transition-all ${
                    isSelected ? colors.selectedBg + " shadow-sm" : colors.bg + " hover:opacity-80"
                  }`}
                >
                  {type.name}
                </button>
              );
            })
          )}
        </div>
      </div>

      {/* Evento (event type) - dropdown */}
      {rule.notification_type_id > 0 && (
        <div className="mb-3">
          <Label className="text-xs text-gray-500 mb-1.5 block">Evento</Label>
          <select
            value={rule.notification_event_type_id}
            onChange={(e) => {
              const eventTypeId = parseInt(e.target.value) || 0;
              onChange({
                ...rule,
                notification_event_type_id: eventTypeId,
                order_status_ids: [],
              });
            }}
            className="w-full px-2.5 py-1.5 border border-gray-300 rounded-md text-sm focus:ring-blue-500 focus:border-blue-500"
            disabled={loadingEvents}
          >
            <option value="0">Seleccionar evento</option>
            {eventTypes.map((event) => (
              <option key={event.id} value={event.id}>
                {event.event_name} ({event.event_code})
              </option>
            ))}
          </select>
        </div>
      )}

      {/* Estados de orden (filtrados por evento) */}
      {rule.notification_event_type_id > 0 && filteredStatuses.length > 0 && (
        <div className="mb-3">
          <Label className="text-xs text-gray-500 mb-1.5 block">
            Estados a notificar
            {allowedStatusIds && allowedStatusIds.length > 0 && (
              <span className="text-xs text-blue-500 ml-1">(filtrados por evento)</span>
            )}
          </Label>
          <div className="flex flex-wrap gap-1.5">
            {filteredStatuses.map((status) => {
              const isChecked = rule.order_status_ids.includes(status.id);
              const statusColor = status.color || "#9CA3AF";
              return (
                <button
                  key={status.id}
                  type="button"
                  onClick={() => {
                    const newIds = isChecked
                      ? rule.order_status_ids.filter((id) => id !== status.id)
                      : [...rule.order_status_ids, status.id];
                    onChange({ ...rule, order_status_ids: newIds });
                  }}
                  className={`flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium border transition-colors ${
                    isChecked
                      ? "border-blue-300 bg-blue-50"
                      : "border-gray-200 bg-white hover:bg-gray-50"
                  }`}
                >
                  <span
                    className="w-2 h-2 rounded-full shrink-0"
                    style={{ backgroundColor: statusColor }}
                  />
                  {status.name}
                  {isChecked && (
                    <svg className="w-3 h-3 text-blue-500" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                    </svg>
                  )}
                </button>
              );
            })}
          </div>
        </div>
      )}

      {/* Descripción */}
      <div>
        <Label className="text-xs text-gray-500 mb-1.5 block">Descripción</Label>
        <Input
          value={rule.description}
          onChange={(e) => onChange({ ...rule, description: e.target.value })}
          placeholder="Descripción opcional"
          className="text-sm h-8"
        />
      </div>
    </div>
  );
}
