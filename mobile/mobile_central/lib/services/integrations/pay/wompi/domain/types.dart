class WompiConfig {
  WompiConfig();
  factory WompiConfig.fromJson(Map<String, dynamic> json) => WompiConfig();
  Map<String, dynamic> toJson() => <String, dynamic>{};
}

class WompiCredentials {
  final String? privateKey;
  final String? environment; // sandbox | production

  WompiCredentials({this.privateKey, this.environment});

  factory WompiCredentials.fromJson(Map<String, dynamic> json) {
    return WompiCredentials(privateKey: json['private_key'], environment: json['environment']);
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (privateKey != null) json['private_key'] = privateKey;
    if (environment != null) json['environment'] = environment;
    return json;
  }
}
