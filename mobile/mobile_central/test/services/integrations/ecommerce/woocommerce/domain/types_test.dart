import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/ecommerce/woocommerce/domain/types.dart';

void main() {
  group('WooCommerceConfig', () {
    test('fromJson parses all fields correctly', () {
      final json = {'store_url': 'https://mystore.com'};

      final config = WooCommerceConfig.fromJson(json);

      expect(config.storeUrl, 'https://mystore.com');
    });

    test('fromJson handles missing fields', () {
      final json = <String, dynamic>{};

      final config = WooCommerceConfig.fromJson(json);

      expect(config.storeUrl, isNull);
    });

    test('fromJson handles explicit null value', () {
      final json = {'store_url': null};

      final config = WooCommerceConfig.fromJson(json);

      expect(config.storeUrl, isNull);
    });

    test('toJson includes non-null fields', () {
      final config = WooCommerceConfig(storeUrl: 'https://mystore.com');

      final json = config.toJson();

      expect(json['store_url'], 'https://mystore.com');
      expect(json.length, 1);
    });

    test('toJson excludes null fields', () {
      final config = WooCommerceConfig();

      final json = config.toJson();

      expect(json, isEmpty);
    });

    test('toJson roundtrip preserves all fields', () {
      final original = WooCommerceConfig(storeUrl: 'https://mystore.com');

      final json = original.toJson();
      final restored = WooCommerceConfig.fromJson(json);

      expect(restored.storeUrl, original.storeUrl);
    });

    test('toJson roundtrip with empty config', () {
      final original = WooCommerceConfig();

      final json = original.toJson();
      final restored = WooCommerceConfig.fromJson(json);

      expect(restored.storeUrl, isNull);
    });
  });

  group('WooCommerceCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'consumer_key': 'ck_abc123',
        'consumer_secret': 'cs_xyz789',
      };

      final credentials = WooCommerceCredentials.fromJson(json);

      expect(credentials.consumerKey, 'ck_abc123');
      expect(credentials.consumerSecret, 'cs_xyz789');
    });

    test('fromJson handles missing fields', () {
      final json = <String, dynamic>{};

      final credentials = WooCommerceCredentials.fromJson(json);

      expect(credentials.consumerKey, isNull);
      expect(credentials.consumerSecret, isNull);
    });

    test('fromJson handles explicit null values', () {
      final json = {
        'consumer_key': null,
        'consumer_secret': null,
      };

      final credentials = WooCommerceCredentials.fromJson(json);

      expect(credentials.consumerKey, isNull);
      expect(credentials.consumerSecret, isNull);
    });

    test('toJson includes all non-null fields', () {
      final credentials = WooCommerceCredentials(
        consumerKey: 'ck_abc123',
        consumerSecret: 'cs_xyz789',
      );

      final json = credentials.toJson();

      expect(json['consumer_key'], 'ck_abc123');
      expect(json['consumer_secret'], 'cs_xyz789');
      expect(json.length, 2);
    });

    test('toJson excludes null fields', () {
      final credentials = WooCommerceCredentials();

      final json = credentials.toJson();

      expect(json, isEmpty);
    });

    test('toJson excludes only null fields', () {
      final credentials = WooCommerceCredentials(consumerKey: 'ck_abc123');

      final json = credentials.toJson();

      expect(json.length, 1);
      expect(json['consumer_key'], 'ck_abc123');
      expect(json.containsKey('consumer_secret'), false);
    });

    test('toJson roundtrip preserves all fields', () {
      final original = WooCommerceCredentials(
        consumerKey: 'ck_abc123',
        consumerSecret: 'cs_xyz789',
      );

      final json = original.toJson();
      final restored = WooCommerceCredentials.fromJson(json);

      expect(restored.consumerKey, original.consumerKey);
      expect(restored.consumerSecret, original.consumerSecret);
    });

    test('toJson roundtrip with empty credentials', () {
      final original = WooCommerceCredentials();

      final json = original.toJson();
      final restored = WooCommerceCredentials.fromJson(json);

      expect(restored.consumerKey, isNull);
      expect(restored.consumerSecret, isNull);
    });
  });
}
