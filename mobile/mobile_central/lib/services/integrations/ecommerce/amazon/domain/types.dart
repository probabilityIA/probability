class AmazonConfig {
  final String? marketplaceId;
  final String? region;

  AmazonConfig({
    this.marketplaceId,
    this.region,
  });

  factory AmazonConfig.fromJson(Map<String, dynamic> json) {
    return AmazonConfig(
      marketplaceId: json['marketplace_id'],
      region: json['region'],
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (marketplaceId != null) json['marketplace_id'] = marketplaceId;
    if (region != null) json['region'] = region;
    return json;
  }
}

class AmazonCredentials {
  final String? sellerId;
  final String? refreshToken;

  AmazonCredentials({
    this.sellerId,
    this.refreshToken,
  });

  factory AmazonCredentials.fromJson(Map<String, dynamic> json) {
    return AmazonCredentials(
      sellerId: json['seller_id'],
      refreshToken: json['refresh_token'],
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (sellerId != null) json['seller_id'] = sellerId;
    if (refreshToken != null) json['refresh_token'] = refreshToken;
    return json;
  }
}
