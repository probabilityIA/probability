import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/invoicing/siigo/domain/types.dart';

void main() {
  group('SiigoConfig', () {
    test('fromJson creates instance from empty json', () {
      final json = <String, dynamic>{};

      final config = SiigoConfig.fromJson(json);

      expect(config, isA<SiigoConfig>());
    });

    test('fromJson ignores unknown fields', () {
      final json = {'unknown_field': 'value', 'another': 123};

      final config = SiigoConfig.fromJson(json);

      expect(config, isA<SiigoConfig>());
    });

    test('toJson returns empty map', () {
      final config = SiigoConfig();

      final json = config.toJson();

      expect(json, isA<Map<String, dynamic>>());
      expect(json, isEmpty);
    });

    test('toJson roundtrip preserves structure', () {
      final config = SiigoConfig();

      final json = config.toJson();
      final restored = SiigoConfig.fromJson(json);

      expect(restored.toJson(), equals(json));
    });
  });

  group('SiigoCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'username': 'siigo_user@empresa.com',
        'access_key': 'siigo_access_key_abc',
        'account_id': 'ACC-12345',
        'partner_id': 'PARTNER-001',
        'base_url': 'https://api.siigo.com',
      };

      final credentials = SiigoCredentials.fromJson(json);

      expect(credentials.username, 'siigo_user@empresa.com');
      expect(credentials.accessKey, 'siigo_access_key_abc');
      expect(credentials.accountId, 'ACC-12345');
      expect(credentials.partnerId, 'PARTNER-001');
      expect(credentials.baseUrl, 'https://api.siigo.com');
    });

    test('fromJson handles missing fields', () {
      final json = <String, dynamic>{};

      final credentials = SiigoCredentials.fromJson(json);

      expect(credentials.username, isNull);
      expect(credentials.accessKey, isNull);
      expect(credentials.accountId, isNull);
      expect(credentials.partnerId, isNull);
      expect(credentials.baseUrl, isNull);
    });

    test('fromJson handles explicit null values', () {
      final json = {
        'username': null,
        'access_key': null,
        'account_id': null,
        'partner_id': null,
        'base_url': null,
      };

      final credentials = SiigoCredentials.fromJson(json);

      expect(credentials.username, isNull);
      expect(credentials.accessKey, isNull);
      expect(credentials.accountId, isNull);
      expect(credentials.partnerId, isNull);
      expect(credentials.baseUrl, isNull);
    });

    test('toJson includes all non-null fields', () {
      final credentials = SiigoCredentials(
        username: 'siigo_user@empresa.com',
        accessKey: 'siigo_access_key_abc',
        accountId: 'ACC-12345',
        partnerId: 'PARTNER-001',
        baseUrl: 'https://api.siigo.com',
      );

      final json = credentials.toJson();

      expect(json['username'], 'siigo_user@empresa.com');
      expect(json['access_key'], 'siigo_access_key_abc');
      expect(json['account_id'], 'ACC-12345');
      expect(json['partner_id'], 'PARTNER-001');
      expect(json['base_url'], 'https://api.siigo.com');
      expect(json.length, 5);
    });

    test('toJson excludes null fields', () {
      final credentials = SiigoCredentials();

      final json = credentials.toJson();

      expect(json, isEmpty);
    });

    test('toJson excludes only null fields', () {
      final credentials = SiigoCredentials(
        username: 'user@test.com',
        accessKey: 'key123',
      );

      final json = credentials.toJson();

      expect(json.length, 2);
      expect(json['username'], 'user@test.com');
      expect(json['access_key'], 'key123');
      expect(json.containsKey('account_id'), false);
      expect(json.containsKey('partner_id'), false);
      expect(json.containsKey('base_url'), false);
    });

    test('toJson roundtrip preserves all fields', () {
      final original = SiigoCredentials(
        username: 'siigo_user@empresa.com',
        accessKey: 'siigo_access_key_abc',
        accountId: 'ACC-12345',
        partnerId: 'PARTNER-001',
        baseUrl: 'https://api.siigo.com',
      );

      final json = original.toJson();
      final restored = SiigoCredentials.fromJson(json);

      expect(restored.username, original.username);
      expect(restored.accessKey, original.accessKey);
      expect(restored.accountId, original.accountId);
      expect(restored.partnerId, original.partnerId);
      expect(restored.baseUrl, original.baseUrl);
    });

    test('toJson roundtrip with empty credentials', () {
      final original = SiigoCredentials();

      final json = original.toJson();
      final restored = SiigoCredentials.fromJson(json);

      expect(restored.username, isNull);
      expect(restored.accessKey, isNull);
      expect(restored.accountId, isNull);
      expect(restored.partnerId, isNull);
      expect(restored.baseUrl, isNull);
    });
  });
}
