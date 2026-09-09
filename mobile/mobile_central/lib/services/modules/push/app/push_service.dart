import 'dart:io';

import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/foundation.dart';

import '../domain/entities.dart';
import '../domain/ports.dart';

typedef PushTapHandler = void Function(PushPayload payload);

class PushService {
  final IPushRepository _repository;

  PushService(this._repository);

  String? _token;
  bool _listening = false;
  PushTapHandler? onTap;

  String? get token => _token;

  Future<bool> requestPermission() async {
    final settings = await FirebaseMessaging.instance.requestPermission(
      alert: true,
      badge: true,
      sound: true,
    );
    return settings.authorizationStatus == AuthorizationStatus.authorized ||
        settings.authorizationStatus == AuthorizationStatus.provisional;
  }

  Future<void> registerAfterLogin({int? businessId, String? appVersion}) async {
    try {
      final granted = await requestPermission();
      if (!granted) return;

      final messaging = FirebaseMessaging.instance;

      if (Platform.isIOS) {
        final apns = await messaging.getAPNSToken();
        if (apns == null) return;
      }

      final token = await messaging.getToken();
      if (token == null || token.isEmpty) return;

      _token = token;
      await _sendToBackend(token, businessId, appVersion);

      messaging.onTokenRefresh.listen((refreshed) async {
        _token = refreshed;
        await _sendToBackend(refreshed, businessId, appVersion);
      });

      _startListening();
    } catch (e) {
      debugPrint('push: no se pudo registrar el dispositivo: $e');
    }
  }

  Future<void> unregister() async {
    final token = _token;
    if (token == null) return;
    try {
      await _repository.unregisterDevice(token);
      await FirebaseMessaging.instance.deleteToken();
    } catch (e) {
      debugPrint('push: no se pudo dar de baja el dispositivo: $e');
    } finally {
      _token = null;
    }
  }

  Future<void> _sendToBackend(
    String token,
    int? businessId,
    String? appVersion,
  ) async {
    await _repository.registerDevice(DeviceRegistration(
      token: token,
      platform: Platform.isIOS ? 'ios' : 'android',
      appVersion: appVersion,
      deviceName: Platform.operatingSystemVersion,
      businessId: businessId,
    ));
  }

  void _startListening() {
    if (_listening) return;
    _listening = true;

    FirebaseMessaging.onMessageOpenedApp.listen((message) {
      onTap?.call(_toPayload(message));
    });

    FirebaseMessaging.instance.getInitialMessage().then((message) {
      if (message != null) {
        onTap?.call(_toPayload(message));
      }
    });
  }

  PushPayload _toPayload(RemoteMessage message) {
    return PushPayload(
      title: message.notification?.title,
      body: message.notification?.body,
      data: message.data.map((k, v) => MapEntry(k, v?.toString() ?? '')),
    );
  }
}
