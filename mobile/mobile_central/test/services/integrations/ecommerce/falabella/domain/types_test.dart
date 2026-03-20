import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/ecommerce/falabella/domain/types.dart';

void main() {
  group('FalabellaConfig', () {
    test('fromJson creates instance from empty json', () {
      final json = <String, dynamic>{};

      final config = FalabellaConfig.fromJson(json);

      expect(config, isA<FalabellaConfig>());
    });

    test('fromJson ignores unknown fields', () {
      final json = {'unknown_field': 'value', 'another': 123};

      final config = FalabellaConfig.fromJson(json);

      expect(config, isA<FalabellaConfig>());
    });

    test('toJson returns empty map', () {
      final config = FalabellaConfig();

      final json = config.toJson();

      expect(json, isA<Map<String, dynamic>>());
      expect(json, isEmpty);
    });

    test('toJson roundtrip preserves structure', () {
      final config = FalabellaConfig();

      final json = config.toJson();
      final restored = FalabellaConfig.fromJson(json);

      expect(restored.toJson(), equals(json));
    });
  });

  group('FalabellaCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'api_key': 'falabella_key_xyz789',
        'user_id': 'USER_42',
      };

      final credentials = FalabellaCredentials.fromJson(json);

      expect(credentials.apiKey, 'falabella_key_xyz789');
      expect(credentials.userId, 'USER_42');
    });

    test('fromJson handles missing fields', () {
      final json = <String, dynamic>{};

      final credentials = FalabellaCredentials.fromJson(json);

      expect(credentials.apiKey, isNull);
      expect(credentials.userId, isNull);
    });

    test('fromJson handles explicit null values', () {
      final json = {
        'api_key': null,
        'user_id': null,
      };

      final credentials = FalabellaCredentials.fromJson(json);

      expect(credentials.apiKey, isNull);
      expect(credentials.userId, isNull);
    });

    test('toJson includes all non-null fields', () {
      final credentials = FalabellaCredentials(
        apiKey: 'falabella_key_xyz789',
        userId: 'USER_42',
      );

      final json = credentials.toJson();

      expect(json['api_key'], 'falabella_key_xyz789');
      expect(json['user_id'], 'USER_42');
      expect(json.length, 2);
    });

    test('toJson excludes null fields', () {
      final credentials = FalabellaCredentials();

      final json = credentials.toJson();

      expect(json, isEmpty);
    });

    test('toJson excludes only null fields', () {
      final credentials = FalabellaCredentials(apiKey: 'my_key');

      final json = credentials.toJson();

      expect(json.length, 1);
      expect(json['api_key'], 'my_key');
      expect(json.containsKey('user_id'), false);
    });

    test('toJson roundtrip preserves all fields', () {
      final original = FalabellaCredentials(
        apiKey: 'falabella_key_xyz789',
        userId: 'USER_42',
      );

      final json = original.toJson();
      final restored = FalabellaCredentials.fromJson(json);

      expect(restored.apiKey, original.apiKey);
      expect(restored.userId, original.userId);
    });

    test('toJson roundtrip with empty credentials', () {
      final original = FalabellaCredentials();

      final json = original.toJson();
      final restored = FalabellaCredentials.fromJson(json);

      expect(restored.apiKey, isNull);
      expect(restored.userId, isNull);
    });
  });
}
