import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/pay/nequi/domain/types.dart';

void main() {
  group('NequiConfig', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'phone_code': '+57',
      };

      final config = NequiConfig.fromJson(json);

      expect(config.phoneCode, '+57');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};

      final config = NequiConfig.fromJson(json);

      expect(config.phoneCode, isNull);
    });

    test('fromJson ignores unknown fields', () {
      final json = {
        'phone_code': '+57',
        'unknown_field': 'value',
      };

      final config = NequiConfig.fromJson(json);

      expect(config.phoneCode, '+57');
    });

    test('toJson includes non-null phoneCode', () {
      final config = NequiConfig(phoneCode: '+57');

      final json = config.toJson();

      expect(json['phone_code'], '+57');
      expect(json.length, 1);
    });

    test('toJson omits null phoneCode', () {
      final config = NequiConfig();

      final json = config.toJson();

      expect(json, isEmpty);
    });

    test('fromJson/toJson roundtrip with phoneCode', () {
      final original = NequiConfig(phoneCode: '+57');

      final json = original.toJson();
      final restored = NequiConfig.fromJson(json);

      expect(restored.phoneCode, original.phoneCode);
    });

    test('fromJson/toJson roundtrip with empty config', () {
      final original = NequiConfig();

      final json = original.toJson();
      final restored = NequiConfig.fromJson(json);

      expect(restored.phoneCode, isNull);
    });

    test('default constructor allows null phoneCode', () {
      final config = NequiConfig();

      expect(config.phoneCode, isNull);
    });
  });

  group('NequiCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'api_key': 'nequi_key_abc123',
        'environment': 'sandbox',
      };

      final creds = NequiCredentials.fromJson(json);

      expect(creds.apiKey, 'nequi_key_abc123');
      expect(creds.environment, 'sandbox');
    });

    test('fromJson handles production environment', () {
      final json = {
        'api_key': 'nequi_key_prod',
        'environment': 'production',
      };

      final creds = NequiCredentials.fromJson(json);

      expect(creds.apiKey, 'nequi_key_prod');
      expect(creds.environment, 'production');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};

      final creds = NequiCredentials.fromJson(json);

      expect(creds.apiKey, isNull);
      expect(creds.environment, isNull);
    });

    test('fromJson handles partial fields - only apiKey', () {
      final json = {'api_key': 'nequi_key_partial'};

      final creds = NequiCredentials.fromJson(json);

      expect(creds.apiKey, 'nequi_key_partial');
      expect(creds.environment, isNull);
    });

    test('fromJson handles partial fields - only environment', () {
      final json = {'environment': 'sandbox'};

      final creds = NequiCredentials.fromJson(json);

      expect(creds.apiKey, isNull);
      expect(creds.environment, 'sandbox');
    });

    test('toJson includes all non-null fields', () {
      final creds = NequiCredentials(
        apiKey: 'nequi_key_abc',
        environment: 'sandbox',
      );

      final json = creds.toJson();

      expect(json['api_key'], 'nequi_key_abc');
      expect(json['environment'], 'sandbox');
      expect(json.length, 2);
    });

    test('toJson omits null fields', () {
      final creds = NequiCredentials();

      final json = creds.toJson();

      expect(json, isEmpty);
    });

    test('toJson omits only null apiKey', () {
      final creds = NequiCredentials(environment: 'production');

      final json = creds.toJson();

      expect(json.containsKey('api_key'), isFalse);
      expect(json['environment'], 'production');
      expect(json.length, 1);
    });

    test('toJson omits only null environment', () {
      final creds = NequiCredentials(apiKey: 'nequi_key_only');

      final json = creds.toJson();

      expect(json['api_key'], 'nequi_key_only');
      expect(json.containsKey('environment'), isFalse);
      expect(json.length, 1);
    });

    test('fromJson/toJson roundtrip with all fields', () {
      final original = NequiCredentials(
        apiKey: 'nequi_key_rt',
        environment: 'sandbox',
      );

      final json = original.toJson();
      final restored = NequiCredentials.fromJson(json);

      expect(restored.apiKey, original.apiKey);
      expect(restored.environment, original.environment);
    });

    test('fromJson/toJson roundtrip with empty credentials', () {
      final original = NequiCredentials();

      final json = original.toJson();
      final restored = NequiCredentials.fromJson(json);

      expect(restored.apiKey, isNull);
      expect(restored.environment, isNull);
    });

    test('default constructor allows all nulls', () {
      final creds = NequiCredentials();

      expect(creds.apiKey, isNull);
      expect(creds.environment, isNull);
    });
  });
}
