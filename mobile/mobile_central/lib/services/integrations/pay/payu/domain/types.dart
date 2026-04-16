class PayUConfig {
  final String? accountId;
  final String? merchantId;

  PayUConfig({this.accountId, this.merchantId});

  factory PayUConfig.fromJson(Map<String, dynamic> json) {
    return PayUConfig(accountId: json['account_id'], merchantId: json['merchant_id']);
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (accountId != null) json['account_id'] = accountId;
    if (merchantId != null) json['merchant_id'] = merchantId;
    return json;
  }
}

class PayUCredentials {
  final String? apiKey;
  final String? apiLogin;
  final String? environment; // sandbox | production

  PayUCredentials({this.apiKey, this.apiLogin, this.environment});

  factory PayUCredentials.fromJson(Map<String, dynamic> json) {
    return PayUCredentials(
      apiKey: json['api_key'],
      apiLogin: json['api_login'],
      environment: json['environment'],
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (apiKey != null) json['api_key'] = apiKey;
    if (apiLogin != null) json['api_login'] = apiLogin;
    if (environment != null) json['environment'] = environment;
    return json;
  }
}
