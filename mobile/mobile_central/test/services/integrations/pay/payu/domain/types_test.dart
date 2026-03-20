import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/pay/payu/domain/types.dart';

void main() {
  group('PayUConfig', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'account_id': 'ACC-12345',
        'merchant_id': 'MERCH-67890',
      };

      final config = PayUConfig.fromJson(json);

      expect(config.accountId, 'ACC-12345');
      expect(config.merchantId, 'MERCH-67890');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};

      final config = PayUConfig.fromJson(json);

      expect(config.accountId, isNull);
      expect(config.merchantId, isNull);
    });

    test('fromJson handles partial fields - only accountId', () {
      final json = {'account_id': 'ACC-111'};

      final config = PayUConfig.fromJson(json);

      expect(config.accountId, 'ACC-111');
      expect(config.merchantId, isNull);
    });

    test('fromJson handles partial fields - only merchantId', () {
      final json = {'merchant_id': 'MERCH-222'};

      final config = PayUConfig.fromJson(json);

      expect(config.accountId, isNull);
      expect(config.merchantId, 'MERCH-222');
    });

    test('toJson includes all non-null fields', () {
      final config = PayUConfig(
        accountId: 'ACC-123',
        merchantId: 'MERCH-456',
      );

      final json = config.toJson();

      expect(json['account_id'], 'ACC-123');
      expect(json['merchant_id'], 'MERCH-456');
      expect(json.length, 2);
    });

    test('toJson omits null fields', () {
      final config = PayUConfig();

      final json = config.toJson();

      expect(json, isEmpty);
    });

    test('toJson omits only null accountId', () {
      final config = PayUConfig(merchantId: 'MERCH-789');

      final json = config.toJson();

      expect(json.containsKey('account_id'), isFalse);
      expect(json['merchant_id'], 'MERCH-789');
      expect(json.length, 1);
    });

    test('fromJson/toJson roundtrip with all fields', () {
      final original = PayUConfig(
        accountId: 'ACC-RT',
        merchantId: 'MERCH-RT',
      );

      final json = original.toJson();
      final restored = PayUConfig.fromJson(json);

      expect(restored.accountId, original.accountId);
      expect(restored.merchantId, original.merchantId);
    });

    test('fromJson/toJson roundtrip with empty config', () {
      final original = PayUConfig();

      final json = original.toJson();
      final restored = PayUConfig.fromJson(json);

      expect(restored.accountId, isNull);
      expect(restored.merchantId, isNull);
    });

    test('default constructor allows all nulls', () {
      final config = PayUConfig();

      expect(config.accountId, isNull);
      expect(config.merchantId, isNull);
    });
  });

  group('PayUCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'api_key': 'key_abc123',
        'api_login': 'login_xyz',
        'environment': 'sandbox',
      };

      final creds = PayUCredentials.fromJson(json);

      expect(creds.apiKey, 'key_abc123');
      expect(creds.apiLogin, 'login_xyz');
      expect(creds.environment, 'sandbox');
    });

    test('fromJson handles production environment', () {
      final json = {
        'api_key': 'key_prod',
        'api_login': 'login_prod',
        'environment': 'production',
      };

      final creds = PayUCredentials.fromJson(json);

      expect(creds.apiKey, 'key_prod');
      expect(creds.apiLogin, 'login_prod');
      expect(creds.environment, 'production');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};

      final creds = PayUCredentials.fromJson(json);

      expect(creds.apiKey, isNull);
      expect(creds.apiLogin, isNull);
      expect(creds.environment, isNull);
    });

    test('fromJson handles partial fields - only apiKey', () {
      final json = {'api_key': 'key_only'};

      final creds = PayUCredentials.fromJson(json);

      expect(creds.apiKey, 'key_only');
      expect(creds.apiLogin, isNull);
      expect(creds.environment, isNull);
    });

    test('fromJson handles partial fields - only apiLogin', () {
      final json = {'api_login': 'login_only'};

      final creds = PayUCredentials.fromJson(json);

      expect(creds.apiKey, isNull);
      expect(creds.apiLogin, 'login_only');
      expect(creds.environment, isNull);
    });

    test('toJson includes all non-null fields', () {
      final creds = PayUCredentials(
        apiKey: 'key_123',
        apiLogin: 'login_456',
        environment: 'sandbox',
      );

      final json = creds.toJson();

      expect(json['api_key'], 'key_123');
      expect(json['api_login'], 'login_456');
      expect(json['environment'], 'sandbox');
      expect(json.length, 3);
    });

    test('toJson omits null fields', () {
      final creds = PayUCredentials();

      final json = creds.toJson();

      expect(json, isEmpty);
    });

    test('toJson omits null apiKey and apiLogin', () {
      final creds = PayUCredentials(environment: 'production');

      final json = creds.toJson();

      expect(json.containsKey('api_key'), isFalse);
      expect(json.containsKey('api_login'), isFalse);
      expect(json['environment'], 'production');
      expect(json.length, 1);
    });

    test('fromJson/toJson roundtrip with all fields', () {
      final original = PayUCredentials(
        apiKey: 'key_rt',
        apiLogin: 'login_rt',
        environment: 'sandbox',
      );

      final json = original.toJson();
      final restored = PayUCredentials.fromJson(json);

      expect(restored.apiKey, original.apiKey);
      expect(restored.apiLogin, original.apiLogin);
      expect(restored.environment, original.environment);
    });

    test('fromJson/toJson roundtrip with empty credentials', () {
      final original = PayUCredentials();

      final json = original.toJson();
      final restored = PayUCredentials.fromJson(json);

      expect(restored.apiKey, isNull);
      expect(restored.apiLogin, isNull);
      expect(restored.environment, isNull);
    });

    test('default constructor allows all nulls', () {
      final creds = PayUCredentials();

      expect(creds.apiKey, isNull);
      expect(creds.apiLogin, isNull);
      expect(creds.environment, isNull);
    });
  });
}
