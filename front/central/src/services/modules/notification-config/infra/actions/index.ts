// ============================================
// RE-EXPORTS - Notification Types & Event Types
// NOTA: Este archivo NO puede tener "use server" porque Next.js no permite re-exports
//       Los archivos notification-types.ts y notification-event-types.ts SÍ tienen "use server"
//       Importa directamente desde esos archivos si necesitas las actions
// ============================================

// Re-export notification types actions
export {
  getNotificationTypesAction,
  getNotificationTypeByIdAction,
  createNotificationTypeAction,
  updateNotificationTypeAction,
  deleteNotificationTypeAction,
} from "./notification-types";

// Re-export notification event types actions
export {
  getNotificationEventTypesAction,
  getNotificationEventTypeByIdAction,
  createNotificationEventTypeAction,
  updateNotificationEventTypeAction,
  deleteNotificationEventTypeAction,
  toggleNotificationEventTypeActiveAction,
} from "./notification-event-types";

// Re-export notification config actions
export {
  createConfigAction,
  updateConfigAction,
  deleteConfigAction,
  listConfigsAction,
  getConfigsAction,
  getConfigAction,
} from "./notification-configs";
