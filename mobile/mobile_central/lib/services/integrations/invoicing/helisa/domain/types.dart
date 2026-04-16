class HelisaConfig {
  HelisaConfig();

  factory HelisaConfig.fromJson(Map<String, dynamic> json) {
    return HelisaConfig();
  }

  Map<String, dynamic> toJson() {
    return <String, dynamic>{};
  }
}

class HelisaCredentials {
  final String? username;
  final String? password;
  final String? companyId;
  final String? baseUrl;

  HelisaCredentials({
    this.username,
    this.password,
    this.companyId,
    this.baseUrl,
  });

  factory HelisaCredentials.fromJson(Map<String, dynamic> json) {
    return HelisaCredentials(
      username: json['username'],
      password: json['password'],
      companyId: json['company_id'],
      baseUrl: json['base_url'],
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (username != null) json['username'] = username;
    if (password != null) json['password'] = password;
    if (companyId != null) json['company_id'] = companyId;
    if (baseUrl != null) json['base_url'] = baseUrl;
    return json;
  }
}
