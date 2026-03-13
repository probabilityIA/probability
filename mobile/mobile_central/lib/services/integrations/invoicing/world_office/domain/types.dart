class WorldOfficeConfig {
  WorldOfficeConfig();

  factory WorldOfficeConfig.fromJson(Map<String, dynamic> json) {
    return WorldOfficeConfig();
  }

  Map<String, dynamic> toJson() {
    return <String, dynamic>{};
  }
}

class WorldOfficeCredentials {
  final String? username;
  final String? password;
  final String? companyCode;
  final String? baseUrl;

  WorldOfficeCredentials({
    this.username,
    this.password,
    this.companyCode,
    this.baseUrl,
  });

  factory WorldOfficeCredentials.fromJson(Map<String, dynamic> json) {
    return WorldOfficeCredentials(
      username: json['username'],
      password: json['password'],
      companyCode: json['company_code'],
      baseUrl: json['base_url'],
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (username != null) json['username'] = username;
    if (password != null) json['password'] = password;
    if (companyCode != null) json['company_code'] = companyCode;
    if (baseUrl != null) json['base_url'] = baseUrl;
    return json;
  }
}
