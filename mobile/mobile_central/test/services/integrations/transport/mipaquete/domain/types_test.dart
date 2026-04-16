import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/transport/mipaquete/domain/types.dart';

void main() {
  group('MiPaqueteConfig', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'base_url': 'https://api.mipaquete.com/v1',
      };

      final config = MiPaqueteConfig.fromJson(json);

      expect(config.baseUrl, 'https://api.mipaquete.com/v1');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};

      final config = MiPaqueteConfig.fromJson(json);

      expect(config.baseUrl, isNull);
    });

    test('fromJson ignores unknown fields', () {
      final json = {
        'base_url': 'https://api.mipaquete.com/v1',
        'unknown_field': 'value',
      };

      final config = MiPaqueteConfig.fromJson(json);

      expect(config.baseUrl, 'https://api.mipaquete.com/v1');
    });

    test('toJson includes non-null baseUrl', () {
      final config = MiPaqueteConfig(baseUrl: 'https://api.mipaquete.com/v1');

      final json = config.toJson();

      expect(json['base_url'], 'https://api.mipaquete.com/v1');
      expect(json.length, 1);
    });

    test('toJson omits null baseUrl', () {
      final config = MiPaqueteConfig();

      final json = config.toJson();

      expect(json, isEmpty);
    });

    test('fromJson/toJson roundtrip with baseUrl', () {
      final original = MiPaqueteConfig(baseUrl: 'https://api.mipaquete.com/v1');

      final json = original.toJson();
      final restored = MiPaqueteConfig.fromJson(json);

      expect(restored.baseUrl, original.baseUrl);
    });

    test('fromJson/toJson roundtrip with empty config', () {
      final original = MiPaqueteConfig();

      final json = original.toJson();
      final restored = MiPaqueteConfig.fromJson(json);

      expect(restored.baseUrl, isNull);
    });

    test('default constructor allows null baseUrl', () {
      final config = MiPaqueteConfig();

      expect(config.baseUrl, isNull);
    });
  });

  group('MiPaqueteCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'api_key': 'mipaquete_key_abc123',
      };

      final creds = MiPaqueteCredentials.fromJson(json);

      expect(creds.apiKey, 'mipaquete_key_abc123');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};

      final creds = MiPaqueteCredentials.fromJson(json);

      expect(creds.apiKey, isNull);
    });

    test('fromJson ignores unknown fields', () {
      final json = {
        'api_key': 'mipaquete_key',
        'extra': 'ignored',
      };

      final creds = MiPaqueteCredentials.fromJson(json);

      expect(creds.apiKey, 'mipaquete_key');
    });

    test('toJson includes non-null apiKey', () {
      final creds = MiPaqueteCredentials(apiKey: 'mipaquete_key_abc');

      final json = creds.toJson();

      expect(json['api_key'], 'mipaquete_key_abc');
      expect(json.length, 1);
    });

    test('toJson omits null apiKey', () {
      final creds = MiPaqueteCredentials();

      final json = creds.toJson();

      expect(json, isEmpty);
    });

    test('fromJson/toJson roundtrip with apiKey', () {
      final original = MiPaqueteCredentials(apiKey: 'mipaquete_key_rt');

      final json = original.toJson();
      final restored = MiPaqueteCredentials.fromJson(json);

      expect(restored.apiKey, original.apiKey);
    });

    test('fromJson/toJson roundtrip with empty credentials', () {
      final original = MiPaqueteCredentials();

      final json = original.toJson();
      final restored = MiPaqueteCredentials.fromJson(json);

      expect(restored.apiKey, isNull);
    });

    test('default constructor allows null apiKey', () {
      final creds = MiPaqueteCredentials();

      expect(creds.apiKey, isNull);
    });
  });
}
