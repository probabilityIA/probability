import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/transport/envioclick/domain/types.dart';

void main() {
  group('EnvioClickConfig', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'use_platform_token': true,
        'base_url_test': 'https://api-test.envioclick.com',
      };

      final config = EnvioClickConfig.fromJson(json);

      expect(config.usePlatformToken, true);
      expect(config.baseUrlTest, 'https://api-test.envioclick.com');
    });

    test('fromJson handles usePlatformToken false', () {
      final json = {
        'use_platform_token': false,
        'base_url_test': 'https://api-test.envioclick.com',
      };

      final config = EnvioClickConfig.fromJson(json);

      expect(config.usePlatformToken, false);
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};

      final config = EnvioClickConfig.fromJson(json);

      expect(config.usePlatformToken, isNull);
      expect(config.baseUrlTest, isNull);
    });

    test('fromJson handles partial fields - only usePlatformToken', () {
      final json = {'use_platform_token': true};

      final config = EnvioClickConfig.fromJson(json);

      expect(config.usePlatformToken, true);
      expect(config.baseUrlTest, isNull);
    });

    test('fromJson handles partial fields - only baseUrlTest', () {
      final json = {'base_url_test': 'https://test.envioclick.com'};

      final config = EnvioClickConfig.fromJson(json);

      expect(config.usePlatformToken, isNull);
      expect(config.baseUrlTest, 'https://test.envioclick.com');
    });

    test('fromJson ignores unknown fields', () {
      final json = {
        'use_platform_token': true,
        'base_url_test': 'https://test.envioclick.com',
        'extra_field': 'ignored',
      };

      final config = EnvioClickConfig.fromJson(json);

      expect(config.usePlatformToken, true);
      expect(config.baseUrlTest, 'https://test.envioclick.com');
    });

    test('toJson includes all non-null fields', () {
      final config = EnvioClickConfig(
        usePlatformToken: true,
        baseUrlTest: 'https://api-test.envioclick.com',
      );

      final json = config.toJson();

      expect(json['use_platform_token'], true);
      expect(json['base_url_test'], 'https://api-test.envioclick.com');
      expect(json.length, 2);
    });

    test('toJson omits null fields', () {
      final config = EnvioClickConfig();

      final json = config.toJson();

      expect(json, isEmpty);
    });

    test('toJson omits only null usePlatformToken', () {
      final config = EnvioClickConfig(baseUrlTest: 'https://test.envioclick.com');

      final json = config.toJson();

      expect(json.containsKey('use_platform_token'), isFalse);
      expect(json['base_url_test'], 'https://test.envioclick.com');
      expect(json.length, 1);
    });

    test('toJson omits only null baseUrlTest', () {
      final config = EnvioClickConfig(usePlatformToken: false);

      final json = config.toJson();

      expect(json['use_platform_token'], false);
      expect(json.containsKey('base_url_test'), isFalse);
      expect(json.length, 1);
    });

    test('fromJson/toJson roundtrip with all fields', () {
      final original = EnvioClickConfig(
        usePlatformToken: true,
        baseUrlTest: 'https://api-test.envioclick.com',
      );

      final json = original.toJson();
      final restored = EnvioClickConfig.fromJson(json);

      expect(restored.usePlatformToken, original.usePlatformToken);
      expect(restored.baseUrlTest, original.baseUrlTest);
    });

    test('fromJson/toJson roundtrip with empty config', () {
      final original = EnvioClickConfig();

      final json = original.toJson();
      final restored = EnvioClickConfig.fromJson(json);

      expect(restored.usePlatformToken, isNull);
      expect(restored.baseUrlTest, isNull);
    });

    test('default constructor allows all nulls', () {
      final config = EnvioClickConfig();

      expect(config.usePlatformToken, isNull);
      expect(config.baseUrlTest, isNull);
    });
  });

  group('EnvioClickCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'api_key': 'envioclick_key_abc123',
      };

      final creds = EnvioClickCredentials.fromJson(json);

      expect(creds.apiKey, 'envioclick_key_abc123');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};

      final creds = EnvioClickCredentials.fromJson(json);

      expect(creds.apiKey, isNull);
    });

    test('fromJson ignores unknown fields', () {
      final json = {
        'api_key': 'envioclick_key',
        'extra': 'ignored',
      };

      final creds = EnvioClickCredentials.fromJson(json);

      expect(creds.apiKey, 'envioclick_key');
    });

    test('toJson includes non-null apiKey', () {
      final creds = EnvioClickCredentials(apiKey: 'envioclick_key_abc');

      final json = creds.toJson();

      expect(json['api_key'], 'envioclick_key_abc');
      expect(json.length, 1);
    });

    test('toJson omits null apiKey', () {
      final creds = EnvioClickCredentials();

      final json = creds.toJson();

      expect(json, isEmpty);
    });

    test('fromJson/toJson roundtrip with apiKey', () {
      final original = EnvioClickCredentials(apiKey: 'envioclick_rt');

      final json = original.toJson();
      final restored = EnvioClickCredentials.fromJson(json);

      expect(restored.apiKey, original.apiKey);
    });

    test('fromJson/toJson roundtrip with empty credentials', () {
      final original = EnvioClickCredentials();

      final json = original.toJson();
      final restored = EnvioClickCredentials.fromJson(json);

      expect(restored.apiKey, isNull);
    });

    test('default constructor allows null apiKey', () {
      final creds = EnvioClickCredentials();

      expect(creds.apiKey, isNull);
    });
  });
}
