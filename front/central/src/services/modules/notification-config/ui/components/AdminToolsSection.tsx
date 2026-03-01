'use client';

import { AccordionItem } from '@/shared/ui/accordion';
import { NotificationTypeList } from './NotificationTypeList';
import { NotificationEventTypeList } from './NotificationEventTypeList';
import type { NotificationType, NotificationEventType } from '../../domain/types';

interface AdminToolsSectionProps {
  onCreateChannel: () => void;
  onEditChannel: (channel: NotificationType) => void;
  channelRefreshKey: number;
  onCreateEventType: () => void;
  onEditEventType: (eventType: NotificationEventType) => void;
  eventTypeRefreshKey: number;
}

export function AdminToolsSection({
  onCreateChannel,
  onEditChannel,
  channelRefreshKey,
  onCreateEventType,
  onEditEventType,
  eventTypeRefreshKey,
}: AdminToolsSectionProps) {
  return (
    <div className="space-y-3">
      <h3 className="text-sm font-medium text-gray-500 uppercase tracking-wider">
        Herramientas de Administración
      </h3>

      <AccordionItem title="Canales de Notificación (WhatsApp, Email, SMS, SSE)">
        <NotificationTypeList
          onEdit={onEditChannel}
          onCreate={onCreateChannel}
          refreshKey={channelRefreshKey}
        />
      </AccordionItem>

      <AccordionItem title="Tipos de Eventos">
        <NotificationEventTypeList
          onEdit={onEditEventType}
          onCreate={onCreateEventType}
          refreshKey={eventTypeRefreshKey}
        />
      </AccordionItem>
    </div>
  );
}
