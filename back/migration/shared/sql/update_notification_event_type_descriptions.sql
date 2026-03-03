-- ============================================
-- ACTUALIZAR DESCRIPCIONES DE EVENTOS DE NOTIFICACION
-- ============================================
-- Agrega descripciones significativas a los notification_event_types
-- para que los usuarios entiendan que hace cada evento.
-- ============================================

-- SSE: Nueva Orden
UPDATE notification_event_types
SET description = 'Notificacion en tiempo real cuando se recibe una nueva orden desde la integracion. Aparece como alerta instantanea en el panel.'
WHERE id = 1 AND event_code = 'order.created' AND notification_type_id = 1;

-- SSE: Cambio de Estado
UPDATE notification_event_types
SET description = 'Notificacion en tiempo real cuando cambia el estado de una orden (ej: pendiente a enviado). Se muestra como alerta en el panel.'
WHERE id = 2 AND event_code = 'order.status_changed' AND notification_type_id = 1;

-- WhatsApp: Confirmacion de Pedido
UPDATE notification_event_types
SET description = 'Envia un mensaje de WhatsApp al cliente confirmando que su pedido fue recibido y esta siendo procesado.'
WHERE id = 3 AND event_code = 'order.created' AND notification_type_id = 2;

-- WhatsApp: Pedido Enviado
UPDATE notification_event_types
SET description = 'Envia un mensaje de WhatsApp al cliente notificando que su pedido ha sido despachado y esta en camino.'
WHERE id = 4 AND event_code = 'order.shipped' AND notification_type_id = 2;

-- WhatsApp: Pedido Entregado
UPDATE notification_event_types
SET description = 'Envia un mensaje de WhatsApp al cliente confirmando que su pedido fue entregado exitosamente.'
WHERE id = 5 AND event_code = 'order.delivered' AND notification_type_id = 2;

-- WhatsApp: Pedido Cancelado
UPDATE notification_event_types
SET description = 'Envia un mensaje de WhatsApp al cliente informando que su pedido ha sido cancelado.'
WHERE id = 6 AND event_code = 'order.canceled' AND notification_type_id = 2;

-- WhatsApp: Factura Generada
UPDATE notification_event_types
SET description = 'Envia un mensaje de WhatsApp al cliente con la informacion de la factura electronica generada para su pedido.'
WHERE id = 7 AND event_code = 'invoice.created' AND notification_type_id = 2;

-- ============================================
-- VERIFICAR
-- ============================================
-- SELECT id, event_code, event_name, description FROM notification_event_types ORDER BY id;
