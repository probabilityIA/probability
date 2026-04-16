class EnvioClickConfig {
  final bool? usePlatformToken;
  final String? baseUrlTest;

  EnvioClickConfig({this.usePlatformToken, this.baseUrlTest});

  factory EnvioClickConfig.fromJson(Map<String, dynamic> json) {
    return EnvioClickConfig(
      usePlatformToken: json['use_platform_token'],
      baseUrlTest: json['base_url_test'],
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (usePlatformToken != null) json['use_platform_token'] = usePlatformToken;
    if (baseUrlTest != null) json['base_url_test'] = baseUrlTest;
    return json;
  }
}

class EnvioClickCredentials {
  final String? apiKey;

  EnvioClickCredentials({this.apiKey});

  factory EnvioClickCredentials.fromJson(Map<String, dynamic> json) {
    return EnvioClickCredentials(apiKey: json['api_key']);
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (apiKey != null) json['api_key'] = apiKey;
    return json;
  }
}
