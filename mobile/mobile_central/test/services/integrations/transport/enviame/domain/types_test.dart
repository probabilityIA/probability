import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/transport/enviame/domain/types.dart';

void main() {
  group('EnviameConfig', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'base_url': 'https://api.enviame.io/v1',
      };

      final config = EnviameConfig.fromJson(json);

      expect(config.baseUrl, 'https://api.enviame.io/v1');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};

      final config = EnviameConfig.fromJson(json);

      expect(config.baseUrl, isNull);
    });

    test('fromJson ignores unknown fields', () {
      final json = {
        'base_url': 'https://api.enviame.io/v1',
        'unknown_field': 'value',
      };

      final config = EnviameConfig.fromJson(json);

      expect(config.baseUrl, 'https://api.enviame.io/v1');
    });

    test('toJson includes non-null baseUrl', () {
      final config = EnviameConfig(baseUrl: 'https://api.enviame.io/v1');

      final json = config.toJson();

      expect(json['base_url'], 'https://api.enviame.io/v1');
      expect(json.length, 1);
    });

    test('toJson omits null baseUrl', () {
      final config = EnviameConfig();

      final json = config.toJson();

      expect(json, isEmpty);
    });

    test('fromJson/toJson roundtrip with baseUrl', () {
      final original = EnviameConfig(baseUrl: 'https://api.enviame.io/v1');

      final json = original.toJson();
      final restored = EnviameConfig.fromJson(json);

      expect(restored.baseUrl, original.baseUrl);
    });

    test('fromJson/toJson roundtrip with empty config', () {
      final original = EnviameConfig();

      final json = original.toJson();
      final restored = EnviameConfig.fromJson(json);

      expect(restored.baseUrl, isNull);
    });

    test('default constructor allows null baseUrl', () {
      final config = EnviameConfig();

      expect(config.baseUrl, isNull);
    });
  });

  group('EnviameCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'api_key': 'enviame_key_abc123',
      };

      final creds = EnviameCredentials.fromJson(json);

      expect(creds.apiKey, 'enviame_key_abc123');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};

      final creds = EnviameCredentials.fromJson(json);

      expect(creds.apiKey, isNull);
    });

    test('fromJson ignores unknown fields', () {
      final json = {
        'api_key': 'enviame_key',
        'extra': 'ignored',
      };

      final creds = EnviameCredentials.fromJson(json);

      expect(creds.apiKey, 'enviame_key');
    });

    test('toJson includes non-null apiKey', () {
      final creds = EnviameCredentials(apiKey: 'enviame_key_abc');

      final json = creds.toJson();

      expect(json['api_key'], 'enviame_key_abc');
      expect(json.length, 1);
    });

    test('toJson omits null apiKey', () {
      final creds = EnviameCredentials();

      final json = creds.toJson();

      expect(json, isEmpty);
    });

    test('fromJson/toJson roundtrip with apiKey', () {
      final original = EnviameCredentials(apiKey: 'enviame_key_rt');

      final json = original.toJson();
      final restored = EnviameCredentials.fromJson(json);

      expect(restored.apiKey, original.apiKey);
    });

    test('fromJson/toJson roundtrip with empty credentials', () {
      final original = EnviameCredentials();

      final json = original.toJson();
      final restored = EnviameCredentials.fromJson(json);

      expect(restored.apiKey, isNull);
    });

    test('default constructor allows null apiKey', () {
      final creds = EnviameCredentials();

      expect(creds.apiKey, isNull);
    });
  });
}
