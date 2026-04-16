import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/pay/domain/entities.dart';

void main() {
  group('PaymentGatewayType', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'name': 'Stripe',
        'code': 'stripe',
        'image_url': 'https://example.com/stripe.png',
        'is_active': true,
        'in_development': false,
      };

      final gateway = PaymentGatewayType.fromJson(json);

      expect(gateway.id, 1);
      expect(gateway.name, 'Stripe');
      expect(gateway.code, 'stripe');
      expect(gateway.imageUrl, 'https://example.com/stripe.png');
      expect(gateway.isActive, true);
      expect(gateway.inDevelopment, false);
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 2,
        'name': 'PayU',
        'code': 'payu',
        'is_active': true,
        'in_development': true,
      };

      final gateway = PaymentGatewayType.fromJson(json);

      expect(gateway.imageUrl, isNull);
    });

    test('fromJson defaults to safe values when fields are missing', () {
      final json = <String, dynamic>{};

      final gateway = PaymentGatewayType.fromJson(json);

      expect(gateway.id, 0);
      expect(gateway.name, '');
      expect(gateway.code, '');
      expect(gateway.imageUrl, isNull);
      expect(gateway.isActive, true);
      expect(gateway.inDevelopment, false);
    });

    test('fromJson handles false is_active', () {
      final json = {
        'id': 3,
        'name': 'Inactive',
        'code': 'inactive',
        'is_active': false,
        'in_development': false,
      };

      final gateway = PaymentGatewayType.fromJson(json);

      expect(gateway.isActive, false);
    });

    test('fromJson handles true in_development', () {
      final json = {
        'id': 4,
        'name': 'Beta Gateway',
        'code': 'beta',
        'is_active': true,
        'in_development': true,
      };

      final gateway = PaymentGatewayType.fromJson(json);

      expect(gateway.inDevelopment, true);
    });
  });

  group('PaymentGatewayTypesResponse', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'success': true,
        'data': [
          {
            'id': 1,
            'name': 'Stripe',
            'code': 'stripe',
            'image_url': 'https://example.com/stripe.png',
            'is_active': true,
            'in_development': false,
          },
          {
            'id': 2,
            'name': 'PayU',
            'code': 'payu',
            'is_active': true,
            'in_development': true,
          },
        ],
        'message': 'Success',
      };

      final response = PaymentGatewayTypesResponse.fromJson(json);

      expect(response.success, true);
      expect(response.data, hasLength(2));
      expect(response.data[0].name, 'Stripe');
      expect(response.data[0].code, 'stripe');
      expect(response.data[1].name, 'PayU');
      expect(response.data[1].code, 'payu');
      expect(response.message, 'Success');
    });

    test('fromJson handles null data list', () {
      final json = {
        'success': true,
        'data': null,
      };

      final response = PaymentGatewayTypesResponse.fromJson(json);

      expect(response.success, true);
      expect(response.data, isEmpty);
    });

    test('fromJson handles empty data list', () {
      final json = {
        'success': true,
        'data': [],
      };

      final response = PaymentGatewayTypesResponse.fromJson(json);

      expect(response.data, isEmpty);
    });

    test('fromJson handles null message', () {
      final json = {
        'success': true,
        'data': [],
      };

      final response = PaymentGatewayTypesResponse.fromJson(json);

      expect(response.message, isNull);
    });

    test('fromJson defaults to safe values when fields are missing', () {
      final json = <String, dynamic>{};

      final response = PaymentGatewayTypesResponse.fromJson(json);

      expect(response.success, false);
      expect(response.data, isEmpty);
      expect(response.message, isNull);
    });

    test('fromJson handles false success', () {
      final json = {
        'success': false,
        'data': [],
        'message': 'Error occurred',
      };

      final response = PaymentGatewayTypesResponse.fromJson(json);

      expect(response.success, false);
      expect(response.message, 'Error occurred');
    });
  });
}
