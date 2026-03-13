import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/notification_config_repository.dart';

class NotificationConfigProvider extends ChangeNotifier {
  final ApiClient _apiClient;

  List<NotificationConfig> _configs = [];
  bool _isLoading = false;
  String? _error;

  NotificationConfigProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  List<NotificationConfig> get configs => _configs;
  bool get isLoading => _isLoading;
  String? get error => _error;

  NotificationConfigUseCases get _useCases =>
      NotificationConfigUseCases(NotificationConfigApiRepository(_apiClient));

  Future<void> fetchConfigs({ConfigFilter? filter}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _configs = await _useCases.list(filter: filter);
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<NotificationConfig?> getById(int id, {int? businessId}) async {
    try {
      return await _useCases.getById(id, businessId: businessId);
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return null;
    }
  }

  Future<NotificationConfig?> createConfig(CreateConfigDTO dto, {int? businessId}) async {
    try {
      final config = await _useCases.create(dto, businessId: businessId);
      return config;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return null;
    }
  }

  Future<NotificationConfig?> updateConfig(int id, UpdateConfigDTO dto, {int? businessId}) async {
    try {
      final config = await _useCases.update(id, dto, businessId: businessId);
      return config;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return null;
    }
  }

  Future<bool> deleteConfig(int id, {int? businessId}) async {
    try {
      await _useCases.delete(id, businessId: businessId);
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  Future<SyncConfigsResponse?> syncByIntegration(SyncConfigsDTO dto, {int? businessId}) async {
    try {
      return await _useCases.syncByIntegration(dto, businessId: businessId);
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return null;
    }
  }
}
