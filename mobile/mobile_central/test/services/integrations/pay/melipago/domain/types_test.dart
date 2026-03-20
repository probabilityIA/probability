import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/pay/melipago/domain/types.dart';

void main() {
  group('MeliPagoConfig', () {
    test('fromJson creates instance from empty json', () {
      final json = <String, dynamic>{};

      final config = MeliPagoConfig.fromJson(json);

      expect(config, isA<MeliPagoConfig>());
    });

    test('fromJson ignores unknown fields', () {
      final json = {'unknown_field': 'value', 'another': 123};

      final config = MeliPagoConfig.fromJson(json);

      expect(config, isA<MeliPagoConfig>());
    });

    test('toJson returns empty map', () {
      final config = MeliPagoConfig();

      final json = config.toJson();

      expect(json, isA<Map<String, dynamic>>());
      expect(json, isEmpty);
    });

    test('fromJson/toJson roundtrip preserves empty state', () {
      final original = MeliPagoConfig();
      final json = original.toJson();
      final restored = MeliPagoConfig.fromJson(json);

      expect(restored.toJson(), equals(original.toJson()));
    });
  });

  group('MeliPagoCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'access_token': 'APP_USR-abc123-xyz789',
        'environment': 'sandbox',
      };

      final creds = MeliPagoCredentials.fromJson(json);

      expect(creds.accessToken, 'APP_USR-abc123-xyz789');
      expect(creds.environment, 'sandbox');
    });

    test('fromJson handles production environment', () {
      final json = {
        'access_token': 'APP_USR-prod-token',
        'environment': 'production',
      };

      final creds = MeliPagoCredentials.fromJson(json);

      expect(creds.accessToken, 'APP_USR-prod-token');
      expect(creds.environment, 'production');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};

      final creds = MeliPagoCredentials.fromJson(json);

      expect(creds.accessToken, isNull);
      expect(creds.environment, isNull);
    });

    test('fromJson handles partial fields - only accessToken', () {
      final json = {'access_token': 'APP_USR-partial'};

      final creds = MeliPagoCredentials.fromJson(json);

      expect(creds.accessToken, 'APP_USR-partial');
      expect(creds.environment, isNull);
    });

    test('fromJson handles partial fields - only environment', () {
      final json = {'environment': 'sandbox'};

      final creds = MeliPagoCredentials.fromJson(json);

      expect(creds.accessToken, isNull);
      expect(creds.environment, 'sandbox');
    });

    test('toJson includes all non-null fields', () {
      final creds = MeliPagoCredentials(
        accessToken: 'APP_USR-token',
        environment: 'sandbox',
      );

      final json = creds.toJson();

      expect(json['access_token'], 'APP_USR-token');
      expect(json['environment'], 'sandbox');
      expect(json.length, 2);
    });

    test('toJson omits null fields', () {
      final creds = MeliPagoCredentials();

      final json = creds.toJson();

      expect(json, isEmpty);
    });

    test('toJson omits only null accessToken', () {
      final creds = MeliPagoCredentials(environment: 'production');

      final json = creds.toJson();

      expect(json.containsKey('access_token'), isFalse);
      expect(json['environment'], 'production');
      expect(json.length, 1);
    });

    test('toJson omits only null environment', () {
      final creds = MeliPagoCredentials(accessToken: 'APP_USR-only');

      final json = creds.toJson();

      expect(json['access_token'], 'APP_USR-only');
      expect(json.containsKey('environment'), isFalse);
      expect(json.length, 1);
    });

    test('fromJson/toJson roundtrip with all fields', () {
      final original = MeliPagoCredentials(
        accessToken: 'APP_USR-roundtrip',
        environment: 'sandbox',
      );

      final json = original.toJson();
      final restored = MeliPagoCredentials.fromJson(json);

      expect(restored.accessToken, original.accessToken);
      expect(restored.environment, original.environment);
    });

    test('fromJson/toJson roundtrip with empty credentials', () {
      final original = MeliPagoCredentials();

      final json = original.toJson();
      final restored = MeliPagoCredentials.fromJson(json);

      expect(restored.accessToken, isNull);
      expect(restored.environment, isNull);
    });

    test('default constructor allows all nulls', () {
      final creds = MeliPagoCredentials();

      expect(creds.accessToken, isNull);
      expect(creds.environment, isNull);
    });
  });
}
