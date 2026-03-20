import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/products/domain/entities.dart';

void main() {
  group('Product', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 42,
        'created_at': '2026-01-01T00:00:00Z',
        'updated_at': '2026-01-02T00:00:00Z',
        'deleted_at': '2026-01-03T00:00:00Z',
        'business_id': 5,
        'integration_id': 10,
        'integration_type': 'shopify',
        'external_id': 'ext-123',
        'sku': 'SKU-001',
        'name': 'Widget',
        'description': 'A fine widget',
        'price': 29.99,
        'compare_at_price': 39.99,
        'cost_price': 10.0,
        'currency': 'USD',
        'stock': 100,
        'stock_status': 'in_stock',
        'manage_stock': true,
        'weight': 1.5,
        'height': 10.0,
        'width': 5.0,
        'length': 8.0,
        'image_url': 'https://img.example.com/1.jpg',
        'images': ['https://img.example.com/1.jpg', 'https://img.example.com/2.jpg'],
        'thumbnail': 'https://img.example.com/thumb.jpg',
        'status': 'active',
        'is_active': true,
        'metadata': {'key': 'value'},
      };

      final product = Product.fromJson(json);

      expect(product.id, '42');
      expect(product.createdAt, '2026-01-01T00:00:00Z');
      expect(product.updatedAt, '2026-01-02T00:00:00Z');
      expect(product.deletedAt, '2026-01-03T00:00:00Z');
      expect(product.businessId, 5);
      expect(product.integrationId, 10);
      expect(product.integrationType, 'shopify');
      expect(product.externalId, 'ext-123');
      expect(product.sku, 'SKU-001');
      expect(product.name, 'Widget');
      expect(product.description, 'A fine widget');
      expect(product.price, 29.99);
      expect(product.compareAtPrice, 39.99);
      expect(product.costPrice, 10.0);
      expect(product.currency, 'USD');
      expect(product.stock, 100);
      expect(product.stockStatus, 'in_stock');
      expect(product.manageStock, true);
      expect(product.weight, 1.5);
      expect(product.height, 10.0);
      expect(product.width, 5.0);
      expect(product.length, 8.0);
      expect(product.imageUrl, 'https://img.example.com/1.jpg');
      expect(product.images, hasLength(2));
      expect(product.thumbnail, 'https://img.example.com/thumb.jpg');
      expect(product.status, 'active');
      expect(product.isActive, true);
      expect(product.metadata, isNotNull);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final product = Product.fromJson(json);

      expect(product.id, '');
      expect(product.createdAt, '');
      expect(product.updatedAt, '');
      expect(product.businessId, 0);
      expect(product.sku, '');
      expect(product.name, '');
      expect(product.price, 0.0);
      expect(product.currency, 'COP');
      expect(product.stock, 0);
      expect(product.manageStock, false);
      expect(product.status, '');
      expect(product.isActive, true);
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'business_id': 1,
        'sku': 'S',
        'name': 'N',
        'price': 10,
        'currency': 'COP',
        'stock': 0,
        'manage_stock': false,
        'status': 'draft',
        'is_active': false,
      };

      final product = Product.fromJson(json);

      expect(product.deletedAt, isNull);
      expect(product.integrationId, isNull);
      expect(product.integrationType, isNull);
      expect(product.externalId, isNull);
      expect(product.description, isNull);
      expect(product.compareAtPrice, isNull);
      expect(product.costPrice, isNull);
      expect(product.stockStatus, isNull);
      expect(product.weight, isNull);
      expect(product.height, isNull);
      expect(product.width, isNull);
      expect(product.length, isNull);
      expect(product.imageUrl, isNull);
      expect(product.images, isNull);
      expect(product.thumbnail, isNull);
      expect(product.metadata, isNull);
    });
  });

  group('ProductIntegration', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'product_id': 42,
        'integration_id': 10,
        'integration_type': 'shopify',
        'integration_name': 'My Shopify',
        'external_product_id': 'ext-prod-1',
        'created_at': '2026-01-01T00:00:00Z',
        'updated_at': '2026-01-02T00:00:00Z',
      };

      final pi = ProductIntegration.fromJson(json);

      expect(pi.id, 1);
      expect(pi.productId, '42');
      expect(pi.integrationId, 10);
      expect(pi.integrationType, 'shopify');
      expect(pi.integrationName, 'My Shopify');
      expect(pi.externalProductId, 'ext-prod-1');
      expect(pi.createdAt, '2026-01-01T00:00:00Z');
      expect(pi.updatedAt, '2026-01-02T00:00:00Z');
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final pi = ProductIntegration.fromJson(json);

      expect(pi.id, 0);
      expect(pi.productId, '');
      expect(pi.integrationId, 0);
      expect(pi.externalProductId, '');
      expect(pi.createdAt, '');
      expect(pi.updatedAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'product_id': '2',
        'integration_id': 3,
        'external_product_id': 'e',
        'created_at': 'c',
        'updated_at': 'u',
      };

      final pi = ProductIntegration.fromJson(json);

      expect(pi.integrationType, isNull);
      expect(pi.integrationName, isNull);
    });
  });

  group('GetProductsParams', () {
    test('toQueryParams includes all non-null fields', () {
      final params = GetProductsParams(
        page: 1,
        pageSize: 20,
        businessId: 5,
        integrationId: 10,
        integrationType: 'shopify',
        sku: 'SKU-001',
        skus: 'SKU-001,SKU-002',
        name: 'Widget',
        externalId: 'ext-1',
        sortBy: 'name',
        sortOrder: 'asc',
        startDate: '2026-01-01',
        endDate: '2026-12-31',
      );

      final qp = params.toQueryParams();

      expect(qp['page'], 1);
      expect(qp['page_size'], 20);
      expect(qp['business_id'], 5);
      expect(qp['integration_id'], 10);
      expect(qp['integration_type'], 'shopify');
      expect(qp['sku'], 'SKU-001');
      expect(qp['skus'], 'SKU-001,SKU-002');
      expect(qp['name'], 'Widget');
      expect(qp['external_id'], 'ext-1');
      expect(qp['sort_by'], 'name');
      expect(qp['sort_order'], 'asc');
      expect(qp['start_date'], '2026-01-01');
      expect(qp['end_date'], '2026-12-31');
    });

    test('toQueryParams excludes null fields', () {
      final params = GetProductsParams(page: 2);

      final qp = params.toQueryParams();

      expect(qp.length, 1);
      expect(qp.containsKey('page'), true);
      expect(qp.containsKey('page_size'), false);
      expect(qp.containsKey('name'), false);
    });

    test('toQueryParams returns empty map when all fields are null', () {
      final params = GetProductsParams();

      final qp = params.toQueryParams();

      expect(qp, isEmpty);
    });
  });

  group('CreateProductDTO', () {
    test('toJson includes required fields', () {
      final dto = CreateProductDTO(
        businessId: 1,
        sku: 'SKU-001',
        name: 'Widget',
        price: 29.99,
        stock: 100,
      );

      final json = dto.toJson();

      expect(json['business_id'], 1);
      expect(json['sku'], 'SKU-001');
      expect(json['name'], 'Widget');
      expect(json['price'], 29.99);
      expect(json['stock'], 100);
    });

    test('toJson includes all non-null optional fields', () {
      final dto = CreateProductDTO(
        businessId: 1,
        sku: 'SKU-001',
        name: 'Widget',
        price: 29.99,
        stock: 100,
        description: 'Desc',
        compareAtPrice: 39.99,
        costPrice: 10.0,
        currency: 'USD',
        stockStatus: 'in_stock',
        manageStock: true,
        weight: 1.5,
        height: 10.0,
        width: 5.0,
        length: 8.0,
        status: 'active',
        isActive: true,
      );

      final json = dto.toJson();

      expect(json['description'], 'Desc');
      expect(json['compare_at_price'], 39.99);
      expect(json['cost_price'], 10.0);
      expect(json['currency'], 'USD');
      expect(json['stock_status'], 'in_stock');
      expect(json['manage_stock'], true);
      expect(json['weight'], 1.5);
      expect(json['height'], 10.0);
      expect(json['width'], 5.0);
      expect(json['length'], 8.0);
      expect(json['status'], 'active');
      expect(json['is_active'], true);
    });

    test('toJson excludes null optional fields', () {
      final dto = CreateProductDTO(
        businessId: 1,
        sku: 'S',
        name: 'N',
        price: 1.0,
        stock: 0,
      );

      final json = dto.toJson();

      expect(json.length, 5);
      expect(json.containsKey('description'), false);
      expect(json.containsKey('compare_at_price'), false);
      expect(json.containsKey('currency'), false);
    });
  });

  group('UpdateProductDTO', () {
    test('toJson includes all non-null fields', () {
      final dto = UpdateProductDTO(
        sku: 'SKU-NEW',
        name: 'NewName',
        description: 'NewDesc',
        price: 49.99,
        compareAtPrice: 59.99,
        costPrice: 20.0,
        currency: 'EUR',
        stock: 50,
        stockStatus: 'low',
        manageStock: false,
        weight: 2.0,
        height: 15.0,
        width: 7.0,
        length: 12.0,
        status: 'draft',
        isActive: false,
      );

      final json = dto.toJson();

      expect(json['sku'], 'SKU-NEW');
      expect(json['name'], 'NewName');
      expect(json['description'], 'NewDesc');
      expect(json['price'], 49.99);
      expect(json['compare_at_price'], 59.99);
      expect(json['cost_price'], 20.0);
      expect(json['currency'], 'EUR');
      expect(json['stock'], 50);
      expect(json['stock_status'], 'low');
      expect(json['manage_stock'], false);
      expect(json['weight'], 2.0);
      expect(json['height'], 15.0);
      expect(json['width'], 7.0);
      expect(json['length'], 12.0);
      expect(json['status'], 'draft');
      expect(json['is_active'], false);
    });

    test('toJson returns empty map when all fields are null', () {
      final dto = UpdateProductDTO();

      final json = dto.toJson();

      expect(json, isEmpty);
    });

    test('toJson includes only provided fields', () {
      final dto = UpdateProductDTO(name: 'JustName', price: 9.99);

      final json = dto.toJson();

      expect(json.length, 2);
      expect(json['name'], 'JustName');
      expect(json['price'], 9.99);
    });
  });

  group('AddProductIntegrationDTO', () {
    test('toJson produces correct structure', () {
      final dto = AddProductIntegrationDTO(
        integrationId: 10,
        externalProductId: 'ext-123',
      );

      final json = dto.toJson();

      expect(json['integration_id'], 10);
      expect(json['external_product_id'], 'ext-123');
    });
  });
}
