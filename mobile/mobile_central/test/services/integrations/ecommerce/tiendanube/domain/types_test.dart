import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/ecommerce/tiendanube/domain/types.dart';

void main() {
  group('TiendanubeConfig', () {
    test('fromJson parses all fields correctly', () {
      final json = {'store_id': '12345'};

      final config = TiendanubeConfig.fromJson(json);

      expect(config.storeId, '12345');
    });

    test('fromJson handles missing fields', () {
      final json = <String, dynamic>{};

      final config = TiendanubeConfig.fromJson(json);

      expect(config.storeId, isNull);
    });

    test('fromJson handles explicit null value', () {
      final json = {'store_id': null};

      final config = TiendanubeConfig.fromJson(json);

      expect(config.storeId, isNull);
    });

    test('toJson includes non-null fields', () {
      final config = TiendanubeConfig(storeId: '12345');

      final json = config.toJson();

      expect(json['store_id'], '12345');
      expect(json.length, 1);
    });

    test('toJson excludes null fields', () {
      final config = TiendanubeConfig();

      final json = config.toJson();

      expect(json, isEmpty);
    });

    test('toJson roundtrip preserves all fields', () {
      final original = TiendanubeConfig(storeId: '12345');

      final json = original.toJson();
      final restored = TiendanubeConfig.fromJson(json);

      expect(restored.storeId, original.storeId);
    });

    test('toJson roundtrip with empty config', () {
      final original = TiendanubeConfig();

      final json = original.toJson();
      final restored = TiendanubeConfig.fromJson(json);

      expect(restored.storeId, isNull);
    });
  });

  group('TiendanubeCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {'access_token': 'tiendanube_token_abc123'};

      final credentials = TiendanubeCredentials.fromJson(json);

      expect(credentials.accessToken, 'tiendanube_token_abc123');
    });

    test('fromJson handles missing fields', () {
      final json = <String, dynamic>{};

      final credentials = TiendanubeCredentials.fromJson(json);

      expect(credentials.accessToken, isNull);
    });

    test('fromJson handles explicit null value', () {
      final json = {'access_token': null};

      final credentials = TiendanubeCredentials.fromJson(json);

      expect(credentials.accessToken, isNull);
    });

    test('toJson includes non-null fields', () {
      final credentials =
          TiendanubeCredentials(accessToken: 'tiendanube_token_abc123');

      final json = credentials.toJson();

      expect(json['access_token'], 'tiendanube_token_abc123');
      expect(json.length, 1);
    });

    test('toJson excludes null fields', () {
      final credentials = TiendanubeCredentials();

      final json = credentials.toJson();

      expect(json, isEmpty);
    });

    test('toJson roundtrip preserves all fields', () {
      final original =
          TiendanubeCredentials(accessToken: 'tiendanube_token_abc123');

      final json = original.toJson();
      final restored = TiendanubeCredentials.fromJson(json);

      expect(restored.accessToken, original.accessToken);
    });

    test('toJson roundtrip with empty credentials', () {
      final original = TiendanubeCredentials();

      final json = original.toJson();
      final restored = TiendanubeCredentials.fromJson(json);

      expect(restored.accessToken, isNull);
    });
  });
}
