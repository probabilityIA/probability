class MeliPagoConfig {
  MeliPagoConfig();
  factory MeliPagoConfig.fromJson(Map<String, dynamic> json) => MeliPagoConfig();
  Map<String, dynamic> toJson() => <String, dynamic>{};
}

class MeliPagoCredentials {
  final String? accessToken;
  final String? environment; // sandbox | production

  MeliPagoCredentials({this.accessToken, this.environment});

  factory MeliPagoCredentials.fromJson(Map<String, dynamic> json) {
    return MeliPagoCredentials(
      accessToken: json['access_token'],
      environment: json['environment'],
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (accessToken != null) json['access_token'] = accessToken;
    if (environment != null) json['environment'] = environment;
    return json;
  }
}
