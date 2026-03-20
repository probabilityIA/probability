import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/pay/bold/domain/types.dart';

void main() {
  group('BoldConfig', () {
    test('fromJson creates instance from empty json', () {
      final json = <String, dynamic>{};

      final config = BoldConfig.fromJson(json);

      expect(config, isA<BoldConfig>());
    });

    test('fromJson ignores unknown fields', () {
      final json = {'unknown_field': 'value', 'another': 123};

      final config = BoldConfig.fromJson(json);

      expect(config, isA<BoldConfig>());
    });

    test('toJson returns empty map', () {
      final config = BoldConfig();

      final json = config.toJson();

      expect(json, isA<Map<String, dynamic>>());
      expect(json, isEmpty);
    });

    test('fromJson/toJson roundtrip preserves empty state', () {
      final original = BoldConfig();
      final json = original.toJson();
      final restored = BoldConfig.fromJson(json);

      expect(restored.toJson(), equals(original.toJson()));
    });
  });

  group('BoldCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'api_key': 'bold_key_abc123',
        'environment': 'sandbox',
      };

      final creds = BoldCredentials.fromJson(json);

      expect(creds.apiKey, 'bold_key_abc123');
      expect(creds.environment, 'sandbox');
    });

    test('fromJson handles production environment', () {
      final json = {
        'api_key': 'bold_key_prod',
        'environment': 'production',
      };

      final creds = BoldCredentials.fromJson(json);

      expect(creds.apiKey, 'bold_key_prod');
      expect(creds.environment, 'production');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};

      final creds = BoldCredentials.fromJson(json);

      expect(creds.apiKey, isNull);
      expect(creds.environment, isNull);
    });

    test('fromJson handles partial fields - only apiKey', () {
      final json = {'api_key': 'bold_key_partial'};

      final creds = BoldCredentials.fromJson(json);

      expect(creds.apiKey, 'bold_key_partial');
      expect(creds.environment, isNull);
    });

    test('fromJson handles partial fields - only environment', () {
      final json = {'environment': 'sandbox'};

      final creds = BoldCredentials.fromJson(json);

      expect(creds.apiKey, isNull);
      expect(creds.environment, 'sandbox');
    });

    test('toJson includes all non-null fields', () {
      final creds = BoldCredentials(
        apiKey: 'bold_key_abc',
        environment: 'sandbox',
      );

      final json = creds.toJson();

      expect(json['api_key'], 'bold_key_abc');
      expect(json['environment'], 'sandbox');
      expect(json.length, 2);
    });

    test('toJson omits null fields', () {
      final creds = BoldCredentials();

      final json = creds.toJson();

      expect(json, isEmpty);
    });

    test('toJson omits only null apiKey', () {
      final creds = BoldCredentials(environment: 'production');

      final json = creds.toJson();

      expect(json.containsKey('api_key'), isFalse);
      expect(json['environment'], 'production');
      expect(json.length, 1);
    });

    test('toJson omits only null environment', () {
      final creds = BoldCredentials(apiKey: 'bold_key_only');

      final json = creds.toJson();

      expect(json['api_key'], 'bold_key_only');
      expect(json.containsKey('environment'), isFalse);
      expect(json.length, 1);
    });

    test('fromJson/toJson roundtrip with all fields', () {
      final original = BoldCredentials(
        apiKey: 'bold_key_rt',
        environment: 'sandbox',
      );

      final json = original.toJson();
      final restored = BoldCredentials.fromJson(json);

      expect(restored.apiKey, original.apiKey);
      expect(restored.environment, original.environment);
    });

    test('fromJson/toJson roundtrip with empty credentials', () {
      final original = BoldCredentials();

      final json = original.toJson();
      final restored = BoldCredentials.fromJson(json);

      expect(restored.apiKey, isNull);
      expect(restored.environment, isNull);
    });

    test('default constructor allows all nulls', () {
      final creds = BoldCredentials();

      expect(creds.apiKey, isNull);
      expect(creds.environment, isNull);
    });
  });
}
