import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/core/domain/entities.dart';

void main() {
  group('IntegrationConfig', () {
    test('fromJson parses data correctly', () {
      final json = {'key1': 'value1', 'key2': 42};

      final config = IntegrationConfig.fromJson(json);

      expect(config.data['key1'], 'value1');
      expect(config.data['key2'], 42);
    });

    test('fromJson handles null by returning empty map', () {
      final config = IntegrationConfig.fromJson(null);

      expect(config.data, isEmpty);
    });

    test('toJson returns the data map', () {
      final config = IntegrationConfig(data: {'enabled': true, 'limit': 100});

      final json = config.toJson();

      expect(json['enabled'], true);
      expect(json['limit'], 100);
    });

    test('default constructor uses empty map', () {
      final config = IntegrationConfig();

      expect(config.data, isEmpty);
    });

    test('fromJson with empty map returns empty data', () {
      final config = IntegrationConfig.fromJson({});

      expect(config.data, isEmpty);
    });
  });

  group('IntegrationCredentials', () {
    test('fromJson parses data correctly', () {
      final json = {'api_key': 'abc123', 'secret': 'xyz'};

      final creds = IntegrationCredentials.fromJson(json);

      expect(creds.data['api_key'], 'abc123');
      expect(creds.data['secret'], 'xyz');
    });

    test('fromJson handles null by returning empty map', () {
      final creds = IntegrationCredentials.fromJson(null);

      expect(creds.data, isEmpty);
    });

    test('toJson returns the data map', () {
      final creds =
          IntegrationCredentials(data: {'token': 'bearer_xyz'});

      final json = creds.toJson();

      expect(json['token'], 'bearer_xyz');
    });

    test('default constructor uses empty map', () {
      final creds = IntegrationCredentials();

      expect(creds.data, isEmpty);
    });

    test('fromJson with empty map returns empty data', () {
      final creds = IntegrationCredentials.fromJson({});

      expect(creds.data, isEmpty);
    });
  });

  group('IntegrationTypeInfo', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 5,
        'name': 'Shopify',
        'code': 'shopify',
        'image_url': 'https://example.com/shopify.png',
      };

      final info = IntegrationTypeInfo.fromJson(json);

      expect(info.id, 5);
      expect(info.name, 'Shopify');
      expect(info.code, 'shopify');
      expect(info.imageUrl, 'https://example.com/shopify.png');
    });

    test('fromJson uses defaults for missing required fields', () {
      final json = <String, dynamic>{};

      final info = IntegrationTypeInfo.fromJson(json);

      expect(info.id, 0);
      expect(info.name, '');
      expect(info.code, '');
    });

    test('fromJson handles null optional imageUrl', () {
      final json = {
        'id': 1,
        'name': 'Test',
        'code': 'test',
      };

      final info = IntegrationTypeInfo.fromJson(json);

      expect(info.imageUrl, isNull);
    });

    test('toJson includes required fields', () {
      final info = IntegrationTypeInfo(id: 3, name: 'Amazon', code: 'amazon');

      final json = info.toJson();

      expect(json['id'], 3);
      expect(json['name'], 'Amazon');
      expect(json['code'], 'amazon');
      expect(json.containsKey('image_url'), false);
    });

    test('toJson includes imageUrl when present', () {
      final info = IntegrationTypeInfo(
        id: 3,
        name: 'Amazon',
        code: 'amazon',
        imageUrl: 'https://example.com/amazon.png',
      );

      final json = info.toJson();

      expect(json['image_url'], 'https://example.com/amazon.png');
    });
  });

  group('Integration', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 10,
        'name': 'My Shopify Store',
        'code': 'my_shopify',
        'integration_type_id': 5,
        'type': 'ecommerce',
        'category': 'sales_channel',
        'category_name': 'Sales Channel',
        'category_color': '#FF5733',
        'business_id': 1,
        'business_name': 'Test Business',
        'store_id': 'store_abc',
        'is_active': true,
        'is_default': false,
        'is_testing': true,
        'config': {'webhook_url': 'https://example.com/webhook'},
        'credentials': {'api_key': 'key123'},
        'description': 'Main store integration',
        'created_by_id': 42,
        'updated_by_id': 43,
        'created_at': '2026-01-01T00:00:00Z',
        'updated_at': '2026-01-02T00:00:00Z',
        'integration_type': {
          'id': 5,
          'name': 'Shopify',
          'code': 'shopify',
          'image_url': 'https://example.com/shopify.png',
        },
      };

      final integration = Integration.fromJson(json);

      expect(integration.id, 10);
      expect(integration.name, 'My Shopify Store');
      expect(integration.code, 'my_shopify');
      expect(integration.integrationTypeId, 5);
      expect(integration.type, 'ecommerce');
      expect(integration.category, 'sales_channel');
      expect(integration.categoryName, 'Sales Channel');
      expect(integration.categoryColor, '#FF5733');
      expect(integration.businessId, 1);
      expect(integration.businessName, 'Test Business');
      expect(integration.storeId, 'store_abc');
      expect(integration.isActive, true);
      expect(integration.isDefault, false);
      expect(integration.isTesting, true);
      expect(integration.config.data['webhook_url'],
          'https://example.com/webhook');
      expect(integration.credentials, isNotNull);
      expect(integration.credentials!.data['api_key'], 'key123');
      expect(integration.description, 'Main store integration');
      expect(integration.createdById, 42);
      expect(integration.updatedById, 43);
      expect(integration.createdAt, '2026-01-01T00:00:00Z');
      expect(integration.updatedAt, '2026-01-02T00:00:00Z');
      expect(integration.integrationType, isNotNull);
      expect(integration.integrationType!.id, 5);
      expect(integration.integrationType!.name, 'Shopify');
      expect(integration.integrationType!.code, 'shopify');
    });

    test('fromJson uses defaults for missing required fields', () {
      final json = <String, dynamic>{};

      final integration = Integration.fromJson(json);

      expect(integration.id, 0);
      expect(integration.name, '');
      expect(integration.code, '');
      expect(integration.integrationTypeId, 0);
      expect(integration.type, '');
      expect(integration.category, '');
      expect(integration.isActive, false);
      expect(integration.isDefault, false);
      expect(integration.isTesting, false);
      expect(integration.config.data, isEmpty);
      expect(integration.createdById, 0);
      expect(integration.createdAt, '');
      expect(integration.updatedAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'name': 'Test',
        'code': 'test',
        'integration_type_id': 1,
        'type': 'ecommerce',
        'category': 'sales',
        'is_active': true,
        'is_default': false,
        'is_testing': false,
        'created_by_id': 1,
        'created_at': '2026-01-01',
        'updated_at': '2026-01-01',
      };

      final integration = Integration.fromJson(json);

      expect(integration.categoryName, isNull);
      expect(integration.categoryColor, isNull);
      expect(integration.businessId, isNull);
      expect(integration.businessName, isNull);
      expect(integration.storeId, isNull);
      expect(integration.credentials, isNull);
      expect(integration.description, isNull);
      expect(integration.updatedById, isNull);
      expect(integration.integrationType, isNull);
    });

    test('fromJson handles null credentials', () {
      final json = {
        'id': 1,
        'name': 'Test',
        'code': 'test',
        'integration_type_id': 1,
        'type': 'ecommerce',
        'category': 'sales',
        'is_active': true,
        'is_default': false,
        'is_testing': false,
        'config': {'key': 'val'},
        'credentials': null,
        'created_by_id': 1,
        'created_at': '2026-01-01',
        'updated_at': '2026-01-01',
      };

      final integration = Integration.fromJson(json);

      expect(integration.credentials, isNull);
    });

    test('toJson includes required fields', () {
      final integration = Integration(
        id: 10,
        name: 'My Store',
        code: 'my_store',
        integrationTypeId: 5,
        type: 'ecommerce',
        category: 'sales_channel',
        isActive: true,
        isDefault: false,
        isTesting: false,
        config: IntegrationConfig(data: {'key': 'value'}),
        createdById: 42,
        createdAt: '2026-01-01',
        updatedAt: '2026-01-02',
      );

      final json = integration.toJson();

      expect(json['id'], 10);
      expect(json['name'], 'My Store');
      expect(json['code'], 'my_store');
      expect(json['integration_type_id'], 5);
      expect(json['type'], 'ecommerce');
      expect(json['category'], 'sales_channel');
      expect(json['is_active'], true);
      expect(json['is_default'], false);
      expect(json['is_testing'], false);
      expect(json['config'], {'key': 'value'});
      expect(json['created_by_id'], 42);
      expect(json['created_at'], '2026-01-01');
      expect(json['updated_at'], '2026-01-02');
    });

    test('toJson excludes null optional fields', () {
      final integration = Integration(
        id: 1,
        name: 'Test',
        code: 'test',
        integrationTypeId: 1,
        type: 'ecommerce',
        category: 'sales',
        isActive: true,
        isDefault: false,
        isTesting: false,
        config: IntegrationConfig(),
        createdById: 1,
        createdAt: '2026-01-01',
        updatedAt: '2026-01-01',
      );

      final json = integration.toJson();

      expect(json.containsKey('category_name'), false);
      expect(json.containsKey('category_color'), false);
      expect(json.containsKey('business_id'), false);
      expect(json.containsKey('business_name'), false);
      expect(json.containsKey('store_id'), false);
      expect(json.containsKey('credentials'), false);
      expect(json.containsKey('description'), false);
      expect(json.containsKey('updated_by_id'), false);
      expect(json.containsKey('integration_type'), false);
    });

    test('toJson includes optional fields when present', () {
      final integration = Integration(
        id: 1,
        name: 'Test',
        code: 'test',
        integrationTypeId: 1,
        type: 'ecommerce',
        category: 'sales',
        categoryName: 'Sales',
        categoryColor: '#00FF00',
        businessId: 5,
        businessName: 'Biz',
        storeId: 'store_1',
        isActive: true,
        isDefault: false,
        isTesting: false,
        config: IntegrationConfig(),
        credentials: IntegrationCredentials(data: {'token': 'abc'}),
        description: 'A description',
        createdById: 1,
        updatedById: 2,
        createdAt: '2026-01-01',
        updatedAt: '2026-01-01',
        integrationType: IntegrationTypeInfo(
          id: 1,
          name: 'Shopify',
          code: 'shopify',
        ),
      );

      final json = integration.toJson();

      expect(json['category_name'], 'Sales');
      expect(json['category_color'], '#00FF00');
      expect(json['business_id'], 5);
      expect(json['business_name'], 'Biz');
      expect(json['store_id'], 'store_1');
      expect(json['credentials'], {'token': 'abc'});
      expect(json['description'], 'A description');
      expect(json['updated_by_id'], 2);
      expect(json['integration_type'], isA<Map<String, dynamic>>());
      expect(json['integration_type']['name'], 'Shopify');
    });
  });

  group('CreateIntegrationDTO', () {
    test('toJson includes all required fields', () {
      final dto = CreateIntegrationDTO(
        name: 'New Integration',
        code: 'new_int',
        integrationTypeId: 3,
        category: 'payment',
      );

      final json = dto.toJson();

      expect(json['name'], 'New Integration');
      expect(json['code'], 'new_int');
      expect(json['integration_type_id'], 3);
      expect(json['category'], 'payment');
    });

    test('toJson includes all fields when provided', () {
      final dto = CreateIntegrationDTO(
        name: 'Full Integration',
        code: 'full_int',
        integrationTypeId: 3,
        type: 'ecommerce',
        category: 'payment',
        businessId: 10,
        storeId: 'store_xyz',
        isActive: true,
        isDefault: false,
        isTesting: true,
        config: IntegrationConfig(data: {'timeout': 30}),
        credentials: IntegrationCredentials(data: {'key': 'val'}),
        description: 'Full description',
      );

      final json = dto.toJson();

      expect(json['name'], 'Full Integration');
      expect(json['code'], 'full_int');
      expect(json['integration_type_id'], 3);
      expect(json['type'], 'ecommerce');
      expect(json['category'], 'payment');
      expect(json['business_id'], 10);
      expect(json['store_id'], 'store_xyz');
      expect(json['is_active'], true);
      expect(json['is_default'], false);
      expect(json['is_testing'], true);
      expect(json['config'], {'timeout': 30});
      expect(json['credentials'], {'key': 'val'});
      expect(json['description'], 'Full description');
    });

    test('toJson excludes null optional fields', () {
      final dto = CreateIntegrationDTO(
        name: 'Minimal',
        code: 'min',
        integrationTypeId: 1,
        category: 'sales',
      );

      final json = dto.toJson();

      expect(json.length, 4);
      expect(json.containsKey('type'), false);
      expect(json.containsKey('business_id'), false);
      expect(json.containsKey('store_id'), false);
      expect(json.containsKey('is_active'), false);
      expect(json.containsKey('is_default'), false);
      expect(json.containsKey('is_testing'), false);
      expect(json.containsKey('config'), false);
      expect(json.containsKey('credentials'), false);
      expect(json.containsKey('description'), false);
    });
  });

  group('UpdateIntegrationDTO', () {
    test('toJson includes all fields when provided', () {
      final dto = UpdateIntegrationDTO(
        name: 'Updated Name',
        code: 'upd_code',
        storeId: 'new_store',
        isActive: false,
        isDefault: true,
        isTesting: false,
        config: IntegrationConfig(data: {'retries': 3}),
        credentials: IntegrationCredentials(data: {'secret': 'new'}),
        description: 'Updated desc',
      );

      final json = dto.toJson();

      expect(json['name'], 'Updated Name');
      expect(json['code'], 'upd_code');
      expect(json['store_id'], 'new_store');
      expect(json['is_active'], false);
      expect(json['is_default'], true);
      expect(json['is_testing'], false);
      expect(json['config'], {'retries': 3});
      expect(json['credentials'], {'secret': 'new'});
      expect(json['description'], 'Updated desc');
    });

    test('toJson returns empty map when all fields are null', () {
      final dto = UpdateIntegrationDTO();

      final json = dto.toJson();

      expect(json, isEmpty);
    });

    test('toJson includes only provided fields', () {
      final dto = UpdateIntegrationDTO(name: 'OnlyName', isActive: true);

      final json = dto.toJson();

      expect(json.length, 2);
      expect(json['name'], 'OnlyName');
      expect(json['is_active'], true);
      expect(json.containsKey('code'), false);
      expect(json.containsKey('store_id'), false);
    });

    test('toJson excludes null config and credentials', () {
      final dto = UpdateIntegrationDTO(name: 'Test');

      final json = dto.toJson();

      expect(json.containsKey('config'), false);
      expect(json.containsKey('credentials'), false);
    });
  });

  group('GetIntegrationsParams', () {
    test('toQueryParams includes all non-null fields', () {
      final params = GetIntegrationsParams(
        page: 2,
        pageSize: 25,
        type: 'ecommerce',
        category: 'sales_channel',
        categoryId: 7,
        businessId: 3,
        isActive: true,
        search: 'shopify',
      );

      final queryParams = params.toQueryParams();

      expect(queryParams['page'], 2);
      expect(queryParams['page_size'], 25);
      expect(queryParams['type'], 'ecommerce');
      expect(queryParams['category'], 'sales_channel');
      expect(queryParams['category_id'], 7);
      expect(queryParams['business_id'], 3);
      expect(queryParams['is_active'], true);
      expect(queryParams['search'], 'shopify');
    });

    test('toQueryParams excludes null fields', () {
      final params = GetIntegrationsParams(page: 1, pageSize: 10);

      final queryParams = params.toQueryParams();

      expect(queryParams.length, 2);
      expect(queryParams.containsKey('page'), true);
      expect(queryParams.containsKey('page_size'), true);
      expect(queryParams.containsKey('type'), false);
      expect(queryParams.containsKey('category'), false);
      expect(queryParams.containsKey('category_id'), false);
      expect(queryParams.containsKey('business_id'), false);
      expect(queryParams.containsKey('is_active'), false);
      expect(queryParams.containsKey('search'), false);
    });

    test('toQueryParams returns empty map when all fields are null', () {
      final params = GetIntegrationsParams();

      final queryParams = params.toQueryParams();

      expect(queryParams, isEmpty);
    });

    test('toQueryParams includes only partial fields', () {
      final params = GetIntegrationsParams(
        category: 'payment',
        isActive: false,
      );

      final queryParams = params.toQueryParams();

      expect(queryParams.length, 2);
      expect(queryParams['category'], 'payment');
      expect(queryParams['is_active'], false);
    });
  });

  group('ActionResponse', () {
    test('fromJson parses success response', () {
      final json = {
        'success': true,
        'message': 'Integration created successfully',
      };

      final response = ActionResponse.fromJson(json);

      expect(response.success, true);
      expect(response.message, 'Integration created successfully');
      expect(response.error, isNull);
    });

    test('fromJson parses failure response with error', () {
      final json = {
        'success': false,
        'message': 'Failed to create integration',
        'error': 'Duplicate code',
      };

      final response = ActionResponse.fromJson(json);

      expect(response.success, false);
      expect(response.message, 'Failed to create integration');
      expect(response.error, 'Duplicate code');
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final response = ActionResponse.fromJson(json);

      expect(response.success, false);
      expect(response.message, '');
      expect(response.error, isNull);
    });

    test('fromJson handles null message', () {
      final json = {
        'success': true,
        'message': null,
      };

      final response = ActionResponse.fromJson(json);

      expect(response.success, true);
      expect(response.message, '');
    });
  });

  group('IntegrationType', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'name': 'Shopify',
        'code': 'shopify',
        'description': 'Shopify e-commerce platform',
        'icon': 'shopify_icon',
        'image_url': 'https://example.com/shopify.png',
        'category': {
          'id': 2,
          'code': 'ecommerce',
          'name': 'E-Commerce',
          'display_order': 1,
          'is_active': true,
          'is_visible': true,
          'created_at': '2026-01-01',
          'updated_at': '2026-01-01',
        },
        'category_id': 2,
        'integration_category': {
          'id': 3,
          'code': 'sales',
          'name': 'Sales',
          'display_order': 2,
          'is_active': true,
          'is_visible': true,
          'created_at': '2026-01-01',
          'updated_at': '2026-01-01',
        },
        'is_active': true,
        'in_development': false,
        'config_schema': {'type': 'object', 'properties': {}},
        'credentials_schema': {'type': 'object', 'required': ['api_key']},
        'setup_instructions': 'Go to Shopify admin...',
        'base_url': 'https://api.shopify.com',
        'base_url_test': 'https://test.shopify.com',
        'has_platform_credentials': true,
        'platform_credential_keys': ['api_key', 'api_secret'],
        'created_at': '2026-01-01T00:00:00Z',
        'updated_at': '2026-01-02T00:00:00Z',
      };

      final intType = IntegrationType.fromJson(json);

      expect(intType.id, 1);
      expect(intType.name, 'Shopify');
      expect(intType.code, 'shopify');
      expect(intType.description, 'Shopify e-commerce platform');
      expect(intType.icon, 'shopify_icon');
      expect(intType.imageUrl, 'https://example.com/shopify.png');
      expect(intType.category, isNotNull);
      expect(intType.category!.id, 2);
      expect(intType.category!.code, 'ecommerce');
      expect(intType.categoryId, 2);
      expect(intType.integrationCategory, isNotNull);
      expect(intType.integrationCategory!.id, 3);
      expect(intType.integrationCategory!.code, 'sales');
      expect(intType.isActive, true);
      expect(intType.inDevelopment, false);
      expect(intType.configSchema, isA<Map>());
      expect(intType.credentialsSchema, isA<Map>());
      expect(intType.setupInstructions, 'Go to Shopify admin...');
      expect(intType.baseUrl, 'https://api.shopify.com');
      expect(intType.baseUrlTest, 'https://test.shopify.com');
      expect(intType.hasPlatformCredentials, true);
      expect(intType.platformCredentialKeys, ['api_key', 'api_secret']);
      expect(intType.createdAt, '2026-01-01T00:00:00Z');
      expect(intType.updatedAt, '2026-01-02T00:00:00Z');
    });

    test('fromJson uses defaults for missing required fields', () {
      final json = <String, dynamic>{};

      final intType = IntegrationType.fromJson(json);

      expect(intType.id, 0);
      expect(intType.name, '');
      expect(intType.code, '');
      expect(intType.isActive, false);
      expect(intType.createdAt, '');
      expect(intType.updatedAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'name': 'Test',
        'code': 'test',
        'is_active': true,
        'created_at': '2026-01-01',
        'updated_at': '2026-01-01',
      };

      final intType = IntegrationType.fromJson(json);

      expect(intType.description, isNull);
      expect(intType.icon, isNull);
      expect(intType.imageUrl, isNull);
      expect(intType.category, isNull);
      expect(intType.categoryId, isNull);
      expect(intType.integrationCategory, isNull);
      expect(intType.inDevelopment, isNull);
      expect(intType.configSchema, isNull);
      expect(intType.credentialsSchema, isNull);
      expect(intType.setupInstructions, isNull);
      expect(intType.baseUrl, isNull);
      expect(intType.baseUrlTest, isNull);
      expect(intType.hasPlatformCredentials, isNull);
      expect(intType.platformCredentialKeys, isNull);
    });

    test('fromJson parses platformCredentialKeys from dynamic list', () {
      final json = {
        'id': 1,
        'name': 'Test',
        'code': 'test',
        'is_active': true,
        'platform_credential_keys': ['key1', 'key2', 'key3'],
        'created_at': '2026-01-01',
        'updated_at': '2026-01-01',
      };

      final intType = IntegrationType.fromJson(json);

      expect(intType.platformCredentialKeys, isNotNull);
      expect(intType.platformCredentialKeys!.length, 3);
      expect(intType.platformCredentialKeys![0], 'key1');
      expect(intType.platformCredentialKeys![2], 'key3');
    });

    test('toJson includes required fields', () {
      final intType = IntegrationType(
        id: 1,
        name: 'Shopify',
        code: 'shopify',
        isActive: true,
        createdAt: '2026-01-01',
        updatedAt: '2026-01-02',
      );

      final json = intType.toJson();

      expect(json['id'], 1);
      expect(json['name'], 'Shopify');
      expect(json['code'], 'shopify');
      expect(json['is_active'], true);
      expect(json['created_at'], '2026-01-01');
      expect(json['updated_at'], '2026-01-02');
    });

    test('toJson excludes null optional fields', () {
      final intType = IntegrationType(
        id: 1,
        name: 'Test',
        code: 'test',
        isActive: true,
        createdAt: '2026-01-01',
        updatedAt: '2026-01-01',
      );

      final json = intType.toJson();

      expect(json.containsKey('description'), false);
      expect(json.containsKey('icon'), false);
      expect(json.containsKey('image_url'), false);
      expect(json.containsKey('category'), false);
      expect(json.containsKey('category_id'), false);
      expect(json.containsKey('in_development'), false);
      expect(json.containsKey('config_schema'), false);
      expect(json.containsKey('credentials_schema'), false);
      expect(json.containsKey('setup_instructions'), false);
      expect(json.containsKey('base_url'), false);
      expect(json.containsKey('base_url_test'), false);
      expect(json.containsKey('has_platform_credentials'), false);
      expect(json.containsKey('platform_credential_keys'), false);
    });

    test('toJson includes all optional fields when present', () {
      final intType = IntegrationType(
        id: 1,
        name: 'Shopify',
        code: 'shopify',
        description: 'A description',
        icon: 'icon_shopify',
        imageUrl: 'https://example.com/img.png',
        category: IntegrationCategory(
          id: 2,
          code: 'ecom',
          name: 'E-Commerce',
          displayOrder: 1,
          isActive: true,
          isVisible: true,
          createdAt: '2026-01-01',
          updatedAt: '2026-01-01',
        ),
        categoryId: 2,
        isActive: true,
        inDevelopment: false,
        configSchema: {'type': 'object'},
        credentialsSchema: {'required': ['key']},
        setupInstructions: 'Instructions here',
        baseUrl: 'https://api.example.com',
        baseUrlTest: 'https://test.example.com',
        hasPlatformCredentials: true,
        platformCredentialKeys: ['api_key'],
        createdAt: '2026-01-01',
        updatedAt: '2026-01-02',
      );

      final json = intType.toJson();

      expect(json['description'], 'A description');
      expect(json['icon'], 'icon_shopify');
      expect(json['image_url'], 'https://example.com/img.png');
      expect(json['category'], isA<Map<String, dynamic>>());
      expect(json['category']['code'], 'ecom');
      expect(json['category_id'], 2);
      expect(json['in_development'], false);
      expect(json['config_schema'], {'type': 'object'});
      expect(json['credentials_schema'], {'required': ['key']});
      expect(json['setup_instructions'], 'Instructions here');
      expect(json['base_url'], 'https://api.example.com');
      expect(json['base_url_test'], 'https://test.example.com');
      expect(json['has_platform_credentials'], true);
      expect(json['platform_credential_keys'], ['api_key']);
    });
  });

  group('CreateIntegrationTypeDTO', () {
    test('toJson includes required fields', () {
      final dto = CreateIntegrationTypeDTO(
        name: 'New Type',
        categoryId: 5,
      );

      final json = dto.toJson();

      expect(json['name'], 'New Type');
      expect(json['category_id'], 5);
    });

    test('toJson includes all fields when provided', () {
      final dto = CreateIntegrationTypeDTO(
        name: 'Full Type',
        code: 'full_type',
        description: 'A full type',
        icon: 'type_icon',
        categoryId: 5,
        isActive: true,
        configSchema: {'type': 'object'},
        credentialsSchema: {'type': 'object', 'properties': {}},
        setupInstructions: 'Setup steps',
        baseUrl: 'https://api.example.com',
        baseUrlTest: 'https://test.example.com',
        platformCredentials: {'key1': 'val1', 'key2': 'val2'},
      );

      final json = dto.toJson();

      expect(json['name'], 'Full Type');
      expect(json['code'], 'full_type');
      expect(json['description'], 'A full type');
      expect(json['icon'], 'type_icon');
      expect(json['category_id'], 5);
      expect(json['is_active'], true);
      expect(json['config_schema'], {'type': 'object'});
      expect(json['credentials_schema'],
          {'type': 'object', 'properties': {}});
      expect(json['setup_instructions'], 'Setup steps');
      expect(json['base_url'], 'https://api.example.com');
      expect(json['base_url_test'], 'https://test.example.com');
      expect(json['platform_credentials'], {'key1': 'val1', 'key2': 'val2'});
    });

    test('toJson excludes null optional fields', () {
      final dto = CreateIntegrationTypeDTO(
        name: 'Minimal',
        categoryId: 1,
      );

      final json = dto.toJson();

      expect(json.length, 2);
      expect(json.containsKey('code'), false);
      expect(json.containsKey('description'), false);
      expect(json.containsKey('icon'), false);
      expect(json.containsKey('is_active'), false);
      expect(json.containsKey('config_schema'), false);
      expect(json.containsKey('credentials_schema'), false);
      expect(json.containsKey('setup_instructions'), false);
      expect(json.containsKey('base_url'), false);
      expect(json.containsKey('base_url_test'), false);
      expect(json.containsKey('platform_credentials'), false);
    });
  });

  group('UpdateIntegrationTypeDTO', () {
    test('toJson includes all fields when provided', () {
      final dto = UpdateIntegrationTypeDTO(
        name: 'Updated Type',
        code: 'upd_type',
        description: 'Updated desc',
        icon: 'new_icon',
        categoryId: 3,
        isActive: false,
        inDevelopment: true,
        configSchema: {'fields': ['a', 'b']},
        credentialsSchema: {'required': ['token']},
        setupInstructions: 'New instructions',
        removeImage: true,
        baseUrl: 'https://new.api.com',
        baseUrlTest: 'https://new.test.com',
        platformCredentials: {'secret': 'abc'},
      );

      final json = dto.toJson();

      expect(json['name'], 'Updated Type');
      expect(json['code'], 'upd_type');
      expect(json['description'], 'Updated desc');
      expect(json['icon'], 'new_icon');
      expect(json['category_id'], 3);
      expect(json['is_active'], false);
      expect(json['in_development'], true);
      expect(json['config_schema'], {'fields': ['a', 'b']});
      expect(json['credentials_schema'], {'required': ['token']});
      expect(json['setup_instructions'], 'New instructions');
      expect(json['remove_image'], true);
      expect(json['base_url'], 'https://new.api.com');
      expect(json['base_url_test'], 'https://new.test.com');
      expect(json['platform_credentials'], {'secret': 'abc'});
    });

    test('toJson returns empty map when all fields are null', () {
      final dto = UpdateIntegrationTypeDTO();

      final json = dto.toJson();

      expect(json, isEmpty);
    });

    test('toJson includes only provided fields', () {
      final dto = UpdateIntegrationTypeDTO(
        name: 'Partial',
        isActive: true,
        removeImage: false,
      );

      final json = dto.toJson();

      expect(json.length, 3);
      expect(json['name'], 'Partial');
      expect(json['is_active'], true);
      expect(json['remove_image'], false);
    });

    test('toJson excludes null config and credentials schemas', () {
      final dto = UpdateIntegrationTypeDTO(name: 'Test');

      final json = dto.toJson();

      expect(json.containsKey('config_schema'), false);
      expect(json.containsKey('credentials_schema'), false);
      expect(json.containsKey('platform_credentials'), false);
    });

    test('toJson includes removeImage flag', () {
      final dto = UpdateIntegrationTypeDTO(removeImage: true);

      final json = dto.toJson();

      expect(json['remove_image'], true);
    });
  });

  group('WebhookInfo', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'url': 'https://example.com/webhook',
        'method': 'POST',
        'description': 'Order webhook',
        'events': ['order.created', 'order.updated'],
      };

      final webhook = WebhookInfo.fromJson(json);

      expect(webhook.url, 'https://example.com/webhook');
      expect(webhook.method, 'POST');
      expect(webhook.description, 'Order webhook');
      expect(webhook.events, isNotNull);
      expect(webhook.events!.length, 2);
      expect(webhook.events![0], 'order.created');
      expect(webhook.events![1], 'order.updated');
    });

    test('fromJson uses defaults for missing required fields', () {
      final json = <String, dynamic>{};

      final webhook = WebhookInfo.fromJson(json);

      expect(webhook.url, '');
      expect(webhook.method, '');
      expect(webhook.description, '');
    });

    test('fromJson handles null events', () {
      final json = {
        'url': 'https://example.com/webhook',
        'method': 'POST',
        'description': 'No events webhook',
      };

      final webhook = WebhookInfo.fromJson(json);

      expect(webhook.events, isNull);
    });

    test('fromJson handles empty events list', () {
      final json = {
        'url': 'https://example.com/webhook',
        'method': 'POST',
        'description': 'Empty events',
        'events': [],
      };

      final webhook = WebhookInfo.fromJson(json);

      expect(webhook.events, isNotNull);
      expect(webhook.events, isEmpty);
    });
  });

  group('SyncOrdersParams', () {
    test('toJson includes all fields when provided', () {
      final params = SyncOrdersParams(
        createdAtMin: '2026-01-01',
        createdAtMax: '2026-01-31',
        status: 'open',
        financialStatus: 'paid',
        fulfillmentStatus: 'unfulfilled',
      );

      final json = params.toJson();

      expect(json['created_at_min'], '2026-01-01');
      expect(json['created_at_max'], '2026-01-31');
      expect(json['status'], 'open');
      expect(json['financial_status'], 'paid');
      expect(json['fulfillment_status'], 'unfulfilled');
    });

    test('toJson returns empty map when all fields are null', () {
      final params = SyncOrdersParams();

      final json = params.toJson();

      expect(json, isEmpty);
    });

    test('toJson excludes null fields', () {
      final params = SyncOrdersParams(
        status: 'closed',
        financialStatus: 'refunded',
      );

      final json = params.toJson();

      expect(json.length, 2);
      expect(json['status'], 'closed');
      expect(json['financial_status'], 'refunded');
      expect(json.containsKey('created_at_min'), false);
      expect(json.containsKey('created_at_max'), false);
      expect(json.containsKey('fulfillment_status'), false);
    });

    test('toJson includes only date range when provided', () {
      final params = SyncOrdersParams(
        createdAtMin: '2026-03-01',
        createdAtMax: '2026-03-15',
      );

      final json = params.toJson();

      expect(json.length, 2);
      expect(json['created_at_min'], '2026-03-01');
      expect(json['created_at_max'], '2026-03-15');
    });
  });

  group('IntegrationSimple', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 7,
        'name': 'Quick Store',
        'type': 'ecommerce',
        'category': 'sales_channel',
        'category_name': 'Sales Channel',
        'category_color': '#3498DB',
        'image_url': 'https://example.com/store.png',
        'business_id': 4,
        'is_active': true,
      };

      final simple = IntegrationSimple.fromJson(json);

      expect(simple.id, 7);
      expect(simple.name, 'Quick Store');
      expect(simple.type, 'ecommerce');
      expect(simple.category, 'sales_channel');
      expect(simple.categoryName, 'Sales Channel');
      expect(simple.categoryColor, '#3498DB');
      expect(simple.imageUrl, 'https://example.com/store.png');
      expect(simple.businessId, 4);
      expect(simple.isActive, true);
    });

    test('fromJson uses defaults for missing required fields', () {
      final json = <String, dynamic>{};

      final simple = IntegrationSimple.fromJson(json);

      expect(simple.id, 0);
      expect(simple.name, '');
      expect(simple.type, '');
      expect(simple.category, '');
      expect(simple.categoryName, '');
      expect(simple.isActive, false);
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'name': 'Test',
        'type': 'ecommerce',
        'category': 'sales',
        'category_name': 'Sales',
        'is_active': true,
      };

      final simple = IntegrationSimple.fromJson(json);

      expect(simple.categoryColor, isNull);
      expect(simple.imageUrl, isNull);
      expect(simple.businessId, isNull);
    });

    test('toJson includes required fields', () {
      final simple = IntegrationSimple(
        id: 3,
        name: 'My Int',
        type: 'payment',
        category: 'finance',
        categoryName: 'Finance',
        isActive: true,
      );

      final json = simple.toJson();

      expect(json['id'], 3);
      expect(json['name'], 'My Int');
      expect(json['type'], 'payment');
      expect(json['category'], 'finance');
      expect(json['category_name'], 'Finance');
      expect(json['is_active'], true);
    });

    test('toJson excludes null optional fields', () {
      final simple = IntegrationSimple(
        id: 1,
        name: 'Test',
        type: 'ecommerce',
        category: 'sales',
        categoryName: 'Sales',
        isActive: false,
      );

      final json = simple.toJson();

      expect(json.containsKey('category_color'), false);
      expect(json.containsKey('image_url'), false);
      expect(json.containsKey('business_id'), false);
    });

    test('toJson includes optional fields when present', () {
      final simple = IntegrationSimple(
        id: 1,
        name: 'Test',
        type: 'ecommerce',
        category: 'sales',
        categoryName: 'Sales',
        categoryColor: '#FF0000',
        imageUrl: 'https://example.com/img.png',
        businessId: 10,
        isActive: true,
      );

      final json = simple.toJson();

      expect(json['category_color'], '#FF0000');
      expect(json['image_url'], 'https://example.com/img.png');
      expect(json['business_id'], 10);
    });
  });

  group('IntegrationCategory', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'code': 'ecommerce',
        'name': 'E-Commerce',
        'description': 'Online sales platforms',
        'icon': 'shopping_cart',
        'color': '#27AE60',
        'display_order': 1,
        'parent_category_id': 10,
        'is_active': true,
        'is_visible': true,
        'created_at': '2026-01-01T00:00:00Z',
        'updated_at': '2026-01-02T00:00:00Z',
      };

      final category = IntegrationCategory.fromJson(json);

      expect(category.id, 1);
      expect(category.code, 'ecommerce');
      expect(category.name, 'E-Commerce');
      expect(category.description, 'Online sales platforms');
      expect(category.icon, 'shopping_cart');
      expect(category.color, '#27AE60');
      expect(category.displayOrder, 1);
      expect(category.parentCategoryId, 10);
      expect(category.isActive, true);
      expect(category.isVisible, true);
      expect(category.createdAt, '2026-01-01T00:00:00Z');
      expect(category.updatedAt, '2026-01-02T00:00:00Z');
    });

    test('fromJson uses defaults for missing required fields', () {
      final json = <String, dynamic>{};

      final category = IntegrationCategory.fromJson(json);

      expect(category.id, 0);
      expect(category.code, '');
      expect(category.name, '');
      expect(category.displayOrder, 0);
      expect(category.isActive, false);
      expect(category.isVisible, false);
      expect(category.createdAt, '');
      expect(category.updatedAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 5,
        'code': 'test',
        'name': 'Test',
        'display_order': 3,
        'is_active': true,
        'is_visible': false,
        'created_at': '2026-01-01',
        'updated_at': '2026-01-01',
      };

      final category = IntegrationCategory.fromJson(json);

      expect(category.description, isNull);
      expect(category.icon, isNull);
      expect(category.color, isNull);
      expect(category.parentCategoryId, isNull);
    });

    test('toJson includes required fields', () {
      final category = IntegrationCategory(
        id: 2,
        code: 'payment',
        name: 'Payment',
        displayOrder: 2,
        isActive: true,
        isVisible: true,
        createdAt: '2026-01-01',
        updatedAt: '2026-01-02',
      );

      final json = category.toJson();

      expect(json['id'], 2);
      expect(json['code'], 'payment');
      expect(json['name'], 'Payment');
      expect(json['display_order'], 2);
      expect(json['is_active'], true);
      expect(json['is_visible'], true);
      expect(json['created_at'], '2026-01-01');
      expect(json['updated_at'], '2026-01-02');
    });

    test('toJson excludes null optional fields', () {
      final category = IntegrationCategory(
        id: 1,
        code: 'test',
        name: 'Test',
        displayOrder: 0,
        isActive: true,
        isVisible: true,
        createdAt: '2026-01-01',
        updatedAt: '2026-01-01',
      );

      final json = category.toJson();

      expect(json.containsKey('description'), false);
      expect(json.containsKey('icon'), false);
      expect(json.containsKey('color'), false);
      expect(json.containsKey('parent_category_id'), false);
    });

    test('toJson includes optional fields when present', () {
      final category = IntegrationCategory(
        id: 1,
        code: 'test',
        name: 'Test',
        description: 'Test description',
        icon: 'test_icon',
        color: '#FF0000',
        displayOrder: 1,
        parentCategoryId: 5,
        isActive: true,
        isVisible: true,
        createdAt: '2026-01-01',
        updatedAt: '2026-01-01',
      );

      final json = category.toJson();

      expect(json['description'], 'Test description');
      expect(json['icon'], 'test_icon');
      expect(json['color'], '#FF0000');
      expect(json['parent_category_id'], 5);
    });
  });
}
