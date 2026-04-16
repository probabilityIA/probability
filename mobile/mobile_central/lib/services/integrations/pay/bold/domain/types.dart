class BoldConfig {
  BoldConfig();
  factory BoldConfig.fromJson(Map<String, dynamic> json) => BoldConfig();
  Map<String, dynamic> toJson() => <String, dynamic>{};
}

class BoldCredentials {
  final String? apiKey;
  final String? environment; // sandbox | production

  BoldCredentials({this.apiKey, this.environment});

  factory BoldCredentials.fromJson(Map<String, dynamic> json) {
    return BoldCredentials(apiKey: json['api_key'], environment: json['environment']);
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (apiKey != null) json['api_key'] = apiKey;
    if (environment != null) json['environment'] = environment;
    return json;
  }
}
