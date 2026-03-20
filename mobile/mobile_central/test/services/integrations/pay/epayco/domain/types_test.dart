import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/pay/epayco/domain/types.dart';

void main() {
  group('EPaycoConfig', () {
    test('fromJson creates instance from empty json', () {
      final json = <String, dynamic>{};

      final config = EPaycoConfig.fromJson(json);

      expect(config, isA<EPaycoConfig>());
    });

    test('fromJson ignores unknown fields', () {
      final json = {'unknown_field': 'value', 'another': 123};

      final config = EPaycoConfig.fromJson(json);

      expect(config, isA<EPaycoConfig>());
    });

    test('toJson returns empty map', () {
      final config = EPaycoConfig();

      final json = config.toJson();

      expect(json, isA<Map<String, dynamic>>());
      expect(json, isEmpty);
    });

    test('fromJson/toJson roundtrip preserves empty state', () {
      final original = EPaycoConfig();
      final json = original.toJson();
      final restored = EPaycoConfig.fromJson(json);

      expect(restored.toJson(), equals(original.toJson()));
    });
  });

  group('EPaycoCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'customer_id': 'cust_12345',
        'key': 'epayco_key_abc',
        'environment': 'test',
      };

      final creds = EPaycoCredentials.fromJson(json);

      expect(creds.customerId, 'cust_12345');
      expect(creds.key, 'epayco_key_abc');
      expect(creds.environment, 'test');
    });

    test('fromJson handles production environment', () {
      final json = {
        'customer_id': 'cust_prod',
        'key': 'epayco_key_prod',
        'environment': 'production',
      };

      final creds = EPaycoCredentials.fromJson(json);

      expect(creds.customerId, 'cust_prod');
      expect(creds.key, 'epayco_key_prod');
      expect(creds.environment, 'production');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};

      final creds = EPaycoCredentials.fromJson(json);

      expect(creds.customerId, isNull);
      expect(creds.key, isNull);
      expect(creds.environment, isNull);
    });

    test('fromJson handles partial fields - only customerId', () {
      final json = {'customer_id': 'cust_only'};

      final creds = EPaycoCredentials.fromJson(json);

      expect(creds.customerId, 'cust_only');
      expect(creds.key, isNull);
      expect(creds.environment, isNull);
    });

    test('fromJson handles partial fields - only key', () {
      final json = {'key': 'key_only'};

      final creds = EPaycoCredentials.fromJson(json);

      expect(creds.customerId, isNull);
      expect(creds.key, 'key_only');
      expect(creds.environment, isNull);
    });

    test('fromJson handles partial fields - only environment', () {
      final json = {'environment': 'test'};

      final creds = EPaycoCredentials.fromJson(json);

      expect(creds.customerId, isNull);
      expect(creds.key, isNull);
      expect(creds.environment, 'test');
    });

    test('toJson includes all non-null fields', () {
      final creds = EPaycoCredentials(
        customerId: 'cust_123',
        key: 'key_456',
        environment: 'test',
      );

      final json = creds.toJson();

      expect(json['customer_id'], 'cust_123');
      expect(json['key'], 'key_456');
      expect(json['environment'], 'test');
      expect(json.length, 3);
    });

    test('toJson omits null fields', () {
      final creds = EPaycoCredentials();

      final json = creds.toJson();

      expect(json, isEmpty);
    });

    test('toJson omits null customerId and key', () {
      final creds = EPaycoCredentials(environment: 'production');

      final json = creds.toJson();

      expect(json.containsKey('customer_id'), isFalse);
      expect(json.containsKey('key'), isFalse);
      expect(json['environment'], 'production');
      expect(json.length, 1);
    });

    test('fromJson/toJson roundtrip with all fields', () {
      final original = EPaycoCredentials(
        customerId: 'cust_rt',
        key: 'key_rt',
        environment: 'test',
      );

      final json = original.toJson();
      final restored = EPaycoCredentials.fromJson(json);

      expect(restored.customerId, original.customerId);
      expect(restored.key, original.key);
      expect(restored.environment, original.environment);
    });

    test('fromJson/toJson roundtrip with empty credentials', () {
      final original = EPaycoCredentials();

      final json = original.toJson();
      final restored = EPaycoCredentials.fromJson(json);

      expect(restored.customerId, isNull);
      expect(restored.key, isNull);
      expect(restored.environment, isNull);
    });

    test('default constructor allows all nulls', () {
      final creds = EPaycoCredentials();

      expect(creds.customerId, isNull);
      expect(creds.key, isNull);
      expect(creds.environment, isNull);
    });
  });
}
