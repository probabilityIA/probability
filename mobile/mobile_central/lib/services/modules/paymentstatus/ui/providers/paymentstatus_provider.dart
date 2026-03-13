import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/paymentstatus_repository.dart';

class PaymentStatusProvider extends ChangeNotifier {
  final ApiClient _apiClient;
  List<PaymentStatusInfo> _statuses = [];
  bool _isLoading = false;
  String? _error;

  PaymentStatusProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  List<PaymentStatusInfo> get statuses => _statuses;
  bool get isLoading => _isLoading;
  String? get error => _error;

  PaymentStatusUseCases get _useCases => PaymentStatusUseCases(PaymentStatusApiRepository(_apiClient));

  Future<void> fetchStatuses({bool? isActive}) async {
    _isLoading = true; _error = null; notifyListeners();
    try { _statuses = await _useCases.getPaymentStatuses(isActive: isActive); } catch (e) { _error = e.toString(); }
    _isLoading = false; notifyListeners();
  }
}
