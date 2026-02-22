/**
 * IntegrationConfigForm - Formulario para crear/editar configuración de notificación por integración
 * (Migrado desde services/integrations/notification-config)
 */

'use client';

import { useState } from 'react';
import type {
  CreateNotificationConfigDTO,
  UpdateNotificationConfigDTO,
  IntegrationNotificationConfig,
  TriggerType,
  IntegrationNotifType,
} from '../../domain/integration-types';
import { PaymentMethodSelector } from './PaymentMethodSelector';
import { StatusSelector } from './StatusSelector';
import { IntegrationSourceSelector } from './IntegrationSourceSelector';
import { useWhatsAppTemplates } from '../hooks/useIntegrationNotificationConfigs';
import { TokenStorage } from '@/shared/utils/token-storage';

interface IntegrationConfigFormProps {
  integrationId: number;
  initialData?: IntegrationNotificationConfig;
  onSubmit: (data: CreateNotificationConfigDTO | UpdateNotificationConfigDTO) => Promise<void>;
  onCancel: () => void;
}

const TRIGGER_OPTIONS: { value: TriggerType; label: string }[] = [
  { value: 'order.created', label: 'Orden Creada' },
  { value: 'order.updated', label: 'Orden Actualizada' },
  { value: 'order.status_changed', label: 'Estado de Orden Cambió' },
];

const NOTIFICATION_TYPE_OPTIONS: { value: IntegrationNotifType; label: string }[] = [
  { value: 'whatsapp', label: 'WhatsApp' },
  { value: 'email', label: 'Email' },
  { value: 'sms', label: 'SMS' },
];

export function IntegrationConfigForm({
  integrationId,
  initialData,
  onSubmit,
  onCancel,
}: IntegrationConfigFormProps) {
  const permissions = TokenStorage.getPermissions();
  const businessId = permissions?.business_id || 0;

  const [formData, setFormData] = useState({
    notification_type: (initialData?.notification_type || 'whatsapp') as IntegrationNotifType,
    trigger: (initialData?.conditions.trigger || 'order.created') as TriggerType,
    statuses: initialData?.conditions.statuses || [],
    payment_methods: initialData?.conditions.payment_methods || [],
    source_integration_id: initialData?.conditions.source_integration_id ?? null,
    template_name: initialData?.config.template_name || '',
    recipient_type: initialData?.config.recipient_type || 'customer',
    language: initialData?.config.language || 'es',
    description: initialData?.description || '',
    priority: initialData?.priority || 10,
    is_active: initialData?.is_active ?? true,
  });

  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const { templates, loading: loadingTemplates } = useWhatsAppTemplates(
    formData.notification_type === 'whatsapp' ? integrationId : null
  );

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setError(null);

    try {
      const data = initialData
        ? {
            notification_type: formData.notification_type,
            is_active: formData.is_active,
            conditions: {
              trigger: formData.trigger,
              statuses: formData.statuses,
              payment_methods: formData.payment_methods,
              source_integration_id: formData.source_integration_id,
            },
            config: {
              template_name: formData.template_name,
              recipient_type: formData.recipient_type,
              language: formData.language,
            },
            description: formData.description,
            priority: formData.priority,
          }
        : {
            integration_id: integrationId,
            notification_type: formData.notification_type,
            is_active: formData.is_active,
            conditions: {
              trigger: formData.trigger,
              statuses: formData.statuses,
              payment_methods: formData.payment_methods,
              source_integration_id: formData.source_integration_id,
            },
            config: {
              template_name: formData.template_name,
              recipient_type: formData.recipient_type,
              language: formData.language,
            },
            description: formData.description,
            priority: formData.priority,
          };

      await onSubmit(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al guardar');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6 bg-white p-6 rounded-lg shadow">
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold text-gray-900">
          {initialData ? 'Editar' : 'Nueva'} Configuración de Notificación
        </h2>
        <label className="flex items-center gap-2">
          <span className="text-sm text-gray-700">Activa</span>
          <input
            type="checkbox"
            checked={formData.is_active}
            onChange={(e) => setFormData({ ...formData, is_active: e.target.checked })}
            className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
          />
        </label>
      </div>

      {error && (
        <div className="p-3 bg-red-50 border border-red-200 rounded text-red-700 text-sm">
          {error}
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Tipo de Notificación */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Tipo de Notificación
          </label>
          <select
            value={formData.notification_type}
            onChange={(e) =>
              setFormData({ ...formData, notification_type: e.target.value as IntegrationNotifType })
            }
            disabled={!!initialData}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100"
          >
            {NOTIFICATION_TYPE_OPTIONS.map((opt) => (
              <option key={opt.value} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </select>
        </div>

        {/* Trigger */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Evento Disparador
          </label>
          <select
            value={formData.trigger}
            onChange={(e) => setFormData({ ...formData, trigger: e.target.value as TriggerType })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
          >
            {TRIGGER_OPTIONS.map((opt) => (
              <option key={opt.value} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </select>
        </div>

        {/* Plantilla WhatsApp */}
        {formData.notification_type === 'whatsapp' && (
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Plantilla de WhatsApp
            </label>
            {loadingTemplates ? (
              <div className="w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-50">
                Cargando plantillas...
              </div>
            ) : (
              <select
                value={formData.template_name}
                onChange={(e) => setFormData({ ...formData, template_name: e.target.value })}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
              >
                <option value="">Seleccionar plantilla</option>
                {templates.map((template) => (
                  <option key={template.name} value={template.name}>
                    {template.name} ({template.language})
                  </option>
                ))}
              </select>
            )}
          </div>
        )}

        {/* Prioridad */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Prioridad
          </label>
          <input
            type="number"
            value={formData.priority}
            onChange={(e) => setFormData({ ...formData, priority: parseInt(e.target.value) })}
            min="0"
            max="100"
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
          />
          <p className="mt-1 text-xs text-gray-500">
            Mayor número = mayor prioridad (0-100)
          </p>
        </div>
      </div>

      {/* Descripción */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Descripción
        </label>
        <textarea
          value={formData.description}
          onChange={(e) => setFormData({ ...formData, description: e.target.value })}
          rows={3}
          placeholder="Describe cuándo se debe enviar esta notificación..."
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
        />
      </div>

      {/* Selectores de condiciones */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <StatusSelector
          selectedStatuses={formData.statuses}
          onChange={(statuses) => setFormData({ ...formData, statuses })}
        />

        <PaymentMethodSelector
          selectedMethods={formData.payment_methods}
          onChange={(payment_methods) => setFormData({ ...formData, payment_methods })}
        />
      </div>

      {/* Selector de Integración Origen */}
      <div className="border-t pt-6">
        <IntegrationSourceSelector
          businessId={businessId}
          value={formData.source_integration_id}
          onChange={(value) => setFormData({ ...formData, source_integration_id: value })}
        />
      </div>

      {/* Botones */}
      <div className="flex justify-end gap-3 pt-4 border-t">
        <button
          type="button"
          onClick={onCancel}
          disabled={submitting}
          className="px-4 py-2 text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50"
        >
          Cancelar
        </button>
        <button
          type="submit"
          disabled={submitting}
          className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {submitting ? 'Guardando...' : initialData ? 'Actualizar' : 'Crear'}
        </button>
      </div>
    </form>
  );
}
