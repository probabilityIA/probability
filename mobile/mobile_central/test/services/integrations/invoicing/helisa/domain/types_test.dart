import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/invoicing/helisa/domain/types.dart';

void main() {
  group('HelisaConfig', () {
    test('fromJson creates instance from empty json', () {
      final json = <String, dynamic>{};

      final config = HelisaConfig.fromJson(json);

      expect(config, isA<HelisaConfig>());
    });

    test('fromJson ignores unknown fields', () {
      final json = {'unknown_field': 'value', 'another': 123};

      final config = HelisaConfig.fromJson(json);

      expect(config, isA<HelisaConfig>());
    });

    test('toJson returns empty map', () {
      final config = HelisaConfig();

      final json = config.toJson();

      expect(json, isA<Map<String, dynamic>>());
      expect(json, isEmpty);
    });

    test('toJson roundtrip preserves structure', () {
      final config = HelisaConfig();

      final json = config.toJson();
      final restored = HelisaConfig.fromJson(json);

      expect(restored.toJson(), equals(json));
    });
  });

  group('HelisaCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'username': 'helisa_admin',
        'password': 'secure_pass_123',
        'company_id': 'COMP-ABC',
        'base_url': 'https://api.helisa.com',
      };

      final credentials = HelisaCredentials.fromJson(json);

      expect(credentials.username, 'helisa_admin');
      expect(credentials.password, 'secure_pass_123');
      expect(credentials.companyId, 'COMP-ABC');
      expect(credentials.baseUrl, 'https://api.helisa.com');
    });

    test('fromJson handles missing fields', () {
      final json = <String, dynamic>{};

      final credentials = HelisaCredentials.fromJson(json);

      expect(credentials.username, isNull);
      expect(credentials.password, isNull);
      expect(credentials.companyId, isNull);
      expect(credentials.baseUrl, isNull);
    });

    test('fromJson handles explicit null values', () {
      final json = {
        'username': null,
        'password': null,
        'company_id': null,
        'base_url': null,
      };

      final credentials = HelisaCredentials.fromJson(json);

      expect(credentials.username, isNull);
      expect(credentials.password, isNull);
      expect(credentials.companyId, isNull);
      expect(credentials.baseUrl, isNull);
    });

    test('toJson includes all non-null fields', () {
      final credentials = HelisaCredentials(
        username: 'helisa_admin',
        password: 'secure_pass_123',
        companyId: 'COMP-ABC',
        baseUrl: 'https://api.helisa.com',
      );

      final json = credentials.toJson();

      expect(json['username'], 'helisa_admin');
      expect(json['password'], 'secure_pass_123');
      expect(json['company_id'], 'COMP-ABC');
      expect(json['base_url'], 'https://api.helisa.com');
      expect(json.length, 4);
    });

    test('toJson excludes null fields', () {
      final credentials = HelisaCredentials();

      final json = credentials.toJson();

      expect(json, isEmpty);
    });

    test('toJson excludes only null fields', () {
      final credentials = HelisaCredentials(
        username: 'admin',
        companyId: 'COMP-1',
      );

      final json = credentials.toJson();

      expect(json.length, 2);
      expect(json['username'], 'admin');
      expect(json['company_id'], 'COMP-1');
      expect(json.containsKey('password'), false);
      expect(json.containsKey('base_url'), false);
    });

    test('toJson roundtrip preserves all fields', () {
      final original = HelisaCredentials(
        username: 'helisa_admin',
        password: 'secure_pass_123',
        companyId: 'COMP-ABC',
        baseUrl: 'https://api.helisa.com',
      );

      final json = original.toJson();
      final restored = HelisaCredentials.fromJson(json);

      expect(restored.username, original.username);
      expect(restored.password, original.password);
      expect(restored.companyId, original.companyId);
      expect(restored.baseUrl, original.baseUrl);
    });

    test('toJson roundtrip with empty credentials', () {
      final original = HelisaCredentials();

      final json = original.toJson();
      final restored = HelisaCredentials.fromJson(json);

      expect(restored.username, isNull);
      expect(restored.password, isNull);
      expect(restored.companyId, isNull);
      expect(restored.baseUrl, isNull);
    });
  });
}
