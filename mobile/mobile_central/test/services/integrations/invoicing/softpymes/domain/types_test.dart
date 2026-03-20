import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/invoicing/softpymes/domain/types.dart';

void main() {
  group('SoftpymesConfig', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'company_nit': '900123456-7',
        'company_name': 'Mi Empresa SAS',
        'referer': '900123456',
        'default_customer_nit': '222222222',
        'resolution_id': 42,
        'branch_code': '001',
        'customer_branch_code': '000',
        'seller_nit': '111111111',
      };

      final config = SoftpymesConfig.fromJson(json);

      expect(config.companyNit, '900123456-7');
      expect(config.companyName, 'Mi Empresa SAS');
      expect(config.referer, '900123456');
      expect(config.defaultCustomerNit, '222222222');
      expect(config.resolutionId, 42);
      expect(config.branchCode, '001');
      expect(config.customerBranchCode, '000');
      expect(config.sellerNit, '111111111');
    });

    test('fromJson handles missing fields', () {
      final json = <String, dynamic>{};

      final config = SoftpymesConfig.fromJson(json);

      expect(config.companyNit, isNull);
      expect(config.companyName, isNull);
      expect(config.referer, isNull);
      expect(config.defaultCustomerNit, isNull);
      expect(config.resolutionId, isNull);
      expect(config.branchCode, isNull);
      expect(config.customerBranchCode, isNull);
      expect(config.sellerNit, isNull);
    });

    test('fromJson handles explicit null values', () {
      final json = {
        'company_nit': null,
        'company_name': null,
        'referer': null,
        'default_customer_nit': null,
        'resolution_id': null,
        'branch_code': null,
        'customer_branch_code': null,
        'seller_nit': null,
      };

      final config = SoftpymesConfig.fromJson(json);

      expect(config.companyNit, isNull);
      expect(config.companyName, isNull);
      expect(config.referer, isNull);
      expect(config.defaultCustomerNit, isNull);
      expect(config.resolutionId, isNull);
      expect(config.branchCode, isNull);
      expect(config.customerBranchCode, isNull);
      expect(config.sellerNit, isNull);
    });

    test('toJson includes all non-null fields', () {
      final config = SoftpymesConfig(
        companyNit: '900123456-7',
        companyName: 'Mi Empresa SAS',
        referer: '900123456',
        defaultCustomerNit: '222222222',
        resolutionId: 42,
        branchCode: '001',
        customerBranchCode: '000',
        sellerNit: '111111111',
      );

      final json = config.toJson();

      expect(json['company_nit'], '900123456-7');
      expect(json['company_name'], 'Mi Empresa SAS');
      expect(json['referer'], '900123456');
      expect(json['default_customer_nit'], '222222222');
      expect(json['resolution_id'], 42);
      expect(json['branch_code'], '001');
      expect(json['customer_branch_code'], '000');
      expect(json['seller_nit'], '111111111');
      expect(json.length, 8);
    });

    test('toJson excludes null fields', () {
      final config = SoftpymesConfig();

      final json = config.toJson();

      expect(json, isEmpty);
    });

    test('toJson excludes only null fields', () {
      final config = SoftpymesConfig(
        companyNit: '900123456-7',
        resolutionId: 42,
      );

      final json = config.toJson();

      expect(json.length, 2);
      expect(json['company_nit'], '900123456-7');
      expect(json['resolution_id'], 42);
      expect(json.containsKey('company_name'), false);
      expect(json.containsKey('referer'), false);
      expect(json.containsKey('default_customer_nit'), false);
      expect(json.containsKey('branch_code'), false);
      expect(json.containsKey('customer_branch_code'), false);
      expect(json.containsKey('seller_nit'), false);
    });

    test('toJson roundtrip preserves all fields', () {
      final original = SoftpymesConfig(
        companyNit: '900123456-7',
        companyName: 'Mi Empresa SAS',
        referer: '900123456',
        defaultCustomerNit: '222222222',
        resolutionId: 42,
        branchCode: '001',
        customerBranchCode: '000',
        sellerNit: '111111111',
      );

      final json = original.toJson();
      final restored = SoftpymesConfig.fromJson(json);

      expect(restored.companyNit, original.companyNit);
      expect(restored.companyName, original.companyName);
      expect(restored.referer, original.referer);
      expect(restored.defaultCustomerNit, original.defaultCustomerNit);
      expect(restored.resolutionId, original.resolutionId);
      expect(restored.branchCode, original.branchCode);
      expect(restored.customerBranchCode, original.customerBranchCode);
      expect(restored.sellerNit, original.sellerNit);
    });

    test('toJson roundtrip with empty config', () {
      final original = SoftpymesConfig();

      final json = original.toJson();
      final restored = SoftpymesConfig.fromJson(json);

      expect(restored.companyNit, isNull);
      expect(restored.companyName, isNull);
      expect(restored.referer, isNull);
      expect(restored.defaultCustomerNit, isNull);
      expect(restored.resolutionId, isNull);
      expect(restored.branchCode, isNull);
      expect(restored.customerBranchCode, isNull);
      expect(restored.sellerNit, isNull);
    });

    test('fromJson handles resolutionId as int', () {
      final json = {'resolution_id': 100};

      final config = SoftpymesConfig.fromJson(json);

      expect(config.resolutionId, 100);
    });
  });

  group('SoftpymesCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'api_key': 'sp_api_key_123',
        'api_secret': 'sp_api_secret_456',
      };

      final credentials = SoftpymesCredentials.fromJson(json);

      expect(credentials.apiKey, 'sp_api_key_123');
      expect(credentials.apiSecret, 'sp_api_secret_456');
    });

    test('fromJson handles missing fields', () {
      final json = <String, dynamic>{};

      final credentials = SoftpymesCredentials.fromJson(json);

      expect(credentials.apiKey, isNull);
      expect(credentials.apiSecret, isNull);
    });

    test('fromJson handles explicit null values', () {
      final json = {
        'api_key': null,
        'api_secret': null,
      };

      final credentials = SoftpymesCredentials.fromJson(json);

      expect(credentials.apiKey, isNull);
      expect(credentials.apiSecret, isNull);
    });

    test('toJson includes all non-null fields', () {
      final credentials = SoftpymesCredentials(
        apiKey: 'sp_api_key_123',
        apiSecret: 'sp_api_secret_456',
      );

      final json = credentials.toJson();

      expect(json['api_key'], 'sp_api_key_123');
      expect(json['api_secret'], 'sp_api_secret_456');
      expect(json.length, 2);
    });

    test('toJson excludes null fields', () {
      final credentials = SoftpymesCredentials();

      final json = credentials.toJson();

      expect(json, isEmpty);
    });

    test('toJson excludes only null fields', () {
      final credentials = SoftpymesCredentials(apiKey: 'my_key');

      final json = credentials.toJson();

      expect(json.length, 1);
      expect(json['api_key'], 'my_key');
      expect(json.containsKey('api_secret'), false);
    });

    test('toJson roundtrip preserves all fields', () {
      final original = SoftpymesCredentials(
        apiKey: 'sp_api_key_123',
        apiSecret: 'sp_api_secret_456',
      );

      final json = original.toJson();
      final restored = SoftpymesCredentials.fromJson(json);

      expect(restored.apiKey, original.apiKey);
      expect(restored.apiSecret, original.apiSecret);
    });

    test('toJson roundtrip with empty credentials', () {
      final original = SoftpymesCredentials();

      final json = original.toJson();
      final restored = SoftpymesCredentials.fromJson(json);

      expect(restored.apiKey, isNull);
      expect(restored.apiSecret, isNull);
    });
  });
}
