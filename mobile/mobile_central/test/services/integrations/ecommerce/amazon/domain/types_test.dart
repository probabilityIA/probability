import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/ecommerce/amazon/domain/types.dart';

void main() {
  group('AmazonConfig', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'marketplace_id': 'ATVPDKIKX0DER',
        'region': 'us-east-1',
      };

      final config = AmazonConfig.fromJson(json);

      expect(config.marketplaceId, 'ATVPDKIKX0DER');
      expect(config.region, 'us-east-1');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};

      final config = AmazonConfig.fromJson(json);

      expect(config.marketplaceId, isNull);
      expect(config.region, isNull);
    });

    test('fromJson handles explicit null values', () {
      final json = {
        'marketplace_id': null,
        'region': null,
      };

      final config = AmazonConfig.fromJson(json);

      expect(config.marketplaceId, isNull);
      expect(config.region, isNull);
    });

    test('toJson includes all non-null fields', () {
      final config = AmazonConfig(
        marketplaceId: 'ATVPDKIKX0DER',
        region: 'us-east-1',
      );

      final json = config.toJson();

      expect(json['marketplace_id'], 'ATVPDKIKX0DER');
      expect(json['region'], 'us-east-1');
      expect(json.length, 2);
    });

    test('toJson excludes null fields', () {
      final config = AmazonConfig();

      final json = config.toJson();

      expect(json, isEmpty);
    });

    test('toJson excludes only null fields', () {
      final config = AmazonConfig(marketplaceId: 'MKTP123');

      final json = config.toJson();

      expect(json.length, 1);
      expect(json['marketplace_id'], 'MKTP123');
      expect(json.containsKey('region'), false);
    });

    test('toJson roundtrip preserves all fields', () {
      final original = AmazonConfig(
        marketplaceId: 'ATVPDKIKX0DER',
        region: 'us-east-1',
      );

      final json = original.toJson();
      final restored = AmazonConfig.fromJson(json);

      expect(restored.marketplaceId, original.marketplaceId);
      expect(restored.region, original.region);
    });

    test('toJson roundtrip with empty config', () {
      final original = AmazonConfig();

      final json = original.toJson();
      final restored = AmazonConfig.fromJson(json);

      expect(restored.marketplaceId, isNull);
      expect(restored.region, isNull);
    });
  });

  group('AmazonCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'seller_id': 'A1B2C3D4E5F6G7',
        'refresh_token': 'Atzr|IwEBIExampleRefreshToken',
      };

      final credentials = AmazonCredentials.fromJson(json);

      expect(credentials.sellerId, 'A1B2C3D4E5F6G7');
      expect(credentials.refreshToken, 'Atzr|IwEBIExampleRefreshToken');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};

      final credentials = AmazonCredentials.fromJson(json);

      expect(credentials.sellerId, isNull);
      expect(credentials.refreshToken, isNull);
    });

    test('fromJson handles explicit null values', () {
      final json = {
        'seller_id': null,
        'refresh_token': null,
      };

      final credentials = AmazonCredentials.fromJson(json);

      expect(credentials.sellerId, isNull);
      expect(credentials.refreshToken, isNull);
    });

    test('toJson includes all non-null fields', () {
      final credentials = AmazonCredentials(
        sellerId: 'A1B2C3D4E5F6G7',
        refreshToken: 'Atzr|IwEBIExampleRefreshToken',
      );

      final json = credentials.toJson();

      expect(json['seller_id'], 'A1B2C3D4E5F6G7');
      expect(json['refresh_token'], 'Atzr|IwEBIExampleRefreshToken');
      expect(json.length, 2);
    });

    test('toJson excludes null fields', () {
      final credentials = AmazonCredentials();

      final json = credentials.toJson();

      expect(json, isEmpty);
    });

    test('toJson excludes only null fields', () {
      final credentials = AmazonCredentials(sellerId: 'SELLER123');

      final json = credentials.toJson();

      expect(json.length, 1);
      expect(json['seller_id'], 'SELLER123');
      expect(json.containsKey('refresh_token'), false);
    });

    test('toJson roundtrip preserves all fields', () {
      final original = AmazonCredentials(
        sellerId: 'A1B2C3D4E5F6G7',
        refreshToken: 'Atzr|IwEBIExampleRefreshToken',
      );

      final json = original.toJson();
      final restored = AmazonCredentials.fromJson(json);

      expect(restored.sellerId, original.sellerId);
      expect(restored.refreshToken, original.refreshToken);
    });

    test('toJson roundtrip with empty credentials', () {
      final original = AmazonCredentials();

      final json = original.toJson();
      final restored = AmazonCredentials.fromJson(json);

      expect(restored.sellerId, isNull);
      expect(restored.refreshToken, isNull);
    });
  });
}
