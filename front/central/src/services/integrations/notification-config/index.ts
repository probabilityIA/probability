/**
 * Notification Config Module - Public API
 * Exporta todo lo necesario para usar el módulo desde fuera
 */

// Domain
export type {
  IntegrationNotificationConfig,
  CreateNotificationConfigDTO,
  UpdateNotificationConfigDTO,
  FilterNotificationConfigDTO,
  NotificationConditions,
  NotificationConfig,
  NotificationType,
  TriggerType,
  RecipientType,
  PaymentMethod,
  OrderStatus,
  WhatsAppTemplate,
} from './domain/types';

// Hooks
export {
  useNotificationConfigs,
  useNotificationConfig,
  usePaymentMethods,
  useOrderStatuses,
  useWhatsAppTemplates,
} from './ui/hooks/useNotificationConfigs';

// Components
export { ConfigList } from './ui/components/ConfigList';
export { ConfigForm } from './ui/components/ConfigForm';
export { PaymentMethodSelector } from './ui/components/PaymentMethodSelector';
export { StatusSelector } from './ui/components/StatusSelector';
