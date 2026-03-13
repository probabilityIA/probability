import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/dashboard_repository.dart';

class DashboardProvider extends ChangeNotifier {
  final ApiClient _apiClient;

  DashboardStats? _stats;
  bool _isLoading = false;
  String? _error;

  DashboardProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  DashboardStats? get stats => _stats;
  bool get isLoading => _isLoading;
  String? get error => _error;

  DashboardUseCases get _useCases =>
      DashboardUseCases(DashboardApiRepository(_apiClient));

  Future<void> fetchStats({int? businessId, int? integrationId}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final response = await _useCases.getStats(
        businessId: businessId,
        integrationId: integrationId,
      );
      _stats = response.data;
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    notifyListeners();
  }
}
