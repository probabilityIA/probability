class EnviameConfig {
  final String? baseUrl;

  EnviameConfig({this.baseUrl});

  factory EnviameConfig.fromJson(Map<String, dynamic> json) {
    return EnviameConfig(baseUrl: json['base_url']);
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (baseUrl != null) json['base_url'] = baseUrl;
    return json;
  }
}

class EnviameCredentials {
  final String? apiKey;

  EnviameCredentials({this.apiKey});

  factory EnviameCredentials.fromJson(Map<String, dynamic> json) {
    return EnviameCredentials(apiKey: json['api_key']);
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (apiKey != null) json['api_key'] = apiKey;
    return json;
  }
}
