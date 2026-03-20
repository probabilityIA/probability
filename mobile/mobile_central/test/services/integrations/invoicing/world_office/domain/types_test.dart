import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/invoicing/world_office/domain/types.dart';

void main() {
  group('WorldOfficeConfig', () {
    test('fromJson creates instance from empty json', () {
      final json = <String, dynamic>{};

      final config = WorldOfficeConfig.fromJson(json);

      expect(config, isA<WorldOfficeConfig>());
    });

    test('fromJson ignores unknown fields', () {
      final json = {'unknown_field': 'value', 'another': 123};

      final config = WorldOfficeConfig.fromJson(json);

      expect(config, isA<WorldOfficeConfig>());
    });

    test('toJson returns empty map', () {
      final config = WorldOfficeConfig();

      final json = config.toJson();

      expect(json, isA<Map<String, dynamic>>());
      expect(json, isEmpty);
    });

    test('toJson roundtrip preserves structure', () {
      final config = WorldOfficeConfig();

      final json = config.toJson();
      final restored = WorldOfficeConfig.fromJson(json);

      expect(restored.toJson(), equals(json));
    });
  });

  group('WorldOfficeCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'username': 'wo_admin',
        'password': 'wo_secure_pass',
        'company_code': 'WO-COMP-001',
        'base_url': 'https://api.worldoffice.com.co',
      };

      final credentials = WorldOfficeCredentials.fromJson(json);

      expect(credentials.username, 'wo_admin');
      expect(credentials.password, 'wo_secure_pass');
      expect(credentials.companyCode, 'WO-COMP-001');
      expect(credentials.baseUrl, 'https://api.worldoffice.com.co');
    });

    test('fromJson handles missing fields', () {
      final json = <String, dynamic>{};

      final credentials = WorldOfficeCredentials.fromJson(json);

      expect(credentials.username, isNull);
      expect(credentials.password, isNull);
      expect(credentials.companyCode, isNull);
      expect(credentials.baseUrl, isNull);
    });

    test('fromJson handles explicit null values', () {
      final json = {
        'username': null,
        'password': null,
        'company_code': null,
        'base_url': null,
      };

      final credentials = WorldOfficeCredentials.fromJson(json);

      expect(credentials.username, isNull);
      expect(credentials.password, isNull);
      expect(credentials.companyCode, isNull);
      expect(credentials.baseUrl, isNull);
    });

    test('toJson includes all non-null fields', () {
      final credentials = WorldOfficeCredentials(
        username: 'wo_admin',
        password: 'wo_secure_pass',
        companyCode: 'WO-COMP-001',
        baseUrl: 'https://api.worldoffice.com.co',
      );

      final json = credentials.toJson();

      expect(json['username'], 'wo_admin');
      expect(json['password'], 'wo_secure_pass');
      expect(json['company_code'], 'WO-COMP-001');
      expect(json['base_url'], 'https://api.worldoffice.com.co');
      expect(json.length, 4);
    });

    test('toJson excludes null fields', () {
      final credentials = WorldOfficeCredentials();

      final json = credentials.toJson();

      expect(json, isEmpty);
    });

    test('toJson excludes only null fields', () {
      final credentials = WorldOfficeCredentials(
        username: 'admin',
        baseUrl: 'https://api.worldoffice.com.co',
      );

      final json = credentials.toJson();

      expect(json.length, 2);
      expect(json['username'], 'admin');
      expect(json['base_url'], 'https://api.worldoffice.com.co');
      expect(json.containsKey('password'), false);
      expect(json.containsKey('company_code'), false);
    });

    test('toJson roundtrip preserves all fields', () {
      final original = WorldOfficeCredentials(
        username: 'wo_admin',
        password: 'wo_secure_pass',
        companyCode: 'WO-COMP-001',
        baseUrl: 'https://api.worldoffice.com.co',
      );

      final json = original.toJson();
      final restored = WorldOfficeCredentials.fromJson(json);

      expect(restored.username, original.username);
      expect(restored.password, original.password);
      expect(restored.companyCode, original.companyCode);
      expect(restored.baseUrl, original.baseUrl);
    });

    test('toJson roundtrip with empty credentials', () {
      final original = WorldOfficeCredentials();

      final json = original.toJson();
      final restored = WorldOfficeCredentials.fromJson(json);

      expect(restored.username, isNull);
      expect(restored.password, isNull);
      expect(restored.companyCode, isNull);
      expect(restored.baseUrl, isNull);
    });
  });
}
