class EPaycoConfig {
  EPaycoConfig();
  factory EPaycoConfig.fromJson(Map<String, dynamic> json) => EPaycoConfig();
  Map<String, dynamic> toJson() => <String, dynamic>{};
}

class EPaycoCredentials {
  final String? customerId;
  final String? key;
  final String? environment; // test | production

  EPaycoCredentials({this.customerId, this.key, this.environment});

  factory EPaycoCredentials.fromJson(Map<String, dynamic> json) {
    return EPaycoCredentials(
      customerId: json['customer_id'],
      key: json['key'],
      environment: json['environment'],
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (customerId != null) json['customer_id'] = customerId;
    if (key != null) json['key'] = key;
    if (environment != null) json['environment'] = environment;
    return json;
  }
}
