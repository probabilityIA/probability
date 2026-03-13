class TiendanubeConfig {
  final String? storeId;

  TiendanubeConfig({this.storeId});

  factory TiendanubeConfig.fromJson(Map<String, dynamic> json) {
    return TiendanubeConfig(storeId: json['store_id']);
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (storeId != null) json['store_id'] = storeId;
    return json;
  }
}

class TiendanubeCredentials {
  final String? accessToken;

  TiendanubeCredentials({this.accessToken});

  factory TiendanubeCredentials.fromJson(Map<String, dynamic> json) {
    return TiendanubeCredentials(accessToken: json['access_token']);
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (accessToken != null) json['access_token'] = accessToken;
    return json;
  }
}
