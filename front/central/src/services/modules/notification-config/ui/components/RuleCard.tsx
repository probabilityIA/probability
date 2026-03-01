"use client";

import { useState, useEffect } from "react";
import {
  NotificationType,
  NotificationEventType,
} from "../../domain/types";
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

const CHANNEL_COLORS: Record<string, { bg: string; selectedBg: string }> = {
  whatsapp: { bg: "bg-green-50 border-green-200 text-green-700", selectedBg: "bg-green-500 border-green-600 text-white" },
  email: { bg: "bg-orange-50 border-orange-200 text-orange-700", selectedBg: "bg-orange-500 border-orange-600 text-white" },
  sms: { bg: "bg-purple-50 border-purple-200 text-purple-700", selectedBg: "bg-purple-500 border-purple-600 text-white" },
  sse: { bg: "bg-blue-50 border-blue-200 text-blue-700", selectedBg: "bg-blue-500 border-blue-600 text-white" },
};

const CHANNEL_BADGE: Record<string, string> = {
  whatsapp: "bg-green-100 text-green-700",
  email: "bg-orange-100 text-orange-700",
  sms: "bg-purple-100 text-purple-700",
  sse: "bg-blue-100 text-blue-700",
};

export function RuleCard({ rule, index, orderStatuses, onChange, onDelete }: RuleCardProps) {
  const [notificationTypes, setNotificationTypes] = useState<NotificationType[]>([]);
  const [eventTypes, setEventTypes] = useState<NotificationEventType[]>([]);
  const [loadingTypes, setLoadingTypes] = useState(false);
  const [loadingEvents, setLoadingEvents] = useState(false);

  const isNew = !rule.id;

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

  const selectedEvent = eventTypes.find((e) => e.id === rule.notification_event_type_id);
  const allowedStatusIds = selectedEvent?.allowed_order_status_ids;
  const filteredStatuses =
    allowedStatusIds && allowedStatusIds.length > 0
      ? orderStatuses.filter((s) => allowedStatusIds.includes(s.id))
      : orderStatuses;

  if (rule._deleted) return null;

  const selectedType = notificationTypes.find((t) => t.id === rule.notification_type_id);
  const channelCode = selectedType?.code?.toLowerCase() || "";
  const channelName = selectedType?.name || "";
  const eventName = selectedEvent?.event_name || "";

  return (
    <tr className="border-b border-gray-100 hover:bg-gray-50/50 transition-colors align-top">
      {/* # */}
      <td className="py-3 px-3 text-center">
        <span className="text-[10px] font-bold text-gray-400 bg-gray-100 rounded-full w-5 h-5 inline-flex items-center justify-center">
          {index + 1}
        </span>
      </td>

      {/* Canal */}
      <td className="py-3 px-3">
        {isNew ? (
          <div className="flex flex-wrap gap-1">
            {loadingTypes ? (
              <span className="text-xs text-gray-400">...</span>
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
                    className={`px-2 py-0.5 rounded border text-[11px] font-medium transition-all ${
                      isSelected ? colors.selectedBg + " shadow-sm" : colors.bg + " hover:opacity-80"
                    }`}
                  >
                    {type.name}
                  </button>
                );
              })
            )}
          </div>
        ) : (
          <span className={`inline-block px-2 py-0.5 rounded text-[11px] font-semibold ${CHANNEL_BADGE[channelCode] || "bg-gray-100 text-gray-600"}`}>
            {channelName || "—"}
          </span>
        )}
      </td>

      {/* Evento */}
      <td className="py-3 px-3">
        {isNew ? (
          rule.notification_type_id > 0 ? (
            <select
              value={rule.notification_event_type_id}
              onChange={(e) => {
                const eventTypeId = parseInt(e.target.value) || 0;
                onChange({ ...rule, notification_event_type_id: eventTypeId, order_status_ids: [] });
              }}
              className="w-full px-2 py-1 border border-gray-300 rounded text-xs focus:ring-blue-500 focus:border-blue-500"
              disabled={loadingEvents}
            >
              <option value="0">Seleccionar...</option>
              {eventTypes.map((event) => (
                <option key={event.id} value={event.id}>
                  {event.event_name}
                </option>
              ))}
            </select>
          ) : (
            <span className="text-xs text-gray-300 italic">Selecciona canal</span>
          )
        ) : (
          <span className="text-xs text-gray-700">{eventName || "—"}</span>
        )}
      </td>

      {/* Estados */}
      <td className="py-3 px-3">
        {(isNew && rule.notification_event_type_id === 0) ? (
          <span className="text-[10px] text-gray-300 italic">
            {isNew ? "Selecciona evento" : "Sin estados"}
          </span>
        ) : filteredStatuses.length > 0 ? (
          <div className="flex flex-wrap gap-1">
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
                  className={`inline-flex items-center gap-1 px-1.5 py-0.5 rounded-full text-[10px] font-medium border transition-colors cursor-pointer ${
                    isChecked ? "border-blue-300 bg-blue-50 text-blue-700" : "border-gray-200 bg-white text-gray-500 hover:bg-gray-50"
                  }`}
                >
                  <span className="w-1.5 h-1.5 rounded-full" style={{ backgroundColor: statusColor }} />
                  {status.name}
                  {isChecked && (
                    <svg className="w-2.5 h-2.5 text-blue-500" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                    </svg>
                  )}
                </button>
              );
            })}
          </div>
        ) : (
          <span className="text-[10px] text-gray-300 italic">Sin estados disponibles</span>
        )}
      </td>

      {/* Activo */}
      <td className="py-3 px-3 text-center">
        <button
          type="button"
          onClick={() => onChange({ ...rule, enabled: !rule.enabled })}
          className="inline-flex items-center gap-1"
          title={rule.enabled ? "Desactivar" : "Activar"}
        >
          <div className={`relative w-8 h-[18px] rounded-full transition-colors ${rule.enabled ? "bg-green-500" : "bg-gray-300"}`}>
            <div className={`absolute top-[2px] w-[14px] h-[14px] rounded-full bg-white shadow transition-transform ${rule.enabled ? "translate-x-[14px]" : "translate-x-[2px]"}`} />
          </div>
        </button>
      </td>

      {/* Eliminar */}
      <td className="py-3 px-3 text-center">
        <button
          type="button"
          onClick={onDelete}
          className="p-1.5 rounded-md bg-red-50 text-red-500 hover:bg-red-100 transition-colors"
          title="Eliminar"
        >
          <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </button>
      </td>
    </tr>
  );
}
