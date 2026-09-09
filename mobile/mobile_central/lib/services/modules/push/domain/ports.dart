import 'entities.dart';

abstract class IPushRepository {
  Future<void> registerDevice(DeviceRegistration registration);
  Future<void> unregisterDevice(String token);
}
