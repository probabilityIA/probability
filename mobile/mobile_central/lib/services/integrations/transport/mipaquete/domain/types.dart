class MiPaqueteConfig {
  final String? baseUrl;

  MiPaqueteConfig({this.baseUrl});

  factory MiPaqueteConfig.fromJson(Map<String, dynamic> json) {
    return MiPaqueteConfig(baseUrl: json['base_url']);
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (baseUrl != null) json['base_url'] = baseUrl;
    return json;
  }
}

class MiPaqueteCredentials {
  final String? apiKey;

  MiPaqueteCredentials({this.apiKey});

  factory MiPaqueteCredentials.fromJson(Map<String, dynamic> json) {
    return MiPaqueteCredentials(apiKey: json['api_key']);
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (apiKey != null) json['api_key'] = apiKey;
    return json;
  }
}
