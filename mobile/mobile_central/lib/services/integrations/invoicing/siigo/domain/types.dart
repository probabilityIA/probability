class SiigoConfig {
  SiigoConfig();

  factory SiigoConfig.fromJson(Map<String, dynamic> json) {
    return SiigoConfig();
  }

  Map<String, dynamic> toJson() {
    return <String, dynamic>{};
  }
}

class SiigoCredentials {
  final String? username;
  final String? accessKey;
  final String? accountId;
  final String? partnerId;
  final String? baseUrl;

  SiigoCredentials({
    this.username,
    this.accessKey,
    this.accountId,
    this.partnerId,
    this.baseUrl,
  });

  factory SiigoCredentials.fromJson(Map<String, dynamic> json) {
    return SiigoCredentials(
      username: json['username'],
      accessKey: json['access_key'],
      accountId: json['account_id'],
      partnerId: json['partner_id'],
      baseUrl: json['base_url'],
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (username != null) json['username'] = username;
    if (accessKey != null) json['access_key'] = accessKey;
    if (accountId != null) json['account_id'] = accountId;
    if (partnerId != null) json['partner_id'] = partnerId;
    if (baseUrl != null) json['base_url'] = baseUrl;
    return json;
  }
}
