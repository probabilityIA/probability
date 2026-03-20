import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/messages/whatsapp/domain/types.dart';

void main() {
  group('WhatsAppConfig', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'whatsapp_url': 'https://graph.facebook.com/v18.0',
        'webhook_callback_url': 'https://api.example.com/webhooks/whatsapp',
      };

      final config = WhatsAppConfig.fromJson(json);

      expect(config.whatsappUrl, 'https://graph.facebook.com/v18.0');
      expect(config.webhookCallbackUrl, 'https://api.example.com/webhooks/whatsapp');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};

      final config = WhatsAppConfig.fromJson(json);

      expect(config.whatsappUrl, isNull);
      expect(config.webhookCallbackUrl, isNull);
    });

    test('fromJson handles partial fields - only whatsappUrl', () {
      final json = {'whatsapp_url': 'https://graph.facebook.com/v18.0'};

      final config = WhatsAppConfig.fromJson(json);

      expect(config.whatsappUrl, 'https://graph.facebook.com/v18.0');
      expect(config.webhookCallbackUrl, isNull);
    });

    test('fromJson handles partial fields - only webhookCallbackUrl', () {
      final json = {'webhook_callback_url': 'https://api.example.com/webhook'};

      final config = WhatsAppConfig.fromJson(json);

      expect(config.whatsappUrl, isNull);
      expect(config.webhookCallbackUrl, 'https://api.example.com/webhook');
    });

    test('fromJson ignores unknown fields', () {
      final json = {
        'whatsapp_url': 'https://graph.facebook.com/v18.0',
        'webhook_callback_url': 'https://api.example.com/webhook',
        'extra_field': 'ignored',
      };

      final config = WhatsAppConfig.fromJson(json);

      expect(config.whatsappUrl, 'https://graph.facebook.com/v18.0');
      expect(config.webhookCallbackUrl, 'https://api.example.com/webhook');
    });

    test('toJson includes all non-null fields', () {
      final config = WhatsAppConfig(
        whatsappUrl: 'https://graph.facebook.com/v18.0',
        webhookCallbackUrl: 'https://api.example.com/webhook',
      );

      final json = config.toJson();

      expect(json['whatsapp_url'], 'https://graph.facebook.com/v18.0');
      expect(json['webhook_callback_url'], 'https://api.example.com/webhook');
      expect(json.length, 2);
    });

    test('toJson omits null fields', () {
      final config = WhatsAppConfig();

      final json = config.toJson();

      expect(json, isEmpty);
    });

    test('toJson omits only null whatsappUrl', () {
      final config = WhatsAppConfig(webhookCallbackUrl: 'https://api.example.com/webhook');

      final json = config.toJson();

      expect(json.containsKey('whatsapp_url'), isFalse);
      expect(json['webhook_callback_url'], 'https://api.example.com/webhook');
      expect(json.length, 1);
    });

    test('toJson omits only null webhookCallbackUrl', () {
      final config = WhatsAppConfig(whatsappUrl: 'https://graph.facebook.com/v18.0');

      final json = config.toJson();

      expect(json['whatsapp_url'], 'https://graph.facebook.com/v18.0');
      expect(json.containsKey('webhook_callback_url'), isFalse);
      expect(json.length, 1);
    });

    test('fromJson/toJson roundtrip with all fields', () {
      final original = WhatsAppConfig(
        whatsappUrl: 'https://graph.facebook.com/v18.0',
        webhookCallbackUrl: 'https://api.example.com/webhook',
      );

      final json = original.toJson();
      final restored = WhatsAppConfig.fromJson(json);

      expect(restored.whatsappUrl, original.whatsappUrl);
      expect(restored.webhookCallbackUrl, original.webhookCallbackUrl);
    });

    test('fromJson/toJson roundtrip with empty config', () {
      final original = WhatsAppConfig();

      final json = original.toJson();
      final restored = WhatsAppConfig.fromJson(json);

      expect(restored.whatsappUrl, isNull);
      expect(restored.webhookCallbackUrl, isNull);
    });

    test('default constructor allows all nulls', () {
      final config = WhatsAppConfig();

      expect(config.whatsappUrl, isNull);
      expect(config.webhookCallbackUrl, isNull);
    });
  });

  group('WhatsAppCredentials', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'phone_number_id': '123456789012345',
        'access_token': 'EAABsbCS1iH0BAJ...',
        'verify_token': 'my_verify_token_123',
        'test_phone_number': '+573001234567',
      };

      final creds = WhatsAppCredentials.fromJson(json);

      expect(creds.phoneNumberId, '123456789012345');
      expect(creds.accessToken, 'EAABsbCS1iH0BAJ...');
      expect(creds.verifyToken, 'my_verify_token_123');
      expect(creds.testPhoneNumber, '+573001234567');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};

      final creds = WhatsAppCredentials.fromJson(json);

      expect(creds.phoneNumberId, isNull);
      expect(creds.accessToken, isNull);
      expect(creds.verifyToken, isNull);
      expect(creds.testPhoneNumber, isNull);
    });

    test('fromJson handles partial fields - only phoneNumberId', () {
      final json = {'phone_number_id': '123456789'};

      final creds = WhatsAppCredentials.fromJson(json);

      expect(creds.phoneNumberId, '123456789');
      expect(creds.accessToken, isNull);
      expect(creds.verifyToken, isNull);
      expect(creds.testPhoneNumber, isNull);
    });

    test('fromJson handles partial fields - only accessToken', () {
      final json = {'access_token': 'token_value'};

      final creds = WhatsAppCredentials.fromJson(json);

      expect(creds.phoneNumberId, isNull);
      expect(creds.accessToken, 'token_value');
      expect(creds.verifyToken, isNull);
      expect(creds.testPhoneNumber, isNull);
    });

    test('fromJson handles partial fields - only verifyToken', () {
      final json = {'verify_token': 'verify_abc'};

      final creds = WhatsAppCredentials.fromJson(json);

      expect(creds.phoneNumberId, isNull);
      expect(creds.accessToken, isNull);
      expect(creds.verifyToken, 'verify_abc');
      expect(creds.testPhoneNumber, isNull);
    });

    test('fromJson handles partial fields - only testPhoneNumber', () {
      final json = {'test_phone_number': '+573009876543'};

      final creds = WhatsAppCredentials.fromJson(json);

      expect(creds.phoneNumberId, isNull);
      expect(creds.accessToken, isNull);
      expect(creds.verifyToken, isNull);
      expect(creds.testPhoneNumber, '+573009876543');
    });

    test('toJson includes all non-null fields', () {
      final creds = WhatsAppCredentials(
        phoneNumberId: '123456789',
        accessToken: 'token_abc',
        verifyToken: 'verify_xyz',
        testPhoneNumber: '+573001112233',
      );

      final json = creds.toJson();

      expect(json['phone_number_id'], '123456789');
      expect(json['access_token'], 'token_abc');
      expect(json['verify_token'], 'verify_xyz');
      expect(json['test_phone_number'], '+573001112233');
      expect(json.length, 4);
    });

    test('toJson omits null fields', () {
      final creds = WhatsAppCredentials();

      final json = creds.toJson();

      expect(json, isEmpty);
    });

    test('toJson omits selective null fields', () {
      final creds = WhatsAppCredentials(
        phoneNumberId: '123',
        accessToken: 'token',
      );

      final json = creds.toJson();

      expect(json['phone_number_id'], '123');
      expect(json['access_token'], 'token');
      expect(json.containsKey('verify_token'), isFalse);
      expect(json.containsKey('test_phone_number'), isFalse);
      expect(json.length, 2);
    });

    test('fromJson/toJson roundtrip with all fields', () {
      final original = WhatsAppCredentials(
        phoneNumberId: '123456789',
        accessToken: 'token_rt',
        verifyToken: 'verify_rt',
        testPhoneNumber: '+573001234567',
      );

      final json = original.toJson();
      final restored = WhatsAppCredentials.fromJson(json);

      expect(restored.phoneNumberId, original.phoneNumberId);
      expect(restored.accessToken, original.accessToken);
      expect(restored.verifyToken, original.verifyToken);
      expect(restored.testPhoneNumber, original.testPhoneNumber);
    });

    test('fromJson/toJson roundtrip with empty credentials', () {
      final original = WhatsAppCredentials();

      final json = original.toJson();
      final restored = WhatsAppCredentials.fromJson(json);

      expect(restored.phoneNumberId, isNull);
      expect(restored.accessToken, isNull);
      expect(restored.verifyToken, isNull);
      expect(restored.testPhoneNumber, isNull);
    });

    test('default constructor allows all nulls', () {
      final creds = WhatsAppCredentials();

      expect(creds.phoneNumberId, isNull);
      expect(creds.accessToken, isNull);
      expect(creds.verifyToken, isNull);
      expect(creds.testPhoneNumber, isNull);
    });
  });
}
