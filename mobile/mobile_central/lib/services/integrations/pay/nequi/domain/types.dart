class NequiConfig {
  final String? phoneCode;

  NequiConfig({this.phoneCode});

  factory NequiConfig.fromJson(Map<String, dynamic> json) {
    return NequiConfig(phoneCode: json['phone_code']);
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (phoneCode != null) json['phone_code'] = phoneCode;
    return json;
  }
}

class NequiCredentials {
  final String? apiKey;
  final String? environment; // sandbox | production

  NequiCredentials({this.apiKey, this.environment});

  factory NequiCredentials.fromJson(Map<String, dynamic> json) {
    return NequiCredentials(apiKey: json['api_key'], environment: json['environment']);
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (apiKey != null) json['api_key'] = apiKey;
    if (environment != null) json['environment'] = environment;
    return json;
  }
}
