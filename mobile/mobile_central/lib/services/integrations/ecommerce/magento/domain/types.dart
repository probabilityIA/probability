class MagentoConfig {
  final String? storeUrl;

  MagentoConfig({this.storeUrl});

  factory MagentoConfig.fromJson(Map<String, dynamic> json) {
    return MagentoConfig(storeUrl: json['store_url']);
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (storeUrl != null) json['store_url'] = storeUrl;
    return json;
  }
}

class MagentoCredentials {
  final String? accessToken;

  MagentoCredentials({this.accessToken});

  factory MagentoCredentials.fromJson(Map<String, dynamic> json) {
    return MagentoCredentials(accessToken: json['access_token']);
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (accessToken != null) json['access_token'] = accessToken;
    return json;
  }
}
