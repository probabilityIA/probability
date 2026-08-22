import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/dashboard_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class DashboardProvider extends ChangeNotifier {
  final ApiClient _apiClient;

  DashboardStats? _stats;
  bool _isLoading = false;
  String? _error;
  DashboardPeriod _period = DashboardPeriod.all;

  DashboardProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  DashboardStats? get stats => _stats;
  bool get isLoading => _isLoading;
  String? get error => _error;
  DashboardPeriod get period => _period;

  DashboardUseCases get _useCases =>
      DashboardUseCases(DashboardApiRepository(_apiClient));

  void setPeriod(DashboardPeriod period) {
    _period = period;
  }

  Future<void> fetchStats({
    int? businessId,
    int? integrationId,
    DateTime? now,
  }) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    final range = _period.rangeFrom(now ?? DateTime.now());

    try {
      final response = await _useCases.getStats(
        businessId: businessId,
        integrationId: integrationId,
        startDate: range?.start,
        endDate: range?.end,
      );
      _stats = response.data;
    } catch (e) {
      _error = parseError(e);
    }

    _isLoading = false;
    notifyListeners();
  }
}
