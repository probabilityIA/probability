class WooCommerceConfig {
  final String? storeUrl;

  WooCommerceConfig({this.storeUrl});

  factory WooCommerceConfig.fromJson(Map<String, dynamic> json) {
    return WooCommerceConfig(storeUrl: json['store_url']);
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (storeUrl != null) json['store_url'] = storeUrl;
    return json;
  }
}

class WooCommerceCredentials {
  final String? consumerKey;
  final String? consumerSecret;

  WooCommerceCredentials({this.consumerKey, this.consumerSecret});

  factory WooCommerceCredentials.fromJson(Map<String, dynamic> json) {
    return WooCommerceCredentials(
      consumerKey: json['consumer_key'],
      consumerSecret: json['consumer_secret'],
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (consumerKey != null) json['consumer_key'] = consumerKey;
    if (consumerSecret != null) json['consumer_secret'] = consumerSecret;
    return json;
  }
}
