class AlegraConfig {
  AlegraConfig();

  factory AlegraConfig.fromJson(Map<String, dynamic> json) {
    return AlegraConfig();
  }

  Map<String, dynamic> toJson() {
    return <String, dynamic>{};
  }
}

class AlegraCredentials {
  final String? email;
  final String? token;
  final String? baseUrl;

  AlegraCredentials({
    this.email,
    this.token,
    this.baseUrl,
  });

  factory AlegraCredentials.fromJson(Map<String, dynamic> json) {
    return AlegraCredentials(
      email: json['email'],
      token: json['token'],
      baseUrl: json['base_url'],
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (email != null) json['email'] = email;
    if (token != null) json['token'] = token;
    if (baseUrl != null) json['base_url'] = baseUrl;
    return json;
  }
}
