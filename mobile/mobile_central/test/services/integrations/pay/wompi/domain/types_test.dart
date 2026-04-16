import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/pay/wompi/domain/types.dart';

void main() {
  group('WompiConfig', () {
    test('fromJson creates instance from empty json', () {
      final json = <String, dynamic>{};

      final config = WompiConfig.fromJson(json);

      expect(config, isA<WompiConfig>());
    });

    test('fromJson ignores unknown fields', () {
      final json = {'unknown_field': 'value', 'another': 123};

      final config = WompiConfig.fromJson(json);

      expect(config, isA<WompiConfig>());
    });

    test('toJson returns empty map', () {
      final config = WompiConfig();

      final json = config.toJson();

      expect(json, isA<Map<String, dynamic>>());
      expect(json, isEmpty);
    });

    test('fromJson/toJson roundtrip preserves empty state', () {
      final original = WompiConfig();
      final json = original.toJson();
      final restored = WompiConfig.fromJson(json);

      expect(restored.toJson(), equals(original.toJson()));
    });
  });

  group('WompiCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'private_key': 'prv_test_abc123',
        'environment': 'sandbox',
      };

      final creds = WompiCredentials.fromJson(json);

      expect(creds.privateKey, 'prv_test_abc123');
      expect(creds.environment, 'sandbox');
    });

    test('fromJson handles production environment', () {
      final json = {
        'private_key': 'prv_prod_xyz789',
        'environment': 'production',
      };

      final creds = WompiCredentials.fromJson(json);

      expect(creds.privateKey, 'prv_prod_xyz789');
      expect(creds.environment, 'production');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};

      final creds = WompiCredentials.fromJson(json);

      expect(creds.privateKey, isNull);
      expect(creds.environment, isNull);
    });

    test('fromJson handles partial fields - only privateKey', () {
      final json = {'private_key': 'prv_test_partial'};

      final creds = WompiCredentials.fromJson(json);

      expect(creds.privateKey, 'prv_test_partial');
      expect(creds.environment, isNull);
    });

    test('fromJson handles partial fields - only environment', () {
      final json = {'environment': 'sandbox'};

      final creds = WompiCredentials.fromJson(json);

      expect(creds.privateKey, isNull);
      expect(creds.environment, 'sandbox');
    });

    test('toJson includes all non-null fields', () {
      final creds = WompiCredentials(
        privateKey: 'prv_test_abc',
        environment: 'sandbox',
      );

      final json = creds.toJson();

      expect(json['private_key'], 'prv_test_abc');
      expect(json['environment'], 'sandbox');
      expect(json.length, 2);
    });

    test('toJson omits null fields', () {
      final creds = WompiCredentials();

      final json = creds.toJson();

      expect(json, isEmpty);
    });

    test('toJson omits only null privateKey', () {
      final creds = WompiCredentials(environment: 'production');

      final json = creds.toJson();

      expect(json.containsKey('private_key'), isFalse);
      expect(json['environment'], 'production');
      expect(json.length, 1);
    });

    test('toJson omits only null environment', () {
      final creds = WompiCredentials(privateKey: 'prv_test_key');

      final json = creds.toJson();

      expect(json['private_key'], 'prv_test_key');
      expect(json.containsKey('environment'), isFalse);
      expect(json.length, 1);
    });

    test('fromJson/toJson roundtrip with all fields', () {
      final original = WompiCredentials(
        privateKey: 'prv_test_roundtrip',
        environment: 'sandbox',
      );

      final json = original.toJson();
      final restored = WompiCredentials.fromJson(json);

      expect(restored.privateKey, original.privateKey);
      expect(restored.environment, original.environment);
    });

    test('fromJson/toJson roundtrip with empty credentials', () {
      final original = WompiCredentials();

      final json = original.toJson();
      final restored = WompiCredentials.fromJson(json);

      expect(restored.privateKey, isNull);
      expect(restored.environment, isNull);
    });

    test('default constructor allows all nulls', () {
      final creds = WompiCredentials();

      expect(creds.privateKey, isNull);
      expect(creds.environment, isNull);
    });
  });
}
