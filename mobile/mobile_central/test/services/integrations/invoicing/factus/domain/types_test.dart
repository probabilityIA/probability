import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/invoicing/factus/domain/types.dart';

void main() {
  group('FactusConfig', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'numbering_range_id': 7,
        'default_tax_rate': '19.00',
        'payment_form': '1',
        'payment_method_code': '10',
        'legal_organization_id': '2',
        'tribute_id': '01',
        'identification_document_id': '13',
        'municipality_id': '11001',
      };

      final config = FactusConfig.fromJson(json);

      expect(config.numberingRangeId, 7);
      expect(config.defaultTaxRate, '19.00');
      expect(config.paymentForm, '1');
      expect(config.paymentMethodCode, '10');
      expect(config.legalOrganizationId, '2');
      expect(config.tributeId, '01');
      expect(config.identificationDocumentId, '13');
      expect(config.municipalityId, '11001');
    });

    test('fromJson handles missing fields', () {
      final json = <String, dynamic>{};

      final config = FactusConfig.fromJson(json);

      expect(config.numberingRangeId, isNull);
      expect(config.defaultTaxRate, isNull);
      expect(config.paymentForm, isNull);
      expect(config.paymentMethodCode, isNull);
      expect(config.legalOrganizationId, isNull);
      expect(config.tributeId, isNull);
      expect(config.identificationDocumentId, isNull);
      expect(config.municipalityId, isNull);
    });

    test('fromJson handles explicit null values', () {
      final json = {
        'numbering_range_id': null,
        'default_tax_rate': null,
        'payment_form': null,
        'payment_method_code': null,
        'legal_organization_id': null,
        'tribute_id': null,
        'identification_document_id': null,
        'municipality_id': null,
      };

      final config = FactusConfig.fromJson(json);

      expect(config.numberingRangeId, isNull);
      expect(config.defaultTaxRate, isNull);
      expect(config.paymentForm, isNull);
      expect(config.paymentMethodCode, isNull);
      expect(config.legalOrganizationId, isNull);
      expect(config.tributeId, isNull);
      expect(config.identificationDocumentId, isNull);
      expect(config.municipalityId, isNull);
    });

    test('toJson includes all non-null fields', () {
      final config = FactusConfig(
        numberingRangeId: 7,
        defaultTaxRate: '19.00',
        paymentForm: '1',
        paymentMethodCode: '10',
        legalOrganizationId: '2',
        tributeId: '01',
        identificationDocumentId: '13',
        municipalityId: '11001',
      );

      final json = config.toJson();

      expect(json['numbering_range_id'], 7);
      expect(json['default_tax_rate'], '19.00');
      expect(json['payment_form'], '1');
      expect(json['payment_method_code'], '10');
      expect(json['legal_organization_id'], '2');
      expect(json['tribute_id'], '01');
      expect(json['identification_document_id'], '13');
      expect(json['municipality_id'], '11001');
      expect(json.length, 8);
    });

    test('toJson excludes null fields', () {
      final config = FactusConfig();

      final json = config.toJson();

      expect(json, isEmpty);
    });

    test('toJson excludes only null fields', () {
      final config = FactusConfig(
        numberingRangeId: 3,
        defaultTaxRate: '19.00',
      );

      final json = config.toJson();

      expect(json.length, 2);
      expect(json['numbering_range_id'], 3);
      expect(json['default_tax_rate'], '19.00');
      expect(json.containsKey('payment_form'), false);
      expect(json.containsKey('payment_method_code'), false);
      expect(json.containsKey('legal_organization_id'), false);
      expect(json.containsKey('tribute_id'), false);
      expect(json.containsKey('identification_document_id'), false);
      expect(json.containsKey('municipality_id'), false);
    });

    test('toJson roundtrip preserves all fields', () {
      final original = FactusConfig(
        numberingRangeId: 7,
        defaultTaxRate: '19.00',
        paymentForm: '1',
        paymentMethodCode: '10',
        legalOrganizationId: '2',
        tributeId: '01',
        identificationDocumentId: '13',
        municipalityId: '11001',
      );

      final json = original.toJson();
      final restored = FactusConfig.fromJson(json);

      expect(restored.numberingRangeId, original.numberingRangeId);
      expect(restored.defaultTaxRate, original.defaultTaxRate);
      expect(restored.paymentForm, original.paymentForm);
      expect(restored.paymentMethodCode, original.paymentMethodCode);
      expect(restored.legalOrganizationId, original.legalOrganizationId);
      expect(restored.tributeId, original.tributeId);
      expect(
          restored.identificationDocumentId,
          original.identificationDocumentId);
      expect(restored.municipalityId, original.municipalityId);
    });

    test('toJson roundtrip with empty config', () {
      final original = FactusConfig();

      final json = original.toJson();
      final restored = FactusConfig.fromJson(json);

      expect(restored.numberingRangeId, isNull);
      expect(restored.defaultTaxRate, isNull);
      expect(restored.paymentForm, isNull);
      expect(restored.paymentMethodCode, isNull);
      expect(restored.legalOrganizationId, isNull);
      expect(restored.tributeId, isNull);
      expect(restored.identificationDocumentId, isNull);
      expect(restored.municipalityId, isNull);
    });

    test('fromJson handles numberingRangeId as int', () {
      final json = {'numbering_range_id': 99};

      final config = FactusConfig.fromJson(json);

      expect(config.numberingRangeId, 99);
    });
  });

  group('FactusCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'client_id': 'factus_client_id_abc',
        'client_secret': 'factus_client_secret_xyz',
        'username': 'factus_user',
        'password': 'factus_pass_123',
        'api_url': 'https://api.factus.com.co',
      };

      final credentials = FactusCredentials.fromJson(json);

      expect(credentials.clientId, 'factus_client_id_abc');
      expect(credentials.clientSecret, 'factus_client_secret_xyz');
      expect(credentials.username, 'factus_user');
      expect(credentials.password, 'factus_pass_123');
      expect(credentials.apiUrl, 'https://api.factus.com.co');
    });

    test('fromJson handles missing fields', () {
      final json = <String, dynamic>{};

      final credentials = FactusCredentials.fromJson(json);

      expect(credentials.clientId, isNull);
      expect(credentials.clientSecret, isNull);
      expect(credentials.username, isNull);
      expect(credentials.password, isNull);
      expect(credentials.apiUrl, isNull);
    });

    test('fromJson handles explicit null values', () {
      final json = {
        'client_id': null,
        'client_secret': null,
        'username': null,
        'password': null,
        'api_url': null,
      };

      final credentials = FactusCredentials.fromJson(json);

      expect(credentials.clientId, isNull);
      expect(credentials.clientSecret, isNull);
      expect(credentials.username, isNull);
      expect(credentials.password, isNull);
      expect(credentials.apiUrl, isNull);
    });

    test('toJson includes all non-null fields', () {
      final credentials = FactusCredentials(
        clientId: 'factus_client_id_abc',
        clientSecret: 'factus_client_secret_xyz',
        username: 'factus_user',
        password: 'factus_pass_123',
        apiUrl: 'https://api.factus.com.co',
      );

      final json = credentials.toJson();

      expect(json['client_id'], 'factus_client_id_abc');
      expect(json['client_secret'], 'factus_client_secret_xyz');
      expect(json['username'], 'factus_user');
      expect(json['password'], 'factus_pass_123');
      expect(json['api_url'], 'https://api.factus.com.co');
      expect(json.length, 5);
    });

    test('toJson excludes null fields', () {
      final credentials = FactusCredentials();

      final json = credentials.toJson();

      expect(json, isEmpty);
    });

    test('toJson excludes only null fields', () {
      final credentials = FactusCredentials(
        clientId: 'client_1',
        username: 'user_1',
      );

      final json = credentials.toJson();

      expect(json.length, 2);
      expect(json['client_id'], 'client_1');
      expect(json['username'], 'user_1');
      expect(json.containsKey('client_secret'), false);
      expect(json.containsKey('password'), false);
      expect(json.containsKey('api_url'), false);
    });

    test('toJson roundtrip preserves all fields', () {
      final original = FactusCredentials(
        clientId: 'factus_client_id_abc',
        clientSecret: 'factus_client_secret_xyz',
        username: 'factus_user',
        password: 'factus_pass_123',
        apiUrl: 'https://api.factus.com.co',
      );

      final json = original.toJson();
      final restored = FactusCredentials.fromJson(json);

      expect(restored.clientId, original.clientId);
      expect(restored.clientSecret, original.clientSecret);
      expect(restored.username, original.username);
      expect(restored.password, original.password);
      expect(restored.apiUrl, original.apiUrl);
    });

    test('toJson roundtrip with empty credentials', () {
      final original = FactusCredentials();

      final json = original.toJson();
      final restored = FactusCredentials.fromJson(json);

      expect(restored.clientId, isNull);
      expect(restored.clientSecret, isNull);
      expect(restored.username, isNull);
      expect(restored.password, isNull);
      expect(restored.apiUrl, isNull);
    });
  });
}
