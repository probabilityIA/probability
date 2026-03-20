import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/transport/tu/domain/types.dart';

void main() {
  group('TuConfig', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'base_url': 'https://api.tu.com/v1',
      };

      final config = TuConfig.fromJson(json);

      expect(config.baseUrl, 'https://api.tu.com/v1');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};

      final config = TuConfig.fromJson(json);

      expect(config.baseUrl, isNull);
    });

    test('fromJson ignores unknown fields', () {
      final json = {
        'base_url': 'https://api.tu.com/v1',
        'unknown_field': 'value',
      };

      final config = TuConfig.fromJson(json);

      expect(config.baseUrl, 'https://api.tu.com/v1');
    });

    test('toJson includes non-null baseUrl', () {
      final config = TuConfig(baseUrl: 'https://api.tu.com/v1');

      final json = config.toJson();

      expect(json['base_url'], 'https://api.tu.com/v1');
      expect(json.length, 1);
    });

    test('toJson omits null baseUrl', () {
      final config = TuConfig();

      final json = config.toJson();

      expect(json, isEmpty);
    });

    test('fromJson/toJson roundtrip with baseUrl', () {
      final original = TuConfig(baseUrl: 'https://api.tu.com/v1');

      final json = original.toJson();
      final restored = TuConfig.fromJson(json);

      expect(restored.baseUrl, original.baseUrl);
    });

    test('fromJson/toJson roundtrip with empty config', () {
      final original = TuConfig();

      final json = original.toJson();
      final restored = TuConfig.fromJson(json);

      expect(restored.baseUrl, isNull);
    });

    test('default constructor allows null baseUrl', () {
      final config = TuConfig();

      expect(config.baseUrl, isNull);
    });
  });

  group('TuCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'api_key': 'tu_key_abc123',
      };

      final creds = TuCredentials.fromJson(json);

      expect(creds.apiKey, 'tu_key_abc123');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};

      final creds = TuCredentials.fromJson(json);

      expect(creds.apiKey, isNull);
    });

    test('fromJson ignores unknown fields', () {
      final json = {
        'api_key': 'tu_key',
        'extra': 'ignored',
      };

      final creds = TuCredentials.fromJson(json);

      expect(creds.apiKey, 'tu_key');
    });

    test('toJson includes non-null apiKey', () {
      final creds = TuCredentials(apiKey: 'tu_key_abc');

      final json = creds.toJson();

      expect(json['api_key'], 'tu_key_abc');
      expect(json.length, 1);
    });

    test('toJson omits null apiKey', () {
      final creds = TuCredentials();

      final json = creds.toJson();

      expect(json, isEmpty);
    });

    test('fromJson/toJson roundtrip with apiKey', () {
      final original = TuCredentials(apiKey: 'tu_key_rt');

      final json = original.toJson();
      final restored = TuCredentials.fromJson(json);

      expect(restored.apiKey, original.apiKey);
    });

    test('fromJson/toJson roundtrip with empty credentials', () {
      final original = TuCredentials();

      final json = original.toJson();
      final restored = TuCredentials.fromJson(json);

      expect(restored.apiKey, isNull);
    });

    test('default constructor allows null apiKey', () {
      final creds = TuCredentials();

      expect(creds.apiKey, isNull);
    });
  });
}
