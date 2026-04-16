import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/ecommerce/shopify/domain/types.dart';

void main() {
  group('ShopifyConfig', () {
    test('fromJson creates instance from empty json', () {
      final json = <String, dynamic>{};

      final config = ShopifyConfig.fromJson(json);

      expect(config, isA<ShopifyConfig>());
    });

    test('fromJson ignores unknown fields', () {
      final json = {'unknown_field': 'value', 'another': 123};

      final config = ShopifyConfig.fromJson(json);

      expect(config, isA<ShopifyConfig>());
    });

    test('toJson returns empty map', () {
      final config = ShopifyConfig();

      final json = config.toJson();

      expect(json, isA<Map<String, dynamic>>());
      expect(json, isEmpty);
    });

    test('toJson roundtrip preserves structure', () {
      final config = ShopifyConfig();

      final json = config.toJson();
      final restored = ShopifyConfig.fromJson(json);

      expect(restored.toJson(), equals(json));
    });
  });

  group('ShopifyCredentials', () {
    test('fromJson creates instance from empty json', () {
      final json = <String, dynamic>{};

      final credentials = ShopifyCredentials.fromJson(json);

      expect(credentials, isA<ShopifyCredentials>());
    });

    test('fromJson ignores unknown fields', () {
      final json = {'unknown_field': 'value', 'another': 123};

      final credentials = ShopifyCredentials.fromJson(json);

      expect(credentials, isA<ShopifyCredentials>());
    });

    test('toJson returns empty map', () {
      final credentials = ShopifyCredentials();

      final json = credentials.toJson();

      expect(json, isA<Map<String, dynamic>>());
      expect(json, isEmpty);
    });

    test('toJson roundtrip preserves structure', () {
      final credentials = ShopifyCredentials();

      final json = credentials.toJson();
      final restored = ShopifyCredentials.fromJson(json);

      expect(restored.toJson(), equals(json));
    });
  });
}
