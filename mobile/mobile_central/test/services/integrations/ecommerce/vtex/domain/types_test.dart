import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/ecommerce/vtex/domain/types.dart';

void main() {
  group('VtexConfig', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'account_name': 'mystore',
        'environment': 'stable',
      };

      final config = VtexConfig.fromJson(json);

      expect(config.accountName, 'mystore');
      expect(config.environment, 'stable');
    });

    test('fromJson handles missing fields', () {
      final json = <String, dynamic>{};

      final config = VtexConfig.fromJson(json);

      expect(config.accountName, isNull);
      expect(config.environment, isNull);
    });

    test('fromJson handles explicit null values', () {
      final json = {
        'account_name': null,
        'environment': null,
      };

      final config = VtexConfig.fromJson(json);

      expect(config.accountName, isNull);
      expect(config.environment, isNull);
    });

    test('toJson includes all non-null fields', () {
      final config = VtexConfig(
        accountName: 'mystore',
        environment: 'stable',
      );

      final json = config.toJson();

      expect(json['account_name'], 'mystore');
      expect(json['environment'], 'stable');
      expect(json.length, 2);
    });

    test('toJson excludes null fields', () {
      final config = VtexConfig();

      final json = config.toJson();

      expect(json, isEmpty);
    });

    test('toJson excludes only null fields', () {
      final config = VtexConfig(accountName: 'mystore');

      final json = config.toJson();

      expect(json.length, 1);
      expect(json['account_name'], 'mystore');
      expect(json.containsKey('environment'), false);
    });

    test('toJson roundtrip preserves all fields', () {
      final original = VtexConfig(
        accountName: 'mystore',
        environment: 'stable',
      );

      final json = original.toJson();
      final restored = VtexConfig.fromJson(json);

      expect(restored.accountName, original.accountName);
      expect(restored.environment, original.environment);
    });

    test('toJson roundtrip with empty config', () {
      final original = VtexConfig();

      final json = original.toJson();
      final restored = VtexConfig.fromJson(json);

      expect(restored.accountName, isNull);
      expect(restored.environment, isNull);
    });
  });

  group('VtexCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'app_key': 'vtexappkey-mystore-ABCDEF',
        'app_token': 'GHIJKLMNOPQRSTUVWXYZ',
      };

      final credentials = VtexCredentials.fromJson(json);

      expect(credentials.appKey, 'vtexappkey-mystore-ABCDEF');
      expect(credentials.appToken, 'GHIJKLMNOPQRSTUVWXYZ');
    });

    test('fromJson handles missing fields', () {
      final json = <String, dynamic>{};

      final credentials = VtexCredentials.fromJson(json);

      expect(credentials.appKey, isNull);
      expect(credentials.appToken, isNull);
    });

    test('fromJson handles explicit null values', () {
      final json = {
        'app_key': null,
        'app_token': null,
      };

      final credentials = VtexCredentials.fromJson(json);

      expect(credentials.appKey, isNull);
      expect(credentials.appToken, isNull);
    });

    test('toJson includes all non-null fields', () {
      final credentials = VtexCredentials(
        appKey: 'vtexappkey-mystore-ABCDEF',
        appToken: 'GHIJKLMNOPQRSTUVWXYZ',
      );

      final json = credentials.toJson();

      expect(json['app_key'], 'vtexappkey-mystore-ABCDEF');
      expect(json['app_token'], 'GHIJKLMNOPQRSTUVWXYZ');
      expect(json.length, 2);
    });

    test('toJson excludes null fields', () {
      final credentials = VtexCredentials();

      final json = credentials.toJson();

      expect(json, isEmpty);
    });

    test('toJson excludes only null fields', () {
      final credentials = VtexCredentials(appKey: 'mykey');

      final json = credentials.toJson();

      expect(json.length, 1);
      expect(json['app_key'], 'mykey');
      expect(json.containsKey('app_token'), false);
    });

    test('toJson roundtrip preserves all fields', () {
      final original = VtexCredentials(
        appKey: 'vtexappkey-mystore-ABCDEF',
        appToken: 'GHIJKLMNOPQRSTUVWXYZ',
      );

      final json = original.toJson();
      final restored = VtexCredentials.fromJson(json);

      expect(restored.appKey, original.appKey);
      expect(restored.appToken, original.appToken);
    });

    test('toJson roundtrip with empty credentials', () {
      final original = VtexCredentials();

      final json = original.toJson();
      final restored = VtexCredentials.fromJson(json);

      expect(restored.appKey, isNull);
      expect(restored.appToken, isNull);
    });
  });
}
