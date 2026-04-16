import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/invoicing/alegra/domain/types.dart';

void main() {
  group('AlegraConfig', () {
    test('fromJson creates instance from empty json', () {
      final json = <String, dynamic>{};

      final config = AlegraConfig.fromJson(json);

      expect(config, isA<AlegraConfig>());
    });

    test('fromJson ignores unknown fields', () {
      final json = {'unknown_field': 'value', 'another': 123};

      final config = AlegraConfig.fromJson(json);

      expect(config, isA<AlegraConfig>());
    });

    test('toJson returns empty map', () {
      final config = AlegraConfig();

      final json = config.toJson();

      expect(json, isA<Map<String, dynamic>>());
      expect(json, isEmpty);
    });

    test('toJson roundtrip preserves structure', () {
      final config = AlegraConfig();

      final json = config.toJson();
      final restored = AlegraConfig.fromJson(json);

      expect(restored.toJson(), equals(json));
    });
  });

  group('AlegraCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'email': 'admin@empresa.com',
        'token': 'alegra_token_abc123',
        'base_url': 'https://api.alegra.com',
      };

      final credentials = AlegraCredentials.fromJson(json);

      expect(credentials.email, 'admin@empresa.com');
      expect(credentials.token, 'alegra_token_abc123');
      expect(credentials.baseUrl, 'https://api.alegra.com');
    });

    test('fromJson handles missing fields', () {
      final json = <String, dynamic>{};

      final credentials = AlegraCredentials.fromJson(json);

      expect(credentials.email, isNull);
      expect(credentials.token, isNull);
      expect(credentials.baseUrl, isNull);
    });

    test('fromJson handles explicit null values', () {
      final json = {
        'email': null,
        'token': null,
        'base_url': null,
      };

      final credentials = AlegraCredentials.fromJson(json);

      expect(credentials.email, isNull);
      expect(credentials.token, isNull);
      expect(credentials.baseUrl, isNull);
    });

    test('toJson includes all non-null fields', () {
      final credentials = AlegraCredentials(
        email: 'admin@empresa.com',
        token: 'alegra_token_abc123',
        baseUrl: 'https://api.alegra.com',
      );

      final json = credentials.toJson();

      expect(json['email'], 'admin@empresa.com');
      expect(json['token'], 'alegra_token_abc123');
      expect(json['base_url'], 'https://api.alegra.com');
      expect(json.length, 3);
    });

    test('toJson excludes null fields', () {
      final credentials = AlegraCredentials();

      final json = credentials.toJson();

      expect(json, isEmpty);
    });

    test('toJson excludes only null fields', () {
      final credentials = AlegraCredentials(email: 'test@test.com');

      final json = credentials.toJson();

      expect(json.length, 1);
      expect(json['email'], 'test@test.com');
      expect(json.containsKey('token'), false);
      expect(json.containsKey('base_url'), false);
    });

    test('toJson roundtrip preserves all fields', () {
      final original = AlegraCredentials(
        email: 'admin@empresa.com',
        token: 'alegra_token_abc123',
        baseUrl: 'https://api.alegra.com',
      );

      final json = original.toJson();
      final restored = AlegraCredentials.fromJson(json);

      expect(restored.email, original.email);
      expect(restored.token, original.token);
      expect(restored.baseUrl, original.baseUrl);
    });

    test('toJson roundtrip with empty credentials', () {
      final original = AlegraCredentials();

      final json = original.toJson();
      final restored = AlegraCredentials.fromJson(json);

      expect(restored.email, isNull);
      expect(restored.token, isNull);
      expect(restored.baseUrl, isNull);
    });
  });
}
