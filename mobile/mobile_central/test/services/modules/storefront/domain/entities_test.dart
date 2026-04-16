import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/storefront/domain/entities.dart';

void main() {
  group('StorefrontProduct', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 42,
        'name': 'Test Product',
        'description': 'A detailed description',
        'short_description': 'Short desc',
        'price': 29999.50,
        'compare_at_price': 39999.00,
        'currency': 'USD',
        'image_url': 'https://example.com/img.png',
        'images': ['https://example.com/1.png', 'https://example.com/2.png'],
        'sku': 'SKU-001',
        'stock_quantity': 100,
        'category': 'Electronics',
        'brand': 'TestBrand',
        'is_featured': true,
        'created_at': '2026-01-01T00:00:00Z',
      };

      final product = StorefrontProduct.fromJson(json);

      expect(product.id, '42');
      expect(product.name, 'Test Product');
      expect(product.description, 'A detailed description');
      expect(product.shortDescription, 'Short desc');
      expect(product.price, 29999.50);
      expect(product.compareAtPrice, 39999.00);
      expect(product.currency, 'USD');
      expect(product.imageUrl, 'https://example.com/img.png');
      expect(product.images, ['https://example.com/1.png', 'https://example.com/2.png']);
      expect(product.sku, 'SKU-001');
      expect(product.stockQuantity, 100);
      expect(product.category, 'Electronics');
      expect(product.brand, 'TestBrand');
      expect(product.isFeatured, true);
      expect(product.createdAt, '2026-01-01T00:00:00Z');
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final product = StorefrontProduct.fromJson(json);

      expect(product.id, '');
      expect(product.name, '');
      expect(product.description, '');
      expect(product.shortDescription, '');
      expect(product.price, 0.0);
      expect(product.compareAtPrice, isNull);
      expect(product.currency, 'COP');
      expect(product.imageUrl, '');
      expect(product.images, isNull);
      expect(product.sku, '');
      expect(product.stockQuantity, 0);
      expect(product.category, '');
      expect(product.brand, '');
      expect(product.isFeatured, false);
      expect(product.createdAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 'abc',
        'name': 'P',
        'compare_at_price': null,
        'images': null,
      };

      final product = StorefrontProduct.fromJson(json);

      expect(product.id, 'abc');
      expect(product.compareAtPrice, isNull);
      expect(product.images, isNull);
    });

    test('fromJson converts id to string', () {
      final json = {'id': 123};
      final product = StorefrontProduct.fromJson(json);
      expect(product.id, '123');
    });
  });

  group('StorefrontOrder', () {
    test('fromJson parses all fields including items', () {
      final json = {
        'id': 'ORD-001',
        'order_number': 'ON-12345',
        'status': 'pending',
        'total_amount': 150000.0,
        'currency': 'COP',
        'created_at': '2026-03-01T10:00:00Z',
        'items': [
          {
            'product_name': 'Item A',
            'quantity': 2,
            'unit_price': 50000.0,
            'total_price': 100000.0,
            'image_url': 'https://example.com/a.png',
          },
          {
            'product_name': 'Item B',
            'quantity': 1,
            'unit_price': 50000.0,
            'total_price': 50000.0,
          },
        ],
      };

      final order = StorefrontOrder.fromJson(json);

      expect(order.id, 'ORD-001');
      expect(order.orderNumber, 'ON-12345');
      expect(order.status, 'pending');
      expect(order.totalAmount, 150000.0);
      expect(order.currency, 'COP');
      expect(order.createdAt, '2026-03-01T10:00:00Z');
      expect(order.items.length, 2);
      expect(order.items[0].productName, 'Item A');
      expect(order.items[0].quantity, 2);
      expect(order.items[1].imageUrl, isNull);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final order = StorefrontOrder.fromJson(json);

      expect(order.id, '');
      expect(order.orderNumber, '');
      expect(order.status, '');
      expect(order.totalAmount, 0.0);
      expect(order.currency, 'COP');
      expect(order.createdAt, '');
      expect(order.items, isEmpty);
    });

    test('fromJson handles null items list', () {
      final json = {'items': null};

      final order = StorefrontOrder.fromJson(json);

      expect(order.items, isEmpty);
    });
  });

  group('StorefrontOrderItem', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'product_name': 'Widget',
        'quantity': 3,
        'unit_price': 10000.0,
        'total_price': 30000.0,
        'image_url': 'https://example.com/widget.png',
      };

      final item = StorefrontOrderItem.fromJson(json);

      expect(item.productName, 'Widget');
      expect(item.quantity, 3);
      expect(item.unitPrice, 10000.0);
      expect(item.totalPrice, 30000.0);
      expect(item.imageUrl, 'https://example.com/widget.png');
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final item = StorefrontOrderItem.fromJson(json);

      expect(item.productName, '');
      expect(item.quantity, 0);
      expect(item.unitPrice, 0.0);
      expect(item.totalPrice, 0.0);
      expect(item.imageUrl, isNull);
    });
  });

  group('CreateStorefrontOrderDTO', () {
    test('toJson includes required items', () {
      final dto = CreateStorefrontOrderDTO(
        items: [
          CreateStorefrontOrderItemDTO(productId: 'p1', quantity: 2),
          CreateStorefrontOrderItemDTO(productId: 'p2', quantity: 1),
        ],
      );

      final json = dto.toJson();

      expect(json['items'], isList);
      expect((json['items'] as List).length, 2);
      expect((json['items'] as List)[0]['product_id'], 'p1');
      expect((json['items'] as List)[0]['quantity'], 2);
      expect(json.containsKey('notes'), false);
      expect(json.containsKey('address'), false);
    });

    test('toJson includes optional notes and address', () {
      final dto = CreateStorefrontOrderDTO(
        items: [CreateStorefrontOrderItemDTO(productId: 'p1', quantity: 1)],
        notes: 'Please deliver fast',
        address: StorefrontAddress(
          firstName: 'John',
          street: 'Calle 100',
          city: 'Bogota',
        ),
      );

      final json = dto.toJson();

      expect(json['notes'], 'Please deliver fast');
      expect(json['address'], isA<Map<String, dynamic>>());
      expect(json['address']['first_name'], 'John');
    });

    test('toJson excludes null notes and address', () {
      final dto = CreateStorefrontOrderDTO(
        items: [],
        notes: null,
        address: null,
      );

      final json = dto.toJson();

      expect(json.containsKey('notes'), false);
      expect(json.containsKey('address'), false);
    });
  });

  group('CreateStorefrontOrderItemDTO', () {
    test('toJson produces correct structure', () {
      final dto = CreateStorefrontOrderItemDTO(productId: 'abc', quantity: 5);

      final json = dto.toJson();

      expect(json['product_id'], 'abc');
      expect(json['quantity'], 5);
    });
  });

  group('StorefrontAddress', () {
    test('toJson includes required fields', () {
      final address = StorefrontAddress(
        firstName: 'Maria',
        street: 'Carrera 7',
        city: 'Medellin',
      );

      final json = address.toJson();

      expect(json['first_name'], 'Maria');
      expect(json['street'], 'Carrera 7');
      expect(json['city'], 'Medellin');
      expect(json.containsKey('last_name'), false);
      expect(json.containsKey('phone'), false);
    });

    test('toJson includes all optional fields when provided', () {
      final address = StorefrontAddress(
        firstName: 'Maria',
        lastName: 'Lopez',
        phone: '+57300123456',
        street: 'Carrera 7',
        street2: 'Apt 201',
        city: 'Medellin',
        state: 'Antioquia',
        country: 'CO',
        postalCode: '050001',
        instructions: 'Ring doorbell',
      );

      final json = address.toJson();

      expect(json['first_name'], 'Maria');
      expect(json['last_name'], 'Lopez');
      expect(json['phone'], '+57300123456');
      expect(json['street'], 'Carrera 7');
      expect(json['street2'], 'Apt 201');
      expect(json['city'], 'Medellin');
      expect(json['state'], 'Antioquia');
      expect(json['country'], 'CO');
      expect(json['postal_code'], '050001');
      expect(json['instructions'], 'Ring doorbell');
    });

    test('toJson excludes null optional fields', () {
      final address = StorefrontAddress(
        firstName: 'Test',
        street: 'St',
        city: 'City',
      );

      final json = address.toJson();

      expect(json.length, 3);
    });
  });

  group('RegisterDTO', () {
    test('toJson includes required fields', () {
      final dto = RegisterDTO(
        name: 'Test User',
        email: 'test@example.com',
        password: 'secret123',
        businessCode: 'BIZ001',
      );

      final json = dto.toJson();

      expect(json['name'], 'Test User');
      expect(json['email'], 'test@example.com');
      expect(json['password'], 'secret123');
      expect(json['business_code'], 'BIZ001');
      expect(json.containsKey('phone'), false);
      expect(json.containsKey('dni'), false);
    });

    test('toJson includes optional phone and dni', () {
      final dto = RegisterDTO(
        name: 'Test',
        email: 'a@b.com',
        password: 'pass',
        businessCode: 'BC',
        phone: '123456',
        dni: '9876543',
      );

      final json = dto.toJson();

      expect(json['phone'], '123456');
      expect(json['dni'], '9876543');
    });

    test('toJson excludes null optional fields', () {
      final dto = RegisterDTO(
        name: 'T',
        email: 'e',
        password: 'p',
        businessCode: 'b',
      );

      final json = dto.toJson();

      expect(json.length, 4);
    });
  });

  group('GetCatalogParams', () {
    test('toQueryParams includes all non-null fields', () {
      final params = GetCatalogParams(
        page: 2,
        pageSize: 20,
        search: 'shoes',
        category: 'footwear',
        businessId: 5,
      );

      final queryParams = params.toQueryParams();

      expect(queryParams['page'], 2);
      expect(queryParams['page_size'], 20);
      expect(queryParams['search'], 'shoes');
      expect(queryParams['category'], 'footwear');
      expect(queryParams['business_id'], 5);
    });

    test('toQueryParams excludes null fields', () {
      final params = GetCatalogParams(page: 1);

      final queryParams = params.toQueryParams();

      expect(queryParams.length, 1);
      expect(queryParams.containsKey('page'), true);
      expect(queryParams.containsKey('search'), false);
    });

    test('toQueryParams returns empty map when all null', () {
      final params = GetCatalogParams();

      final queryParams = params.toQueryParams();

      expect(queryParams, isEmpty);
    });
  });

  group('GetOrdersParams', () {
    test('toQueryParams includes all non-null fields', () {
      final params = GetOrdersParams(
        page: 1,
        pageSize: 10,
        businessId: 3,
      );

      final queryParams = params.toQueryParams();

      expect(queryParams['page'], 1);
      expect(queryParams['page_size'], 10);
      expect(queryParams['business_id'], 3);
    });

    test('toQueryParams excludes null fields', () {
      final params = GetOrdersParams();

      final queryParams = params.toQueryParams();

      expect(queryParams, isEmpty);
    });

    test('toQueryParams includes only provided fields', () {
      final params = GetOrdersParams(page: 5);

      final queryParams = params.toQueryParams();

      expect(queryParams.length, 1);
      expect(queryParams['page'], 5);
    });
  });
}
