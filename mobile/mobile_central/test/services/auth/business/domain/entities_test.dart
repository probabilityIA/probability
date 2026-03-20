import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/auth/business/domain/entities.dart';

void main() {
  group('Business', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'name': 'Test Business',
        'logo_url': 'https://example.com/logo.png',
        'primary_color': '#FF0000',
        'secondary_color': '#00FF00',
        'accent_color': '#0000FF',
        'navbar_color': '#FFFFFF',
        'navbar_image_url': 'https://example.com/navbar.png',
        'domain': 'test.com',
        'has_delivery': true,
        'has_pickup': false,
        'business_type_id': 2,
        'business_type_name': 'Retail',
        'is_active': true,
      };

      final business = Business.fromJson(json);

      expect(business.id, 1);
      expect(business.name, 'Test Business');
      expect(business.logoUrl, 'https://example.com/logo.png');
      expect(business.primaryColor, '#FF0000');
      expect(business.secondaryColor, '#00FF00');
      expect(business.accentColor, '#0000FF');
      expect(business.navbarColor, '#FFFFFF');
      expect(business.navbarImageUrl, 'https://example.com/navbar.png');
      expect(business.domain, 'test.com');
      expect(business.hasDelivery, true);
      expect(business.hasPickup, false);
      expect(business.businessTypeId, 2);
      expect(business.businessTypeName, 'Retail');
      expect(business.isActive, true);
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 5,
        'name': 'Minimal',
        'is_active': false,
      };

      final business = Business.fromJson(json);

      expect(business.id, 5);
      expect(business.name, 'Minimal');
      expect(business.logoUrl, isNull);
      expect(business.primaryColor, isNull);
      expect(business.secondaryColor, isNull);
      expect(business.accentColor, isNull);
      expect(business.navbarColor, isNull);
      expect(business.navbarImageUrl, isNull);
      expect(business.domain, isNull);
      expect(business.hasDelivery, isNull);
      expect(business.hasPickup, isNull);
      expect(business.businessTypeId, isNull);
      expect(business.businessTypeName, isNull);
      expect(business.isActive, false);
    });

    test('fromJson defaults id to 0 when missing', () {
      final json = <String, dynamic>{'name': 'No ID'};
      final business = Business.fromJson(json);
      expect(business.id, 0);
    });

    test('fromJson defaults name to empty string when missing', () {
      final json = <String, dynamic>{'id': 1};
      final business = Business.fromJson(json);
      expect(business.name, '');
    });

    test('fromJson defaults isActive to true when missing', () {
      final json = <String, dynamic>{'id': 1, 'name': 'Test'};
      final business = Business.fromJson(json);
      expect(business.isActive, true);
    });
  });

  group('BusinessSimple', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 3,
        'name': 'Simple Biz',
        'logo_url': 'https://example.com/logo.png',
        'primary_color': '#123456',
        'secondary_color': '#654321',
      };

      final business = BusinessSimple.fromJson(json);

      expect(business.id, 3);
      expect(business.name, 'Simple Biz');
      expect(business.logoUrl, 'https://example.com/logo.png');
      expect(business.primaryColor, '#123456');
      expect(business.secondaryColor, '#654321');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 7,
        'name': 'Bare Minimum',
      };

      final business = BusinessSimple.fromJson(json);

      expect(business.id, 7);
      expect(business.name, 'Bare Minimum');
      expect(business.logoUrl, isNull);
      expect(business.primaryColor, isNull);
      expect(business.secondaryColor, isNull);
    });

    test('fromJson defaults id to 0 and name to empty when missing', () {
      final json = <String, dynamic>{};
      final business = BusinessSimple.fromJson(json);
      expect(business.id, 0);
      expect(business.name, '');
    });
  });

  group('BusinessType', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 10,
        'name': 'Retail',
        'code': 'RET',
        'description': 'Retail businesses',
        'icon': 'store',
      };

      final type = BusinessType.fromJson(json);

      expect(type.id, 10);
      expect(type.name, 'Retail');
      expect(type.code, 'RET');
      expect(type.description, 'Retail businesses');
      expect(type.icon, 'store');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'name': 'Food',
      };

      final type = BusinessType.fromJson(json);

      expect(type.id, 1);
      expect(type.name, 'Food');
      expect(type.code, isNull);
      expect(type.description, isNull);
      expect(type.icon, isNull);
    });

    test('fromJson defaults id to 0 and name to empty when missing', () {
      final json = <String, dynamic>{};
      final type = BusinessType.fromJson(json);
      expect(type.id, 0);
      expect(type.name, '');
    });
  });

  group('ConfiguredResource', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 42,
        'name': 'Inventory',
        'is_active': true,
      };

      final resource = ConfiguredResource.fromJson(json);

      expect(resource.id, 42);
      expect(resource.name, 'Inventory');
      expect(resource.isActive, true);
    });

    test('fromJson handles null name', () {
      final json = {
        'id': 1,
        'is_active': false,
      };

      final resource = ConfiguredResource.fromJson(json);

      expect(resource.id, 1);
      expect(resource.name, isNull);
      expect(resource.isActive, false);
    });

    test('fromJson defaults id to 0 and isActive to false when missing', () {
      final json = <String, dynamic>{};
      final resource = ConfiguredResource.fromJson(json);
      expect(resource.id, 0);
      expect(resource.isActive, false);
    });
  });

  group('GetBusinessesParams', () {
    test('toQueryParams includes all set fields', () {
      final params = GetBusinessesParams(
        page: 2,
        pageSize: 20,
        name: 'test',
        businessTypeId: 5,
      );

      final query = params.toQueryParams();

      expect(query['page'], 2);
      expect(query['page_size'], 20);
      expect(query['name'], 'test');
      expect(query['business_type_id'], 5);
    });

    test('toQueryParams omits null fields', () {
      final params = GetBusinessesParams();
      final query = params.toQueryParams();
      expect(query, isEmpty);
    });

    test('toQueryParams omits empty name', () {
      final params = GetBusinessesParams(name: '');
      final query = params.toQueryParams();
      expect(query.containsKey('name'), false);
    });

    test('toQueryParams includes only provided fields', () {
      final params = GetBusinessesParams(page: 1, name: 'shop');
      final query = params.toQueryParams();

      expect(query['page'], 1);
      expect(query['name'], 'shop');
      expect(query.containsKey('page_size'), false);
      expect(query.containsKey('business_type_id'), false);
    });
  });

  group('CreateBusinessDTO', () {
    test('toJson includes all set fields', () {
      final dto = CreateBusinessDTO(
        name: 'New Business',
        primaryColor: '#FF0000',
        secondaryColor: '#00FF00',
        accentColor: '#0000FF',
        navbarColor: '#FFFFFF',
        domain: 'new.com',
        hasDelivery: true,
        hasPickup: false,
        businessTypeId: 3,
      );

      final json = dto.toJson();

      expect(json['name'], 'New Business');
      expect(json['primary_color'], '#FF0000');
      expect(json['secondary_color'], '#00FF00');
      expect(json['accent_color'], '#0000FF');
      expect(json['navbar_color'], '#FFFFFF');
      expect(json['domain'], 'new.com');
      expect(json['has_delivery'], true);
      expect(json['has_pickup'], false);
      expect(json['business_type_id'], 3);
    });

    test('toJson includes only name when optional fields are null', () {
      final dto = CreateBusinessDTO(name: 'Minimal');
      final json = dto.toJson();

      expect(json, {'name': 'Minimal'});
      expect(json.containsKey('primary_color'), false);
      expect(json.containsKey('secondary_color'), false);
      expect(json.containsKey('accent_color'), false);
      expect(json.containsKey('navbar_color'), false);
      expect(json.containsKey('domain'), false);
      expect(json.containsKey('has_delivery'), false);
      expect(json.containsKey('has_pickup'), false);
      expect(json.containsKey('business_type_id'), false);
    });

    test('toJson always includes name', () {
      final dto = CreateBusinessDTO(name: 'Test');
      final json = dto.toJson();
      expect(json.containsKey('name'), true);
      expect(json['name'], 'Test');
    });
  });

  group('UpdateBusinessDTO', () {
    test('toJson includes all set fields', () {
      final dto = UpdateBusinessDTO(
        name: 'Updated',
        primaryColor: '#111111',
        secondaryColor: '#222222',
        accentColor: '#333333',
        navbarColor: '#444444',
        domain: 'updated.com',
        hasDelivery: false,
        hasPickup: true,
        businessTypeId: 7,
      );

      final json = dto.toJson();

      expect(json['name'], 'Updated');
      expect(json['primary_color'], '#111111');
      expect(json['secondary_color'], '#222222');
      expect(json['accent_color'], '#333333');
      expect(json['navbar_color'], '#444444');
      expect(json['domain'], 'updated.com');
      expect(json['has_delivery'], false);
      expect(json['has_pickup'], true);
      expect(json['business_type_id'], 7);
    });

    test('toJson returns empty map when all fields are null', () {
      final dto = UpdateBusinessDTO();
      final json = dto.toJson();
      expect(json, isEmpty);
    });

    test('toJson includes only provided fields', () {
      final dto = UpdateBusinessDTO(name: 'Partial', hasDelivery: true);
      final json = dto.toJson();

      expect(json.length, 2);
      expect(json['name'], 'Partial');
      expect(json['has_delivery'], true);
    });
  });
}
