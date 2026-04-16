import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/ecommerce/magento/domain/types.dart';

void main() {
  group('MagentoConfig', () {
    test('fromJson parses all fields correctly', () {
      final json = {'store_url': 'https://magento.mystore.com'};

      final config = MagentoConfig.fromJson(json);

      expect(config.storeUrl, 'https://magento.mystore.com');
    });

    test('fromJson handles missing fields', () {
      final json = <String, dynamic>{};

      final config = MagentoConfig.fromJson(json);

      expect(config.storeUrl, isNull);
    });

    test('fromJson handles explicit null value', () {
      final json = {'store_url': null};

      final config = MagentoConfig.fromJson(json);

      expect(config.storeUrl, isNull);
    });

    test('toJson includes non-null fields', () {
      final config = MagentoConfig(storeUrl: 'https://magento.mystore.com');

      final json = config.toJson();

      expect(json['store_url'], 'https://magento.mystore.com');
      expect(json.length, 1);
    });

    test('toJson excludes null fields', () {
      final config = MagentoConfig();

      final json = config.toJson();

      expect(json, isEmpty);
    });

    test('toJson roundtrip preserves all fields', () {
      final original = MagentoConfig(storeUrl: 'https://magento.mystore.com');

      final json = original.toJson();
      final restored = MagentoConfig.fromJson(json);

      expect(restored.storeUrl, original.storeUrl);
    });

    test('toJson roundtrip with empty config', () {
      final original = MagentoConfig();

      final json = original.toJson();
      final restored = MagentoConfig.fromJson(json);

      expect(restored.storeUrl, isNull);
    });
  });

  group('MagentoCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {'access_token': 'magento_access_token_123'};

      final credentials = MagentoCredentials.fromJson(json);

      expect(credentials.accessToken, 'magento_access_token_123');
    });

    test('fromJson handles missing fields', () {
      final json = <String, dynamic>{};

      final credentials = MagentoCredentials.fromJson(json);

      expect(credentials.accessToken, isNull);
    });

    test('fromJson handles explicit null value', () {
      final json = {'access_token': null};

      final credentials = MagentoCredentials.fromJson(json);

      expect(credentials.accessToken, isNull);
    });

    test('toJson includes non-null fields', () {
      final credentials =
          MagentoCredentials(accessToken: 'magento_access_token_123');

      final json = credentials.toJson();

      expect(json['access_token'], 'magento_access_token_123');
      expect(json.length, 1);
    });

    test('toJson excludes null fields', () {
      final credentials = MagentoCredentials();

      final json = credentials.toJson();

      expect(json, isEmpty);
    });

    test('toJson roundtrip preserves all fields', () {
      final original =
          MagentoCredentials(accessToken: 'magento_access_token_123');

      final json = original.toJson();
      final restored = MagentoCredentials.fromJson(json);

      expect(restored.accessToken, original.accessToken);
    });

    test('toJson roundtrip with empty credentials', () {
      final original = MagentoCredentials();

      final json = original.toJson();
      final restored = MagentoCredentials.fromJson(json);

      expect(restored.accessToken, isNull);
    });
  });
}
