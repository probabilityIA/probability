enum AppEnvironment { development, staging, production }

class Environment {
  static const String _env = String.fromEnvironment(
    'APP_ENV',
    defaultValue: 'production',
  );

  static const String _apiBaseUrlOverride = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: '',
  );

  static AppEnvironment get current {
    switch (_env) {
      case 'development':
      case 'dev':
        return AppEnvironment.development;
      case 'staging':
      case 'stg':
        return AppEnvironment.staging;
      case 'production':
      case 'prod':
      default:
        return AppEnvironment.production;
    }
  }

  static String get apiBaseUrl {
    if (_apiBaseUrlOverride.isNotEmpty) return _apiBaseUrlOverride;

    switch (current) {
      case AppEnvironment.development:
        return 'http://10.0.2.2:3050/api/v1';
      case AppEnvironment.staging:
        return 'https://staging.probabilityia.com.co/api/v1';
      case AppEnvironment.production:
        return 'https://www.probabilityia.com.co/api/v1';
    }
  }

  static bool get isDevelopment => current == AppEnvironment.development;
  static bool get isStaging => current == AppEnvironment.staging;
  static bool get isProduction => current == AppEnvironment.production;

  /// Credenciales de desarrollo (solo disponibles con --dart-define)
  static const String devEmail = String.fromEnvironment('DEV_EMAIL', defaultValue: '');
  static const String devPassword = String.fromEnvironment('DEV_PASSWORD', defaultValue: '');
  static bool get hasDevCredentials => devEmail.isNotEmpty && devPassword.isNotEmpty;
}
