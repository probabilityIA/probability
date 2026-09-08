import 'package:flutter/services.dart';
import 'package:local_auth/local_auth.dart';

enum BiometricStatus { available, noHardware, notEnrolled, unavailable }

class BiometricAuth {
  BiometricAuth({LocalAuthentication? auth})
      : _auth = auth ?? LocalAuthentication();

  final LocalAuthentication _auth;

  Future<BiometricStatus> status() async {
    try {
      if (!await _auth.isDeviceSupported()) return BiometricStatus.noHardware;
      if (!await _auth.canCheckBiometrics) return BiometricStatus.noHardware;
      final enrolled = await _auth.getAvailableBiometrics();
      if (enrolled.isEmpty) return BiometricStatus.notEnrolled;
      return BiometricStatus.available;
    } on LocalAuthException {
      return BiometricStatus.unavailable;
    } on PlatformException {
      return BiometricStatus.unavailable;
    }
  }

  Future<bool> isAvailable() async =>
      await status() == BiometricStatus.available;

  Future<bool> authenticate(String reason) async {
    try {
      return await _auth.authenticate(
        localizedReason: reason,
        biometricOnly: false,
        persistAcrossBackgrounding: true,
      );
    } on LocalAuthException {
      return false;
    } on PlatformException {
      return false;
    }
  }
}
