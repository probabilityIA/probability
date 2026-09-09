import 'package:firebase_core/firebase_core.dart' show FirebaseOptions;
import 'package:flutter/foundation.dart'
    show TargetPlatform, defaultTargetPlatform, kIsWeb;

class DefaultFirebaseOptions {
  static FirebaseOptions? get currentPlatform {
    if (kIsWeb) return null;
    switch (defaultTargetPlatform) {
      case TargetPlatform.android:
        return android;
      case TargetPlatform.iOS:
        return ios;
      default:
        return null;
    }
  }

  static const FirebaseOptions android = FirebaseOptions(
    apiKey: 'AIzaSyDKgJr8slND26i3EMKz44EWRYGUdN7Uufc',
    appId: '1:176535439905:android:7910eb3f5acb3acc12f160',
    messagingSenderId: '176535439905',
    projectId: 'probability-app-6de7c',
    storageBucket: 'probability-app-6de7c.firebasestorage.app',
  );

  static const FirebaseOptions ios = FirebaseOptions(
    apiKey: 'AIzaSyA-ux-xmHZnh5KvkYQVnfUL9eA4A4gTW7I',
    appId: '1:176535439905:ios:960f016c0ad72e0412f160',
    messagingSenderId: '176535439905',
    projectId: 'probability-app-6de7c',
    storageBucket: 'probability-app-6de7c.firebasestorage.app',
    iosBundleId: 'com.probabilityia.mobileCentral',
  );
}
