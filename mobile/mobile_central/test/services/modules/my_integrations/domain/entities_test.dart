import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/my_integrations/domain/entities.dart';

void main() {
  // =========================================================================
  // MyIntegration
  // =========================================================================
  group('MyIntegration', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'created_at': '2026-01-01T00:00:00Z',
        'updated_at': '2026-01-02T00:00:00Z',
        'deleted_at': '2026-01-03T00:00:00Z',
        'business_id': 5,
        'integration_type_id': 10,
        'integration_type_name': 'Shopify',
        'integration_type_code': 'shopify',
        'category_code': 'ecommerce',
        'name': 'My Shopify Store',
        'is_active': true,
        'credentials': {'api_key': 'key123'},
        'config': {'sync_interval': 60},
      };

      final integration = MyIntegration.fromJson(json);

      expect(integration.id, 1);
      expect(integration.createdAt, '2026-01-01T00:00:00Z');
      expect(integration.updatedAt, '2026-01-02T00:00:00Z');
      expect(integration.deletedAt, '2026-01-03T00:00:00Z');
      expect(integration.businessId, 5);
      expect(integration.integrationTypeId, 10);
      expect(integration.integrationTypeName, 'Shopify');
      expect(integration.integrationTypeCode, 'shopify');
      expect(integration.categoryCode, 'ecommerce');
      expect(integration.name, 'My Shopify Store');
      expect(integration.isActive, true);
      expect(integration.credentials, {'api_key': 'key123'});
      expect(integration.config, {'sync_interval': 60});
    });

    test('fromJson handles defaults for missing required fields', () {
      final json = <String, dynamic>{};

      final integration = MyIntegration.fromJson(json);

      expect(integration.id, 0);
      expect(integration.createdAt, '');
      expect(integration.updatedAt, '');
      expect(integration.businessId, 0);
      expect(integration.integrationTypeId, 0);
      expect(integration.name, '');
      expect(integration.isActive, false);
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'created_at': '2026-01-01',
        'updated_at': '2026-01-01',
        'business_id': 1,
        'integration_type_id': 1,
        'name': 'Test',
        'is_active': true,
      };

      final integration = MyIntegration.fromJson(json);

      expect(integration.deletedAt, isNull);
      expect(integration.integrationTypeName, isNull);
      expect(integration.integrationTypeCode, isNull);
      expect(integration.categoryCode, isNull);
      expect(integration.credentials, isNull);
      expect(integration.config, isNull);
    });
  });

  // =========================================================================
  // IntegrationCategory
  // =========================================================================
  group('IntegrationCategory', () {
    test('constructor sets fields correctly', () {
      final category = IntegrationCategory(code: 'ecommerce', icon: 'cart');

      expect(category.code, 'ecommerce');
      expect(category.icon, 'cart');
    });
  });

  // =========================================================================
  // Constants
  // =========================================================================
  group('Constants', () {
    test('channelCodes has expected values', () {
      expect(channelCodes, ['platform', 'ecommerce']);
    });

    test('serviceCodes has expected values', () {
      expect(serviceCodes, ['messaging', 'invoicing', 'shipping', 'payment']);
    });

    test('categoryIcons has entries for all categories', () {
      expect(categoryIcons.containsKey('platform'), true);
      expect(categoryIcons.containsKey('ecommerce'), true);
      expect(categoryIcons.containsKey('invoicing'), true);
      expect(categoryIcons.containsKey('messaging'), true);
      expect(categoryIcons.containsKey('payment'), true);
      expect(categoryIcons.containsKey('shipping'), true);
    });

    test('categoryIcons has correct icon values', () {
      expect(categoryIcons['platform'], 'puzzle');
      expect(categoryIcons['ecommerce'], 'cart');
      expect(categoryIcons['invoicing'], 'receipt');
      expect(categoryIcons['messaging'], 'chat');
      expect(categoryIcons['payment'], 'credit_card');
      expect(categoryIcons['shipping'], 'local_shipping');
    });
  });

  // =========================================================================
  // GetMyIntegrationsParams
  // =========================================================================
  group('GetMyIntegrationsParams', () {
    test('toQueryParams includes all set fields', () {
      final params = GetMyIntegrationsParams(
        page: 2,
        pageSize: 20,
        businessId: 5,
        categoryCode: 'ecommerce',
        isActive: true,
      );

      final query = params.toQueryParams();

      expect(query['page'], 2);
      expect(query['page_size'], 20);
      expect(query['business_id'], 5);
      expect(query['category_code'], 'ecommerce');
      expect(query['is_active'], true);
    });

    test('toQueryParams omits null fields', () {
      final params = GetMyIntegrationsParams();
      final query = params.toQueryParams();
      expect(query, isEmpty);
    });

    test('toQueryParams includes only provided fields', () {
      final params = GetMyIntegrationsParams(page: 1, isActive: false);
      final query = params.toQueryParams();

      expect(query.length, 2);
      expect(query['page'], 1);
      expect(query['is_active'], false);
      expect(query.containsKey('page_size'), false);
      expect(query.containsKey('business_id'), false);
      expect(query.containsKey('category_code'), false);
    });
  });
}
