import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class TokenStorage {
  static const _tokenKey = 'session_token';
  static const _userKey = 'user_data';
  static const _sessionKey = 'session_data';
  static const _biometricKey = 'biometric_enabled';

  final FlutterSecureStorage _storage;

  TokenStorage({FlutterSecureStorage? storage})
      : _storage = storage ?? const FlutterSecureStorage();

  Future<void> saveToken(String token) async {
    await _storage.write(key: _tokenKey, value: token);
  }

  Future<String?> getToken() async {
    return _storage.read(key: _tokenKey);
  }

  Future<void> saveUserData(String userData) async {
    await _storage.write(key: _userKey, value: userData);
  }

  Future<String?> getUserData() async {
    return _storage.read(key: _userKey);
  }

  Future<void> saveSessionData(String sessionData) async {
    await _storage.write(key: _sessionKey, value: sessionData);
  }

  Future<String?> getSessionData() async {
    return _storage.read(key: _sessionKey);
  }

  Future<void> setBiometricEnabled(bool enabled) async {
    await _storage.write(key: _biometricKey, value: enabled ? '1' : '0');
  }

  Future<bool> isBiometricEnabled() async {
    return await _storage.read(key: _biometricKey) == '1';
  }

  Future<void> clearAll() async {
    await _storage.deleteAll();
  }
}
