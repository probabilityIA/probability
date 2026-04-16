import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/notification_config/domain/entities.dart';

void main() {
  // =========================================================================
  // NotificationType
  // =========================================================================
  group('NotificationType', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'name': 'WhatsApp',
        'code': 'whatsapp',
        'description': 'WhatsApp notifications',
        'icon': 'chat',
        'is_active': true,
        'config_schema': {'template': 'string'},
        'created_at': '2026-01-01',
        'updated_at': '2026-01-02',
      };

      final type = NotificationType.fromJson(json);

      expect(type.id, 1);
      expect(type.name, 'WhatsApp');
      expect(type.code, 'whatsapp');
      expect(type.description, 'WhatsApp notifications');
      expect(type.icon, 'chat');
      expect(type.isActive, true);
      expect(type.configSchema, {'template': 'string'});
      expect(type.createdAt, '2026-01-01');
      expect(type.updatedAt, '2026-01-02');
    });

    test('fromJson handles defaults for missing fields', () {
      final json = <String, dynamic>{};

      final type = NotificationType.fromJson(json);

      expect(type.id, 0);
      expect(type.name, '');
      expect(type.code, '');
      expect(type.isActive, true);
      expect(type.createdAt, '');
      expect(type.updatedAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'name': 'Email',
        'code': 'email',
        'is_active': true,
        'created_at': '',
        'updated_at': '',
      };

      final type = NotificationType.fromJson(json);

      expect(type.description, isNull);
      expect(type.icon, isNull);
      expect(type.configSchema, isNull);
    });
  });

  // =========================================================================
  // NotificationEventType
  // =========================================================================
  group('NotificationEventType', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 10,
        'notification_type_id': 1,
        'event_code': 'order_created',
        'event_name': 'Order Created',
        'description': 'When an order is created',
        'template_config': {'body': 'Hello {{name}}'},
        'is_active': true,
        'allowed_order_status_ids': [1, 2, 3],
        'created_at': '2026-01-01',
        'updated_at': '2026-01-02',
        'notification_type': {
          'id': 1,
          'name': 'WhatsApp',
          'code': 'whatsapp',
          'is_active': true,
          'created_at': '',
          'updated_at': '',
        },
      };

      final event = NotificationEventType.fromJson(json);

      expect(event.id, 10);
      expect(event.notificationTypeId, 1);
      expect(event.eventCode, 'order_created');
      expect(event.eventName, 'Order Created');
      expect(event.description, 'When an order is created');
      expect(event.templateConfig, {'body': 'Hello {{name}}'});
      expect(event.isActive, true);
      expect(event.allowedOrderStatusIds, [1, 2, 3]);
      expect(event.createdAt, '2026-01-01');
      expect(event.updatedAt, '2026-01-02');
      expect(event.notificationType, isNotNull);
      expect(event.notificationType!.name, 'WhatsApp');
    });

    test('fromJson handles defaults for missing fields', () {
      final json = <String, dynamic>{};

      final event = NotificationEventType.fromJson(json);

      expect(event.id, 0);
      expect(event.notificationTypeId, 0);
      expect(event.eventCode, '');
      expect(event.eventName, '');
      expect(event.isActive, true);
      expect(event.createdAt, '');
      expect(event.updatedAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'notification_type_id': 1,
        'event_code': 'test',
        'event_name': 'Test',
        'is_active': true,
        'created_at': '',
        'updated_at': '',
      };

      final event = NotificationEventType.fromJson(json);

      expect(event.description, isNull);
      expect(event.templateConfig, isNull);
      expect(event.allowedOrderStatusIds, isNull);
      expect(event.notificationType, isNull);
    });
  });

  // =========================================================================
  // NotificationConfig
  // =========================================================================
  group('NotificationConfig', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'business_id': 5,
        'integration_id': 10,
        'notification_type_id': 1,
        'notification_event_type_id': 2,
        'enabled': true,
        'filters': {'min_amount': 1000},
        'description': 'Test config',
        'created_at': '2026-01-01',
        'updated_at': '2026-01-02',
        'deleted_at': null,
        'notification_type': {
          'id': 1,
          'name': 'WhatsApp',
          'code': 'whatsapp',
          'is_active': true,
          'created_at': '',
          'updated_at': '',
        },
        'notification_event_type': {
          'id': 2,
          'notification_type_id': 1,
          'event_code': 'order_created',
          'event_name': 'Order Created',
          'is_active': true,
          'created_at': '',
          'updated_at': '',
        },
        'order_status_ids': [1, 2, 3],
        'notification_type_name': 'WhatsApp',
        'notification_event_name': 'Order Created',
        'event_type': 'order_created',
        'channels': ['whatsapp', 'email'],
      };

      final config = NotificationConfig.fromJson(json);

      expect(config.id, 1);
      expect(config.businessId, 5);
      expect(config.integrationId, 10);
      expect(config.notificationTypeId, 1);
      expect(config.notificationEventTypeId, 2);
      expect(config.enabled, true);
      expect(config.filters, {'min_amount': 1000});
      expect(config.description, 'Test config');
      expect(config.createdAt, '2026-01-01');
      expect(config.updatedAt, '2026-01-02');
      expect(config.deletedAt, isNull);
      expect(config.notificationType, isNotNull);
      expect(config.notificationType!.name, 'WhatsApp');
      expect(config.notificationEventType, isNotNull);
      expect(config.notificationEventType!.eventCode, 'order_created');
      expect(config.orderStatusIds, [1, 2, 3]);
      expect(config.notificationTypeName, 'WhatsApp');
      expect(config.notificationEventName, 'Order Created');
      expect(config.eventType, 'order_created');
      expect(config.channels, ['whatsapp', 'email']);
    });

    test('fromJson handles defaults for missing fields', () {
      final json = <String, dynamic>{};

      final config = NotificationConfig.fromJson(json);

      expect(config.id, 0);
      expect(config.businessId, 0);
      expect(config.integrationId, 0);
      expect(config.notificationTypeId, 0);
      expect(config.notificationEventTypeId, 0);
      expect(config.enabled, false);
      expect(config.createdAt, '');
      expect(config.updatedAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'business_id': 1,
        'integration_id': 1,
        'notification_type_id': 1,
        'notification_event_type_id': 1,
        'enabled': true,
        'created_at': '',
        'updated_at': '',
      };

      final config = NotificationConfig.fromJson(json);

      expect(config.filters, isNull);
      expect(config.description, isNull);
      expect(config.deletedAt, isNull);
      expect(config.notificationType, isNull);
      expect(config.notificationEventType, isNull);
      expect(config.orderStatusIds, isNull);
      expect(config.notificationTypeName, isNull);
      expect(config.notificationEventName, isNull);
      expect(config.eventType, isNull);
      expect(config.channels, isNull);
    });
  });

  // =========================================================================
  // CreateConfigDTO (notification_config)
  // =========================================================================
  group('CreateConfigDTO', () {
    test('toJson includes all required fields', () {
      final dto = CreateConfigDTO(
        businessId: 1,
        integrationId: 2,
        notificationTypeId: 3,
        notificationEventTypeId: 4,
      );

      final json = dto.toJson();

      expect(json['business_id'], 1);
      expect(json['integration_id'], 2);
      expect(json['notification_type_id'], 3);
      expect(json['notification_event_type_id'], 4);
    });

    test('toJson includes optional fields when set', () {
      final dto = CreateConfigDTO(
        businessId: 1,
        integrationId: 2,
        notificationTypeId: 3,
        notificationEventTypeId: 4,
        enabled: true,
        filters: {'min': 100},
        description: 'Test',
        orderStatusIds: [1, 2],
      );

      final json = dto.toJson();

      expect(json['enabled'], true);
      expect(json['filters'], {'min': 100});
      expect(json['description'], 'Test');
      expect(json['order_status_ids'], [1, 2]);
    });

    test('toJson omits null optional fields', () {
      final dto = CreateConfigDTO(
        businessId: 1,
        integrationId: 2,
        notificationTypeId: 3,
        notificationEventTypeId: 4,
      );

      final json = dto.toJson();

      expect(json.containsKey('enabled'), false);
      expect(json.containsKey('filters'), false);
      expect(json.containsKey('description'), false);
      expect(json.containsKey('order_status_ids'), false);
    });
  });

  // =========================================================================
  // UpdateConfigDTO (notification_config)
  // =========================================================================
  group('UpdateConfigDTO', () {
    test('toJson includes all set fields', () {
      final dto = UpdateConfigDTO(
        integrationId: 5,
        notificationTypeId: 6,
        notificationEventTypeId: 7,
        enabled: false,
        filters: {'key': 'val'},
        description: 'Updated',
        orderStatusIds: [3, 4],
      );

      final json = dto.toJson();

      expect(json['integration_id'], 5);
      expect(json['notification_type_id'], 6);
      expect(json['notification_event_type_id'], 7);
      expect(json['enabled'], false);
      expect(json['filters'], {'key': 'val'});
      expect(json['description'], 'Updated');
      expect(json['order_status_ids'], [3, 4]);
    });

    test('toJson returns empty map when all fields null', () {
      final dto = UpdateConfigDTO();
      final json = dto.toJson();
      expect(json, isEmpty);
    });

    test('toJson includes only provided fields', () {
      final dto = UpdateConfigDTO(enabled: true);
      final json = dto.toJson();

      expect(json.length, 1);
      expect(json['enabled'], true);
    });
  });

  // =========================================================================
  // ConfigFilter
  // =========================================================================
  group('ConfigFilter', () {
    test('toQueryParams includes all set fields', () {
      final filter = ConfigFilter(
        businessId: 1,
        integrationId: 2,
        notificationTypeId: 3,
        notificationEventTypeId: 4,
      );

      final query = filter.toQueryParams();

      expect(query['business_id'], 1);
      expect(query['integration_id'], 2);
      expect(query['notification_type_id'], 3);
      expect(query['notification_event_type_id'], 4);
    });

    test('toQueryParams omits null fields', () {
      final filter = ConfigFilter();
      final query = filter.toQueryParams();
      expect(query, isEmpty);
    });

    test('toQueryParams includes only provided fields', () {
      final filter = ConfigFilter(businessId: 5);
      final query = filter.toQueryParams();

      expect(query.length, 1);
      expect(query['business_id'], 5);
    });
  });

  // =========================================================================
  // SyncRule
  // =========================================================================
  group('SyncRule', () {
    test('toJson includes all required fields', () {
      final rule = SyncRule(
        notificationTypeId: 1,
        notificationEventTypeId: 2,
        enabled: true,
        description: 'Test rule',
        orderStatusIds: [1, 2, 3],
      );

      final json = rule.toJson();

      expect(json['notification_type_id'], 1);
      expect(json['notification_event_type_id'], 2);
      expect(json['enabled'], true);
      expect(json['description'], 'Test rule');
      expect(json['order_status_ids'], [1, 2, 3]);
      expect(json.containsKey('id'), false);
    });

    test('toJson includes id when set', () {
      final rule = SyncRule(
        id: 42,
        notificationTypeId: 1,
        notificationEventTypeId: 2,
        enabled: true,
        description: 'Test',
        orderStatusIds: [],
      );

      final json = rule.toJson();

      expect(json['id'], 42);
    });

    test('toJson omits null id', () {
      final rule = SyncRule(
        notificationTypeId: 1,
        notificationEventTypeId: 2,
        enabled: false,
        description: '',
        orderStatusIds: [],
      );

      final json = rule.toJson();

      expect(json.containsKey('id'), false);
    });
  });

  // =========================================================================
  // SyncConfigsDTO
  // =========================================================================
  group('SyncConfigsDTO', () {
    test('toJson includes integrationId and rules', () {
      final dto = SyncConfigsDTO(
        integrationId: 5,
        rules: [
          SyncRule(
            notificationTypeId: 1,
            notificationEventTypeId: 2,
            enabled: true,
            description: 'Rule 1',
            orderStatusIds: [1],
          ),
          SyncRule(
            notificationTypeId: 3,
            notificationEventTypeId: 4,
            enabled: false,
            description: 'Rule 2',
            orderStatusIds: [2, 3],
          ),
        ],
      );

      final json = dto.toJson();

      expect(json['integration_id'], 5);
      expect(json['rules'], isA<List>());
      expect((json['rules'] as List).length, 2);
    });
  });

  // =========================================================================
  // SyncConfigsResponse
  // =========================================================================
  group('SyncConfigsResponse', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'created': 2,
        'updated': 1,
        'deleted': 0,
        'configs': [
          {
            'id': 1,
            'business_id': 1,
            'integration_id': 1,
            'notification_type_id': 1,
            'notification_event_type_id': 1,
            'enabled': true,
            'created_at': '',
            'updated_at': '',
          },
        ],
      };

      final response = SyncConfigsResponse.fromJson(json);

      expect(response.created, 2);
      expect(response.updated, 1);
      expect(response.deleted, 0);
      expect(response.configs.length, 1);
      expect(response.configs.first.id, 1);
    });

    test('fromJson handles defaults for missing fields', () {
      final json = <String, dynamic>{};

      final response = SyncConfigsResponse.fromJson(json);

      expect(response.created, 0);
      expect(response.updated, 0);
      expect(response.deleted, 0);
      expect(response.configs, isEmpty);
    });
  });

  // =========================================================================
  // OrderStatus
  // =========================================================================
  group('OrderStatus', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'code': 'pending',
        'name': 'Pending',
        'description': 'Order is pending',
        'category': 'processing',
        'is_active': true,
        'icon': 'clock',
        'color': '#FFA500',
        'created_at': '2026-01-01',
        'updated_at': '2026-01-02',
      };

      final status = OrderStatus.fromJson(json);

      expect(status.id, 1);
      expect(status.code, 'pending');
      expect(status.name, 'Pending');
      expect(status.description, 'Order is pending');
      expect(status.category, 'processing');
      expect(status.isActive, true);
      expect(status.icon, 'clock');
      expect(status.color, '#FFA500');
      expect(status.createdAt, '2026-01-01');
      expect(status.updatedAt, '2026-01-02');
    });

    test('fromJson handles defaults for missing fields', () {
      final json = <String, dynamic>{};

      final status = OrderStatus.fromJson(json);

      expect(status.id, 0);
      expect(status.code, '');
      expect(status.name, '');
      expect(status.category, '');
      expect(status.isActive, true);
      expect(status.createdAt, '');
      expect(status.updatedAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'code': 'c',
        'name': 'n',
        'category': 'cat',
        'is_active': true,
        'created_at': '',
        'updated_at': '',
      };

      final status = OrderStatus.fromJson(json);

      expect(status.description, isNull);
      expect(status.icon, isNull);
      expect(status.color, isNull);
    });
  });

  // =========================================================================
  // Integration
  // =========================================================================
  group('Integration', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'name': 'My Store',
        'code': 'my_store',
        'type': 'ecommerce',
        'business_id': 5,
        'is_active': true,
        'integration_type_id': 10,
        'integration_type_name': 'Shopify',
        'integration_type_icon': 'shop',
        'created_at': '2026-01-01',
        'updated_at': '2026-01-02',
      };

      final integration = Integration.fromJson(json);

      expect(integration.id, 1);
      expect(integration.name, 'My Store');
      expect(integration.code, 'my_store');
      expect(integration.type, 'ecommerce');
      expect(integration.businessId, 5);
      expect(integration.isActive, true);
      expect(integration.integrationTypeId, 10);
      expect(integration.integrationTypeName, 'Shopify');
      expect(integration.integrationTypeIcon, 'shop');
      expect(integration.createdAt, '2026-01-01');
      expect(integration.updatedAt, '2026-01-02');
    });

    test('fromJson handles defaults for missing fields', () {
      final json = <String, dynamic>{};

      final integration = Integration.fromJson(json);

      expect(integration.id, 0);
      expect(integration.name, '');
      expect(integration.code, '');
      expect(integration.type, '');
      expect(integration.isActive, true);
      expect(integration.createdAt, '');
      expect(integration.updatedAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'name': 'n',
        'code': 'c',
        'type': 't',
        'is_active': true,
        'created_at': '',
        'updated_at': '',
      };

      final integration = Integration.fromJson(json);

      expect(integration.businessId, isNull);
      expect(integration.integrationTypeId, isNull);
      expect(integration.integrationTypeName, isNull);
      expect(integration.integrationTypeIcon, isNull);
    });
  });

  // =========================================================================
  // MessageAuditLog
  // =========================================================================
  group('MessageAuditLog', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 'uuid-123',
        'conversation_id': 'conv-1',
        'message_id': 'msg-1',
        'direction': 'outbound',
        'template_name': 'order_confirmation',
        'content': 'Your order has been confirmed',
        'status': 'delivered',
        'delivered_at': '2026-01-01T12:00:00Z',
        'read_at': '2026-01-01T12:05:00Z',
        'created_at': '2026-01-01T11:59:00Z',
        'phone_number': '+573001234567',
        'order_number': 'ORD-001',
        'business_id': 5,
      };

      final log = MessageAuditLog.fromJson(json);

      expect(log.id, 'uuid-123');
      expect(log.conversationId, 'conv-1');
      expect(log.messageId, 'msg-1');
      expect(log.direction, 'outbound');
      expect(log.templateName, 'order_confirmation');
      expect(log.content, 'Your order has been confirmed');
      expect(log.status, 'delivered');
      expect(log.deliveredAt, '2026-01-01T12:00:00Z');
      expect(log.readAt, '2026-01-01T12:05:00Z');
      expect(log.createdAt, '2026-01-01T11:59:00Z');
      expect(log.phoneNumber, '+573001234567');
      expect(log.orderNumber, 'ORD-001');
      expect(log.businessId, 5);
    });

    test('fromJson handles defaults for missing fields', () {
      final json = <String, dynamic>{};

      final log = MessageAuditLog.fromJson(json);

      expect(log.id, '');
      expect(log.conversationId, '');
      expect(log.messageId, '');
      expect(log.direction, '');
      expect(log.templateName, '');
      expect(log.content, '');
      expect(log.status, '');
      expect(log.createdAt, '');
      expect(log.phoneNumber, '');
      expect(log.orderNumber, '');
      expect(log.businessId, 0);
    });

    test('fromJson converts id to string', () {
      final json = {
        'id': 12345,
        'conversation_id': '',
        'message_id': '',
        'direction': '',
        'template_name': '',
        'content': '',
        'status': '',
        'created_at': '',
        'phone_number': '',
        'order_number': '',
        'business_id': 0,
      };

      final log = MessageAuditLog.fromJson(json);
      expect(log.id, '12345');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': '1',
        'conversation_id': '',
        'message_id': '',
        'direction': '',
        'template_name': '',
        'content': '',
        'status': '',
        'created_at': '',
        'phone_number': '',
        'order_number': '',
        'business_id': 0,
      };

      final log = MessageAuditLog.fromJson(json);

      expect(log.deliveredAt, isNull);
      expect(log.readAt, isNull);
    });
  });

  // =========================================================================
  // MessageAuditStats
  // =========================================================================
  group('MessageAuditStats', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'total_sent': 100,
        'total_delivered': 95,
        'total_read': 80,
        'total_failed': 5,
        'success_rate': 95.0,
      };

      final stats = MessageAuditStats.fromJson(json);

      expect(stats.totalSent, 100);
      expect(stats.totalDelivered, 95);
      expect(stats.totalRead, 80);
      expect(stats.totalFailed, 5);
      expect(stats.successRate, 95.0);
    });

    test('fromJson handles defaults for missing fields', () {
      final json = <String, dynamic>{};

      final stats = MessageAuditStats.fromJson(json);

      expect(stats.totalSent, 0);
      expect(stats.totalDelivered, 0);
      expect(stats.totalRead, 0);
      expect(stats.totalFailed, 0);
      expect(stats.successRate, 0.0);
    });
  });

  // =========================================================================
  // MessageAuditFilter
  // =========================================================================
  group('MessageAuditFilter', () {
    test('toQueryParams includes all set fields', () {
      final filter = MessageAuditFilter(
        businessId: 5,
        status: 'delivered',
        direction: 'outbound',
        templateName: 'order_confirmation',
        dateFrom: '2026-01-01',
        dateTo: '2026-12-31',
        page: 2,
        pageSize: 20,
      );

      final query = filter.toQueryParams();

      expect(query['business_id'], 5);
      expect(query['status'], 'delivered');
      expect(query['direction'], 'outbound');
      expect(query['template_name'], 'order_confirmation');
      expect(query['date_from'], '2026-01-01');
      expect(query['date_to'], '2026-12-31');
      expect(query['page'], 2);
      expect(query['page_size'], 20);
    });

    test('toQueryParams always includes businessId', () {
      final filter = MessageAuditFilter(businessId: 5);
      final query = filter.toQueryParams();

      expect(query['business_id'], 5);
      expect(query.length, 1);
    });

    test('toQueryParams omits null optional fields', () {
      final filter = MessageAuditFilter(businessId: 1);
      final query = filter.toQueryParams();

      expect(query.containsKey('status'), false);
      expect(query.containsKey('direction'), false);
      expect(query.containsKey('template_name'), false);
      expect(query.containsKey('date_from'), false);
      expect(query.containsKey('date_to'), false);
      expect(query.containsKey('page'), false);
      expect(query.containsKey('page_size'), false);
    });
  });

  // =========================================================================
  // PaginatedMessageAuditResponse
  // =========================================================================
  group('PaginatedMessageAuditResponse', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'data': [
          {
            'id': '1',
            'conversation_id': '',
            'message_id': '',
            'direction': 'outbound',
            'template_name': 'test',
            'content': 'Hello',
            'status': 'sent',
            'created_at': '',
            'phone_number': '',
            'order_number': '',
            'business_id': 1,
          },
        ],
        'total': 100,
        'page': 2,
        'page_size': 20,
        'total_pages': 5,
      };

      final response = PaginatedMessageAuditResponse.fromJson(json);

      expect(response.data.length, 1);
      expect(response.data.first.direction, 'outbound');
      expect(response.total, 100);
      expect(response.page, 2);
      expect(response.pageSize, 20);
      expect(response.totalPages, 5);
    });

    test('fromJson handles defaults for missing fields', () {
      final json = <String, dynamic>{};

      final response = PaginatedMessageAuditResponse.fromJson(json);

      expect(response.data, isEmpty);
      expect(response.total, 0);
      expect(response.page, 1);
      expect(response.pageSize, 10);
      expect(response.totalPages, 0);
    });
  });
}
