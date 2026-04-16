class TuConfig {
  final String? baseUrl;

  TuConfig({this.baseUrl});

  factory TuConfig.fromJson(Map<String, dynamic> json) {
    return TuConfig(baseUrl: json['base_url']);
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (baseUrl != null) json['base_url'] = baseUrl;
    return json;
  }
}

class TuCredentials {
  final String? apiKey;

  TuCredentials({this.apiKey});

  factory TuCredentials.fromJson(Map<String, dynamic> json) {
    return TuCredentials(apiKey: json['api_key']);
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (apiKey != null) json['api_key'] = apiKey;
    return json;
  }
}
