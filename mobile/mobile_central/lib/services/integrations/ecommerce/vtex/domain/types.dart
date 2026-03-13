class VtexConfig {
  final String? accountName;
  final String? environment;

  VtexConfig({
    this.accountName,
    this.environment,
  });

  factory VtexConfig.fromJson(Map<String, dynamic> json) {
    return VtexConfig(
      accountName: json['account_name'],
      environment: json['environment'],
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (accountName != null) json['account_name'] = accountName;
    if (environment != null) json['environment'] = environment;
    return json;
  }
}

class VtexCredentials {
  final String? appKey;
  final String? appToken;

  VtexCredentials({
    this.appKey,
    this.appToken,
  });

  factory VtexCredentials.fromJson(Map<String, dynamic> json) {
    return VtexCredentials(
      appKey: json['app_key'],
      appToken: json['app_token'],
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (appKey != null) json['app_key'] = appKey;
    if (appToken != null) json['app_token'] = appToken;
    return json;
  }
}
