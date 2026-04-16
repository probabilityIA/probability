import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/pay/stripe/domain/types.dart';

void main() {
  group('StripeConfig', () {
    test('fromJson creates instance from empty json', () {
      final json = <String, dynamic>{};

      final config = StripeConfig.fromJson(json);

      expect(config, isA<StripeConfig>());
    });

    test('fromJson ignores unknown fields', () {
      final json = {'unknown_field': 'value', 'another': 123};

      final config = StripeConfig.fromJson(json);

      expect(config, isA<StripeConfig>());
    });

    test('toJson returns empty map', () {
      final config = StripeConfig();

      final json = config.toJson();

      expect(json, isA<Map<String, dynamic>>());
      expect(json, isEmpty);
    });

    test('fromJson/toJson roundtrip preserves empty state', () {
      final original = StripeConfig();
      final json = original.toJson();
      final restored = StripeConfig.fromJson(json);

      expect(restored.toJson(), equals(original.toJson()));
    });
  });

  group('StripeCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'secret_key': 'sk_test_abc123',
        'environment': 'test',
      };

      final creds = StripeCredentials.fromJson(json);

      expect(creds.secretKey, 'sk_test_abc123');
      expect(creds.environment, 'test');
    });

    test('fromJson handles live environment', () {
      final json = {
        'secret_key': 'sk_live_xyz789',
        'environment': 'live',
      };

      final creds = StripeCredentials.fromJson(json);

      expect(creds.secretKey, 'sk_live_xyz789');
      expect(creds.environment, 'live');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};

      final creds = StripeCredentials.fromJson(json);

      expect(creds.secretKey, isNull);
      expect(creds.environment, isNull);
    });

    test('fromJson handles partial fields - only secretKey', () {
      final json = {'secret_key': 'sk_test_partial'};

      final creds = StripeCredentials.fromJson(json);

      expect(creds.secretKey, 'sk_test_partial');
      expect(creds.environment, isNull);
    });

    test('fromJson handles partial fields - only environment', () {
      final json = {'environment': 'test'};

      final creds = StripeCredentials.fromJson(json);

      expect(creds.secretKey, isNull);
      expect(creds.environment, 'test');
    });

    test('toJson includes all non-null fields', () {
      final creds = StripeCredentials(
        secretKey: 'sk_test_abc',
        environment: 'test',
      );

      final json = creds.toJson();

      expect(json['secret_key'], 'sk_test_abc');
      expect(json['environment'], 'test');
      expect(json.length, 2);
    });

    test('toJson omits null fields', () {
      final creds = StripeCredentials();

      final json = creds.toJson();

      expect(json, isEmpty);
    });

    test('toJson omits only null secretKey', () {
      final creds = StripeCredentials(environment: 'live');

      final json = creds.toJson();

      expect(json.containsKey('secret_key'), isFalse);
      expect(json['environment'], 'live');
      expect(json.length, 1);
    });

    test('toJson omits only null environment', () {
      final creds = StripeCredentials(secretKey: 'sk_test_key');

      final json = creds.toJson();

      expect(json['secret_key'], 'sk_test_key');
      expect(json.containsKey('environment'), isFalse);
      expect(json.length, 1);
    });

    test('fromJson/toJson roundtrip with all fields', () {
      final original = StripeCredentials(
        secretKey: 'sk_test_roundtrip',
        environment: 'test',
      );

      final json = original.toJson();
      final restored = StripeCredentials.fromJson(json);

      expect(restored.secretKey, original.secretKey);
      expect(restored.environment, original.environment);
    });

    test('fromJson/toJson roundtrip with empty credentials', () {
      final original = StripeCredentials();

      final json = original.toJson();
      final restored = StripeCredentials.fromJson(json);

      expect(restored.secretKey, isNull);
      expect(restored.environment, isNull);
    });

    test('default constructor allows all nulls', () {
      final creds = StripeCredentials();

      expect(creds.secretKey, isNull);
      expect(creds.environment, isNull);
    });
  });
}
