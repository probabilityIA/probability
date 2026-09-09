import '../../../../../core/network/api_client.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class PushApiRepository implements IPushRepository {
  final ApiClient _client;

  PushApiRepository(this._client);

  @override
  Future<void> registerDevice(DeviceRegistration registration) async {
    await _client.post('/push/devices', data: registration.toJson());
  }

  @override
  Future<void> unregisterDevice(String token) async {
    await _client.delete('/push/devices', queryParameters: {'token': token});
  }
}
