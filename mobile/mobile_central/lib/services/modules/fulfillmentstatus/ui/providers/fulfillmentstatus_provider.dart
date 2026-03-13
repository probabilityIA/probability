import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/fulfillmentstatus_repository.dart';

class FulfillmentStatusProvider extends ChangeNotifier {
  final ApiClient _apiClient;
  List<FulfillmentStatusInfo> _statuses = [];
  bool _isLoading = false;
  String? _error;

  FulfillmentStatusProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  List<FulfillmentStatusInfo> get statuses => _statuses;
  bool get isLoading => _isLoading;
  String? get error => _error;

  FulfillmentStatusUseCases get _useCases => FulfillmentStatusUseCases(FulfillmentStatusApiRepository(_apiClient));

  Future<void> fetchStatuses() async {
    _isLoading = true; _error = null; notifyListeners();
    try { _statuses = await _useCases.getFulfillmentStatuses(); } catch (e) { _error = e.toString(); }
    _isLoading = false; notifyListeners();
  }
}
