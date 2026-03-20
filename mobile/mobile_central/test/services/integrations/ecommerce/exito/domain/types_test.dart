import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/ecommerce/exito/domain/types.dart';

void main() {
  group('ExitoConfig', () {
    test('fromJson creates instance from empty json', () {
      final json = <String, dynamic>{};

      final config = ExitoConfig.fromJson(json);

      expect(config, isA<ExitoConfig>());
    });

    test('fromJson ignores unknown fields', () {
      final json = {'unknown_field': 'value', 'another': 123};

      final config = ExitoConfig.fromJson(json);

      expect(config, isA<ExitoConfig>());
    });

    test('toJson returns empty map', () {
      final config = ExitoConfig();

      final json = config.toJson();

      expect(json, isA<Map<String, dynamic>>());
      expect(json, isEmpty);
    });

    test('toJson roundtrip preserves structure', () {
      final config = ExitoConfig();

      final json = config.toJson();
      final restored = ExitoConfig.fromJson(json);

      expect(restored.toJson(), equals(json));
    });
  });

  group('ExitoCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'api_key': 'exito_key_abc123',
        'seller_id': 'SELLER_001',
      };

      final credentials = ExitoCredentials.fromJson(json);

      expect(credentials.apiKey, 'exito_key_abc123');
      expect(credentials.sellerId, 'SELLER_001');
    });

    test('fromJson handles missing fields', () {
      final json = <String, dynamic>{};

      final credentials = ExitoCredentials.fromJson(json);

      expect(credentials.apiKey, isNull);
      expect(credentials.sellerId, isNull);
    });

    test('fromJson handles explicit null values', () {
      final json = {
        'api_key': null,
        'seller_id': null,
      };

      final credentials = ExitoCredentials.fromJson(json);

      expect(credentials.apiKey, isNull);
      expect(credentials.sellerId, isNull);
    });

    test('toJson includes all non-null fields', () {
      final credentials = ExitoCredentials(
        apiKey: 'exito_key_abc123',
        sellerId: 'SELLER_001',
      );

      final json = credentials.toJson();

      expect(json['api_key'], 'exito_key_abc123');
      expect(json['seller_id'], 'SELLER_001');
      expect(json.length, 2);
    });

    test('toJson excludes null fields', () {
      final credentials = ExitoCredentials();

      final json = credentials.toJson();

      expect(json, isEmpty);
    });

    test('toJson excludes only null fields', () {
      final credentials = ExitoCredentials(apiKey: 'my_key');

      final json = credentials.toJson();

      expect(json.length, 1);
      expect(json['api_key'], 'my_key');
      expect(json.containsKey('seller_id'), false);
    });

    test('toJson roundtrip preserves all fields', () {
      final original = ExitoCredentials(
        apiKey: 'exito_key_abc123',
        sellerId: 'SELLER_001',
      );

      final json = original.toJson();
      final restored = ExitoCredentials.fromJson(json);

      expect(restored.apiKey, original.apiKey);
      expect(restored.sellerId, original.sellerId);
    });

    test('toJson roundtrip with empty credentials', () {
      final original = ExitoCredentials();

      final json = original.toJson();
      final restored = ExitoCredentials.fromJson(json);

      expect(restored.apiKey, isNull);
      expect(restored.sellerId, isNull);
    });
  });
}
