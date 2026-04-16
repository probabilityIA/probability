class WhatsAppConfig {
  final String? whatsappUrl;
  final String? webhookCallbackUrl;

  WhatsAppConfig({this.whatsappUrl, this.webhookCallbackUrl});

  factory WhatsAppConfig.fromJson(Map<String, dynamic> json) {
    return WhatsAppConfig(
      whatsappUrl: json['whatsapp_url'],
      webhookCallbackUrl: json['webhook_callback_url'],
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (whatsappUrl != null) json['whatsapp_url'] = whatsappUrl;
    if (webhookCallbackUrl != null) json['webhook_callback_url'] = webhookCallbackUrl;
    return json;
  }
}

class WhatsAppCredentials {
  final String? phoneNumberId;
  final String? accessToken;
  final String? verifyToken;
  final String? testPhoneNumber;

  WhatsAppCredentials({
    this.phoneNumberId,
    this.accessToken,
    this.verifyToken,
    this.testPhoneNumber,
  });

  factory WhatsAppCredentials.fromJson(Map<String, dynamic> json) {
    return WhatsAppCredentials(
      phoneNumberId: json['phone_number_id'],
      accessToken: json['access_token'],
      verifyToken: json['verify_token'],
      testPhoneNumber: json['test_phone_number'],
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (phoneNumberId != null) json['phone_number_id'] = phoneNumberId;
    if (accessToken != null) json['access_token'] = accessToken;
    if (verifyToken != null) json['verify_token'] = verifyToken;
    if (testPhoneNumber != null) json['test_phone_number'] = testPhoneNumber;
    return json;
  }
}
