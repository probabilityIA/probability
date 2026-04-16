class StripeConfig {
  StripeConfig();
  factory StripeConfig.fromJson(Map<String, dynamic> json) => StripeConfig();
  Map<String, dynamic> toJson() => <String, dynamic>{};
}

class StripeCredentials {
  final String? secretKey;
  final String? environment; // test | live

  StripeCredentials({this.secretKey, this.environment});

  factory StripeCredentials.fromJson(Map<String, dynamic> json) {
    return StripeCredentials(secretKey: json['secret_key'], environment: json['environment']);
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (secretKey != null) json['secret_key'] = secretKey;
    if (environment != null) json['environment'] = environment;
    return json;
  }
}
