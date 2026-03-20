import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/orderstatus/domain/entities.dart';

void main() {
  group('OrderStatusInfo', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'code': 'pending',
        'name': 'Pending',
        'description': 'Order is pending',
        'category': 'open',
        'color': '#FFA500',
        'priority': 10,
        'is_active': true,
      };

      final status = OrderStatusInfo.fromJson(json);

      expect(status.id, 1);
      expect(status.code, 'pending');
      expect(status.name, 'Pending');
      expect(status.description, 'Order is pending');
      expect(status.category, 'open');
      expect(status.color, '#FFA500');
      expect(status.priority, 10);
      expect(status.isActive, true);
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 2,
        'code': 'shipped',
        'name': 'Shipped',
      };

      final status = OrderStatusInfo.fromJson(json);

      expect(status.description, isNull);
      expect(status.category, isNull);
      expect(status.color, isNull);
      expect(status.priority, isNull);
      expect(status.isActive, isNull);
    });

    test('fromJson defaults to safe values when fields are missing', () {
      final json = <String, dynamic>{};

      final status = OrderStatusInfo.fromJson(json);

      expect(status.id, 0);
      expect(status.code, '');
      expect(status.name, '');
    });
  });

  group('IntegrationTypeInfo', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'code': 'shopify',
        'name': 'Shopify',
        'image_url': 'https://example.com/shopify.png',
      };

      final info = IntegrationTypeInfo.fromJson(json);

      expect(info.id, 1);
      expect(info.code, 'shopify');
      expect(info.name, 'Shopify');
      expect(info.imageUrl, 'https://example.com/shopify.png');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 2,
        'code': 'amazon',
        'name': 'Amazon',
      };

      final info = IntegrationTypeInfo.fromJson(json);

      expect(info.imageUrl, isNull);
    });

    test('fromJson defaults to safe values when fields are missing', () {
      final json = <String, dynamic>{};

      final info = IntegrationTypeInfo.fromJson(json);

      expect(info.id, 0);
      expect(info.code, '');
      expect(info.name, '');
    });
  });

  group('OrderStatusMapping', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'integration_type_id': 5,
        'integration_type': {
          'id': 5,
          'code': 'shopify',
          'name': 'Shopify',
          'image_url': 'https://example.com/shopify.png',
        },
        'original_status': 'paid',
        'order_status_id': 2,
        'order_status': {
          'id': 2,
          'code': 'confirmed',
          'name': 'Confirmed',
        },
        'is_active': true,
        'description': 'Maps paid to confirmed',
        'created_at': '2026-01-01T00:00:00Z',
        'updated_at': '2026-01-02T00:00:00Z',
      };

      final mapping = OrderStatusMapping.fromJson(json);

      expect(mapping.id, 1);
      expect(mapping.integrationTypeId, 5);
      expect(mapping.integrationType, isNotNull);
      expect(mapping.integrationType!.code, 'shopify');
      expect(mapping.originalStatus, 'paid');
      expect(mapping.orderStatusId, 2);
      expect(mapping.orderStatus, isNotNull);
      expect(mapping.orderStatus!.code, 'confirmed');
      expect(mapping.isActive, true);
      expect(mapping.description, 'Maps paid to confirmed');
      expect(mapping.createdAt, '2026-01-01T00:00:00Z');
      expect(mapping.updatedAt, '2026-01-02T00:00:00Z');
    });

    test('fromJson handles null nested objects', () {
      final json = {
        'id': 1,
        'integration_type_id': 5,
        'original_status': 'paid',
        'order_status_id': 2,
        'is_active': true,
        'description': '',
        'created_at': '',
        'updated_at': '',
      };

      final mapping = OrderStatusMapping.fromJson(json);

      expect(mapping.integrationType, isNull);
      expect(mapping.orderStatus, isNull);
    });

    test('fromJson defaults to safe values when fields are missing', () {
      final json = <String, dynamic>{};

      final mapping = OrderStatusMapping.fromJson(json);

      expect(mapping.id, 0);
      expect(mapping.integrationTypeId, 0);
      expect(mapping.originalStatus, '');
      expect(mapping.orderStatusId, 0);
      expect(mapping.isActive, true);
      expect(mapping.description, '');
      expect(mapping.createdAt, '');
      expect(mapping.updatedAt, '');
    });
  });

  group('ChannelStatusInfo', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'integration_type_id': 5,
        'integration_type': {
          'id': 5,
          'code': 'shopify',
          'name': 'Shopify',
        },
        'code': 'fulfilled',
        'name': 'Fulfilled',
        'description': 'Order fulfilled on channel',
        'is_active': true,
        'display_order': 3,
      };

      final status = ChannelStatusInfo.fromJson(json);

      expect(status.id, 1);
      expect(status.integrationTypeId, 5);
      expect(status.integrationType, isNotNull);
      expect(status.integrationType!.code, 'shopify');
      expect(status.code, 'fulfilled');
      expect(status.name, 'Fulfilled');
      expect(status.description, 'Order fulfilled on channel');
      expect(status.isActive, true);
      expect(status.displayOrder, 3);
    });

    test('fromJson handles null nested integration type', () {
      final json = {
        'id': 1,
        'integration_type_id': 5,
        'code': 'shipped',
        'name': 'Shipped',
        'is_active': true,
        'display_order': 1,
      };

      final status = ChannelStatusInfo.fromJson(json);

      expect(status.integrationType, isNull);
      expect(status.description, isNull);
    });

    test('fromJson defaults to safe values when fields are missing', () {
      final json = <String, dynamic>{};

      final status = ChannelStatusInfo.fromJson(json);

      expect(status.id, 0);
      expect(status.integrationTypeId, 0);
      expect(status.code, '');
      expect(status.name, '');
      expect(status.isActive, true);
      expect(status.displayOrder, 0);
    });
  });

  group('GetOrderStatusMappingsParams', () {
    test('toQueryParams includes all non-null fields', () {
      final params = GetOrderStatusMappingsParams(
        page: 2,
        pageSize: 25,
        integrationTypeId: 5,
        isActive: true,
      );

      final query = params.toQueryParams();

      expect(query['page'], 2);
      expect(query['page_size'], 25);
      expect(query['integration_type_id'], 5);
      expect(query['is_active'], true);
    });

    test('toQueryParams returns empty map when all fields are null', () {
      final params = GetOrderStatusMappingsParams();

      final query = params.toQueryParams();

      expect(query, isEmpty);
    });

    test('toQueryParams includes only provided fields', () {
      final params = GetOrderStatusMappingsParams(page: 1);

      final query = params.toQueryParams();

      expect(query.length, 1);
      expect(query['page'], 1);
    });
  });

  group('CreateOrderStatusDTO', () {
    test('toJson includes required fields', () {
      final dto = CreateOrderStatusDTO(code: 'new_status', name: 'New Status');

      final json = dto.toJson();

      expect(json['code'], 'new_status');
      expect(json['name'], 'New Status');
    });

    test('toJson includes all optional fields when provided', () {
      final dto = CreateOrderStatusDTO(
        code: 'new_status',
        name: 'New Status',
        description: 'A new status',
        category: 'active',
        color: '#FF0000',
        priority: 5,
        isActive: true,
      );

      final json = dto.toJson();

      expect(json['code'], 'new_status');
      expect(json['name'], 'New Status');
      expect(json['description'], 'A new status');
      expect(json['category'], 'active');
      expect(json['color'], '#FF0000');
      expect(json['priority'], 5);
      expect(json['is_active'], true);
    });

    test('toJson excludes null optional fields', () {
      final dto = CreateOrderStatusDTO(code: 'test', name: 'Test');

      final json = dto.toJson();

      expect(json.containsKey('description'), false);
      expect(json.containsKey('category'), false);
      expect(json.containsKey('color'), false);
      expect(json.containsKey('priority'), false);
      expect(json.containsKey('is_active'), false);
    });
  });

  group('CreateOrderStatusMappingDTO', () {
    test('toJson includes required fields', () {
      final dto = CreateOrderStatusMappingDTO(
        integrationTypeId: 5,
        originalStatus: 'paid',
        orderStatusId: 2,
      );

      final json = dto.toJson();

      expect(json['integration_type_id'], 5);
      expect(json['original_status'], 'paid');
      expect(json['order_status_id'], 2);
    });

    test('toJson includes description when provided', () {
      final dto = CreateOrderStatusMappingDTO(
        integrationTypeId: 5,
        originalStatus: 'paid',
        orderStatusId: 2,
        description: 'Maps paid to confirmed',
      );

      final json = dto.toJson();

      expect(json['description'], 'Maps paid to confirmed');
    });

    test('toJson excludes null description', () {
      final dto = CreateOrderStatusMappingDTO(
        integrationTypeId: 5,
        originalStatus: 'paid',
        orderStatusId: 2,
      );

      final json = dto.toJson();

      expect(json.containsKey('description'), false);
    });
  });

  group('UpdateOrderStatusMappingDTO', () {
    test('toJson includes required fields', () {
      final dto = UpdateOrderStatusMappingDTO(
        originalStatus: 'shipped',
        orderStatusId: 3,
      );

      final json = dto.toJson();

      expect(json['original_status'], 'shipped');
      expect(json['order_status_id'], 3);
    });

    test('toJson includes description when provided', () {
      final dto = UpdateOrderStatusMappingDTO(
        originalStatus: 'shipped',
        orderStatusId: 3,
        description: 'Updated mapping',
      );

      final json = dto.toJson();

      expect(json['description'], 'Updated mapping');
    });

    test('toJson excludes null description', () {
      final dto = UpdateOrderStatusMappingDTO(
        originalStatus: 'shipped',
        orderStatusId: 3,
      );

      final json = dto.toJson();

      expect(json.containsKey('description'), false);
    });
  });

  group('CreateChannelStatusDTO', () {
    test('toJson includes all required fields', () {
      final dto = CreateChannelStatusDTO(
        integrationTypeId: 5,
        code: 'fulfilled',
        name: 'Fulfilled',
        isActive: true,
        displayOrder: 3,
      );

      final json = dto.toJson();

      expect(json['integration_type_id'], 5);
      expect(json['code'], 'fulfilled');
      expect(json['name'], 'Fulfilled');
      expect(json['is_active'], true);
      expect(json['display_order'], 3);
    });

    test('toJson includes description when provided', () {
      final dto = CreateChannelStatusDTO(
        integrationTypeId: 5,
        code: 'fulfilled',
        name: 'Fulfilled',
        isActive: true,
        displayOrder: 3,
        description: 'Order is fulfilled',
      );

      final json = dto.toJson();

      expect(json['description'], 'Order is fulfilled');
    });

    test('toJson excludes null description', () {
      final dto = CreateChannelStatusDTO(
        integrationTypeId: 5,
        code: 'fulfilled',
        name: 'Fulfilled',
        isActive: true,
        displayOrder: 3,
      );

      final json = dto.toJson();

      expect(json.containsKey('description'), false);
    });
  });

  group('UpdateChannelStatusDTO', () {
    test('toJson includes all required fields', () {
      final dto = UpdateChannelStatusDTO(
        code: 'shipped',
        name: 'Shipped',
        isActive: true,
        displayOrder: 2,
      );

      final json = dto.toJson();

      expect(json['code'], 'shipped');
      expect(json['name'], 'Shipped');
      expect(json['is_active'], true);
      expect(json['display_order'], 2);
    });

    test('toJson includes description when provided', () {
      final dto = UpdateChannelStatusDTO(
        code: 'shipped',
        name: 'Shipped',
        isActive: true,
        displayOrder: 2,
        description: 'Updated description',
      );

      final json = dto.toJson();

      expect(json['description'], 'Updated description');
    });

    test('toJson excludes null description', () {
      final dto = UpdateChannelStatusDTO(
        code: 'shipped',
        name: 'Shipped',
        isActive: true,
        displayOrder: 2,
      );

      final json = dto.toJson();

      expect(json.containsKey('description'), false);
    });
  });
}
