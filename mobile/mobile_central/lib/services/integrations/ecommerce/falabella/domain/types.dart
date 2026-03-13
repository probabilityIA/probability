class FalabellaConfig {
  FalabellaConfig();

  factory FalabellaConfig.fromJson(Map<String, dynamic> json) {
    return FalabellaConfig();
  }

  Map<String, dynamic> toJson() => <String, dynamic>{};
}

class FalabellaCredentials {
  final String? apiKey;
  final String? userId;

  FalabellaCredentials({this.apiKey, this.userId});

  factory FalabellaCredentials.fromJson(Map<String, dynamic> json) {
    return FalabellaCredentials(
      apiKey: json['api_key'],
      userId: json['user_id'],
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (apiKey != null) json['api_key'] = apiKey;
    if (userId != null) json['user_id'] = userId;
    return json;
  }
}
