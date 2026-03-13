class ExitoConfig {
  ExitoConfig();

  factory ExitoConfig.fromJson(Map<String, dynamic> json) {
    return ExitoConfig();
  }

  Map<String, dynamic> toJson() => <String, dynamic>{};
}

class ExitoCredentials {
  final String? apiKey;
  final String? sellerId;

  ExitoCredentials({this.apiKey, this.sellerId});

  factory ExitoCredentials.fromJson(Map<String, dynamic> json) {
    return ExitoCredentials(
      apiKey: json['api_key'],
      sellerId: json['seller_id'],
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (apiKey != null) json['api_key'] = apiKey;
    if (sellerId != null) json['seller_id'] = sellerId;
    return json;
  }
}
